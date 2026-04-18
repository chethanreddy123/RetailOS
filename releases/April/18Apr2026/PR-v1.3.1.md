# RetailOS v1.3.1 — Release Notes

### Distributor Management — Know Where Your Stock Comes From

**Release Date:** April 2026
**Branch:** `feat/subhanu` → `main`

---

## Overview

RetailOS v1.3.1 introduces **Distributors** as a first-class entity. Every batch of stock you receive can now be traced to a specific distributor and invoice number. You get a dedicated management page, per-distributor stock views, and a new dashboard section showing your top suppliers by stock value.

This replaces the old free-text distributor fields on the Add Stock form — all existing data is **automatically migrated** with no manual effort required.

---

## What's New

### 1. Distributor Management Page
**Route:** `/distributors` | **Commit:** `599cf7c`

A dedicated page to manage all your suppliers in one place.

- **Directory Table** — Lists all distributors with Name, Phone, Email, Address, and Active/Inactive status badge
- **Search** — Filter distributors by name instantly
- **Add Distributor** — Create a new distributor with name, phone, email, and address via a modal form
- **Edit Distributor** — Update any distributor's details or toggle their active/inactive status
- **Safe Delete** — Deleting a distributor is blocked if any batches are linked to them, preventing data loss. You'll see a clear warning message explaining why
- **Sidebar Navigation** — New "Distributors" entry (truck icon) in the sidebar between Inventory and Orders

**Why it matters:** Before this, distributor information was scattered across individual batch entries as free text — unsearchable, inconsistent, and impossible to update in bulk. Now distributors are a managed entity: update a distributor's phone once and it's correct everywhere.

---

### 2. Stock by Distributor View
**Route:** `/distributors` (lower section) | **Commit:** `599cf7c`

See all the stock you've ever received from a specific supplier.

- **Distributor Dropdown** — Select any distributor from a dropdown on the Distributors page
- **Linked Batches Table** — Shows every batch tied to that distributor: product name, company, batch number, expiry date, MRP, buying price, selling price, purchase quantity, available stock, and purchase invoice number
- **Traceability** — Instantly answer "where did this batch come from?" and "how much of this distributor's stock do I still have?"

**Why it matters:** When a distributor calls about an invoice dispute, or when you need to trace a product recall, you can pull up all batches from that supplier in seconds — no manual searching through batch records.

---

### 3. Distributor Linking on Add Stock & Edit Batch
**Commit:** `599cf7c`

Structuring how distributor information is captured at the point of stock entry.

**Add Stock page:**
- Old 4 freeform text fields (name, phone, location, invoice) → replaced with a **structured distributor dropdown** (active distributors only) + a separate **Purchase Invoice Number** field
- Type a name in the dropdown to search, or select from the list

**Edit Batch modal:**
- Same replacement — dropdown pre-populated with the batch's current distributor
- Change the linked distributor or invoice number anytime

**Why it matters:** Free-text distributor entry led to duplicates ("Reddy Pharma", "reddy pharma", "Reddy pharma co") and broken links. A structured dropdown enforces consistency — every batch points to the same canonical distributor record.

---

### 4. Dashboard — Top Distributors by Stock Value
**Commit:** `599cf7c`

Your dashboard now tells you which suppliers represent the most capital tied up in current stock.

- **Top 5 Distributors** — Ranked by total stock value of non-expired, in-stock batches
- **Columns** — Distributor name, number of linked batches, total purchase quantity, and current stock value
- **Live Data** — Calculated on page load from real-time inventory
- **Conditional Display** — Section only appears when distributor data exists; no empty tables

**Why it matters:** Supplier concentration is a real business risk. If your top distributor has delivery issues, it impacts a large portion of your inventory. This view makes that concentration visible at a glance so you can diversify or plan ahead.

---

### 5. Automatic Data Migration
**Commit:** `599cf7c`

All existing distributor data from previous stock entries is migrated automatically on deploy — no action required from shop owners.

- Old free-text distributor entries (stored as JSONB) are read, deduplicated by name, and inserted as proper distributor records
- Each batch is then linked to its corresponding distributor record
- Purchase invoice numbers are carried over from the old data
- The migration is idempotent — running it multiple times (e.g., across server restarts) produces no duplicate records

**Why it matters:** Zero migration burden on shops. Open the app after the update and your distributors are already there, linked to the right batches.

---

## Summary of Changes

### New Pages & Components
- `app/(dashboard)/distributors/page.tsx` — Distributor management page with search, CRUD table, and stock-by-distributor view
- `components/distributors/AddDistributorModal.tsx` — Create distributor modal
- `components/distributors/EditDistributorModal.tsx` — Edit distributor modal with active toggle
- `lib/distributorCache.ts` — Session-level cache for distributor list (avoids redundant API calls across page navigations)

### Modified Frontend
- `components/shared/Sidebar.tsx` — New Distributors nav entry (Truck icon)
- `app/(dashboard)/inventory/add/page.tsx` — Distributor dropdown + invoice field replaces 4 freeform fields
- `components/inventory/EditBatchModal.tsx` — Same replacement in edit flow
- `app/(dashboard)/dashboard/page.tsx` — Top Distributors table section
- `lib/api.ts` — Added `listDistributors`, `createDistributor`, `updateDistributor`, `deleteDistributor`, `listBatchesByDistributor`
- `types/index.ts` — Added `Distributor`, `DistributorBatchRow` interfaces; updated `InventoryRow`, `Batch`, `DashboardData`

### New Backend
- `handlers/distributors.go` — Full CRUD handler with delete guard (blocks delete if batches are linked)
- `queries/distributors.sql` — `CreateDistributor`, `ListDistributors`, `UpdateDistributor`, `DeleteDistributor`, `CountBatchesByDistributor`, `ListBatchesByDistributor`, `GetDistributorStats`
- `queries/dashboard.sql` — `DashboardDistributorStats` query (top 5 by stock value)
- 5 new routes registered under JWT-authenticated tenant group: `GET/POST /distributors`, `PUT/DELETE /distributors/{id}`, `GET /distributors/{id}/batches`

### Database Migration
- `000009_create_distributors_and_migrate.up.sql` — Creates `distributors` table; adds `distributor_id` FK and `purchase_invoice_no` to `batches`; auto-migrates existing JSONB data with deduplication

---

## Technical Notes

- Distributor FK on `batches` is `ON DELETE SET NULL` — deleting a distributor nulls the FK on linked batches rather than cascade-deleting stock records. The application-level guard (HTTP 400 if batches exist) prevents accidental orphaning
- The session-level distributor cache (`distributorCache.ts`) uses `sessionStorage` — persists across Next.js client-side navigations within a session but clears on page refresh, keeping the dropdown fresh without redundant API calls
- Dashboard distributor stats query computes stock value using `landing_price` (if available) falling back to `buying_price` — accounts for purchase GST being factored into landed cost
- Migration is guarded with `NOT EXISTS (SELECT 1 FROM distributors)` so the JSONB auto-import only runs on first deploy; subsequent restarts go through the deduplication block only
- Migration idempotency was tested across multiple server restarts on live tenant schemas including ones with pre-existing JSONB distributor data

---

## Test Plan

- [ ] Open `/distributors` — verify existing distributor records appear (auto-migrated from old data)
- [ ] Add a new distributor → verify it appears in the table and in the dropdown on Add Stock
- [ ] Edit a distributor's phone and email → verify changes persist on reload
- [ ] Toggle a distributor to Inactive → verify it no longer appears in the Add Stock dropdown
- [ ] Try to delete a distributor that has batches linked → verify delete is blocked with a clear error message
- [ ] Delete a distributor with no linked batches → verify it's removed cleanly
- [ ] Open Add Stock page → verify distributor dropdown shows active distributors only, with search working
- [ ] Add stock and link it to a distributor with an invoice number → verify both appear in the batch list
- [ ] Open Edit Batch modal → verify existing distributor is pre-selected; change it and save → verify update persists
- [ ] Select a distributor from the Distributors page dropdown → verify linked batches table shows correct rows
- [ ] Open the Dashboard → verify "Top Distributors by Stock Value" section appears and shows correct ranking
- [ ] Confirm version shows as `v1.3.1` in the bottom-right corner of the UI

---

**RetailOS v1.3.1** — Your supply chain, finally organized.
