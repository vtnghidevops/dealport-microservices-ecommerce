-- Rollback migration to remove category product count triggers

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_product_insert_update_category_count ON products;
DROP TRIGGER IF EXISTS trigger_product_update_update_category_count ON products;
DROP TRIGGER IF EXISTS trigger_product_delete_update_category_count ON products;

-- Drop function
DROP FUNCTION IF EXISTS update_category_product_count(); 