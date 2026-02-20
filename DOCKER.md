# 🐳 Docker Setup Guide

## Prerequisites
- Docker 20.10+
- Docker Compose 2.8+

## Quick Start

### Option 1: Without RabbitMQ (Memory Event Bus)
```bash
docker-compose up -d
```

### Option 2: With RabbitMQ (Distributed Event System)
```bash
docker-compose --profile with-rabbitmq up -d
```

## Services

| Service | Container Name | Port | Description |
|---------|---------------|------|-------------|
| MySQL | mikrops-mysql | 3306 | Database |
| Backend | mikrops-backend | 8080 | Go API server |
| RabbitMQ | mikrops-rabbitmq | 5672, 15672 | Message broker (optional) |

## Access Points

- **API Server**: http://localhost:8080
- **MySQL**: localhost:3306
- **RabbitMQ Management**: http://localhost:15672 (if enabled)

## Default Credentials

### Superadmin
- Username: `superadmin`
- Password: `admin123`
- Email: `superadmin@mikrops.local`

### Admin
- Username: `admin`
- Password: `admin123`
- Email: `admin@mikrops.local`

### MySQL
- Root Password: `root`
- Database: `mikrops_simple`
- User: `root`

### RabbitMQ (if enabled)
- Username: `guest`
- Password: `guest`

## Environment Variables

The backend service accepts these environment variables:

### Database
- `DB_HOST` - MySQL host (default: mysql)
- `DB_PORT` - MySQL port (default: 3306)
- `DB_NAME` - Database name (default: mikrops_simple)
- `DB_USER` - MySQL user (default: root)
- `DB_PASSWORD` - MySQL password (default: root)

### Server
- `SERVER_PORT` - Server port (default: 8080)
- `SERVER_MODE` - Server mode (debug/release)

### JWT
- `JWT_SECRET` - JWT secret key
- `JWT_EXPIRATION` - JWT token expiration

### RBAC
- `RBAC_EVENT_SYSTEM_TYPE` - Event system type (memory/rabbitmq)
- `RBAC_AUDIT_ENABLED` - Enable audit logging (true/false)
- `RBAC_AUDIT_LOG_LEVEL` - Audit log level (debug/info/error)
- `RBAC_DEFAULT_SUPERADMIN_USERNAME` - Default superadmin username
- `RBAC_DEFAULT_SUPERADMIN_PASSWORD` - Default superadmin password
- `RBAC_DEFAULT_SUPERADMIN_EMAIL` - Default superadmin email

## Common Commands

### Start services
```bash
docker-compose up -d
```

### Stop services
```bash
docker-compose down
```

### Stop and remove volumes
```bash
docker-compose down -v
```

### View logs
```bash
# All services
docker-compose logs -f

# Backend only
docker-compose logs -f backend

# MySQL only
docker-compose logs -f mysql
```

### Restart backend
```bash
docker-compose restart backend
```

### Rebuild after code changes
```bash
docker-compose up -d --build
```

### Access backend shell
```bash
docker-compose exec backend sh
```

### Access MySQL shell
```bash
docker-compose exec mysql mysql -u root -p
```

## Troubleshooting

### Backend can't connect to MySQL
1. Check MySQL health: `docker-compose ps`
2. View MySQL logs: `docker-compose logs mysql`
3. Wait for MySQL to be healthy before starting backend

### Container won't start
1. Check port conflicts: `lsof -i:8080` or `lsof -i:3306`
2. Stop conflicting services
3. Clean up: `docker system prune -a`

### Database errors
1. Rebuild from scratch: `docker-compose down -v && docker-compose up -d`
2. Check MySQL logs: `docker-compose logs mysql`
3. Run migration manually: `docker-compose exec backend sh -c "./mikrops_simple migrate"`

### Permission errors
1. Fix file permissions on host: `chmod +x database/init/01_init.sql`

### RBAC Not Working
1. Verify superadmin created: Check MySQL for admin_users table
2. Check audit_log table for entries
3. Check casbin_rule table is created
4. Login with superadmin to manage policies

## Development

### Hot Reload with Air (Optional)
```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with Air
docker-compose up mysql
docker-compose exec backend go run main.go
```

### Debugging with Delve (Optional)
```bash
# Build with debugging
go build -gcflags "all=-N -l" -o mikrops_debug cmd/server/main.go

# Run Delve inside container
docker-compose exec backend /mikrops_debug
```

## Production Deployment

1. Change default passwords in environment variables
2. Set `SERVER_MODE=release`
3. Use external MySQL instead of Docker MySQL
4. Use RabbitMQ for distributed event system
5. Enable HTTPS/TLS termination at load balancer
6. Set up health checks and monitoring
7. Configure backups for MySQL

## Backup

### Database Backup
```bash
docker-compose exec mysql mysqldump -u root -pmikrops_simple mikrops_simple > backup_$(date +%Y%m%d).sql
```

### Volume Backup
```bash
docker run --rm -v mikrops_mysql_data:/data -v $(pwd):/backup alpine tar czf backup.tar.gz /data
```

## Monitoring

### Check Container Health
```bash
docker-compose ps
```

### Resource Usage
```bash
docker stats
```

## API Testing

### Login as Superadmin
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "superadmin",
    "password": "admin123"
  }'
```

### Get Casbin Policies
```bash
# First login to get token
TOKEN="YOUR_JWT_TOKEN_HERE"

curl -X GET http://localhost:8080/api/casbin/policies \
  -H "Authorization: Bearer $TOKEN"
```

### Check Permission
```bash
curl -X POST http://localhost:8080/api/casbin/check \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin",
    "resource": "/api/dashboard",
    "action": "GET",
    "owner": "*"
  }'
```

## Security Notes

⚠️ **IMPORTANT**: Change default passwords before deploying to production!

### Recommended Actions
1. Change superadmin password on first login
2. Update `RBAC_DEFAULT_SUPERADMIN_PASSWORD` in docker-compose.yml
3. Use strong JWT_SECRET in production
4. Set `SERVER_MODE=release`
5. Enable firewall rules to restrict access
6. Use separate MySQL instance for production
7. Enable SSL/TLS for all endpoints
8. Set up rate limiting
9. Monitor audit logs regularly
10. Enable RabbitMQ for distributed deployments

## Support

For issues or questions, check:
1. Container logs: `docker-compose logs -f <service>`
2. Database connection: Verify MySQL is healthy
3. Port availability: Check if ports are already in use
4. RBAC setup: Ensure Casbin tables are created
