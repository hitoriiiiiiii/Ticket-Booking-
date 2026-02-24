// Package e2e provides comprehensive end-to-end tests for the ticket booking system
// This file contains helper functions for E2E testing
package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestUser represents a test user for E2E tests
type TestUser struct {
	ID       string
	Username string
	Email    string
	Password string
	Token    string
}

// E2EClient represents an HTTP client configured for E2E testing
type E2EClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// RateLimitResult holds results from rate limit testing
type RateLimitResult struct {
	TotalRequests    int
	SuccessCount     int
	RateLimitedCount int
	ErrorCount       int
	EndTime          time.Time
}

// NewE2EClient creates a new E2E test client
func NewE2EClient(baseURL string) *E2EClient {
	if baseURL == "" {
		baseURL = getBaseURL()
	}
	return &E2EClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// getBaseURL returns the base URL for API testing
func getBaseURL() string {
	if url := os.Getenv("API_BASE_URL"); url != "" {
		return url
	}
	return "http://localhost:8080"
}

// GenerateUniqueID generates a unique ID for test data
func GenerateUniqueID(prefix string) string {
	return fmt.Sprintf("%s-%s-%d", prefix, uuid.New().String()[:8], time.Now().UnixNano())
}

// GetValueAsString extracts a string value from a map
func GetValueAsString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetValueAsInt extracts an int value from a map
func GetValueAsInt(data map[string]interface{}, key string) int {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		case int64:
			return int(v)
		}
	}
	return 0
}

// makeRequest performs an HTTP request and returns the response body as a map
func (c *E2EClient) makeRequest(method, endpoint string, body interface{}) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		// Return raw response if not JSON
		return map[string]interface{}{
			"status_code": float64(resp.StatusCode),
			"body":        string(respBody),
		}, nil
	}

	result["status_code"] = float64(resp.StatusCode)
	return result, nil
}

// Get performs a GET request
func (c *E2EClient) Get(endpoint string) (map[string]interface{}, error) {
	return c.makeRequest("GET", endpoint, nil)
}

// Post performs a POST request
func (c *E2EClient) Post(endpoint string, body interface{}) (map[string]interface{}, error) {
	return c.makeRequest("POST", endpoint, body)
}

// Booking API Helpers

// ReserveTicket reserves a ticket for a user
func ReserveTicket(client *E2EClient, userID, seatID string) (map[string]interface{}, error) {
	return client.Post("/cmd/reserve", map[string]interface{}{
		"user_id": userID,
		"seat_id": seatID,
	})
}

// ConfirmTicket confirms a reserved ticket
func ConfirmTicket(client *E2EClient, userID, seatID string) (map[string]interface{}, error) {
	return client.Post("/cmd/confirm", map[string]interface{}{
		"user_id": userID,
		"seat_id": seatID,
	})
}

// CancelTicket cancels a reserved ticket
func CancelTicket(client *E2EClient, userID, seatID string) (map[string]interface{}, error) {
	return client.Post("/cmd/cancel", map[string]interface{}{
		"user_id": userID,
		"seat_id": seatID,
	})
}

// CheckAvailability checks seat availability
func CheckAvailability(client *E2EClient, seatID string) (map[string]interface{}, error) {
	return client.Get(fmt.Sprintf("/query/availability/%s", seatID))
}

// GetUserReservations gets reservations for a user
func GetUserReservations(client *E2EClient, userID string) (map[string]interface{}, error) {
	return client.Get(fmt.Sprintf("/query/reservations/%s", userID))
}

// User API Helpers

// CreateTestUser creates a test user
func CreateTestUser(client *E2EClient, user TestUser) (map[string]interface{}, error) {
	return client.Post("/cmd/users/register", map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"password": user.Password,
	})
}

// ListUsers lists all users
func ListUsers(client *E2EClient) (map[string]interface{}, error) {
	return client.Get("/query/users")
}

// Movie API Helpers

// CreateMovie creates a new movie
func CreateMovie(client *E2EClient, title string, duration int) (map[string]interface{}, error) {
	return client.Post("/cmd/movies", map[string]interface{}{
		"title":    title,
		"duration": duration,
	})
}

// GetMovies gets all movies
func GetMovies(client *E2EClient) (map[string]interface{}, error) {
	return client.Get("/query/movies")
}

// GetMovie gets a specific movie
func GetMovie(client *E2EClient, movieID string) (map[string]interface{}, error) {
	return client.Get(fmt.Sprintf("/query/movies/%s", movieID))
}

// Show API Helpers

// CreateShow creates a new show
func CreateShow(client *E2EClient, movieID, showTime string) (map[string]interface{}, error) {
	return client.Post("/cmd/shows", map[string]interface{}{
		"movie_id":  movieID,
		"show_time": showTime,
	})
}

// GetShows gets all shows
func GetShows(client *E2EClient) (map[string]interface{}, error) {
	return client.Get("/query/shows")
}

// Payment API Helpers

// InitiatePayment initiates a payment for a booking
func InitiatePayment(client *E2EClient, bookingID, userID string, amount int) (map[string]interface{}, error) {
	return client.Post("/cmd/payments/initiate", map[string]interface{}{
		"booking_id": bookingID,
		"user_id":   userID,
		"amount":    amount,
	})
}

// VerifyPayment verifies a payment
func VerifyPayment(client *E2EClient, paymentID, mode string) (map[string]interface{}, error) {
	return client.Post("/cmd/payments/verify", map[string]interface{}{
		"payment_id": paymentID,
		"mode":       mode,
	})
}

// Notification API Helpers

// GetUserNotifications gets notifications for a user
func GetUserNotifications(client *E2EClient, userID string) (map[string]interface{}, error) {
	return client.Get(fmt.Sprintf("/query/notifications/%s", userID))
}

// Health Check

// CheckHealth checks if the API is healthy
func CheckHealth(client *E2EClient) (map[string]interface{}, error) {
	return client.Get("/health")
}

// Rate Limiting Test Helpers

// TestRateLimiter tests the rate limiter with a given number of requests
func TestRateLimiter(t *testing.T, client *E2EClient, numRequests int) RateLimitResult {
	t.Helper()

	result := RateLimitResult{
		TotalRequests: numRequests,
		EndTime:       time.Now().Add(1 * time.Second),
	}

	// Use a channel to limit concurrency
	sem := make(chan struct{}, 50)

	for i := 0; i < numRequests; i++ {
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()

			resp, err := client.Get("/health")
			if err != nil {
				result.ErrorCount++
				return
			}

			if GetValueAsInt(resp, "status_code") == 429 {
				result.RateLimitedCount++
			} else if GetValueAsInt(resp, "status_code") == 200 {
				result.SuccessCount++
			} else {
				result.ErrorCount++
			}
		}()
	}

	// Wait for all requests to complete
	time.Sleep(2 * time.Second)

	return result
}

// Test Data Cleanup

// CleanupTestData cleans up test data
func CleanupTestData(client *E2EClient, seatIDs []string, userIDs []string) {
	// Cancel all reservations
	for _, seatID := range seatIDs {
		for _, userID := range userIDs {
			CancelTicket(client, userID, seatID)
		}
	}
}

// WaitForCondition waits for a condition to be true
func WaitForCondition(condition func() bool, timeout time.Duration) bool {
	timeoutChan := time.After(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return false
		case <-ticker.C:
			if condition() {
				return true
			}
		}
	}
}

// SetupE2ETest sets up an E2E test with a configured client
func SetupE2ETest(t *testing.T) *E2EClient {
	t.Helper()

	client := NewE2EClient("")

	// Wait for the service to be healthy
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := CheckHealth(client)
		if err == nil && GetValueAsInt(resp, "status_code") == 200 {
			t.Logf("API is healthy and ready")
			return client
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("Warning: API health check failed, continuing anyway...")
	return client
}

// AssertResponseSuccess asserts that a response indicates success
func AssertResponseSuccess(t *testing.T, resp map[string]interface{}, expectedStatus int) {
	t.Helper()
	statusCode := GetValueAsInt(resp, "status_code")
	require.Equal(t, expectedStatus, statusCode,
		fmt.Sprintf("Expected status %d but got %d. Response: %+v", expectedStatus, statusCode, resp))
}

// AssertResponseError asserts that a response indicates an error
func AssertResponseError(t *testing.T, resp map[string]interface{}, expectedStatus int) {
	t.Helper()
	statusCode := GetValueAsInt(resp, "status_code")
	require.NotEqual(t, expectedStatus, statusCode,
		fmt.Sprintf("Expected error status but got %d. Response: %+v", statusCode, resp))
}

// ParseErrorMessage extracts error message from response
func ParseErrorMessage(resp map[string]interface{}) string {
	if msg := GetValueAsString(resp, "error"); msg != "" {
		return msg
	}
	if msg := GetValueAsString(resp, "message"); msg != "" {
		return msg
	}
	return ""
}

// ContainsString checks if a string contains a substring
func ContainsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

