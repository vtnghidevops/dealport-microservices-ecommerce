-- Insert a sample user with password 'password123' (hashed)
INSERT INTO users (id, email, password_hash, first_name, last_name, username, role, status, active)
VALUES (
    'f47ac10b-58cc-4372-a567-0e02b2c3d479',
    'admin@example.com', 
    '$2a$10$tRNEG16z9o2kVQs0Cq.dZuXlKqLjNbhGKdZ/GG6nxYSxpZyK6.sPe', 
    'Admin', 
    'User',
    'UserAdmin', -- username = lastName + firstName
    'admin',
    'active',
    true
);

INSERT INTO users (id, email, password_hash, first_name, last_name, username, role, status, active)
VALUES (
    'a47ac10b-58cc-4372-a567-0e02b2c3d480',
    'user@example.com', 
    '$2a$10$tRNEG16z9o2kVQs0Cq.dZuXlKqLjNbhGKdZ/GG6nxYSxpZyK6.sPe', 
    'Regular', 
    'User',
    'UserRegular', -- username = lastName + firstName
    'user', -- Đổi từ 'customer' sang 'user'
    'active',
    true
);