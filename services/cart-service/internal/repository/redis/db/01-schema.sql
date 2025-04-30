-- Redis doesn't use traditional SQL schemas, but we can use this file for documentation purposes
-- This file serves as a reference for the key structures and patterns used in the Redis database

-- Cart keys are stored with the following pattern:
-- cart:{user_id} - Hash containing cart metadata
--   Fields:
--     id - The cart ID
--     user_id - The user ID
--     coupon_code - The applied coupon code (if any)
--     discount_amount - The discount amount from the coupon
--     created_at - Timestamp when the cart was created
--     updated_at - Timestamp when the cart was last updated

-- Cart items are stored with the following pattern:
-- cart:{user_id}:items - Set containing item IDs
-- cart:{user_id}:item:{item_id} - Hash containing item data
--   Fields:
--     id - The item ID
--     product_id - The product ID
--     name - The product name
--     price - The current price
--     original_price - The original price (before any discounts)
--     quantity - The quantity of this item
--     image_url - The URL to the product image

-- Cart totals are stored with the following pattern:
-- cart:{user_id}:totals - Hash containing cart totals
--   Fields:
--     subtotal - The subtotal (sum of all items)
--     shipping - The shipping cost
--     discount - The discount amount
--     tax - The tax amount
--     total - The total amount (subtotal + shipping + tax - discount)

-- In Redis, we'll primarily be using the following commands:
-- HSET - Set hash fields
-- HGET, HGETALL - Get hash fields
-- SADD, SMEMBERS - Manage sets
-- DEL - Delete keys
-- EXPIRE - Set key expiration 