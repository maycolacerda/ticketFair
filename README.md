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
- [Production Guide](#production-guide)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

TicketFair is a ticket management platform designed to ensure secondary market integrity and combat predatory scalping. Powered by a robust CockroachDB backend, the system implements database-level security constraints to manage non-transferable gift tickets and price-capped resale. The infrastructure leverages Docker, Caddy, and a full observability stack (Loki/Prometheus/Grafana) to ensure high availability and traceability during peak demand.

The platform supports three user types — **clients**, **merchants** (event producers) and **administrators** — each with their own authentication flow and permissions. Payments are processed through **Stripe**, emulated locally via the **Stripe CLI** Docker container. Event images are stored in **AWS S3**, emulated via **LocalStack**.

### Key Features

- Registration and authentication for users, merchants and representatives
- Event creation with multiple ticket types (General, VIP, Early Bird, Reserved, Group, Day Pass, Tiered, Complimentary, Demographic)
- Per-ticket-type capacity management, sale windows and order limits
- Ticket purchase via Stripe PaymentIntents with webhook-driven fulfillment
- Full refund flow through Stripe with automatic capacity restoration
- Ticket validation at event entry
- Email and phone verification
- Password reset via 6-digit email code
- Admin panel with user and merchant management
- Per-IP rate limiting for brute force protection
- Structured logging with Loki + Grafana
- Event cover image upload to S3

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
│                       Caddy                          │
│              (Reverse Proxy / TLS)                   │
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
│   ├── password_reset_controller.go
│   ├── payment_controller.go
│   ├── profile_controller.go
│   ├── ticket_controller.go
│   ├── ticket_type_controller.go
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
│   ├── password_reset_dto.go
│   ├── payment_dto.go
│   ├── profile_dto.go
│   ├── ticket_dto.go
│   ├── ticket_type_dto.go
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
│   ├── password_reset_model.go
│   ├── payment_model.go
│   ├── profile_model.go
│   ├── ticket_model.go
│   ├── ticket_type_model.go
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
│   ├── password_reset_service.go
│   ├── payment_service.go
│   ├── profile_service.go
│   ├── s3_service.go
│   ├── ticket_service.go
│   ├── ticket_type_service.go
│   ├── token.go
│   ├── transaction_service.go
│   ├── user_service.go
│   └── verification_service.go
├── tests/
│   ├── testutil/
│   │   ├── db.go
│   │   ├── engine.go
│   │   ├── fixtures.go
│   │   ├── request.go
│   │   └── router.go
│   ├── integration/
│   │   ├── admin_test.go
│   │   ├── auth_test.go
│   │   ├── event_test.go
│   │   ├── main_test.go
│   │   ├── profile_test.go
│   │   └── purchase_test.go
│   └── unit/
│       ├── basics_test.go
│       ├── token_test.go
│       ├── users_test.go
│       └── validators_test.go
├── .env
├── .env.example
├── .gitignore
├── Caddyfile
├── docker-compose.yaml
├── docker-compose.test.yml
├── Dockerfile
├── Dockerfile.test
├── Makefile
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
STRIPE_SECRET_KEY=sk_test_xxxx
```

### 3. Start the stack

```bash
docker compose up --build
```

### 4. Get the Stripe webhook secret

```bash
docker compose logs stripe-cli | grep "webhook signing secret"
# → Your webhook signing secret is whsec_xxxxxxxxxxxxxxxxxxxxx
```

Add it to `.env` and restart:

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
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxxxxxxxxx
STRIPE_DEVICE_NAME=ticketfair-local
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
POST /api/v1/public/auth/password/forgot
POST /api/v1/public/auth/password/reset
POST /api/v1/public/merchant/register
GET  /api/v1/public/events
GET  /api/v1/public/events/:id
GET  /api/v1/public/events/:id/ticket-types
POST /api/v1/public/webhooks/stripe
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
GET  /api/v1/private/payments
POST /api/v1/private/payments/intent
POST /api/v1/private/payments/:id/refund
POST /api/v1/private/logout
```

### 🏪 Merchant (requires merchant token)

```
PUT    /api/v1/merchant/update
POST   /api/v1/merchant/events/new
PUT    /api/v1/merchant/events/:id
POST   /api/v1/merchant/events/:id/image
DELETE /api/v1/merchant/events/:id/image
GET    /api/v1/merchant/events/:id/ticket-types
POST   /api/v1/merchant/events/:id/ticket-types
PUT    /api/v1/merchant/events/:id/ticket-types/:ttid
DELETE /api/v1/merchant/events/:id/ticket-types/:ttid
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
| Verification / Password Reset | 3 req | 1 every 60s | too many verification attempts |
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
events           — Events created by merchants
ticket_types     — Ticket tiers per event (GA, VIP, Early Bird, etc.)
transactions     — Ticket purchases
tickets          — Individual tickets linked to transactions
verifications    — Email/phone verification codes
password_resets  — Password reset codes
payments         — Stripe PaymentIntent records
```

### Ticket Type Categories

```
general        — Standard general admission
vip            — Premium access with perks
early_bird     — Discounted, time-limited
reserved       — Specific seat selection
group          — Discounted multi-ticket packs
day_pass       — Single or multi-day festival passes
tiered         — First/second release pricing
complimentary  — Free (sponsors, press, VIPs)
demographic    — Student / junior / senior pricing
```

### SQL Functions (Atomicity)

```sql
create_profile_with_address(...)  -- Creates profile + address in one transaction
purchase_ticket(...)              -- Validates sale window, decrements ticket_type.available
                                  -- + event.capacity + creates transaction atomically
refund_ticket(...)                -- Restores capacity on both ticket_type and event,
                                  -- marks transaction as refunded
```

### Payment Status Lifecycle

```
pending → succeeded → refunded
        ↘ failed
        ↘ canceled
```

### Ticket Lifecycle

```
Purchase (Stripe webhook)  →  active
Validate at door           →  used
Refund                     →  refunded
```

---

## Image Storage (S3 / LocalStack)

Event cover images are stored in AWS S3. In development, LocalStack emulates S3 locally.

```bash
# Upload
curl -X POST http://localhost:8000/api/v1/merchant/events/<event_id>/image \
  -H "Authorization: Bearer <merchant_token>" \
  -F "image=@photo.jpg"

# Inspect bucket
aws --endpoint-url=http://localhost:4566 s3 ls s3://ticketfair-images/ --recursive
```

| Rule | Value |
|---|---|
| Allowed formats | JPEG, PNG, WebP |
| Max size | 5 MB |
| Old image | Automatically deleted when replaced |

---

## Stripe Payments

### Flow

```
POST /private/payments/intent  →  Stripe creates PaymentIntent
                                  (returns client_secret)
                               ↓
Client confirms with card      →  Stripe fires webhook
                               ↓
/public/webhooks/stripe        →  payment_intent.succeeded
                               ↓
purchase_ticket() SQL fn       →  transaction + ticket created
                               ↓
Confirmation email sent        →  async goroutine
```

### Webhook Events Handled

| Event | Action |
|---|---|
| `payment_intent.succeeded` | Creates transaction + ticket, sends email |
| `payment_intent.payment_failed` | Marks payment as failed |
| `payment_intent.canceled` | Marks payment as canceled |
| `charge.refunded` | Restores capacity, marks ticket as refunded |

### Stripe Test Cards

| Card | Behavior |
|---|---|
| `4242 4242 4242 4242` | Success |
| `4000 0000 0000 9995` | Declined |
| `4000 0025 0000 3155` | 3D Secure required |

### Trigger Webhook Manually

```bash
docker compose exec stripe-cli stripe trigger payment_intent.succeeded
docker compose logs -f stripe-cli
```

---

## Observability

### Structured Logging

```json
{
  "time":    "2026-03-29T12:00:00Z",
  "level":   "INFO",
  "msg":     "Payment succeeded — ticket created",
  "pi_id":   "pi_xxxxxxxxxxxxxxxx",
  "user_id": "uuid"
}
```

### Grafana

```
URL:      http://localhost:3000
Username: admin
Password: admin
```

Add Loki at `http://loki:3100` as a data source.

---

## Admin Dashboard

**URL:** http://localhost:3001

```
Email:    admin@ticketfair.com
Password: PassW0rd!
```

Features: live stats overview, user and merchant management, event listing.

---

## Testing

```bash
make test-docker       # everything in Docker, zero local setup
make test-unit         # unit tests only, no DB required
make test-integration  # requires: make test-db first
make test-cover        # HTML coverage report
```

### Integration test coverage

| Suite | What it covers |
|---|---|
| `auth_test.go` | Register, login, disabled accounts, rep cascade, full password reset flow |
| `event_test.go` | Event CRUD, all 9 ticket type categories, capacity adjustment, delete guard |
| `purchase_test.go` | Full buy→validate→refund flow, sold out, sale window, max per order |
| `profile_test.go` | Profile CRUD, email/phone verification end-to-end |
| `admin_test.go` | User/merchant activate/deactivate, merchant→rep cascade, role access control |

---

## Seed Data

| Entity | Field | Value |
|---|---|---|
| Merchant | Email | contato@ticketfairprod.com |
| Merchant | Password | `PassW0rd!` |
| Merchant | ID | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| Rep | Email | carlos@ticketfairprod.com |
| Rep | Role | admin |
| Event | Name | Summer Festival 2026 |
| Event | ID | `c3d4e5f6-a7b8-9012-cdef-123456789012` |
| Event | Capacity | 1000 |
| Admin | Email | admin@ticketfair.com |
| Admin | Password | `PassW0rd!` |
| Admin | ID | `d4e5f6a7-b8c9-0123-defa-234567890123` |
| All | Password | `PassW0rd!` |

---

## Production Guide

This section covers everything that must change before deploying TicketFair to a production environment.

### 1. Secrets and Environment Variables

Never commit `.env` to version control. In production use a secrets manager.

| Variable | Action |
|---|---|
| `JWT_SECRET` | Generate a cryptographically random string of at least 64 characters: `openssl rand -hex 64` |
| `STRIPE_SECRET_KEY` | Replace `sk_test_...` with `sk_live_...` from the Stripe dashboard |
| `STRIPE_WEBHOOK_SECRET` | Create a production webhook endpoint in the Stripe dashboard and copy the live signing secret |
| `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` | Replace with real IAM credentials (use a dedicated IAM user with S3-only permissions) |
| `AWS_ENDPOINT_URL` | **Remove entirely** — this variable must not be set in production, or the SDK will point to LocalStack |
| `COCKROACH_USER` | Create a dedicated DB user instead of using root |
| `SMTP_PASSWORD` | Use an app password or dedicated SMTP account, never your personal login |
| `GIN_MODE` | Must be `release` |

### 2. Remove Development Services

In `docker-compose.yaml`, remove or comment out these services before deploying:

```yaml
# Remove for production:
stripe-cli:    # webhook forwarding — Stripe handles this via dashboard
localstack:    # S3 emulator — use real AWS S3
```

And remove the LocalStack init script mount from any remaining services.

### 3. Database

**CockroachDB in production:**

```bash
# Create a dedicated user (not root)
cockroach sql --insecure --execute="
  CREATE USER ticketfair_app WITH PASSWORD 'strong_password';
  GRANT ALL ON DATABASE ticketfair TO ticketfair_app;
"

# Enable TLS — do not run insecure in production
# Replace: start-single-node --insecure
# With:    start-single-node --certs-dir=/certs
```

Update your DSN to include the password and SSL mode:

```bash
DB_HOST=your-cockroach-host
DB_PORT=26257
COCKROACH_USER=ticketfair_app
COCKROACH_PASSWORD=strong_password
COCKROACH_DB=ticketfair
DB_SSLMODE=require
```

**For managed CockroachDB Cloud**, use the connection string from the dashboard directly and set `DB_SSLMODE=verify-full`.

**Backups:** Enable automatic backups:

```sql
ALTER DATABASE ticketfair SET SCHEDULE BACKUP = '0 2 * * *' INTO 's3://your-backup-bucket';
```

### 4. AWS S3

```bash
# Create the production bucket
aws s3 mb s3://ticketfair-images-prod --region us-east-1

# Apply a public-read policy for event images
aws s3api put-bucket-policy \
  --bucket ticketfair-images-prod \
  --policy file://scripts/s3-bucket-policy.json

# Apply CORS
aws s3api put-bucket-cors \
  --bucket ticketfair-images-prod \
  --cors-configuration file://scripts/s3-cors.json
```

Create an IAM user with only the permissions TicketFair needs:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["s3:PutObject", "s3:DeleteObject", "s3:GetObject"],
      "Resource": "arn:aws:s3:::ticketfair-images-prod/*"
    }
  ]
}
```

Optional but recommended: put **CloudFront** in front of the bucket for CDN caching and to avoid exposing the S3 URL directly. Update `buildImageURL()` in `services/s3_service.go` to return the CloudFront domain.

### 5. Stripe

**Switch to live mode:**

1. Replace `sk_test_...` with `sk_live_...` in `.env`
2. Remove the `stripe-cli` service from `docker-compose.yaml`
3. Go to [Stripe Dashboard → Webhooks](https://dashboard.stripe.com/webhooks) and add a new endpoint:
   - URL: `https://yourdomain.com/api/v1/public/webhooks/stripe`
   - Events to listen for: `payment_intent.succeeded`, `payment_intent.payment_failed`, `payment_intent.canceled`, `charge.refunded`
4. Copy the **Signing Secret** from the webhook endpoint and set it as `STRIPE_WEBHOOK_SECRET`

**Webhook raw body:** Ensure Gin does not consume the request body before your webhook handler reads it. Add the raw body middleware specifically to the webhook route (see the "Critical" section in the next steps document).

### 6. TLS / HTTPS

Caddy handles TLS automatically via Let's Encrypt. Update your `Caddyfile` to use real domains:

```caddy
ticketfair.yourdomain.com {
  reverse_proxy ticketfair-app:8000
}

dashboard.yourdomain.com {
  reverse_proxy dashboard:3001

  # Restrict dashboard to your IP or use basic auth
  basicauth {
    admin $2a$14$...  # bcrypt hash of your password
  }
}
```

For the certificate to be issued, your domain must point to the server's public IP and ports 80/443 must be open.

### 7. CORS

Add CORS headers to the Gin router if your clients are served from a different domain:

```go
// main.go or routes.go
import "github.com/gin-contrib/cors"

r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
}))
```

### 8. Rate Limiting at Scale

The current token bucket rate limiter is in-memory and per-instance. With multiple API replicas, each instance has its own bucket and the effective limit multiplies by the number of instances.

For multi-instance deployments, replace with a Redis-backed limiter:

```bash
go get github.com/go-redis/redis/v8
go get github.com/go-redis/redis_rate/v10
```

Add a Redis service to `docker-compose.yaml` and update `middlewares/rate_limiter.go` to use `redis_rate.Limiter`.

### 9. Email (SMTP)

For production volume, replace the plain `net/smtp` sender with a transactional provider:

| Provider | Notes |
|---|---|
| **SendGrid** | `go get github.com/sendgrid/sendgrid-go` — reliable deliverability, webhooks, analytics |
| **AWS SES** | Cheapest at scale, use existing AWS credentials |
| **Resend** | Simple API, good developer experience |
| **Mailgun** | Strong EU compliance |

Set `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME` and `SMTP_PASSWORD` to your provider's values. No code changes needed as long as the provider supports SMTP.

### 10. Observability

**Loki** — configure retention in `configs/loki/loki-config.yml`:

```yaml
limits_config:
  retention_period: 744h  # 31 days
```

**Prometheus** — add alerting rules for high error rates, slow response times and DB connection exhaustion. Connect to **Alertmanager** to send alerts via Slack or PagerDuty.

**Grafana** — create dashboards for:
- Request rate and latency per endpoint
- Payment success/failure rate
- Active tickets and capacity remaining per event
- DB connection pool usage

Lock down Grafana behind authentication before exposing it publicly.

### 11. Seed Data

The `docker-database-init.sql` seed data is for development only. In production:

1. Remove all `INSERT` statements from the init SQL file — they include known credentials
2. Create the first admin account manually via the DB console after first deploy:

```sql
INSERT INTO admins (name, email, password, active)
VALUES ('Super Admin', 'admin@yourdomain.com', '<bcrypt_hash>', true);
```

Generate the bcrypt hash:

```bash
go run -e 'import "golang.org/x/crypto/bcrypt"; h,_:=bcrypt.GenerateFromPassword([]byte("YourStrongPassword"),10); println(string(h))'
```

### 12. Resource Limits

Review the `deploy.resources` limits in `docker-compose.yaml` and tune for your server size. For a VPS with 4 GB RAM:

```yaml
ticketfair-app:
  deploy:
    resources:
      limits:
        cpus: "2"
        memory: 512M
      reservations:
        cpus: "0.5"
        memory: 256M
```

CockroachDB needs the most memory — allocate at least 1 GB for anything beyond light traffic.

### 13. Health Checks and Restart Policies

All services already have `restart: unless-stopped`. Add a process supervisor like **Watchtower** to automatically pull and redeploy updated Docker images:

```yaml
watchtower:
  image: containrrr/watchtower
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  command: --interval 300 --cleanup
```

### 14. Pre-Deploy Checklist

```
□ JWT_SECRET is at least 64 random characters
□ GIN_MODE=release
□ AWS_ENDPOINT_URL is removed (not set to LocalStack)
□ S3_BUCKET points to production bucket
□ STRIPE_SECRET_KEY starts with sk_live_
□ STRIPE_WEBHOOK_SECRET is from the production dashboard endpoint
□ Stripe CLI service removed from docker-compose.yaml
□ LocalStack service removed from docker-compose.yaml
□ CockroachDB runs with TLS (--certs-dir, not --insecure)
□ Database user is not root
□ Seed data INSERT statements removed from init SQL
□ Admin account created via DB console with strong password
□ Dashboard protected behind auth or IP allowlist
□ Grafana protected behind auth
□ CORS configured for your actual frontend domain
□ Loki retention period set
□ Prometheus alerting rules configured
□ All integration tests pass: make test-docker
□ DNS records point to your server
□ Ports 80 and 443 open in firewall
□ Caddy Caddyfile uses real domain names
```

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