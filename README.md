# 🎟️ TicketFair

**Event ticketing platform built with Go, Gin, GORM, CockroachDB, AWS S3 (LocalStack) and Stripe.**

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
- [Stripe Payments](#stripe-payments)
- [Observability](#observability)
- [Admin Dashboard](#admin-dashboard)
- [Testing](#testing)
- [Seed Data](#seed-data)

---

## Overview

TicketFair is a REST API for selling and managing event tickets. The platform supports three user types — **clients**, **merchants** (event producers) and **administrators** — each with their own authentication flow and permissions.

Payments are processed through **Stripe**, emulated locally via the **Stripe CLI** Docker container. Event images are stored in **AWS S3**, emulated via **LocalStack**.

### Key Features

- Registration and authentication for users, merchants and representatives
- Event creation and management with cover image upload to S3
- Ticket purchase via Stripe PaymentIntents with webhook-driven fulfillment
- Full refund flow through Stripe with automatic capacity restoration
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
| Image Storage | AWS S3 (LocalStack in development) |
| S3 SDK | aws-sdk-go-v2 |
| Payments | Stripe (stripe-go/v79) |
| Stripe local dev | Stripe CLI (Docker container) |
| Logging | slog + Loki + Promtail |
| Metrics | Prometheus + Grafana |
| Containerization | Docker + Docker Compose |

---

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                       Caddy                           │
│              (Reverse Proxy / TLS)                    │
└──────────────┬──────────────────┬────────────────────┘
               │                  │
       ┌───────▼──────┐   ┌───────▼──────┐
       │  TicketFair  │   │  Dashboard   │
       │     API      │   │   (admin)    │
       │   :8000      │   │   :3001      │
       └───────┬──────┘   └──────────────┘
               │
    ┌──────────┼──────────────────┐
    │          │                  │
    ▼          ▼                  ▼
CockroachDB  LocalStack        Stripe
  :26257      :4566 (S3)    (via CLI :4242)
                                  ▲
                         stripe-cli container
                         (forwards webhooks →
                          /public/webhooks/stripe)

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
configs/       Email (SMTP) + S3 client + Stripe init
```

---

## Project Structure

```
ticketfair/
├── configs/
│   ├── email.go
│   ├── s3.go
│   ├── stripe.go                      
│   ├── loki/loki-config.yml
│   ├── promtail/promtail-config.yml
│   └── prometheus/prometheus.yml
├── controllers/
│   ├── admin_controller.go
│   ├── auth_controller.go
│   ├── base_controller.go
│   ├── basics_controller.go
│   ├── events_controller.go
│   ├── image_controller.go
│   ├── merchant_controller.go
│   ├── merchant_rep_controller.go
│   ├── payment_controller.go                     
│   ├── profile_controller.go
│   ├── ticket_controller.go
│   ├── transaction_controller.go
│   ├── users_controller.go
│   └── verification_controller.go
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
│   ├── payment_dto.go                 
│   ├── profile_dto.go
│   ├── ticket_dto.go
│   ├── transaction_dto.go
│   ├── user_dto.go
│   ├── validators.go
│   └── verification_dto.go
├── middlewares/
│   ├── admin_middleware.go
│   ├── base_middleware.go
│   ├── client_middleware.go
│   ├── merchant_middleware.go
│   ├── merchant_rep_middleware.go
│   ├── public_middleware.go
│   └── rate_limiter_middleware.go
├── migration/
│   └── docker-database-init.sql      
├── models/
│   ├── address_model.go
│   ├── admin_model.go
│   ├── event_model.go
│   ├── merchant_model.go
│   ├── merchant_rep_model.go
│   ├── payment_model.go                     
│   ├── profile_model.go
│   ├── ticket_model.go
│   ├── transaction_model.go
│   ├── user_model.go
│   └── verification_model.go
├── routes/
│   └── routes.go
├── scripts/
│   └── localstack-init.sh
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
│   ├── payment_service.go             
│   ├── profile_service.go
│   ├── s3_service.go
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
- A Stripe account (free) — get test keys at https://dashboard.stripe.com/test/apikeys

### 1. Clone the repository

```bash
git clone https://github.com/maycolacerda/ticketFair
cd ticketFair
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` and at minimum set:

```bash
JWT_SECRET=some_long_random_string
STRIPE_SECRET_KEY=sk_test_xxxx   # from Stripe dashboard
```

### 3. Start the stack

```bash
docker compose up --build
```

### 4. Get the Stripe webhook secret

After the stack starts, the `stripe-cli` container begins forwarding webhooks. Grab the signing secret from its logs:

```bash
docker compose logs stripe-cli | grep "webhook signing secret"
# → Your webhook signing secret is whsec_xxxxxxxxxxxxxxxxxxxxx
```

Add it to `.env`:

```bash
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxxxxxxxxxx
```

Then restart the app:

```bash
docker compose restart ticketfair-app
```

### 5. Access

| Service | URL |
|---|---|
| API | http://localhost:8000 |
| Swagger | http://localhost:8000/swagger/index.html |
| Admin Dashboard | http://localhost:3001 |
| LocalStack / S3 | http://localhost:4566 |
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
JWT_SECRET=your_very_long_secret_key_change_this

# SMTP (optional — without this, emails are only logged)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your@email.com
SMTP_PASSWORD=your_app_password
SMTP_FROM=your@email.com
SMTP_FROM_NAME=TicketFair

# AWS S3 / LocalStack
AWS_ENDPOINT_URL=http://localstack:4566
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
S3_BUCKET=ticketfair-images

# Stripe
STRIPE_SECRET_KEY=sk_test_xxxxxxxxxxxxxxxxxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxxxxxxxxx  # from stripe-cli logs
STRIPE_DEVICE_NAME=ticketfair-local               # optional
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
POST /api/v1/public/webhooks/stripe      ← Stripe webhook receiver
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
POST /api/v1/private/tickets/purchase    ← legacy direct purchase (no Stripe)
POST /api/v1/private/tickets/refund
GET  /api/v1/private/transactions
GET  /api/v1/private/payments            ← list Stripe payments
POST /api/v1/private/payments/intent     ← create PaymentIntent
POST /api/v1/private/payments/:id/refund ← refund via Stripe
POST /api/v1/private/logout
```

### 🏪 Merchant (requires merchant token)

```
PUT    /api/v1/merchant/update
POST   /api/v1/merchant/events/new
PUT    /api/v1/merchant/events/:id
POST   /api/v1/merchant/events/:id/image
DELETE /api/v1/merchant/events/:id/image
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
payments         — Stripe PaymentIntent records
```

### SQL Functions (Atomicity)

```sql
create_profile_with_address(...)  -- Creates profile + address in one transaction
purchase_ticket(...)              -- Decrements capacity + creates transaction atomically
refund_ticket(...)                -- Restores capacity + marks transaction as refunded
```

### Payment status lifecycle

```
pending → succeeded → refunded
        ↘ failed
        ↘ canceled
```

### Ticket lifecycle

```
Purchase (via Stripe webhook)  →  active
Validate at door               →  used
Refund (via Stripe)            →  refunded
```

---

## Image Storage (S3 / LocalStack)

Event cover images are stored in AWS S3. In development, LocalStack emulates S3 locally.

### Upload an event image

```bash
curl -X POST http://localhost:8000/api/v1/merchant/events/<event_id>/image \
  -H "Authorization: Bearer <merchant_token>" \
  -F "image=@/path/to/photo.jpg"
```

### Image validation rules

| Rule | Value |
|---|---|
| Allowed formats | JPEG, PNG, WebP |
| Max size | 5 MB |
| Old image | Automatically deleted from S3 when replaced |

### Inspect the S3 bucket

```bash
aws --endpoint-url=http://localhost:4566 s3 ls s3://ticketfair-images/ --recursive
```

---

## Stripe Payments

### How it works

```
Client calls POST /private/payments/intent
             │
             ▼
  Stripe creates PaymentIntent
  (returns client_secret)
             │
             ▼
  Client confirms payment on their side
  (using Stripe.js or mobile SDK)
             │
             ▼
  Stripe sends webhook → /public/webhooks/stripe
  (forwarded by stripe-cli in development)
             │
             ▼
  payment_intent.succeeded handler:
    ├── calls purchase_ticket() SQL function
    ├── creates Ticket record
    ├── updates Payment.status = "succeeded"
    └── sends confirmation email (async)
```

### Stripe webhook events handled

| Event | Action |
|---|---|
| `payment_intent.succeeded` | Creates transaction + ticket, sends email |
| `payment_intent.payment_failed` | Marks payment as failed |
| `payment_intent.canceled` | Marks payment as canceled |
| `charge.refunded` | Restores capacity, marks ticket as refunded |

### Payment flow (step by step)

**1. Create a PaymentIntent:**

```bash
curl -X POST http://localhost:8000/api/v1/private/payments/intent \
  -H "Authorization: Bearer <client_token>" \
  -H "Content-Type: application/json" \
  -d '{ "event_id": "c3d4e5f6-a7b8-9012-cdef-123456789012", "quantity": 1 }'
```

Response:

```json
{
  "data": {
    "client_secret":      "pi_xxx_secret_yyy",
    "payment_intent_id":  "pi_xxxxxxxxxxxxxxxx",
    "amount":             5000,
    "currency":           "brl",
    "event_id":           "c3d4e5f6-...",
    "event_name":         "Summer Festival 2026",
    "quantity":           1
  }
}
```

**2. Confirm using Stripe test card** (for API testing — use Stripe CLI or your own client):

```bash
# Trigger a test payment_intent.succeeded event directly via Stripe CLI
docker compose exec stripe-cli stripe trigger payment_intent.succeeded
```

**3. Check payment status:**

```bash
curl http://localhost:8000/api/v1/private/payments \
  -H "Authorization: Bearer <client_token>"
```

**4. Issue a refund:**

```bash
curl -X POST http://localhost:8000/api/v1/private/payments/<payment_id>/refund \
  -H "Authorization: Bearer <client_token>"
```

### Stripe test cards

| Card number | Behavior |
|---|---|
| `4242 4242 4242 4242` | Payment succeeds |
| `4000 0000 0000 9995` | Payment declined (insufficient funds) |
| `4000 0025 0000 3155` | Requires 3D Secure authentication |

Use any future expiry date, any 3-digit CVC and any 5-digit postal code.

### Trigger webhook events manually

```bash
# Simulate a succeeded payment
docker compose exec stripe-cli stripe trigger payment_intent.succeeded

# Simulate a failed payment
docker compose exec stripe-cli stripe trigger payment_intent.payment_failed

# Simulate a refund
docker compose exec stripe-cli stripe trigger charge.refunded

# Watch live webhook events
docker compose logs -f stripe-cli
```

### Moving to production

1. Replace `STRIPE_SECRET_KEY=sk_test_...` with `sk_live_...`
2. Remove the `stripe-cli` service from `docker-compose.yaml`
3. Set up a real webhook endpoint in the Stripe dashboard pointing to `https://yourdomain.com/api/v1/public/webhooks/stripe`
4. Copy the production webhook signing secret into `STRIPE_WEBHOOK_SECRET`

### Idempotency

The webhook handler checks `payment.Status == "succeeded"` before processing. Stripe may deliver the same event more than once — duplicate processing is safe.

---

## Observability

### Structured Logging

In production (`GIN_MODE=release`) all logs are emitted as JSON and collected by Promtail for indexing in Loki.

```json
{
  "time":       "2026-03-29T12:00:00Z",
  "level":      "INFO",
  "msg":        "Payment succeeded — ticket created",
  "pi_id":      "pi_xxxxxxxxxxxxxxxx",
  "transaction_id": "uuid",
  "user_id":    "uuid",
  "event_id":   "uuid"
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
go test ./... -v
go test ./controllers/... -v
```

### Add the Stripe Go SDK

```bash
go get github.com/stripe/stripe-go/v79
```

### Postman Collection

Import `ticketfair.postman_collection.json` and run in this order:

```
1.  Admin Login
2.  Merchant Login (seeded account)
3.  Rep Login (seeded account)
4.  Register User
5.  Client Login
6.  Create Profile
7.  Send Email Verification → check logs → Verify Email
8.  Create Event (merchant)
9.  Upload Event Image
10. Create PaymentIntent          ← POST /private/payments/intent
11. Trigger webhook               ← docker compose exec stripe-cli stripe trigger payment_intent.succeeded
12. List My Payments              ← verify status = succeeded
13. List My Tickets               ← verify ticket was created
14. Validate Ticket (merchant)
15. Refund Payment                ← POST /private/payments/:id/refund
16. List My Payments              ← verify status = refunded
```

### Full Stripe test with curl

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8000/api/v1/public/auth/client/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"PassW0rd!"}' | jq -r '.data.token')

# 2. Create PaymentIntent
PI=$(curl -s -X POST http://localhost:8000/api/v1/private/payments/intent \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"event_id":"c3d4e5f6-a7b8-9012-cdef-123456789012","quantity":1}')

echo $PI | jq '.data.payment_intent_id'

# 3. Trigger succeeded webhook
docker compose exec stripe-cli stripe trigger payment_intent.succeeded

# 4. Check payment status
curl http://localhost:8000/api/v1/private/payments \
  -H "Authorization: Bearer $TOKEN" | jq '.data[0].status'
# → "succeeded"

# 5. Check ticket was created
curl http://localhost:8000/api/v1/private/tickets \
  -H "Authorization: Bearer $TOKEN" | jq '.data[0].status'
# → "active"
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
| Ticket Price | R$ 50.00 (5000 cents) |
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
git checkout -b feature/my-feature
git add .
git commit -m "feat: description of the feature"
git push origin feature/my-feature
```

---

## License

MIT — see [LICENSE](LICENSE) for details.