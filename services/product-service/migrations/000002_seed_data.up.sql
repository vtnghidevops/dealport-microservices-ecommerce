-- Insert categories with UUIDs
INSERT INTO categories (id, name, slug, description, image_url, product_count, is_active, is_visible) 
VALUES 
(1, 'Grocery', 'grocery', 'Grocery', '/images/category/grocery.png', 0, true, true),
(2, 'Home', 'home', 'Home', '/images/category/home.png', 0, true, true),
(3, 'Fashion', 'fashion', 'Fashion', '/images/category/fashion.png', 0, true, true),
(4, 'Electronics', 'electronics', 'Electronics', '/images/category/electronics.png', 0, true, true),
(5, 'Toys', 'toys', 'Toys', '/images/category/toys.png', 0, true, true),
(6, 'Laptops', 'laptops', 'Laptops', 'https://images.unsplash.com/photo-1517336714731-489689fd1ca4', 0, true, true),
(7, 'Men', 'men', 'Men', 'https://images.unsplash.com/photo-1617137968427-85924c800a22', 0, true, true), 
(8, 'Limited', 'limited', 'Limited', 'https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da', 0, true, true),
(9, 'Top Sell', 'top-sell', 'Top Sell', 'https://images.unsplash.com/photo-1607082349566-187342175e2f', 0, true, true),
(10, 'Women', 'women', 'Women', 'https://images.unsplash.com/photo-1483985988355-763728e1935b', 0, true, true),
(11, 'Kids', 'kids', 'Kids', 'https://images.unsplash.com/photo-1519340241574-2cec6aef0c01', 0, true, true),
(12, 'Piano', 'piano', 'Piano', 'https://images.unsplash.com/photo-1520523839897-bd0b52f945a0', 0, true, true);

-- Insert products
INSERT INTO products (
  id, type, name, description, slug, price, original_price, discount,
  category_id, category_slug, stock_quantity, brand,
  features, shipping_info, ui_metadata, orders
) VALUES
(
    1, 'normal', '2020 Apple MacBook Pro with Apple M1 Chip (13-inch, 8GB RAM, 256GB SSD Storage) - Space Gray', 'The MacBook Air 13-inch is a slim and lightweight laptop from Apple with a sharp Retina display and impressive performance powered by the M1 chip. It''s known for its portability, sleek design, and reliable performance.', 'apple-macbook-pro-m1-chip-13-inch-8gb-ram-256gb-ssd-storage-space-gray', 100.25, 125.99, 25,
    1, 'grocery', 5, 'Apple',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    2, 'normal', 'Organic Bananas from sustainable farms', 'Fresh organic bananas from sustainable farms', 'organic-bananas', 2.99, 3.99, 25,
    1, 'grocery', 19, 'Apple',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    3, 'normal', 'Organic Bananas from sustainable farms', 'Fresh organic bananas from sustainable farms', 'organic-bananas1', 2.99, 3.99, 25,
    1, 'grocery', 0, 'Samsung',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    4, 'normal', 'Organic Bananas from sustainable farms', 'Fresh organic bananas from sustainable farms', 'organic-bananas2', 2.99, 3.99, 25,
    1, 'grocery', 150, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    5, 'normal', 'Whole Grain Bread', 'Freshly baked whole grain bread', 'whole-grain-bread', 3.49, 3.99, 25,
    1, 'grocery', 75, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    6, 'normal', 'Free Range Eggs (12pk)', 'Farm fresh free-range eggs', 'free-range-eggs', 4.99, 3.99, 25,
    1, 'grocery', 120, 'xiaomi',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    7, 'normal', 'Organic Milk', 'Fresh organic whole milk', 'organic-milk', 3.99, 3.99, 25,
    1, 'grocery', 95, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    8, 'normal', 'Avocados (3pk)', 'Ripe and ready to eat avocados', 'avocados1', 5.99, 3.99, 25,
    1, 'grocery', 65, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    9, 'normal', 'Avocados (3pk)', 'Ripe and ready to eat avocados', 'avocados2', 5.99, 3.99, 25,
    1, 'grocery', 65, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    10, 'normal', 'Scented Candle Set', 'Set of 3 aromatic scented candles', 'scented-candle-set', 24.99, 3.99, 25,
    1, 'home', 45, 'xiaomi',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    12, 'normal', 'Cotton Bed Sheets', 'Soft 100% Egyptian cotton bed sheets', 'cotton-bed-sheets', 49.99, 3.99, 25,
    2, 'home', 30, 'xiaomi',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    13, 'normal', 'Kitchen Utensil Set', 'Complete set of silicone kitchen utensils', 'kitchen-utensil-set1', 35.99, 3.99, 25,
    2, 'home', 55, 'apple',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    14, 'normal', 'Decorative Throw Pillows', 'Set of 2 decorative throw pillows for sofa or bed', 'decorative-throw-pillows1', 22.99, 3.99, 25,
    2, 'home', 0, 'panasonic',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    15, 'normal', 'LED Table Lamp', 'Modern LED table lamp with adjustable brightness', 'led-table-lamp', 39.99, 3.99, 25,
    2, 'home', 25, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    16, 'normal', 'Women''s Denim Jacket', 'Classic denim jacket for women', 'womens-denim-jacket', 59.99, 3.99, 25,
    3, 'fashion', 40, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    17, 'normal', 'Men''s Casual T-Shirt', '100% cotton casual t-shirt for men', 'mens-casual-tshirt', 19.99, 3.99, 25,
    3, 'fashion', 120, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    18, 'normal', 'Leather Sneakers', 'Comfortable leather sneakers for everyday use', 'leather-sneakers', 89.99, 3.99, 25,
    3, 'fashion', 35, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    19, 'normal', 'Summer Dress', 'Lightweight cotton summer dress', 'summer-dress', 45.99, 3.99, 25,
    3, 'fashion', 0, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    20, 'normal', 'Unisex Beanie', 'Warm knitted beanie for all seasons', 'unisex-beanie', 14.99, 3.99, 25,
    3, 'fashion', 60, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    21, 'normal', 'Wireless Earbuds', 'Bluetooth wireless earbuds with noise cancellation', 'wireless-earbuds', 129.99, 3.99, 25,
    4, 'electronic', 45, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    22, 'normal', '4K Smart TV', '55-inch 4K Ultra HD Smart LED TV', '4k-smart-tv', 699.99, 3.99, 25,
    4, 'electronic', 15, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    23, 'normal', 'Portable Bluetooth Speaker', 'Waterproof portable bluetooth speaker with 20-hour battery life', 'portable-bluetooth-speaker', 89.99, 3.99, 25,
    4, 'electronic', 30, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    24, 'normal', 'Digital Camera', '24MP digital camera with 4K video recording', 'digital-camera', 449.99, 3.99, 25,
    4, 'electronic', 0, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    25, 'normal', 'Smartwatch', 'Fitness tracking smartwatch with heart rate monitor', 'smartwatch', 199.99, 3.99, 25,
    4, 'electronic', 25, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    26, 'limited', 'Building Blocks Set', 'Educational building blocks set for children', 'building-blocks-set', 29.99, 3.99, 25,
    5, 'toys', 50, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    27, 'trending', 'Remote Control Car', 'High-speed remote control race car', 'remote-control-car', 49.99, 3.99, 25,
    5, 'toys', 35, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 10
  ),
(
    28, 'trending', 'Stuffed Teddy Bear', 'Soft and huggable teddy bear for kids', 'stuffed-teddy-bear', 19.99, 3.99, 25,
    5, 'toys', 80, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    29, 'trending', 'Educational Board Game', 'Family-friendly educational board game', 'educational-board-game', 34.99, 3.99, 25,
    5, 'toys', 0, 'Educational Games Co.',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    30, 'normal', 'Art and Craft Kit', 'Creative art and craft kit for children', 'art-and-craft-kit', 24.99, 3.99, 25,
    5, 'toys', 40, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    31, 'trending', 'Instant Coffee', 'Premium instant coffee for quick preparation', 'instant-coffee', 7.99, 3.99, 25,
    6, 'laptop', 70, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    32, 'normal', 'Chocolate Cookies', 'Crunchy chocolate chip cookies', 'chocolate-cookies', 3.99, 3.99, 25,
    6, 'laptop', 90, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    33, 'normal', 'Fresh Orange Juice', '100% freshly squeezed orange juice', 'fresh-orange-juice', 4.49, 3.99, 25,
    6, 'laptop', 55, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    34, 'trending', 'Radiant Glow Hydrating Serum', 'Gentle yet effective, our Radiant Boosting Foaming our Radiant Boosting Foaming our Radiant Boosting Foaming', 'radiant-glow-hydrating-serum', 29.99, 39.99, 20,
    1, 'grocery', 5, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    35, 'trending', 'Modern Minimalist Vase', 'Track your workouts, heart rate, sleep quality and receive notifications. Water resistant up to 50m with 7-day battery life.', 'modern-minimalist-vase', 40.99, 49.99, 18,
    1, 'grocery', 5, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    36, 'trending', 'FitPro 3000 Smart Watch', 'Fast-charging power bank with dual USB ports and USB-C compatibility. Charge multiple devices simultaneously on the go.', 'fitpro-3000-smart-watch', 119.99, 0, 0,
    1, 'grocery', 5, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    37, 'trending', 'Pants', 'Pants', 'pants', 25.95, 45.95, 20,
    7, 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    38, 'trending', 'Shirt', 'Shirt', 'shirt', 12.99, 12.99, 0,
    7, 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    39, 'trending', 'Hat', 'Hat', 'hat', 104, 104, 0,
    '7', 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    40, 'trending', 'Shoe', 'Shoe', 'shoe', 10.56, 10.56, 0,
    '7', 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    41, 'limited', 'Samsung Galaxy S24', 'Gentle yet effective, our hydrating serum infuses skin with essential moisture while brightening and plumping for a radiant complexion.', 'samsung-galaxy-s24', 29.99, 39.99, 25,
    8, 'limited', 0, 'samsung',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    42, 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud', 34.99, 42.99, 18,
    8, 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    43, 'limited', 'Winter fashion jacket', 'Delivers extreme hydration and helps strengthen the skin barrier. Perfect for both daytime wear and overnight rejuvenation. [[2]]', 'winter-fashion-jacket', 39.99, 49.99, 20,
    8, 'limited', 0, 'piano',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    44, 'limited', 'New Balance 574 Senekers', 'Illuminating serum infused with glow-boosting intelligent botanicals that leave the skin visibly radiant while reducing redness. [[4]]', 'new-balance-574-senekers', 44.99, 59.99, 25,
    8, 'limited', 0, 'new balance',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    45, 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud1', 34.99, 42.99, 18,
    8, 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    46, 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud2', 34.99, 42.99, 18,
    8, 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    47, 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud3', 34.99, 42.99, 18,
    8, 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    48, 'top-sale', 'Computer Accessories', 'category', 'computer-accessories', 12.99, 12.99, 0,
    10, 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isComingSoon":true}', 100
  ),
(
    49, 'top-sale', 'Men''s Casual Outfit', 'category', 'men-casual-outfit', 200, 200, 0,
    10, 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isComingSoon":false}', 100
  ),
(
    50, 'top-sale', 'Pome Granate Juice', 'category', 'pome-granate-juice', 49, 49, 0,
    10, 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isComingSoon":false}', 100
  ),
(
    51, 'top-sale', 'Dog Food Made With Love', 'category', 'dog-food-made-with-love', 12.99, 12.99, 0,
    10, 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"col","isComingSoon":false}', 100
  ),
(
    52, 'top-sale', 'Security Camera System', 'category', 'security-camera-system', 420, 420, 0,
    10, 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"double","isComingSoon":false}', 100
),
(
    53, 'top-sale', 'Premium Cosmetic Set', 'category', 'premium-cosmetic-set', 149, 149, 0,
    10, 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isComingSoon":false}', 100
);



-- Add product images
INSERT INTO product_images (id, product_id, url, is_primary, display_order) VALUES
(1, 1, '/images/products/grocery/bananas.png', true, 0),
(2, 1, '/images/products/grocery/02.png', false, 1),
(3, 1, '/images/products/grocery/laptop.png', false, 2),
(4, 1, '/images/products/grocery/03.png', false, 3),
(5, 1, '/images/products/grocery/04.png', false, 4),
(6, 1, '/images/products/grocery/05.png', false, 5),
(7, 1, '/images/products/grocery/06.png', false, 6),
(8, 1, '/images/products/grocery/bananas.png', true, 7),
(9, 2, '/images/products/grocery/bananas.png', false, 0),
(10, 2, '/images/products/grocery/laptop.png', true, 1),
(11, 2, '/images/products/grocery/bananas.png', false, 2),
(12, 2, '/images/products/grocery/bananas.png', false, 3),
(13, 2, '/images/products/grocery/bananas.png', false, 4),
(14, 3, '/images/products/grocery/bananas.png', false, 0),
(15, 3, '/images/products/grocery/bananas.png', false, 1),
(16, 3, '/images/products/grocery/bananas.png', false, 2),
(17, 3, '/images/products/grocery/bananas.png', false, 3),
(18, 3, '/images/products/grocery/bananas.png', false, 4),
(19, 3, '/images/products/grocery/laptop.png', true, 5),
(20, 4, '/images/products/grocery/bananas.png', true, 0),
(21, 4, '/images/products/grocery/bananas.png', true, 1),
(22, 4, '/images/products/grocery/bananas.png', true, 2),
(23, 4, '/images/products/grocery/bananas.png', true, 3),
(24, 4, '/images/products/grocery/bananas.png', true, 4),
(25, 5, '/images/products/grocery/bananas.png', true, 0),
(26, 5, '/images/products/grocery/bananas.png', true, 1),
(27, 5, '/images/products/grocery/bananas.png', true, 2),
(28, 5, '/images/products/grocery/bananas.png', true, 3),
(29, 5, '/images/products/grocery/bananas.png', true, 4),
(30, 6, '/images/products/grocery/bananas.png', true, 0),
(31, 6, '/images/products/grocery/bananas.png', true, 1),
(32, 6, '/images/products/grocery/bananas.png', true, 2),
(33, 6, '/images/products/grocery/bananas.png', true, 3),
(34, 6, '/images/products/grocery/bananas.png', true, 4),
(35, 7, '/images/products/grocery/bananas.png', true, 0),
(36, 7, '/images/products/grocery/bananas.png', true, 1),
(37, 7, '/images/products/grocery/bananas.png', true, 2),
(38, 7, '/images/products/grocery/bananas.png', true, 3),
(39, 7, '/images/products/grocery/bananas.png', true, 4),
(40, 8, '/images/products/grocery/bananas.png', true, 0),
(41, 8, '/images/products/grocery/bananas.png', true, 1),
(42, 8, '/images/products/grocery/bananas.png', true, 2),
(43, 8, '/images/products/grocery/bananas.png', true, 3),
(44, 8, '/images/products/grocery/bananas.png', true, 4),
(45, 9, '/images/products/grocery/bananas.png', true, 0),
(46, 9, '/images/products/grocery/bananas.png', true, 1),
(47, 9, '/images/products/grocery/bananas.png', true, 2),
(48, 9, '/images/products/grocery/bananas.png', true, 3),
(49, 9, '/images/products/grocery/bananas.png', true, 4),

-- Products n10-n28 
(50, 10, '/images/products/grocery/bananas.png', true, 0),
(51, 12, '/images/products/grocery/bananas.png', true, 1),
(52, 13, '/images/products/grocery/bananas.png', true, 2),
(53, 14, '/images/products/grocery/bananas.png', true, 3),
(54, 15, '/images/products/grocery/bananas.png', true, 4),

-- Product images for n16-n25
(55, 16, '/images/products/grocery/bananas.png', true, 0),
(56, 17, '/images/products/grocery/bananas.png', true, 1),
(57, 18, '/images/products/grocery/bananas.png', true, 2),
(58, 19, '/images/products/grocery/bananas.png', true, 3),
(59, 20, '/images/products/grocery/bananas.png', true, 4),
(60, 21, '/images/products/grocery/bananas.png', true, 0),
(61, 22, '/images/products/grocery/bananas.png', true, 1),
(62, 23, '/images/products/grocery/bananas.png', true, 2),
(63, 24, '/images/products/grocery/bananas.png', true, 3),
(64, 25, '/images/products/grocery/bananas.png', true, 4),

-- Products n26-n28
(65, 26, '/images/products/grocery/bananas.png', false, 0),
(66, 27, '/images/products/grocery/bananas.png', false, 1),
(67, 28, '/images/products/grocery/bananas.png', false, 2),

-- Home category products with specialized images
(68, 29, '/images/products/home/lamp.jpg', true, 5),

-- Fashion category products with specialized images
(69, 30, '/images/products/fashion/denim-jacket.jpg', true, 5),
(70, 31, '/images/products/fashion/tshirt.jpg', true, 5),
(71, 32, '/images/products/fashion/sneakers.jpg', true, 5),
(72, 33, '/images/products/fashion/summer-dress.jpg', true, 5),
(73, 34, '/images/products/fashion/beanie.jpg', true, 5),

-- Electronic category products with specialized images
(74, 35, '/images/products/electronics/earbuds.jpg', true, 5),
(75, 36, '/images/products/electronics/smart-tv.jpg', true, 5),
(76, 37, '/images/products/electronics/bluetooth-speaker.jpg', true, 5),
(77, 38, '/images/products/electronics/digital-camera.jpg', true, 5),
(78, 39, '/images/products/electronics/smartwatch.jpg', true, 5),

-- Trending products (t1-t4, t6)
(80, 40, '/images/products/grocery/bananas.png', true, 0),
(81, 40, '/images/products/grocery/bananas.png', true, 1),
(82, 40, '/images/products/grocery/bananas.png', true, 2),
(83, 40, '/images/products/grocery/bananas.png', true, 3),
(84, 40, '/images/products/grocery/bananas.png', true, 4),
(85, 41, '/images/products/grocery/bananas.png', true, 0),
(86, 41, '/images/products/grocery/bananas.png', true, 1),
(87, 41, '/images/products/grocery/bananas.png', true, 2),
(88, 41, '/images/products/grocery/bananas.png', true, 3),
(89, 41, '/images/products/grocery/bananas.png', true, 4),
(90, 42, '/images/products/grocery/bananas.png', true, 0),
(91, 42, '/images/products/grocery/bananas.png', true, 1),
(92, 42, '/images/products/grocery/bananas.png', true, 2),
(93, 42, '/images/products/grocery/bananas.png', true, 3),
(94, 42, '/images/products/grocery/bananas.png', true, 4),
(95, 43, '/images/products/grocery/bananas.png', true, 0),
(96, 43, '/images/products/grocery/bananas.png', true, 1),
(97, 43, '/images/products/grocery/bananas.png', true, 2),
(98, 43, '/images/products/grocery/bananas.png', true, 3),
(99, 43, '/images/products/grocery/bananas.png', true, 4),

-- n28 (toys)
(100, 44, '/images/products/toys/art-kit.jpg', true, 5),

-- t6 (trending)
(101, 45, '/images/products/grocery/bananas.png', false, 0),
(102, 45, '/images/products/grocery/bananas.png', false, 1),
(103, 45, '/images/products/grocery/bananas.png', false, 2),
(104, 45, '/images/products/grocery/bananas.png', false, 3),
(105, 45, '/images/products/grocery/bananas.png', false, 4),
(106, 45, '/images/products/grocery/laptop.png', true, 5),

-- n26, n27 (laptop products)
(107, 46, '/images/products/grocery/cookies.jpg', true, 5),
(108, 47, '/images/products/grocery/orange-juice.jpg', true, 5),

-- Trending products with specialized images (t7-t14)
(109, 48, '/images/trending/shirt.png', false, 0),
(110, 48, '/images/trending/shirt.png', false, 1),
(111, 48, '/images/trending/shirt.png', false, 2),
(112, 48, '/images/trending/serum.png', true, 3),
(113, 49, '/images/trending/shirt.png', false, 0),
(114, 49, '/images/trending/shirt.png', false, 1),
(115, 49, '/images/trending/shirt.png', false, 2),
(116, 49, '/images/trending/vase.png', true, 3),
(117, 50, '/images/trending/shirt.png', false, 0),
(118, 50,'/images/trending/shirt.png', false, 1),
(119, 50, '/images/trending/shirt.png', false, 2),
(120, 50, '/images/trending/smartwatch.png', true, 3),
(121, 51, '/images/trending/shirt.png', false, 0),
(122, 51, '/images/trending/shirt.png', false, 1),
(123, 51, '/images/trending/shirt.png', false, 2),
(124, 51, '/images/trending/pants.png', true, 3),
(125, 52, '/images/trending/shirt.png', false, 0),
(126, 52, '/images/trending/shirt.png', false, 1),
(127, 52 , 'images/trending/shirt.png', false, 2),
(128, 52, '/images/trending/shirt.png', true, 3),
(129, 52, '/images/products/grocery/bananas.png', false, 4),
(130, 53, '/images/products/grocery/laptop.png', false, 5);



-- Add product tags
INSERT INTO product_tags (product_id, tag) VALUES
-- Regular products n1-n9 with appropriate tags
(1, 'laptop'),  
(1, 'apple'),
(1, 'macbook'),
(2, 'banana'),
(2, 'organic'),
(2, 'apple'),
(3, 'samsung'),
(3, 'banana'),
(4, 'product'),
(4, 'sony'),
(5, 'product'),
(5, 'sony'),
(6, 'product'),
(6, 'xiaomi'),
(7, 'product'),
(7, 'sony'),
(8, 'product'),
(8, 'sony'),
(9, 'product'),
(9, 'sony'),

-- Products n10-n28 with brand tags
(10, 'home'),
(12, 'home'),
(13, 'home'),
(14, 'home'),
(15, 'home'),
(16, 'fashion'),
(17, 'fashion'),
(18, 'fashion'),
(19, 'fashion'),
(20, 'fashion'),
(21, 'electronics'),
(22, 'electronics'),
(23, 'electronics'),
(24, 'electronics'),
(25, 'electronics'),
(26, 'grocery'),
(27, 'grocery'),
(28, 'toys'),

-- Trending products (t1-t6)
(29, 'toy'),
(29, 'building'),
(29, 'educational'),
(30, 'toy'),
(30, 'remote'),
(30, 'piano'),
(31, 'toy'),
(31, 'teddy'),
(31, 'piano'),
(32, 'educational'),
(32, 'board game'),
(32, 'family'),
(33, 'coffee'),
(33, 'instant'),
(33, 'piano'),

-- More trending products (t7-t14)
(34, 'beauty'),
(34, 'serum'),
(34, 'skincare'),
(35, 'home'),
(35, 'vase'),
(35, 'decoration'),
(36, 'electronics'),
(36, 'smartwatch'),
(36, 'fitness'),
(37, 'fashion'),
(37, 'pants'),
(37, 'clothing'),
(37, 'shirt'),


-- Limited edition products (l1-l7)
(38, 'samsung'),
(38, 'smartphone'),
(38, 'electronics'),
(39, 'xiaomi'),
(39, 'earbuds'),
(39, 'audio'),
(39, 'piano'),
(39, 'jacket'),
(39, 'winter'),
(39, 'new balance'),
(39, 'sneakers'),
(39, 'footwear'),
(40, 'xiaomi'),
(40, 'earbuds'),
(40, 'audio'),
(41, 'electronics'),
(41, 'xiaomi'),

-- Top sale products (top1-top6)
(42, 'computer'),
(42, 'accessories'),
(42, 'apple'),
(43, 'outfit'),
(43, 'men'),
(43, 'casual'),
(43, 'juice'),
(43, 'drink'),
(43, 'healthy'),
(43, 'pet'),
(43, 'dog'),
(43, 'food'),
(45, 'security'),
(45, 'camera'),
(45, 'electronics'),
(46, 'cosmetic'),
(46, 'beauty'),
(46, 'skincare');



-- Add product reviews
INSERT INTO product_reviews (id, product_id, user_id, user_name, rating, comment, created_at) VALUES
-- Product reviews for existing products
(1, 1, 'u123456', 'John Smith', 5, 'Best MacBook I''ve ever had! The performance is amazing.', CURRENT_TIMESTAMP - INTERVAL '3 days'),
(2, 1, 'u234567', 'Emily Johnson', 4, 'Great laptop overall, but battery life could be better.', CURRENT_TIMESTAMP - INTERVAL '7 days');


-- Update product counts for categories
UPDATE categories c
SET product_count = (
    SELECT COUNT(*) FROM products p
    WHERE p.category_id = c.id
);
-- Insert banner slider data
INSERT INTO banners (
    title, subtitle, description, discount, highlight_text, 
    image_url, link_url, action_text, background_color, 
    text_color, animation_type, is_active, priority, type, 
    product_id, category_id
) VALUES 
(
    'Discover the Latest Deals', 
    'Limited Time Offer', 
    'Explore our collection of premium products with exclusive discounts',
    'Up to 50% Off!',
    'Shop Now',
    'https://images.unsplash.com/photo-1607082348824-0a96f2a4b9da?w=1440&h=500&fit=crop',
    '/collections/deals',
    'Shop Now',
    '#1e3a8a',
    '#ffffff',
    'fade',
    true,
    1,
    'hero',
    48, -- Computer Accessories product
    4 -- Electronics category
),
(
    'New Season Collection', 
    'Spring/Summer 2024',
    'Refresh your wardrobe with our latest arrivals',
    '20% Off for Members',
    'Exclusive',
    'https://images.unsplash.com/photo-1483985988355-763728e1935b?w=1440&h=500&fit=crop',
    '/collections/summer',
    'Explore',
    '#e9ecef',
    '#343a40',
    'slide',
    true,
    2,
    'hero',
    49, -- Men's Casual Outfit product
    7 -- Men category
),
(
    'Tech Innovations', 
    'Next Generation Gadgets',
    'Discover cutting-edge technology for your lifestyle',
    'Free Shipping on Orders $100+',
    'New Arrivals',
    'https://images.unsplash.com/photo-1499951360447-b19be8fe80f5?w=1440&h=500&fit=crop',
    '/collections/tech',
    'Discover',
    '#d8e2dc',
    '#2b2d42',
    'zoom',
    true,
    3,
    'hero',
    25, -- Smartwatch product
    4 -- Electronics category
),
(
    'Fashion Flash Sale', 
    'Limited Time Only',
    'Get the latest fashion trends at unbeatable prices',
    'Buy One Get One Free',
    'Hot Deal',
    'https://images.unsplash.com/photo-1487222477894-8943e31ef7b2?w=1440&h=500&fit=crop',
    '/category/fashion',
    'Shop Now',
    '#023047',
    '#ffb703',
    'fade',
    true,
    4,
    'hero',
    16, -- Women's Denim Jacket product
    3 -- Fashion category
),
(
    'Home Essentials', 
    'Upgrade Your Living Space',
    'Discover beautiful furniture and decor for your home',
    'Up to 30% Off Select Items',
    'Limited Stock',
    'https://images.unsplash.com/photo-1556228453-efd6c1ff04f6?w=1440&h=500&fit=crop',
    '/category/home',
    'Explore Collection',
    '#264653',
    '#e9c46a',
    'slide',
    true,
    5,
    'promotional',
    12, -- Cotton Bed Sheets product
    2 -- Home category
),
(
    'Holiday Gift Guide', 
    'Perfect Presents for Everyone',
    'Find the ideal gifts for your loved ones this holiday season',
    'Special Holiday Pricing',
    'Gift Ideas',
    'https://images.unsplash.com/photo-1512909006721-3d6018887383?w=1440&h=500&fit=crop',
    '/collections/holiday',
    'Explore Gifts',
    '#3a0ca3',
    '#f72585',
    'zoom',
    false, -- Inactive for now, can be activated during holiday season
    6,
    'seasonal',
    NULL, -- Not tied to specific product
    NULL -- Not tied to specific category
);


-- =====================================================
-- New Fashion item (NewFashion.tsx component)
-- =====================================================
INSERT INTO ads_placement (
    location, reference_type, reference_id, display_order, 
    custom_title, custom_image_url, ui_settings, is_active
) VALUES 
-- New Fashion banner
(
    'new_fashion', 'product', 1, 1, -- Using Fashion category (ID=3)
    'New Year! New Fashion', 'images/ads/products/new-fashion.png',
    '{"button_type": "secondary"}',
    true
);


-- Add special seasonal banners
INSERT INTO banners (
    title, subtitle, description, discount, highlight_text, 
    image_url, link_url, action_text, background_color, 
    text_color, animation_type, is_active, priority, type
) VALUES 
(
    'Back to School', 
    'Gear Up for Success',
    'Everything students need for the new academic year',
    'Student Discount: 15% Off',
    'Limited Time',
    'https://images.unsplash.com/photo-1503676260728-1c00da094a0b?w=1440&h=500&fit=crop',
    '/collections/back-to-school',
    'Shop Essentials',
    '#184e77',
    '#d9ed92',
    'fade',
    false, -- Will be activated during back-to-school season
    10,
    'seasonal'
),
(
    'Summer Vibes', 
    'Hot Days, Cool Deals',
    'Get ready for summer with our exclusive collection',
    'Seasonal Offers',
    'Summer Must-Haves',
    'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=1440&h=500&fit=crop',
    '/collections/summer',
    'Shop Collection',
    '#00b4d8',
    '#ffffff',
    'slide',
    false, -- Will be activated during summer
    11,
    'seasonal'
);

-- Add category-specific banners
INSERT INTO banners (
    title, subtitle, description, discount, highlight_text, 
    image_url, link_url, action_text, background_color, 
    text_color, animation_type, is_active, priority, type, 
    category_id
) VALUES 
(
    'Electronics Bonanza', 
    'Smart Devices for Modern Living',
    'Discover the latest gadgets and electronics',
    'Up to 40% Off Select Items',
    'Tech Deals',
    'https://images.unsplash.com/photo-1550009158-9ebf69173e03?w=1440&h=500&fit=crop',
    '/category/electronics',
    'Shop Electronics',
    '#2c3e50',
    '#ecf0f1',
    'fade',
    true,
    20,
    'category',
    4 -- Electronics category
),
(
    'Fashion Forward', 
    'Express Your Style',
    'Trendy outfits for every occasion',
    'Buy 2 Get 1 Free',
    'Trending Styles',
    'https://images.unsplash.com/photo-1490481651871-ab68de25d43d?w=1440&h=500&fit=crop',
    '/category/fashion',
    'Shop Fashion',
    '#5e548e',
    '#ffffff',
    'slide',
    true,
    21,
    'category',
    3 -- Fashion category
);

-- Add product-specific banners
INSERT INTO banners (
    title, subtitle, description, discount, highlight_text, 
    image_url, link_url, action_text, background_color, 
    text_color, animation_type, is_active, priority, type, 
    product_id
) VALUES 
(
    'Limited Edition Smartwatch', 
    'Track Your Health in Style',
    'Advanced health monitoring and sleek design',
    '25% Launch Discount',
    'New Release',
    'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=1440&h=500&fit=crop',
    '/products/smartwatch',
    'Shop Now',
    '#1a1a2e',
    '#e3f6f5',
    'zoom',
    true,
    30,
    'product',
    25 -- Smartwatch product
);


-- =====================================================
-- Banner Showcase items (BannerShowcase.tsx component)
-- =====================================================
INSERT INTO ads_placement (
    location, reference_type, reference_id, display_order, 
    custom_image_url, ui_settings, is_active
) VALUES 
-- Banner 1: Trousers fashion image
(
    'banner', 'category', 7, 1, -- Using Men category (ID=7)
    'images/ads/banner/trousers_fashion.png',
    '{"has_more": false}',
    true
),
-- Banner 2: Watchmen fashion with gray Shop Now button
(
    'banner', 'category', 7, 2, -- Using Men category (ID=7)
    'images/ads/banner/watchmen_fashion.png',
    '{"button_type": "gray", "button_text": "Shop Now", "has_more": false}',
    true
),
-- Banner 3: Denim fashion with "see more" option
(
    'banner', 'category', 10, 3, -- Using Women category (ID=10)
    'images/ads/banner/denim_fashion.png',
    '{"has_more": true}',
    true
),
-- Banner 4: Domestic item with black Shop Now button
(
    'banner', 'category', 2, 4, -- Using Home category (ID=2)
    'images/ads/banner/dometic.png',
    '{"button_type": "black", "button_text": "Shop Now", "has_more": false}',
    true
);

-- =====================================================
-- Display Grid items (DisplayGrid.tsx component)
-- =====================================================
INSERT INTO ads_placement (
    location, reference_type, reference_id, display_order, 
    custom_image_url, ui_settings, is_active
) VALUES 
-- Display item 1: Be Winner promotional banner
(
    'display', 'product', 22, 1, -- Using 4K Smart TV (ID=22)
    'images/ads/display/be-winner.png',
    '{}',
    true
),
-- Display item 2: Redmi Y3 with gradient button
(
    'display', 'product', 21, 2, -- Using Wireless Earbuds (ID=21)
    'images/ads/display/redmi-y3.png',
    '{"button_type": "gradient"}',
    true
),
-- Display item 3: Philips Ambilight TV with price and discount
(
    'display', 'product', 22, 3, -- Using 4K Smart TV (ID=22)
    'images/ads/display/ambilighttv.png',
    '{"button_type": "secondary", "discount_img": "images/ads/display/discount_img.png"}',
    true
);

-- =====================================================
-- Gaming Banner items (GamingBanner.tsx component)
-- =====================================================
INSERT INTO ads_placement (
  location,
  reference_type,
  reference_id,
  display_order,
  custom_title,
  custom_image_url,
  is_active
) VALUES
  -- Gaming item 1: Headsets
  (
    'gaming', 'category', 4, 1,
    'Headsets', 'images/ads/gaming/headsets.png',
    true
  ),
  -- Gaming item 2: Mouse
  (
    'gaming', 'category', 4, 2,
    'Mouse', 'images/ads/gaming/mouse.png',
    true
  ),
  -- Gaming item 3: Controller
  (
    'gaming', 'category', 4, 3,
    'Controller', 'images/ads/gaming/controller.png',
    true
  ),
  -- Gaming item 4: Chair
  (
    'gaming', 'category', 4, 4,
    'Chair', 'images/ads/gaming/chair.png',
    true
  );

