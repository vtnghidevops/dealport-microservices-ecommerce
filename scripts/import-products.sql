INSERT INTO products (
  id, type, name, description, slug, price, original_price, discount,
  category_id, category_slug, stock_quantity, brand,
  features, shipping_info, ui_metadata, orders
) VALUES
(
    'n1', 'normal', '2020 Apple MacBook Pro with Apple M1 Chip (13-inch, 8GB RAM, 256GB SSD Storage) - Space Gray', 'The MacBook Air 13-inch is a slim and lightweight laptop from Apple with a sharp Retina display and impressive performance powered by the M1 chip. It''s known for its portability, sleek design, and reliable performance.', 'apple-macbook-pro-m1-chip-13-inch-8gb-ram-256gb-ssd-storage-space-gray', 100.25, 125.99, 25,
    'g1', 'grocery', 5, 'Apple',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n2', 'normal', 'Organic Bananas from sustainable farms', 'Fresh organic bananas from sustainable farms', 'organic-bananas', 2.99, 3.99, 25,
    'g1', 'grocery', 19, 'Apple',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n3', 'normal', 'Organic Bananas from sustainable farms', 'Fresh organic bananas from sustainable farms', 'organic-bananas1', 2.99, 3.99, 25,
    'g1', 'grocery', 0, 'Samsung',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n4', 'normal', 'Organic Bananas from sustainable farms', 'Fresh organic bananas from sustainable farms', 'organic-bananas2', 2.99, 3.99, 25,
    'g1', 'grocery', 150, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n5', 'normal', 'Whole Grain Bread', 'Freshly baked whole grain bread', 'whole-grain-bread', 3.49, 3.99, 25,
    'g1', 'grocery', 75, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n6', 'normal', 'Free Range Eggs (12pk)', 'Farm fresh free-range eggs', 'free-range-eggs', 4.99, 3.99, 25,
    'g1', 'grocery', 120, 'xiaomi',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n7', 'normal', 'Organic Milk', 'Fresh organic whole milk', 'organic-milk', 3.99, 3.99, 25,
    'g1', 'grocery', 95, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n8', 'normal', 'Avocados (3pk)', 'Ripe and ready to eat avocados', 'avocados1', 5.99, 3.99, 25,
    'g1', 'grocery', 65, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n9', 'normal', 'Avocados (3pk)', 'Ripe and ready to eat avocados', 'avocados2', 5.99, 3.99, 25,
    'g1', 'grocery', 65, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n1', 'normal', 'Scented Candle Set', 'Set of 3 aromatic scented candles', 'scented-candle-set', 24.99, 3.99, 25,
    'h2', 'home', 45, 'xiaomi',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n2', 'normal', 'Cotton Bed Sheets', 'Soft 100% Egyptian cotton bed sheets', 'cotton-bed-sheets', 49.99, 3.99, 25,
    'h2', 'home', 30, 'xiaomi',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n3', 'normal', 'Kitchen Utensil Set', 'Complete set of silicone kitchen utensils', 'kitchen-utensil-set1', 35.99, 3.99, 25,
    'h2', 'home', 55, 'apple',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n4', 'normal', 'Decorative Throw Pillows', 'Set of 2 decorative throw pillows for sofa or bed', 'decorative-throw-pillows1', 22.99, 3.99, 25,
    'h2', 'home', 0, 'panasonic',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n5', 'normal', 'LED Table Lamp', 'Modern LED table lamp with adjustable brightness', 'led-table-lamp', 39.99, 3.99, 25,
    'h2', 'home', 25, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n1', 'normal', 'Women''s Denim Jacket', 'Classic denim jacket for women', 'womens-denim-jacket', 59.99, 3.99, 25,
    'f3', 'fashion', 40, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n2', 'normal', 'Men''s Casual T-Shirt', '100% cotton casual t-shirt for men', 'mens-casual-tshirt', 19.99, 3.99, 25,
    'f3', 'fashion', 120, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n3', 'normal', 'Leather Sneakers', 'Comfortable leather sneakers for everyday use', 'leather-sneakers', 89.99, 3.99, 25,
    'f3', 'fashion', 35, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n4', 'normal', 'Summer Dress', 'Lightweight cotton summer dress', 'summer-dress', 45.99, 3.99, 25,
    'f3', 'fashion', 0, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n5', 'normal', 'Unisex Beanie', 'Warm knitted beanie for all seasons', 'unisex-beanie', 14.99, 3.99, 25,
    'f3', 'fashion', 60, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n1', 'normal', 'Wireless Earbuds', 'Bluetooth wireless earbuds with noise cancellation', 'wireless-earbuds', 129.99, 3.99, 25,
    'e4', 'electronic', 45, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n2', 'normal', '4K Smart TV', '55-inch 4K Ultra HD Smart LED TV', '4k-smart-tv', 699.99, 3.99, 25,
    'e4', 'electronic', 15, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n3', 'normal', 'Portable Bluetooth Speaker', 'Waterproof portable bluetooth speaker with 20-hour battery life', 'portable-bluetooth-speaker', 89.99, 3.99, 25,
    'e4', 'electronic', 30, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n4', 'normal', 'Digital Camera', '24MP digital camera with 4K video recording', 'digital-camera', 449.99, 3.99, 25,
    'e4', 'electronic', 0, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n5', 'normal', 'Smartwatch', 'Fitness tracking smartwatch with heart rate monitor', 'smartwatch', 199.99, 3.99, 25,
    'e4', 'electronic', 25, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    't1', 'litmited', 'Building Blocks Set', 'Educational building blocks set for children', 'building-blocks-set', 29.99, 3.99, 25,
    't5', 'toys', 50, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    't2', 'trending', 'Remote Control Car', 'High-speed remote control race car', 'remote-control-car', 49.99, 3.99, 25,
    't5', 'toys', 35, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 10
  ),
(
    't3', 'trending', 'Stuffed Teddy Bear', 'Soft and huggable teddy bear for kids', 'stuffed-teddy-bear', 19.99, 3.99, 25,
    't5', 'toys', 80, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    't4', 'trending', 'Educational Board Game', 'Family-friendly educational board game', 'educational-board-game', 34.99, 3.99, 25,
    't5', 'toys', 0, 'Educational Games Co.',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n5', 'normal', 'Art and Craft Kit', 'Creative art and craft kit for children', 'art-and-craft-kit', 24.99, 3.99, 25,
    't5', 'toys', 40, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    't6', 'trending', 'Instant Coffee', 'Premium instant coffee for quick preparation', 'instant-coffee', 7.99, 3.99, 25,
    'l6', 'laptop', 70, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n7', 'normal', 'Chocolate Cookies', 'Crunchy chocolate chip cookies', 'chocolate-cookies', 3.99, 3.99, 25,
    'l6', 'laptop', 90, 'piano',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'n8', 'normal', 'Fresh Orange Juice', '100% freshly squeezed orange juice', 'fresh-orange-juice', 4.49, 3.99, 25,
    'l6', 'laptop', 55, 'sony',
    '["Free Shipping","Free 1 Year Warranty","100% Money-back guarantee","Secure payment method","24/7 Customer support"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    't1', 'trending', 'Radiant Glow Hydrating Serum', 'Gentle yet effective, our Radiant Boosting Foaming our Radiant Boosting Foaming our Radiant Boosting Foaming', 'radiant-glow-hydrating-serum', 29.99, 39.99, 20,
    '1', 'grocery', 5, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    't2', 'trending', 'Modern Minimalist Vase', 'Track your workouts, heart rate, sleep quality and receive notifications. Water resistant up to 50m with 7-day battery life.', 'modern-minimalist-vase', 40.99, 49.99, 18,
    '1', 'grocery', 5, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    't3', 'trending', 'FitPro 3000 Smart Watch', 'Fast-charging power bank with dual USB ports and USB-C compatibility. Charge multiple devices simultaneously on the go.', 'fitpro-3000-smart-watch', 119.99, 0, 0,
    '1', 'grocery', 5, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    '1', 'trending', 'Pants', 'Pants', 'pants', 25.95, 45.95, 20,
    '1', 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    '2', 'trending', 'Shirt', 'Shirt', 'shirt', 12.99, 12.99, 0,
    '1', 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    '3', 'trending', 'Hat', 'Hat', 'hat', 104, 104, 0,
    '1', 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    '4', 'trending', 'Shoe', 'Shoe', 'shoe', 10.56, 10.56, 0,
    '1', 'men', 100, 'Brand 1',
    '["feature1","feature2","feature3"]', '{"courier":"Courier 1","local":"Local 1","ups":"UPS 1","global":"Global 1"}', '{}', 100
  ),
(
    'l1', 'limited', 'Samsung Galaxy S24', 'Gentle yet effective, our hydrating serum infuses skin with essential moisture while brightening and plumping for a radiant complexion.', 'samsung-galaxy-s24', 29.99, 39.99, 25,
    'l1', 'limited', 0, 'samsung',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'l2', 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud', 34.99, 42.99, 18,
    'l1', 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'l3', 'limited', 'Winter fashion jacket', 'Delivers extreme hydration and helps strengthen the skin barrier. Perfect for both daytime wear and overnight rejuvenation. [[2]]', 'winter-fashion-jacket', 39.99, 49.99, 20,
    'l1', 'limited', 0, 'piano',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'l4', 'limited', 'New Balance 574 Senekers', 'Illuminating serum infused with glow-boosting intelligent botanicals that leave the skin visibly radiant while reducing redness. [[4]]', 'new-balance-574-senekers', 44.99, 59.99, 25,
    'l1', 'limited', 0, 'new balance',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'l5', 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud', 34.99, 42.99, 18,
    'l1', 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'l6', 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud', 34.99, 42.99, 18,
    'l1', 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    'l7', 'limited', 'Ui TWS 7002 Earbud', 'Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]', 'ui-tws-7002-earbud', 34.99, 42.99, 18,
    'l1', 'limited', 0, 'xiaomi',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{}', 100
  ),
(
    't1', 'top-sale', 'Computer Accessories', 'category', 'computer-accessories', 12.99, 12.99, 0,
    't1', 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isCommingSoon":true}', 100
  ),
(
    't2', 'top-sale', 'Men''s Casual Outfit', 'category', 'men-casual-outfit', 200, 200, 0,
    't1', 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isCommingSoon":false}', 100
  ),
(
    't3', 'top-sale', 'Pome Granate Juice', 'category', 'pome-granate-juice', 49, 49, 0,
    't1', 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isCommingSoon":false}', 100
  ),
(
    't4', 'top-sale', 'Dog Food Made With Love', 'category', 'dog-food-made-with-love', 12.99, 12.99, 0,
    't1', 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"col","isCommingSoon":false}', 100
  ),
(
    't5', 'top-sale', 'Security Camera System', 'category', 'security-camera-system', 420, 420, 0,
    't1', 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"double","isCommingSoon":false}', 100
  ),
(
    't6', 'top-sale', 'Premium Cosmetic Set', 'category', 'premium-cosmetic-set', 149, 149, 0,
    't1', 'top-sell', 0, 'apple',
    '["feature1","feature2","feature3"]', '{"courier":"2-4 days, free shipping","local":"up to one week, $19.00","ups":"4-6 days, $29.00","global":"3-4 days, $39.00"}', '{"setUpDesign":"row","isCommingSoon":false}', 100
  )
ON CONFLICT (id, slug) DO UPDATE SET
  name = EXCLUDED.name,
  type = EXCLUDED.type,
  description = EXCLUDED.description,
  slug = EXCLUDED.slug,
  price = EXCLUDED.price,
  original_price = EXCLUDED.original_price,
  discount = EXCLUDED.discount,
  category_id = EXCLUDED.category_id,
  category_slug = EXCLUDED.category_slug,
  stock_quantity = EXCLUDED.stock_quantity,
  brand = EXCLUDED.brand,
  features = EXCLUDED.features,
  shipping_info = EXCLUDED.shipping_info,
  ui_metadata = EXCLUDED.ui_metadata,
  orders = EXCLUDED.orders,
  updated_at = CURRENT_TIMESTAMP;
