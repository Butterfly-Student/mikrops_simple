//go:build integration

package impl

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	mysqlContainer, err := mysqlcontainer.Run(ctx,
		"mysql:8.0",
		mysqlcontainer.WithDatabase("gembok_test"),
		mysqlcontainer.WithUsername("test"),
		mysqlcontainer.WithPassword("test"),
	)
	if err != nil {
		fmt.Printf("Failed to start MySQL container: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			fmt.Printf("Failed to terminate container: %v\n", err)
		}
	}()

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		fmt.Printf("Failed to get container host: %v\n", err)
		os.Exit(1)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		fmt.Printf("Failed to get container port: %v\n", err)
		os.Exit(1)
	}

	dsn := fmt.Sprintf("test:test@tcp(%s:%s)/gembok_test?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())

	// Retry connection a few times since MySQL might need a moment
	for i := 0; i < 10; i++ {
		testDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}

	// Create tables using raw SQL to avoid GORM longtext+unique index issues
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS admin_users (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			role VARCHAR(50) DEFAULT 'admin',
			status VARCHAR(50) DEFAULT 'active',
			created_at DATETIME(3),
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS packages (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			price DOUBLE NOT NULL,
			speed VARCHAR(100),
			description TEXT,
			profile_normal VARCHAR(255),
			profile_isolir VARCHAR(255),
			status VARCHAR(50) DEFAULT 'active',
			created_at DATETIME(3),
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS customers (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			phone VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(255),
			address TEXT,
			package_id BIGINT UNSIGNED,
			pppoe_username VARCHAR(255) UNIQUE,
			pppoe_password VARCHAR(255) NOT NULL,
			status VARCHAR(50) DEFAULT 'active',
			router_id BIGINT UNSIGNED DEFAULT 0,
			onu_id VARCHAR(255),
			onu_serial VARCHAR(255),
			onu_mac_address VARCHAR(255),
			onu_ip_address VARCHAR(255),
			latitude DOUBLE DEFAULT 0,
			longitude DOUBLE DEFAULT 0,
			isolation_date DATETIME(3),
			activation_date DATETIME(3),
			created_at DATETIME(3),
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS invoices (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			customer_id BIGINT UNSIGNED NOT NULL,
			number VARCHAR(255) NOT NULL UNIQUE,
			amount DOUBLE NOT NULL,
			period VARCHAR(50) NOT NULL,
			due_date DATETIME(3),
			status VARCHAR(50) DEFAULT 'unpaid',
			paid_at DATETIME(3),
			payment_method VARCHAR(100),
			payment_reference VARCHAR(255),
			created_at DATETIME(3),
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS routers (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			host VARCHAR(255) NOT NULL,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			port INT DEFAULT 8728,
			is_active BOOLEAN DEFAULT FALSE,
			created_at DATETIME(3),
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS onu_locations (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			customer_id BIGINT UNSIGNED NOT NULL,
			router_id BIGINT UNSIGNED DEFAULT 0,
			onu_id VARCHAR(255) UNIQUE,
			serial_number VARCHAR(255),
			mac_address VARCHAR(255),
			ip_address VARCHAR(255),
			port_number VARCHAR(50),
			signal_strength VARCHAR(50),
			latitude DOUBLE DEFAULT 0,
			longitude DOUBLE DEFAULT 0,
			address TEXT,
			notes TEXT,
			created_at DATETIME(3),
			updated_at DATETIME(3),
			INDEX idx_onu_locations_customer_id (customer_id)
		)`,
		`CREATE TABLE IF NOT EXISTS trouble_tickets (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			customer_id BIGINT UNSIGNED NOT NULL,
			subject VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			priority VARCHAR(50) DEFAULT 'medium',
			status VARCHAR(50) DEFAULT 'open',
			assigned_to VARCHAR(255),
			resolved_at DATETIME(3),
			created_at DATETIME(3),
			updated_at DATETIME(3),
			INDEX idx_trouble_tickets_customer_id (customer_id)
		)`,
		`CREATE TABLE IF NOT EXISTS settings (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			setting_key VARCHAR(255) NOT NULL UNIQUE,
			setting_value TEXT,
			description TEXT,
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS cron_schedules (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			task_type VARCHAR(255) NOT NULL,
			schedule_time VARCHAR(100),
			schedule_days VARCHAR(255),
			is_active BOOLEAN DEFAULT TRUE,
			last_run_at DATETIME(3),
			next_run_at DATETIME(3),
			created_at DATETIME(3),
			updated_at DATETIME(3)
		)`,
		`CREATE TABLE IF NOT EXISTS cron_logs (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			schedule_id BIGINT UNSIGNED,
			task_type VARCHAR(255),
			status VARCHAR(50),
			output TEXT,
			error TEXT,
			created_at DATETIME(3),
			INDEX idx_cron_logs_schedule_id (schedule_id)
		)`,
		`CREATE TABLE IF NOT EXISTS webhook_logs (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			event VARCHAR(255) NOT NULL,
			url VARCHAR(255) NOT NULL,
			payload LONGTEXT,
			response LONGTEXT,
			status_code INT,
			duration INT,
			created_at DATETIME(3),
			INDEX idx_webhook_logs_event (event)
		)`,
	}

	for _, schema := range schemas {
		if result := testDB.Exec(schema); result.Error != nil {
			fmt.Printf("Failed to create table: %v\n", result.Error)
			os.Exit(1)
		}
	}

	os.Exit(m.Run())
}

// cleanTable truncates a table between tests.
func cleanTable(t *testing.T, tableName string) {
	t.Helper()
	testDB.Exec("DELETE FROM " + tableName)
	testDB.Exec("ALTER TABLE " + tableName + " AUTO_INCREMENT = 1")
}
