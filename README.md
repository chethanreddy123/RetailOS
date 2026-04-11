# RetailOS

A minimalist point-of-sale system built for Indian medical shops. Multi-tenant, GST-compliant, fast.

---

## Stack

| Layer | Technology |
|---|---|
| **Frontend** | Next.js 16, Redux Toolkit, Tailwind CSS, shadcn/ui |
| **Backend** | Go 1.22, chi router, sqlc, pgx v5 |
| **Database** | PostgreSQL (NeonDB Serverless) |
| **Auth** | Custom JWT (golang-jwt) — no third-party auth |
| **Hosting** | Vercel (frontend) + Render (backend) |

---

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                          Client Browser                          │
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTPS
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Vercel (Next.js 16)                           │
│                                                                  │
│  /(auth)          /(dashboard)              /super-admin         │
│  └── /login       ├── /billing             ├── /login           │
│                   ├── /inventory           └── /dashboard        │
│                   ├── /orders                                    │
│                   ├── /reports             Redux Store           │
│                   └── /admin              (authSlice, cartSlice) │
└────────────────────────────┬────────────────────────────────────┘
                             │ REST + Bearer JWT
                             │ NEXT_PUBLIC_API_URL
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Render (Go / chi router)                       │
│                                                                  │
│  Public routes              Tenant routes (JWT + search_path)    │
│  POST /auth/login           GET/POST /products                   │
│  POST /super-admin/login    GET/POST /batches                    │
│  POST /super-admin/verify   GET      /inventory                  │
│                             GET      /customers                  │
│  Admin routes (JWT role)    POST/GET /orders                     │
│  POST /super-admin/tenants  GET      /reports/gst                │
│  GET  /super-admin/tenants                                       │
│  PATCH /super-admin/tenants/:id                                  │
│                                                                  │
│  Middleware stack:                                               │
│  Logger → Recoverer → CORS → RateLimit → JWTAuth → TenantCtx    │
└────────────────────────────┬────────────────────────────────────┘
                             │ pgx pool (TLS)
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                  NeonDB (PostgreSQL Serverless)                  │
│                                                                  │
│  public schema                  tenant_<id> schema (per shop)    │
│  ├── tenants                    ├── products                     │
│  └── super_admins               ├── batches                      │
│                                 ├── customers                    │
│                                 ├── orders                       │
│                                 └── order_items                  │
└─────────────────────────────────────────────────────────────────┘
```

### Multi-Tenancy Model

```
Login (username + password)
        │
        ▼
  Lookup tenant in public.tenants
        │
        ▼
  Issue JWT  { tenant_id, schema_name, order_prefix, ... }
        │
        ▼
  Every API request → middleware reads schema_name from JWT
        │
        ▼
  SET search_path = tenant_<id>   ← all queries run inside this schema
        │
        ▼
  Tenant data is fully isolated — no cross-tenant queries possible
```

### Auth Flow

```
Shop Owner                       Super Admin
─────────────────                ──────────────────────────
POST /auth/login                 POST /super-admin/auth/login
  username + password              username + password
        │                                │
        ▼                                ▼
  JWT { role: tenant }           OTP sent to registered email
        │                                │
        │                        POST /super-admin/auth/verify-otp
        │                                │
        │                                ▼
        │                          JWT { role: super_admin }
        │                                │
        ▼                                ▼
  Stored in Redux              Stored in localStorage (sa_token)
  Access: /billing etc.        Access: /super-admin/dashboard
```

### Order Creation (ACID Transaction)

```
POST /orders
    │
    ├── BEGIN transaction
    ├── INSERT orders row
    ├── For each item:
    │     ├── SELECT batch FOR UPDATE  ← prevents overselling
    │     ├── Check available_stock >= qty
    │     ├── INSERT order_items row
    │     └── UPDATE batch SET available_stock -= qty
    ├── COMMIT
    └── Return order with GST breakdown
```

---

## Features

- **Billing** — product search, batch selection, CGST/SGST/IGST auto-computation, print bill
- **Inventory** — batch management, expiry tracking, low-stock alerts
- **Orders** — paginated history, search by bill/name/phone, soft-delete, print view
- **Reports** — GST slab breakdown by date range, CSV export
- **Multi-tenant** — isolated PostgreSQL schema per shop, JWT auth
- **Super Admin** — OTP-secured dashboard to create and manage shops

---

## Running Locally

```bash
# Backend
cd backend
cp .env.example .env   # fill in your values
go run ./cmd/server    # → :8080

# Frontend
cd frontend
npm install
npm run dev            # → :3000
```

See [DEV.md](DEV.md) for full setup, env vars, and credentials.

---

## Project Structure

```
retail-os/
├── backend/
│   ├── cmd/
│   │   ├── server/           # main entry point
│   │   └── seed-admin/       # one-time super admin seed
│   ├── internal/
│   │   ├── config/           # env var loading
│   │   ├── db/               # pool, migrations runner
│   │   ├── email/            # SMTP OTP sender
│   │   ├── handlers/         # HTTP handlers (auth, inventory, orders, reports, admin)
│   │   ├── middleware/        # JWT auth, tenant context, rate limiter
│   │   ├── migrations/
│   │   │   ├── public/       # tenants, super_admins tables
│   │   │   └── tenant/       # per-shop schema migrations
│   │   ├── queries/          # sqlc SQL source files
│   │   └── generated/        # sqlc generated Go code
│   └── sqlc.yaml
└── frontend/
    ├── app/
    │   ├── (auth)/           # login page
    │   ├── (dashboard)/      # billing, inventory, orders, reports, admin
    │   └── super-admin/      # super admin login + dashboard
    ├── components/
    │   ├── billing/          # AddItemBar, LineItem
    │   └── shared/           # Sidebar, Pagination, TableSkeleton, AuthGuard
    ├── store/                # Redux slices (auth, cart)
    └── lib/                  # API client, GST helpers
```

---

Built by [AIOverflow](https://github.com/AIOverflow-in)
