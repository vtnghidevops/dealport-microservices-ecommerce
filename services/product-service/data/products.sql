-- products.sql
-- Database schema for products service

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    original_price DECIMAL(10, 2),
    discount DECIMAL(5, 2),
    image_url VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    category_slug VARCHAR(255) NOT NULL,
    stock_quantity INTEGER NOT NULL DEFAULT 0,
    brand VARCHAR(100),
    features JSONB,
    shipping_info JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Product images table
CREATE TABLE IF NOT EXISTS product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Product tags table
CREATE TABLE IF NOT EXISTS product_tags (
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    PRIMARY KEY (product_id, tag)
);

-- Product reviews table
CREATE TABLE IF NOT EXISTS product_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(product_id, user_id)
);

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    image_url VARCHAR(255),
    icon VARCHAR(100),
    banner_url VARCHAR(255),
    product_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_visible BOOLEAN DEFAULT true,
    display_order INTEGER DEFAULT 0,
    meta_title VARCHAR(255),
    meta_description TEXT,
    parent_id UUID REFERENCES categories(id),
    level INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Category attributes table
CREATE TABLE IF NOT EXISTS category_attributes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'number', 'boolean', 'select')),
    required BOOLEAN DEFAULT false,
    options TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_slug ON products(slug);
CREATE INDEX idx_products_brand ON products(brand);
CREATE INDEX idx_product_images_product_id ON product_images(product_id);
CREATE INDEX idx_product_reviews_product_id ON product_reviews(product_id);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_category_attributes_category_id ON category_attributes(category_id);

-- Insert some sample categories
INSERT INTO categories (name, slug, description, level, is_active, is_visible)
VALUES 
('Electronics', 'electronics', 'Electronic devices and gadgets', 0, true, true),
('Clothing', 'clothing', 'Apparel and fashion items', 0, true, true);

-- Insert subcategories
INSERT INTO categories (name, slug, description, parent_id, level, is_active, is_visible)
SELECT 'Smartphones', 'smartphones', 'Mobile phones and accessories', id, 1, true, true
FROM categories WHERE slug = 'electronics';

INSERT INTO categories (name, slug, description, parent_id, level, is_active, is_visible)
SELECT 'Laptops', 'laptops', 'Portable computers', id, 1, true, true
FROM categories WHERE slug = 'electronics';

INSERT INTO categories (name, slug, description, parent_id, level, is_active, is_visible)
SELECT 'Men', 'men', 'Men''s clothing', id, 1, true, true
FROM categories WHERE slug = 'clothing';

INSERT INTO categories (name, slug, description, parent_id, level, is_active, is_visible)
SELECT 'Women', 'women', 'Women''s clothing', id, 1, true, true
FROM categories WHERE slug = 'clothing';