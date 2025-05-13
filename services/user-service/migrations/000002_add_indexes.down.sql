-- Drop indexes in reverse order
DROP INDEX IF EXISTS idx_payment_methods_user_id;
DROP INDEX IF EXISTS idx_wishlist_product_id;
DROP INDEX IF EXISTS idx_wishlist_user_id;
DROP INDEX IF EXISTS idx_addresses_user_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_payment_methods_user_default;
DROP INDEX IF EXISTS idx_addresses_user_default_type; 