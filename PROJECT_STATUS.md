# PROJECT_STATUS.md

## Current Objective

- Implement cross-app UX polish, navigation cleanup, title/meta consistency, Iconify adoption, and auth identifier login (`email` or `username`) on top of the existing production-ready auth, member, wallet, wifi, smart mess, profile, settings, import, and PWA runtime.

## Current Phase

- Phase 10 — UX Polish, Navigation, and Auth Refinement

## Summary Status

- Observed: protected SvelteKit server actions now resolve the current user through a shared server-side auth helper before checking role-sensitive mutations, fixing the false "Sesi login tidak ditemukan" error on settings saves and other admin/member actions while still redirecting expired sessions back to `/login`.
- Observed: the shared app shell now uses navigation-start aware route state, closes menus consistently on click, and keeps the desktop sidebar bounded to the viewport with its own scroll area so the sidebar no longer bleeds down over the content column.
- Observed: `/members` is now a fuller admin workspace with search/filter controls, member create/edit sheets, activate/deactivate actions, role/status updates, admin password reset, and clearer empty/error feedback on both desktop and mobile.
- Observed: backend user management now supports richer editable member fields plus `PATCH /api/v1/users/:id/password`, adds create/update/reset audit logs, and protects against self-demotion, self-deactivation, and removing the last active admin.
- Observed: empty-state UX is now more consistent across members, wallet, feed, notifications, contributions, and settings through a richer shared `StatePanel` and reusable action-sheet surface.
- Observed: frontend now uses a shared navigation/meta source for route-aware titles, menu visibility, active states, and browser tab metadata, with consistent `• MessHub` titles for login, dashboard, wallet, wifi, feed, members, profile, settings, admin, and nested import routes.
- Observed: frontend authenticated UX now has a warm interior-inspired theme with reusable CSS tokens, Iconify icons, a new responsive shell with desktop sidebar + mobile bottom navigation, safer profile dropdown behavior, smoother page transitions, and richer feedback/loading/skeleton states across the main routes.
- Observed: login now accepts email or username on the frontend, with a refined identifier-based form, better submit states, and production-facing copy that no longer reads like helper/dev text.
- Observed: backend now supports login by `identifier` (`email` or `username`), includes `username` in auth/session payloads, auto-generates unique usernames for new/seed/imported users, and ships a safe `009_username_auth_step10.sql` migration that backfills existing members before enforcing uniqueness.
- Observed: production-facing copy across login, dashboard, members, wallet, profile, settings, the shared app shell, and the global error boundary has been cleaned up so the UI no longer exposes seed-account, `.env`, backend-auth, or phase-internal helper text, while admin-facing messaging is now more user-friendly.
- Observed: backend now has `backend/db/migrations/008_import_step9.sql`, an `import_jobs` tracking table, wallet `transaction_date` + `proof_url` + `source` metadata, and admin-only CSV preview/commit endpoints for `/api/v1/import/members/...` and `/api/v1/import/wallet/...`, with audit logs for preview and commit actions.
- Observed: frontend now includes an admin import hub at `/admin/import`, dedicated `/admin/import/members` and `/admin/import/wallet` flows, downloadable CSV templates, mobile-friendly preview cards, duplicate/validation summaries, and commit guards before data is persisted.
- Observed: `docs/data-import.md` now documents Google Sheets CSV export steps, legacy spreadsheet-to-template column mapping, member temporary-password handling, sample preview payloads, and the rule that wallet balances are recalculated from imported transactions instead of spreadsheet saldo values.
- Observed: frontend manifest now targets a standalone `/dashboard` start experience, includes install-ready `192x192` and `512x512` PNG icons, safe-area mobile meta tags, and a conditional `Install MessHub App` prompt in the authenticated shell.
- Observed: frontend now has a richer native PWA runtime with versioned service-worker caches for static assets plus safe dashboard/feed/wallet/wifi routes, an `/offline` fallback page, Web Push notification handlers, notification-click routing, logout-time runtime cache clearing, and a lightweight IndexedDB outbox for background sync of offline activity/comment submissions.
- Observed: frontend mobile UX now includes an upgraded bottom navigation for `Dashboard`, `Feed`, `Wallet`, `Wifi`, and `Profile`, pull-to-refresh interactions on dashboard/feed, an offline-mode banner, a push-permission prompt, and clearer outbox/offline messaging when feed writes are queued.
- Observed: backend now supports authenticated browser/service-worker calls through the existing `mh_access_token` cookie in addition to bearer tokens, includes `backend/db/migrations/007_push_step8.sql` plus a `push_subscriptions` table and `/api/v1/push/subscribe` + `/api/v1/push/unsubscribe`, and fans out Web Push delivery for wifi bill creation, wifi payment verification, new activities, and new comments.
- Observed: backend now applies request-scoped `X-Request-ID`, structured JSON request logging, JSON panic recovery, API security headers, and route-level rate limiting for login plus high-write feed endpoints.
- Observed: backend error responses are now standardized to `{ error, message, details? }`, while request validation now covers login, users, wallet creation, wifi bill/proof/review payloads, activity/comment/reaction writes, profile updates, password changes, settings updates, and notifications read payloads.
- Observed: backend role enforcement is now hardened through the shared role middleware on sensitive routes, `/api/v1/system/status` is admin-only, and `/health` now returns readiness-oriented app/database/version/time data with `503` when the database is unreachable.
- Observed: frontend API helpers now parse the hardened backend error payload, preserve backend request IDs for debugging, distinguish network failures, and redirect to `/login` when authenticated route loads/actions receive `401`.
- Observed: frontend now has a shared `StatePanel` for loading/error/empty/forbidden UX across dashboard, members, wallet, wifi, feed, profile, settings, and wallet-create flows, plus a global `+error.svelte` boundary instead of blank/default crash output.
- Observed: backend now has `backend/db/migrations/006_settings_profile_step6.sql`, extending `users` with `phone` and `avatar_url` plus a singleton `mess_settings` table with wifi and bank defaults.
- Observed: backend now exposes authenticated `GET /api/v1/profile`, `PATCH /api/v1/profile`, `PATCH /api/v1/profile/password`, and `GET /api/v1/settings`, plus admin-only `GET /api/v1/system/status` and `PATCH /api/v1/settings`.
- Observed: profile, password-change, settings, role-change, and activation flows now record the required audit actions, and API errors now use a shared `{ error, message, details? }` contract from the response helper.
- Observed: frontend `/profile` is now a live account page for editing name, phone, avatar URL, and password, while `/settings` is now an admin-only configuration screen with live system status.
- Observed: frontend `/members` is no longer read-only for admins; it now supports role updates and activate/deactivate controls, while treasurers remain read-only viewers.
- Observed: wifi bill creation defaults now read from live mess settings instead of duplicating the initial `Rp20.000` / day-10 defaults only in the wifi page form.
- Observed: backend now has `backend/db/migrations/005_smart_mess_step5.sql`, introducing the Step 5 smart mess schema for `activities`, comments, reactions, food claims, rice responses, and normalized notifications.
- Observed: backend now exposes authenticated activity feed APIs for list/create/comment/reaction plus food claim and rice response actions, contribution leaderboard APIs, and notifications list/read APIs.
- Observed: wifi bill creation and wifi payment verification now generate in-app notifications, while new activities and new comments also fan out notifications to active members.
- Observed: audit logging now records `food_claim`, `rice_response`, and `notification_read` in addition to the earlier wifi, wallet, and user-audit actions.
- Observed: frontend now includes an interactive `/feed` page, `/contributions` leaderboard page, `/notifications` inbox, dashboard leaderboard section, and a notification badge in the shared app shell header.
- Observed: monorepo baseline now exists with `frontend/`, `backend/`, shared docs, and the database migration set tracked under `backend/db/migrations/`.
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
- Observed: service-scoped `.env.example` files under `frontend/` and `backend/` are the active checked-in env references for local and VPS runtime.

## Done

- Cross-app UX polish, navigation cleanup, title/meta consistency, Iconify adoption, warm-theme refresh, and username login support are now implemented across the main authenticated shell and key frontend/backend auth paths.
- STEP 9 data import, UI cleanup, and migration tooling is now implemented end-to-end: admin-only CSV preview/commit APIs for members and wallet, import job tracking plus audit logs, import templates and docs, admin import hub/pages, wallet history date/proof preservation, and production UI copy cleanup on key pages.
- Masterplan exists in `messhub-masterplan.md`.
- Shared context baseline created: `AGENTS.md`, `PROJECT_STATUS.md`, `TASKS.md`, `docs/decisions.md`, `docs/handoffs/HANDOFF_TEMPLATE.md`.
- Agent runtime ignore baseline created in `.gitignore`.
- Monorepo root created with `README.md`, shared project context files, and service-scoped env examples under `frontend/` and `backend/`.
- Frontend scaffold created with SvelteKit, Tailwind config, native PWA baseline, auth flow baseline, and placeholder routes.
- Frontend baseline UI polished with shared global styles, refined login screen, and a cleaner mobile-first AppShell.
- Frontend install/runtime path is now stable on Node 24, with deprecated Workbox install warnings removed and `.svelte-kit` sync automated for fresh installs.
- Backend scaffold created with Go Fiber, env config, PostgreSQL bootstrap, `/api/v1`, auth endpoints, middleware, migration SQL, and seed admin.
- Runtime/deploy adaptation completed for separated frontend/backend services on ports `4101` and `4100`.
- GitHub Actions CI/CD baseline added for `main` branch deploys to the VPS through SSH, reusing the existing GAS CLI + PM2 runtime commands and backend health check.
- STEP 1 auth/session foundation is now connected end-to-end: login stores a single access-token cookie, SvelteKit verifies it through `auth/me`, and protected routes now rely on verified server-side user state.
- Member management foundation is now available: users schema is hardened, admin seed is non-duplicating, and backend/frontend now support the initial members list flow.
- Dashboard now has meaningful summary cards for member counts, wallet, wifi, and contribution leaderboard data when the current role is allowed to read the supporting APIs.
- STEP 2 wallet foundation is now available end-to-end: database migration, backend summary/list/create flows, role-aware authorization, wallet pages, and dashboard balance summary are connected.
- STEP 3 wifi billing foundation is now available end-to-end: migration, backend billing/proof/review flows, dashboard wifi summary, frontend `/wifi` screen, and reusable audit logging are connected.
- STEP 5 smart mess features are now available end-to-end: unified activity feed, comments/reactions, contribution leaderboard, food claim, rice response, notifications UI, and notification triggers for activity/comment/wifi events.
- STEP 6 system settings, profile management, and admin controls are now available end-to-end: backend settings/profile/system status APIs, standardized API errors, frontend `/profile` and `/settings`, and admin member role/activation controls are connected.
- STEP 7 production hardening and reliability is now available end-to-end: backend request validation, error standardization, request IDs, structured logs, rate limiting, security headers, readiness-aware health checks, admin-only system status, and frontend unauthorized/forbidden/error/loading containment are connected.
- STEP 8 PWA upgrade and mobile experience is now implemented in the repo: installable manifest/icons, install prompt UI, versioned service worker, offline fallback, safe runtime caching, feed outbox background sync, Web Push subscription/delivery plumbing, mobile bottom navigation, and dashboard/feed pull-to-refresh.

## In Progress

- Live validation with real exported Google Sheets/legacy CSV files is still pending, including final confirmation of member duplicate handling and wallet category inference against actual mess data.
- Live browser/device validation for Step 8 is still pending for standalone install, Android push delivery, offline cache behavior after logout/login, and background-sync replay after reconnect.
- Shared expenses and proposals are still placeholder-only on the frontend and do not yet have live runtime integration.
- The new GitHub Actions deploy workflow has been authored, but it has not been exercised against the live VPS from this workspace yet.

## Blockers / Risks

- Step 9 rollout now depends on applying `backend/db/migrations/008_import_step9.sql` after the earlier migrations in each environment before the new import screens or wallet transaction metadata can be used safely.
- Migration is manual; schema application is not yet automated.
- Wallet rollout now depends on applying `backend/db/migrations/003_wallet_step2.sql` after the existing migrations in each environment.
- Wifi rollout now depends on applying `backend/db/migrations/004_wifi_audit_step3.sql` after the existing migrations in each environment.
- Smart mess rollout now depends on applying `backend/db/migrations/005_smart_mess_step5.sql` after the earlier migrations in each environment.
- Step 6 rollout now depends on applying `backend/db/migrations/006_settings_profile_step6.sql` after the earlier migrations in each environment.
- Step 8 rollout now depends on applying `backend/db/migrations/007_push_step8.sql`, setting `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, `VAPID_SUBJECT` on the backend, and matching `PUBLIC_PUSH_VAPID_PUBLIC_KEY` on the frontend.
- Production rollout should set `CORS_ORIGIN` to the real frontend origin list and replace the default `JWT_SECRET`, `APP_VERSION`, and optional `CONTENT_SECURITY_POLICY` values before exposing the hardened runtime publicly.
- Production deploy has not been validated with live GAS build or Nginx apply yet.
- Risk: GitHub Actions deploy depends on repository secrets plus non-interactive `git pull` access already working on the VPS checkout.
- Risk: Web Push, homescreen installability, and service-worker background sync still need live validation on HTTPS/Android Chrome because they cannot be fully proven from static builds alone.
- Risk: Authenticated offline caches are cleared on logout, but same-device account switching and stale cached pages should still be regression-checked in the browser after deploy.
- Residual: `frontend/npm audit` still reports 3 low severity vulnerabilities from SvelteKit's current `cookie@0.6.0` dependency chain; no safe local override has been applied yet.
- Risk: shared expenses and proposals still remain out of sequence relative to the original phase order; repo planning must keep that explicit rather than implying they are already done.
- Risk: `frontend/.env` now needs a valid `PRIVATE_API_BASE_URL` for server-side auth and data loads outside the Nginx split runtime.
- Assumption risk: avatar is currently stored as a string URL/reference only; no upload pipeline or image storage flow was added in this step.

## Recently Touched Areas

- frontend/src/lib/auth/server.ts
- frontend/src/lib/components/ActionSheet.svelte
- frontend/src/lib/components/StatePanel.svelte
- frontend/src/lib/components/AppShell.svelte
- frontend/src/app.css
- frontend/src/routes/login/
- frontend/src/routes/dashboard/
- frontend/src/routes/wallet/
- frontend/src/routes/wifi/
- frontend/src/routes/feed/
- frontend/src/routes/members/
- frontend/src/routes/profile/
- frontend/src/routes/settings/
- backend/internal/services/auth_service.go
- backend/internal/handlers/auth_handler.go
- backend/internal/repository/user_repository.go
- backend/internal/models/models.go
- backend/db/migrations/
- `backend/db/migrations/008_import_step9.sql`
- `backend/internal/handlers/import_handler.go`
- `backend/internal/repository/import_job_repository.go`
- `backend/internal/services/import_service.go`
- `docs/data-import.md`
- `frontend/src/routes/admin/import/`
- `frontend/static/templates/`
- `.gitignore`
- `.github/workflows/deploy.yml`
- `frontend/.env.example`
- `frontend/ecosystem.config.cjs`
- `frontend/src/app.css`
- `frontend/src/app.html`
- `frontend/src/service-worker.ts`
- `frontend/static/icons/icon-192.png`
- `frontend/static/icons/icon-512.png`
- `backend/.air.toml`
- `frontend/src/lib/components/AppShell.svelte`
- `frontend/src/lib/components/PwaControlBar.svelte`
- `frontend/src/lib/components/PullToRefresh.svelte`
- `frontend/src/lib/components/PageCard.svelte`
- `frontend/src/routes/+layout.svelte`
- `frontend/src/routes/offline/+page.svelte`
- `frontend/src/routes/+page.svelte`
- `frontend/src/routes/login/+page.svelte`
- `frontend/src/routes/login/+page.server.ts`
- `frontend/static/manifest.webmanifest`
- `frontend/package.json`
- `frontend/package-lock.json`
- `frontend/tailwind.config.ts`
- `backend/.env.example`
- `backend/cmd/api/main.go`
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
- `backend/db/migrations/005_smart_mess_step5.sql`
- `backend/db/migrations/006_settings_profile_step6.sql`
- `backend/db/migrations/007_push_step8.sql`
- `backend/internal/app/app.go`
- `backend/internal/handlers/profile_handler.go`
- `backend/internal/handlers/health_handler.go`
- `backend/internal/handlers/settings_handler.go`
- `backend/internal/handlers/system_handler.go`
- `backend/internal/handlers/push_handler.go`
- `backend/internal/handlers/user_handler.go`
- `backend/internal/handlers/wallet_handler.go`
- `backend/internal/handlers/wifi_handler.go`
- `backend/internal/handlers/activity_handler.go`
- `backend/internal/handlers/notification_handler.go`
- `backend/internal/handlers/validation_helpers.go`
- `backend/internal/middleware/rate_limit.go`
- `backend/internal/middleware/recovery.go`
- `backend/internal/middleware/request_context.go`
- `backend/internal/middleware/request_logger.go`
- `backend/internal/middleware/security.go`
- `backend/internal/observability/logger.go`
- `backend/internal/response/json.go`
- `backend/internal/repository/audit_repository.go`
- `backend/internal/repository/activity_repository.go`
- `backend/internal/repository/notification_repository.go`
- `backend/internal/repository/push_subscription_repository.go`
- `backend/internal/repository/settings_repository.go`
- `backend/internal/repository/wallet_repository.go`
- `backend/internal/repository/wifi_repository.go`
- `backend/internal/services/wallet_service.go`
- `backend/internal/services/activity_service.go`
- `backend/internal/services/audit_service.go`
- `backend/internal/services/notification_service.go`
- `backend/internal/services/push_service.go`
- `backend/internal/services/settings_service.go`
- `backend/internal/services/system_service.go`
- `backend/internal/services/user_service.go`
- `backend/internal/services/wifi_service.go`
- `backend/internal/validation/validation.go`
- `frontend/src/lib/api/client.ts`
- `frontend/src/lib/api/http.ts`
- `frontend/src/lib/api/server.ts`
- `frontend/src/lib/api/types.ts`
- `frontend/src/lib/auth/guard.ts`
- `frontend/src/lib/auth/session.ts`
- `frontend/src/lib/components/StatePanel.svelte`
- `frontend/src/lib/pwa/`
- `frontend/src/lib/server/api-errors.ts`
- `frontend/src/routes/+error.svelte`
- `frontend/src/routes/dashboard/+page.server.ts`
- `frontend/src/routes/dashboard/+page.svelte`
- `frontend/src/routes/feed/+page.server.ts`
- `frontend/src/routes/feed/+page.svelte`
- `frontend/src/routes/contributions/+page.server.ts`
- `frontend/src/routes/contributions/+page.svelte`
- `frontend/src/routes/notifications/+page.server.ts`
- `frontend/src/routes/notifications/+page.svelte`
- `frontend/src/routes/members/+page.server.ts`
- `frontend/src/routes/profile/+page.server.ts`
- `frontend/src/routes/profile/+page.svelte`
- `frontend/src/routes/settings/+page.server.ts`
- `frontend/src/routes/settings/+page.svelte`
- `frontend/src/routes/wallet/+page.server.ts`
- `frontend/src/routes/wallet/+page.svelte`
- `frontend/src/routes/wallet/new/+page.server.ts`
- `frontend/src/routes/wallet/new/+page.svelte`
- `frontend/src/routes/wifi/+page.server.ts`
- `frontend/src/routes/wifi/+page.svelte`
- `frontend/src/hooks.server.ts`

## Assumptions / Unknowns

- Assumption: `messhub-masterplan.md` is the current approved product source.
- Assumption: stack choices in the masterplan are accepted for the bootstrap.
- Assumption: deploy target remains VPS Linux with PM2 and Nginx split frontend/backend.
- Assumption: frontend and backend continue to share one public origin with `/api/v1` routed through the same domain so browser/service-worker requests can reuse the auth cookie safely.
- Assumption: frontend can reach the backend from the server runtime through `PRIVATE_API_BASE_URL`, defaulting to `http://127.0.0.1:4100/api/v1`.
- Unknown: final public domain and SSL mode for production deploy.
- Unknown: preferred migration tool and CI workflow.
- Assumption: the VPS checkout at `/home/ubuntu/projects/messhub` already has Git remote credentials configured so `git pull --ff-only origin main` succeeds non-interactively.

## Next Recommended Steps

- Apply migrations through `backend/db/migrations/008_import_step9.sql`, then run one dry-run import with real member and kas CSV exports before using the new admin import flow in production.
- Validate Step 8 on a real mobile browser and homescreen install: install prompt, standalone opening, push permission/subscription, wifi/activity/comment push delivery, offline dashboard/feed loads, and activity/comment outbox replay after reconnect.
- Re-check logout/login behavior with the new authenticated runtime caches so offline pages from one user are not retained into another user session.
- Return to the still-pending shared expenses and proposals modules after Step 8 rollout is confirmed.
- Add the required `VPS_HOST`, `VPS_USER`, and `VPS_SSH_KEY` repository secrets in GitHub, then run the workflow with a test push to `main`.
- Validate the updated runtime on the VPS with the new `PRIVATE_API_BASE_URL` frontend env.
- Re-check the `cookie` advisory after the next SvelteKit release before applying any auth-related dependency override.
- Record any stack or architecture changes in `docs/decisions.md`.
