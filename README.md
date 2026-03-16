# MessHub

MessHub adalah monorepo aplikasi operasional mess dengan runtime terpisah:

- `frontend/`: SvelteKit + Tailwind + PWA
- `backend/`: Go Fiber + PostgreSQL

Deploy utama ditujukan untuk VPS Linux dengan:

- GAS CLI
- PM2
- Nginx
- frontend dan backend berjalan terpisah

Docker bukan jalur utama deploy. `docker-compose.yml` dipertahankan hanya untuk Postgres lokal.

## Production Workflow

Production deploy memakai:

- GAS CLI untuk build/run metadata
- PM2 untuk process manager
- Nginx untuk reverse proxy

Frontend dan backend tidak dijalankan via container pada workflow production ini.

## Runtime Default

- Frontend port: `4101`
- Backend port: `4100`
- PM2 frontend: `messhub-frontend`
- PM2 backend: `messhub-backend`

Port ini dipilih karena `3000`, `4000`, `4001`, `5000`, dan `5001` sudah terpakai.

## Struktur Project

```text
.
├── frontend/
│   ├── .env.example
│   ├── ecosystem.config.cjs
│   ├── package.json
│   ├── svelte.config.js
│   ├── vite.config.ts
│   └── src/
├── backend/
│   ├── .env.example
│   ├── cmd/
│   ├── db/
│   ├── go.mod
│   └── internal/
├── docs/
├── docker-compose.yml
├── .env.example
└── README.md
```

## File Penting

- `frontend/ecosystem.config.cjs`: default PM2 config untuk GAS/PM2
- `frontend/.env.example`: env frontend untuk VPS/local run
- `frontend/package.json`: script `dev`, `build`, `preview`, `start`
- `frontend/svelte.config.js`: pakai `@sveltejs/adapter-node`
- `frontend/vite.config.ts`: dev proxy ke backend `4100`
- `backend/.env.example`: env backend untuk VPS/local run
- `backend/internal/config/config.go`: backend sekarang membaca `PORT` lebih dulu
- `backend/db/migrations/001_init.sql`: schema awal database
- `docker-compose.yml`: Postgres lokal saja

## Frontend

Frontend sekarang diarahkan ke runtime Node production, bukan preview server:

- adapter: `@sveltejs/adapter-node`
- script start: `npm run start`
- env utama: `PORT`, `HOST`, `ORIGIN`, `PUBLIC_API_BASE_URL`, `PRIVATE_API_BASE_URL`
- default API URL: `/api/v1`

Ini membuat frontend cocok untuk:

```bash
cd frontend
gas build --no-ui --type svelte --pm2-name messhub-frontend --port 4101 --yes
```

## Backend

Backend sekarang membaca port dari `PORT`, lalu fallback ke `BACKEND_PORT`.

Default env backend:

- `PORT=4100`
- `BACKEND_HOST=0.0.0.0`
- `DATABASE_URL=postgres://...`
- `JWT_SECRET=...`
- `CORS_ORIGIN=http://127.0.0.1:4101,http://localhost:4101`

Ini membuat backend cocok untuk:

```bash
cd backend
gas build --no-ui --type go --pm2-name messhub-backend --port 4100 --yes
```

## Setup Local

### 1. Jalankan Postgres Lokal

```bash
cp .env.example .env
docker compose up -d
```

### 2. Setup Backend

```bash
cd backend
cp .env.example .env
go mod tidy
psql "postgres://messhub:messhub@127.0.0.1:5432/messhub?sslmode=disable" -f db/migrations/001_init.sql
psql "postgres://messhub:messhub@127.0.0.1:5432/messhub?sslmode=disable" -f db/migrations/002_auth_foundation.sql
go run ./cmd/seed-admin
go run ./cmd/api
```

Backend akan listen di `http://127.0.0.1:4100`.

### 3. Setup Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Frontend akan listen di `http://127.0.0.1:4101`.

Dev proxy akan mengarahkan `/api/*` dari frontend ke backend `4100`.
Server-side auth/data loads akan memakai `PRIVATE_API_BASE_URL`, yang default-nya diarahkan ke `http://127.0.0.1:4100/api/v1`.

## Workflow VPS Dengan GAS CLI

### Frontend

```bash
cd frontend
cp .env.example .env
npm install
gas build --no-ui --type svelte --pm2-name messhub-frontend --port 4101 --yes
```

### Backend

```bash
cd backend
cp .env.example .env
go mod tidy
gas build --no-ui --type go --pm2-name messhub-backend --port 4100 --yes
```

### Deploy Nginx Split Frontend/Backend

```bash
gas deploy --no-ui \
  --frontend messhub-frontend \
  --backend messhub-backend \
  --domain messhub.example.com \
  --mode frontend-backend-split \
  --ssl certbot-nginx \
  --yes
```

Ekspektasi route deploy:

- `/` -> frontend `messhub-frontend`
- `/api/` -> backend `messhub-backend`

Karena frontend default ke `/api/v1`, tidak perlu hardcode domain backend di browser.

## Seed Admin

Default seed:

- Email: `admin@messhub.local`
- Password: lihat `SEED_ADMIN_PASSWORD` di `backend/.env`

Command:

```bash
cd backend
go run ./cmd/seed-admin
```

## Docker

`docker-compose.yml` dipertahankan hanya untuk:

- Postgres lokal
- bootstrap development cepat

Docker tidak dipakai untuk menjalankan frontend atau backend di VPS production.

## Next Steps

1. Tambahkan create/edit members UI untuk admin bila STEP 1 perlu diperdalam di frontend.
2. Implement `wallet_transactions`.
3. Implement `wifi_bills` dan `wifi_bill_members`.
4. Tambahkan flow upload bukti pembayaran.
5. Tambahkan hardening env production per domain VPS, termasuk final `PRIVATE_API_BASE_URL`.
