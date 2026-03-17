# MessHub

MessHub adalah monorepo aplikasi operasional mess dengan runtime terpisah:

- `frontend/`: SvelteKit + Tailwind + PWA
- `backend/`: Go Fiber + PostgreSQL

Deploy utama ditujukan untuk VPS Linux dengan:

- GAS CLI
- PM2
- Nginx
- frontend dan backend berjalan terpisah

Docker bukan jalur utama deploy. Saat ini repo hanya menyimpan runtime frontend/backend dan service-scoped `.env.example`; provisioning Postgres lokal perlu disiapkan terpisah.

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
- `backend/db/migrations/`: urutan migrasi database yang harus diterapkan berurutan

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

### 1. Siapkan Postgres Lokal

Pastikan ada instance PostgreSQL yang bisa diakses backend, misalnya lewat service lokal/VPS dev, lalu sesuaikan `DATABASE_URL` di `backend/.env`.

### 2. Setup Backend

```bash
cd backend
cp .env.example .env
go mod tidy
for file in db/migrations/*.sql; do
  psql "postgres://messhub:messhub@127.0.0.1:5432/messhub?sslmode=disable" -f "$file"
done
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

## CI/CD GitHub Actions

Deploy otomatis production sekarang memakai workflow GitHub Actions di [`.github/workflows/deploy.yml`](/Users/ryanprayoga/Kerjaan/Pribadi/messhub/.github/workflows/deploy.yml).

Trigger:

- push ke branch `main`

Secrets yang wajib di-set di GitHub repository:

- `VPS_HOST`: host atau IP VPS Ubuntu
- `VPS_USER`: user SSH VPS, mis. `ubuntu`
- `VPS_SSH_KEY`: private key SSH untuk login ke VPS

Alur workflow:

1. checkout repository
2. load private key ke `ssh-agent`
3. tambah host VPS ke `known_hosts`
4. SSH ke VPS dan masuk ke `/home/ubuntu/projects/messhub`
5. `git checkout main` lalu `git pull --ff-only origin main`
6. build backend dengan GAS CLI dan restart PM2 `messhub-backend`
7. build frontend dengan GAS CLI dan restart PM2 `messhub-frontend`
8. jalankan health check `curl http://127.0.0.1:4100/health`

Command yang dijalankan di VPS:

```bash
cd /home/ubuntu/projects/messhub
git checkout main
git pull --ff-only origin main

cd backend
gas build --no-ui --type go --pm2-name messhub-backend --port 4100 --yes

cd ../frontend
gas build --no-ui --type svelte --pm2-name messhub-frontend --port 4101 --yes

curl --fail --silent --show-error http://127.0.0.1:4100/health
```

Catatan penting:

- workflow deploy saat ini **belum** menjalankan migrasi database otomatis
- setiap perubahan schema tetap perlu memastikan file di `backend/db/migrations/` sudah diterapkan berurutan di environment target sebelum atau saat deploy

Contoh setup secrets di GitHub:

1. Buka repository `Settings` -> `Secrets and variables` -> `Actions`
2. Tambah secret `VPS_HOST`
3. Tambah secret `VPS_USER`
4. Tambah secret `VPS_SSH_KEY`

Contoh generate SSH key khusus deploy:

```bash
ssh-keygen -t ed25519 -C "github-actions@messhub" -f ~/.ssh/messhub_github_actions
cat ~/.ssh/messhub_github_actions.pub
```

Tambahkan public key ke VPS pada `~/.ssh/authorized_keys` untuk user deploy, lalu simpan isi private key `~/.ssh/messhub_github_actions` sebagai secret `VPS_SSH_KEY`.

Contoh log deploy yang diharapkan:

```text
Run ssh "${VPS_USER}@${VPS_HOST}"
+ cd /home/ubuntu/projects/messhub
+ git checkout main
+ git pull --ff-only origin main
Already up to date.
+ cd backend
+ gas build --no-ui --type go --pm2-name messhub-backend --port 4100 --yes
[PM2] Applying action restartProcessId on app [messhub-backend]
+ cd ../frontend
+ gas build --no-ui --type svelte --pm2-name messhub-frontend --port 4101 --yes
[PM2] Applying action restartProcessId on app [messhub-frontend]
+ curl --fail --silent --show-error http://127.0.0.1:4100/health
{"message":"service ready","data":{"status":"ok","database_reachable":true}}
```

Cara test CI/CD:

1. pastikan secret GitHub sudah terpasang
2. pastikan VPS bisa `git pull origin main` tanpa prompt interaktif
3. commit perubahan kecil ke branch `main`
4. push: `git push origin main`
5. buka tab `Actions` dan verifikasi job `Deploy to VPS` sukses
6. cek aplikasi di `https://messhub.ryannn.net` dan endpoint health di VPS

Rollback jika deploy gagal:

1. SSH ke VPS
2. buka folder `/home/ubuntu/projects/messhub`
3. cek commit sebelumnya yang stabil: `git log --oneline -n 5`
4. checkout commit stabil atau branch/tag yang diinginkan
5. jalankan ulang command deploy manual yang sama:

```bash
cd /home/ubuntu/projects/messhub/backend
gas build --no-ui --type go --pm2-name messhub-backend --port 4100 --yes

cd /home/ubuntu/projects/messhub/frontend
gas build --no-ui --type svelte --pm2-name messhub-frontend --port 4101 --yes

curl --fail http://127.0.0.1:4100/health
```

Rollback tetap memakai arsitektur deploy yang sama: git checkout ke revision stabil lalu rebuild kedua service dengan GAS CLI.

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

Docker tidak dipakai untuk menjalankan frontend atau backend di VPS production, dan file compose lokal tidak termasuk dalam state repo saat ini.

## Next Steps

1. Validasi rollout STEP 7 di VPS production, terutama `CORS_ORIGIN`, `JWT_SECRET`, health/readiness response, dan request logging di PM2/Nginx.
2. Kembali ke modul yang masih placeholder: shared expenses dan proposals.
3. Tambahkan create-member flow di frontend admin agar manajemen anggota tidak berhenti di role/activation update.
4. Putuskan kebutuhan upload/storage nyata untuk bukti transfer dan avatar bila string reference sudah tidak cukup.
5. Lanjutkan STEP 8 bila PWA upgrade, push notifications, dan offline support sudah diprioritaskan.
