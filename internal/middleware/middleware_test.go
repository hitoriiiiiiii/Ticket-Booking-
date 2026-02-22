package middleware

import (
	"testing"
	"time"
)

// TestRateLimiter tests the RateLimiter struct
func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(5, time.Minute)
	
	if limiter == nil {
		t.Error("Expected non-nil RateLimiter")
	}
	if limiter.limit != 5 {
		t.Errorf("Expected limit 5, got %d", limiter.limit)
	}
	if limiter.window != time.Minute {
		t.Errorf("Expected window 1m, got %v", limiter.window)
	}
}

// TestRateLimiterAllow tests the Allow method
func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)
	
	// First 3 requests should be allowed
	for i := 0; i < 3; i++ {
		if !limiter.Allow("test_key") {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}
	
	// 4th request should be denied
	if limiter.Allow("test_key") {
		t.Error("Expected 4th request to be denied")
	}
}

// TestRateLimiterDifferentKeys tests different IP keys
func TestRateLimiterDifferentKeys(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute)
	
	// Each key should have its own limit
	if !limiter.Allow("key1") {
		t.Error("Expected first request for key1 to be allowed")
	}
	if !limiter.Allow("key2") {
		t.Error("Expected first request for key2 to be allowed")
	}
	
	// Second request for key1 should be denied
	if limiter.Allow("key1") {
		t.Error("Expected second request for key1 to be denied")
	}
}

// TestLogger tests the Logger middleware
func TestLogger(t *testing.T) {
	handler := Logger()
	if handler == nil {
		t.Error("Expected non-nil handler from Logger")
	}
}

// TestRateLimit tests the RateLimit middleware
func TestRateLimit(t *testing.T) {
	handler := RateLimit(60)
	if handler == nil {
		t.Error("Expected non-nil handler from RateLimit")
	}
}

// TestAuthMiddleware tests that AuthMiddleware exists
func TestAuthMiddleware(t *testing.T) {
	// Just verify the function exists and can be referenced
	// Actual testing would require a database mock
	_ = AuthMiddleware
}

// TestAdminMiddleware tests that AdminMiddleware exists
func TestAdminMiddleware(t *testing.T) {
	// Just verify the function exists
	_ = AdminMiddleware
}
