# decisions.md

## Decision 1
- Date: Unknown in source; normalized on 2026-03-16
- Context: The project needs a practical stack for a mobile-first internal PWA.
- Decision: Use SvelteKit + TailwindCSS for frontend, Go + Fiber for backend, PostgreSQL for database, and VPS + Nginx for deployment.
- Rationale: This stack is explicitly specified in the masterplan and matches the need for a lightweight, staged build.
- Impact: Initial repo structure, tooling, and implementation tasks should follow this stack unless superseded by a new decision.
- Follow-up: Confirm whether these choices are final before implementation starts.

## Decision 2
- Date: Unknown in source; normalized on 2026-03-16
- Context: Scope must stay realistic for v1.
- Decision: Keep v1 focused on auth, members, wallet, wifi billing and verification, shared expenses, contributions, simple feed, and PWA installability.
- Rationale: The masterplan explicitly excludes multi-mess, payment gateway, complex chat, AI, and deep accounting for the initial release.
- Impact: New work should be phase-gated; deferred items should not enter active tasks without a deliberate scope change.
- Follow-up: Re-evaluate deferred items after MVP is stable.

## Decision 3
- Date: Unknown in source; normalized on 2026-03-16
- Context: Access control affects both security and operational workflow.
- Decision: Accounts should not be self-registered freely; creation should be admin-driven or invite-based.
- Rationale: The masterplan marks this as safer for an internal, role-based system.
- Impact: Auth design, onboarding flow, and admin tooling must reflect controlled account creation.
- Follow-up: Decide the exact invite/onboarding mechanism during implementation.

## Decision 4
- Date: Unknown in source; normalized on 2026-03-16
- Context: Financial records cover different categories with different operational meanings.
- Decision: Keep wallet transactions, wifi billing, and shared non-cash expenses as separate domains.
- Rationale: The masterplan explicitly states that non-cash shared expenses must not reduce wallet balance, while wifi uses its own monthly billing workflow.
- Impact: Data model and UI must avoid merging these records into a single ambiguous ledger.
- Follow-up: Preserve this separation in schema and reporting.

## Decision 5
- Date: Unknown in source; normalized on 2026-03-16
- Context: The app is intended for daily use by residents on phones.
- Decision: Prioritize mobile-first UX and PWA installability from the start.
- Rationale: This is a core product requirement, not a later enhancement.
- Impact: Layout, routing, form flows, and performance work should optimize for Android phone usage first.
- Follow-up: Add installability and minimal offline shell support during Phase 1.

## Decision 6
- Date: 2026-03-16
- Context: Initial implementation needs a clean baseline without splitting the codebase into many repositories.
- Decision: Use a simple monorepo structure with `frontend/`, `backend/`, `docs/`, root `.env.example`, and root `docker-compose.yml`.
- Rationale: This keeps setup, local development, and context sharing simple while the product is still in early implementation.
- Impact: Shared configuration and bootstrap docs live at repo root; service-specific logic stays isolated in each folder.
- Follow-up: Revisit tooling only after the core modules are stable.

## Decision 7
- Date: 2026-03-16
- Context: Auth must be usable early without building full session infrastructure first.
- Decision: Start with email/password login on the backend using JWT, and use frontend cookie storage plus route guards for the initial web session.
- Rationale: This keeps auth operational and production-aligned enough for early module work without adding OAuth or session store complexity.
- Impact: Future auth work should preserve the existing role model while improving cookie security, refresh flow, and onboarding.
- Follow-up: Harden cookie settings and add refresh/logout invalidation strategy before production release.

## Decision 8
- Date: 2026-03-16
- Context: The project must fit an existing VPS workflow based on GAS CLI, PM2, and Nginx instead of Docker-first deployment.
- Decision: Make frontend and backend runnable as separate PM2 apps, standardize default ports to `4101` and `4100`, keep Docker only for local Postgres, and use per-service `.env.example` files.
- Rationale: This matches the established server workflow, avoids port collisions on the VPS, and removes ambiguity between local dev tooling and production deploy.
- Impact: Frontend now uses `adapter-node` plus `ecosystem.config.cjs`; backend now prefers `PORT`; deployment docs and env handling are centered on service directories instead of root Docker runtime.
- Follow-up: Validate the new workflow end-to-end with `gas build` and `gas deploy` on the target VPS.

## Decision 9
- Date: 2026-03-16
- Context: The frontend needs a usable baseline before dashboard and domain modules are implemented, but the codebase should avoid a large design-system rewrite.
- Decision: Use a light mobile-first visual baseline with shared global component classes in `frontend/src/app.css`, then layer route-specific UI improvements on the existing SvelteKit structure.
- Rationale: This fixes the immediate default-browser look, keeps Tailwind usage consistent, and provides reusable styling for later dashboard/member/wallet/wifi/feed work without introducing heavy abstraction.
- Impact: App shell, cards, buttons, inputs, helper boxes, badges, and empty states should reuse the shared classes first before adding one-off patterns.
- Follow-up: Apply the same baseline to dashboard summary cards and module list/detail screens as feature work continues.

## Decision 10
- Date: 2026-03-16
- Context: Fresh installs on Node 24 were showing deprecated dependency warnings from the Workbox-based PWA plugin chain, and `npm run dev` on a clean checkout could warn that `./.svelte-kit/tsconfig.json` was missing.
- Decision: Keep the PWA baseline with a static `manifest.webmanifest` and a native SvelteKit service worker, and standardize frontend scripts to run `svelte-kit sync` on install/dev/build/preview instead of depending on `@vite-pwa/sveltekit`.
- Rationale: This preserves baseline installability, removes deprecated install noise caused by the Workbox dependency chain, and makes the generated SvelteKit TypeScript config available on fresh installs.
- Impact: Future PWA work should extend the native service worker/manifest baseline first; reintroducing a plugin should only happen if it solves a concrete feature gap without reintroducing the Node 24 install warnings.
- Follow-up: Revisit richer offline caching or push-related features only when the actual product requirements exceed the native baseline.

## Decision 11
- Date: 2026-03-16
- Context: STEP 1 needs stable auth/session behavior without adding a separate session store or trusting user profile cookies copied into the browser.
- Decision: Keep a single `mh_access_token` httpOnly cookie on the frontend, then verify the current user server-side through `GET /api/v1/auth/me` on SvelteKit loads by calling the backend through `PRIVATE_API_BASE_URL`.
- Rationale: This keeps the runtime aligned with the existing split frontend/backend deploy, avoids trusting tamperable identity cookies, and gives later modules a verified current-user source without introducing refresh-token/session-store complexity yet.
- Impact: Protected frontend routes should fetch user-aware data from server loads/actions, and future domain work should reuse the same verified session flow instead of reading role/name/email directly from client-stored cookies.
- Follow-up: Add refresh/logout invalidation strategy before production hardening if session duration or device switching becomes more complex.

## Decision 12
- Date: 2026-03-16
- Context: Production deployment is already stable through manual SSH plus GAS CLI, but releases from `main` still require manual login to the VPS.
- Decision: Automate deploys with GitHub Actions on pushes to `main`, using SSH into the existing VPS checkout, `git pull --ff-only origin main`, the current GAS CLI build commands for backend and frontend, and a post-deploy backend health check.
- Rationale: This removes manual login from the release path while preserving the existing GAS CLI, PM2 app names, ports, Nginx split routing, and VPS directory structure that are already proven in the project.
- Impact: Repository secrets must provide `VPS_HOST`, `VPS_USER`, and `VPS_SSH_KEY`; the VPS checkout must already support non-interactive Git pulls; failed health checks should fail the workflow after the rebuild.
- Follow-up: Validate the workflow against the live VPS and decide later whether notifications or manual rollback helpers are needed.

## Decision 13
- Date: 2026-03-16
- Context: STEP 3 needs monthly wifi billing that stays separate from wallet accounting while remaining easy to verify from mobile and easy to audit on the backend.
- Decision: Keep one wifi bill per month-year, snapshot active users into `wifi_bill_members` at bill creation time, use a simple `proof_url`/reference field plus optional note for transfer proof submission, and record important mutations through a reusable audit helper inside service-level database transactions.
- Rationale: This matches the product rule that wifi is a monthly obligation, avoids mixing wallet and wifi ledgers, works even before a dedicated upload service exists, and keeps audit writes consistent with the main state changes.
- Impact: Future wifi work should build on the existing monthly bill/member snapshot model, and new auditable financial/member mutations should reuse the same audit helper pattern instead of bespoke logging.
- Follow-up: If file upload storage is added later, keep the `proof_url` contract stable so current submit/review flows do not need a breaking API change.

## Decision 14
- Date: 2026-03-16
- Context: STEP 5 needs contribution scoring, temporary food/rice posts, comments, reactions, and notification triggers, but the legacy `contributions` and `posts` tables from the initial schema were not connected to any live runtime.
- Decision: Implement smart mess interactions on top of a unified `activities` runtime with dedicated `activity_comments`, `activity_reactions`, `food_claims`, and `rice_responses` tables, and generate in-app notifications through a shared notification service inside domain transactions.
- Rationale: This matches the requested `/api/v1/activities/...` API shape, keeps leaderboard and feed logic on one canonical model, and avoids splitting social interactions across multiple half-used legacy tables.
- Impact: Future feed/contribution extensions should build on `activities` first, and notification-triggering domain events should reuse the shared notification service rather than inserting notification rows ad hoc.
- Follow-up: Decide later whether the unused legacy `contributions` and `posts` tables should be migrated into `activities` or removed once production data strategy is defined.

## Decision 15
- Date: 2026-03-16
- Context: STEP 6 introduces runtime-configurable mess defaults, user self-service profile flows, and more frontend/server action handling, so duplicated wifi defaults and inconsistent API errors would become harder to maintain.
- Decision: Store global operational configuration in a singleton `mess_settings` row, drive wifi default price/deadline from that row, and standardize backend API errors through the shared `{ message, data, error.code }` response envelope.
- Rationale: This keeps admin-editable settings in one canonical place, removes drift between settings and wifi bill creation defaults, and gives SvelteKit server actions one consistent error contract across modules.
- Impact: Future configurable mess-wide values should extend `mess_settings` first, and new backend handlers should reuse the shared response helper instead of returning ad hoc error JSON.
- Follow-up: If media upload/storage is added later, keep the current avatar string contract stable or add a backward-compatible migration path.

## Decision 16
- Date: 2026-03-16
- Context: STEP 7 needs production-oriented reliability without rewriting the current SvelteKit + Fiber architecture or adding heavy observability/security infrastructure.
- Decision: Standardize backend error payloads to `{ error, message, details? }`, attach `X-Request-ID` to every backend request/response, log API traffic as structured JSON, apply targeted Fiber rate limits on login/feed write routes, keep `/health` readiness-oriented and public, and restrict `/api/v1/system/status` to admin access.
- Rationale: This keeps debugging and permission checks centered in the existing app runtime, improves misuse resistance, and gives frontend/server actions one stable contract for unauthorized/forbidden/network handling without introducing a larger platform migration.
- Impact: Future backend handlers should validate request payloads before service calls, reuse the shared error helper, preserve request IDs in logs and client-visible failures, and treat backend route guards as the source of truth for access control.
- Follow-up: Revisit fuller observability, alerting, and reverse-proxy header policy once live VPS validation confirms the current hardening defaults do not conflict with Nginx.

## Decision 17
- Date: 2026-03-16
- Context: STEP 8 adds installability, browser-side push subscription, service-worker background sync, and offline queue replay, but the backend previously accepted authenticated requests only through bearer headers while the frontend session token is intentionally stored as an httpOnly cookie.
- Decision: Keep the existing split SvelteKit + Fiber architecture, extend the backend auth middleware to accept the existing `mh_access_token` cookie in addition to bearer tokens, and let browser/service-worker PWA features call the same-origin `/api/v1` routes directly while service-worker caches remain versioned and restricted to whitelisted static assets plus safe GET routes.
- Rationale: This preserves the current cookie security model, avoids exposing tokens to client JavaScript just to support Push/Background Sync, and keeps PWA/offline behavior incremental on top of the native SvelteKit service worker baseline instead of introducing a parallel proxy/auth stack.
- Impact: Future client-side PWA features should prefer same-origin `/api/v1` calls that rely on the existing auth cookie, and any new service-worker caching should stay on explicit allowlists rather than broad authenticated API caching.
- Follow-up: Validate the cookie-auth + service-worker flow on the live HTTPS origin and revisit the cache allowlist if later phases add more offline-safe routes or more sensitive browser-side data.

## Decision 18
- Date: 2026-03-16
- Context: STEP 9 needs a practical migration path from legacy spreadsheets without introducing unstable live integrations or allowing imports to write directly into production data without review.
- Decision: Use admin-only CSV imports with a two-step preview/commit flow, store preview metadata in `import_jobs`, preserve wallet `transaction_date` and optional `proof_url`, skip spreadsheet saldo as a source of truth, and recalculate wallet balance from imported transactions.
- Rationale: CSV works with Google Sheets exports and old spreadsheet structures, keeps migration deterministic, allows row-level validation plus duplicate warnings before writes happen, and leaves a durable audit/import trail without adding direct Google Sheets coupling.
- Impact: Future migration work should extend the existing CSV preview/commit pipeline first, and any wallet reporting or re-import logic must continue treating transactions as canonical while leaving spreadsheet balances informational only.
- Follow-up: Validate the import rules against real exported mess spreadsheets and refine category inference or duplicate handling only if actual data shows gaps.

## Decision 19
- Date: 2026-03-16
- Context: The app shell and login flow needed a production-ready UX refresh without rewriting the existing route structure or breaking the current auth/role model.
- Decision: Keep the existing SvelteKit routes, but centralize navigation and browser metadata in a shared frontend config, adopt warm interior-inspired theme tokens plus Iconify-based iconography, and extend auth so users can sign in with either `email` or `username` while preserving the existing response envelope and role guards.
- Rationale: This improves consistency across desktop/mobile/PWA navigation, avoids title/menu drift between routes, gives the frontend a maintainable tokenized design baseline, and adds a more practical login identifier without forcing a breaking auth rewrite.
- Impact: Future frontend routes should derive titles/navigation state from the shared metadata helper first, new UI work should reuse the warm theme tokens and Iconify wrapper components, and user creation/import flows must continue assigning stable unique usernames.
- Follow-up: Validate the new shell and username migration against live device installs plus real production user data before rollout.
