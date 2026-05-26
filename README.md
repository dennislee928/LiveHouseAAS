# LiveHouseAAS — Live House 展演整合平台

A full-stack SaaS platform for Taiwanese Live House venues and independent musicians. Built with Go + Gin + Next.js + PostgreSQL.

---

## Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Next.js 14  │────▶│   Go API     │────▶│  PostgreSQL  │
│  (Vercel)    │     │  (WSO2/Docker)│    │  (AlwaysData)│
├──────────────┤     ├──────────────┤     ├──────────────┤
│  Tailwind    │     │  Gin + JWT   │     │  pgx pool    │
│  shadcn/ui   │     │  sqlc types  │     │  Redis 7     │
│  WebSocket   │     │  Payments    │     ├──────────────┤
└──────────────┘     │  Blockchain  │     │  Redis       │
                     │  Analytics   │     │  (session/   │
                     │  Admin       │     │   cache)     │
                     │  Notif/SMTP  │     └──────────────┘
                     └──────────────┘
```

### Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.26+, Gin, pgx/v5, sqlc |
| **Frontend** | Next.js 14 (App Router), Tailwind CSS, shadcn/ui |
| **Database** | PostgreSQL 16, Redis 7 |
| **Queue** | NATS (declared, pending integration) |
| **File Storage** | MinIO (dev) / S3-compatible (prod) |
| **Auth** | JWT (HS256, 72h expiry) |
| **Payments** | ECPay, NewebPay, Binance Pay, EVM Smart Contract |
| **NFT** | Solidity ERC-721 (LiveHouseTicket.sol) |
| **Infrastructure** | Docker Compose, Cloudflare CDN |

---

## Quick Start

### Prerequisites
- Go 1.26+, Node.js 22+, Docker + Docker Compose

### 1. Clone & configure
```bash
git clone <repo-url> && cd LiveHouseAAS
cp backend/.env.example backend/.env
```

### 2. Start dependencies
```bash
docker compose -f docker-compose.dev.yml up -d postgres redis minio nats
```

### 3. Run database migrations
```bash
cd backend && go run ./cmd/migrate up
```

### 4. Start API (with hot-reload)
```bash
cd backend && air
```

### 5. Start frontend
```bash
cd frontend && npm install && npm run dev
```

### 6. Open
- Frontend: http://localhost:3000
- API Health: http://localhost:8080/health

---

## Project Structure

```
LiveHouseAAS/
├── backend/
│   ├── cmd/
│   │   ├── api/              # API server entry point
│   │   └── migrate/          # Database migration runner
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handler/      # HTTP handlers (18 files)
│   │   │   └── middleware/   # Auth, CORS, Rate limit
│   │   ├── auth/             # JWT generation/validation
│   │   ├── blockchain/       # NFT + POAP service
│   │   ├── config/           # Environment config
│   │   ├── domain/           # Domain types (booking, event, etc.)
│   │   ├── infra/
│   │   │   ├── cache/        # Redis client + helpers
│   │   │   ├── db/           # PostgreSQL connection pool
│   │   │   └── queue/        # NATS (pending)
│   │   ├── notification/     # Email/SMTP service
│   │   └── payment/          # Payment providers (ECPay, NewebPay, Crypto)
│   ├── migrations/           # 15 migration pairs (up/down)
│   ├── Dockerfile.prod       # Multi-stage production build (~15MB)
│   └── .env.example
├── frontend/
│   ├── app/
│   │   ├── (dashboard)/      # Dashboard pages (admin, venues, analytics, etc.)
│   │   ├── login/            # Login page
│   │   ├── register/         # Registration page
│   │   ├── forgot-password/  # Password reset request
│   │   ├── reset-password/   # Password reset form
│   │   ├── layout.tsx        # Root layout (SEO metadata)
│   │   ├── error.tsx         # Error boundary
│   │   ├── loading.tsx       # Root loading state
│   │   └── not-found.tsx     # 404 page
│   ├── components/
│   │   ├── seats/SeatMap.tsx # Interactive seat map
│   │   ├── ui/               # shadcn/ui components
│   │   └── Pagination.tsx    # Shared pagination component
│   ├── hooks/
│   │   └── useWebSocket.ts  # WebSocket connection hook
│   └── lib/
│       └── api.ts           # API client (GET/POST/PUT/DELETE)
├── scripts/
│   ├── backup.sh            # pg_dump backup (keeps last 7)
│   └── restore.sh           # pg_restore with --clean
├── docker-compose.dev.yml   # Dev environment
├── docker-compose.prod.yml  # Production setup
└── .github/workflows/
    ├── ci.yml               # Go vet/build/test + Next.js build
    └── deploy.yml           # Docker build & push to WSO2
```

---

## API Endpoints (80+)

### Public
| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/ws` | WebSocket (JWT in query param) |
| POST | `/api/v1/auth/register` | User registration |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password with token |
| GET | `/api/v1/search/events` | Search events (keyword/city/date) |
| GET | `/api/v1/search/venues` | Search venues (keyword/city/capacity) |
| GET | `/api/v1/verify-email` | Verify email with token |
| POST | `/api/v1/payments/callback` | Generic payment callback |
| POST | `/api/v1/payments/ecpay/notify` | ECPay notification |
| POST | `/api/v1/payments/newebpay/notify` | NewebPay notification |

### Authenticated
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/me` | Get current user |
| GET | `/api/v1/dashboard/stats` | Role-based dashboard stats |
| PUT | `/api/v1/me/profile` | Update profile |
| POST | `/api/v1/me/change-password` | Change password |
| POST | `/api/v1/me/avatar` | Update avatar |
| POST | `/api/v1/me/verify-email` | Request email verification |
| GET | `/api/v1/notifications` | List user notifications |
| GET | `/api/v1/notifications/unread` | Unread notification count |
| PUT | `/api/v1/notifications/:id/read` | Mark notification as read |
| POST | `/api/v1/upload` | Generic file upload |
| POST | `/api/v1/upload/kyb` | KYB document upload |
| POST | `/api/v1/events/:id/upload` | Event image upload |

### Venues (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/api/v1/venues` | List / Create venues |
| GET/PUT/DELETE | `/api/v1/venues/:id` | Get / Update / Delete venue |
| GET/POST/DELETE | `/api/v1/venues/:id/specs` | Manage venue specs |
| GET/POST/DELETE | `/api/v1/venues/:venueId/slots` | Manage time slots |
| POST | `/api/v1/venues/:venueId/slots/batch` | Batch create slots |
| GET/PUT | `/api/v1/venues/:venueId/seats` | Seat layout management |

### Bookings (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/bookings` | Create booking request |
| PUT | `/api/v1/bookings/:id/status` | Approve/reject/confirm/cancel |
| GET | `/api/v1/bookings/artist` | Artist's bookings |
| GET | `/api/v1/bookings/owner` | Venue owner's bookings |

### Events (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| GET/POST | `/api/v1/events/:id/ticket-types` | Manage ticket types |
| PUT | `/api/v1/events/:id` | Update event |
| POST | `/api/v1/events/:id/publish` | Publish event |
| POST | `/api/v1/events/:id/purchase` | Purchase tickets |

### Orders & Tickets (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/orders` | List user orders |
| POST | `/api/v1/orders/:orderId/refund` | Request refund |
| GET | `/api/v1/tickets` | List user tickets |
| POST | `/api/v1/tickets/verify` | Verify/check-in ticket |
| GET | `/api/v1/tickets/lookup` | Lookup ticket by code |

### NFT (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/tickets/:ticketId/nft/claim` | Claim NFT ticket |
| POST | `/api/v1/tickets/:ticketId/nft/poap` | Claim POAP |
| GET | `/api/v1/tickets/:ticketId/nft` | NFT status |

### KYB (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/kyb` | Submit KYB (venue only) |
| GET | `/api/v1/kyb` | Get KYB status |

### Analytics (authenticated)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/analytics/summary` | Platform summary |
| GET | `/api/v1/analytics/revenue` | Revenue over time |
| GET | `/api/v1/analytics/top-venues` | Top venues |
| GET | `/api/v1/analytics/top-events` | Top events |
| GET | `/api/v1/analytics/booking-trends` | Booking trends |
| GET | `/api/v1/analytics/venue-performance` | Per-venue performance |

### Admin (authenticated + admin role)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/admin/stats` | Platform stats |
| GET/PUT | `/api/v1/admin/users` | List / Update user role |
| GET/PUT | `/api/v1/admin/venues` | List / Update venue status |
| GET | `/api/v1/admin/events` | List all events |
| GET | `/api/v1/admin/bookings` | List all bookings |
| GET | `/api/v1/admin/orders` | List all orders |
| GET/PUT | `/api/v1/admin/kyb` | List pending / Review KYB |
| GET | `/api/v1/admin/notifications` | All notifications (cross-user) |
| POST | `/api/v1/admin/notifications/broadcast` | Broadcast notification |

---

## Payment Providers

| Provider | Type | Status |
|----------|------|--------|
| **Mock** | Test auto-complete | ✅ Active |
| **ECPay** | CheckMacValue HMAC-SHA256, AES encrypt, callback/refund | ✅ Implemented (needs merchant keys) |
| **NewebPay** | AES-CBC encrypt/decrypt, SHA256, MPG gateway | ✅ Implemented (needs merchant keys) |
| **Binance Pay** | HMAC-SHA256 signing, nonce, order create/refund | ✅ Implemented (needs API keys) |
| **Crypto EVM** | Payment address generation, network selection | ✅ Implemented (needs RPC + contract) |

---

## Frontend Pages (31 total)

| Route | Role | Description |
|-------|------|-------------|
| `/login` | Public | Login form |
| `/register` | Public | Registration form |
| `/forgot-password` | Public | Request password reset |
| `/reset-password` | Public | Reset password with token |
| `/dashboard` | All | Role-based dashboard |
| `/settings` | All | Profile & password management |
| `/notifications` | All | Notification list |
| `/search` | All | Search venues & events |
| `/analytics` | Admin/Venue | Data analytics dashboard |
| `/venues` | All | Venue list |
| `/venues/new` | Venue | Create venue |
| `/venues/[id]` | All | Venue detail + slots + bookings |
| `/events` | All | Event list |
| `/events/[id]` | All | Event detail + purchase |
| `/bookings` | Venue/Artist | Booking requests |
| `/orders` | All | Order list |
| `/tickets` | Artist | Ticket list |
| `/verify` | Venue/Admin | QR ticket verification |
| `/nft` | All | NFT ticket management |
| `/kyb` | Venue | KYB submission & status |
| `/admin` | Admin | Admin dashboard |
| `/admin/users` | Admin | User management |
| `/admin/venues` | Admin | Venue management |
| `/admin/events` | Admin | Event management |
| `/admin/bookings` | Admin | Booking management |
| `/admin/orders` | Admin | Order management |
| `/admin/kyb` | Admin | KYB review |
| `/admin/notifications` | Admin | Notification history |
| `/admin/broadcast` | Admin | Broadcast notifications |

---

## Database Migrations (15 pairs)

| # | Migration | Description |
|---|-----------|-------------|
| 001 | `create_users` | Users table (id, email, password_hash, name, role, avatar_url) |
| 002 | `create_venues` | Venues + venue_specs tables |
| 003 | `create_slots` | Slots with conflict detection |
| 004 | `create_bookings` | Booking requests + auto event creation |
| 005 | `create_ticketing` | Ticket types, orders, tickets |
| 006 | `create_kyb` | Business verifications |
| 007 | `create_nft` | NFT tickets + POAPs |
| 008 | `create_notifications` | Notifications (user_id, type, title, body, JSONB data) |
| 009 | `create_seat_layout` | Seat layouts + ticket type extensions |
| 010 | `create_user_tokens` | User tokens + email_verified column |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | API server port |
| `DATABASE_URL` | postgres://... | PostgreSQL connection string |
| `REDIS_URL` | redis://... | Redis connection string |
| `JWT_SECRET` | (dev default) | JWT signing key |
| `GIN_MODE` | `debug` | Gin framework mode |
| `UPLOAD_DIR` | `/app/uploads` | File upload directory |
| `NATS_URL` | nats://localhost:4222 | NATS server URL |
| `CORS_ORIGIN` | `*` | Allowed CORS origin |
| `DB_MAX_CONNS` | `25` | PostgreSQL max connections |
| `DB_MIN_CONNS` | `5` | PostgreSQL min connections |
| `SMTP_HOST` | (empty) | SMTP server host |
| `SMTP_PORT` | `587` | SMTP server port |
| `SMTP_USER` | (empty) | SMTP username |
| `SMTP_PASSWORD` | (empty) | SMTP password |
| `FROM_EMAIL` | `noreply@livehouseaas.com` | Sender email address |

---

## Testing

```bash
# Backend: all tests
cd backend && go test ./... -count=1 -timeout 120s

# With verbose output
go test ./... -v

# Specific package
go test ./internal/payment/ -v

# Frontend build check
cd frontend && npm run build
```

### Test Coverage

| Package | Tests | Description |
|---------|-------|-------------|
| `handler` | 3 | Health check, auth endpoints, protected paths |
| `payment` | 6 | Mock provider, ECPay CheckMacValue, NewebPay AES/SHA256 |
| `middleware` | 3 | Rate limiter logic |
| `auth` | 4 | JWT creation, validation, error cases |

---

## Deployment

### Production Docker
```bash
docker compose -f docker-compose.prod.yml up -d api
```

### CI/CD Pipeline
- **CI** (`.github/workflows/ci.yml`): Go vet → build → migrations → tests → Next.js build
- **Deploy** (`.github/workflows/deploy.yml`): Docker build & push → SSH deploy to WSO2

### Key Production Configurations
- Multi-stage Docker build (~15MB alpine image)
- HEALTHCHECK on API container
- Structured JSON logging via `slog`
- Graceful shutdown (30s timeout for active connections)
- Configurable DB pool (25 max / 5 min connections)
- Rate limiting (100 req/min per user/IP)

---

## Features

### B2B (Venue Management)
- Venue CRUD with specs (capacity, equipment, etc.)
- Time slot calendar with conflict detection
- Booking request lifecycle (pending → approved → confirmed → rejected)
- Automatic event creation on booking approval
- Seat map layout editor

### B2C (Ticketing)
- Ticket type management (name, price, quantity, seat sections)
- Purchase flow with QR code generation
- Multi-provider payment abstraction
- QR ticket verification with check-in

### NFT Tickets
- ERC-721 smart contract (OpenZeppelin-based)
- Mint ticket as NFT after purchase
- POAP (Proof of Attendance Protocol) claims
- Metadata generation and storage

### KYB (Know Your Business)
- Business verification submission (tax ID, documents)
- Admin review workflow (verify/reject with reason)
- Venue gating (only verified venues can publish/payout)

### Notifications
- In-app notifications (database-backed)
- Real-time WebSocket delivery
- Admin broadcast to all/venues/artists/user
- Email via SMTP (configurable)
- Auto-notify: KYB submission, KYB review, booking create, booking status

### Analytics
- Revenue over time (daily/weekly/monthly)
- Top venues and events by orders/revenue
- Booking trends
- Per-venue performance breakdown
- Growth percentage calculations

### Admin Panel
- Platform stats (users, venues, events, orders, revenue)
- Full CRUD management for users, venues, events, bookings, orders
- KYB review queue
- System notification broadcast
- Cross-user notification history

---

## License

Proprietary — All rights reserved.
