# 🎬 Ticket-Booking System

A high-performance backend for ticket booking applications, inspired by platforms like BookMyShow. Designed to handle 50K+ concurrent users with robust concurrency control and scalability.

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
## 🚀 System Design Technologies

### 🏗 Backend Architecture

| Technology | Description |
|------------|-------------|
| **Golang (Go)** | High-performance backend service with native concurrency support |
| **RESTful API** | API layer for booking operations using Gin web framework |
| **Layered Architecture** | Clean separation between Controllers, Services, and Repository layers |

```
cmd/
├── api/main.go          # API server entry point
└── worker/main.go       # Background job worker

internal/
├── booking/            # Booking service (command side)
│   ├── handler.go      # HTTP handlers
│   ├── projection.go   # Read models (query side)
│   └── service.go     # Business logic
├── movie/             # Movie service
├── show/               # Showtime service
├── user/               # User & authentication
└── middleware/         # Auth, logging, rate limiting
```

---

### 🔀 CQRS Pattern

> **Command Query Responsibility Segregation**

```
┌─────────────────────────────────────────────────────────────┐
│                        CQRS Architecture                    │
├──────────────────────────┬──────────────────────────────────┤
│     COMMAND SIDE         │         QUERY SIDE               │
│  (Write Operations)      │    (Read Operations)             │
├──────────────────────────┼──────────────────────────────────┤
│ • Seat Booking           │ • Fetch Shows                    │
│ • Payment Processing    │ • Seat Availability              │
│ • Reservations           │ • User Bookings                  │
│ • Event Sourcing         │ • Projections                    │
└──────────────────────────┴──────────────────────────────────┘
```

- **Command Side**: Handles seat booking, payment, reservation updates
- **Query Side**: Optimized for fetching shows, seat availability, user bookings
- **Benefits**: Improves performance under heavy concurrent traffic

---

### 🗄 Database & Storage

| Feature | Implementation |
|---------|----------------|
| **PostgreSQL 15** | Relational DB for shows, seats, bookings |
| **ACID Transactions** | Ensures safe booking without double reservation |
| **Database Indexing** | Faster show/search queries |
| **GORM + pgx** | ORM and native PostgreSQL driver |

```
sql
-- Sample: Events table for event sourcing
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

---

### 🔒 Concurrency & Seat Locking

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
│ │ (10 min)│        │ (FAIL!) │        │ (OK!)   │           │
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

- **Seat Locking**: Prevents double booking with configurable timeout
- **Optimistic Locking**: Version-based conflict detection
- **Pessimistic Locking**: Database-level row locking
- **Go Concurrency**: Goroutines & Mutex Locks for thread-safe operations

---

### ⚡ Caching Layer

> **Redis** for real-time seat availability and temporary locks

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

- **Real-time seat availability**
- **Temporary seat locks**
- **Session management**
- **Rate limiting**
- Reduces database load significantly

---

### 📩 Job Queues & Background Processing

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

**Future Scope:**
- RabbitMQ / Kafka / Redis Streams

---

### 📡 Event-Driven System Design

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
- `ReservationCreated` - Initial seat reservation
- `SeatLocked` - Seat temporarily locked
- `BookingConfirmed` - Payment successful
- `SeatReleased` - Lock expired or cancelled

---

### 🔐 Authentication & Security

| Feature | Implementation |
|---------|----------------|
| **JWT Authentication** | Secure session handling with token-based auth |
| **Role-Based Access Control (RBAC)** | Admin & User roles |
| **Rate Limiting** | Prevent API abuse |
| **Structured Logging** | Debugging & tracing |

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

---

### 📊 Scalability Concepts

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

### 🐳 DevOps & Deployment

```
yaml
# docker-compose.yml
services:
  postgres:
    image: postgres:15
    ports: [5433:5432]
    
  redis:
    image: redis:7-alpine
    ports: [6379:6379]
    
  api:
    build: ./docker/DockerFile.api
    ports: [8081:8081]
    
  worker:
    build: ./docker/DockerFile.worker
```

| Tool | Purpose |
|------|---------|
| **Docker** | Containerized backend |
| **Docker Compose** | Multi-service setup (DB + Redis + API) |
| **GitHub Actions** | CI/CD Pipeline Ready |

---

### 📈 Monitoring & Logging

| Feature | Status |
|---------|--------|
| **Structured Logging** | ✅ Implemented |
| **Request Logging** | ✅ Implemented |
| **Error Tracing** | ✅ Implemented |
| **Prometheus + Grafana** | 🔜 Future Scope |

---

### 🧪 Testing & Reliability

| Type | Description |
|------|-------------|
| **Unit Testing** | Go testing package |
| **Integration Testing** | Booking workflow tests |
| **API Testing** | Postman collection ready |


---

## 📁 Repository Structure

```
.
├── .gitignore
├── .prettierignore
├── .prettierrc
├── docker-compose.yml
├── go.mod
├── go.sum
├── package-lock.json
├── package.json
├── READMe.Md
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── docker/
│   ├── DockerFile.api
│   └── DockerFile.worker
├── internal/
│   ├── booking/
│   │   ├── handler.go
│   │   ├── projection.go
│   │   └── service.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── postgres.go
│   ├── events/
│   │   ├── model.go
│   │   └── store.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── rateLimiter.go
│   ├── movie/
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── service.go
│   ├── queue/
│   │   ├── jobs.go
│   │   ├── redis.go
│   │   └── stream.go
│   ├── show/
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── service.go
│   ├── user/
│   │   ├── handler.go
│   │   └── model.go
│   └── utils/
│       └── response.go
├── migrations/
│   ├── 001_events.sql
│   └── 002_users.sql
└── postgres-data/
```

---

## 🛠 Quick Start

### Prerequisites
- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 15
- Redis 7

### Run with Docker Compose

## Database Access

To access the database using psql:

```bash
docker exec -it ticket-booking-db psql -U ticket -d ticket_db
```
```
bash
# Clone the repository
git clone https://github.com/hitorii/ticket-booking.git
cd ticket-booking

# Start all services
docker-compose up -d

# Access API
http://localhost:8081

# Access database
docker exec -it ticket-booking-db psql -U ticket -d ticket_db
```

### Run Locally

```
bash
# Install dependencies
go mod download

# Run API server
go run cmd/api/main.go

# Run background worker
go run cmd/worker/main.go
```

---

## 📝 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | User login |
| GET | `/api/movies` | List all movies |
| GET | `/api/shows/:movieId` | Get showtimes |
| POST | `/api/booking/reserve` | Reserve seats |
| POST | `/api/booking/confirm` | Confirm booking |
| GET | `/api/booking/user/:id` | Get user bookings |

---

## 🔧 Technology Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.25 |
| **Framework** | Gin |
| **Database** | PostgreSQL 15 |
| **Cache** | Redis 7 |
| **ORM** | GORM |
| **Auth** | JWT |
| **Container** | Docker |

---

## 📄 License

MIT License - Feel free to use this project for learning and development.

---

<div align="center">

**Built with ❤️ using Go + PostgreSQL + Redis**

</div>
