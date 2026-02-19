# Database Migrations

This directory contains SQL migration files for the GEMBOK ISP Management System.

## Migration Files

### 20240119120000_initial_schema.sql
Initial schema creation with all required tables:
- `admin_users` - Admin user accounts
- `packages` - Internet packages
- `routers` - MikroTik routers
- `customers` - Customer accounts
- `invoices` - Invoice records
- `onu_locations` - ONU device locations
- `trouble_tickets` - Support tickets
- `settings` - Application settings
- `cron_schedules` - Scheduled tasks
- `cron_logs` - Cron execution logs
- `webhook_logs` - Webhook call logs

## How to Run Migrations

### Using MySQL Command Line
```bash
mysql -u root -p gembok_db < database/migrations/20240119120000_initial_schema.sql
```

### Using Docker Compose
```bash
# First, copy the migration file to the docker-entrypoint-initdb.d directory
docker-compose exec db mysql -u root -prootpassword gembok_db < database/migrations/20240119120000_initial_schema.sql
```

### Using PHPMyAdmin
1. Open PHPMyAdmin
2. Select your database
3. Click on "Import" tab
4. Choose the migration file
5. Click "Go"

## Default Data

### Default Admin User
- Username: `admin`
- Password: `admin123`
- Note: Change this password immediately after first login!

### Default Settings
The following settings are pre-configured:
- MIKROTIK_* - MikroTik router settings
- GENIEACS_* - GenieACS server settings
- TRIPAY_* - Tripay payment gateway settings
- DEFAULT_WHATSAPP_GATEWAY - Default WhatsApp gateway
- INVOICE_PREFIX, INVOICE_START - Invoice number settings
- CURRENCY_SYMBOL - Currency display

## Schema Overview

### Relationships
- `customers.package_id` → `packages.id`
- `customers.router_id` → `routers.id`
- `invoices.customer_id` → `customers.id`
- `onu_locations.customer_id` → `customers.id`
- `onu_locations.router_id` → `routers.id`
- `trouble_tickets.customer_id` → `customers.id`
- `cron_logs.schedule_id` → `cron_schedules.id`

### Indexes
All tables have appropriate indexes for common queries:
- Primary keys on all tables
- Unique indexes on usernames, phone numbers, etc.
- Foreign key indexes
- Status and date indexes for filtering

## Rollback

Each migration file includes a `-- Down` section that can be used to rollback changes.

```bash
# Run only the Down section
mysql -u root -p gembok_db < database/migrations/20240119120000_initial_schema.sql
```

Note: The Down section drops all tables, so use with caution!

## Adding New Migrations

1. Create a new SQL file with timestamp prefix: `YYYYMMDDHHMMSS_description.sql`
2. Include both `-- Up` and `-- Down` sections
3. Add your migration logic to the `-- Up` section
4. Add rollback logic to the `-- Down` section

Example:
```sql
-- Migration: Add new column
-- Up

ALTER TABLE customers ADD COLUMN referral_code VARCHAR(50) DEFAULT NULL;

-- Down

ALTER TABLE customers DROP COLUMN referral_code;
```

## Notes

- All tables use InnoDB engine
- All text fields use utf8mb4_unicode_ci collation
- All datetime fields use precision (3) for milliseconds
- Foreign keys are set to cascade on update
- Foreign keys on customer-related tables cascade on delete
