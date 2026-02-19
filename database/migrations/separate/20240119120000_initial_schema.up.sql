-- Migration: Initial Schema - Up

CREATE TABLE IF NOT EXISTS `admin_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `role` varchar(50) DEFAULT 'admin',
  `status` varchar(50) DEFAULT 'active',
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admin_users_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `packages` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `price` double NOT NULL,
  `speed` varchar(50) DEFAULT NULL,
  `description` text DEFAULT NULL,
  `profile_normal` varchar(255) DEFAULT NULL COMMENT 'Profile name for active customers',
  `profile_isolir` varchar(255) DEFAULT NULL COMMENT 'Profile name for isolated customers',
  `status` varchar(50) DEFAULT 'active',
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_packages_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `routers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `host` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `port` int DEFAULT 8728,
  `is_active` tinyint(1) DEFAULT 0,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `customers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `address` text DEFAULT NULL,
  `package_id` bigint unsigned DEFAULT NULL,
  `pppoe_username` varchar(255) DEFAULT NULL,
  `pppoe_password` varchar(255) NOT NULL,
  `status` varchar(50) DEFAULT 'active',
  `router_id` bigint unsigned DEFAULT NULL,
  `onu_id` varchar(255) DEFAULT NULL,
  `onu_serial` varchar(255) DEFAULT NULL,
  `onu_mac_address` varchar(255) DEFAULT NULL,
  `onu_ip_address` varchar(45) DEFAULT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `isolation_date` datetime(3) DEFAULT NULL,
  `activation_date` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_customers_phone` (`phone`),
  UNIQUE KEY `idx_customers_pppoe_username` (`pppoe_username`),
  KEY `idx_customers_package_id` (`package_id`),
  KEY `idx_customers_router_id` (`router_id`),
  CONSTRAINT `fk_customers_package` FOREIGN KEY (`package_id`) REFERENCES `packages` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `fk_customers_router` FOREIGN KEY (`router_id`) REFERENCES `routers` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `invoices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` bigint unsigned NOT NULL,
  `number` varchar(255) NOT NULL,
  `amount` double NOT NULL,
  `period` varchar(50) NOT NULL,
  `due_date` datetime(3) NOT NULL,
  `status` varchar(50) DEFAULT 'unpaid',
  `paid_at` datetime(3) DEFAULT NULL,
  `payment_method` varchar(100) DEFAULT NULL,
  `payment_reference` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_invoices_number` (`number`),
  KEY `idx_invoices_customer_id` (`customer_id`),
  KEY `idx_invoices_status` (`status`),
  KEY `idx_invoices_due_date` (`due_date`),
  CONSTRAINT `fk_invoices_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `onu_locations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` bigint unsigned NOT NULL,
  `router_id` bigint unsigned DEFAULT NULL,
  `onu_id` varchar(255) DEFAULT NULL,
  `serial_number` varchar(255) DEFAULT NULL,
  `mac_address` varchar(255) DEFAULT NULL,
  `ip_address` varchar(45) DEFAULT NULL,
  `port_number` varchar(50) DEFAULT NULL,
  `signal_strength` varchar(50) DEFAULT NULL,
  `latitude` double DEFAULT NULL,
  `longitude` double DEFAULT NULL,
  `address` text DEFAULT NULL,
  `notes` text DEFAULT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_onu_locations_onu_id` (`onu_id`),
  KEY `idx_onu_locations_customer_id` (`customer_id`),
  KEY `idx_onu_locations_router_id` (`router_id`),
  CONSTRAINT `fk_onu_locations_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_onu_locations_router` FOREIGN KEY (`router_id`) REFERENCES `routers` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `trouble_tickets` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `customer_id` bigint unsigned NOT NULL,
  `subject` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `priority` varchar(50) DEFAULT 'medium',
  `status` varchar(50) DEFAULT 'open',
  `assigned_to` varchar(255) DEFAULT NULL,
  `resolved_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_trouble_tickets_customer_id` (`customer_id`),
  KEY `idx_trouble_tickets_status` (`status`),
  KEY `idx_trouble_tickets_priority` (`priority`),
  CONSTRAINT `fk_trouble_tickets_customer` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `setting_key` varchar(255) NOT NULL,
  `setting_value` text DEFAULT NULL,
  `description` text DEFAULT NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_settings_key` (`setting_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cron_schedules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `task_type` varchar(100) NOT NULL,
  `schedule_time` varchar(50) DEFAULT NULL,
  `schedule_days` varchar(50) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT 1,
  `last_run_at` datetime(3) DEFAULT NULL,
  `next_run_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cron_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `schedule_id` bigint unsigned DEFAULT NULL,
  `task_type` varchar(100) DEFAULT NULL,
  `status` varchar(50) DEFAULT NULL,
  `output` text DEFAULT NULL,
  `error` text DEFAULT NULL,
  `created_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_cron_logs_schedule_id` (`schedule_id`),
  KEY `idx_cron_logs_task_type` (`task_type`),
  CONSTRAINT `fk_cron_logs_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `cron_schedules` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `webhook_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `event` varchar(100) NOT NULL,
  `url` varchar(500) NOT NULL,
  `payload` longtext DEFAULT NULL,
  `response` longtext DEFAULT NULL,
  `status_code` int DEFAULT NULL,
  `duration` int DEFAULT NULL,
  `created_at` datetime(3) NULL,
  PRIMARY KEY (`id`),
  KEY `idx_webhook_logs_event` (`event`),
  KEY `idx_webhook_logs_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default admin user (password: admin123)
INSERT INTO `admin_users` (`username`, `password`, `email`, `role`, `status`, `created_at`, `updated_at`) 
VALUES ('admin', '$2a$10$8K1p/a0dLQl1d2Q.2JY.KfX3.3eB.1.6.3.3.3.3.3.3.3.3', 'admin@gembok.com', 'admin', 'active', NOW(), NOW());

-- Insert default settings
INSERT INTO `settings` (`setting_key`, `setting_value`, `description`, `updated_at`) VALUES
('MIKROTIK_HOST', '', 'MikroTik Router Host', NOW()),
('MIKROTIK_USER', '', 'MikroTik Router Username', NOW()),
('MIKROTIK_PASS', '', 'MikroTik Router Password', NOW()),
('MIKROTIK_PORT', '8728', 'MikroTik API Port', NOW()),
('GENIEACS_URL', '', 'GenieACS Server URL', NOW()),
('GENIEACS_USERNAME', '', 'GenieACS Username', NOW()),
('GENIEACS_PASSWORD', '', 'GenieACS Password', NOW()),
('TRIPAY_API_KEY', '', 'Tripay API Key', NOW()),
('TRIPAY_PRIVATE_KEY', '', 'Tripay Private Key', NOW()),
('TRIPAY_MERCHANT_CODE', '', 'Tripay Merchant Code', NOW()),
('TRIPAY_MODE', 'production', 'Tripay Mode (sandbox/production)', NOW()),
('DEFAULT_WHATSAPP_GATEWAY', 'fonnte', 'Default WhatsApp Gateway', NOW()),
('INVOICE_PREFIX', 'INV-', 'Invoice Number Prefix', NOW()),
('INVOICE_START', '1', 'Invoice Number Start', NOW()),
('CURRENCY_SYMBOL', 'Rp', 'Currency Symbol', NOW());
