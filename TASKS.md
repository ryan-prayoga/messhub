# TASKS.md

## Todo
- ID: T-007
  Title: Implement shared expenses module
  Priority: P2
  Owner/Agent: Unassigned
  Dependencies: T-003
  Notes: Track payer, fronting, reimbursement status, and proof without affecting wallet balance.

- ID: T-010
  Title: Implement proposals and simple voting
  Priority: P3
  Owner/Agent: Unassigned
  Dependencies: T-003
  Notes: Candidate for v1.1 if scope tightens.

## Doing
- None.

## Blocked
- None observed.

## Done
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
  Notes: Users schema is now hardened with `joined_at`, `phone`, and `avatar_url`; admin seed remains non-duplicating; backend users/profile APIs are live; and the frontend members list now supports admin role and activation controls.

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
  Notes: Added `backend/db/migrations/004_wifi_audit_step3.sql`, wifi billing APIs for create/list/detail/active/my/submit/verify/reject, auto-generated monthly member obligations from active users, frontend `/wifi` page with proof submission and review UI, dashboard wifi summary, and reusable audit logging for wifi, wallet transaction creation, and user role/activation updates.

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

- ID: CTX-001
  Title: Normalize shared project context
  Priority: P1
  Owner/Agent: Codex
  Dependencies: `messhub-masterplan.md`
  Notes: Created minimal cross-agent operating context from the existing masterplan, including safe runtime ignore rules.
