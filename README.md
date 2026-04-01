# RetailOS

A minimalist point-of-sale system built for Indian medical shops. Multi-tenant, GST-compliant, fast.

---

## Stack

| | |
|---|---|
| **Frontend** | Next.js 16, Redux Toolkit, Tailwind CSS |
| **Backend** | Go, chi, sqlc, pgx |
| **Database** | PostgreSQL (NeonDB) |

---

## Features

- **Billing** — product search, batch selection, CGST/SGST/IGST auto-computation
- **Inventory** — batch management, expiry tracking, stock levels
- **Orders** — paginated history, search, soft-delete, print view
- **Reports** — GST slab breakdown by date range, CSV export
- **Multi-tenant** — isolated schema per shop, JWT auth
- **Admin** — create and manage shops

---

## Running Locally

```bash
# Backend
cd backend
cp .env.example .env   # fill in your values
go run ./cmd/server    # → :8080

# Frontend
cd frontend
cp .env.example .env.local
npm install && npm run dev   # → :3000
```

See [DEV.md](DEV.md) for full setup details and credentials.

---

## Project Structure

```
retail-os/
├── backend/
│   ├── cmd/server/        # entry point
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/     # auth, tenant isolation
│   │   ├── migrations/    # SQL migrations
│   │   ├── queries/       # sqlc SQL files
│   │   └── generated/     # sqlc generated code
│   └── sqlc.yaml
└── frontend/
    ├── app/               # Next.js App Router pages
    ├── components/        # shared + feature components
    ├── store/             # Redux slices
    └── lib/               # API client, GST helpers
```

---

Built by [Chethanreddy](https://github.com/chethanreddy123)
