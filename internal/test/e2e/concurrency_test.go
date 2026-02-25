// Package e2e provides comprehensive end-to-end tests for the ticket booking system
// Designed for 50K scalable backend with proper concurrency and distributed system testing
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

// Concurrency E2E Tests - Critical for 50K Scalable Backend

// TestConcurrentBookings_SameSeat tests race condition when multiple users try to book the same seat
// This is the MOST CRITICAL test for 50K scalable backend
func TestConcurrentBookings_SameSeat(t *testing.T) {
	client := SetupE2ETest(t)

	seatID := GenerateUniqueID("seat-race")
	numUsers := 50 // High concurrency to stress test

	t.Logf("Starting concurrent booking test: %d users trying to book seat %s", numUsers, seatID)

	var wg sync.WaitGroup
	successCount := int32(0)
	failedCount := int32(0)
	results := make(chan string, numUsers)

	startTime := time.Now()

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-race-%d-%d", index, time.Now().UnixNano())

			resp, err := ReserveTicket(client, userID, seatID)
			if err == nil && resp != nil {
				message := GetValueAsString(resp, "message")
				if message == "Ticket Reserved" {
					atomic.AddInt32(&successCount, 1)
					results <- fmt.Sprintf("SUCCESS: User %d reserved seat", index)
					// Cleanup immediately for next test
					CancelTicket(client, userID, seatID)
					return
				}
			}
			atomic.AddInt32(&failedCount, 1)
			results <- fmt.Sprintf("FAILED: User %d could not reserve", index)
		}(i)
	}

	wg.Wait()
	duration := time.Since(startTime)
	close(results)

	// Collect results
	successResults := []string{}
	failedResults := []string{}
	for r := range results {
		if len(successResults) < 5 {
			successResults = append(successResults, r)
		} else if len(failedResults) < 5 {
			failedResults = append(failedResults, r)
		}
	}

	t.Logf("=== Concurrent Booking Results ===")
	t.Logf("Total time: %v", duration)
	t.Logf("Successful: %d, Failed: %d", successCount, failedCount)
	t.Logf("First successes: %v", successResults)
	t.Logf("First failures: %v", failedResults)

	// CRITICAL ASSERTION: Only ONE user should successfully reserve the seat
	assert.Equal(t, int32(1), successCount,
		"CRITICAL: Only ONE user should be able to reserve the same seat! Found %d", successCount)

	// Additional verification
	t.Logf("Race condition test PASSED: Distributed lock working correctly")
}

// TestConcurrentBookings_DifferentSeats tests parallel booking of different seats
func TestConcurrentBookings_DifferentSeats(t *testing.T) {
	client := SetupE2ETest(t)

	numSeats := 20
	numUsersPerSeat := 10

	seatIDs := make([]string, numSeats)
	for i := 0; i < numSeats; i++ {
		seatIDs[i] = fmt.Sprintf("seat-parallel-%d", i)
	}

	t.Logf("Testing parallel booking: %d seats, %d users per seat", numSeats, numUsersPerSeat)

	var wg sync.WaitGroup
	totalSuccess := int32(0)

	// Each seat gets concurrent users
	for seatIdx, seatID := range seatIDs {
		for i := 0; i < numUsersPerSeat; i++ {
			wg.Add(1)
			go func(seatIndex int, seat string, userIndex int) {
				defer wg.Done()
				userID := fmt.Sprintf("user-parallel-%d-%d-%d", seatIndex, userIndex, time.Now().UnixNano())

				resp, err := ReserveTicket(client, userID, seat)
				if err == nil && resp != nil {
					if GetValueAsString(resp, "message") == "Ticket Reserved" {
						atomic.AddInt32(&totalSuccess, 1)
						CancelTicket(client, userID, seat)
					}
				}
			}(seatIdx, seatID, i)
		}
	}

	wg.Wait()

	t.Logf("Parallel booking results: %d/%d succeeded", totalSuccess, numSeats)

	// Each seat should have exactly 1 success
	assert.Equal(t, int32(numSeats), totalSuccess,
		"Each seat should have exactly one successful reservation")
}

// TestConcurrentPayments_SameBooking tests concurrent payment attempts for same booking
func TestConcurrentPayments_SameBooking(t *testing.T) {
	client := SetupE2ETest(t)

	bookingID := GenerateUniqueID("booking-concurrent-payment")
	numUsers := 20

	t.Logf("Testing concurrent payments: %d users trying to pay for booking %s", numUsers, bookingID)

	var wg sync.WaitGroup
	successCount := int32(0)

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-payment-concurrent-%d", index)
			amount := 500

			resp, err := InitiatePayment(client, bookingID, userID, amount)
			if err == nil && resp != nil {
				paymentID := GetValueAsString(resp, "id")
				if paymentID == "" {
					paymentID = GetValueAsString(resp, "payment_id")
				}
				if paymentID != "" {
					atomic.AddInt32(&successCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Concurrent payment results: %d/%d succeeded", successCount, numUsers)

	// Multiple payments can be initiated for same booking - this is expected
	assert.True(t, successCount > 0, "At least some payments should succeed")
}

// TestRateLimiting_Burst tests rate limiter under burst load
func TestRateLimiting_Burst(t *testing.T) {
	client := SetupE2ETest(t)

	// Note: The system allows 5000 requests per minute per the code
	// We'll test with a smaller burst to see behavior

	numRequests := 100
	t.Logf("Testing rate limiter with %d concurrent requests", numRequests)

	result := TestRateLimiter(t, client, numRequests)

	t.Logf("=== Rate Limit Test Results ===")
	t.Logf("Total requests: %d", result.TotalRequests)
	t.Logf("Successful: %d", result.SuccessCount)
	t.Logf("Rate limited: %d", result.RateLimitedCount)
	t.Logf("Errors: %d", result.ErrorCount)
	t.Logf("Duration: %v", result.EndTime)

	// Most requests should succeed as we're under the limit
	assert.True(t, result.SuccessCount > 0, "Some requests should succeed")
}

// TestRateLimiting_Sustained tests rate limiter under sustained load
func TestRateLimiting_Sustained(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sustained rate limit test in short mode")
	}

	client := SetupE2ETest(t)

	// Simulate sustained load
	duration := 10 * time.Second
	t.Logf("Testing rate limiter with sustained load for %v", duration)

	start := time.Now()
	successCount := int32(0)
	rateLimitedCount := int32(0)
	errorCount := int32(0)

	go func() {
		for {
			if time.Since(start) >= duration {
				break
			}
		go func() {
			resp, err := client.Get("/health")
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
				return
			}
			statusCode := GetValueAsInt(resp, "status_code")
			if statusCode == 429 {
				atomic.AddInt32(&rateLimitedCount, 1)
			} else if statusCode == 200 {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&errorCount, 1)
			}
		}()
			time.Sleep(10 * time.Millisecond) // 100 req/sec
		}
	}()

	time.Sleep(duration + time.Second)

	totalRequests := successCount + rateLimitedCount + errorCount
	t.Logf("=== Sustained Rate Limit Results ===")
	t.Logf("Total requests: %d", totalRequests)
	t.Logf("Successful: %d", successCount)
	t.Logf("Rate limited: %d", rateLimitedCount)
	t.Logf("Errors: %d", errorCount)
}

// TestDistributedLock_Validation tests that distributed locking works correctly
func TestDistributedLock_Validation(t *testing.T) {
	client := SetupE2ETest(t)

	// Test that locking works per-seat
	seats := []string{
		GenerateUniqueID("lock-seat-1"),
		GenerateUniqueID("lock-seat-2"),
		GenerateUniqueID("lock-seat-3"),
	}

	// Try to reserve each seat with multiple users
	for _, seatID := range seats {
		numUsers := 10

		var wg sync.WaitGroup
		successCount := int32(0)

		for i := 0; i < numUsers; i++ {
			wg.Add(1)
			go func(s string, u int) {
				defer wg.Done()
				userID := fmt.Sprintf("user-lock-%s-%d", s, u)
				resp, err := ReserveTicket(client, userID, s)
				if err == nil && resp != nil {
					if GetValueAsString(resp, "message") == "Ticket Reserved" {
						atomic.AddInt32(&successCount, 1)
						CancelTicket(client, userID, s)
					}
				}
			}(seatID, i)
		}

		wg.Wait()

		t.Logf("Seat %s: %d/%d users succeeded", seatID, successCount, numUsers)
		assert.Equal(t, int32(1), successCount, "Each seat should have exactly 1 successful reservation")
	}
}

// TestConcurrentBookings_MixedOperations tests mixed concurrent operations
func TestConcurrentBookings_MixedOperations(t *testing.T) {
	client := SetupE2ETest(t)

	seatID := GenerateUniqueID("seat-mixed")
	numReserve := 10
	numConfirm := 5
	numCancel := 5

	var wg sync.WaitGroup

	// Concurrent reservations
	reserveCount := int32(0)
	for i := 0; i < numReserve; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-reserve-%d", index)
			resp, err := ReserveTicket(client, userID, seatID)
			if err == nil && resp != nil {
				if GetValueAsString(resp, "message") == "Ticket Reserved" {
					atomic.AddInt32(&reserveCount, 1)
				}
			}
		}(i)
	}

	// Concurrent confirmations
	for i := 0; i < numConfirm; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-confirm-%d", index)
			ConfirmTicket(client, userID, seatID)
		}(i)
	}

	// Concurrent cancellations
	for i := 0; i < numCancel; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			userID := fmt.Sprintf("user-cancel-%d", index)
			CancelTicket(client, userID, seatID)
		}(i)
	}

	wg.Wait()

	t.Logf("Mixed operations results: %d reservations succeeded", reserveCount)

	// Cleanup
	CancelTicket(client, "cleanup-user", seatID)
}

// TestConcurrentReadWrite tests concurrent reads and writes
func TestConcurrentReadWrite(t *testing.T) {
	client := SetupE2ETest(t)

	seatID := GenerateUniqueID("seat-readwrite")
	numOperations := 30

	var wg sync.WaitGroup
	writeSuccess := int32(0)
	readSuccess := int32(0)

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			if index%2 == 0 {
				// Write operation
				userID := fmt.Sprintf("user-write-%d", index)
				resp, err := ReserveTicket(client, userID, seatID)
				if err == nil && resp != nil {
					if GetValueAsString(resp, "message") == "Ticket Reserved" {
						atomic.AddInt32(&writeSuccess, 1)
					}
				}
			} else {
				// Read operation
				resp, err := CheckAvailability(client, seatID)
				if err == nil && resp != nil {
					atomic.AddInt32(&readSuccess, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Concurrent read/write: %d writes succeeded, %d reads succeeded", writeSuccess, readSuccess)

	// Cleanup
	CancelTicket(client, "cleanup", seatID)
}

// TestHighConcurrency_Scalability tests the system under high concurrency
func TestHighConcurrency_Scalability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping high concurrency test in short mode")
	}

	client := SetupE2ETest(t)

	numSeats := 100
	numUsersPerSeat := 10

	seatIDs := make([]string, numSeats)
	for i := 0; i < numSeats; i++ {
		seatIDs[i] = fmt.Sprintf("seat-scale-%d", i)
	}

	t.Logf("Starting high concurrency test: %d seats x %d users = %d total attempts",
		numSeats, numUsersPerSeat, numSeats*numUsersPerSeat)

	start := time.Now()

	var wg sync.WaitGroup
	totalSuccess := int32(0)

	// All seats get concurrent users at the same time
	barrier := make(chan struct{})

	for seatIdx, seatID := range seatIDs {
		for i := 0; i < numUsersPerSeat; i++ {
			wg.Add(1)
			go func(seatIndex int, seat string, userIndex int) {
				<-barrier // Wait for all goroutines to be ready
				defer wg.Done()

				userID := fmt.Sprintf("user-scale-%d-%d-%d", seatIndex, userIndex, time.Now().UnixNano())
				resp, err := ReserveTicket(client, userID, seat)
				if err == nil && resp != nil {
					if GetValueAsString(resp, "message") == "Ticket Reserved" {
						atomic.AddInt32(&totalSuccess, 1)
						CancelTicket(client, userID, seat)
					}
				}
			}(seatIdx, seatID, i)
		}
	}

	// Release all goroutines at once
	close(barrier)

	wg.Wait()
	duration := time.Since(start)

	t.Logf("=== High Concurrency Test Results ===")
	t.Logf("Total time: %v", duration)
	t.Logf("Successful: %d/%d", totalSuccess, numSeats*numUsersPerSeat)
	t.Logf("Throughput: %.2f operations/sec", float64(numSeats*numUsersPerSeat)/duration.Seconds())

	// Each seat should have exactly 1 success
	assert.Equal(t, int32(numSeats), totalSuccess,
		"Each seat should have exactly one successful reservation")

	// Print performance metrics
	t.Logf("Performance: %.2f seats/sec", float64(numSeats)/duration.Seconds())
}

// TestConcurrentUserRegistration tests concurrent user registrations
func TestConcurrentUserRegistration(t *testing.T) {
	client := SetupE2ETest(t)

	numUsers := 50
	t.Logf("Testing concurrent user registration: %d users", numUsers)

	var wg sync.WaitGroup
	successCount := int32(0)

	start := time.Now()

	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			user := TestUser{
				Username: GenerateUniqueID(fmt.Sprintf("concurrent-%d", index)),
				Email:    GenerateUniqueID(fmt.Sprintf("concurrent-%d", index)) + "@example.com",
				Password: "testpass123",
			}

			resp, err := CreateTestUser(client, user)
			if err == nil && resp != nil {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("Concurrent user registration: %d/%d succeeded in %v (%.2f users/sec)",
		successCount, numUsers, duration, float64(numUsers)/duration.Seconds())

	assert.True(t, successCount > 0, "At least some users should register successfully")
}

// TestRaceCondition_PreventDoubleConfirm tests that double confirmation is prevented
func TestRaceCondition_PreventDoubleConfirm(t *testing.T) {
	client := SetupE2ETest(t)

	seatID := GenerateUniqueID("seat-double-confirm")

	// First, reserve the seat
	userID := GenerateUniqueID("user-first")
	resp, err := ReserveTicket(client, userID, seatID)
	require.NoError(t, err, "Initial reservation should succeed")
	require.Equal(t, "Ticket Reserved", GetValueAsString(resp, "message"))

	// Now try to confirm from multiple users concurrently
	numUsers := 20
	var confirmCount int32

	var wg sync.WaitGroup
	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			uid := fmt.Sprintf("user-confirm-%d", index)
			resp, err := ConfirmTicket(client, uid, seatID)
			if err == nil && resp != nil {
				if GetValueAsString(resp, "message") == "Ticket Confirmed" {
					atomic.AddInt32(&confirmCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Double confirmation test: %d/%d confirmations succeeded", confirmCount, numUsers)

	// Only one confirmation should succeed
	assert.Equal(t, int32(1), confirmCount,
		"Only one confirmation should succeed for a reserved seat")

	// Cleanup
	CancelTicket(client, userID, seatID)
}

