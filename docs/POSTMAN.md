# Ticket Booking API - Postman Documentation

This document provides comprehensive documentation for all available endpoints in the Ticket Booking API.

## Base URL

```
http://localhost:8081
```

---

## Table of Contents

- [API Endpoints Overview](#api-endpoints-overview)
- [Command Endpoints (Write Operations)](#command-endpoints-write-operations)
- [Query Endpoints (Read Operations)](#query-endpoints-read-operations)
- [System Endpoints](#system-endpoints)
- [Testing Workflow Example](#testing-workflow-example)
- [Error Responses](#error-responses)

---

## API Endpoints Overview

### Command Endpoints (Write Operations)

| Method | Endpoint                       | Description               |
| ------ | ------------------------------ | ------------------------- |
| POST   | `/cmd/reserve`                 | Reserve a ticket          |
| POST   | `/cmd/confirm`                 | Confirm ticket booking    |
| POST   | `/cmd/cancel`                  | Cancel ticket reservation |
| POST   | `/cmd/users/register`          | Register new user         |
| POST   | `/cmd/movies`                  | Create new movie          |
| PUT    | `/cmd/movies/:id`              | Update movie              |
| DELETE | `/cmd/movies/:id`              | Delete movie              |
| POST   | `/cmd/shows`                   | Create new show           |
| PUT    | `/cmd/shows/:id`               | Update show               |
| DELETE | `/cmd/shows/:id`               | Delete show               |
| POST   | `/cmd/payments/initiate`       | Initiate payment          |
| POST   | `/cmd/payments/verify`         | Verify payment            |
| POST   | `/cmd/payments/:id/refund`     | Refund payment            |

### Query Endpoints (Read Operations)

| Method | Endpoint                              | Description                |
| ------ | ------------------------------------- | -------------------------- |
| GET    | `/query/movies`                       | List all movies            |
| GET    | `/query/movies/:id`                   | Get movie by ID           |
| GET    | `/query/shows`                         | List all shows            |
| GET    | `/query/shows/:id`                     | Get show by ID            |
| GET    | `/query/shows/movie/:movieID`          | Get shows by movie        |
| GET    | `/query/availability/:seat_id`         | Check seat availability   |
| GET    | `/query/reservations/:user_id`         | Get user reservations     |
| GET    | `/query/users`                         | List all users            |
| GET    | `/query/users/:id`                     | Get user by ID            |
| POST   | `/query/users/login`                   | User login                |
| GET    | `/query/events`                        | Get all events            |
| GET    | `/query/notifications/:user_id`        | Get user notifications    |
| GET    | `/query/notifications/:user_id/unread` | Get unread notifications  |
| GET    | `/query/notifications/single/:id`      | Get notification by ID    |
| GET    | `/query/payments/:id`                  | Get payment by ID         |
| GET    | `/query/payments/booking/:bookingID`   | Get payment by booking    |
| GET    | `/query/payments/user/:userID`         | Get payments by user      |

### System Endpoints

| Method | Endpoint    | Description            |
| ------ | ----------- | ---------------------- |
| GET    | `/health`   | Health check           |
| GET    | `/metrics`  | Prometheus metrics     |

---

## Command Endpoints (Write Operations)

### 1. Reserve Ticket

**Endpoint:** `POST /cmd/reserve`

**Request Body:**

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "seat_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response (200 OK):**

```json
{
  "message": "Ticket Reserved"
}
```

---

### 2. Confirm Ticket

**Endpoint:** `POST /cmd/confirm`

**Request Body:**

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "seat_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response (200 OK):**

```json
{
  "message": "Ticket Confirmed"
}
```

---

### 3. Cancel Ticket

**Endpoint:** `POST /cmd/cancel`

**Request Body:**

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "seat_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response (200 OK):**

```json
{
  "message": "Ticket Cancelled"
}
```

---

### 4. Register User

**Endpoint:** `POST /cmd/users/register`

**Request Body:**

```json
{
  "username": "johndoe",
  "email": "johndoe@example.com",
  "password": "securepassword123",
  "is_admin": false
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "johndoe@example.com",
  "is_admin": false
}
```

---

### 5. Create Movie

**Endpoint:** `POST /cmd/movies`

**Request Body:**

```json
{
  "name": "Inception",
  "genre": "Sci-Fi",
  "duration": 148
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "name": "Inception",
  "genre": "Sci-Fi",
  "duration": 148
}
```

---

### 6. Update Movie

**Endpoint:** `PUT /cmd/movies/:id`

**Request Body:**

```json
{
  "name": "Inception",
  "genre": "Sci-Fi",
  "duration": 150
}
```

**Response (200 OK):**

```json
{
  "message": "Movie updated successfully"
}
```

---

### 7. Delete Movie

**Endpoint:** `DELETE /cmd/movies/:id`

**Response (200 OK):**

```json
{
  "message": "Movie deleted successfully"
}
```

---

### 8. Create Show

**Endpoint:** `POST /cmd/shows`

**Request Body:**

```json
{
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z"
}
```

**Response (201 Created):**

```json
{
  "id": 1,
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z",
  "end_time": "2024-01-20T16:30:00Z"
}
```

---

### 9. Update Show

**Endpoint:** `PUT /cmd/shows/:id`

**Request Body:**

```json
{
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z",
  "end_time": "2024-01-20T16:30:00Z"
}
```

**Response (200 OK):**

```json
{
  "message": "Show updated successfully"
}
```

---

### 10. Delete Show

**Endpoint:** `DELETE /cmd/shows/:id`

**Response (200 OK):**

```json
{
  "message": "Show deleted successfully"
}
```

---

### 11. Initiate Payment

**Endpoint:** `POST /cmd/payments/initiate`

**Request Body:**

```json
{
  "booking_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "amount": 500
}
```

**Response (200 OK):**

```json
{
  "id": "pay_123abc",
  "booking_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "amount": 500,
  "status": "PENDING"
}
```

---

### 12. Verify Payment

**Endpoint:** `POST /cmd/payments/verify`

**Request Body:**

```json
{
  "payment_id": "550e8400-e29b-41d4-a716-446655440099",
  "mode": "success"
}
```

**Note:** `mode` can be: `success`, `fail`, or `random`

**Response (200 OK):**

```json
{
  "status": "Payment verification completed"
}
```

---

### 13. Refund Payment

**Endpoint:** `POST /cmd/payments/:id/refund`

**Response (200 OK):**

```json
{
  "status": "Payment refund completed"
}
```

---

## Query Endpoints (Read Operations)

### 14. List Movies

**Endpoint:** `GET /query/movies`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "name": "Inception",
    "genre": "Sci-Fi",
    "duration": 148
  }
]
```

---

### 15. Get Movie by ID

**Endpoint:** `GET /query/movies/:id`

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "Inception",
  "genre": "Sci-Fi",
  "duration": 148
}
```

---

### 16. List Shows

**Endpoint:** `GET /query/shows`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "movie_id": 1,
    "theater": "Theater A",
    "start_time": "2024-01-20T14:00:00Z",
    "end_time": "2024-01-20T16:30:00Z"
  }
]
```

---

### 17. Get Show by ID

**Endpoint:** `GET /query/shows/:id`

**Response (200 OK):**

```json
{
  "id": 1,
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z",
  "end_time": "2024-01-20T16:30:00Z"
}
```

---

### 18. Get Shows by Movie

**Endpoint:** `GET /query/shows/movie/:movieID`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "movie_id": 1,
    "theater": "Theater A",
    "start_time": "2024-01-20T14:00:00Z"
  }
]
```

---

### 19. Check Seat Availability

**Endpoint:** `GET /query/availability/:seat_id`

**Response (200 OK):**

```json
{
  "seat_id": "550e8400-e29b-41d4-a716-446655440001",
  "status": "AVAILABLE"
}
```

---

### 20. Get User Reservations

**Endpoint:** `GET /query/reservations/:user_id`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "seat_id": "550e8400-e29b-41d4-a716-446655440001",
    "status": "BOOKED"
  }
]
```

---

### 21. List Users

**Endpoint:** `GET /query/users`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "username": "johndoe",
    "email": "johndoe@example.com",
    "is_admin": false
  }
]
```

---

### 22. Get User by ID

**Endpoint:** `GET /query/users/:id`

**Response (200 OK):**

```json
{
  "id": 1,
  "username": "johndoe",
  "email": "johndoe@example.com",
  "is_admin": false
}
```

---

### 23. User Login

**Endpoint:** `POST /query/users/login`

**Request Body:**

```json
{
  "email": "johndoe@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 24. Get All Events

**Endpoint:** `GET /query/events`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "aggregate_id": "550e8400-e29b-41d4-a716-446655440000",
    "event_type": "TicketReserved",
    "payload": {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "seat_id": "550e8400-e29b-41d4-a716-446655440001"
    },
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

---

### 25. Get User Notifications

**Endpoint:** `GET /query/notifications/:user_id`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Your ticket has been confirmed!",
    "type": "booking_confirmation",
    "created_at": "2024-01-15T10:30:00Z",
    "read": false
  }
]
```

---

### 26. Get Unread Notifications

**Endpoint:** `GET /query/notifications/:user_id/unread`

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Your ticket has been confirmed!",
    "type": "booking_confirmation",
    "created_at": "2024-01-15T10:30:00Z",
    "read": false
  }
]
```

---

### 27. Get Notification by ID

**Endpoint:** `GET /query/notifications/single/:id`

**Response (200 OK):**

```json
{
  "id": 1,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Your ticket has been confirmed!",
  "type": "booking_confirmation",
  "created_at": "2024-01-15T10:30:00Z",
  "read": false
}
```

---

### 28. Get Payment by ID

**Endpoint:** `GET /query/payments/:id`

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440099",
  "booking_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "amount": 500,
  "status": "COMPLETED"
}
```

---

### 29. Get Payment by Booking

**Endpoint:** `GET /query/payments/booking/:bookingID`

**Response (200 OK):**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440099",
  "booking_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "amount": 500,
  "status": "COMPLETED"
}
```

---

### 30. Get Payments by User

**Endpoint:** `GET /query/payments/user/:userID`

**Response (200 OK):**

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440099",
    "booking_id": "550e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "amount": 500,
    "status": "COMPLETED"
  }
]
```

---

## System Endpoints

### 31. Health Check

**Endpoint:** `GET /health`

**Response (200 OK):**

```json
{
  "status": "healthy"
}
```

---

### 32. Prometheus Metrics

**Endpoint:** `GET /metrics`

**Response:** Prometheus metrics in text format

---

## Testing Workflow Example

Here's a typical workflow for testing the booking system:

```bash
# 1. Register a new user
POST /cmd/users/register
{
  "username": "johndoe",
  "email": "johndoe@example.com",
  "password": "securepassword123",
  "is_admin": false
}

# 2. Login to get token
POST /query/users/login
{
  "email": "johndoe@example.com",
  "password": "securepassword123"
}

# 3. View available movies
GET /query/movies

# 4. Create a movie
POST /cmd/movies
{
  "name": "Inception",
  "genre": "Sci-Fi",
  "duration": 148
}

# 5. Create a show
POST /cmd/shows
{
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z"
}

# 6. View shows
GET /query/shows

# 7. Check seat availability
GET /query/availability/550e8400-e29b-41d4-a716-446655440001

# 8. Reserve a ticket
POST /cmd/reserve
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "seat_id": "550e8400-e29b-41d4-a716-446655440001"
}

# 9. Confirm the booking
POST /cmd/confirm
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "seat_id": "550e8400-e29b-41d4-a716-446655440001"
}

# 10. Initiate payment
POST /cmd/payments/initiate
{
  "booking_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "amount": 500
}

# 11. Verify payment
POST /cmd/payments/verify
{
  "payment_id": "550e8400-e29b-41d4-a716-446655440099",
  "mode": "success"
}

# 12. Check user reservations
GET /query/reservations/550e8400-e29b-41d4-a716-446655440000

# 13. Check notifications
GET /query/notifications/550e8400-e29b-41d4-a716-446655440000
```

---

## Error Responses

The API returns standard HTTP status codes:

| Status Code | Description           |
| ----------- | --------------------- |
| 200         | Success               |
| 201         | Created               |
| 400         | Bad Request           |
| 401         | Unauthorized          |
| 404         | Not Found             |
| 409         | Conflict              |
| 500         | Internal Server Error |

---

## Environment Variables

Make sure to set up your environment variables:

```
DATABASE_URL=postgres://user:password@localhost:5432/ticket_booking
PORT=8081
REDIS_URL=redis://localhost:6379
```

---

## Importing Postman Collection

1. Open Postman
2. Click "Import" button
3. Select the `postman-collection.json` file from the `docs/` folder
4. The collection will be imported with all endpoints
5. Set the `baseUrl` variable to `http://localhost:8081`

