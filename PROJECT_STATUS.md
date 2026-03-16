# PROJECT_STATUS.md

## Current Objective

- Complete STEP 3 wifi billing and audit log foundation on top of the existing auth, session, member, and wallet runtime.

## Current Phase

- Phase 3 — Wifi Billing and Audit Log Foundation

## Summary Status

- Observed: monorepo baseline now exists with `frontend/`, `backend/`, root env/compose/readme, and initial database migration.
- Observed: GitHub Actions deployment automation is now defined for pushes to `main`, using SSH into the VPS, `git pull --ff-only origin main`, GAS CLI rebuilds for backend/frontend, and a backend health check.
- Observed: frontend now verifies the access token server-side through `GET /api/v1/auth/me` instead of trusting copied identity cookies.
- Observed: frontend includes a protected `/dashboard` route, dashboard summary cards, and a `/members` list screen with empty/error/access-denied states.
- Observed: Tailwind global stylesheet is now imported from the root layout, so utility classes and custom component layers render instead of falling back to browser defaults.
- Observed: login and app shell now share a light mobile-first UI baseline with reusable card, button, input, badge, helper, and empty-state styles.
- Observed: frontend install/runtime scripts now run `svelte-kit sync` before `dev`/`build`/`preview` and after `npm install`, preventing the missing `./.svelte-kit/tsconfig.json` warning on a fresh checkout.
- Observed: the PWA baseline now uses a static web manifest plus a native SvelteKit service worker instead of the Workbox plugin chain that emitted deprecated package warnings during `npm install` on Node 24.
- Observed: backend now resolves authenticated users from the database on each protected request, exposes `GET /api/v1/users`, `POST /api/v1/users`, and `PATCH /api/v1/users/:id`, and keeps `/health` plus `/api/v1/health` available.
- Observed: backend now has a follow-up migration for users foundation hardening and a seed admin flow that creates the initial admin only when it does not already exist.
- Observed: backend wallet foundation is now live with `GET /api/v1/wallet`, `GET /api/v1/wallet/transactions`, and `POST /api/v1/wallet/transactions`, plus a follow-up migration that normalizes the wallet table to the STEP 2 transaction shape.
- Observed: wallet reads are now available to all authenticated roles, while wallet transaction creation is limited to `admin` and `treasurer`.
- Observed: frontend now includes a real `/wallet` page with balance cards and paginated transactions, a `/wallet/new` form for authorized roles, and a dashboard wallet summary card instead of a placeholder.
- Observed: backend now has STEP 3 wifi billing endpoints for create/list/detail/active/my/submit/verify/reject, a follow-up migration for wifi schema hardening, and reusable audit logging wired into wifi, wallet transaction creation, and role/activation changes on users.
- Observed: frontend now includes a role-aware `/wifi` page for monthly bill creation, payment proof submission, verification, rejection, status history, and a dashboard wifi summary card.
- Observed: frontend is now prepared for PM2 runtime with `adapter-node`, `start` script, local `.env.example`, and `ecosystem.config.cjs`.
- Observed: backend now prefers `PORT` and has its own `.env.example` for GAS/PM2 usage.
- Observed: Docker has been reduced to local Postgres only.

## Done

- Masterplan exists in `messhub-masterplan.md`.
- Shared context baseline created: `AGENTS.md`, `PROJECT_STATUS.md`, `TASKS.md`, `docs/decisions.md`, `docs/handoffs/HANDOFF_TEMPLATE.md`.
- Agent runtime ignore baseline created in `.gitignore`.
- Monorepo root created with `.env.example`, `docker-compose.yml`, and `README.md`.
- Frontend scaffold created with SvelteKit, Tailwind config, native PWA baseline, auth flow baseline, and placeholder routes.
- Frontend baseline UI polished with shared global styles, refined login screen, and a cleaner mobile-first AppShell.
- Frontend install/runtime path is now stable on Node 24, with deprecated Workbox install warnings removed and `.svelte-kit` sync automated for fresh installs.
- Backend scaffold created with Go Fiber, env config, PostgreSQL bootstrap, `/api/v1`, auth endpoints, middleware, migration SQL, and seed admin.
- Runtime/deploy adaptation completed for separated frontend/backend services on ports `4101` and `4100`.
- GitHub Actions CI/CD baseline added for `main` branch deploys to the VPS through SSH, reusing the existing GAS CLI + PM2 runtime commands and backend health check.
- STEP 1 auth/session foundation is now connected end-to-end: login stores a single access-token cookie, SvelteKit verifies it through `auth/me`, and protected routes now rely on verified server-side user state.
- Member management foundation is now available: users schema is hardened, admin seed is non-duplicating, and backend/frontend now support the initial members list flow.
- Dashboard now has a meaningful placeholder with auth status, current role, and member count summary when the user role is allowed to read the members API.
- STEP 2 wallet foundation is now available end-to-end: database migration, backend summary/list/create flows, role-aware authorization, wallet pages, and dashboard balance summary are connected.
- STEP 3 wifi billing foundation is now available end-to-end: migration, backend billing/proof/review flows, dashboard wifi summary, frontend `/wifi` screen, and reusable audit logging are connected.

## In Progress

- Domain modules beyond auth, members, wallet, and wifi are still placeholder-only or read-only placeholders on the frontend.
- The new GitHub Actions deploy workflow has been authored, but it has not been exercised against the live VPS from this workspace yet.
- Admin create/edit UI for members is still backend-only; frontend currently ships the members list only.

## Blockers / Risks

- Migration is manual; schema application is not yet automated.
- Wallet rollout now depends on applying `backend/db/migrations/003_wallet_step2.sql` after the existing migrations in each environment.
- Wifi rollout now depends on applying `backend/db/migrations/004_wifi_audit_step3.sql` after the existing migrations in each environment.
- Production deploy has not been validated with live GAS build or Nginx apply yet.
- Risk: GitHub Actions deploy depends on repository secrets plus non-interactive `git pull` access already working on the VPS checkout.
- Residual: `frontend/npm audit` still reports 3 low severity vulnerabilities from SvelteKit's current `cookie@0.6.0` dependency chain; no safe local override has been applied yet.
- Risk: scope can expand too early if phase order is not enforced.
- Risk: `frontend/.env` now needs a valid `PRIVATE_API_BASE_URL` for server-side auth and data loads outside the Nginx split runtime.

## Recently Touched Areas

- `.gitignore`
- `.github/workflows/deploy.yml`
- `.env.example`
- `frontend/.env.example`
- `frontend/ecosystem.config.cjs`
- `frontend/src/app.css`
- `frontend/src/app.html`
- `frontend/src/service-worker.ts`
- `backend/.air.toml`
- `frontend/src/lib/components/AppShell.svelte`
- `frontend/src/lib/components/PageCard.svelte`
- `frontend/src/routes/+layout.svelte`
- `frontend/src/routes/+page.svelte`
- `frontend/src/routes/login/+page.svelte`
- `frontend/src/routes/login/+page.server.ts`
- `frontend/static/manifest.webmanifest`
- `frontend/package.json`
- `frontend/package-lock.json`
- `frontend/tailwind.config.ts`
- `backend/.env.example`
- `docker-compose.yml`
- `README.md`
- `AGENTS.md`
- `PROJECT_STATUS.md`
- `TASKS.md`
- `frontend/`
- `backend/`
- `docs/decisions.md`
- `docs/handoffs/HANDOFF_TEMPLATE.md`
- `backend/db/migrations/002_auth_foundation.sql`
- `backend/db/migrations/003_wallet_step2.sql`
- `backend/db/migrations/004_wifi_audit_step3.sql`
- `backend/internal/handlers/user_handler.go`
- `backend/internal/handlers/wallet_handler.go`
- `backend/internal/handlers/wifi_handler.go`
- `backend/internal/response/json.go`
- `backend/internal/repository/audit_repository.go`
- `backend/internal/repository/wallet_repository.go`
- `backend/internal/repository/wifi_repository.go`
- `backend/internal/services/wallet_service.go`
- `backend/internal/services/audit_service.go`
- `backend/internal/services/user_service.go`
- `backend/internal/services/wifi_service.go`
- `frontend/src/lib/api/client.ts`
- `frontend/src/lib/api/server.ts`
- `frontend/src/lib/api/types.ts`
- `frontend/src/routes/dashboard/+page.server.ts`
- `frontend/src/routes/dashboard/+page.svelte`
- `frontend/src/routes/members/+page.server.ts`
- `frontend/src/routes/wallet/+page.server.ts`
- `frontend/src/routes/wallet/+page.svelte`
- `frontend/src/routes/wallet/new/+page.server.ts`
- `frontend/src/routes/wallet/new/+page.svelte`
- `frontend/src/routes/wifi/+page.server.ts`
- `frontend/src/routes/wifi/+page.svelte`

## Assumptions / Unknowns

- Assumption: `messhub-masterplan.md` is the current approved product source.
- Assumption: stack choices in the masterplan are accepted for the bootstrap.
- Assumption: deploy target remains VPS Linux with PM2 and Nginx split frontend/backend.
- Assumption: frontend can reach the backend from the server runtime through `PRIVATE_API_BASE_URL`, defaulting to `http://127.0.0.1:4100/api/v1`.
- Unknown: final public domain and SSL mode for production deploy.
- Unknown: preferred migration tool and CI workflow.
- Assumption: the VPS checkout at `/home/ubuntu/projects/messhub` already has Git remote credentials configured so `git pull --ff-only origin main` succeeds non-interactively.

## Next Recommended Steps

- Apply `backend/db/migrations/002_auth_foundation.sql` and `backend/db/migrations/003_wallet_step2.sql`, then run the seed admin command if the environment is still new.
- Validate the wallet flow against a real Postgres database by creating sample income and expense transactions through the new frontend form.
- Apply `backend/db/migrations/004_wifi_audit_step3.sql`, then validate create bill, submit proof, and verify/reject flows against the live Postgres runtime.
- Continue to STEP 4 modules only after confirming the wifi flow and audit log entries on the VPS.
- Add the required `VPS_HOST`, `VPS_USER`, and `VPS_SSH_KEY` repository secrets in GitHub, then run the workflow with a test push to `main`.
- Validate the updated runtime on the VPS with the new `PRIVATE_API_BASE_URL` frontend env.
- Re-check the `cookie` advisory after the next SvelteKit release before applying any auth-related dependency override.
- Record any stack or architecture changes in `docs/decisions.md`.
