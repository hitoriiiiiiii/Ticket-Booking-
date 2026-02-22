# Ticket Booking System - Testing Documentation

## Overview

This document outlines the comprehensive testing strategy for the Ticket Booking System. The project uses Go 1.25.5 with Gin framework, PostgreSQL databases (Command and Query), Redis for caching/queue/distributed locking, and follows CQRS (Command Query Responsibility Segregation) pattern with event-driven architecture.

## Project Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           TICKET BOOKING SYSTEM                         │
├─────────────────────────────────────────────────────────────────────────┤
│  API LAYER (Gin HTTP)                                                  │
│  ├── Command Handlers: /cmd/*                                          │
│  └── Query Handlers: /query/*                                          │
├─────────────────────────────────────────────────────────────────────────┤
│  SERVICE LAYER                                                         │
│  ├── Command Services: booking, user, movie, show, payment            │
│  └── Query Services: booking, user, movie, show, payment, notification│
├─────────────────────────────────────────────────────────────────────────┤
│  INFRASTRUCTURE LAYER                                                  │
│  ├── Command DB (PostgreSQL + pgx)                                     │
│  ├── Query DB (PostgreSQL + GORM)                                      │
│  ├── Redis: Cache, Queue, Distributed Lock                            │
│  ├── Event Dispatcher & Store                                          │
│  └── Notification Worker                                                │
└─────────────────────────────────────────────────────────────────────────┘
```

## Table of Contents

1. [Testing Philosophy](#testing-philosophy)
2. [Testing Types](#testing-types)
3. [Testing Dependencies](#testing-dependencies)
4. [Test Organization](#test-organization)
5. [Test Coverage Strategy](#test-coverage-strategy)
6. [Service-by-Service Testing Plan](#service-by-service-testing-plan)
7. [Integration Testing](#integration-testing)
8. [End-to-End Testing](#end-to-end-testing)
9. [Performance & Load Testing](#performance--load-testing)
10. [Test Fixtures & Utilities](#test-fixtures--utilities)
11. [Running Tests](#running-tests)
12. [CI/CD Integration](#cicd-integration)

---

## 1. Testing Philosophy

We follow the **Testing Pyramid** approach:

- **Unit Tests (70%)**: Fast, isolated tests for business logic
- **Integration Tests (20%)**: Testing interactions between components
- **E2E Tests (10%)**: Full workflow validation

### Core Principles

1. **Test Behavior, Not Implementation**: Focus on public interfaces
2. **Fast Feedback**: Unit tests should run in milliseconds
3. **Isolation**: Each test should be independent
4. **Clear Assertions**: Use descriptive assertion messages
5. **Mock External Dependencies**: Database, Redis, HTTP calls

---

## 2. Testing Types

### 2.1 Unit Tests
- Test individual functions and methods in isolation
- Use mocks for external dependencies (DB, Redis, event dispatcher)
- Focus on business logic and edge cases

### 2.2 Integration Tests
- Test multiple components working together
- Use test containers for database and Redis
- Test actual HTTP endpoints

### 2.3 Contract Tests
- Verify API contracts between services
- Ensure backward compatibility

### 2.4 Performance Tests
- Load testing with k6 or Gatling
- Concurrent booking stress tests
- Database connection pool testing

### 2.5 Security Tests
- Input validation
- Authentication/Authorization
- SQL injection prevention

---

## 3. Testing Dependencies

### Required Go Testing Packages

```
go
require (
    github.com/stretchr/testify v1.9.0          // Assertions and mocking
    github.com/DATA-DOG/go-sqlmock v1.5.2       // SQL mock for PostgreSQL
    github.com/go-faker/faker/v4 v4.4.1         // Test data generation
    github.com/gin-gonic/gin v1.11.0            // HTTP testing
    github.com/testcontainers/testcontainers-go v0.34.0 // Test containers
    github.com/redis/go-redis/v9 v9.18.0        // Redis client (for mocking)
    gorm.io/gorm v1.31.1                        // GORM (for mocking)
    github.com/jackc/pgx/v5 v5.8.0              // pgx (for mocking)
)
```

---

## 4. Test Organization

### Directory Structure



### Naming Conventions

- **Unit Tests**: `*_test.go`
- **Test Files**: `package_test` (black-box testing)
- **Fixtures**: `fixtures.go` or `testdata/`
- **Helper Functions**: `test_utils.go`

---

## 5. Test Coverage Strategy

### Target Coverage by Layer

| Layer | Minimum Coverage | Target Coverage |
|-------|------------------|------------------|
| Command Services | 80% | 90% |
| Query Services | 75% | 85% |
| Handlers | 70% | 80% |
| Middleware | 85% | 95% |
| Utils | 80% | 90% |
| Event System | 75% | 85% |

### Critical Paths to Cover

1. **Booking Flow**: Reserve → Confirm → Payment → Notification
2. **User Registration**: Register → Login → Token Generation
3. **Concurrency**: Multiple users booking same seat
4. **Error Handling**: All error cases and edge conditions

---

## 6. Service-by-Service Testing Plan

### 6.1 User Service

#### Command Service Tests (`user/command_service_test.go`)

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestRegisterUser_Success` | Valid user registration | User created, no password in response |
| `TestRegisterUser_DuplicateEmail` | Same email registration | Error: "email already exists" |
| `TestRegisterUser_InvalidEmail` | Invalid email format | Error: "invalid email format" |
| `TestRegisterUser_EmptyPassword` | Empty password | Error: "password required" |
| `TestRegisterUser_ShortPassword` | Password < 8 chars | Error: "password too short" |

#### Query Service Tests (`user/query_service_test.go`)

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestGetUserByID_Exists` | Valid user ID | Returns user data |
| `TestGetUserByID_NotFound` | Non-existent ID | Returns nil |
| `TestListUsers_Pagination` | List with limit/offset | Correct pagination |

#### Handler Tests

```
go
// Example handler test
func TestRegisterHandler(t *testing.T) {
    // Setup gin router
    // Mock service
    // Make HTTP request
    // Assert status code and response
}
```

---

### 6.2 Movie Service

#### Command Service Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestCreateMovie_Success` | Valid movie data | Movie created with ID |
| `TestCreateMovie_InvalidTitle` | Empty title | Error: "title required" |
| `TestCreateMovie_NegativeDuration` | Duration < 0 | Error: "invalid duration" |
| `TestCreateMovie_EmitEvent` | Movie creation | EventUserRegistered emitted |

#### Query Service Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestGetMovieByID_Exists` | Valid ID | Returns movie details |
| `TestGetMovies_List` | List all movies | Returns movie list |
| `TestGetMovie_NotFound` | Invalid ID | Returns nil |

---

### 6.3 Show Service

#### Command Service Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestCreateShow_Success` | Valid show data | Show created |
| `TestCreateShow_Overlap` | Overlapping shows | Error: "showtime conflict" |
| `TestCreateShow_PastDate` | Past show date | Error: "cannot create show in past" |
| `TestCreateShow_InvalidMovie` | Non-existent movie | Error: "movie not found" |

---

### 6.4 Booking Service (CRITICAL)

#### Command Service Tests

```
go
// Test file: booking/command_service_test.go
```

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestReserveTicket_Success` | Valid reservation | Ticket reserved, status=HELD |
| `TestReserveTicket_AlreadyReserved` | Double booking | Error: "seat already reserved" |
| `TestReserveTicket_EmitEvent` | Reservation made | EventTicketReserved emitted |
| `TestReserveTicket_WithLock` | Concurrent booking | Lock acquired, prevent race |
| `TestConfirmTicket_Success` | Valid confirmation | Status changes to BOOKED |
| `TestConfirmTicket_NotHeld` | Confirm without reserve | Error: "seat not held" |
| `TestConfirmTicket_EmitEvent` | Confirmation made | EventTicketConfirmed emitted |
| `TestCancelTicket_Success` | Valid cancellation | Reservation deleted |
| `TestCancelTicket_NotFound` | Cancel non-existent | Error: "no reservation found" |

#### Distributed Lock Tests

```
go
func TestDistributedLock_ConcurrentBooking(t *testing.T) {
    // Setup Redis mock
    // Create 10 concurrent reservation requests for same seat
    // Only 1 should succeed, others should wait/retry
}
```

---

### 6.5 Payment Service

#### Command Service Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestInitiatePayment_Success` | Valid payment | Payment record created |
| `TestInitiatePayment_InvalidAmount` | Amount <= 0 | Error: "invalid amount" |
| `TestVerifyPayment_Success` | Valid verification | Payment status = VERIFIED |
| `TestVerifyPayment_AlreadyVerified` | Double verify | Error: "already verified" |
| `TestVerifyPayment_EmitEvent` | Payment verified | EventPaymentVerified emitted |

#### Query Service Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestGetPaymentByID` | Valid payment ID | Returns payment details |
| `TestGetUserPayments` | User payment history | Returns payment list |

---

### 6.6 Notification Service

#### Query Service Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestGetUserNotifications` | Get user notifications | Returns notification list |
| `TestGetNotifications_Unread` | Filter unread only | Returns unread count |

#### Worker Tests

```
go
func TestNotificationWorker_ProcessJob(t *testing.T) {
    // Enqueue test job
    // Process job
    // Verify notification sent
}
```

---

### 6.7 Middleware Tests

#### Rate Limiter Tests (`middleware/rate_limiter_test.go`)

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestRateLimit_UnderLimit` | Normal requests | All succeed (200) |
| `TestRateLimit_OverLimit` | Excessive requests | Some return 429 |
| `TestRateLimit_Reset` | After window reset | Requests succeed again |
| `TestRateLimit_Distributed` | Multiple instances | Shared counter |

#### Auth Middleware Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestAuth_ValidToken` | Valid JWT | Request proceeds |
| `TestAuth_ExpiredToken` | Expired JWT | Error: 401 |
| `TestAuth_MissingToken` | No token | Error: 401 |
| `TestAuth_InvalidToken` | Malformed token | Error: 401 |

---

### 6.8 Event System Tests

#### Dispatcher Tests (`events/dispatcher_test.go`)

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestPublish_Success` | Valid event | Event stored and published |
| `TestPublish_Async` | Async publish | Non-blocking publish |
| `TestSubscribe_Called` | Subscriber receives event | Callback invoked |
| `TestSubscribe_Multiple` | Multiple subscribers | All callbacks invoked |

#### Store Tests

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestEventStore_Save` | Save event | Event persisted |
| `TestEventStore_GetByAggregate` | Get aggregate events | Returns event list |

---

### 6.9 Utils Tests

#### Distributed Lock Tests (`utils/lock_test.go`)

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestLock_Acquire` | Lock acquisition | Lock acquired |
| `TestLock_Release` | Lock release | Lock released |
| `TestLock_Timeout` | Lock timeout | Returns error after timeout |
| `TestLock_Reacquire` | Reacquire after release | Success |

#### Cache Tests (`utils/cache_test.go`)

| Test Case | Description | Expected Behavior |
|-----------|-------------|-------------------|
| `TestCache_SetGet` | Basic set/get | Correct value returned |
| `TestCache_Expiration` | TTL expiration | Value expires |
| `TestCache_Delete` | Delete key | Key removed |

---

## 7. Integration Testing

### Database Setup

Use testcontainers for isolated testing:

```
go
func TestMain(m *testing.M) {
    // Start PostgreSQL container
    // Start Redis container
    // Run tests
    // Cleanup
}
```

### Integration Test Examples

```
go
func TestBookingIntegration_ReserveAndConfirm(t *testing.T) {
    // Setup test database
    // Create user
    // Reserve ticket
    // Confirm ticket
    // Verify in database
    // Cleanup
}
```

### Test Fixtures

```
go
// fixtures/user_fixtures.go
var (
    TestUser = &user.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }
)
```

---

## 8. End-to-End Testing

### Test Scenarios

1. **Happy Path**: User registers → Searches movie → Reserves seat → Confirms → Pays
2. **Cancellation**: User reserves → Cancels → Seat available
3. **Concurrent Booking**: Multiple users try to book same seat
4. **Payment Failure**: Payment verification fails

### E2E Test Example

```
go
func TestE2E_BookingFlow(t *testing.T) {
    // Register user
    // Login and get token
    // Get available shows
    // Reserve ticket
    // Confirm ticket
    // Initiate payment
    // Verify payment
    // Get notification
}
```

---

## 9. Performance & Load Testing

### Tools

- **k6**: Load testing (recommended)
- **Gatling**: Alternative
- **Apache Bench**: Simple load testing

### Load Test Scenarios

| Scenario | Users | Duration | Target |
|----------|-------|----------|--------|
| Normal Load | 1000 | 5 min | < 200ms p95 |
| Peak Load | 10000 | 2 min | < 500ms p95 |
| Stress Test | 50000 | 1 min | No errors |

### Specific Load Tests

1. **Concurrent Booking**: 100 users booking same seat simultaneously
2. **Rate Limiting**: Verify rate limiter works under load
3. **Database Pool**: Test connection pool exhaustion
4. **Redis Cache**: Test cache hit/miss ratio

---

## 10. Test Fixtures & Utilities

### Test Helper Package

Create `internal/testutils/` with:

```
go
// testutils/database.go
func SetupTestDB() (*sql.DB, func()) {
    // Create test database
    // Run migrations
    // Return cleanup function
}

// testutils/mock.go
func MockUserService() *MockUserService {
    // Return mock implementation
}

// testutils/fixtures.go
func LoadFixtures(path string) map[string]interface{} {
    // Load JSON fixtures
}
```

---

## 11. Running Tests

### Run All Tests

```
bash
# Run all unit tests
go test ./internal/... -v

# Run with coverage
go test ./internal/... -coverprofile=coverage.out -covermode=atomic

# Run specific package
go test ./internal/booking/... -v
```

### Run Specific Test Types

```
bash
# Unit tests only
go test ./internal/... -tags=unit -v

# Integration tests
go test ./internal/... -tags=integration -v

# Run with race detector
go test ./internal/... -race -v
```

### Test Flags

```
bash
# Verbose output
go test -v

# Stop on first failure
go test -failfast

# Run specific test
go test -run TestReserveTicket_Success

# Timeout
go test -timeout 5m
```

---

## 12. CI/CD Integration

### GitHub Actions Example

```
yaml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          
      - name: Run Unit Tests
        run: go test ./internal/... -tags=unit -race -coverprofile=coverage.out
        
      - name: Run Integration Tests
        run: go test ./internal/... -tags=integration
        
      - name: Upload Coverage
        uses: codecov/codecov-action@v4
        
      - name: Run Security Tests
        run: go run golang.org/x/vuln/cmd/govulncheck ./...
```

---

## Testing Checklist

Before each release, verify:

- [ ] All unit tests pass (> 80% coverage)
- [ ] Integration tests pass
- [ ] E2E tests pass
- [ ] Load tests meet performance targets
- [ ] No race conditions detected (`go test -race`)
- [ ] Security tests pass
- [ ] Documentation updated

---

## Additional Resources

- [Go Testing Documentation](https://go.dev/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [TestContainers Go](https://github.com/testcontainers/testcontainers-go)
- [k6 Load Testing](https://k6.io/docs/)

---

*Document Version: 1.0*  

