# SellOS v1.5.0 — Release Notes

### Out-of-stock visibility in billing search, and a full rebrand from RetailOS to SellOS.

**Release Date:** June 2026
**Branch:** `feat/subhanu` → `main`

---

## Overview

Two changes ship together. The first is a billing UX fix that prevents shopkeepers from selecting out-of-stock products during checkout. The second is a full rename across every user-facing surface — landing page, app, login pages, sidebar, super-admin OTP email, and backend boot log — from RetailOS to SellOS. The rebrand is cosmetic-only: no API contracts, behaviour, or database schema changed.

---

## What's New

### 1. Out-of-stock flagging in billing search
**Commit:** `f4a3d5c`

The billing screen's product search now surfaces stock state at the row level. Products with zero available stock across all batches are visually marked and disabled, so a shopkeeper can no longer accidentally add an out-of-stock item to a bill.

**Why it matters:** Previously a search result for an out-of-stock product would still appear selectable. Clicking it led to an empty batch picker — confusing UX and lost time at the counter. The bad path is now impossible to reach.

**Technical Notes:**
- Stock aggregate computed in the product search query (`internal/queries/products.sql`)
- Generated sqlc output regenerated to match
- `BillingRowItem` renders the out-of-stock disabled state
- New `isOutOfStock` field added to the product type in `frontend/types/index.ts`

---

### 2. SellOS rebrand (UI-facing only)
**Commit:** `4475934`

Every user-facing reference to "RetailOS" replaced with "SellOS" across the landing site (`sellos.in`), the app (`app.sellos.in`), the super-admin OTP email, and the backend boot log. Logo badge letters on the landing navbar and footer flipped from "R" to "S". `BASE_URL` and `APP_URL` constants now point at the new sellos.in domains, so canonical URLs and OG-tag metadata stay consistent for SEO and social sharing.

**Why it matters:** "RetailOS" was unavailable as an owned brand name and matching domains. The team purchased `sellos.in`, `app.sellos.in`, and `api.sellos.in` and switched the production deployments over. The product surface needs to match the brand customers will see in their browser, their inbox, and on social previews.

**Technical Notes:**
- Cosmetic-only: zero API endpoints, zero DB columns, zero behaviour changed
- Go module path (`github.com/retail-os/backend`) and internal Go import paths deliberately left untouched (invisible to users, pure churn)
- `JWT_SECRET` and admin-key rotation NOT in this release — out of scope, to be coordinated separately at a low-traffic time
- Dead file `frontend/app/(landing)/LandingPage.tsx` (523 lines, zero imports anywhere in the tree, no `page.tsx` to make it a route) deleted alongside the rebrand
- Real email account `retailos.noreply@gmail.com` kept as-is; only the display name in `SMTP_FROM` / `RESEND_FROM` changes to "SellOS"

---

## Summary of Changes

### Deleted Files
- `frontend/app/(landing)/LandingPage.tsx` — orphaned near-duplicate of `landing/app/page.tsx` with zero imports; the empty parent route-group folder was removed too

### Modified Files (highlights)

**Frontend:**
- `app/layout.tsx` — `BASE_URL` → `https://app.sellos.in`; SEO, OG, and JSON-LD metadata
- `app/(auth)/login/page.tsx` — login wordmark
- `app/super-admin/login/page.tsx` — super-admin login wordmark
- `components/shared/Sidebar.tsx` — fallback shop name + footer label
- `components/billing/BillingRowItem.tsx` — out-of-stock disabled state
- `types/index.ts` — `isOutOfStock` field
- `package.json` — version bump 1.4.3 → 1.5.0

**Landing:**
- `landing/app/page.tsx` — `APP_URL` → `https://app.sellos.in`; R→S badge letters in navbar and footer; wordmarks; all body copy; stats version badge v1.3 → v1.5
- `landing/app/layout.tsx` — `BASE_URL` → `https://sellos.in`; SEO, OG, and JSON-LD metadata

**Backend:**
- `cmd/server/main.go` — boot log line
- `internal/email/smtp.go` — super-admin OTP email subject and body
- `internal/queries/products.sql` — stock aggregate in search query
- `internal/generated/products.sql.go` — regenerated sqlc output

---

## Technical Notes

- Frontend ESLint and TypeScript (`tsc --noEmit`) pass cleanly on both `frontend/` and `landing/`.
- Backend `go test ./...` passes all packages except a pre-existing failure in `TestLookupCustomer_PhoneValidation` — verified to fail identically against `main` with this branch's changes stashed. Unrelated to this release; tracked as a follow-up.
- Vercel auto-deploys both landing and app on merge to main; Render auto-deploys the backend.

---

## Known Considerations

- **OG image asset still RetailOS-branded.** `/og-image.png` in `frontend/public/` and `landing/public/` is unchanged. Social previews on Twitter, Facebook, and LinkedIn will continue showing the old image until the asset is replaced — separate design task.
- **Social-card cache.** Twitter, Facebook, and LinkedIn cache OG previews for hours to days. After deploy, force re-fetch via each platform's card debugger.
- **Google re-crawl.** Search results may continue showing "RetailOS" in titles for one to two weeks until re-indexing completes. Self-resolving.
- **Render env vars.** `SMTP_FROM` and `RESEND_FROM` display names must be updated to "SellOS" in the Render dashboard for production OTP emails. Local `.env` change in this release is dev-only.
- **`JWT_SECRET` and admin-key rotation deliberately out of scope.** Rotation will log out every active session — must be coordinated at a low-traffic time, separately from this release.
- **Pre-existing handler test failure.** `TestLookupCustomer_PhoneValidation` panics on a nil pgxpool connection because the test does not mock the DB pool. Pre-existing, unrelated to this release — follow-up to fix.

---

## Test Plan

### Out-of-stock billing search
- [ ] Open billing screen, search for a product with at least one in-stock batch — row appears selectable as normal
- [ ] Search for a product whose batches are all depleted — row is visually marked as out of stock and cannot be selected
- [ ] Existing in-stock workflow still completes a bill end-to-end

### SellOS rebrand — visual
- [ ] `https://sellos.in` — navbar wordmark reads "SellOS"; badge shows "S"
- [ ] `https://sellos.in` — footer wordmark reads "SellOS"; badge shows "S"
- [ ] `https://sellos.in` — body copy in every section references SellOS, not RetailOS
- [ ] `https://sellos.in` — stats badge displays v1.5 — June 2026
- [ ] `https://app.sellos.in/login` — login wordmark reads "SellOS"
- [ ] `https://app.sellos.in/super-admin/login` — super-admin login wordmark reads "SellOS"
- [ ] Log in to the app — sidebar fallback shop name reads "SellOS"; sidebar footer label reads "SellOS"

### SellOS rebrand — system
- [ ] Trigger a super-admin OTP — email subject reads "SellOS Super Admin Login OTP" (after Render env-var update)
- [ ] Render backend logs show "SellOS API listening on :8080" after deploy
- [ ] `view-source:https://sellos.in` — `<title>`, OG `siteName`, JSON-LD `name` all read SellOS
- [ ] `view-source:https://app.sellos.in` — same SEO/OG metadata reads SellOS

### Regression
- [ ] Bill creation flow works end-to-end
- [ ] OTP login for super admin still works
- [ ] No new console errors or 404s introduced

### Version
- [ ] Confirm `frontend/package.json` shows `1.5.0`
- [ ] Confirm landing stats badge displays `v1.5`

---

**SellOS v1.5.0** — A cleaner billing search, and the brand the product was always meant to wear.
