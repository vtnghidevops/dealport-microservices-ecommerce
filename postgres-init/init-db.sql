-- Products database schema and initial data
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'postgres-products'
    ) THEN
        CREATE ROLE "postgres-products" LOGIN PASSWORD 'password';
    END IF;
END
$$;

-- Tạo DB nếu chưa có
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_database WHERE datname = 'products'
    ) THEN
        CREATE DATABASE products OWNER "postgres-products";
    END IF;
END
$$;

\connect products


-- Create tables for products database
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(50),
    price DECIMAL(10, 2) NOT NULL,
    original_price DECIMAL(10, 2),
    discount DECIMAL(10, 2),
    image_url VARCHAR(255),
    category_id VARCHAR(36) NOT NULL,
    category_slug VARCHAR(255) NOT NULL,
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    brand VARCHAR(100),
    features JSONB,
    shipping_info JSONB,
    highlighted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

GRANT ALL PRIVILEGES ON DATABASE products TO "postgres-products";


CREATE TABLE IF NOT EXISTS product_images (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_tags (
    product_id VARCHAR(36) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag VARCHAR(100) NOT NULL,
    PRIMARY KEY (product_id, tag)
);

CREATE TABLE IF NOT EXISTS categories (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    image_url VARCHAR(255),
    icon VARCHAR(100),
    banner_url VARCHAR(255),
    product_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_visible BOOLEAN DEFAULT TRUE, 
    display_order INTEGER DEFAULT 0,
    meta_title VARCHAR(255),
    meta_description TEXT,
    parent_id VARCHAR(36) REFERENCES categories(id),
    level INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS category_attributes (
    id VARCHAR(36) PRIMARY KEY,
    category_id VARCHAR(36) NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    required BOOLEAN DEFAULT FALSE,
    options JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS product_reviews (
    id VARCHAR(36) PRIMARY KEY, 
    product_id VARCHAR(36) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL,
    user_name VARCHAR(100),
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_category_slug ON products(category_slug);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_products_type ON products(type);
CREATE INDEX idx_product_reviews_product_id ON product_reviews(product_id);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);

-- Insert categories with UUIDs
INSERT INTO categories (id, name, slug, description, image_url, icon, level, is_active, is_visible) 
VALUES 
('be1fb352-cb58-4d47-b03b-814278b9abdc', 'Electronics', 'electronics', 'Electronic devices and gadgets', 'https://images.unsplash.com/photo-1550009158-9ebf69173e03', 'IoTvOutline', 0, true, true),
('c2d7e3a8-75f6-4409-ab92-5be6c2c9c0ef', 'Smartphones', 'smartphones', 'Latest smartphones', 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9', 'IoPhonePortraitOutline', 1, true, true),
('1f1cfc46-7a96-4190-941e-fd9d940c7eb1', 'Laptops', 'laptops', 'Powerful laptops for all needs', 'https://images.unsplash.com/photo-1602080858428-57174f9431cf', 'IoLaptopOutline', 1, true, true),
('a7e98f31-bd6d-4bea-9b23-34ca7f38917f', 'Clothing', 'clothing', 'Clothing and fashion items', 'https://images.unsplash.com/photo-1523381210434-271e8be1f52b', 'IoShirtOutline', 0, true, true),
('e8d4e326-86be-4196-9175-a0f7c9413b2c', 'Men', 'men', 'Men''s fashion', 'https://images.unsplash.com/photo-1620012253295-c15cc3e65df4', 'IoManOutline', 1, true, true),
('b78e3c0d-6c9e-4cc7-85a3-2e91d3833943', 'Women', 'women', 'Women''s fashion', 'https://images.unsplash.com/photo-1618244972960-6db44be685ad', 'IoWomanOutline', 1, true, true),
('9c6ad65d-8715-45d6-ba12-f1c20fdc6dd1', 'Home & Garden', 'home-garden', 'Home decor and garden supplies', 'https://images.unsplash.com/photo-1581783342308-f792dbdd27c5', 'IoHomeOutline', 0, true, true);

-- Update parent IDs for subcategories
UPDATE categories SET parent_id = 'be1fb352-cb58-4d47-b03b-814278b9abdc' WHERE id IN ('c2d7e3a8-75f6-4409-ab92-5be6c2c9c0ef', '1f1cfc46-7a96-4190-941e-fd9d940c7eb1');
UPDATE categories SET parent_id = 'a7e98f31-bd6d-4bea-9b23-34ca7f38917f' WHERE id IN ('e8d4e326-86be-4196-9175-a0f7c9413b2c', 'b78e3c0d-6c9e-4cc7-85a3-2e91d3833943');

-- Add more detailed product data with UUIDs
INSERT INTO products (
    id, 
    name, 
    slug, 
    description, 
    price, 
    original_price, 
    discount,
    image_url, 
    category_id, 
    category_slug, 
    stock_quantity, 
    type,
    brand,
    features,
    shipping_info,
    highlighted
) VALUES
-- Electronics - Smartphones
('7d1e7337-2722-4a20-a685-3e53eaa51669', 'iPhone 14 Pro', 'iphone-14-pro', 
'Experience the next level of iPhone with dynamic island, advanced camera system, and A16 Bionic chip.', 
999.99, 1099.99, 100.00, 
'https://images.unsplash.com/photo-1551641506-ee5bf4cb45f3', 
'c2d7e3a8-75f6-4409-ab92-5be6c2c9c0ef', 'smartphones', 50, 'trending', 'Apple',
'{"display": "6.1-inch Super Retina XDR", "chip": "A16 Bionic", "camera": "48MP Main | Ultra Wide | Telephoto", "battery": "Up to 23 hours video playback"}',
'{"express": true, "standard": true, "international": true}',
true),

('f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'Samsung Galaxy S23 Ultra', 'samsung-galaxy-s23-ultra', 
'The ultimate smartphone experience with S Pen, 200MP camera, and powerful performance.', 
1199.99, 1299.99, 100.00, 
'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c', 
'c2d7e3a8-75f6-4409-ab92-5be6c2c9c0ef', 'smartphones', 30, 'new', 'Samsung',
'{"display": "6.8-inch Dynamic AMOLED 2X", "chip": "Snapdragon 8 Gen 2", "camera": "200MP Main | Ultra Wide | 2x Telephoto", "battery": "5000mAh"}',
'{"express": true, "standard": true, "international": true}',
true),

-- Electronics - Laptops
('c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'MacBook Pro M2', 'macbook-pro-m2', 
'Supercharged by the next-generation M2 chip, the redesigned MacBook Pro is the ultimate pro laptop.', 
1499.99, 1599.99, 100.00, 
'https://images.unsplash.com/photo-1517336714731-489689fd1ca8', 
'1f1cfc46-7a96-4190-941e-fd9d940c7eb1', 'laptops', 20, 'top-sale', 'Apple',
'{"processor": "Apple M2", "memory": "16GB unified memory", "storage": "512GB SSD", "display": "14-inch Liquid Retina XDR"}',
'{"express": true, "standard": true, "international": false}',
false),

('825d29e1-c960-4237-86f6-f6c96bcd8af3', 'Dell XPS 15', 'dell-xps-15', 
'A powerhouse laptop with stunning InfinityEdge display and high-performance components.', 
1899.99, 2099.99, 200.00, 
'https://images.unsplash.com/photo-1593642702749-b7d2a804fbcf', 
'1f1cfc46-7a96-4190-941e-fd9d940c7eb1', 'laptops', 15, 'normal', 'Dell',
'{"processor": "Intel Core i9-12900HK", "memory": "32GB DDR5", "storage": "1TB NVMe SSD", "display": "15.6-inch 4K OLED"}',
'{"express": true, "standard": true, "international": true}',
false),

-- Clothing - Men
('a0e5b6a1-86f4-4380-8d38-a10e1786916f', 'Premium Cotton T-Shirt', 'premium-cotton-tshirt', 
'Soft, breathable cotton t-shirt perfect for everyday wear.', 
29.99, 39.99, 10.00, 
'https://images.unsplash.com/photo-1581655353564-df123a1eb820', 
'e8d4e326-86be-4196-9175-a0f7c9413b2c', 'men', 100, 'normal', 'Essentials',
'{"material": "100% Organic Cotton", "fit": "Regular", "care": "Machine wash cold", "origin": "Made in Vietnam"}',
'{"express": true, "standard": true, "international": true}',
false),

-- Clothing - Women
('9d9ed1e8-c2f8-40d8-b8fc-6c0971a048b3', 'Summer Floral Dress', 'summer-floral-dress', 
'Elegant floral pattern dress, perfect for summer events.', 
59.99, 79.99, 20.00, 
'https://images.unsplash.com/photo-1492707892479-7bc8d5a4ee93', 
'b78e3c0d-6c9e-4cc7-85a3-2e91d3833943', 'women', 75, 'trending', 'Urban Style',
'{"material": "Viscose", "length": "Midi", "pattern": "Floral print", "style": "V-neck with short sleeves"}',
'{"express": true, "standard": true, "international": true}',
true),

-- Home & Garden
('b1e67fe9-93bb-4df0-8743-e6d28648396d', 'Smart Home Speaker', 'smart-home-speaker', 
'Voice-controlled smart speaker with advanced AI assistant.', 
129.99, 149.99, 20.00, 
'https://images.unsplash.com/photo-1543512214-318c7553f230', 
'9c6ad65d-8715-45d6-ba12-f1c20fdc6dd1', 'home-garden', 40, 'new', 'EchoTech',
'{"connectivity": "WiFi, Bluetooth", "assistant": "Multi-platform compatible", "audio": "360° premium sound", "microphones": "Far-field voice recognition"}',
'{"express": true, "standard": true, "international": false}',
true);

-- Add product images
INSERT INTO product_images (id, product_id, url, is_primary, display_order) VALUES
-- iPhone images
('d4a9e8c7-2d20-4fc4-87cb-c32101b11d1c', '7d1e7337-2722-4a20-a685-3e53eaa51669', 'https://images.unsplash.com/photo-1551641506-ee5bf4cb45f3', true, 0),
('e8a9d2b6-3f17-4e0a-bc9d-7743a8f8c31a', '7d1e7337-2722-4a20-a685-3e53eaa51669', 'https://images.unsplash.com/photo-1592899677977-9c10ca588bbd', false, 1),
('f7c2e5d1-8a46-42b9-ae36-9f51a9b832c0', '7d1e7337-2722-4a20-a685-3e53eaa51669', 'https://images.unsplash.com/photo-1592286927505-1def25115558', false, 2),

-- Samsung images
('a1c3e5b7-9f28-40d6-8e19-2b7a3c641d0f', 'f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'https://images.unsplash.com/photo-1610945415295-d9bbf067e59c', true, 0),
('b2d4f6c8-0e39-41e7-9f20-3c8b4d752e1g', 'f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'https://images.unsplash.com/photo-1565849904461-04a58ad377e0', false, 1),
('c3e5g7d9-1f40-42f8-0g21-4d9c5e863f2h', 'f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'https://images.unsplash.com/photo-1598327105666-5b89351aff97', false, 2),

-- MacBook images
('h6j8k0l2-5g34-49h6-7j38-9m10n11o12p', 'c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8', true, 0),
('i7j9k1l3-6h45-50i7-8k49-0n11o12p13q', 'c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'https://images.unsplash.com/photo-1537498425277-c283d32ef9db', false, 1),
('j8k0l2m4-7i56-61j8-9l50-1o12p13q14r', 'c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'https://images.unsplash.com/photo-1611186871348-b1ce696e52c9', false, 2);

-- Add product tags
INSERT INTO product_tags (product_id, tag) VALUES
('7d1e7337-2722-4a20-a685-3e53eaa51669', 'phone'),
('7d1e7337-2722-4a20-a685-3e53eaa51669', 'premium'),
('7d1e7337-2722-4a20-a685-3e53eaa51669', 'apple'),
('7d1e7337-2722-4a20-a685-3e53eaa51669', 'ios'),

('f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'phone'),
('f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'premium'),
('f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'samsung'),
('f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'android'),

('c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'laptop'),
('c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'premium'),
('c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'apple'),
('c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'macos'),

('825d29e1-c960-4237-86f6-f6c96bcd8af3', 'laptop'),
('825d29e1-c960-4237-86f6-f6c96bcd8af3', 'premium'),
('825d29e1-c960-4237-86f6-f6c96bcd8af3', 'dell'),
('825d29e1-c960-4237-86f6-f6c96bcd8af3', 'windows'),

('a0e5b6a1-86f4-4380-8d38-a10e1786916f', 'clothing'),
('a0e5b6a1-86f4-4380-8d38-a10e1786916f', 'casual'),
('a0e5b6a1-86f4-4380-8d38-a10e1786916f', 'men'),
('a0e5b6a1-86f4-4380-8d38-a10e1786916f', 'cotton'),

('9d9ed1e8-c2f8-40d8-b8fc-6c0971a048b3', 'clothing'),
('9d9ed1e8-c2f8-40d8-b8fc-6c0971a048b3', 'dress'),
('9d9ed1e8-c2f8-40d8-b8fc-6c0971a048b3', 'women'),
('9d9ed1e8-c2f8-40d8-b8fc-6c0971a048b3', 'summer'),

('b1e67fe9-93bb-4df0-8743-e6d28648396d', 'smart-home'),
('b1e67fe9-93bb-4df0-8743-e6d28648396d', 'speaker'),
('b1e67fe9-93bb-4df0-8743-e6d28648396d', 'tech'),
('b1e67fe9-93bb-4df0-8743-e6d28648396d', 'voice-control');

-- Add product reviews
INSERT INTO product_reviews (id, product_id, user_id, user_name, rating, comment, created_at) VALUES
-- iPhone reviews
('r1c3e5g7-9i0k-2m4o-6q8s-0u2w4y6a8c0', '7d1e7337-2722-4a20-a685-3e53eaa51669', 'u123456', 'John Smith', 5, 'Best iPhone I''ve ever had! The camera quality is amazing.', CURRENT_TIMESTAMP - INTERVAL '3 days'),
('r2d4f6h8-0j1l-3n5p-7r9t-1v3x5z7b9d1', '7d1e7337-2722-4a20-a685-3e53eaa51669', 'u234567', 'Emily Johnson', 4, 'Great phone overall, but battery life could be better.', CURRENT_TIMESTAMP - INTERVAL '7 days'),
('r3e5g7i9-1k2m-4o6q-8s0u-2w4y6a8c0e2', '7d1e7337-2722-4a20-a685-3e53eaa51669', 'u345678', 'Michael Brown', 5, 'The dynamic island is a game changer! Love this phone.', CURRENT_TIMESTAMP - INTERVAL '10 days'),

-- Samsung reviews
('r4f6h8j0-2l3n-5p7r-9t1v-3x5z7b9d1f3', 'f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'u456789', 'Sarah Wilson', 5, 'The camera on this phone is incredible! The 200MP sensor captures amazing detail.', CURRENT_TIMESTAMP - INTERVAL '5 days'),
('r5g7i9k1-3m4o-6q8s-0u2w-4y6a8c0e2g4', 'f2186f89-b646-453c-bdb9-30f7e68a4e9b', 'u567890', 'David Garcia', 4, 'Really enjoying the S Pen functionality, makes this phone unique.', CURRENT_TIMESTAMP - INTERVAL '12 days'),

-- MacBook reviews
('r6h8j0l2-4n5p-7r9t-1v3x-5z7b9d1f3h5', 'c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'u678901', 'Jessica Martinez', 5, 'The M2 chip is blazing fast. Battery life is incredible too!', CURRENT_TIMESTAMP - INTERVAL '8 days'),
('r7i9k1m3-5o6q-8s0u-2w4y-6a8c0e2g4i6', 'c55c93b4-648a-4729-ba78-bf2e913a6b1f', 'u789012', 'Robert Taylor', 4, 'Great laptop for development work. The screen is beautiful.', CURRENT_TIMESTAMP - INTERVAL '15 days'),

-- Update product counts for categories
UPDATE categories c
SET product_count = (
    SELECT COUNT(*) FROM products p
    WHERE p.category_id = c.id
);

-- Update parent category product counts to include subcategories
UPDATE categories c
SET product_count = (
    SELECT COUNT(*) FROM products p
    WHERE p.category_id IN (
        SELECT id FROM categories
        WHERE id = c.id OR parent_id = c.id
    )
)
WHERE c.parent_id IS NULL; 