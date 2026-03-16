# PROJECT_STATUS.md

## Current Objective

- Complete STEP 1 auth, session, and user management foundation on top of the existing VPS-ready runtime.

## Current Phase

- Phase 1 — Auth + Session + User Management Foundation

## Summary Status

- Observed: monorepo baseline now exists with `frontend/`, `backend/`, root env/compose/readme, and initial database migration.
- Observed: frontend now verifies the access token server-side through `GET /api/v1/auth/me` instead of trusting copied identity cookies.
- Observed: frontend includes a protected `/dashboard` route, dashboard summary cards, and a `/members` list screen with empty/error/access-denied states.
- Observed: Tailwind global stylesheet is now imported from the root layout, so utility classes and custom component layers render instead of falling back to browser defaults.
- Observed: login and app shell now share a light mobile-first UI baseline with reusable card, button, input, badge, helper, and empty-state styles.
- Observed: frontend install/runtime scripts now run `svelte-kit sync` before `dev`/`build`/`preview` and after `npm install`, preventing the missing `./.svelte-kit/tsconfig.json` warning on a fresh checkout.
- Observed: the PWA baseline now uses a static web manifest plus a native SvelteKit service worker instead of the Workbox plugin chain that emitted deprecated package warnings during `npm install` on Node 24.
- Observed: backend now resolves authenticated users from the database on each protected request, exposes `GET /api/v1/users`, `POST /api/v1/users`, and `PATCH /api/v1/users/:id`, and keeps `/health` plus `/api/v1/health` available.
- Observed: backend now has a follow-up migration for users foundation hardening and a seed admin flow that creates the initial admin only when it does not already exist.
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
- STEP 1 auth/session foundation is now connected end-to-end: login stores a single access-token cookie, SvelteKit verifies it through `auth/me`, and protected routes now rely on verified server-side user state.
- Member management foundation is now available: users schema is hardened, admin seed is non-duplicating, and backend/frontend now support the initial members list flow.
- Dashboard now has a meaningful placeholder with auth status, current role, and member count summary when the user role is allowed to read the members API.

## In Progress

- Domain modules beyond auth and members are still placeholder-only on the frontend.
- VPS build/deploy has not been executed from this workspace yet.
- Admin create/edit UI for members is still backend-only; frontend currently ships the members list only.

## Blockers / Risks

- Migration is manual; schema application is not yet automated.
- Production deploy has not been validated with live GAS build or Nginx apply yet.
- Residual: `frontend/npm audit` still reports 3 low severity vulnerabilities from SvelteKit's current `cookie@0.6.0` dependency chain; no safe local override has been applied yet.
- Risk: scope can expand too early if phase order is not enforced.
- Risk: `frontend/.env` now needs a valid `PRIVATE_API_BASE_URL` for server-side auth and data loads outside the Nginx split runtime.

## Recently Touched Areas

- `.gitignore`
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
- `backend/internal/handlers/user_handler.go`
- `backend/internal/response/json.go`
- `backend/internal/services/user_service.go`
- `frontend/src/lib/api/server.ts`
- `frontend/src/lib/api/types.ts`
- `frontend/src/routes/dashboard/+page.server.ts`
- `frontend/src/routes/dashboard/+page.svelte`
- `frontend/src/routes/members/+page.server.ts`

## Assumptions / Unknowns

- Assumption: `messhub-masterplan.md` is the current approved product source.
- Assumption: stack choices in the masterplan are accepted for the bootstrap.
- Assumption: deploy target remains VPS Linux with PM2 and Nginx split frontend/backend.
- Assumption: frontend can reach the backend from the server runtime through `PRIVATE_API_BASE_URL`, defaulting to `http://127.0.0.1:4100/api/v1`.
- Unknown: final public domain and SSL mode for production deploy.
- Unknown: preferred migration tool and CI workflow.

## Next Recommended Steps

- Apply `backend/db/migrations/002_auth_foundation.sql` after the initial schema migration, then run the seed admin command.
- Validate the updated runtime on the VPS with the new `PRIVATE_API_BASE_URL` frontend env.
- Continue to STEP 2 with wallet transaction foundation, reusing the verified auth/session and user role checks.
- Re-check the `cookie` advisory after the next SvelteKit release before applying any auth-related dependency override.
- Record any stack or architecture changes in `docs/decisions.md`.
