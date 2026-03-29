# 🎟️ TicketFair

**Event ticketing platform built with Go, Gin, GORM and CockroachDB.**

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [API Routes](#api-routes)
- [Authentication](#authentication)
- [Rate Limiting](#rate-limiting)
- [Database](#database)
- [Observability](#observability)
- [Admin Dashboard](#admin-dashboard)
- [Testing](#testing)
- [Seed Data](#seed-data)

---

## Overview

TicketFair is a REST API for selling and managing event tickets. The platform supports three user types — **clients**, **merchants** (event producers) and **administrators** — each with their own authentication flow and permissions.

### Key Features

- Registration and authentication for users, merchants and representatives
- Event creation and management
- Ticket purchase and refund with atomic capacity control
- Ticket validation at event entry
- Email and phone verification
- Admin panel with user and merchant management
- Per-IP rate limiting for brute force protection
- Structured logging with Loki + Grafana

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23 |
| HTTP Framework | Gin |
| ORM | GORM |
| Database | CockroachDB (Postgres-compatible) |
| Authentication | JWT HS256 (golang-jwt/jwt v5) |
| Passwords | bcrypt (DefaultCost) |
| Validation | go-playground/validator v10 |
| Documentation | Swagger (swaggo) |
| Reverse Proxy | Caddy |
| Logging | slog + Loki + Promtail |
| Metrics | Prometheus + Grafana |
| Containerization | Docker + Docker Compose |

---

## Architecture

```
┌─────────────────────────────────────────────────┐
│                    Caddy                        │
│         (Reverse Proxy / TLS)                   │
└──────────┬──────────────────┬───────────────────┘
           │                  │
    ┌──────▼──────┐   ┌───────▼──────┐
    │  TicketFair │   │  Dashboard   │
    │     API     │   │    Admin     │
    │  :8000      │   │   :3001      │
    └──────┬──────┘   └──────────────┘
           │
    ┌──────▼──────┐
    │ CockroachDB │
    │   :26257    │
    └─────────────┘

Observability:
Promtail → Loki → Grafana
App      → Prometheus → Grafana
```

### Application Layers

```
controllers/   HTTP only — bind, validate, respond
services/      Business logic — no gin, no net/http
models/        GORM structs
dto/           Input/output shapes, custom validators
middlewares/   JWT, roles, rate limiting, logging
routes/        Route registration
database/      Connection and migration
configs/       Email (SMTP)
```

---

## Project Structure

```
ticketfair/
├── configs/
│   ├── email.go
│   ├── loki/loki-config.yml
│   ├── promtail/promtail-config.yml
│   └── prometheus/prometheus.yml
├── controllers/
│   ├── admin.go
│   ├── auth.go
│   ├── base.go
│   ├── basics.go
│   ├── events.go
│   ├── merchant.go
│   ├── merchant_rep.go
│   ├── profile.go
│   ├── ticket.go
│   ├── transaction.go
│   ├── users.go
│   └── verification.go
├── dashboard/
│   ├── Dockerfile
│   ├── entrypoint.sh
│   ├── index.html
│   └── nginx.conf
├── database/
│   └── db.go
├── dto/
│   ├── address_dto.go
│   ├── admin_dto.go
│   ├── auth_dto.go
│   ├── event_dto.go
│   ├── merchant_dto.go
│   ├── merchant_rep_dto.go
│   ├── profile_dto.go
│   ├── ticket_dto.go
│   ├── transaction_dto.go
│   ├── user_dto.go
│   ├── validators.go
│   └── verification_dto.go
├── middlewares/
│   ├── admin.go
│   ├── base.go
│   ├── client.go
│   ├── merchant.go
│   ├── merchant_rep.go
│   ├── public.go
│   └── rate_limiter.go
├── migration/
│   └── docker-database-init.sql
├── models/
│   ├── address.go
│   ├── admin.go
│   ├── event.go
│   ├── merchant.go
│   ├── merchant_rep.go
│   ├── profile.go
│   ├── ticket.go
│   ├── transaction.go
│   ├── user.go
│   └── verification.go
├── routes/
│   └── routes.go
├── services/
│   ├── admin_auth.go
│   ├── admin_service.go
│   ├── auth.go
│   ├── email_service.go
│   ├── email_templates.go
│   ├── errors.go
│   ├── event_service.go
│   ├── logger.go
│   ├── merchant_rep_auth.go
│   ├── merchant_rep_service.go
│   ├── merchant_service.go
│   ├── profile_service.go
│   ├── ticket_service.go
│   ├── token.go
│   ├── transaction_service.go
│   ├── user_service.go
│   └── verification_service.go
├── .env
├── .env.example
├── .gitignore
├── Caddyfile
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── main.go
```

---

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### 1. Clone the repository

```bash
git clone https://github.com/maycolacerda/ticketFair
cd ticketFair
```

### 2. Set up environment variables

```bash
cp .env.example .env
# Edit .env with your values
```

### 3. Start the stack

```bash
docker compose up --build
```

### 4. Generate Swagger docs (first time only)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go
docker compose up --build ticketfair-app
```

### 5. Access

| Service | URL |
|---|---|
| API | http://localhost:8000 |
| Swagger | http://localhost:8000/swagger/index.html |
| Admin Dashboard | http://localhost:3001 |
| CockroachDB UI | http://localhost:8081 |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |

---

## Environment Variables

```bash
# App
APP_VERSION=1.0.0
GIN_MODE=release

# Database (CockroachDB)
DB_HOST=ticketfair-db
DB_PORT=26257
COCKROACH_USER=root
COCKROACH_DB=ticketfair

# JWT
JWT_SECRET=your_very_long_secret_key

# SMTP (optional — without this, emails are only logged)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your@email.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=your@email.com
SMTP_FROM_NAME=TicketFair
```

---

## API Routes

### 🌐 Public

```
GET  /api/v1/public/health
POST /api/v1/public/auth/register
POST /api/v1/public/auth/client/login
POST /api/v1/public/auth/merchant/login
POST /api/v1/public/auth/rep/login
POST /api/v1/public/auth/logout
POST /api/v1/public/merchant/register
GET  /api/v1/public/events
GET  /api/v1/public/events/:id
```

### 🔒 Private (requires client token)

```
GET  /api/v1/private/users
GET  /api/v1/private/users/me
GET  /api/v1/private/users/:id
GET  /api/v1/private/profile
POST /api/v1/private/profile
PUT  /api/v1/private/profile
POST /api/v1/private/verify/email/send
POST /api/v1/private/verify/email
POST /api/v1/private/verify/phone/send
POST /api/v1/private/verify/phone
GET  /api/v1/private/tickets
GET  /api/v1/private/tickets/:id
POST /api/v1/private/tickets/purchase
POST /api/v1/private/tickets/refund
GET  /api/v1/private/transactions
POST /api/v1/private/logout
```

### 🏪 Merchant (requires merchant token)

```
PUT  /api/v1/merchant/update
POST /api/v1/merchant/events/new
PUT  /api/v1/merchant/events/:id
POST /api/v1/merchant/tickets/:id/validate
POST /api/v1/merchant/rep/new
PUT  /api/v1/merchant/rep/:id
POST /api/v1/merchant/logout
```

### 🔑 Admin (requires admin token)

```
POST /api/v1/admin/auth/login
GET  /api/v1/admin/users
POST /api/v1/admin/users
PUT  /api/v1/admin/users/:id
POST /api/v1/admin/users/:id/deactivate
POST /api/v1/admin/users/:id/activate
GET  /api/v1/admin/merchants
POST /api/v1/admin/merchants
PUT  /api/v1/admin/merchants/:id
POST /api/v1/admin/merchants/:id/deactivate
POST /api/v1/admin/merchants/:id/activate
```

---

## Authentication

All protected endpoints require a JWT in the header:

```
Authorization: Bearer <token>
```

### Available Roles

| Role | Access |
|---|---|
| `client` | `/private/*` routes |
| `merchant` | `/merchant/*` routes |
| `admin` / `manager` / `staff` | Merchant rep routes |
| `superadmin` | `/admin/*` routes |

### JWT Claims

```json
{
  "user_id":     "uuid",
  "role":        "client | merchant | admin | manager | staff | superadmin",
  "merchant_id": "uuid (reps only)",
  "exp":         1234567890,
  "iss":         "ticketfair"
}
```

---

## Rate Limiting

Per-IP protection using the token bucket algorithm.

| Endpoint | Burst | Refill Rate | Message |
|---|---|---|---|
| Login (all) | 5 req | 1 every 3s | too many login attempts |
| Register | 3 req | 1 every 10s | too many registration attempts |
| Verification | 3 req | 1 every 60s | too many verification attempts |
| Public events | 30 req | 10/s | too many requests |
| Admin | 20 req | 5/s | too many requests |

When the limit is hit:

```
HTTP 429 Too Many Requests
Retry-After: 60
{ "error": "too many login attempts — try again later" }
```

---

## Database

### Tables

```
admins           — Platform administrators
users            — Client users
merchants        — Event producers
merchant_reps    — Merchant representatives
profiles         — User profiles (1:1 with users)
addresses        — Profile addresses (1:1 with profiles)
events           — Events created by merchants
transactions     — Ticket purchases
tickets          — Tickets linked to transactions
verifications    — Email/phone verification codes
```

### SQL Functions (Atomicity)

```sql
create_profile_with_address(...)  -- Creates profile + address in one transaction
purchase_ticket(...)              -- Decrements capacity + creates transaction atomically
refund_ticket(...)                -- Restores capacity + marks transaction as refunded
```

### Ticket Lifecycle

```
Purchase  →  active
Validate  →  used
Refund    →  refunded
```

---

## Observability

### Structured Logging

In production (`GIN_MODE=release`) all logs are emitted as JSON and collected by Promtail for indexing in Loki.

```json
{
  "time":    "2026-03-29T12:00:00Z",
  "level":   "INFO",
  "msg":     "Client login successful",
  "user_id": "uuid"
}
```

### Grafana Access

```
URL:      http://localhost:3000
Username: admin
Password: admin
```

Add Loki as a data source at `http://loki:3100` to query logs.

---

## Admin Dashboard

Web interface for platform management.

**URL:** http://localhost:3001

**Login:**
```
Email:    admin@ticketfair.com
Password: PassW0rd!
```

**Features:**
- Overview with real-time statistics
- List, create, activate and deactivate users
- List, create, activate and deactivate merchants
- View all active events

---

## Testing

### Run Unit Tests

```bash
go test ./...

# Verbose output
go test ./... -v

# Specific package
go test ./controllers/... -v
```

### Postman Collection

Import `ticketfair.postman_collection.json` into Postman to test all endpoints with automated tests and pre-configured variables.

**Recommended order:**

```
1.  Admin Login
2.  Merchant Login (seeded account)
3.  Rep Login (seeded account)
4.  Register User
5.  Client Login
6.  Create Profile
7.  Send Email Verification → check logs → Verify Email
8.  Purchase Ticket (Summer Festival 2026)
9.  List My Tickets
10. Validate Ticket (merchant)
```

---

## Seed Data

The database is automatically initialized with test data on first run.

### Merchant

| Field | Value |
|---|---|
| Name | TicketFair Productions |
| Email | contato@ticketfairprod.com |
| Password | `PassW0rd!` |
| ID | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |

### Merchant Rep (Admin)

| Field | Value |
|---|---|
| Name | Carlos Admin |
| Email | carlos@ticketfairprod.com |
| Password | `PassW0rd!` |
| Role | admin |

### Event

| Field | Value |
|---|---|
| Name | Summer Festival 2026 |
| Location | Cianorte Exhibition Park — PR, Brazil |
| Date | Dec 20, 2026 at 6:00 PM UTC |
| Capacity | 1000 |
| ID | `c3d4e5f6-a7b8-9012-cdef-123456789012` |

### Admin

| Field | Value |
|---|---|
| Email | admin@ticketfair.com |
| Password | `PassW0rd!` |
| ID | `d4e5f6a7-b8c9-0123-defa-234567890123` |

---

## Contributing

```bash
# 1. Create a branch
git checkout -b feature/my-feature

# 2. Make your changes and commit
git add .
git commit -m "feat: description of the feature"

# 3. Push and open a PR
git push origin feature/my-feature
```

---

## License

MIT — see [LICENSE](LICENSE) for details.
