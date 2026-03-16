# MessHub

MessHub adalah monorepo awal untuk aplikasi PWA internal mess dengan frontend `SvelteKit + Tailwind + PWA` dan backend `Go Fiber + PostgreSQL`.

## Stack

- Frontend: SvelteKit, TailwindCSS, Vite PWA
- Backend: Go Fiber, JWT auth, PostgreSQL
- Infra: Docker Compose

## Struktur Folder

```text
.
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/
│   │   │   ├── auth/
│   │   │   ├── components/
│   │   │   ├── config/
│   │   │   └── stores/
│   │   ├── routes/
│   │   │   ├── login/
│   │   │   ├── members/
│   │   │   ├── wallet/
│   │   │   ├── wifi/
│   │   │   ├── shared-expenses/
│   │   │   ├── contributions/
│   │   │   ├── feed/
│   │   │   ├── proposals/
│   │   │   ├── profile/
│   │   │   └── settings/
│   │   ├── app.css
│   │   ├── app.d.ts
│   │   ├── app.html
│   │   └── hooks.server.ts
│   ├── static/
│   ├── package.json
│   ├── svelte.config.js
│   ├── tailwind.config.ts
│   └── vite.config.ts
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   └── seed-admin/
│   ├── db/
│   │   ├── migrations/
│   │   └── seeds/
│   ├── internal/
│   │   ├── app/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── routes/
│   │   ├── services/
│   │   └── types/
│   └── go.mod
├── docs/
├── docker-compose.yml
├── .env.example
└── README.md
```

## Yang Sudah Disiapkan

- Mobile-first AppShell dan placeholder routes:
  - `/login`
  - `/`
  - `/members`
  - `/wallet`
  - `/wifi`
  - `/shared-expenses`
  - `/contributions`
  - `/feed`
  - `/proposals`
  - `/profile`
  - `/settings`
- API client frontend
- Basic auth guard berbasis cookie
- Backend `/api/v1`
- Health check: `GET /api/v1/health`
- Auth endpoints:
  - `POST /api/v1/auth/login`
  - `GET /api/v1/auth/me`
- Middleware auth JWT
- Middleware role untuk `admin` dan `treasurer`
- Migration SQL awal untuk semua tabel inti
- Seed admin command

## Setup

1. Salin env:

```bash
cp .env.example .env
```

2. Jalankan PostgreSQL:

```bash
docker compose up -d postgres
```

3. Export env untuk shell saat ini:

```bash
set -a
source .env
set +a
```

4. Jalankan migration:

```bash
psql "$DATABASE_URL" -f backend/db/migrations/001_init.sql
```

5. Seed admin:

```bash
cd backend
go run ./cmd/seed-admin
```

6. Jalankan backend:

```bash
cd backend
go run ./cmd/api
```

7. Jalankan frontend:

```bash
cd frontend
npm install
npm run dev
```

8. Buka aplikasi:

- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8080/api/v1/health`

## Setup Dengan Docker Compose

```bash
cp .env.example .env
docker compose up
```

Catatan:
- Compose ini cocok untuk bootstrap dev.
- Migration masih dijalankan manual agar flow schema tetap jelas.
- Frontend membaca root `.env` lewat `envDir`.
- Backend akan membaca `.env` dari root repo atau dari folder backend jika file tersedia.

## Kredensial Awal

- Email: `admin@messhub.local`
- Password: ambil dari `SEED_ADMIN_PASSWORD` di `.env`

## File Penting

- `frontend/src/routes/login/+page.server.ts`: login action dan set auth cookies
- `frontend/src/routes/+layout.server.ts`: auth guard dasar
- `frontend/src/lib/api/client.ts`: API client untuk backend
- `frontend/src/lib/components/AppShell.svelte`: shell mobile-first
- `backend/internal/services/auth_service.go`: validasi password dan issue JWT
- `backend/internal/middleware/auth.go`: middleware auth dan role
- `backend/db/migrations/001_init.sql`: schema awal database
- `backend/cmd/seed-admin/main.go`: seed admin idempotent

## Next Steps

1. Tambahkan loader migration otomatis atau Makefile task supaya bootstrap lebih singkat.
2. Implement CRUD member management.
3. Implement wallet transactions dan wifi billing flow.
4. Tambahkan password reset dan invite/onboarding flow.
5. Ganti adapter frontend dan setup Dockerfile jika mulai harden untuk production deploy.
