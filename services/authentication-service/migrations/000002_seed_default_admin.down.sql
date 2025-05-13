-- 000002_seed_default_admin.down.sql

-- Remove default admin user
DELETE FROM users WHERE username = 'admin';