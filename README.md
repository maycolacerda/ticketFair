# 🎟️ TicketFair

**Event ticketing platform built with Go, Gin, GORM, CockroachDB and AWS S3 (LocalStack).**

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
- [Image Storage (S3 / LocalStack)](#image-storage-s3--localstack)
- [Observability](#observability)
- [Admin Dashboard](#admin-dashboard)
- [User Frontend](#user-frontend)
- [Testing](#testing)
- [Seed Data](#seed-data)

---

## Overview

TicketFair is a REST API for selling and managing event tickets. The platform supports three user types — **clients**, **merchants** (event producers) and **administrators** — each with their own authentication flow and permissions.

Event images are stored in **AWS S3**, emulated locally via **LocalStack**, so the full cloud storage workflow runs entirely on your machine with no AWS account needed in development.

### Key Features

- Registration and authentication for users, merchants and representatives
- Event creation and management with cover image upload to S3
- Ticket purchase and refund with atomic capacity control
- Ticket validation at event entry
- Email and phone verification
- Admin panel with user and merchant management
- Per-IP rate limiting for brute force protection
- Structured logging with Loki + Grafana
- User-facing frontend for browsing events and buying tickets
- Admin dashboard for platform management

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
| Image Storage | AWS S3 (LocalStack in development) |
| S3 SDK | aws-sdk-go-v2 |
| Logging | slog + Loki + Promtail |
| Metrics | Prometheus + Grafana |
| Containerization | Docker + Docker Compose |

---

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                        Caddy                             │
│              (Reverse Proxy / TLS)                       │
└─────┬──────────────┬───────────────┬────────────────────┘
      │              │               │
┌─────▼──────┐ ┌─────▼──────┐ ┌─────▼──────┐
│ TicketFair │ │  Frontend  │ │ Dashboard  │
│    API     │ │  (users)   │ │  (admin)   │
│  :8000     │ │  :3002     │ │  :3001     │
└─────┬──────┘ └────────────┘ └────────────┘
      │
 ┌────┴────────────┐
 │                 │
 ▼                 ▼
CockroachDB     LocalStack
 :26257           :4566
                  (S3)

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
configs/       Email (SMTP) + S3 client
```

---

## Project Structure

```
ticketfair/
├── configs/
│   ├── email.go
│   ├── s3.go                          ← S3/LocalStack client init
│   ├── loki/loki-config.yml
│   ├── promtail/promtail-config.yml
│   └── prometheus/prometheus.yml
├── controllers/
│   ├── admin.go
│   ├── auth.go
│   ├── base.go
│   ├── basics.go
│   ├── events.go
│   ├── image.go                       ← image upload/delete handlers
│   ├── merchant.go
│   ├── merchant_rep.go
│   ├── profile.go
│   ├── ticket.go
│   ├── transaction.go
│   ├── users.go
│   └── verification.go
├── dashboard/                         ← admin SPA
│   ├── Dockerfile
│   ├── entrypoint.sh
│   ├── index.html
│   └── nginx.conf
├── frontend/                          ← user-facing SPA
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
│   ├── event_dto.go                   ← includes image_url field
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
│   └── docker-database-init.sql      ← includes image_url column
├── models/
│   ├── address.go
│   ├── admin.go
│   ├── event.go                       ← includes ImageURL field
│   ├── merchant.go
│   ├── merchant_rep.go
│   ├── profile.go
│   ├── ticket.go
│   ├── transaction.go
│   ├── user.go
│   └── verification.go
├── routes/
│   └── routes.go
├── scripts/
│   └── localstack-init.sh            ← creates S3 bucket on startup
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
│   ├── s3_service.go                 ← upload, delete, presign
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
# Edit .env — at minimum set JWT_SECRET
```

### 3. Start the stack

```bash
docker compose up --build
```

Docker Compose will:

1. Start CockroachDB and run `docker-database-init.sql`
2. Start LocalStack and create the `ticketfair-images` S3 bucket
3. Start the Go API (waits for DB + LocalStack to be healthy)
4. Start the user frontend and admin dashboard
5. Start Caddy, Loki, Promtail, Prometheus and Grafana

### 4. Generate Swagger docs (first time only)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go
docker compose up --build ticketfair-app
```

### 5. Access

| Service | Direct URL | Via Caddy |
|---|---|---|
| API | http://localhost:8000 | http://ticketfair.localhost |
| Swagger | http://localhost:8000/swagger/index.html | — |
| User Frontend | http://localhost:3002 | http://app.localhost |
| Admin Dashboard | http://localhost:3001 | http://dashboard.localhost |
| LocalStack / S3 | http://localhost:4566 | — |
| CockroachDB UI | http://localhost:8081 | — |
| Grafana | http://localhost:3000 | http://grafana.localhost |
| Prometheus | http://localhost:9090 | — |

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
JWT_SECRET=your_very_long_secret_key_change_this

# SMTP (optional — without this, emails are only logged)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your@email.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=your@email.com
SMTP_FROM_NAME=TicketFair

# AWS S3 / LocalStack
# Development (LocalStack):
AWS_ENDPOINT_URL=http://localstack:4566
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
S3_BUCKET=ticketfair-images

# Production (real AWS — remove AWS_ENDPOINT_URL):
# AWS_REGION=us-east-1
# AWS_ACCESS_KEY_ID=AKIAxxxxxxxxxxxxxxxxx
# AWS_SECRET_ACCESS_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
# S3_BUCKET=ticketfair-images-prod
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
PUT    /api/v1/merchant/update
POST   /api/v1/merchant/events/new
PUT    /api/v1/merchant/events/:id
POST   /api/v1/merchant/events/:id/image    ← upload cover image (multipart)
DELETE /api/v1/merchant/events/:id/image    ← remove cover image
POST   /api/v1/merchant/tickets/:id/validate
POST   /api/v1/merchant/rep/new
PUT    /api/v1/merchant/rep/:id
POST   /api/v1/merchant/logout
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
events           — Events created by merchants (includes image_url)
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

## Image Storage (S3 / LocalStack)

Event cover images are stored in AWS S3. In development, [LocalStack](https://localstack.cloud) emulates the S3 service locally — no AWS account or credentials required.

### How it works

```
Merchant uploads image
        │
        ▼
POST /merchant/events/:id/image
(multipart/form-data, field: "image")
        │
        ▼
Validation (JPEG/PNG/WebP, max 5MB, real image check)
        │
        ▼
Upload to S3: events/<uuid>.<ext>
        │
        ▼
Save URL to events.image_url in DB
        │
        ▼
URL returned in all EventResponse payloads
        │
        ▼
Frontend renders image from S3 URL
```

### Upload an event image

```bash
curl -X POST http://localhost:8000/api/v1/merchant/events/<event_id>/image \
  -H "Authorization: Bearer <merchant_token>" \
  -F "image=@/path/to/photo.jpg"
```

Response:

```json
{
  "message": "image uploaded successfully",
  "image_url": "http://localhost:4566/ticketfair-images/events/abc123.jpg"
}
```

### Delete an event image

```bash
curl -X DELETE http://localhost:8000/api/v1/merchant/events/<event_id>/image \
  -H "Authorization: Bearer <merchant_token>"
```

### Image validation rules

| Rule | Value |
|---|---|
| Allowed formats | JPEG, PNG, WebP |
| Max size | 5 MB |
| Validation | Real image decode check (not just MIME sniffing) |
| Key format | `events/<uuid>.<ext>` |
| Old image | Automatically deleted from S3 when replaced |

### S3 bucket

| Setting | Value |
|---|---|
| Bucket name | `ticketfair-images` |
| Region | `us-east-1` |
| Access | Public read (images are publicly accessible) |
| LocalStack endpoint | `http://localhost:4566` |

### LocalStack setup

LocalStack starts automatically via Docker Compose. The `scripts/localstack-init.sh` script runs on startup and creates the bucket with public read access and CORS headers.

To inspect the bucket manually:

```bash
# List bucket contents
aws --endpoint-url=http://localhost:4566 s3 ls s3://ticketfair-images/ --recursive

# Upload a test file
aws --endpoint-url=http://localhost:4566 s3 cp test.jpg s3://ticketfair-images/test.jpg

# Or use awslocal (LocalStack CLI wrapper)
awslocal s3 ls s3://ticketfair-images/
```

### Moving to production (real AWS S3)

1. Remove `AWS_ENDPOINT_URL` from `.env`
2. Set real `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
3. Create the bucket in AWS console and set the bucket policy to allow public reads
4. Optionally put CloudFront in front for CDN caching

---

## Observability

### Structured Logging

In production (`GIN_MODE=release`) all logs are emitted as JSON and collected by Promtail for indexing in Loki.

```json
{
  "time":    "2026-03-29T12:00:00Z",
  "level":   "INFO",
  "msg":     "Image uploaded",
  "key":     "events/abc123.jpg",
  "url":     "http://localhost:4566/ticketfair-images/events/abc123.jpg"
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

## Admin Dashboard (not publicly avaiable yet)

Web interface for platform management.

**URL:** http://localhost:3001

**Login:**
```
Email:    admin@ticketfair.com
Password: PassW0rd!
```

**Features:**
- Overview with real-time statistics (users, merchants, events, capacity)
- List, create, activate and deactivate users
- List, create, activate and deactivate merchants
- View all active events

---

## User Frontend (Not publicly avaiable yet)

Public-facing web application for browsing events and buying tickets.

**URL:** http://localhost:3002

**Features:**

| Feature | Description |
|---|---|
| Event listing | Browse all upcoming events with cover images from S3 |
| Search | Filter events by name, location or description |
| Event detail | View full details, capacity and buy tickets inline |
| Authentication | Register and sign in without leaving the page |
| Ticket purchase | Enter amount and buy a ticket directly |
| My Tickets | View all purchased tickets and their status |
| Responsive | Works on mobile, tablet and desktop |

**Design:** Editorial serif aesthetic — warm cream palette, Fraunces display font, Cabinet Grotesk for UI text.

### User flow

```
1. Browse events on the homepage
2. Click an event card to open the detail modal
3. Sign in or register (inline, no page reload)
4. Enter purchase amount and click "Buy Now"
5. View purchased tickets under "My Tickets" in the nav
```

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
8.  Create Event
9.  Upload Event Image       ← POST /merchant/events/:id/image (multipart)
10. List Events              ← verify image_url is present
11. Purchase Ticket
12. List My Tickets
13. Validate Ticket (merchant)
14. Delete Event Image       ← DELETE /merchant/events/:id/image
```

### Test image upload with curl

```bash
# 1. Login as merchant
TOKEN=$(curl -s -X POST http://localhost:8000/api/v1/public/auth/merchant/login \
  -H "Content-Type: application/json" \
  -d '{"email":"contato@ticketfairprod.com","password":"PassW0rd!"}' \
  | jq -r '.data.token')

# 2. Upload image to seeded event
curl -X POST \
  http://localhost:8000/api/v1/merchant/events/c3d4e5f6-a7b8-9012-cdef-123456789012/image \
  -H "Authorization: Bearer $TOKEN" \
  -F "image=@./my-event.jpg"
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
| Image | Upload via `POST /merchant/events/:id/image` after first boot |
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