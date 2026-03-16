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
