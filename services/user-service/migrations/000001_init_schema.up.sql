CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table for profile information
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    display_name VARCHAR(100),
    username VARCHAR(100), -- Tự động tạo từ last_name + first_name
    phone VARCHAR(20),
    profile_image VARCHAR(255),
    role VARCHAR(20) NOT NULL DEFAULT 'user', -- Default role là 'user'
    status VARCHAR(20) DEFAULT 'active',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    date_of_birth DATE,
    gender VARCHAR(20),
    cart_id UUID,
    order_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Addresses table
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100), -- Tên người nhận
    phone VARCHAR(20), -- SĐT liên hệ
    line1 VARCHAR(255) NOT NULL, -- Địa chỉ
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    postal_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL DEFAULT 'Vietnam',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    address_type VARCHAR(20) DEFAULT 'shipping',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Payment methods
CREATE TABLE IF NOT EXISTS payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- 'credit_card', 'mono', 'other'
    provider VARCHAR(50),
    account_number VARCHAR(50),
    expiry_date VARCHAR(10),
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User wishlist items
CREATE TABLE IF NOT EXISTS user_wishlist_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id INT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    notes TEXT,
    UNIQUE (user_id, product_id)
);

-- Ensure only one default address per user and type
CREATE UNIQUE INDEX IF NOT EXISTS idx_addresses_user_default_type 
ON addresses (user_id, address_type) 
WHERE is_default = TRUE;

-- Ensure only one default payment method per user
CREATE UNIQUE INDEX IF NOT EXISTS idx_payment_methods_user_default
ON payment_methods (user_id)
WHERE is_default = TRUE; 