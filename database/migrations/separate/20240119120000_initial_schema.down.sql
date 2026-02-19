-- Migration: Initial Schema - Down

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS `webhook_logs`;
DROP TABLE IF EXISTS `cron_logs`;
DROP TABLE IF EXISTS `cron_schedules`;
DROP TABLE IF EXISTS `settings`;
DROP TABLE IF EXISTS `trouble_tickets`;
DROP TABLE IF EXISTS `onu_locations`;
DROP TABLE IF EXISTS `invoices`;
DROP TABLE IF EXISTS `customers`;
DROP TABLE IF EXISTS `routers`;
DROP TABLE IF EXISTS `packages`;
DROP TABLE IF EXISTS `admin_users`;
