# Ticket-Booking System

A high-performance backend for ticket booking applications, inspired by platforms like BookMyShow. Designed to handle 50K+ concurrent users with robust concurrency control and scalability.

---

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Docker Compose](#docker-compose)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [Testing](#testing)
- [Scaling](#scaling)
- [License](#license)

---

## Features

- Users browse movies
- Select theater + showtime
- Pick seats
- Seat lock for few minutes
- Payment
- Booking confirmation
- Prevent double booking
- Handle 50K+ concurrent users

---

## Services

| Service              | Responsibility    |
| -------------------- | ----------------- |
| User Service         | Auth, profiles    |
| Movie Service        | Movies, metadata  |
| Show Service         | Showtimes         |
| Booking Service      | Reserve + confirm |
| Payment Service      | Payments          |
| Notification Service | SMS/email         |

---

## System Design Technologies

### Backend Architecture

| Technology               | Description                                                           |
| ------------------------ | --------------------------------------------------------------------- |
| **Golang (Go)**          | High-performance backend service with native concurrency support      |
| **RESTful API**          | API layer for booking operations using Gin web framework              |
| **Layered Architecture** | Clean separation between Controllers, Services, and Repository layers |

---

### CQRS Pattern

> **Command Query Responsibility Segregation**

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                    CQRS ARCHITECTURE                                                │
├───────────────────────────┬─────────────────────────────┬──────────────────────────────────────────┤
│      COMMANDS             │      EVENT STORE            │           QUERIES                        │
│   (Write Side)            │                             │         (Read Side)                     │
├───────────────────────────┼─────────────────────────────┼──────────────────────────────────────────┤
│                           │                             │                                          │
│   ┌─────────────┐         │   ┌─────────────────┐       │   ┌─────────────┐                        │
│   │ Seat        │         │   │  events table   │       │   │ Fetch Shows │                        │
│   │ Reservation │────────────►│                 │       │   │             │                        │
│   └─────────────┘         │   │ • aggregate_id  │       │   └──────┬──────┘                        │
│           │               │   │ • event_type   │       │          │                               │
│           ▼               │   │ • payload      │       │          ▼                               │
│   ┌─────────────┐         │   │ • created_at   │       │   ┌─────────────┐                        │
│   │  Payment    │────────────►│                 │       │   │   Seat      │                        │
│   │  Processing │         │   └────────┬────────┘       │   │Availability │                        │
│   └─────────────┘         │            │                 │   └──────┬──────┘                        │
│           │               │            │                 │          │                               │
│           ▼               │            ▼                 │          ▼                               │
│   ┌─────────────┐         │   ┌─────────────────┐       │   ┌─────────────┐                        │
│   │ Reservation │         │   │ Projection      │       │   │    User     │                        │
│   │ Confirmation│         │   │ Worker          │       │   │  Bookings   │                        │
│   └─────────────┘         │   └────────┬────────┘       │   └─────────────┘                        │
│                           │            │                 │                                          │
└───────────────────────────┼────────────┼─────────────────┼──────────────────────────────────────────┘
                            │            │                 │
                            │            ▼                 │
                            │   ┌─────────────────┐       │
                            │   │ Reservation     │       │
                            │   │ Projection      │       │
                            │   │                 │       │
                            │   │ • TicketReserved│       │
                            │   │ • TicketCancel  │       │
                            │   └────────┬────────┘       │
                            │            │                 │
                            │            ▼                 │
                            │   ┌─────────────────┐       │
                            │   │   Read Model    │       │
                            │   │ (reservations)  │       │
                            │   │                 │       │
                            │   │ • seat_id       │       │
                            │   │ • user_id       │       │
                            │   │ • status        │       │
                            │   └─────────────────┘       │
                            │                             │
                            └─────────────────────────────┘
```

**CQRS Flow:**

1. **Commands (Write Side)**: User actions like seat reservation, payment processing are handled as commands
2. **Event Store**: Commands are stored as events in the `events` table (event sourcing)
3. **Projection Worker**: `ReservationProjection` reads events from the event store and updates the read model
4. **Read Model**: The `reservations` table is updated with the latest booking status
5. **Queries (Read Side)**: Read operations like fetching shows, checking seat availability, and retrieving user bookings query the optimized read model

- **Command Side** (`/cmd/*`): Handles seat booking, payment, reservation updates → writes to Event Store
- **Query Side** (`/query/*`): Optimized for fetching shows, seat availability, user bookings → reads from Read Model
- **Benefits**: Improves performance under heavy concurrent traffic, separates read/write concerns

---

### CQRS Data Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐
│                                                CQRS DATA FLOW                                                             │
├─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                                                 │
│    ┌──────────────────────┐                                    ┌──────────────────────┐                                    │
│    │      COMMAND API     │                                    │      QUERY API      │                                    │
│    │  (Reserve/Cancel)    │                                    │   (Availability)    │                                    │
│    └──────────┬───────────┘                                    └──────────┬───────────┘                                    │
│               │                                                        │                                                │
│               │  Write DB (Postgres)                                  │  Read DB (Postgres Replica)                   │
│               ▼                                                        ▼                                                │
│    ┌──────────────────────┐                                    ┌──────────────────────┐                                    │
│    │    Primary DB        │                                    │   Read Replica      │                                    │
│    │  (Write Operations) │                                    │  (Read Operations)  │                                    │
│    │  • Reservations     │                                    │  • Shows            │                                    │
│    │  • Payments         │                                    │  • Seat Availability│                                    │
│    │  • Users            │                                    │  • User Bookings    │                                    │
│    └──────────┬───────────┘                                    └──────────────────────┘                                    │
│               │                                                                                                            │
│               │                                                                                                            │
│               ▼                                                                                                            │
│    ┌──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┐    │
│    │                                         KAFKA TOPIC: booking.events                                                │    │
│    ├──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤    │
│    │                                                                                                                                  │
│    │   ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐                   │    │
│    │   │ ReservationCreated│    │  SeatLocked      │    │BookingConfirmed  │    │  SeatReleased    │                   │    │
│    │   └────────┬─────────┘    └────────┬─────────┘    └────────┬─────────┘    └────────┬─────────┘                   │    │
│    │            │                      │                      │                      │                             │    │
│    │            ▼                      ▼                      ▼                      ▼                             │    │
│    │   ┌────────────────────────────────────────────────────────────────────────────────────────────────────────┐   │    │
│    │   │                                    EVENT STORE                                                      │   │    │
│    │   │   • aggregate_id    • event_type    • payload (JSON)    • created_at                              │   │    │
│    │   └────────────────────────────────────────────────────────────────────────────────────────────────────────┘   │    │
│    │                                                                                                            │    │
│    └──────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘    │
│               │                                                                                                            │
│               │                                                                                                            │
│               ▼                                                                                                            │
│    ┌──────────────────────┐                                                                                                │
│    │  PROJECTION WORKER   │                                                                                                │
│    │    (Consumer)        │                                                                                                │
│    ├──────────────────────┤                                                                                                │
│    │  • Reads events      │                                                                                                │
│    │  • Updates read model│                                                                                                │
│    │  • Transforms data   │                                                                                                │
│    └──────────┬───────────┘                                                                                                │
│               │                                                                                                            │
│               │  Updates Read DB                                                                                           │
│               ▼                                                                                                            │
│    ┌──────────────────────┐                                                                                                │
│    │      READ DB         │                                                                                                │
│    │  (Projections)       │                                                                                                │
│    │  • seat_availability │                                                                                                │
│    │  • user_bookings     │                                                                                                │
│    │  • show_stats        │                                                                                                │
│    └──────────────────────┘                                                                                                │
│                                                                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┘
```

**Data Flow:**

1. **Command API** handles Reserve/Cancel operations → writes to Primary DB (Postgres)
2. **Query API** handles Availability queries → reads from Read Replica (Postgres)
3. **Events** are published to Kafka Topic: `booking.events`
4. **Projection Worker** consumes events from Kafka and updates the Read DB
5. **Read DB** stores projections for optimized query performance

---

### Database & Storage

| Feature               | Implementation                                                       |
| --------------------- | -------------------------------------------------------------------- |
| **PostgreSQL 15**     | Relational DB for shows, seats, bookings (Command & Query databases) |
| **ACID Transactions** | Ensures safe booking without double reservation                      |
| **Database Indexing** | Faster show/search queries                                           |
| **GORM**              | ORM for database operations                                          |

---

### Concurrency & Seat Locking

```
┌─────────────────────────────────────────────────────────────┐
│              SEAT LOCKING MECHANISM                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   User A              User B              User C             │
│     │                   │                   │                │
│     ▼                   ▼                   ▼                │
│ ┌─────────┐        ┌─────────┐        ┌─────────┐           │
│ │ Lock    │        │ Lock    │        │ Lock    │           │
│ │ Seat 5  │        │ Seat 5  │        │ Seat 6  │           │
│ │ (10 sec)│        │ (FAIL!) │        │ (OK!)   │           │
│ └─────────┘        └─────────┘        └─────────┘           │
│     │                                      │                │
│     ▼                                      ▼                │
│ ┌─────────┐                          ┌─────────┐            │
│ │ Book/   │                          │ Book/   │            │
│ │ Release │                          │ Release │            │
│ └─────────┘                          └─────────┘            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

- **Seat Locking**: Prevents double booking with configurable timeout (10 seconds)
- **Optimistic Locking**: Version-based conflict detection
- **Pessimistic Locking**: Database-level row locking
- **Go Concurrency**: Goroutines & Mutex Locks for thread-safe operations
- **Distributed Locking**: Redis-based distributed locks for 50K+ concurrent users

---

### Caching Layer

```
┌─────────────────────────────────────────────────────────────┐
│                    REDIS CACHE LAYER                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌─────────────────┐    ┌─────────────────┐                │
│  │  Seat Locks     │    │  Show Cache    │                │
│  │  ─────────────  │    │  ────────────  │                │
│  │  seat:5:user:1  │    │  show:123       │                │
│  │  TTL: 10 min    │    │  TTL: 5 min     │                │
│  └─────────────────┘    └─────────────────┘                │
│                                                              │
│  ┌─────────────────┐    ┌─────────────────┐                │
│  │  Rate Limiting  │    │  Session Cache  │                │
│  │  ─────────────  │    │  ─────────────  │                │
│  │  ip:192.168.1.1 │    │  jwt:token123   │                │
│  │  count: 100/min │    │  TTL: 24 hrs    │                │
│  └─────────────────┘    └─────────────────┘                │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

> **Redis** for real-time seat availability and temporary locks

- **Real-time seat availability**
- **Temporary seat locks** (TTL: 10 seconds)
- **Session management**
- **Rate limiting** (5000 requests/minute)
- **Job Queue**: Redis-based job queue for background processing
- Reduces database load significantly

---

### Job Queues & Background Processing

```
┌─────────────────────────────────────────────────────────────┐
│                   JOB QUEUE SYSTEM                           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│    API Server              Worker              External      │
│       │                     │                    │          │
│       │  ┌──────────┐       │                    │          │
│       ├──►│ Enqueue  │──────┤                    │          │
│       │  │  Job     │       │                    │          │
│       │  └──────────┘       │                    │          │
│       │                     │  ┌──────────────┐   │          │
│       │                     ├──►│ Process Job  │   │          │
│       │                     │  └──────────────┘   │          │
│       │                     │         │            │          │
│       │                     │         ▼            ▼          │
│       │                     │   ┌─────────────────────┐      │
│       │                     │   │ • Email/SMS        │      │
│       │                     │   │ • Payment Verify   │      │
│       │                     │   │ • Expiry Cleanup   │      │
│       │                     │   │ • Analytics        │      │
│       │                     │   └─────────────────────┘      │
│       │                     │                                 │
│       ▼                     ▼                                 │
│   (Response)          (Background)                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Async Operations:**

- ✅ Booking confirmation emails/SMS
- ✅ Payment verification
- ✅ Reservation expiry cleanup
- ✅ Analytics logging
- ✅ Event dispatching via Kafka

---

### Event-Driven System Design

```
┌─────────────────────────────────────────────────────────────┐
│                   EVENT-DRIVEN FLOW                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │  Reservation │───►│   SeatLock   │───►│    Payment   │ │
│  │   Created    │    │    Event     │    │   Verified   │ │
│  └──────────────┘    └──────────────┘    └──────────────┘ │
│         │                   │                   │            │
│         ▼                   ▼                   ▼            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   EVENT STORE                          │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐      │  │
│  │  │Event #1 │ │Event #2 │ │Event #3 │ │Event #4 │ ...  │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘      │  │
│  └──────────────────────────────────────────────────────┘  │
│         │                   │                   │            │
│         ▼                   ▼                   ▼            │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐ │
│  │   Booking    │    │    User      │    │   Payment    │ │
│  │  Confirmed   │    │  Notification│    │    Record    │ │
│  └──────────────┘    └──────────────┘    └──────────────┘ │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Events:**

- `TicketReserved` - Initial seat reservation
- `SeatLocked` - Seat temporarily locked
- `TicketConfirmed` - Payment successful
- `TicketCancelled` - Lock expired or cancelled
- `PaymentVerified` - Payment verification completed
- `UserRegistered` - New user registration

---

### Authentication & Security

```
┌─────────────────────────────────────────────────────────────┐
│                    AUTHENTICATION FLOW                       │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   Client                                                    │
│     │                                                       │
│     │ 1. Login (email/password)                             │
│     ▼                                                       │
│   ┌─────────────────────┐                                   │
│   │   Generate JWT      │                                   │
│   │   Token             │                                   │
│   └─────────────────────┘                                   │
│     │                                                       │
│     │ 2. Bearer Token                                       │
│     ▼                                                       │
│   ┌─────────────────────┐    ┌─────────────────────┐       │
│   │   Auth Middleware   │───►│   Verify Token      │       │
│   │   (Gin)             │    │   Extract User ID   │       │
│   └─────────────────────┘    └─────────────────────┘       │
│     │                               │                        │
│     │ 3. Request + User Context     │                        │
│     ▼                               ▼                        │
│   ┌─────────────────────┐    ┌─────────────────────┐       │
│   │   Admin?            │    │   Regular User?    │       │
│   │   Full Access       │    │   Limited Access   │       │
│   └─────────────────────┘    └─────────────────────┘       │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

| Feature                              | Implementation                                |
| ------------------------------------ | --------------------------------------------- |
| **JWT Authentication**               | Secure session handling with token-based auth |
| **Role-Based Access Control (RBAC)** | Admin & User roles                            |
| **Rate Limiting**                    | Prevent API abuse (5000 req/min)              |
| **Structured Logging**               | Debugging & tracing                           |

---

### Scalability Concepts

```
                    ┌─────────────────────────────────────┐
                    │         LOAD BALANCER               │
                    │        (nginx/haproxy)              │
                    └──────────────┬──────────────────────┘
                                   │
           ┌───────────────────────┼───────────────────────┐
           │                       │                       │
           ▼                       ▼                       ▼
    ┌─────────────┐        ┌─────────────┐        ┌─────────────┐
    │   API Pod 1 │        │   API Pod 2 │        │   API Pod N │
    │  (Go)       │        │  (Go)       │        │  (Go)       │
    └──────┬──────┘        └──────┬──────┘        └──────┬──────┘
           │                       │                       │
           └───────────────────────┼───────────────────────┘
                                   │
                    ┌──────────────┴──────────────┐
                    │      SHARED INFRASTRUCTURE  │
                    ├─────────────────────────────┤
                    │  • PostgreSQL (Primary)     │
                    │  • Redis Cluster           │
                    │  • Message Queue           │
                    └─────────────────────────────┘
```

- **Horizontal Scaling**: Supports multiple backend instances
- **Load Balancing**: Handles high traffic during peak booking times
- **Microservice-Friendly**: Booking, Payment, Movies can be separated

---

### DevOps & Deployment

| Tool               | Purpose                                        |
| ------------------ | ---------------------------------------------- |
| **Docker**         | Containerized backend                          |
| **Docker Compose** | Multi-service setup (DB + Redis + API + Kafka) |
| **Kubernetes**     | Production deployment with HPA                 |
| **Nginx**          | Load balancer                                  |

---

### Monitoring & Logging

| Feature                  | Status          |
| ------------------------ | --------------- |
| **Structured Logging**   | ✅ Implemented  |
| **Request Logging**      | ✅ Implemented  |
| **Error Tracing**        | ✅ Implemented  |
| **Prometheus + Grafana** | Future Scope |

---

### Testing & Reliability

| Type                    | Description              |
| ----------------------- | ------------------------ |
| **Unit Testing**        | Go testing package       |
| **Integration Testing** | Booking workflow tests   |
| **API Testing**         | Postman collection ready |

---

## Technology Stack

| Category       | Technology    |
| -------------- | ------------- |
| Language       | Go 1.21+      |
| Framework      | Gin           |
| Database       | PostgreSQL 15 |
| Cache          | Redis 7       |
| ORM            | GORM          |
| Authentication | JWT           |
| Message Queue  | Kafka         |
| Container      | Docker        |
| Orchestration  | Kubernetes    |

---

## Getting Started

### Prerequisites

- **Go** 1.21 or higher
- **Docker** & Docker Compose
- **PostgreSQL** 15
- **Redis** 7

### Run Locally

```
bash
# Install dependencies
go mod download

# Run API server
go run cmd/api/main.go

# Run background worker (in separate terminal)
go run cmd/worker/main.go
```

### Environment Variables

Create a `.env` file in the root directory:

````
# Server
PORT=8080

# Kafka
KAFKA_BROKER=localhost:9092
 
````

## Docker Compose

### Services Overview

The Docker Compose setup includes the following services:

| Service          | Container Name        | Port | Description              |
| ---------------- | --------------------- | ---- | ------------------------ |
| `postgres_cmd`   | ticket-cmd-db         | 5433 | Comma+nd Database (Write) |
| `postgres_query` | ticket-query-db       | 5434 | Query Database (Read)    |
| `redis`          | ticket-booking-redis  | 6379 | Cache & Queue            |
| `zookeeper`      | ticket-zookeeper      | 2181 | Kafka Manager            |
| `kafka`          | ticket-kafka          | 9092 | Message Broker           |
| `api`            | ticket-booking-api    | 8081 | API Server               |
| `worker`         | ticket-booking-worker | -    | Background Worker        |

### Quick Start with Docker Compose

```
bash
# Clone the repository
git clone https://github.com/hitorii/ticket-booking.git
cd ticket-booking

# Create .env file
echo "DB_USER=ticket" > .env
echo "DB_PASSWORD=ticket123" >> .env

# Start all services
docker-compose up -d

# Access API
http://localhost:8081

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Database Access

```
bash
# Connect to Command Database
docker exec -it ticket-cmd-db psql -U ticket -d ticket_cmd_db

# Connect to Query Database
docker exec -it ticket-query-db psql -U ticket -d ticket_query_db
```

### Database Connection Details

| Database       | Host          | Port | Database Name     | Username | Password   |
| ---------------|---------------|------|------------------|----------|------------|
| Command DB     | localhost     | 5433 | ticket_cmd_db    | ticket   | ticket123  |
| Query DB       | localhost     | 5434 | ticket_query_db  | ticket   | ticket123  |
| Redis          | localhost     | 6379 | -                | -        | -          |

### Environment Variables

```
bash
# Command Database (Write)
export COMMAND_DATABASE_URL="postgres://ticket:ticket123@localhost:5433/ticket_cmd_db?sslmode=disable"

# Query Database (Read)
export QUERY_DATABASE_URL="postgres://ticket:ticket123@localhost:5434/ticket_query_db?sslmode=disable"

# Redis
export REDIS_URL="redis://localhost:6379"
```

---

## API Endpoints

> **For complete API documentation with request/response examples, refer to [Postman Documentation](./docs/POSTMAN.md)**

### Base URL

```
http://localhost:8081
```

### Command Endpoints (Write Operations)

| Method | Endpoint                 | Description               | Request Body                                                                                  |
| ------ | ------------------------ | ------------------------- | --------------------------------------------------------------------------------------------- |
| POST   | `/cmd/reserve`           | Reserve a ticket          | `{"user_id": "1", "seat_id": "A1"}`                                                           |
| POST   | `/cmd/confirm`           | Confirm ticket booking    | `{"user_id": "1", "seat_id": "A1"}`                                                           |
| POST   | `/cmd/cancel`            | Cancel ticket reservation | `{"user_id": "1", "seat_id": "A1"}`                                                           |
| POST   | `/cmd/users/register`    | Register new user         | `{"username": "john", "email": "john@example.com", "password": "pass123", "is_admin": false}` |
| POST   | `/cmd/movies`            | Create new movie          | `{"name": "Inception", "genre": "Sci-Fi", "duration": 148}`                                   |
| POST   | `/cmd/shows`             | Create new show           | `{"movie_id": 1, "theater": "Theater A", "start_time": "2024-01-20T14:00:00Z"}`               |
| POST   | `/cmd/payments/initiate` | Initiate payment          | `{"booking_id": "1", "user_id": "1", "amount": 500}`                                          |
| POST   | `/cmd/payments/verify`   | Verify payment            | `{"payment_id": "pay_xxx", "mode": "success"}`                                                |

### Query Endpoints (Read Operations)

| Method | Endpoint                        | Description             |
| ------ | ------------------------------- | ----------------------- |
| GET    | `/query/movies`                 | List all movies         |
| GET    | `/query/movies/:id`             | Get movie by ID         |
| GET    | `/query/shows`                  | List all shows          |
| GET    | `/query/availability/:seat_id`  | Check seat availability |
| GET    | `/query/reservations/:user_id`  | Get user reservations   |
| GET    | `/query/users`                  | List all users          |
| GET    | `/query/events`                 | Get all events          |
| GET    | `/query/notifications/:user_id` | Get user notifications  |

### System Endpoints

| Method | Endpoint  | Description  |
| ------ | --------- | ------------ |
| GET    | `/health` | Health check |

---

## Project Structure

```
.
├── .gitignore
├── .prettierignore
├── .prettierrc
├── docker-compose.yml
├── docker-compose.scaling.yml
├── go.mod
├── go.sum
├── package-lock.json
├── package.json
├── README.md
├── cmd/
│   ├── api/
│   │   └── main.go              # API server entry point
│   └── worker/
│       └── main.go              # Background worker entry point
├── docker/
│   ├── Dockerfile.api           # API container image
│   ├── Dockerfile.worker       # Worker container image
│   └── nginx.conf               # Nginx load balancer config
├── internal/
│   ├── booking/
│   │   ├── command_handler.go   # Booking command handlers
│   │   ├── command_service.go  # Booking business logic
│   │   ├── projection.go        # Event projection
│   │   ├── query_handler.go     # Booking query handlers
│   │   └── query_service.go    # Booking query service
│   ├── config/
│   │   └── config.go            # Configuration
│   ├── db/
│   │   ├── command.go           # Command database connection
│   │   └── query.go            # Query database connection
│   ├── events/
│   │   ├── dispatcher.go       # Event dispatcher
│   │   ├── store.go            # Event store
│   │   └── types.go            # Event types
│   ├── middleware/
│   │   ├── auth.go             # Authentication middleware
│   │   ├── logger.go           # Logging middleware
│   │   └── rateLimiter.go      # Rate limiting middleware
│   ├── movie/
│   │   ├── command_handler.go  # Movie command handlers
│   │   ├── command_service.go  # Movie business logic
│   │   ├── model.go            # Movie model
│   │   ├── query_handler.go    # Movie query handlers
│   │   └── query_service.go    # Movie query service
│   ├── notification/
│   │   ├── command_handler.go
│   │   ├── command_service.go
│   │   ├── model.go
│   │   ├── query_handler.go
│   │   ├── query_service.go
│   │   ├── queue.go
│   │   ├── repository.go
│   │   └── worker.go
│   ├── payments/
│   │   ├── command_handler.go  # Payment command handlers
│   │   ├── command_service.go  # Payment business logic
│   │   ├── model.go            # Payment model
│   │   ├── query_handler.go    # Payment query handlers
│   │   ├── query_service.go    # Payment query service
│   │   └── repository.go       # Payment repository
│   ├── queue/
│   │   ├── jobs.go             # Job definitions
│   │   ├── redis.go            # Redis queue
│   │   └── stream.go           # Stream processing
│   ├── show/
│   │   ├── command_handler.go  # Show command handlers
│   │   ├── command_service.go  # Show business logic
│   │   ├── model.go            # Show model
│   │   ├── query_handler.go    # Show query handlers
│   │   └── query_service.go    # Show query service
│   ├── user/
│   │   ├── command_handler.go  # User command handlers
│   │   ├── command_service.go  # User business logic
│   │   ├── model.go            # User model
│   │   ├── query_handler.go    # User query handlers
│   │   └── query_service.go    # User query service
│   └── utils/
│       ├── cache.go            # Cache utilities
│       ├── lock.go             # Distributed lock
│       └── response.go         # Response helpers
├── migrations/
│   ├── Command-DB/
│   │   ├── 001_event.sql
│   │   ├── 002_reservation.sql
│   │   ├── 003_user.sql
│   │   ├── 004_movies.sql
│   │   ├── 005_shows.sql
│   │   ├── 006_payments.sql
│   │   └── 007_notifications.sql
│   └── Query-DB/
│       ├── 001_project.sql
│       ├── 002_reservation.sql
│       ├── 003_shows&movies.sql
│       ├── 004_payment.sql
│       ├── 005_notification.sql
│       └── 006_user.sql
├── k8s/
│   ├── api-deployment.yaml
│   ├── hpa.yaml
│   ├── ingress.yaml
│   └── namespace.yaml
├── docs/
│   ├── POSTMAN.md
│   ├── SCALING.md
│   ├── TEST.md
│   ├── postman-collection.json
│   └── README.md
└── postgres-data/
```

---

## Database Schema

### Command Database (Write)

```
sql
-- Events table for event sourcing
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_aggregate ON events(aggregate_id);
CREATE INDEX idx_events_type ON events(event_type);
```

### Query Database (Read)

Optimized for fast queries:

- `reservations` - User bookings
- `shows` & `movies` - Showtime data
- `payments` - Payment records
- `notifications` - User alerts

---

## Testing

### Test Commands

Run all tests using the following commands:

The project includes comprehensive unit tests following the testing strategy outlined in [TEST.md](./docs/TEST.md).

```
bash
# Run all unit tests
go test ./internal/... -v

# Run with coverage
go test ./internal/... -coverprofile=coverage.out -covermode=atomic

# Run specific package tests
go test ./internal/booking/... -v
go test ./internal/user/... -v
go test ./internal/payments/... -v
go test ./internal/movie/... -v
go test ./internal/show/... -v
go test ./internal/notification/... -v
go test ./internal/events/... -v
go test ./internal/middleware/... -v
go test ./internal/queue/... -v
go test ./internal/utils/... -v

# Run integration tests (requires Docker)
go test ./internal/test/integration/... -v

# Skip integration tests (for CI/CD)
go test ./internal/... -short -v
```

### Test Types

- **Unit Tests**: Located in each package (e.g., `booking_test.go`, `user_test.go`)
- **Integration Tests**: Located in `internal/test/integration/`
- **Model Tests**: Testing struct definitions and field values
- **Validation Tests**: Testing input validation logic

### Test Coverage

| Package      | Tests                                | Status  |
| ------------ | ------------------------------------ | ------- |
| booking      | 5+ test cases                        | ✅ PASS |
| user         | 4+ test cases                        | ✅ PASS |
| movie        | 4+ test cases                        | ✅ PASS |
| show         | 5+ test cases                        | ✅ PASS |
| payments     | 7+ test cases                        | ✅ PASS |
| notification | 3+ test cases                        | ✅ PASS |
| events       | 5+ test cases                        | ✅ PASS |
| middleware   | 7+ test cases                        | ✅ PASS |
| queue        | 4+ test cases                        | ✅ PASS |
| utils        | 4+ test cases                        | ✅ PASS |

### Using Postman

1. Import `docs/postman-collection.json` into Postman
2. Set the base URL to `http://localhost:8081`
3. Follow the testing workflow in [Postman Documentation](./docs/POSTMAN.md)

### Example Workflow

```
bash
# 1. Register a new user
POST /cmd/users/register
{"username": "john", "email": "john@example.com", "password": "pass123", "is_admin": false}

# 2. List movies
GET /query/movies

# 3. Create a show (admin)
POST /cmd/shows
{"movie_id": 1, "theater": "Theater A", "start_time": "2024-01-20T14:00:00Z"}

# 4. Check seat availability
GET /query/availability/A1

# 5. Reserve a ticket
POST /cmd/reserve
{"user_id": "1", "seat_id": "A1"}

# 6. Confirm the booking
POST /cmd/confirm
{"user_id": "1", "seat_id": "A1"}

# 7. Initiate payment
POST /cmd/payments/initiate
{"booking_id": "1", "user_id": "1", "amount": 500}

# 8. Verify payment
POST /cmd/payments/verify
{"payment_id": "pay_xxx", "mode": "success"}

# 9. Check user reservations
GET /query/reservations/1

# 10. Check notifications
GET /query/notifications/1
```

---

## Scaling

The system supports horizontal scaling with:

- **Load Balancing**: Nginx
- **Auto-scaling**: Kubernetes HPA
- **Session Management**: Redis
- **Rate Limiting**: Per-IP token bucket (5000 req/min)
- **Distributed Locking**: Redis-based locks for concurrency control

For detailed scaling documentation, see [SCALING.md](./docs/SCALING.md).

---

## License

MIT License - Feel free to use this project for learning and development.

---

---

<div align="center">

**Built with Go + PostgreSQL + Redis**

</div>
