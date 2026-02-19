# GEMBOK Backend - Go Implementation

Backend API untuk GEMBOK ISP Management System yang dibangun dengan Go dan Clean Architecture.

## 🏗️ Architecture

```
gembok-backend/
├── cmd/server/              # Application entry point
├── internal/
│   ├── domain/              # Domain layer (entities, repositories)
│   ├── usecase/             # Business logic
│   ├── infrastructure/      # External implementations (DB, APIs)
│   └── interface/           # HTTP handlers, DTOs, middleware
└── pkg/                     # Shared packages (config, logger, utils)
```

## 📦 Tech Stack

- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL
- **Auth**: JWT
- **Config**: Viper
- **Logger**: Zap
- **Password Hashing**: bcrypt

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- MySQL 5.7 or higher

### Installation

```bash
# Install dependencies
cd gembok-backend
go mod tidy

# Setup configuration
cp configs/config.yaml.example configs/config.yaml
# Edit configs/config.yaml with your settings

# Setup environment variables
cp .env.example .env
# Edit .env with your database credentials

# Run database migrations
cd database/migrations
./migrate.sh up

# Or manually using MySQL
mysql -u root -p gembok_db < 20240119120000_initial_schema.sql

# Go back to project root and run server
cd ../..
go run cmd/server/main.go
```

### Database Migrations

See `database/migrations/README.md` for complete documentation.

Quick reference:
```bash
# Run UP migration (create tables)
cd database/migrations
./migrate.sh up

# Check migration status
./migrate.sh status

# Run DOWN migration (drop tables) - USE WITH CAUTION!
./migrate.sh down
```

Default admin user:
- Username: `admin`
- Password: `admin123`
- **IMPORTANT**: Change this password after first login!

### Running the Server

```bash
# Development
go run cmd/server/main.go

# Production
go build -o gembok-backend cmd/server/main.go
./gembok-backend
```

## 📊 API Endpoints

### Dashboard
- `GET /api/dashboard` - Get dashboard statistics

### Customers
- `GET /api/customers` - Get all customers (with pagination)
- `GET /api/customers/:id` - Get customer by ID
- `POST /api/customers` - Create new customer
- `PUT /api/customers/:id` - Update customer
- `DELETE /api/customers/:id` - Delete customer
- `POST /api/customers/:id/isolate` - Isolate customer (change MikroTik profile)
- `POST /api/customers/:id/activate` - Activate customer (change MikroTik profile)
- `POST /api/customers/bulk-isolate` - Bulk isolate customers
- `POST /api/customers/bulk-activate` - Bulk activate customers
- `POST /api/customers/:id/sync` - Sync customer to MikroTik

### Invoices
- `GET /api/invoices` - Get all invoices (with pagination)
- `GET /api/invoices/:id` - Get invoice by ID
- `POST /api/invoices` - Create new invoice
- `PUT /api/invoices/:id` - Update invoice
- `DELETE /api/invoices/:id` - Delete invoice

### Routers
- `GET /api/routers` - Get all routers
- `GET /api/routers/:id` - Get router by ID
- `POST /api/routers` - Create new router
- `PUT /api/routers/:id` - Update router
- `DELETE /api/routers/:id` - Delete router
- `GET /api/routers/active` - Get active router
- `POST /api/routers/:id/test` - Test router connection
- `PUT /api/routers/:id/activate` - Set router as active
- `GET /api/routers/:id/status` - Get router status
- `GET /api/routers/status/all` - Get all routers status

### MikroTik PPPoE
- `GET /api/mikrotik/ppp/users` - Get all PPPoE users
- `POST /api/mikrotik/ppp/users` - Add PPPoE user
- `PUT /api/mikrotik/ppp/users/:username` - Update PPPoE user
- `DELETE /api/mikrotik/ppp/users/:username` - Remove PPPoE user
- `GET /api/mikrotik/ppp/active` - Get active PPPoE sessions
- `GET /api/mikrotik/ppp/profiles` - Get PPPoE profiles
- `POST /api/mikrotik/ppp/users/:username/disconnect` - Disconnect PPPoE user

### GenieACS
- `GET /api/genieacs/devices` - Get all GenieACS devices
- `GET /api/genieacs/devices/:serial` - Get device by serial
- `GET /api/genieacs/devices/find?username=xxx` - Find device by PPPoE username
- `POST /api/genieacs/devices/reboot` - Reboot device
- `POST /api/genieacs/devices/parameter` - Set device parameter

## ⚙️ Configuration

### Config File (configs/config.yaml)

```yaml
server:
  port: "8080"
  mode: "debug"  # or "release"

database:
  host: "localhost"
  port: 3306
  name: "gembok_db"
  user: "root"
  password: "your-password"
  max_idle_conns: 10
  max_open_conns: 100
  max_lifetime: 3600

jwt:
  secret: "your-secret-key"
  expiration: 3600s

mikrotik:
  host: "192.168.1.1"
  user: "admin"
  password: "admin"
  port: 8728

genieacs:
  url: "http://localhost:7557"
  username: "admin"
  password: "admin"

whatsapp:
  api_url: "https://api.fonnte.com/send"
  token: "your-whatsapp-token"

tripay:
  api_key: "your-tripay-api-key"
  private_key: "your-tripay-private-key"
  merchant_code: "your-merchant-code"
  mode: "sandbox"  # or "production"

app:
  name: "GEMBOK ISP Management"
  version: "1.0.0"
  url: "http://localhost:8080"
```

### Environment Variables (.env)

Copy `.env.example` to `.env` and configure:

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gembok_db

# Server
SERVER_PORT=8080
SERVER_MODE=debug

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRATION=3600

# GenieACS
GENIEACS_URL=http://localhost:7557
GENIEACS_USERNAME=admin
GENIEACS_PASSWORD=admin

# Tripay
TRIPAY_API_KEY=your-api-key
TRIPAY_PRIVATE_KEY=your-private-key
TRIPAY_MERCHANT_CODE=your-code
TRIPAY_MODE=sandbox
```

## 🔒 Authentication

JWT-based authentication with the following flow:

1. Send login request
2. Receive JWT token
3. Include token in Authorization header: `Bearer <token>`

## 📝 Features Implemented

### Core Features
- ✅ Clean Architecture structure
- ✅ Domain entities
- ✅ Repository pattern
- ✅ Use cases (business logic)
- ✅ HTTP handlers
- ✅ Middleware (CORS, Logging, Auth)
- ✅ Database migrations
- ✅ Configuration management
- ✅ Structured logging

### Database
- ✅ Complete database schema
- ✅ Migration scripts (Up/Down)
- ✅ Migration bash script with safety checks
- ✅ Default admin user (username: admin, password: admin123)
- ✅ Default settings

### Customer Management
- ✅ CRUD operations for customers
- ✅ Customer pagination and search
- ✅ Customer isolation (change MikroTik profile)
- ✅ Customer activation (change MikroTik profile)
- ✅ Bulk isolation/activation
- ✅ Sync customer to MikroTik

### MikroTik Integration
- ✅ Multi-router support
- ✅ Connection pooling
- ✅ PPPoE user management (CRUD)
- ✅ Active PPPoE sessions monitoring
- ✅ PPPoE profiles listing
- ✅ User disconnection
- ✅ Router status monitoring
- ✅ Connection testing
- ✅ Set active router

### GenieACS Integration
- ✅ Device listing
- ✅ Device retrieval by serial
- ✅ Device info with status
- ✅ Find device by PPPoE username
- ✅ Device reboot
- ✅ Parameter setting
- ✅ WiFi SSID/password access
- ✅ ONU power monitoring

## 🚧 TODO

- [ ] Admin authentication (JWT)
- [ ] Payment Gateway (Tripay) integration
- [ ] ONU locations management
- [ ] WiFi settings control
- [ ] WhatsApp notifications
- [ ] Portal password management
- [ ] Hotspot management
- [ ] Cron job scheduler
- [ ] API documentation (Swagger)
- [ ] Unit tests
- [ ] Integration tests

## 📈 Performance

Compared to PHP:
- Faster response time
- Lower memory usage
- Better concurrency handling
- Type-safe code

## 🔧 Development

### Adding New Feature

1. Add entity to `internal/domain/entities/`
2. Add repository interface to `internal/domain/repositories/`
3. Implement repository in `internal/infrastructure/repositories/`
4. Create use case in `internal/usecase/`
5. Create handler in `internal/interface/http/handlers/`
6. Add route in `internal/interface/http/router.go`

### Running Tests

```bash
go test ./...
```

## 📞 Support

For issues or questions, please contact:
- GitHub: https://github.com/alijayanet/gembok-backend
- Email: alijayanet@gmail.com

## 📄 License

MIT License
