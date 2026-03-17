# TASKS.md

## Todo
- None.

## Doing
- None.

## Blocked
- None observed.

## Done
- ID: OPS-009
  Title: Replace native confirms, fix destructive-action correctness, and polish shell/page spacing
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-008, UX-001
  Notes: Replaced the last native browser confirm flow with a reusable accessible confirmation dialog plus shared confirmable submit controller, fixed the Cancel path that could still leak into enhanced form submission, added loading/double-submit protection for destructive member and wifi lifecycle actions, refined modal/page spacing across the key routes, made desktop/mobile navigation feel immediately responsive while keeping active state tied to the actual route, and followed up with optimistic shell menu switching, route-code preloading, and a compact two-column mobile overflow panel under the header that only shows non-navbar destinations as icon-first tiles.

- ID: OPS-008
  Title: Ship sprint hardening, shared expenses, proposals, and safe member lifecycle
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-007, T-007, T-010
  Notes: Added JWT invalidation through `auth_version`, fixed the wrong-current-password logout path, restricted settings reads to manager roles, split feed into active/history with expired interaction guards, added wifi bill status lifecycle updates plus no-overwrite pending proof rules and reject notifications, introduced a bootstrap-aware Go migration runner and migration-backed deploy step, completed shared expenses and proposals end-to-end on the existing schema, and finished member lifecycle with archive/reactivate/delete-safe behavior plus richer admin UI.

- ID: OPS-007
  Title: Complete admin member management, fix protected action session resolution, and refine shell states
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-006, T-003, T-012
  Notes: Fixed the false missing-session error on protected SvelteKit actions by resolving auth server-side, tightened shell/sidebar bounds plus route-active behavior, added richer empty-state/sheet UI, completed admin member create/edit/role/status/reset-password flows with audit logging and safety guards, fixed the `create_user` audit-log scan bug that was turning successful member inserts into `500` responses, and hardened member-create feedback so duplicate username/email conflicts surface cleanly in centered sheets plus global toasts instead of generic hidden errors.

- ID: OPS-006
  Title: Implement navigation polish, warm-theme UI refresh, and username login
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-005, T-002, UX-001
  Notes: Refreshed the authenticated shell and key pages for desktop/mobile/PWA, standardized titles/meta, adopted Iconify, added animation/loading polish, and extended auth to accept either email or username while preserving existing roles and response contracts.

- ID: OPS-005
  Title: Implement STEP 9 data import, UI cleanup, and migration tools
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-004, T-003, T-004
  Notes: Added admin-only CSV preview/commit flows for members and wallet, import job tracking plus audit logs, downloadable templates, spreadsheet migration docs, wallet transaction date/proof preservation, and production-facing UI copy cleanup across key pages.

- ID: OPS-004
  Title: Implement STEP 8 PWA upgrade and mobile experience
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-003
  Notes: Added install-ready manifest/icons, install prompt UI, versioned service worker caching with `/offline` fallback, Web Push subscriptions plus backend fan-out for wifi/feed events, lightweight activity/comment outbox background sync, cookie-auth browser/service-worker API access, mobile bottom navigation, and pull-to-refresh on dashboard/feed.

- ID: OPS-003
  Title: Implement STEP 7 production hardening and reliability
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-012
  Notes: Standardized backend API error payloads, added request validation, stricter role middleware usage, request IDs, structured request logging, rate limits for login/feed writes, JSON panic recovery, health/readiness status improvements, security headers, frontend API error helpers, route-level unauthorized redirects, forbidden/error/empty/loading states, and updated env examples.

- ID: OPS-002
  Title: Automate VPS deploys from GitHub Actions
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-001
  Notes: Added `.github/workflows/deploy.yml` for push-to-main SSH deploys to `/home/ubuntu/projects/messhub`, reusing the existing GAS CLI commands for `messhub-backend` on `4100` and `messhub-frontend` on `4101`, with a backend `/health` check and documented secrets/rollback flow.

- ID: DX-001
  Title: Stabilize frontend install/runtime warnings on Node 24
  Priority: P1
  Owner/Agent: Codex
  Dependencies: OPS-001
  Notes: Replaced the Workbox-based PWA plugin with a static manifest plus native SvelteKit service worker, added sync hooks around install/dev/build, removed deprecated `npm install` warnings, and eliminated the fresh-checkout `.svelte-kit/tsconfig.json` warning.

- ID: UX-001
  Title: Polish frontend baseline UI for login and app shell
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-002
  Notes: Fixed root Tailwind stylesheet import, added shared global UI classes, refined login states, and refreshed the mobile-first AppShell/navigation baseline.

- ID: OPS-001
  Title: Adapt deploy/runtime workflow to GAS CLI, PM2, and Nginx
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-001, T-002
  Notes: Frontend now uses adapter-node + ecosystem config on port 4101, backend prefers PORT on 4100, service env files were split, and Docker was reduced to local Postgres only.

- ID: T-001
  Title: Initialize frontend and backend repo structure
  Priority: P1
  Owner/Agent: Codex
  Dependencies: None
  Notes: Monorepo root, frontend, backend, env example, Docker Compose, and README are in place.

- ID: T-002
  Title: Implement auth and role guard
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-001
  Notes: Basic login endpoint, JWT middleware, role middleware, frontend cookie auth guard, and login/logout flow are scaffolded.

- ID: T-003
  Title: Implement member management
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-002
  Notes: Users schema is now hardened with `joined_at`, `phone`, `avatar_url`, `auth_version`, and `archived_at`; admin seed remains non-duplicating; backend users/profile APIs are live; and the frontend members list now supports admin role/status/lifecycle controls including archive, reactivate, and delete-safe guards.

- ID: T-004
  Title: Implement wallet transactions module
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-002
  Notes: Added wallet summary and paginated transaction APIs, create-transaction role guard for admin/treasurer, a wallet migration follow-up, `/wallet` and `/wallet/new` frontend flows, and dashboard wallet summary integration while preserving the separation from wifi and non-cash shared expenses.

- ID: T-005
  Title: Implement monthly wifi billing and proof verification
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-003
  Notes: Added `backend/db/migrations/004_wifi_audit_step3.sql`, wifi billing APIs for create/list/detail/active/my/submit/verify/reject plus bill status lifecycle updates, auto-generated monthly member obligations from active non-archived users, frontend `/wifi` page with proof submission and review UI, dashboard wifi summary, and reusable audit logging for wifi, wallet transaction creation, and user role/activation updates.

- ID: T-006
  Title: Implement dashboard summary
  Priority: P2
  Owner/Agent: Codex
  Dependencies: T-004, T-005
  Notes: Mobile-first dashboard summary cards are now live for members, wallet, wifi, and contribution leaderboard data.

- ID: T-008
  Title: Implement contributions and simple leaderboard base
  Priority: P2
  Owner/Agent: Codex
  Dependencies: T-003
  Notes: Added unified `activities` runtime for contribution posts, `GET /api/v1/contributions/leaderboard` with `period=month|all`, dashboard and `/contributions` leaderboard UI, and simple point scoring from contribution activities.

- ID: T-009
  Title: Implement temporary feed/info module
  Priority: P2
  Owner/Agent: Codex
  Dependencies: T-003
  Notes: Added `/api/v1/activities` feed runtime with create/comment/reaction support, food claim and rice response flows, and a mobile-first `/feed` UI for smart mess interactions.

- ID: T-011
  Title: Implement in-app notification system
  Priority: P2
  Owner/Agent: Codex
  Dependencies: T-005, T-009
  Notes: Added notification list/read APIs, activity/comment/wifi notification triggers, header badge, and `/notifications` UI with unread management.

- ID: T-012
  Title: Implement system settings, profile, and admin panel
  Priority: P1
  Owner/Agent: Codex
  Dependencies: T-003, T-005, T-011
  Notes: Added `backend/db/migrations/006_settings_profile_step6.sql`, live profile/password/settings/system-status APIs, standardized API errors, frontend `/profile` and admin `/settings`, wifi defaults sourced from settings, and admin member role/activation controls on `/members`.

- ID: T-007
  Title: Implement shared expenses module
  Priority: P2
  Owner/Agent: Codex
  Dependencies: T-003
  Notes: Implemented list/create/update runtime on the existing `shared_expenses` table, added `updated_at` migration hardening, summary cards and admin forms on `/shared-expenses`, dashboard integration, and audit logs for create/update/reimbursement status changes without affecting wallet balance.

- ID: T-010
  Title: Implement proposals and simple voting
  Priority: P3
  Owner/Agent: Codex
  Dependencies: T-003
  Notes: Implemented create/list/detail/vote/close/finalize runtime on the existing proposals tables, added real `/proposals` UI with voting and admin finalization flows, dashboard proposal summary, and audit logs for proposal lifecycle actions.

- ID: CTX-001
  Title: Normalize shared project context
  Priority: P1
  Owner/Agent: Codex
  Dependencies: `messhub-masterplan.md`
  Notes: Created minimal cross-agent operating context from the existing masterplan, including safe runtime ignore rules.
