-- Migration to add triggers for automatic category product count updates

-- Function to update category product count
CREATE OR REPLACE FUNCTION update_category_product_count()
RETURNS TRIGGER AS $$
BEGIN
    -- For INSERT: increment count for new category
    IF TG_OP = 'INSERT' THEN
        UPDATE categories 
        SET product_count = (
            SELECT COUNT(*) 
            FROM products 
            WHERE category_id = NEW.category_id
        )
        WHERE id = NEW.category_id;
        RETURN NEW;
    END IF;
    
    -- For UPDATE: update counts for both old and new categories (if different)
    IF TG_OP = 'UPDATE' THEN
        -- Update old category count if category changed
        IF OLD.category_id != NEW.category_id THEN
            UPDATE categories 
            SET product_count = (
                SELECT COUNT(*) 
                FROM products 
                WHERE category_id = OLD.category_id
            )
            WHERE id = OLD.category_id;
        END IF;
        
        -- Update new category count
        UPDATE categories 
        SET product_count = (
            SELECT COUNT(*) 
            FROM products 
            WHERE category_id = NEW.category_id
        )
        WHERE id = NEW.category_id;
        RETURN NEW;
    END IF;
    
    -- For DELETE: decrement count for old category
    IF TG_OP = 'DELETE' THEN
        UPDATE categories 
        SET product_count = (
            SELECT COUNT(*) 
            FROM products 
            WHERE category_id = OLD.category_id
        )
        WHERE id = OLD.category_id;
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for product insert
CREATE TRIGGER trigger_product_insert_update_category_count
    AFTER INSERT ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_category_product_count();

-- Create trigger for product update
CREATE TRIGGER trigger_product_update_update_category_count
    AFTER UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_category_product_count();

-- Create trigger for product delete
CREATE TRIGGER trigger_product_delete_update_category_count
    AFTER DELETE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_category_product_count();

-- Initialize correct product counts for existing categories
UPDATE categories 
SET product_count = (
    SELECT COUNT(*) 
    FROM products 
    WHERE products.category_id = categories.id
);

-- Add comment for documentation
COMMENT ON FUNCTION update_category_product_count() IS 'Automatically updates category product counts when products are inserted, updated, or deleted'; 