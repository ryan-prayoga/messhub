# PROJECT_STATUS.md

## Current Objective
- Adapt the existing MessHub monorepo to the VPS runtime workflow based on GAS CLI, PM2, and Nginx.

## Current Phase
- Phase 1 — Foundation bootstrap

## Summary Status
- Observed: monorepo baseline now exists with `frontend/`, `backend/`, root env/compose/readme, and initial database migration.
- Observed: frontend includes mobile-first AppShell, placeholder routes, API client, and cookie-based auth guard.
- Observed: backend includes Fiber app bootstrap, PostgreSQL connection, JWT auth, middleware auth/role, migration, and admin seed command.
- Observed: frontend is now prepared for PM2 runtime with `adapter-node`, `start` script, local `.env.example`, and `ecosystem.config.cjs`.
- Observed: backend now prefers `PORT` and has its own `.env.example` for GAS/PM2 usage.
- Observed: Docker has been reduced to local Postgres only.

## Done
- Masterplan exists in `messhub-masterplan.md`.
- Shared context baseline created: `AGENTS.md`, `PROJECT_STATUS.md`, `TASKS.md`, `docs/decisions.md`, `docs/handoffs/HANDOFF_TEMPLATE.md`.
- Agent runtime ignore baseline created in `.gitignore`.
- Monorepo root created with `.env.example`, `docker-compose.yml`, and `README.md`.
- Frontend scaffold created with SvelteKit, Tailwind config, PWA config, auth flow baseline, and placeholder routes.
- Backend scaffold created with Go Fiber, env config, PostgreSQL bootstrap, `/api/v1`, auth endpoints, middleware, migration SQL, and seed admin.
- Runtime/deploy adaptation completed for separated frontend/backend services on ports `4101` and `4100`.

## In Progress
- Domain modules are still placeholder-only on the frontend.
- VPS build/deploy has not been executed from this workspace yet.

## Blockers / Risks
- Dependencies have not been installed in this environment.
- Migration is manual; schema application is not yet automated.
- Production deploy has not been validated with live GAS build or Nginx apply yet.
- Risk: scope can expand too early if phase order is not enforced.

## Recently Touched Areas
- `.gitignore`
- `.env.example`
- `frontend/.env.example`
- `frontend/ecosystem.config.cjs`
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

## Assumptions / Unknowns
- Assumption: `messhub-masterplan.md` is the current approved product source.
- Assumption: stack choices in the masterplan are accepted for the bootstrap.
- Assumption: deploy target remains VPS Linux with PM2 and Nginx split frontend/backend.
- Unknown: final public domain and SSL mode for production deploy.
- Unknown: preferred migration tool and CI workflow.

## Next Recommended Steps
- Install dependencies and run frontend/backend on `4101` and `4100`.
- Apply migration and seed the initial admin user.
- Run `gas build` from `frontend/` and `backend/`, then preview `gas deploy` in split mode.
- Continue with member management and wallet/wifi modules.
- Record any stack or architecture changes in `docs/decisions.md`.
