# TASKS.md

## Todo
- ID: T-003
  Title: Implement member management
  Priority: P1
  Owner/Agent: Unassigned
  Dependencies: T-002
  Notes: Support dynamic active members with `joined_at` and optional `left_at`.

- ID: T-004
  Title: Implement wallet transactions module
  Priority: P1
  Owner/Agent: Unassigned
  Dependencies: T-002
  Notes: Must keep auditable separation from wifi and non-cash shared expenses.

- ID: T-005
  Title: Implement monthly wifi billing and proof verification
  Priority: P1
  Owner/Agent: Unassigned
  Dependencies: T-003
  Notes: Default nominal Rp20.000, deadline before day 10, statuses: unpaid, pending_verification, verified, rejected.

- ID: T-006
  Title: Implement dashboard summary
  Priority: P2
  Owner/Agent: Unassigned
  Dependencies: T-004, T-005
  Notes: Mobile-first summary of wallet, wifi, shared expenses, contributions, feed, proposals.

- ID: T-007
  Title: Implement shared expenses module
  Priority: P2
  Owner/Agent: Unassigned
  Dependencies: T-003
  Notes: Track payer, fronting, reimbursement status, and proof without affecting wallet balance.

- ID: T-008
  Title: Implement contributions and simple leaderboard base
  Priority: P2
  Owner/Agent: Unassigned
  Dependencies: T-003
  Notes: Keep scoring simple in v1.

- ID: T-009
  Title: Implement temporary feed/info module
  Priority: P2
  Owner/Agent: Unassigned
  Dependencies: T-003
  Notes: Support expiring posts and basic interactions.

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

- ID: CTX-001
  Title: Normalize shared project context
  Priority: P1
  Owner/Agent: Codex
  Dependencies: `messhub-masterplan.md`
  Notes: Created minimal cross-agent operating context from the existing masterplan, including safe runtime ignore rules.
