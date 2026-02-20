.PHONY: help up down restart logs build clean db shell backup restore

help:
	@echo "MikroTik Simple - Docker Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  up           - Start all services"
	@echo "  up-rabbitmq - Start services with RabbitMQ"
	@echo "  down         - Stop all services"
	@echo "  restart      - Restart all services"
	@echo "  logs         - View logs (all services)"
	@echo "  logs-mysql   - View MySQL logs"
	@echo "  logs-backend - View backend logs"
	@echo "  logs-rabbitmq - View RabbitMQ logs"
	@echo "  build        - Rebuild backend image"
	@echo "  clean        - Remove all containers and volumes"
	@echo " db           - Connect to MySQL shell"
	@echo " shell        - Connect to backend container shell"
	@echo " backup       - Backup MySQL database"
	@ echo " restore      - Restore MySQL database from backup"
	@echo " status       - Show containers status"
	@echo " help         - Show this help message"
	@echo ""

up:
	docker compose up -d

up-rabbitmq:
	docker compose --profile with-rabbitmq up -d

down:
	docker compose down

restart:
	docker compose restart backend

logs:
	docker compose logs -f

logs-mysql:
	docker compose logs -f mysql

logs-backend:
	docker compose logs -f backend

logs-rabbitmq:
	docker compose logs -f rabbitmq

build:
	docker compose build backend

clean:
	docker compose down -v
	docker system prune -f
	@echo "✅ Cleaned containers and images"

db:
	docker exec mikrops-mysql mysql -u root -p

shell:
	docker exec mikrops-backend sh

backup:
	@echo "Creating backup..."
	docker exec mikrops-mysql mysqldump -u root -pmikrops_simple mikrops_simple > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup created: backup_$(shell date +%Y%m%d_%H%M%S).sql"

restore:
	@read -p "Enter backup file path: " backup_file
	docker exec -i mikrops-mysql mysql -u root -p mikrops_simple < $(backup_file)
	@echo "✅ Database restored"

status:
	@docker compose ps

health:
	@echo "Checking service health..."
	docker compose ps
	@echo ""
	@echo "MySQL Health Check:"
	docker exec mikrops-mysql mysqladmin ping -h localhost -u root -proot || echo "❌ MySQL not responding"
