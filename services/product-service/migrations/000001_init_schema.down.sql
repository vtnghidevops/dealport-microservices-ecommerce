-- Drop indexes first
DROP INDEX IF EXISTS idx_product_reviews_user_id;
DROP INDEX IF EXISTS idx_product_reviews_product_id;
DROP INDEX IF EXISTS idx_product_images_product_id;
DROP INDEX IF EXISTS idx_products_brand;
DROP INDEX IF EXISTS idx_products_type;
DROP INDEX IF EXISTS idx_products_category_slug;
DROP INDEX IF EXISTS idx_products_category_id;

DROP INDEX IF EXISTS idx_banners_category_id;
DROP INDEX IF EXISTS idx_banners_product_id;
DROP INDEX IF EXISTS idx_banners_priority;
DROP INDEX IF EXISTS idx_banners_is_active;
DROP INDEX IF EXISTS idx_banners_type;

DROP INDEX IF EXISTS idx_ads_placement_active;
DROP INDEX IF EXISTS idx_ads_placement_order;
DROP INDEX IF EXISTS idx_ads_placement_reference;
DROP INDEX IF EXISTS idx_ads_placement_location;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS product_reviews;
DROP TABLE IF EXISTS product_tags;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS ads_placement;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories; 