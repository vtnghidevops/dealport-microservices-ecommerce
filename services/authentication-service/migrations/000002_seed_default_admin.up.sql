-- 000002_seed_default_admin.up.sql

-- Add default admin user
-- This migration creates a default admin user with a predefined password
-- The default password is 'admin123' and should be changed immediately after deployment
INSERT INTO users (id, email, password_hash, first_name, last_name, username, role, status, active)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'admin-sapogo@gmail.com', 
    '$2a$10$DgTbgv8SmctaZa/QIvWrouRCUvPrvoHcJAG1tpJ2R344c/IuG0bqC', 
    'Admin', 
    'User',
    'UserAdmin', -- username = lastName + firstName
    'admin',
    'active',
    true
);
