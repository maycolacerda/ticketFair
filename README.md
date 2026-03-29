# 🎟️ TicketFair

**Plataforma de venda de ingressos construída com Go, Gin, GORM e CockroachDB.**

---

## Índice

- [Visão Geral](#visão-geral)
- [Stack Tecnológica](#stack-tecnológica)
- [Arquitetura](#arquitetura)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Como Rodar](#como-rodar)
- [Variáveis de Ambiente](#variáveis-de-ambiente)
- [Rotas da API](#rotas-da-api)
- [Autenticação](#autenticação)
- [Rate Limiting](#rate-limiting)
- [Banco de Dados](#banco-de-dados)
- [Observabilidade](#observabilidade)
- [Dashboard Admin](#dashboard-admin)
- [Testes](#testes)
- [Dados de Seed](#dados-de-seed)

---

## Visão Geral

TicketFair é uma API REST para venda e gerenciamento de ingressos para eventos. A plataforma suporta três tipos de usuários — **clientes**, **merchants** (produtoras) e **administradores** — cada um com seu próprio fluxo de autenticação e permissões.

### Funcionalidades principais

- Cadastro e autenticação de usuários, merchants e representantes
- Criação e gerenciamento de eventos
- Compra e reembolso de ingressos com controle de capacidade atômico
- Validação de ingressos na entrada do evento
- Verificação de email e telefone
- Painel administrativo com controle de usuários e merchants
- Rate limiting por IP para proteção contra brute force
- Logs estruturados com Loki + Grafana

---

## Stack Tecnológica

| Camada | Tecnologia |
|---|---|
| Linguagem | Go 1.23 |
| Framework HTTP | Gin |
| ORM | GORM |
| Banco de Dados | CockroachDB (compatível com Postgres) |
| Autenticação | JWT HS256 (golang-jwt/jwt v5) |
| Senhas | bcrypt (DefaultCost) |
| Validação | go-playground/validator v10 |
| Documentação | Swagger (swaggo) |
| Reverse Proxy | Caddy |
| Logs | slog + Loki + Promtail |
| Métricas | Prometheus + Grafana |
| Containerização | Docker + Docker Compose |

---

## Arquitetura

```
┌─────────────────────────────────────────────────┐
│                    Caddy                         │
│         (Reverse Proxy / TLS)                    │
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

Observabilidade:
Promtail → Loki → Grafana
App      → Prometheus → Grafana
```

### Camadas da aplicação

```
controllers/   HTTP apenas — bind, validação, resposta
services/      Lógica de negócio — sem gin, sem net/http
models/        Structs GORM
dto/           Shapes de entrada/saída, validadores customizados
middlewares/   JWT, roles, rate limiting, logging
routes/        Registro de rotas
database/      Conexão e migração
configs/       Email (SMTP)
```

---

## Estrutura do Projeto

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

## Como Rodar

### Pré-requisitos

- Docker
- Docker Compose

### 1. Clone o repositório

```bash
git clone https://github.com/maycolacerda/ticketFair
cd ticketFair
```

### 2. Configure as variáveis de ambiente

```bash
cp .env.example .env
# Edite o .env com seus valores
```

### 3. Suba o stack

```bash
docker compose up --build
```

### 4. Gere a documentação Swagger (primeira vez)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g main.go
docker compose up --build ticketfair-app
```

### 5. Acesse

| Serviço | URL |
|---|---|
| API | http://localhost:8000 |
| Swagger | http://localhost:8000/swagger/index.html |
| Dashboard Admin | http://localhost:3001 |
| CockroachDB UI | http://localhost:8081 |
| Grafana | http://localhost:3000 |
| Prometheus | http://localhost:9090 |

---

## Variáveis de Ambiente

```bash
# App
APP_VERSION=1.0.0
GIN_MODE=release

# Banco de Dados (CockroachDB)
DB_HOST=ticketfair-db
DB_PORT=26257
COCKROACH_USER=root
COCKROACH_DB=ticketfair

# JWT
JWT_SECRET=sua_chave_secreta_muito_longa

# SMTP (opcional — sem isso emails são apenas logados)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=seu@email.com
SMTP_PASSWORD=sua_app_password
SMTP_FROM=seu@email.com
SMTP_FROM_NAME=TicketFair
```

---

## Rotas da API

### 🌐 Públicas

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

### 🔒 Privadas (requer token de cliente)

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

### 🏪 Merchant (requer token de merchant)

```
PUT  /api/v1/merchant/update
POST /api/v1/merchant/events/new
PUT  /api/v1/merchant/events/:id
POST /api/v1/merchant/tickets/:id/validate
POST /api/v1/merchant/rep/new
PUT  /api/v1/merchant/rep/:id
POST /api/v1/merchant/logout
```

### 🔑 Admin (requer token de admin)

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

## Autenticação

Todos os endpoints protegidos exigem um JWT no header:

```
Authorization: Bearer <token>
```

### Roles disponíveis

| Role | Acesso |
|---|---|
| `client` | Rotas `/private/*` |
| `merchant` | Rotas `/merchant/*` |
| `admin` / `manager` / `staff` | Rotas de merchant rep |
| `superadmin` | Rotas `/admin/*` |

### Claims do JWT

```json
{
  "user_id":    "uuid",
  "role":       "client | merchant | admin | manager | staff | superadmin",
  "merchant_id": "uuid (apenas para reps)",
  "exp":        1234567890,
  "iss":        "ticketfair"
}
```

---

## Rate Limiting

Proteção por IP usando token bucket algorithm.

| Endpoint | Burst | Taxa de recarga | Mensagem |
|---|---|---|---|
| Login (todos) | 5 req | 1 a cada 3s | too many login attempts |
| Register | 3 req | 1 a cada 10s | too many registration attempts |
| Verificação | 3 req | 1 a cada 60s | too many verification attempts |
| Eventos públicos | 30 req | 10/s | too many requests |
| Admin | 20 req | 5/s | too many requests |

Quando o limite é atingido:

```
HTTP 429 Too Many Requests
Retry-After: 60
{ "error": "too many login attempts — try again later" }
```

---

## Banco de Dados

### Tabelas

```
admins           — Administradores da plataforma
users            — Usuários clientes
merchants        — Produtoras de eventos
merchant_reps    — Representantes das produtoras
profiles         — Perfil dos usuários (1:1 com users)
addresses        — Endereços dos perfis (1:1 com profiles)
events           — Eventos criados pelos merchants
transactions     — Compras de ingressos
tickets          — Ingressos vinculados às transações
verifications    — Códigos de verificação de email/telefone
```

### Funções SQL (atomicidade)

```sql
create_profile_with_address(...)  -- Cria perfil + endereço em uma transação
purchase_ticket(...)              -- Decrementa capacidade + cria transação atomicamente
refund_ticket(...)                -- Restaura capacidade + marca transação como reembolsada
```

### Ciclo de vida dos ingressos

```
Compra  →  active
Validar →  used
Reembolso → refunded
```

---

## Observabilidade

### Logs estruturados

Em produção (`GIN_MODE=release`) todos os logs são emitidos em JSON e coletados pelo Promtail para indexação no Loki.

```json
{
  "time": "2026-03-29T12:00:00Z",
  "level": "INFO",
  "msg": "Client login successful",
  "user_id": "uuid"
}
```

### Acesso ao Grafana

```
URL:   http://localhost:3000
User:  admin
Pass:  admin
```

Configure o Loki como data source em `http://loki:3100` para visualizar os logs.

---

## Dashboard Admin

Interface web para gerenciamento da plataforma.

**Acesso:** http://localhost:3001

**Login:**
```
Email: admin@ticketfair.com
Senha: PassW0rd!
```

**Funcionalidades:**
- Overview com estatísticas em tempo real
- Listar, criar, ativar e desativar usuários
- Listar, criar, ativar e desativar merchants
- Visualizar todos os eventos ativos

---

## Testes

### Rodar testes unitários

```bash
go test ./...

# Com output detalhado
go test ./... -v

# Pacote específico
go test ./controllers/... -v
```

### Coleção Postman

Importe o arquivo `ticketfair.postman_collection.json` no Postman para testar todos os endpoints com testes automatizados e variáveis pré-configuradas.

**Ordem recomendada:**

```
1.  Admin Login
2.  Merchant Login (conta seed)
3.  Rep Login (conta seed)
4.  Register User
5.  Client Login
6.  Create Profile
7.  Send Email Verification → checar logs → Verify Email
8.  Purchase Ticket (Festival de Verão 2026)
9.  List My Tickets
10. Validate Ticket (merchant)
```

---

## Dados de Seed

O banco é inicializado automaticamente com dados de teste.

### Merchant

| Campo | Valor |
|---|---|
| Nome | TicketFair Produções |
| Email | contato@ticketfairprod.com |
| Senha | `PassW0rd!` |
| ID | `a1b2c3d4-e5f6-7890-abcd-ef1234567890` |

### Merchant Rep (Admin)

| Campo | Valor |
|---|---|
| Nome | Carlos Admin |
| Email | carlos@ticketfairprod.com |
| Senha | `PassW0rd!` |
| Role | admin |

### Evento

| Campo | Valor |
|---|---|
| Nome | Festival de Verão 2026 |
| Local | Parque de Exposições de Cianorte — PR |
| Data | 20/12/2026 às 18h |
| Capacidade | 1000 |
| ID | `c3d4e5f6-a7b8-9012-cdef-123456789012` |

### Admin

| Campo | Valor |
|---|---|
| Email | admin@ticketfair.com |
| Senha | `PassW0rd!` |
| ID | `d4e5f6a7-b8c9-0123-defa-234567890123` |

---

## Contribuindo

```bash
# 1. Crie uma branch
git checkout -b feature/minha-feature

# 2. Faça suas alterações e commit
git add .
git commit -m "feat: descrição da feature"

# 3. Push e abra um PR
git push origin feature/minha-feature
```

---

## Licença

MIT — veja [LICENSE](LICENSE) para detalhes.