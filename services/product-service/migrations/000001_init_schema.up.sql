-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    image_url VARCHAR(255),
    product_count INT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    type VARCHAR(20) DEFAULT 'normal',
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    original_price DECIMAL(10, 2),
    discount DECIMAL(5, 2) DEFAULT 0,
    image_url VARCHAR(255),
    category_id INTEGER REFERENCES categories(id),
    category_slug VARCHAR(100),
    stock_quantity INTEGER DEFAULT 0,
    brand VARCHAR(100),
    features JSONB,
    shipping_info JSONB,
    orders INTEGER DEFAULT 0,
    ui_metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create ads_placement table
CREATE TABLE IF NOT EXISTS ads_placement (
    id SERIAL PRIMARY KEY,
    location VARCHAR(50) NOT NULL, -- 'banner', 'display', 'gaming', 'new_fashion'
    reference_type VARCHAR(20) NOT NULL, -- 'product', 'category'
    reference_id INTEGER NOT NULL,
    display_order INTEGER DEFAULT 0,
    custom_title VARCHAR(255),
    custom_image_url VARCHAR(255),
    ui_settings JSONB, -- For special UI customizations
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create product_images table
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create product_tags table
CREATE TABLE IF NOT EXISTS product_tags (
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    tag VARCHAR(50) NOT NULL,
    PRIMARY KEY (product_id, tag)
);

-- Create product_reviews table
CREATE TABLE IF NOT EXISTS product_reviews (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    user_id VARCHAR(36) NOT NULL,
    user_name VARCHAR(100),
    rating DECIMAL(3, 1) NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create banners table
CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255),
    description TEXT,
    discount VARCHAR(100),
    highlight_text VARCHAR(100),
    image_url VARCHAR(255) NOT NULL,
    link_url VARCHAR(255),
    action_text VARCHAR(100) DEFAULT 'Shop Now',
    background_color VARCHAR(50) DEFAULT '#1e3a8a',
    text_color VARCHAR(50) DEFAULT '#ffffff',
    animation_type VARCHAR(20) DEFAULT 'fade',
    is_active BOOLEAN DEFAULT true,
    priority INTEGER DEFAULT 0,
    type VARCHAR(50) NOT NULL DEFAULT 'hero',
    product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for ads_placement table
CREATE INDEX IF NOT EXISTS idx_ads_placement_location ON ads_placement(location);
CREATE INDEX IF NOT EXISTS idx_ads_placement_reference ON ads_placement(reference_type, reference_id);
CREATE INDEX IF NOT EXISTS idx_ads_placement_order ON ads_placement(display_order);
CREATE INDEX IF NOT EXISTS idx_ads_placement_active ON ads_placement(is_active);

-- Indexes for banners table
CREATE INDEX IF NOT EXISTS idx_banners_type ON banners(type);
CREATE INDEX IF NOT EXISTS idx_banners_is_active ON banners(is_active);
CREATE INDEX IF NOT EXISTS idx_banners_priority ON banners(priority);
CREATE INDEX IF NOT EXISTS idx_banners_product_id ON banners(product_id);
CREATE INDEX IF NOT EXISTS idx_banners_category_id ON banners(category_id);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_category_slug ON products(category_slug);
CREATE INDEX IF NOT EXISTS idx_products_type ON products(type);
CREATE INDEX IF NOT EXISTS idx_products_brand ON products(brand);
CREATE INDEX IF NOT EXISTS idx_product_images_product_id ON product_images(product_id);
CREATE INDEX IF NOT EXISTS idx_product_reviews_product_id ON product_reviews(product_id);
CREATE INDEX IF NOT EXISTS idx_product_reviews_user_id ON product_reviews(user_id); 