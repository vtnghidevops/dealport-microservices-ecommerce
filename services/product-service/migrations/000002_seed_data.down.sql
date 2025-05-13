-- Delete all sample data
DELETE FROM products WHERE id BETWEEN 1 AND 5;
DELETE FROM categories WHERE id BETWEEN 1 AND 5;

-- Reset sequences
ALTER SEQUENCE products_id_seq RESTART WITH 1;
ALTER SEQUENCE categories_id_seq RESTART WITH 1; 