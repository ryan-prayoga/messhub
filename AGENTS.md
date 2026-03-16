# AGENTS.md

## Project Intent
- Build `MessHub`, a mobile-first PWA for internal mess operations.
- Primary v1 scope: auth, member management, wallet transactions, wifi billing and proof verification, shared expenses, contributions, simple feed.
- Source of truth for product intent: `messhub-masterplan.md`.

## Read First
1. `messhub-masterplan.md`
2. `PROJECT_STATUS.md`
3. `TASKS.md`
4. `docs/decisions.md`
5. `docs/handoffs/`

## Working Rules
- Prefer observed repo state over assumptions.
- Keep changes scoped to the active task.
- Do not expand scope beyond current phase unless `TASKS.md` or `PROJECT_STATUS.md` is updated first.
- Reuse existing files and structure; avoid redundant docs.
- Mark uncertain statements as `Assumption`, `Unknown`, or `Placeholder`.

## Context Check Before Work
- Confirm current objective in `PROJECT_STATUS.md`.
- Confirm active task status in `TASKS.md`.
- Read relevant decisions in `docs/decisions.md`.
- Check latest handoff if continuing prior work.
- If repo state conflicts with docs, update docs to match observed state.

## Definition of Done
- Code or doc change matches the active task.
- Relevant validation is run or explicitly marked as not run.
- `PROJECT_STATUS.md` reflects current reality when scope/status changed.
- `TASKS.md` status is updated when task state changed.
- `docs/decisions.md` is updated only for durable decisions.

## Finishing Ritual
- Update touched task entries.
- Update status summary if progress changed.
- Add a decision entry only if the choice affects future work.
- Create a handoff note from the template when handing work to another agent or ending incomplete work.

## Required Output Format
- `Summary`: what changed.
- `Validation`: what was checked.
- `Docs Updated`: which context files changed.
- `Open Items`: blockers, assumptions, or next step.

## Update Rules
- `PROJECT_STATUS.md`: update on any meaningful progress, blocker, or scope change.
- `TASKS.md`: update whenever a task moves between sections.
- `docs/decisions.md`: append-only for durable decisions; do not rewrite history.
- `docs/handoffs/`: add a new handoff file only when a real handoff is needed; keep the template unchanged.
