-- Create database and tables if they don't exist
CREATE DATABASE IF NOT EXISTS mikrops_simple;
USE mikrops_simple;

-- Admin Users table
CREATE TABLE IF NOT EXISTS `admin_users` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `username` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `email` varchar(255) DEFAULT NULL,
    `role` varchar(50) DEFAULT 'admin' NOT NULL,
    `status` varchar(50) DEFAULT 'active' NOT NULL,
    `is_active` tinyint(1) DEFAULT 1 NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Insert default superadmin (password: admin123)
INSERT IGNORE INTO `admin_users` (`username`, `password`, `email`, `role`, `status`, `is_active`) 
VALUES ('superadmin', '$2a$10$vG8hBQbJ4aWvzLx4Z/2$2$10$u8JvG/9$7$9$9$Jf$9$1$v$8$7$I$J$8$8$8$8$8$I$I$I$I$I$I$I$I$I$Y', 'superadmin@mikrops.local', 'superadmin', 'active', 1),
       ('admin', '$2a$10$G8hBQbJ4aWvzLx4Z/2$2$10$u8JvG/9$7$9$9$Jf$9$1$v$8$7$I$J$8$8$8$8$8$I$I$I$I$I$I$I$I$I$Y', 'admin@mikrops.local', 'admin', 'active', 1);

-- Customers table
CREATE TABLE IF NOT EXISTS `customers` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `phone` varchar(255) NOT NULL,
    `email` varchar(255) DEFAULT NULL,
    `address` text,
    `package_id` int(11) unsigned DEFAULT NULL,
    `pppoe_username` varchar(255) DEFAULT NULL,
    `pppoe_password` varchar(255) DEFAULT NULL,
    `status` varchar(50) DEFAULT 'active' NOT NULL,
    `router_id` int(11) unsigned DEFAULT NULL,
    `onu_id` varchar(255) DEFAULT NULL,
    `onu_serial` varchar(255) DEFAULT NULL,
    `onu_mac_address` varchar(255) DEFAULT NULL,
    `onu_ip_address` varchar(255) DEFAULT NULL,
    `latitude` double DEFAULT NULL,
    `longitude` double DEFAULT NULL,
    `isolation_date` timestamp NULL DEFAULT NULL,
    `activation_date` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `phone` (`phone`),
    UNIQUE KEY `pppoe_username` (`pppoe_username`),
    KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Invoices table
CREATE TABLE IF NOT EXISTS `invoices` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `customer_id` int(11) unsigned NOT NULL,
    `number` varchar(255) NOT NULL,
    `amount` decimal(10,2) NOT NULL,
    `period` varchar(255) NOT NULL,
    `due_date` timestamp NOT NULL,
    `status` varchar(50) DEFAULT 'unpaid' NOT NULL,
    `paid_at` timestamp NULL DEFAULT NULL,
    `payment_method` varchar(255) DEFAULT NULL,
    `payment_reference` varchar(255) DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `number` (`number`),
    KEY `customer_id` (`customer_id`),
    KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Packages table
CREATE TABLE IF NOT EXISTS `packages` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `price` decimal(10,2) NOT NULL,
    `speed` varchar(255) NOT NULL,
    `description` text,
    `profile_normal` varchar(255) DEFAULT NULL,
    `profile_isolir` varchar(255) DEFAULT NULL,
    `status` varchar(50) DEFAULT 'active' NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Routers table
CREATE TABLE IF NOT EXISTS `routers` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `host` varchar(255) NOT NULL,
    `username` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `port` int(11) DEFAULT 8728 NOT NULL,
    `is_active` tinyint(1) DEFAULT 0 NOT NULL,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ONU Locations table
CREATE TABLE IF NOT EXISTS `onu_locations` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `customer_id` int(11) unsigned NOT NULL,
    `router_id` int(11) unsigned DEFAULT NULL,
    `onu_id` varchar(255) NOT NULL,
    `serial_number` varchar(255) DEFAULT NULL,
    `mac_address` varchar(255) DEFAULT NULL,
    `ip_address` varchar(255) DEFAULT NULL,
    `port_number` varchar(255) DEFAULT NULL,
    `signal_strength` varchar(255) DEFAULT NULL,
    `latitude` double DEFAULT NULL,
    `longitude` double DEFAULT NULL,
    `address` text,
    `notes` text,
    `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `customer_id` (`customer_id`),
    UNIQUE KEY `onu_id` (`onu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Trouble Tickets table
CREATE TABLE IF NOT EXISTS `trouble_tickets` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `customer_id` int(11) unsigned NOT NULL,
    `subject` varchar(255) NOT NULL,
    `description` text NOT NULL,
    `priority` varchar(50) DEFAULT 'medium' NOT NULL,
    `status` varchar(50) DEFAULT 'open' NOT NULL,
    `assigned_to` varchar(255) DEFAULT NULL,
    `resolved_at` timestamp NULL DEFAULT NULL,
    `created_at` timestamp NULL DEFAULT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `customer_id` (`customer_id`),
    KEY `status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Settings table
CREATE TABLE IF NOT EXISTS `settings` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `setting_key` varchar(255) NOT NULL,
    `setting_value` text,
    `description` text,
    `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `setting_key` (`setting_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- RBAC: Audit Log table
CREATE TABLE IF NOT EXISTS `audit_log` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `username` VARCHAR(255) NOT NULL,
    `action` VARCHAR(100) NOT NULL,
    `resource_type` VARCHAR(50) NOT NULL,
    `resource_id` VARCHAR(255),
    `details` JSON,
    `ip_address` VARCHAR(45),
    `user_agent` VARCHAR(500),
    `status` VARCHAR(20) NOT NULL DEFAULT 'success',
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_action` (`action`),
    INDEX `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Note: Casbin will create casbin_rule table automatically via GORM adapter
