# 🎟️ TicketFair

**Event ticketing platform built with Go, Gin, GORM, CockroachDB, AWS S3, Stripe and React.**

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
- [Connections & Gift Tickets](#connections--gift-tickets)
- [Image Storage (S3 / LocalStack)](#image-storage-s3--localstack)
- [Stripe Payments](#stripe-payments)
- [Observability](#observability)
- [Frontend (React)](#frontend-react)
- [Admin Dashboard](#admin-dashboard)
- [Testing](#testing)
- [Seed Data](#seed-data)
- [Production Guide](#production-guide)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

TicketFair is a ticket management platform designed to ensure secondary market integrity and combat predatory scalping. Powered by a robust CockroachDB backend, the system implements database-level security constraints to manage non-transferable gift tickets and price-capped resale. The infrastructure leverages Docker, Caddy, and a full observability stack (Loki/Prometheus/Grafana) to ensure high availability and traceability during peak demand.

The platform supports three user types — **clients**, **merchants** (event producers) and **administrators** — each with their own authentication flow and permissions. Payments are processed through **Stripe**, emulated locally via the **Stripe CLI** Docker container. Event images are stored in **AWS S3**, emulated via **LocalStack**. A React/Vite frontend provides a complete user-facing interface covering every feature.

### Key Features

- Registration and authentication for users, merchants and representatives
- Event creation with multiple ticket types (General, VIP, Early Bird, Reserved, Group, Day Pass, Tiered, Complimentary, Demographic)
- Per-ticket-type capacity management, sale windows and per-order limits
- Connections system — send, accept and manage friend connections
- Gift tickets to accepted connections with an optional message
- Friends feed — see which events your connections are attending
- Ticket purchase via Stripe PaymentIntents with webhook-driven fulfillment
- Full refund flow through Stripe with automatic capacity restoration
- Ticket validation at event entry (merchant door scan)
- Email and phone verification with 6-digit codes
- Password reset via 6-digit email code
- Admin panel with full user and merchant lifecycle management
- Per-IP rate limiting for brute force protection
- Structured JSON logging with Loki + Grafana
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
| Frontend | React 18 + Vite |
| Frontend UI | Custom design system (no component library) |

---

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                          Caddy                               │
│                  (Reverse Proxy / TLS)                       │
└──────┬───────────────────────┬─────────────────┬─────────────┘
       │                       │                 │
┌──────▼──────┐  ┌─────────────▼───┐  ┌──────────▼──────┐
│ TicketFair  │  │    Dashboard    │  │    React UI     │
│    API      │  │    (admin SPA)  │  │  (user-facing)  │
│   :8000     │  │     :3001       │  │     :3002       │
└──────┬──────┘  └─────────────────┘  └─────────────────┘
       │
  ┌────┴──────────────────────┐
  │                           │                    │
  ▼                           ▼                    ▼
CockroachDB              LocalStack            Stripe
  :26257                  :4566 (S3)        (via CLI)
                                                   ▲
                                          stripe-cli container
                                          (forwards webhooks →
                                           /public/webhooks/stripe)

Observability:
Promtail → Loki (:3100) → Grafana (:3000)
App      → Prometheus (:9090) → Grafana
```

### Application Layers (API)

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
│   ├── connection_controller.go
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
├── database/
│   └── db.go
├── dto/
│   ├── address_dto.go
│   ├── admin_dto.go
│   ├── auth_dto.go
│   ├── connection_dto.go
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
│   ├── connection_model.go
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
│   ├── connection_service.go
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
├── ticketfair-react/              ← React frontend
│   ├── src/
│   │   ├── api/client.js          ← all API calls in one place
│   │   ├── context/
│   │   │   ├── AuthContext.jsx
│   │   │   └── ToastContext.jsx
│   │   ├── components/
│   │   │   ├── ui.jsx             ← full design system primitives
│   │   │   └── Sidebar.jsx
│   │   ├── pages/
│   │   │   ├── AuthPage.jsx
│   │   │   ├── EventsPage.jsx
│   │   │   ├── TicketsPage.jsx    ← also exports PaymentsPage, FeedPage
│   │   │   ├── ConnectionsPage.jsx
│   │   │   ├── ProfilePage.jsx
│   │   │   ├── MerchantPage.jsx
│   │   │   └── AdminPage.jsx
│   │   ├── styles/global.css
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── vite.config.js
│   └── package.json
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
- Node.js 20+ (for frontend development only)
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

### 3. Start the backend stack

```bash
docker compose up --build
```

### 4. Get the Stripe webhook secret

After the stack starts, the `stripe-cli` container begins forwarding webhooks. Grab the signing secret from its logs:

```bash
docker compose logs stripe-cli | grep "webhook signing secret"
# → Your webhook signing secret is whsec_xxxxxxxxxxxxxxxxxxxxx
```

Add it to `.env` and restart the app:

```bash
docker compose restart ticketfair-app
```

### 5. Start the React frontend (development)

```bash
cd ticketfair-react
echo "VITE_API_URL=http://localhost:8000/api/v1" > .env.local
npm install
npm run dev
# → http://localhost:3002
```

### 6. Access

| Service | URL | Notes |
|---|---|---|
| React Frontend | http://localhost:3002 | Full user interface |
| API | http://localhost:8000 | REST API |
| Swagger | http://localhost:8000/swagger/index.html | API docs |
| Admin Dashboard | http://localhost:3001 | Admin-only SPA |
| LocalStack / S3 | http://localhost:4566 | Dev S3 emulator |
| CockroachDB UI | http://localhost:8081 | DB console |
| Grafana | http://localhost:3000 | Logs + metrics |
| Prometheus | http://localhost:9090 | Raw metrics |

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
POST /api/v1/private/tickets/:id/gift
GET  /api/v1/private/transactions
GET  /api/v1/private/payments
POST /api/v1/private/payments/intent
POST /api/v1/private/payments/:id/refund
GET  /api/v1/private/connections
POST /api/v1/private/connections
GET  /api/v1/private/connections/requests
GET  /api/v1/private/connections/events
POST /api/v1/private/connections/:id/respond
DELETE /api/v1/private/connections/:id
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
tickets          — Individual tickets (gift fields: is_gift, gifted_by, gifted_at, gift_message)
verifications    — Email/phone verification codes
password_resets  — Password reset codes
payments         — Stripe PaymentIntent records
connections      — User connections (requester_id, addressee_id, status)
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
create_profile_with_address(...)   -- Creates profile + address in one transaction
purchase_ticket(p_user_id, p_event_id, p_ticket_type_id, p_quantity, p_amount)
  -- Validates sale window, min/max per order, availability
  -- Decrements ticket_type.available + event.capacity atomically
  -- Creates transaction, returns transaction_id
refund_ticket(p_transaction_id)
  -- Restores ticket_type.available + event.capacity
  -- Marks transaction as refunded
```

### Payment Status Lifecycle

```
pending → succeeded → refunded
        ↘ failed
        ↘ canceled
```

### Ticket Lifecycle

```
Purchase (Stripe webhook)   →  active
Gift to connection          →  active (ownership transferred)
Validate at door            →  used
Refund (Stripe)             →  refunded
```

---

## Connections & Gift Tickets

### Connection Lifecycle

```
User A sends request  →  pending
User B responds:
  accept  →  accepted   ← can now gift tickets, see each other's events
  decline →  declined
Either party removes  →  soft deleted (connection_id preserved in audit)
```

### API Usage

```bash
# Send a connection request
POST /api/v1/private/connections
{ "addressee_id": "uuid-of-target-user" }

# See incoming pending requests
GET /api/v1/private/connections/requests

# Accept or decline
POST /api/v1/private/connections/:id/respond
{ "action": "accept" }   # or "decline"

# See events your connections are attending
GET /api/v1/private/connections/events

# Remove a connection
DELETE /api/v1/private/connections/:id
```

### Gift Tickets

Tickets can be transferred to accepted connections only. This is the platform's primary anti-scalping mechanism at the social layer — bulk resale is infeasible because it requires an established connection with every buyer.

```bash
POST /api/v1/private/tickets/:id/gift
{
  "recipient_id": "uuid-of-connection",
  "message": "Enjoy the show! 🎸"
}
```

**Rules enforced at the service layer:**

- Sender must own the ticket (`user_id = sender`)
- Ticket must have `status = active`
- Sender and recipient must have an `accepted` connection
- Sender cannot gift a ticket to themselves
- A ticket that has already been gifted (`is_gift = true`) cannot be gifted again
- Used or refunded tickets cannot be gifted

**Response includes:** `ticket_id`, `event_name`, `recipient_username`, `message`, `gifted_at`

---

## Image Storage (S3 / LocalStack)

Event cover images are stored in AWS S3. In development, LocalStack emulates S3 locally.

```bash
# Upload a cover image
curl -X POST http://localhost:8000/api/v1/merchant/events/<event_id>/image \
  -H "Authorization: Bearer <merchant_token>" \
  -F "image=@photo.jpg"

# Delete cover image
DELETE /api/v1/merchant/events/<event_id>/image

# Inspect the bucket locally
aws --endpoint-url=http://localhost:4566 s3 ls s3://ticketfair-images/ --recursive
```

| Rule | Value |
|---|---|
| Allowed formats | JPEG, PNG, WebP |
| Max size | 5 MB |
| Validation | Real image decode check (not just MIME header) |
| Old image | Automatically deleted from S3 when replaced |

---

## Stripe Payments

### Flow

```
POST /private/payments/intent       →  Stripe creates PaymentIntent
                                       (returns client_secret)
                                    ↓
Client confirms with test card      →  Stripe fires webhook
                                    ↓
/public/webhooks/stripe             →  payment_intent.succeeded
                                    ↓
purchase_ticket() SQL function      →  transaction + tickets created atomically
                                    ↓
Confirmation email sent             →  async goroutine
```

### Webhook Events Handled

| Event | Action |
|---|---|
| `payment_intent.succeeded` | Calls `purchase_ticket()`, creates ticket records, sends email |
| `payment_intent.payment_failed` | Marks payment as failed |
| `payment_intent.canceled` | Marks payment as canceled |
| `charge.refunded` | Calls `refund_ticket()`, restores capacity, marks tickets refunded |

### Test Cards

| Card number | Behavior |
|---|---|
| `4242 4242 4242 4242` | Payment succeeds |
| `4000 0000 0000 9995` | Declined — insufficient funds |
| `4000 0025 0000 3155` | Requires 3D Secure |

Use any future expiry, any 3-digit CVC, any 5-digit postal code.

### Trigger Webhook Events Manually

```bash
# Simulate a successful payment
docker compose exec stripe-cli stripe trigger payment_intent.succeeded

# Simulate a failed payment
docker compose exec stripe-cli stripe trigger payment_intent.payment_failed

# Simulate a refund
docker compose exec stripe-cli stripe trigger charge.refunded

# Watch events live
docker compose logs -f stripe-cli
```

### Idempotency

The webhook handler checks `payment.Status == "succeeded"` before processing. Stripe may deliver the same event more than once — duplicate processing is safe.

### Moving to Production

1. Replace `STRIPE_SECRET_KEY=sk_test_...` with `sk_live_...`
2. Remove the `stripe-cli` service from `docker-compose.yaml`
3. Create a webhook endpoint in the Stripe dashboard → `https://yourdomain.com/api/v1/public/webhooks/stripe`
4. Listen for: `payment_intent.succeeded`, `payment_intent.payment_failed`, `payment_intent.canceled`, `charge.refunded`
5. Copy the signing secret into `STRIPE_WEBHOOK_SECRET`

---

## Observability

### Structured Logging

In `GIN_MODE=release` all logs are emitted as JSON and indexed by Loki:

```json
{
  "time":           "2026-03-29T12:00:00Z",
  "level":          "INFO",
  "msg":            "Payment succeeded — ticket created",
  "pi_id":          "pi_xxxxxxxxxxxxxxxx",
  "transaction_id": "uuid",
  "user_id":        "uuid",
  "event_id":       "uuid"
}
```

### Grafana

```
URL:      http://localhost:3000
Username: admin
Password: admin
```

Add Loki as a data source at `http://loki:3100` to query application logs.

---

## Frontend (React)

A full React 18 + Vite SPA located in `ticketfair-react/`. Covers every platform feature with a dark, industrial concert-poster aesthetic (Bebas Neue + Outfit + JetBrains Mono).

### Running locally

```bash
cd ticketfair-react
echo "VITE_API_URL=http://localhost:8000/api/v1" > .env.local
npm install
npm run dev
# → http://localhost:3002
```

### Building for production

```bash
npm run build
# Output: ticketfair-react/dist/
```

### Docker

```bash
docker build \
  --build-arg VITE_API_URL=https://api.yourdomain.com/api/v1 \
  -t ticketfair-ui \
  ./ticketfair-react

docker run -p 3002:3002 ticketfair-ui
```

### Docker Compose integration

Add to `docker-compose.yaml`:

```yaml
frontend:
  build:
    context: ./ticketfair-react
    args:
      VITE_API_URL: http://localhost:8000/api/v1
  ports:
    - "3002:3002"
  networks:
    - ticketfair-network
  depends_on:
    ticketfair-app:
      condition: service_healthy
```

Add to `Caddyfile`:

```caddy
app.yourdomain.com {
  reverse_proxy frontend:3002
}
```

### Pages and features

| Page | Features |
|---|---|
| **Auth** | Login, register, forgot/reset password (email code, 2-step) |
| **Discover** | Event grid with search, cover images, ticket type selector with category badges, quantity picker, Stripe PaymentIntent creation |
| **My Tickets** | Status badges (active/used/refunded), gift ticket to a connection (modal with connection picker + message), refund button |
| **Payments** | Stripe payment history, status badges, refund via API |
| **Friends Feed** | Upcoming events that accepted connections have tickets for, attendee chips |
| **Connections** | Connected/pending tabs, send request by user ID, accept/decline, remove, pending badge on nav |
| **Profile** | View/edit profile and address, email + phone verification (inline code flow), change password (modal with email code) |
| **My Venue** (Merchant) | Create events, manage all 9 ticket type categories, live ticket validator by UUID |
| **Admin** | Create/activate/deactivate users and merchants in tabs |

### Architecture

```
src/
├── api/client.js          — single API module, handles client/merchant/admin tokens
├── context/
│   ├── AuthContext.jsx    — login, register, auto merchant detection, logout
│   └── ToastContext.jsx   — useToast() hook, renders toast stack
├── components/
│   ├── ui.jsx             — Btn, Badge, Modal, Field, Alert, Spinner, Empty,
│   │                        Skeleton, Card, Tabs, Stat, Avatar, Row, PanelSection,
│   │                        SearchInput, SectionHeader, Divider
│   └── Sidebar.jsx        — collapsible icon/label nav, pending badge, sign-out chip
├── pages/
│   ├── AuthPage.jsx
│   ├── EventsPage.jsx
│   ├── TicketsPage.jsx    — default export TicketsPage; named: PaymentsPage, FeedPage
│   ├── ConnectionsPage.jsx
│   ├── ProfilePage.jsx
│   ├── MerchantPage.jsx
│   └── AdminPage.jsx
├── styles/global.css      — CSS custom properties, animations, base resets
├── App.jsx                — shell, topbar, page switching, pending count polling
└── main.jsx               — entry point
```

---

## Admin Dashboard

A separate lightweight SPA at `dashboard/` for platform administrators only.

**URL:** http://localhost:3001

**Login:**

```
Email:    admin@ticketfair.com
Password: PassW0rd!
```

**Features:** live statistics overview, user management (list, create, activate, deactivate), merchant management, event listing.

---

## Testing

```bash
# All tests in Docker — zero local setup required
make test-docker

# Unit tests only (no DB needed)
make test-unit

# Integration tests (requires: make test-db first)
make test-db
make test-integration

# Coverage report (generates coverage.html)
make test-cover
```

### Integration test coverage

| Suite | What it covers |
|---|---|
| `auth_test.go` | Register (success, duplicate email/username, weak password), client login, merchant login, rep login blocked when merchant disabled, full password reset flow |
| `event_test.go` | Event CRUD, all 9 ticket type categories, capacity adjustment on update, delete blocked when tickets sold |
| `purchase_test.go` | Full buy → validate → refund flow, sold out (409), sale window enforcement, max per order, refund restores availability, unauthorized refund |
| `profile_test.go` | Profile create/get/update, phone conflict, email + phone verification end-to-end |
| `admin_test.go` | User/merchant activate/deactivate, merchant deactivation cascades to reps, client token rejected on admin routes |

### Postman Collection

Import `ticketfair.postman_collection.json` and run the requests in the documented order. All requests include automated tests that check status codes, response shape and key field values.

---

## Seed Data

The database is automatically initialized with the following test data on first run.

| Entity | Field | Value |
|---|---|---|
| Merchant | Email | contato@ticketfairprod.com |
| Merchant | Password | `PassW0rd!` |
| Merchant | ID | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |
| Rep (admin role) | Email | carlos@ticketfairprod.com |
| Rep (admin role) | Password | `PassW0rd!` |
| Event | Name | Summer Festival 2026 |
| Event | Location | Cianorte Exhibition Park — PR, Brazil |
| Event | Date | Dec 20, 2026 at 6:00 PM UTC |
| Event | Capacity | 1000 |
| Event | ID | `c3d4e5f6-a7b8-9012-cdef-123456789012` |
| Admin | Email | admin@ticketfair.com |
| Admin | Password | `PassW0rd!` |
| Admin | ID | `d4e5f6a7-b8c9-0123-defa-234567890123` |

Seeded ticket types for the event: Early Bird (R$25), General Admission (R$50), VIP (R$150), Group Pack (R$160/4).

---

## Production Guide

### 1. Secrets and Environment Variables

Never commit `.env` to version control. In production use a secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.).

| Variable | Action |
|---|---|
| `JWT_SECRET` | Generate at least 64 random characters: `openssl rand -hex 64` |
| `STRIPE_SECRET_KEY` | Replace `sk_test_...` with `sk_live_...` |
| `STRIPE_WEBHOOK_SECRET` | Create a production webhook in the Stripe dashboard, copy the signing secret |
| `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` | Use a dedicated IAM user with S3-only permissions |
| `AWS_ENDPOINT_URL` | **Remove entirely** — must not point to LocalStack in production |
| `COCKROACH_USER` | Use a dedicated DB user, not root |
| `SMTP_PASSWORD` | Use an app password or a dedicated transactional email account |
| `GIN_MODE` | Must be `release` |

### 2. Remove Development Services

Remove from `docker-compose.yaml` before deploying:

```yaml
# Remove for production:
stripe-cli:   # webhook forwarding — handled by Stripe dashboard
localstack:   # S3 emulator — use real AWS S3
```

### 3. Database

Enable TLS and use a dedicated user:

```bash
# Create a dedicated DB user
cockroach sql --insecure --execute="
  CREATE USER ticketfair_app WITH PASSWORD 'strong_password';
  GRANT ALL ON DATABASE ticketfair TO ticketfair_app;
"
```

Switch from `--insecure` to `--certs-dir=/certs` in the CockroachDB Docker command. Update the DSN:

```bash
DB_HOST=your-cockroach-host
COCKROACH_USER=ticketfair_app
COCKROACH_PASSWORD=strong_password
DB_SSLMODE=require
```

Enable automatic backups:

```sql
ALTER DATABASE ticketfair SET SCHEDULE BACKUP = '0 2 * * *' INTO 's3://your-backup-bucket';
```

### 4. AWS S3

```bash
aws s3 mb s3://ticketfair-images-prod --region us-east-1
```

Restrict the IAM user to only what TicketFair needs:

```json
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": ["s3:PutObject", "s3:DeleteObject", "s3:GetObject"],
    "Resource": "arn:aws:s3:::ticketfair-images-prod/*"
  }]
}
```

Consider CloudFront in front of the bucket for CDN caching. Update `buildImageURL()` in `services/s3_service.go` to return the CloudFront domain.

### 5. Stripe

1. Replace `sk_test_...` with `sk_live_...` in `.env`
2. Remove `stripe-cli` from `docker-compose.yaml`
3. Stripe Dashboard → Webhooks → Add endpoint: `https://yourdomain.com/api/v1/public/webhooks/stripe`
4. Events: `payment_intent.succeeded`, `payment_intent.payment_failed`, `payment_intent.canceled`, `charge.refunded`
5. Copy the signing secret → `STRIPE_WEBHOOK_SECRET`

### 6. TLS / HTTPS

Caddy issues Let's Encrypt certificates automatically. Update `Caddyfile`:

```caddy
api.yourdomain.com {
  reverse_proxy ticketfair-app:8000
}

app.yourdomain.com {
  reverse_proxy frontend:3002
}

dashboard.yourdomain.com {
  reverse_proxy dashboard:3001
  basicauth {
    admin $2a$14$...  # bcrypt hash
  }
}
```

Ports 80 and 443 must be open and the domain must resolve to your server's IP.

### 7. CORS

Add CORS middleware to Gin if the frontend is served from a different origin:

```go
import "github.com/gin-contrib/cors"

r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"https://app.yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Authorization", "Content-Type"},
    AllowCredentials: true,
}))
```

### 8. Rate Limiting at Scale

The current in-memory token bucket limiter is per-instance. With multiple API replicas, replace with a Redis-backed limiter:

```bash
go get github.com/go-redis/redis/v8
go get github.com/go-redis/redis_rate/v10
```

### 9. Email (SMTP)

For production volume, swap `net/smtp` for a transactional provider. All require only updating `SMTP_*` env vars — no code changes needed:

| Provider | Notes |
|---|---|
| **SendGrid** | Reliable, strong analytics |
| **AWS SES** | Cheapest at scale |
| **Resend** | Simple API, great DX |
| **Mailgun** | Strong EU compliance |

### 10. Observability

Set Loki retention:

```yaml
limits_config:
  retention_period: 744h  # 31 days
```

Set up Prometheus alerting for high error rates, P95 latency and DB connection exhaustion. Lock down Grafana with authentication before making it publicly accessible.

### 11. Seed Data

Remove all `INSERT` statements from `migration/docker-database-init.sql` in production — they contain known credentials. Create the first admin account directly:

```sql
INSERT INTO admins (name, email, password, active)
VALUES ('Super Admin', 'admin@yourdomain.com', '<bcrypt_hash>', true);
```

### 12. Frontend Build

For production, build the React app with the correct API URL:

```bash
cd ticketfair-react
VITE_API_URL=https://api.yourdomain.com/api/v1 npm run build
```

Or via Docker build arg as shown in the [Frontend](#frontend-react) section.

### 13. Pre-Deploy Checklist

```
□ JWT_SECRET is at least 64 random characters
□ GIN_MODE=release
□ AWS_ENDPOINT_URL is not set (removed from .env)
□ S3_BUCKET points to production bucket
□ STRIPE_SECRET_KEY starts with sk_live_
□ STRIPE_WEBHOOK_SECRET is from the production Stripe dashboard
□ stripe-cli service removed from docker-compose.yaml
□ localstack service removed from docker-compose.yaml
□ CockroachDB runs with TLS (--certs-dir, not --insecure)
□ Database user is not root
□ Seed data INSERT statements removed from migration SQL
□ Admin account created via DB console with strong password
□ Dashboard protected behind auth or IP allowlist
□ Grafana protected behind authentication
□ CORS configured for production frontend domain
□ Loki retention period configured
□ Prometheus alerting rules configured
□ React frontend built with production VITE_API_URL
□ All integration tests pass: make test-docker
□ DNS records point to your server
□ Ports 80 and 443 open in firewall
□ Caddyfile uses real domain names
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