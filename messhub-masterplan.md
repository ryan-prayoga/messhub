# MessHub — Project Masterplan

## 1. Ringkasan Produk

**MessHub** adalah aplikasi PWA untuk manajemen kehidupan mess yang fokus pada:
- pencatatan **Kantong Duafa / kas**
- manajemen **tagihan wifi bulanan**
- pencatatan **pengeluaran talangan / non-kas**
- sistem **usulan dan persetujuan**
- pencatatan **kontribusi penghuni**
- sistem **info sementara** untuk makanan, nasi, dan pengumuman ringan

Aplikasi ini dipakai internal untuk **Mess Traspac Menyala 🔥** dengan jumlah anggota yang **dinamis**.

Target utama versi awal:
- bisa di-install di HP
- ringan dipakai harian
- transparan
- tidak ribet
- cocok untuk admin, bendahara, dan anggota biasa

---

## 2. Tujuan Utama

### Tujuan bisnis / operasional
- menggantikan spreadsheet manual menjadi sistem terstruktur
- membuat kas dan wifi lebih transparan
- memudahkan verifikasi bukti transfer
- mencatat kontribusi penghuni secara adil
- memudahkan koordinasi info sehari-hari di mess

### Tujuan teknis
- PWA installable
- mobile-first
- backend stabil di VPS
- mudah dikembangkan bertahap
- role-based access
- data auditable

---

## 3. Batasan Scope Awal

Agar realistis, **v1 jangan terlalu melebar**.

### Yang masuk v1
- login email/password
- role admin / bendahara / anggota
- manajemen anggota
- transaksi kantong duafa
- wifi bulanan + upload bukti + verifikasi
- pengeluaran non-kas / talangan
- usulan + voting sederhana
- kontribusi penghuni
- info sementara / post dengan masa aktif
- notifikasi dalam aplikasi
- PWA installable

### Yang ditunda dulu
- multi-mess / SaaS
- payment gateway
- integrasi bank otomatis
- push notification native penuh
- chat real-time kompleks
- laporan akuntansi kompleks
- AI summary / OCR struk
- integrasi WhatsApp

---

## 4. Stakeholder dan Role

### Admin
Biasanya ketua mess.

Hak akses:
- kelola anggota
- tunjuk bendahara
- lihat semua data
- verifikasi pembayaran wifi
- buat dan kelola tagihan
- catat/edit transaksi kas
- catat/edit pengeluaran non-kas
- approve usulan
- moderasi post/info

### Bendahara
Orang yang ditunjuk.

Hak akses:
- hampir sama dengan admin untuk urusan keuangan
- verifikasi pembayaran wifi
- input transaksi kantong duafa
- kelola data tagihan
- lihat laporan

### Anggota
Penghuni biasa.

Hak akses:
- lihat dashboard
- lihat status wifi
- upload bukti pembayaran wifi
- buat usulan
- ikut voting
- buat post kontribusi
- buat info makanan / nasi / info ringan
- beri reaksi / komentar

---

## 5. Aturan Operasional Produk

### 5.1 Kantong Duafa
- hanya untuk pencatatan
- uang asli tetap disimpan di:
  - **Bank Jago**
  - **Ryan Prayoga**
  - **104987106615**
- aplikasi tidak memproses pembayaran, hanya mencatat bukti dan transaksi
- pemasukan “sumbangan sukarela” dicatat sebagai:
  - **Hibah Orang Baik**
- kantong duafa **tidak wajib**
- saldo aplikasi harus merepresentasikan catatan pemasukan dan pengeluaran terkait kantong duafa

### 5.2 Wifi
- wajib bulanan
- nominal tetap awal: **Rp20.000 per orang**
- deadline: **sebelum tanggal 10 setiap bulan**
- anggota upload bukti transfer
- admin/bendahara memverifikasi
- status pembayaran minimal:
  - belum bayar
  - menunggu verifikasi
  - terverifikasi
  - ditolak

### 5.3 Pengeluaran Mess Non-Kas
Contoh:
- keamanan
- kebersihan

Aturan:
- tidak mengurangi saldo kantong duafa
- tetap harus dicatat untuk transparansi
- perlu diketahui:
  - siapa yang bayar
  - siapa yang nalangi
  - nominal
  - status penggantian

Status contoh:
- ditanggung pribadi
- ditalangi
- sudah diganti sebagian
- sudah lunas

### 5.4 Kontribusi
Kontribusi bukan jadwal piket, tetapi aksi sukarela yang benar-benar dilakukan.

Contoh:
- buang sampah
- bersih dapur
- beres-beres
- isi galon
- beberes area umum

Tujuan:
- tercatat
- bisa dilihat leaderboard
- memberi apresiasi sosial

### 5.5 Info Sementara
Dipakai untuk:
- makanan bebas ambil
- pengumuman ringan
- tanya kepemilikan makanan
- rencana masak nasi
- info sesaat lain

Karakter:
- punya masa aktif
- default bisa 24 jam
- bisa diatur manual
- setelah expire:
  - hilang dari feed aktif
  - masuk history

---

## 6. Struktur Menu Aplikasi

### 6.1 Dashboard
Isi ringkasan cepat:
- saldo kantong duafa
- status wifi bulan berjalan
- transaksi terbaru
- pengeluaran non-kas terbaru
- kontribusi terbaru
- info aktif
- usulan yang masih terbuka

### 6.2 Kantong Duafa
Submenu:
- daftar transaksi
- tambah transaksi
- detail transaksi
- saldo berjalan
- filter per bulan
- statistik sederhana

Jenis transaksi:
- pemasukan
- pengeluaran

Kategori pemasukan:
- Hibah Orang Baik
- pemasukan lain

Kategori pengeluaran:
- gas
- plastik sampah
- keperluan bersama lain yang memang dibebankan ke kantong duafa

### 6.3 Wifi
Submenu:
- tagihan bulan ini
- daftar anggota dan status pembayaran
- upload bukti
- verifikasi admin/bendahara
- riwayat bulan sebelumnya

Informasi:
- nominal per orang
- deadline
- jumlah yang sudah bayar
- jumlah yang belum bayar
- total terkumpul

### 6.4 Pengeluaran Mess
Untuk biaya di luar kantong duafa.

Submenu:
- daftar pengeluaran
- tambah catatan pengeluaran
- detail
- status talangan/penggantian

### 6.5 Usulan & Voting
Dipakai anggota untuk mengusulkan hal tertentu.

Contoh:
- beli sapu
- beli perlengkapan
- usulan keputusan bersama

Field:
- judul
- deskripsi
- pembuat usulan
- masa voting
- hasil voting
- status akhir

### 6.6 Kontribusi
Submenu:
- tambah kontribusi
- feed kontribusi
- leaderboard
- riwayat pribadi

### 6.7 Feed / Info Mess
Satu feed untuk post sementara:
- makanan bebas ambil
- siapa mau nasi
- ini makanan siapa?
- pengumuman sesaat

Interaksi:
- reaction
- komentar
- ikut / mau / ambil
- penanda expire

### 6.8 Anggota
Submenu:
- daftar anggota
- detail anggota
- status aktif/nonaktif
- role
- statistik ringan

### 6.9 Profil
Isi:
- data akun
- role
- histori pembayaran wifi
- histori kontribusi
- histori usulan

### 6.10 Pengaturan
Untuk admin:
- nominal wifi default
- deadline wifi default
- kategori transaksi
- kategori kontribusi
- masa aktif default feed
- daftar role

---

## 7. Fitur Detail per Modul

### 7.1 Modul Auth
#### Fitur
- register akun oleh admin atau undangan
- login email/password
- logout
- reset password
- proteksi route berdasarkan role

#### Catatan keputusan
Lebih aman kalau:
- anggota tidak bisa daftar bebas
- akun dibuat admin atau melalui invite

### 7.2 Modul Anggota
#### Fitur
- tambah anggota
- edit anggota
- nonaktifkan anggota
- ubah role
- lihat riwayat partisipasi

#### Catatan
Karena anggota dinamis, jangan hardcode 15 orang.

Tambahkan field:
- active status
- joined_at
- left_at opsional

### 7.3 Modul Kantong Duafa
#### Fitur
- create transaksi
- edit transaksi
- hapus transaksi terbatas
- upload bukti
- lihat saldo berjalan
- filter per bulan
- audit trail pembuat dan pengubah

#### Aturan
- hanya admin/bendahara yang bisa membuat transaksi resmi
- anggota hanya bisa mengusulkan transaksi jika nanti mau dikembangkan
- transaksi wajib punya kategori

#### Data minimum
- tanggal
- tipe
- kategori
- deskripsi
- nominal
- bukti opsional
- created_by
- verified_by opsional

### 7.4 Modul Wifi
#### Fitur
- generate tagihan bulanan
- nominal default Rp20.000 per anggota aktif
- anggota upload bukti transfer
- admin/bendahara verifikasi
- status per orang
- rekap bulanan

#### Aturan
- hanya anggota aktif pada bulan itu yang ditagih
- bukti transfer disimpan sebagai lampiran
- verifikasi menghasilkan status final

#### Workflow
1. sistem generate tagihan bulan baru
2. tiap anggota melihat kewajibannya
3. anggota upload bukti
4. status jadi “menunggu verifikasi”
5. admin/bendahara approve atau reject
6. jika approve, masuk rekap pembayaran wifi bulan itu

### 7.5 Modul Pengeluaran Non-Kas
#### Fitur
- catat siapa bayar
- catat siapa nalangi
- nominal
- kategori
- status penggantian
- lampiran opsional

#### Tujuan
- bukan saldo kas
- tapi transparansi pengeluaran bersama

### 7.6 Modul Usulan & Voting
#### Fitur
- anggota buat usulan
- anggota vote setuju / tidak
- admin lihat hasil
- admin putuskan status akhir

#### Rule sederhana v1
- satu user satu vote
- voting dibuka sampai tanggal tertentu
- hasil persentase ditampilkan
- keputusan final tetap di admin

Status:
- draft
- aktif
- selesai
- disetujui
- ditolak

### 7.7 Modul Kontribusi
#### Fitur
- anggota atau admin membuat catatan kontribusi
- kategori kontribusi
- bisa diberi poin
- leaderboard

#### Contoh kategori
- buang sampah
- bersih dapur
- beres area umum
- buang barang
- bantu logistik

#### Rule v1
- satu kontribusi = satu entri
- poin default per kontribusi bisa sama dulu
- nanti baru dikembangkan kalau perlu bobot berbeda

### 7.8 Modul Feed / Info Sementara
#### Tipe post
- makanan
- nasi
- tanya kepemilikan
- pengumuman
- lainnya

#### Fitur
- buat post
- set expired_at
- komentar
- reaction
- history post expired

#### Alur khusus
##### Makanan bebas ambil
User post:
- ada makanan
- lokasi
- catatan
- bebas ambil

User lain:
- komentar “aku ambil”
- reaction

##### Masak nasi
User post:
- mau masak nasi
- user lain klik “mau”
- pembuat post bisa lihat estimasi jumlah orang

##### Tanya kepemilikan
User post:
- “ini makanan siapa?”
- beri batas waktu
- jika tidak ada respon sampai expire, status bisa ditandai:
  - **milik bersama**

Ini bisa dibuat manual dulu oleh pembuat/admin di v1, belum perlu otomatis kompleks.

---

## 8. Non-Functional Requirements

### Performa
- mobile-first
- halaman utama cepat dibuka
- query sederhana dan efisien

### Keamanan
- password di-hash
- access control per role
- validasi upload file
- audit log untuk aksi sensitif

### Ketersediaan
- berjalan di VPS yang sudah ada
- backup database berkala

### UX
- sederhana
- tidak terlalu banyak klik
- nyaman di HP Android

### PWA
- installable
- app icon
- splash screen
- offline fallback minimal untuk shell app

---

## 9. Rekomendasi Stack

### Frontend
- **SvelteKit**
- TailwindCSS
- PWA plugin
- form validation ringan
- state management secukupnya, jangan over-engineering

### Backend
- **Go**
- **Fiber**
- PostgreSQL
- JWT auth atau session token sederhana
- object storage lokal / S3-compatible untuk upload bukti

### Deployment
- frontend dan backend di VPS yang sama
- reverse proxy via Nginx
- domain/subdomain terpisah jika perlu

Contoh:
- `messhub.domain.com` untuk frontend
- `api.messhub.domain.com` untuk backend

---

## 10. Arsitektur Tingkat Tinggi

### Frontend
Tanggung jawab:
- UI
- auth state
- list/detail page
- form upload bukti
- install PWA

### Backend API
Tanggung jawab:
- autentikasi
- otorisasi role
- CRUD semua modul
- verifikasi wifi
- voting
- feed/comment/reaction
- file handling

### Database
Tanggung jawab:
- source of truth
- auditability
- query laporan sederhana

---

## 11. Struktur Entitas Data

### users
- id
- name
- email
- password_hash
- role
- is_active
- joined_at
- left_at
- created_at
- updated_at

### roles
Bisa hardcoded dulu:
- admin
- treasurer
- member

### wallet_transactions
Untuk kantong duafa.
- id
- transaction_date
- type (`income`, `expense`)
- category
- description
- amount
- proof_url
- created_by
- updated_by
- created_at
- updated_at

### wifi_bills
Satu tagihan bulanan global.
- id
- month
- year
- nominal_per_person
- deadline_date
- status
- created_by
- created_at

### wifi_bill_members
Status per anggota per bulan.
- id
- wifi_bill_id
- user_id
- amount
- payment_status
- proof_url
- submitted_at
- verified_at
- verified_by
- rejection_reason

### shared_expenses
Pengeluaran non-kas.
- id
- expense_date
- category
- description
- amount
- paid_by_user_id
- status
- notes
- proof_url
- created_by
- created_at

### proposals
- id
- title
- description
- created_by
- voting_start
- voting_end
- status
- final_decision_by
- final_decision_note
- created_at
- updated_at

### proposal_votes
- id
- proposal_id
- user_id
- vote_type
- created_at

### contributions
- id
- user_id
- category
- description
- points
- contributed_at
- created_by
- created_at

### posts
Feed/info sementara.
- id
- type
- title
- content
- location_note
- created_by
- expires_at
- status
- created_at
- updated_at

### post_reactions
- id
- post_id
- user_id
- reaction_type
- created_at

### post_comments
- id
- post_id
- user_id
- comment
- created_at
- updated_at

### notifications
- id
- user_id
- title
- body
- type
- reference_type
- reference_id
- is_read
- created_at

### audit_logs
- id
- user_id
- action
- entity_type
- entity_id
- old_value
- new_value
- created_at

---

## 12. Status dan Enum yang Disarankan

### user role
- admin
- treasurer
- member

### wallet transaction type
- income
- expense

### wifi payment status
- unpaid
- pending_verification
- verified
- rejected

### proposal vote type
- agree
- disagree

### proposal status
- active
- closed
- approved
- rejected

### shared expense status
- personal
- fronted
- partially_reimbursed
- reimbursed

### post type
- food
- rice
- ownership_question
- announcement
- other

### post status
- active
- expired
- archived

---

## 13. API Domain Map

### Auth API
- login
- logout
- me
- reset password

### Users API
- list users
- create user
- update user
- change role
- deactivate user

### Wallet API
- list transactions
- create transaction
- update transaction
- get balance

### Wifi API
- create monthly bill
- list bill members
- submit payment proof
- verify payment
- reject payment
- monthly recap

### Shared Expenses API
- list expenses
- create expense
- update expense

### Proposals API
- create proposal
- list proposals
- vote proposal
- close proposal
- finalize proposal

### Contributions API
- create contribution
- list contributions
- leaderboard

### Feed API
- create post
- list active posts
- list history
- comment
- react
- expire post

### Notifications API
- list notifications
- mark as read

---

## 14. Flow Utama Pengguna

### 14.1 Flow Login
1. user buka app
2. login email/password
3. sistem cek role
4. masuk ke dashboard sesuai akses

### 14.2 Flow Pembayaran Wifi
1. anggota buka menu wifi
2. lihat tagihan bulan aktif
3. upload bukti transfer
4. status jadi menunggu verifikasi
5. admin/bendahara menerima notifikasi
6. admin/bendahara approve atau reject
7. anggota melihat status final

### 14.3 Flow Pencatatan Kantong Duafa
1. admin/bendahara buka kantong duafa
2. tambah transaksi
3. isi tipe, kategori, nominal, deskripsi, bukti opsional
4. sistem simpan dan update saldo berjalan

### 14.4 Flow Pengeluaran Talangan
1. admin/bendahara buat catatan pengeluaran
2. isi siapa bayar / nalangi
3. isi status
4. data tampil di daftar transparansi pengeluaran

### 14.5 Flow Usulan
1. anggota buat usulan
2. usulan tampil di menu voting
3. anggota lain vote
4. admin lihat hasil
5. admin finalisasi

### 14.6 Flow Kontribusi
1. user tambah kontribusi
2. sistem simpan
3. leaderboard terupdate

### 14.7 Flow Masak Nasi
1. user buat post tipe rice
2. user lain klik reaction “mau”
3. pembuat post lihat jumlah peminat
4. setelah expire, post masuk history

---

## 15. Prioritas Pengerjaan

### Phase 1 — Fondasi
- setup repo frontend/backend
- auth
- user management
- role guard
- layout dashboard
- PWA setup

### Phase 2 — Core Finance
- kantong duafa
- wifi bill monthly
- upload bukti transfer
- verifikasi admin/bendahara
- shared expenses

### Phase 3 — Social Operation
- kontribusi
- usulan & voting
- feed/info sementara
- komentar & reaction

### Phase 4 — Finishing
- notifikasi in-app
- leaderboard
- filter laporan
- audit log
- polish UI/UX

---

## 16. MVP yang Disarankan

Supaya cepat jadi, definisikan **MVP** seperti ini:

### MVP wajib
- login
- manajemen anggota
- kantong duafa CRUD
- wifi bulanan
- upload bukti
- verifikasi pembayaran
- dashboard ringkas
- pengeluaran non-kas
- feed info sederhana
- kontribusi

### Belum wajib di MVP
- voting
- komentar kompleks
- reaction banyak jenis
- leaderboard canggih
- history analytics lengkap

Kalau mau super realistis, voting bahkan bisa masuk **v1.1**.

---

## 17. UI Sitemap Sederhana

- `/login`
- `/`
- `/wallet`
- `/wallet/new`
- `/wifi`
- `/wifi/:billId`
- `/shared-expenses`
- `/contributions`
- `/feed`
- `/proposals`
- `/members`
- `/profile`
- `/settings`

---

## 18. Prinsip UX yang Harus Dijaga

- semua fitur utama bisa diakses maksimal 2–3 tap dari dashboard
- tombol aksi besar dan jelas di HP
- status penting pakai badge
- bukti transfer mudah diupload
- jangan terlalu banyak form field wajib
- feed dibuat familiar seperti social feed sederhana

---

## 19. Risiko dan Mitigasi

### Risiko: scope membesar
Mitigasi:
- pegang phase
- jangan semua dikerjakan sekaligus

### Risiko: data kacau karena banyak edit
Mitigasi:
- audit log
- batasi edit/hapus ke admin & bendahara

### Risiko: anggota malas pakai
Mitigasi:
- UI simpel
- fokus manfaat harian
- dashboard jelas
- notifikasi internal

### Risiko: upload file berantakan
Mitigasi:
- struktur folder file
- validasi ukuran dan tipe file

---

## 20. Definisi Selesai untuk v1

MessHub dianggap **layak pakai v1** jika:
- bisa login dari HP
- bisa di-install sebagai PWA
- admin bisa kelola anggota
- admin/bendahara bisa catat kas
- anggota bisa submit bukti wifi
- admin/bendahara bisa verifikasi wifi
- pengeluaran non-kas bisa dicatat
- kontribusi bisa dicatat
- feed info sementara bisa dipakai
- dashboard menampilkan ringkasan utama

---

## 21. Brief Singkat untuk Agent / Developer

> Bangun aplikasi PWA bernama MessHub untuk Mess Traspac Menyala. Stack: SvelteKit frontend, Go Fiber backend, PostgreSQL database. Sistem memiliki role admin, treasurer, dan member. Fitur utama: auth email/password, manajemen anggota dinamis, pencatatan kantong duafa, tagihan wifi bulanan Rp20.000/orang dengan deadline sebelum tanggal 10 dan upload bukti transfer yang diverifikasi admin/bendahara, pencatatan pengeluaran non-kas/talangan, modul usulan dan voting sederhana, pencatatan kontribusi penghuni, serta feed info sementara untuk makanan, nasi, dan pengumuman yang bisa expire dan masuk history. Fokus v1 adalah mobile-first, installable PWA, sederhana, transparan, dan realistis untuk dipakai harian di mess.

---

## 22. Saran Implementasi Nyata

Urutan implementasi yang disarankan:
1. auth + role
2. member management
3. wallet transactions
4. wifi billing + upload proof + verify
5. dashboard
6. shared expenses
7. contributions
8. feed/info
9. proposals/voting
10. notifikasi & audit log
