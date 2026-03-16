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
