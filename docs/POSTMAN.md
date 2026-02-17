# Ticket Booking API - Postman Documentation

This document provides comprehensive Postman collection for testing all available endpoints in the Ticket Booking API.

## Base URL
```
http://localhost:8080
```

## Authentication

### Login to get JWT Token
**Endpoint:** `POST /users/login`

**Request Body:**
```
json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

**Response:**
```
json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Note:** For protected endpoints, include the JWT token in the Authorization header:
```
Authorization: Bearer <your_token>
```

---

## Health Check

### Check API Health
**Endpoint:** `GET /health`

**Response:**
```
json
{
  "status": "ok"
}
```

---

## User Endpoints

### 1. Register New User
**Endpoint:** `POST /users/register`

**Request Body:**
```
json
{
  "username": "johndoe",
  "email": "johndoe@example.com",
  "password": "securepassword123",
  "is_admin": false
}
```

**Response (201 Created):**
```
json
{
  "id": 1,
  "username": "johndoe",
  "email": "johndoe@example.com",
  "is_admin": false,
  "created_at": "2024-01-15T10:30:00Z"
}
```

---

### 2. Login User
**Endpoint:** `POST /users/login`

**Request Body:**
```
json
{
  "email": "johndoe@example.com",
  "password": "securepassword123"
}
```

**Response (200 OK):**
```
json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### 3. List All Users
**Endpoint:** `GET /users`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response (200 OK):**
```
json
[
  {
    "id": 1,
    "username": "johndoe",
    "email": "johndoe@example.com",
    "is_admin": false,
    "created_at": "2024-01-15T10:30:00Z"
  }
]
```

---

## Movie Endpoints

### 4. Get All Movies
**Endpoint:** `GET /movies`

**Response (200 OK):**
```
json
[
  {
    "id": 1,
    "name": "Inception",
    "genre": "Sci-Fi",
    "duration": 148
  },
  {
    "id": 2,
    "name": "The Dark Knight",
    "genre": "Action",
    "duration": 152
  }
]
```

---

### 5. Get Single Movie
**Endpoint:** `GET /movies/:id`

**Example:** `GET /movies/1`

**Response (200 OK):**
```
json
{
  "id": 1,
  "name": "Inception",
  "genre": "Sci-Fi",
  "duration": 148
}
```

**Error Response (404 Not Found):**
```
json
{
  "error": "Movie not found"
}
```

---

### 6. Create New Movie
**Endpoint:** `POST /movies`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Request Body:**
```
json
{
  "name": "Interstellar",
  "genre": "Sci-Fi",
  "duration": 169
}
```

**Response (201 Created):**
```
json
{
  "id": 3,
  "name": "Interstellar",
  "genre": "Sci-Fi",
  "duration": 169
}
```

---

## Show Endpoints

### 7. Get All Shows
**Endpoint:** `GET /shows`

**Response (200 OK):**
```
json
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

### 8. Create New Show
**Endpoint:** `POST /shows`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Request Body:**
```
json
{
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z"
}
```

**Note:** The `end_time` is automatically calculated based on the movie duration.

**Response (201 Created):**
```
json
{
  "id": 2,
  "movie_id": 1,
  "theater": "Theater A",
  "start_time": "2024-01-20T14:00:00Z",
  "end_time": "2024-01-20T16:30:00Z"
}
```

---

## Booking Endpoints

### 9. Reserve Ticket
**Endpoint:** `POST /reserve`

**Request Body:**
```
json
{
  "user_id": "1",
  "seat_id": "A1"
}
```

**Response (200 OK):**
```
json
{
  "message": "Ticket Reserved"
}
```

**Error Response (409 Conflict):**
```
json
{
  "error": "Seat already reserved"
}
```

---

### 10. Confirm Ticket Booking
**Endpoint:** `POST /confirm`

**Request Body:**
```
json
{
  "user_id": "1",
  "seat_id": "A1"
}
```

**Response (200 OK):**
```
json
{
  "message": "Ticket Confirmed"
}
```

**Error Response (409 Conflict):**
```
json
{
  "error": "Seat not held or already booked"
}
```

---

### 11. Cancel Ticket Reservation
**Endpoint:** `POST /cancel`

**Request Body:**
```
json
{
  "user_id": "1",
  "seat_id": "A1"
}
```

**Response (200 OK):**
```
json
{
  "message": "Ticket Cancelled"
}
```

**Error Response (404 Not Found):**
```
json
{
  "error": "No reservation found"
}
```

---

### 12. Check Seat Availability
**Endpoint:** `GET /availability/:seat_id`

**Example:** `GET /availability/A1`

**Response (200 OK):**
```
json
{
  "seat_id": "A1",
  "status": "AVAILABLE"
}
```

Or if reserved:
```
json
{
  "seat_id": "A1",
  "status": "HELD"
}
```

---

### 13. Get All Events (Admin)
**Endpoint:** `GET /events`

**Headers:**
```
Authorization: Bearer <admin_token>
```

**Response (200 OK):**
```
json
[
  {
    "seat_id": "A1",
    "type": "TicketReserved",
    "payload": "{\"UserID\":\"1\",\"SeatID\":\"A1\"}",
    "time": "2024-01-15T10:30:00Z"
  }
]
```

---

## Payment Endpoints

### 14. Initiate Payment
**Endpoint:** `POST /payments/initiate`

**Request Body:**
```
json
{
  "booking_id": "1",
  "user_id": "1",
  "amount": 500
}
```

**Response (200 OK):**
```
json
{
  "id": "pay_123abc",
  "booking_id": "1",
  "user_id": "1",
  "amount": 500,
  "status": "PENDING"
}
```

---

### 15. Verify Payment
**Endpoint:** `POST /payments/verify`

**Request Body:**
```
json
{
  "payment_id": "pay_123abc",
  "mode": "success"
}
```

**Note:** `mode` can be: `success`, `fail`, or `random`

**Response (200 OK):**
```
json
{
  "status": "Payment verification completed"
}
```

---

## Notification Endpoints

### 16. Get User Notifications
**Endpoint:** `GET /notifications/:user_id`

**Example:** `GET /notifications/1`

**Response (200 OK):**
```
json
[
  {
    "id": 1,
    "user_id": "1",
    "message": "Your ticket has been confirmed!",
    "type": "booking_confirmation",
    "created_at": "2024-01-15T10:30:00Z",
    "read": false
  }
]
```

---

## Testing Workflow Example

Here's a typical workflow for testing the booking system:

1. **Register a new user**
   
```
   POST /users/register
   
```

2. **Login to get token**
   
```
   POST /users/login
   
```

3. **View available movies**
   
```
   GET /movies
   
```

4. **Create a show**
   
```
   POST /shows (with admin token)
   
```

5. **Check seat availability**
   
```
   GET /availability/A1
   
```

6. **Reserve a ticket**
   
```
   POST /reserve
   
```

7. **Confirm the booking**
   
```
   POST /confirm
   
```

8. **Initiate payment**
   
```
   POST /payments/initiate
   
```

9. **Verify payment**
   
```
   POST /payments/verify
   
```

10. **Check notifications**
    
```
    GET /notifications/:user_id
    
```

---

## Error Responses

The API returns standard HTTP status codes:

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 404 | Not Found |
| 409 | Conflict |
| 500 | Internal Server Error |

---

## Environment Variables

Make sure to set up your environment variables in Postman:

```
DATABASE_URL=postgres://user:password@localhost:5432/ticket_booking
PORT=8080
