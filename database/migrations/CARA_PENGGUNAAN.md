# Cara Menggunakan Migrasi Database

## Struktur File Migrasi

Database migrasi terletak di `database/migrations/` dengan struktur:

```
database/migrations/
├── README.md                           # Dokumentasi lengkap
├── migrate.sh                          # Script bash untuk eksekusi
├── 20240119120000_initial_schema.sql  # File migrasi lengkap (up + down)
└── separate/                            # File terpisah (opsional)
    ├── 20240119120000_initial_schema.up.sql
    └── 20240119120000_initial_schema.down.sql
```

## Cara Menjalankan Migrasi

### 1. Menggunakan Script Bash (Disarankan)

```bash
# Masuk ke direktori migrasi
cd gembok-backend/database/migrations

# Jalankan migrasi UP (buat tabel)
./migrate.sh up

# Cek status migrasi
./migrate.sh status

# Jalankan migrasi DOWN (hapus tabel) - HATI-HATI!
./migrate.sh down
```

### 2. Menggunakan MySQL Command Line

```bash
# Masuk ke direktori migrasi
cd gembok-backend/database/migrations

# Jalankan migrasi lengkap
mysql -u root -p gembok_db < 20240119120000_initial_schema.sql

# Atau jalankan hanya bagian UP
mysql -u root -p gembok_db < separate/20240119120000_initial_schema.up.sql
```

### 3. Menggunakan Docker Compose

```bash
# Jika database berjalan di Docker
docker-compose exec db mysql -u root -prootpassword gembok_db < database/migrations/20240119120000_initial_schema.sql
```

### 4. Menggunakan PHPMyAdmin

1. Buka PHPMyAdmin di browser
2. Pilih database `gembok_db`
3. Klik tab "Import"
4. Pilih file migrasi SQL
5. Klik tombol "Go"

## Environment Variables

Sebelum menjalankan aplikasi, pastikan untuk mengkonfigurasi environment variables. Copy file `.env.example` ke `.env`:

```bash
cp .env.example .env
```

Edit file `.env` sesuai dengan konfigurasi database Anda:

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gembok_db

# ... konfigurasi lainnya
```

## Default Data

Setelah migrasi berhasil, beberapa data default akan dibuat:

### Default Admin User
- **Username**: `admin`
- **Password**: `admin123`
- **Email**: `admin@gembok.com`
- **Role**: `admin`

> **PENTING**: Ganti password default ini setelah login pertama!

### Default Settings
Beberapa setting default tersimpan di tabel `settings`:
- `MIKROTIK_*` - Pengaturan router MikroTik
- `GENIEACS_*` - Pengaturan GenieACS server
- `TRIPAY_*` - Pengaturan Tripay payment gateway
- `INVOICE_PREFIX`, `INVOICE_START` - Pengaturan nomor invoice
- `CURRENCY_SYMBOL` - Simbol mata uang (Rp)

## Struktur Database

### Tabel Utama

1. **admin_users** - Akun admin
2. **packages** - Paket internet
3. **routers** - Router MikroTik (support multi-router)
4. **customers** - Data pelanggan
5. **invoices** - Tagihan/invoice
6. **onu_locations** - Lokasi perangkat ONU
7. **trouble_tickets** - Tiket support
8. **settings** - Pengaturan aplikasi
9. **cron_schedules** - Jadwal tugas cron
10. **cron_logs** - Log eksekusi cron
11. **webhook_logs** - Log webhook

### Relasi Tabel

- `customers.package_id` → `packages.id`
- `customers.router_id` → `routers.id`
- `invoices.customer_id` → `customers.id`
- `onu_locations.customer_id` → `customers.id`
- `onu_locations.router_id` → `routers.id`
- `trouble_tickets.customer_id` → `customers.id`
- `cron_logs.schedule_id` → `cron_schedules.id`

## Troubleshooting

### Error: Access Denied
```
ERROR 1045 (28000): Access denied for user 'root'@'localhost'
```
**Solusi**: Pastikan username dan password MySQL benar di file `.env`

### Error: Unknown Database
```
ERROR 1049 (42000): Unknown database 'gembok_db'
```
**Solusi**: Buat database terlebih dahulu:
```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS gembok_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"
```

### Error: Foreign Key Constraint
```
ERROR 1215 (HY000): Cannot add foreign key constraint
```
**Solusi**: Pastikan migrasi dijalankan secara berurutan, dan tabel referensi sudah ada.

## Menambah Migrasi Baru

Untuk menambah migrasi baru:

1. Buat file SQL baru dengan nama format: `YYYYMMDDHHMMSS_description.sql`
2. Tambahkan bagian `-- Up` untuk perubahan
3. Tambahkan bagian `-- Down` untuk rollback
4. Jalankan migrasi baru

Contoh:
```sql
-- Migration: Add referral_code to customers
-- Up

ALTER TABLE customers ADD COLUMN referral_code VARCHAR(50) DEFAULT NULL AFTER onu_ip_address;

-- Down

ALTER TABLE customers DROP COLUMN referral_code;
```

## Notes

- Semua tabel menggunakan engine InnoDB
- Semua text field menggunakan utf8mb4_unicode_ci collation
- Semua datetime field memiliki presisi (3) untuk milidetik
- Foreign keys cascade on update
- Foreign keys pada tabel customer cascade on delete
- Gunakan script `migrate.sh` untuk kemudahan dan safety checks
