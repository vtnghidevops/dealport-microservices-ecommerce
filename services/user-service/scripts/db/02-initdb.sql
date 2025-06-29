-- Insert a sample admin user matching the authentication service
INSERT INTO users (id, email, first_name, last_name, display_name, username, role, status, active)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'admin@example.com',
    'Admin',
    'User',
    'Admin User',
    'UserAdmin', -- username = lastName + firstName
    'admin',
    'active',
    true
);

-- Insert a sample regular user
INSERT INTO users (id, email, first_name, last_name, display_name, username, role, status, active)
VALUES (
    'a47ac10b-58cc-4372-a567-0e02b2c3d480',
    'user@example.com',
    'Regular',
    'User',
    'Regular User',
    'UserRegular', -- username = lastName + firstName
    'user',
    'active',
    true
);

-- Add sample address for admin
INSERT INTO addresses (user_id, name, phone, line1, city, state, postal_code, country, is_default, address_type)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'Admin User', -- Tên người nhận
    '0123456789', -- Số điện thoại liên hệ
    '123 Admin Street',
    'Admin City',
    'Admin State',
    '10000',
    'Vietnam',
    true,
    'shipping'
);

-- Add sample address for regular user
INSERT INTO addresses (user_id, name, phone, line1, city, state, postal_code, country, is_default, address_type)
VALUES (
    'a47ac10b-58cc-4372-a567-0e02b2c3d480',
    'Regular User', -- Tên người nhận
    '0987654321', -- Số điện thoại liên hệ
    '456 User Street',
    'User City',
    'User State',
    '20000',
    'Vietnam',
    true,
    'shipping'
);

-- Add payment methods for admin
INSERT INTO payment_methods (user_id, type, provider, account_number, expiry_date, is_default)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'credit_card',
    'Visa',
    '4111111111111111',
    '12/2025',
    true
);

-- Add payment methods for regular user
INSERT INTO payment_methods (user_id, type, provider, account_number, expiry_date, is_default)
VALUES (
    'a47ac10b-58cc-4372-a567-0e02b2c3d480',
    'credit_card',
    'Mastercard',
    '5555555555554444',
    '10/2024',
    true
);