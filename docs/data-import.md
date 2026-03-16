# Data Import Guide

## Purpose

Dokumen ini menjelaskan cara memigrasikan data spreadsheet lama ke MessHub menggunakan alur CSV preview -> commit.

Fitur terkait:
- Admin import hub: `/admin/import`
- Member import: `/admin/import/members`
- Wallet import: `/admin/import/wallet`

## Export Google Sheets to CSV

1. Buka sheet yang ingin diexport.
2. Rapikan header kolom terlebih dahulu.
3. Pilih `File` -> `Download` -> `Comma-separated values (.csv, current sheet)`.
4. Simpan file hasil export tanpa mengubah isi numerik atau tanggal secara manual bila tidak perlu.

## Member Import

### Template resmi

File template tersedia di UI:
- `/templates/member-import-template.csv`

Header minimal:

```csv
name,email,role,is_active
```

Contoh:

```csv
name,email,role,is_active
Ryan Prayoga,ryan@example.com,admin,true
Budi Santoso,budi@example.com,member,true
```

### Mapping dari spreadsheet lama

Gunakan mapping berikut bila sumber lama belum sesuai template:

| Sumber lama | Template MessHub |
| --- | --- |
| `Nama` | `name` |
| `Email` | `email` |
| Role belum ada | isi `member` untuk penghuni biasa |
| Status aktif / catatan aktif | `is_active` dengan nilai `true` atau `false` |

Catatan:
- Role yang diterima: `admin`, `treasurer`, `member`
- Email harus unik
- `is_active` menerima bentuk umum seperti `true/false`, `ya/tidak`, `aktif/nonaktif`
- Saat commit, admin harus menentukan satu password sementara untuk akun baru yang berhasil diimpor

## Wallet Import

### Template resmi

File template tersedia di UI:
- `/templates/wallet-import-template.csv`

Header template:

```csv
transaction_date,description,income,expense,proof
```

Contoh:

```csv
transaction_date,description,income,expense,proof
2026-03-01,Hibah Orang Baik,50000,,https://contoh.com/bukti-hibah
2026-03-03,Bayar wifi bulanan,,20000,
```

### Mapping dari spreadsheet kas lama

Struktur spreadsheet lama yang umum:

| Spreadsheet lama | Template MessHub |
| --- | --- |
| `Tanggal` | `transaction_date` |
| `Deskripsi` | `description` |
| `Pemasukan (Rp)` | `income` |
| `Pengeluaran (Rp)` | `expense` |
| `Bukti` | `proof` |
| `Saldo` | tidak diimpor |

Catatan:
- Sistem menerima header template MessHub dan alias umum seperti `Tanggal`, `Deskripsi`, `Pemasukan`, `Pengeluaran`, dan `Bukti`
- Isi salah satu kolom `income` atau `expense`, bukan keduanya
- Nilai harus positif
- Kategori diinfer sederhana dari deskripsi:
  - `wifi` -> `wifi`
  - `hibah`, `donasi`, `sumbangan` -> `hibah`
  - `galon` -> `galon`
  - `plastik`, `sabun`, `kebersihan`, `pel`, `sapu` -> `kebersihan`
  - selain itu -> `lainnya`
- `proof` boleh kosong
- `Saldo` spreadsheet lama tidak dipakai sebagai source of truth

## Import Flow

1. Admin unggah CSV.
2. Backend parse file dan membuat preview job.
3. UI menampilkan:
   - total rows
   - valid rows
   - invalid rows
   - duplicate rows untuk member import
   - total income dan total expense untuk wallet import
4. Admin meninjau error atau warning per baris.
5. Admin commit import.
6. Sistem hanya menyimpan baris yang siap diimpor.

Tracking yang tercatat:
- `import_jobs`
- audit log:
  - `member_import_preview`
  - `member_import_commit`
  - `wallet_import_preview`
  - `wallet_import_commit`

## Example Preview Response

### Member preview

```json
{
  "message": "member import preview ready",
  "data": {
    "job_id": "9c50b2ab-9f76-4b50-8ed7-3fb5ff0b4b58",
    "file_name": "members.csv",
    "summary": {
      "total_rows": 3,
      "valid_rows": 1,
      "invalid_rows": 1,
      "duplicate_rows": 1,
      "importable_rows": 1
    },
    "rows": [
      {
        "row_number": 2,
        "status": "valid",
        "name": "Budi Santoso",
        "email": "budi@example.com",
        "role": "member",
        "normalized_role": "member",
        "is_active": "true",
        "normalized_is_active": true,
        "errors": [],
        "warnings": []
      },
      {
        "row_number": 3,
        "status": "duplicate",
        "name": "Ryan Prayoga",
        "email": "ryan@example.com",
        "role": "admin",
        "normalized_role": "admin",
        "is_active": "true",
        "normalized_is_active": true,
        "errors": [],
        "warnings": ["Email ini sudah terdaftar di MessHub."]
      }
    ],
    "warnings": [],
    "can_commit": true,
    "requires_temporary_password": true
  }
}
```

### Wallet preview

```json
{
  "message": "wallet import preview ready",
  "data": {
    "job_id": "b498fbd7-dbf7-47a4-84c8-332b98ef1426",
    "file_name": "wallet.csv",
    "summary": {
      "total_rows": 2,
      "valid_rows": 1,
      "invalid_rows": 1,
      "importable_rows": 1,
      "total_income": 50000,
      "total_expense": 0
    },
    "rows": [
      {
        "row_number": 2,
        "status": "valid",
        "transaction_date": "2026-03-01",
        "normalized_transaction_date": "2026-03-01",
        "description": "Hibah Orang Baik",
        "income": "50000",
        "expense": "",
        "type": "income",
        "amount": 50000,
        "category": "hibah",
        "proof": "https://contoh.com/bukti-hibah",
        "errors": [],
        "warnings": []
      }
    ],
    "warnings": [],
    "can_commit": true
  }
}
```

## Wallet Balance Rule

Saldo kas tidak diimpor dari spreadsheet lama.

Aturan yang dipakai MessHub:
- saldo dihitung ulang dari semua transaksi wallet yang tersimpan
- pemasukan menambah saldo
- pengeluaran mengurangi saldo

Dengan aturan ini, preview dan commit impor tetap aman walau sheet lama punya kolom `Saldo`.

## Production UI Cleanup Note

Pada STEP 9, helper text yang bersifat dev-only sudah dibersihkan dari UI production, termasuk:
- petunjuk seed account di halaman login
- referensi `.env` atau backend auth di UI user
- kartu dashboard yang menjelaskan detail fase internal
- copy admin/settings yang terlalu teknis untuk pemakaian harian

Pesan error yang tampil ke user juga dibuat lebih ramah, misalnya untuk kasus login gagal atau koneksi server terputus.
