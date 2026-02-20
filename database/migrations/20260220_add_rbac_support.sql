-- Add is_active field to admin_users
ALTER TABLE `admin_users` 
ADD COLUMN `is_active` BOOLEAN DEFAULT TRUE NOT NULL AFTER `status`;

-- Update existing admin records
UPDATE `admin_users` SET `is_active` = TRUE WHERE `is_active` IS NULL;

-- Ensure valid roles (set to 'admin' if empty)
UPDATE `admin_users` SET `role` = 'admin' 
WHERE `role` IS NULL OR `role` = '';

-- Create audit_log table
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

-- Note: Casbin's casbin_rule table will be created automatically by GORM adapter
