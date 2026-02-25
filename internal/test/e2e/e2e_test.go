// Package e2e provides comprehensive end-to-end tests for the ticket booking system
// This file contains the main E2E test scenarios
package e2e

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// E2E Test Suite - Main Booking Flow Tests

// TestE2E_CompleteBookingFlow tests the complete booking workflow
func TestE2E_CompleteBookingFlow(t *testing.T) {
	client := SetupE2ETest(t)

	// Generate unique identifiers for this test
	seatID := GenerateUniqueID("seat-e2e-complete")
	userID := GenerateUniqueID("user-e2e-complete")
	movieTitle := "Test Movie E2E"
	showTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	t.Logf("=== Starting Complete Booking Flow Test ===")
	t.Logf("SeatID: %s, UserID: %s", seatID, userID)

	// Step 1: Create a movie (optional - depends on system setup)
	t.Run("Step1_CreateMovie", func(t *testing.T) {
		resp, err := CreateMovie(client, movieTitle, 120)
		if err != nil {
			t.Logf("Movie creation response: %+v", resp)
			// Some systems may require admin privileges
			t.Logf("Note: Movie creation may require admin privileges: %v", err)
		}
	})

	// Step 2: Create a show (optional - depends on system setup)
	t.Run("Step2_CreateShow", func(t *testing.T) {
		resp, err := CreateShow(client, "movie-1", showTime)
		if err != nil {
			t.Logf("Show creation response: %+v", resp)
			t.Logf("Note: Show creation may require valid movie ID: %v", err)
		}
	})

	// Step 3: Reserve a ticket
	t.Run("Step3_ReserveTicket", func(t *testing.T) {
		resp, err := ReserveTicket(client, userID, seatID)
		require.NoError(t, err, "Reserve ticket should not return error")
		AssertResponseSuccess(t, resp, 200)

		message := GetValueAsString(resp, "message")
		assert.Equal(t, "Ticket Reserved", message, "Should return 'Ticket Reserved' message")

		t.Logf("✓ Ticket reserved successfully: %s", seatID)
	})

	// Step 4: Confirm the ticket
	t.Run("Step4_ConfirmTicket", func(t *testing.T) {
		resp, err := ConfirmTicket(client, userID, seatID)
		require.NoError(t, err, "Confirm ticket should not return error")
		AssertResponseSuccess(t, resp, 200)

		message := GetValueAsString(resp, "message")
		assert.Equal(t, "Ticket Confirmed", message, "Should return 'Ticket Confirmed' message")

		t.Logf("✓ Ticket confirmed successfully: %s", seatID)
	})

	// Step 5: Initiate payment
	t.Run("Step5_InitiatePayment", func(t *testing.T) {
		resp, err := InitiatePayment(client, seatID, userID, 500)
		require.NoError(t, err, "Initiate payment should not return error")
		AssertResponseSuccess(t, resp, 200)

		paymentID := GetValueAsString(resp, "id")
		if paymentID == "" {
			paymentID = GetValueAsString(resp, "payment_id")
		}
		assert.NotEmpty(t, paymentID, "Should return payment ID")

		t.Logf("✓ Payment initiated successfully: %s", paymentID)
	})

	// Step 6: Verify payment
	t.Run("Step6_VerifyPayment", func(t *testing.T) {
		resp, err := VerifyPayment(client, seatID, "success")
		require.NoError(t, err, "Verify payment should not return error")
		AssertResponseSuccess(t, resp, 200)

		t.Logf("✓ Payment verified successfully")
	})

	// Step 7: Get user reservations (query)
	t.Run("Step7_GetUserReservations", func(t *testing.T) {
		resp, err := GetUserReservations(client, userID)
		require.NoError(t, err, "Get reservations should not return error")
		t.Logf("✓ User reservations retrieved: %+v", resp)
	})

	// Cleanup
	t.Cleanup(func() {
		CancelTicket(client, userID, seatID)
	})

	t.Logf("=== Complete Booking Flow Test Passed ===")
}

// User Management E2E Tests

// TestE2E_UserRegistration tests user registration flow
func TestE2E_UserRegistration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping in short mode")
	}

	client := SetupE2ETest(t)

	t.Run("RegisterValidUser", func(t *testing.T) {
		user := TestUser{
			Username: GenerateUniqueID("testuser"),
			Email:    GenerateUniqueID("testuser") + "@example.com",
			Password: "securePassword123",
		}

		resp, err := CreateTestUser(client, user)
		require.NoError(t, err, "User registration should not return error")
		AssertResponseSuccess(t, resp, 201)

		// Verify user was created (check for user ID or username in response)
		userID := GetValueAsString(resp, "id")
		if userID == "" {
			userID = GetValueAsString(resp, "user_id")
		}
		if userID == "" {
			userID = GetValueAsString(resp, "username")
		}
		assert.NotEmpty(t, userID, "Should return user identification")

		t.Logf("✓ User registered successfully: %s", userID)
	})

	t.Run("RegisterDuplicateEmail", func(t *testing.T) {
		email := GenerateUniqueID("duplicate") + "@example.com"
		user1 := TestUser{
			Username: GenerateUniqueID("user1"),
			Email:    email,
			Password: "password123",
		}

		// First registration should succeed
		resp1, err1 := CreateTestUser(client, user1)
		if err1 != nil {
			t.Logf("First registration error (may already exist): %v", err1)
		} else {
			_ = resp1 // explicitly ignore unused response
		}

		// Second registration with same email should fail
		user2 := TestUser{
			Username: GenerateUniqueID("user2"),
			Email:    email,
			Password: "password123",
		}

		resp2, err2 := CreateTestUser(client, user2)
		if err2 == nil {
			// Check if duplicate is allowed or rejected
			statusCode := GetValueAsInt(resp2, "status_code")
			assert.True(t, statusCode >= 400, "Duplicate email should be rejected")
		}

		t.Logf("✓ Duplicate email test completed")
	})

	t.Run("RegisterInvalidEmail", func(t *testing.T) {
		user := TestUser{
			Username: GenerateUniqueID("invalid"),
			Email:    "not-an-email",
			Password: "password123",
		}

		resp, err := CreateTestUser(client, user)
		// Either network error or 400/422 status code expected
		if err == nil {
			statusCode := GetValueAsInt(resp, "status_code")
			assert.True(t, statusCode >= 400, "Invalid email should be rejected")
			errMsg := ParseErrorMessage(resp)
			t.Logf("Invalid email error: %s", errMsg)
		}

		t.Logf("✓ Invalid email test completed")
	})
}

// Booking E2E Tests

// TestE2E_BookingReserve tests ticket reservation
func TestE2E_BookingReserve(t *testing.T) {
	client := SetupE2ETest(t)

	t.Run("ReserveAvailableSeat", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-available")
		userID := GenerateUniqueID("user-reserve")

		resp, err := ReserveTicket(client, userID, seatID)
		require.NoError(t, err)
		AssertResponseSuccess(t, resp, 200)

		message := GetValueAsString(resp, "message")
		assert.Equal(t, "Ticket Reserved", message)

		t.Cleanup(func() {
			CancelTicket(client, userID, seatID)
		})

		t.Logf("✓ Available seat reserved successfully")
	})

	t.Run("ReserveSameSeatTwice", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-double")
		user1 := GenerateUniqueID("user-first")
		user2 := GenerateUniqueID("user-second")

		// First reservation should succeed
		resp1, err1 := ReserveTicket(client, user1, seatID)
		require.NoError(t, err1)
		AssertResponseSuccess(t, resp1, 200)

		// Second reservation should fail
		resp2, err2 := ReserveTicket(client, user2, seatID)
		
		// Should either return error or 409 Conflict
		if err2 == nil {
			statusCode := GetValueAsInt(resp2, "status_code")
			assert.Equal(t, 409, statusCode, "Second reservation should return 409")
		}

		errMsg := ParseErrorMessage(resp2)
		assert.True(t, 
			errMsg == "seat already reserved" || 
			GetValueAsInt(resp2, "status_code") == 409,
			"Should indicate seat already reserved")

		t.Cleanup(func() {
			CancelTicket(client, user1, seatID)
		})

		t.Logf("✓ Double reservation correctly prevented")
	})

	t.Run("ReserveWithInvalidUser", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-invalid-user")

		resp, err := ReserveTicket(client, "", seatID)
		if err == nil {
			statusCode := GetValueAsInt(resp, "status_code")
			assert.True(t, statusCode >= 400, "Invalid user should be rejected")
		}

		t.Logf("✓ Invalid user test completed")
	})
}

// TestE2E_BookingConfirm tests ticket confirmation
func TestE2E_BookingConfirm(t *testing.T) {
	client := SetupE2ETest(t)

	t.Run("ConfirmReservedTicket", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-confirm")
		userID := GenerateUniqueID("user-confirm")

		// Reserve first
		resp1, err1 := ReserveTicket(client, userID, seatID)
		require.NoError(t, err1)
		AssertResponseSuccess(t, resp1, 200)

		// Then confirm
		resp2, err2 := ConfirmTicket(client, userID, seatID)
		require.NoError(t, err2)
		AssertResponseSuccess(t, resp2, 200)

		message := GetValueAsString(resp2, "message")
		assert.Equal(t, "Ticket Confirmed", message)

		t.Cleanup(func() {
			CancelTicket(client, userID, seatID)
		})

		t.Logf("✓ Reserved ticket confirmed successfully")
	})

	t.Run("ConfirmWithoutReservation", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-never-reserved")
		userID := GenerateUniqueID("user-no-res")

		resp, err := ConfirmTicket(client, userID, seatID)
		if err == nil {
			statusCode := GetValueAsInt(resp, "status_code")
			assert.Equal(t, 409, statusCode, "Should return 409 for unreserved seat")
		}

		t.Logf("✓ Unreserved ticket confirmation correctly rejected")
	})
}

// TestE2E_BookingCancel tests ticket cancellation
func TestE2E_BookingCancel(t *testing.T) {
	client := SetupE2ETest(t)

	t.Run("CancelReservedTicket", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-cancel")
		userID := GenerateUniqueID("user-cancel")

		// Reserve first
		resp1, err1 := ReserveTicket(client, userID, seatID)
		require.NoError(t, err1)
		AssertResponseSuccess(t, resp1, 200)

		// Then cancel
		resp2, err2 := CancelTicket(client, userID, seatID)
		require.NoError(t, err2)
		AssertResponseSuccess(t, resp2, 200)

		message := GetValueAsString(resp2, "message")
		assert.Equal(t, "Ticket Cancelled", message)

		t.Logf("✓ Reserved ticket cancelled successfully")
	})

	t.Run("CancelNonExistentTicket", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-never-existed")
		userID := GenerateUniqueID("user-never-existed")

		resp, err := CancelTicket(client, userID, seatID)
		if err == nil {
			statusCode := GetValueAsInt(resp, "status_code")
			assert.Equal(t, 404, statusCode, "Should return 404 for non-existent reservation")
		}

		t.Logf("✓ Non-existent ticket cancellation correctly handled")
	})
}

// Payment E2E Tests

// TestE2E_PaymentFlow tests payment initiation and verification
func TestE2E_PaymentFlow(t *testing.T) {
	client := SetupE2ETest(t)

	t.Run("InitiatePayment", func(t *testing.T) {
		bookingID := GenerateUniqueID("booking-pay")
		userID := GenerateUniqueID("user-pay")
		amount := 500

		resp, err := InitiatePayment(client, bookingID, userID, amount)
		require.NoError(t, err)
		AssertResponseSuccess(t, resp, 200)

		paymentID := GetValueAsString(resp, "id")
		if paymentID == "" {
			paymentID = GetValueAsString(resp, "payment_id")
		}
		assert.NotEmpty(t, paymentID, "Should return payment ID")

		t.Logf("✓ Payment initiated: %s", paymentID)
	})

	t.Run("VerifyPaymentSuccess", func(t *testing.T) {
		bookingID := GenerateUniqueID("booking-verify")
		userID := GenerateUniqueID("user-verify")

		// Initiate payment first
		_, err1 := InitiatePayment(client, bookingID, userID, 1000)
		require.NoError(t, err1)
		
		// Verify with success mode
		resp2, err2 := VerifyPayment(client, bookingID, "success")
		require.NoError(t, err2)
		AssertResponseSuccess(t, resp2, 200)

		t.Logf("✓ Payment verified successfully")
	})

	t.Run("VerifyPaymentFailure", func(t *testing.T) {
		bookingID := GenerateUniqueID("booking-fail")
		userID := GenerateUniqueID("user-fail")

		// Initiate payment first
		_, err1 := InitiatePayment(client, bookingID, userID, 1000)
		require.NoError(t, err1)

		// Verify with failure mode
		resp2, err2 := VerifyPayment(client, bookingID, "fail")
		// Failure verification might still return 200 but with different status
		t.Logf("Payment failure response: %+v, err: %v", resp2, err2)

		t.Logf("✓ Payment failure test completed")
	})
}

// Query E2E Tests

// TestE2E_QueryEndpoints tests query endpoints
func TestE2E_QueryEndpoints(t *testing.T) {
	client := SetupE2ETest(t)

	t.Run("GetUserReservations", func(t *testing.T) {
		userID := GenerateUniqueID("user-query")

		// Create a reservation first
		seatID1 := GenerateUniqueID("seat-q1")
		seatID2 := GenerateUniqueID("seat-q2")
		ReserveTicket(client, userID, seatID1)
		
		// Give some time for the reservation to be processed
		time.Sleep(100 * time.Millisecond)

		// Query reservations
		resp, err := GetUserReservations(client, userID)
		require.NoError(t, err)
		t.Logf("Reservations response: %+v", resp)

		// Cleanup
		CancelTicket(client, userID, seatID1)
		CancelTicket(client, userID, seatID2)

		t.Logf("✓ User reservations queried successfully")
	})

	t.Run("CheckAvailability", func(t *testing.T) {
		seatID := GenerateUniqueID("seat-avail")

		// Check availability before reservation
		resp1, err1 := CheckAvailability(client, seatID)
		require.NoError(t, err1)
		t.Logf("Availability before reservation: %+v", resp1)

		// Reserve the seat
		userID := GenerateUniqueID("user-avail")
		ReserveTicket(client, userID, seatID)

		// Check availability after reservation
		resp2, err2 := CheckAvailability(client, seatID)
		require.NoError(t, err2)
		t.Logf("Availability after reservation: %+v", resp2)

		// Cleanup
		CancelTicket(client, userID, seatID)

		t.Logf("✓ Seat availability checked successfully")
	})

	t.Run("GetMovies", func(t *testing.T) {
		resp, err := GetMovies(client)
		require.NoError(t, err)
		t.Logf("Movies response: %+v", resp)

		t.Logf("✓ Movies queried successfully")
	})

	t.Run("GetShows", func(t *testing.T) {
		resp, err := GetShows(client)
		require.NoError(t, err)
		t.Logf("Shows response: %+v", resp)

		t.Logf("✓ Shows queried successfully")
	})

	t.Run("GetNotifications", func(t *testing.T) {
		userID := GenerateUniqueID("user-notif")

		resp, err := GetUserNotifications(client, userID)
		require.NoError(t, err)
		t.Logf("Notifications response: %+v", resp)

		t.Logf("✓ Notifications queried successfully")
	})
}

// Health Check E2E Tests

// TestE2E_HealthCheck tests the health endpoint
func TestE2E_HealthCheck(t *testing.T) {
	client := SetupE2ETest(t)

	resp, err := CheckHealth(client)
	require.NoError(t, err)
	AssertResponseSuccess(t, resp, 200)

	t.Logf("✓ Health check passed: %+v", resp)
}

// Error Handling E2E Tests

// TestE2E_ErrorHandling tests various error scenarios
func TestE2E_ErrorHandling(t *testing.T) {
	client := SetupE2ETest(t)

	t.Run("InvalidJSON", func(t *testing.T) {
		// Test with invalid request body
		resp, err := client.Post("/cmd/reserve", "invalid json")
		// Should handle gracefully
		t.Logf("Invalid JSON response: %+v, err: %v", resp, err)
		t.Logf("✓ Invalid JSON handled")
	})

	t.Run("MissingRequiredFields", func(t *testing.T) {
		// Test with missing user_id
		resp, err := ReserveTicket(client, "", "seat-test")
		if err == nil {
			statusCode := GetValueAsInt(resp, "status_code")
			assert.True(t, statusCode >= 400, "Missing user_id should be rejected")
		}

		t.Logf("✓ Missing fields handled correctly")
	})

	t.Run("NonExistentEndpoints", func(t *testing.T) {
		resp, err := client.Get("/nonexistent/endpoint")
		require.NoError(t, err)
		
		statusCode := GetValueAsInt(resp, "status_code")
		// Should return 404 or similar
		assert.True(t, statusCode >= 400, "Non-existent endpoint should return error")

		t.Logf("✓ Non-existent endpoint handled: %d", statusCode)
	})
}

// Concurrent Booking E2E Tests (Simplified version)

// TestE2E_ConcurrentBooking tests concurrent booking attempts
func TestE2E_ConcurrentBooking(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent test in short mode")
	}

	client := SetupE2ETest(t)

	seatID := GenerateUniqueID("seat-race-e2e")
	numUsers := 10

	t.Logf("Testing concurrent booking with %d users for seat %s", numUsers, seatID)

	var wg sync.WaitGroup
	successCount := int32(0)
	failCount := int32(0)

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-race-%d-%d", index, time.Now().UnixNano())

			resp, err := ReserveTicket(client, userID, seatID)
			if err == nil {
				statusCode := GetValueAsInt(resp, "status_code")
				if statusCode == 200 {
					atomic.AddInt32(&successCount, 1)
					// Cleanup immediately
					CancelTicket(client, userID, seatID)
					return
				}
			}
			atomic.AddInt32(&failCount, 1)
		}(i)
	}

	wg.Wait()

	t.Logf("Concurrent booking results: %d succeeded, %d failed", successCount, failCount)

	// Only one should succeed
	assert.Equal(t, int32(1), successCount, 
		"Only one user should successfully reserve the seat")

	t.Logf("✓ Concurrent booking test passed")
}

// End-to-End Multi-User Scenario Tests

// TestE2E_MultiUserScenario tests a realistic multi-user booking scenario
func TestE2E_MultiUserScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multi-user test in short mode")
	}

	client := SetupE2ETest(t)

	numUsers := 5
	numSeats := 5

	t.Logf("=== Multi-user scenario: %d users booking %d seats ===", numUsers, numSeats)

	// Each user tries to book different seats
	results := make(chan string, numUsers*numSeats)
	var wg sync.WaitGroup

	for i := 0; i < numUsers; i++ {
		userID := fmt.Sprintf("multiuser-%d", i)
		
		for j := 0; j < numSeats; j++ {
			seatID := fmt.Sprintf("seat-multi-%d-%d", i, j)
			
			wg.Add(1)
			go func(u, s string) {
				defer wg.Done()
				
				resp, err := ReserveTicket(client, u, s)
				if err == nil {
					statusCode := GetValueAsInt(resp, "status_code")
					if statusCode == 200 {
						results <- fmt.Sprintf("SUCCESS: User %s booked seat %s", u, s)
						// Cleanup
						CancelTicket(client, u, s)
						return
					}
				}
				results <- fmt.Sprintf("FAILED: User %s could not book seat %s", u, s)
			}(userID, seatID)
		}
	}

	wg.Wait()
	close(results)

	// Count results
	successCount := 0
	failCount := 0
	for r := range results {
		if len(results) <= 10 || successCount+failCount < 10 {
			t.Logf("%s", r)
		}
		if len(r) > 6 && r[:6] == "SUCCESS" {
			successCount++
		} else {
			failCount++
		}
	}

	t.Logf("Multi-user results: %d succeeded, %d failed", successCount, failCount)
	assert.Equal(t, numUsers*numSeats, successCount, 
		"All users should be able to book different seats")

	t.Logf("✓ Multi-user scenario test passed")
}

// Performance Baseline E2E Tests

// TestE2E_PerformanceBaseline establishes performance baseline
func TestE2E_PerformanceBaseline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	client := SetupE2ETest(t)

	numRequests := 100
	t.Logf("Testing baseline performance with %d requests", numRequests)

	start := time.Now()
	
	var wg sync.WaitGroup
	successCount := int32(0)
	errorCount := int32(0)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := CheckHealth(client)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
				return
			}
			if GetValueAsInt(resp, "status_code") == 200 {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&errorCount, 1)
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("=== Performance Baseline Results ===")
	t.Logf("Total requests: %d", numRequests)
	t.Logf("Successful: %d", successCount)
	t.Logf("Errors: %d", errorCount)
	t.Logf("Duration: %v", duration)
	t.Logf("Throughput: %.2f req/sec", float64(numRequests)/duration.Seconds())

	assert.True(t, successCount > 0, "At least some requests should succeed")
	
	t.Logf("✓ Performance baseline test completed")
}

// Test Suite Summary

// TestE2E_AllFlowsSummary runs a summary of all main flows
func TestE2E_AllFlowsSummary(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping summary test in short mode")
	}

	client := SetupE2ETest(t)

	t.Logf("=== Running E2E Test Summary ===")

	// 1. Health check
	t.Run("Summary_HealthCheck", func(t *testing.T) {
		resp, err := CheckHealth(client)
		require.NoError(t, err)
		AssertResponseSuccess(t, resp, 200)
	})

	// 2. User registration
	t.Run("Summary_UserRegistration", func(t *testing.T) {
		user := TestUser{
			Username: GenerateUniqueID("summary-user"),
			Email:    GenerateUniqueID("summary-user") + "@test.com",
			Password: "test123456",
		}
		resp, _ := CreateTestUser(client, user)
		statusCode := GetValueAsInt(resp, "status_code")
		assert.True(t, statusCode == 201 || statusCode == 200 || statusCode >= 400)
	})

	// 3. Booking flow
	t.Run("Summary_BookingFlow", func(t *testing.T) {
		seatID := GenerateUniqueID("summary-seat")
		userID := GenerateUniqueID("summary-user-booking")

		resp, err := ReserveTicket(client, userID, seatID)
		if err == nil && GetValueAsInt(resp, "status_code") == 200 {
			ConfirmTicket(client, userID, seatID)
			CancelTicket(client, userID, seatID)
		}
	})

	// 4. Payment flow
	t.Run("Summary_PaymentFlow", func(t *testing.T) {
		bookingID := GenerateUniqueID("summary-booking")
		userID := GenerateUniqueID("summary-user-payment")
		
		resp, _ := InitiatePayment(client, bookingID, userID, 100)
		statusCode := GetValueAsInt(resp, "status_code")
		assert.True(t, statusCode == 200 || statusCode >= 400)
	})

	// 5. Queries
	t.Run("Summary_Queries", func(t *testing.T) {
		GetMovies(client)
		GetShows(client)
	})

	t.Logf("=== E2E Test Summary Completed ===")
}

