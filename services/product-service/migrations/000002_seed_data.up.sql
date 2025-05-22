-- Insert categories with UUIDs
INSERT INTO categories (id, name, slug, description, image_url, product_count, is_active, is_visible) 
VALUES 
(1, 'Grocery', 'grocery', 'Grocery', '/images/category/grocery.png', 0, true, true),
(2, 'Household Appliances', 'household-appliances', 'Household Appliances', '/images/category/household-appliances.png', 0, true, true),
(3, 'Fashion', 'fashion', 'Fashion', '/images/category/fashion.png', 0, true, true),
(4, 'Electronics', 'electronics', 'Electronics', '/images/category/electronics.png', 0, true, true),
(5, 'Toys', 'toys', 'Toys', '/images/category/toys.png', 0, true, true),
(6, 'Laptops', 'laptops', 'Laptops', '/images/category/laptops.png', 0, true, true),
(7, 'Men', 'men', 'Men', '/images/category/men.png', 0, true, true), 
(8, 'Fruits', 'fruits', 'Fruits', '/images/category/fruits.png', 0, true, true),
(9, 'Vegetables', 'vegetables', 'Vegetables', '/images/category/vegetables.png', 0, true, true),
(10, 'Women', 'women', 'Women', '/images/category/women.png', 0, true, true),
(11, 'Kids', 'kids', 'Kids', '/images/category/kids.png', 0, true, true),
(12, 'Piano', 'piano', 'Piano', '/images/category/piano.png', 0, true, true);

-- Insert products
INSERT INTO products (
  id, type, name, description, slug, price, original_price, discount,
  category_id, category_slug, stock_quantity, brand,
  features, shipping_info, ui_metadata, orders
) VALUES
(
    1, 'normal', 'Apple MacBook Pro M1 (2020) - 13-inch, 8GB RAM, 256GB SSD', 'Powerful and efficient laptop featuring the revolutionary M1 chip with 8-core CPU, delivering up to 2.8x faster performance. Includes a stunning Retina display, all-day battery life, and macOS.', 'apple-macbook-pro-m1-chip-13-inch-8gb-ram-256gb-ssd-storage-space-gray', 999.99, 1299.99, 23,
    1, 'grocery', 25, 'Apple',
    '["Free Express Shipping","1 Year Apple Warranty","30-day Return Policy","Secure Checkout","24/7 Technical Support"]', '{"courier":"1-2 days, free shipping","local":"3-5 days, $19.00","ups":"2-3 days, $29.00","global":"5-7 days, $49.00"}', '{}', 245
  ),
(
    2, 'normal', 'Organic Fair Trade Bananas (Bunch of 7)', 'Perfectly ripened organic bananas from sustainable, fair trade certified farms. Rich in potassium and naturally sweet flavor. Each bunch contains approximately 7 bananas.', 'organic-bananas', 3.99, 4.99, 20,
    8, 'fruits', 154, 'Nature''s Best',
    '["Certified Organic","Fair Trade Certified","No Pesticides","Sustainably Grown","Farm to Table"]', '{"courier":"Same-day delivery available","local":"Next-day delivery, $5.00","pickup":"Free store pickup"}', '{}', 1240
  ),
(
    3, 'normal', 'Premium Red Delicious Apples (Pack of 6)', 'Sweet and crisp red delicious apples, carefully selected and packed fresh. Perfect for snacking, baking, or adding to salads. Each apple is hand-picked at peak ripeness.', 'premium-red-apples', 5.49, 6.99, 21,
    8, 'fruits', 89, 'Orchard Fresh',
    '["Pesticide-Free","Hand-Picked","Premium Selection","Farm Fresh","High in Fiber"]', '{"courier":"Same-day delivery, $2.99","local":"Next-day delivery, $5.00","pickup":"Free store pickup"}', '{}', 876
  ),
(
    4, 'normal', 'Seasonal Organic Berries Mix (16oz)', 'Delicious assortment of organic strawberries, blueberries, and raspberries. Perfect for smoothies, desserts, or enjoying as a healthy snack. Packed with antioxidants and vitamins.', 'seasonal-berries-mix', 7.99, 9.99, 20,
    8, 'fruits', 67, 'Berry Farms',
    '["USDA Organic","Hand-Picked","Seasonal Selection","Rich in Antioxidants","Vitamin C-Rich"]', '{"courier":"Same-day cold delivery, $3.99","local":"Next-day cold delivery, $6.00","pickup":"Free store pickup"}', '{}', 543
  ),
(
    5, 'normal', 'Artisan Sourdough Bread (24oz Loaf)', 'Freshly baked artisan sourdough with a perfect crust and tender, tangy interior. Made with a 100-year-old starter and slow-fermented for 24 hours. No preservatives or additives.', 'artisan-sourdough-bread', 6.49, 7.99, 19,
    1, 'grocery', 42, 'Artisan Bakery',
    '["Handcrafted Daily","No Preservatives","Traditional Recipe","Natural Ingredients","Source of Prebiotics"]', '{"courier":"Same-day delivery, $2.99","local":"Next-day delivery, $5.00","pickup":"Free store pickup"}', '{}', 789
  ),
(
    6, 'normal', 'Free-Range Organic Eggs (Dozen)', 'Farm-fresh, certified organic eggs from free-range chickens. Chickens are raised on open pastures with access to natural diets and no antibiotics. Rich, golden yolks and exceptional flavor.', 'free-range-organic-eggs', 5.99, 7.49, 20,
    1, 'grocery', 103, 'Happy Hens Farm',
    '["Certified Organic","Free-Range","No Antibiotics","High in Protein","Hormone-Free"]', '{"courier":"Same-day delivery, $2.99","local":"Next-day delivery, $5.00","pickup":"Free store pickup"}', '{}', 932
  ),
(
    7, 'normal', 'Grass-Fed Organic Whole Milk (64oz)', 'Premium organic whole milk from grass-fed cows raised on sustainable family farms. Non-homogenized with cream on top for a rich, natural taste. Packaged in recyclable glass bottles.', 'grass-fed-whole-milk', 5.49, 6.99, 21,
    1, 'grocery', 78, 'Green Pastures',
    '["Grass-Fed","Certified Organic","Non-Homogenized","No Hormones","Glass Bottled"]', '{"courier":"Same-day cold delivery, $3.99","local":"Next-day cold delivery, $6.00","pickup":"Free store pickup"}', '{}', 654
  ),
(
    8, 'normal', 'Organic Hass Avocados (Pack of 4)', 'Premium organic Hass avocados, perfectly ripened and ready to eat. Rich, creamy texture and buttery flavor. High in healthy fats, fiber, and potassium.', 'organic-hass-avocados', 6.99, 8.49, 18,
    8, 'fruits', 56, 'Avocado Farms',
    '["Certified Organic","Perfectly Ripened","Rich in Nutrients","Heart-Healthy","Sustainably Grown"]', '{"courier":"Same-day delivery, $2.99","local":"Next-day delivery, $5.00","pickup":"Free store pickup"}', '{}', 765
  ),
(
    9, 'normal', 'Fresh Tropical Mango (Pack of 3)', 'Sweet and juicy tropical mangoes, handpicked at peak ripeness. Vibrant orange flesh with a rich, aromatic flavor. Perfect for smoothies, salads, or enjoying fresh.', 'fresh-tropical-mangoes', 7.49, 8.99, 17,
    8, 'fruits', 48, 'Tropical Harvest',
    '["Naturally Ripened","High in Vitamin A","Rich in Fiber","No Preservatives","Hand-Selected"]', '{"courier":"Same-day delivery, $2.99","local":"Next-day delivery, $5.00","pickup":"Free store pickup"}', '{}', 432
  ),
(
    10, 'normal', 'Luxury Scented Soy Candle Set (Set of 3)', 'Hand-poured soy wax candles in three signature scents: Lavender Dream, Vanilla Comfort, and Fresh Linen. Made with essential oils and cotton wicks for a clean, long-lasting burn up to 45 hours each.', 'luxury-scented-candle-set', 29.99, 39.99, 25,
    2, 'household-appliances', 62, 'Tranquil Home',
    '["100% Soy Wax","Essential Oil Blends","Lead-free Cotton Wicks","45-Hour Burn Time","Recyclable Glass Jars"]', '{"courier":"2-3 days delivery, $4.99","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 321
  ),
(
    12, 'normal', 'Premium Egyptian Cotton Bed Sheet Set (Queen)', 'Luxurious 1000-thread count Egyptian cotton sheet set including flat sheet, fitted sheet, and 2 pillowcases. Ultra-soft, breathable, and durable with a subtle sateen finish for a hotel-quality sleep experience.', 'premium-egyptian-cotton-bed-sheets', 89.99, 119.99, 25,
    2, 'household-appliances', 34, 'Luxe Linens',
    '["1000 Thread Count","100% Egyptian Cotton","OEKO-TEX Certified","Deep Pocket Fitted Sheet","Wrinkle-Resistant"]', '{"courier":"2-3 days delivery, $4.99","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 267
  ),
(
    13, 'normal', 'Professional Stainless Steel Kitchen Utensil Set (10-Piece)', 'Complete kitchen utensil set featuring ergonomic handles and premium stainless steel construction. Includes spatula, slotted spoon, solid spoon, pasta server, ladle, whisk, tongs, grater, peeler, and measuring spoons.', 'professional-kitchen-utensil-set', 42.99, 59.99, 28,
    2, 'household-appliances', 89, 'Chef''s Elite',
    '["Premium Stainless Steel","Dishwasher Safe","Heat-Resistant Handles","Non-Toxic Materials","Hanging Storage Loops"]', '{"courier":"2-3 days delivery, $4.99","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 456
  ),
(
    14, 'normal', 'Designer Throw Pillow Collection (Set of 4)', 'Elegant decorative throw pillows in complementary patterns and textures. Set includes 2 geometric patterned and 2 solid velvet pillows with hidden zipper closures and plush polyester filling.', 'designer-throw-pillow-collection', 59.99, 79.99, 25,
    2, 'household-appliances', 42, 'Modern Décor',
    '["Removable Covers","Machine Washable","Premium Fabrics","Hidden Zippers","Hypoallergenic Filling"]', '{"courier":"2-3 days delivery, $4.99","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 198
  ),
(
    15, 'normal', 'Smart LED Desk Lamp with Wireless Charging', 'Modern desk lamp featuring adjustable brightness levels, color temperature control, and built-in wireless charging pad. Includes USB port, touch controls, and timer function.', 'smart-led-desk-lamp', 64.99, 84.99, 24,
    2, 'household-appliances', 56, 'TechLight',
    '["Wireless Charging Pad","USB Charging Port","5 Brightness Levels","4 Color Temperatures","Auto-Off Timer"]', '{"courier":"2-3 days delivery, $4.99","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 345
  ),
(
    16, 'normal', 'Vintage-Inspired Women''s Denim Jacket', 'Classic denim jacket with a modern fit and vintage-inspired detailing. Features adjustable button cuffs, chest flap pockets, and subtle distressing for an authentic look. Made from premium stretch denim for comfort and durability.', 'vintage-womens-denim-jacket', 79.99, 99.99, 20,
    3, 'fashion', 78, 'Urban Heritage',
    '["Premium Stretch Denim","Button Closure","Adjustable Waist Tabs","Machine Washable","Sustainably Manufactured"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 432
  ),
(
    17, 'normal', 'Men''s Premium Cotton Graphic T-Shirt', 'Soft and comfortable men''s t-shirt made from 100% organic combed cotton. Features a contemporary fit, reinforced seams, and original graphic print. Pre-shrunk and colorfast for lasting quality.', 'mens-premium-graphic-tshirt', 34.99, 44.99, 22,
    7, 'men', 124, 'Modern Basics',
    '["100% Organic Cotton","Original Artwork","Pre-Shrunk Fabric","Reinforced Stitching","Tagless Comfort"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 567
  ),
(
    18, 'normal', 'Italian Leather Minimalist Sneakers', 'Handcrafted minimalist sneakers made from premium Italian leather. Features cushioned insoles, durable rubber outsoles, and waxed cotton laces. Versatile design perfect for casual or semi-formal occasions.', 'italian-leather-minimalist-sneakers', 129.99, 159.99, 19,
    3, 'fashion', 63, 'Milano Footwear',
    '["Genuine Italian Leather","Cushioned Memory Foam Insole","Antibacterial Lining","Handcrafted Construction","Durable Rubber Outsole"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 389
  ),
(
    19, 'normal', 'Bohemian Floral Maxi Summer Dress', 'Flowing maxi dress featuring a vibrant floral print on lightweight rayon fabric. Includes adjustable spaghetti straps, smocked bodice, and tiered skirt with subtle high-low hem. Perfect for beach days or summer events.', 'bohemian-floral-maxi-dress', 69.99, 89.99, 22,
    10, 'women', 87, 'Boho Chic',
    '["Lightweight Rayon Fabric","Adjustable Straps","Smocked Bodice","Side Pockets","Machine Washable"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 412
  ),
(
    20, 'normal', 'Merino Wool Knit Beanie Hat', 'Premium beanie hat knitted from 100% merino wool. Naturally temperature-regulating, moisture-wicking, and odor-resistant. Features a modern ribbed design and comes in a variety of colors. One size fits most.', 'merino-wool-knit-beanie', 29.99, 39.99, 25,
    3, 'fashion', 98, 'Nordic Essentials',
    '["100% Merino Wool","Temperature Regulating","Moisture-Wicking","Odor-Resistant","Soft Non-Itch Fabric"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 276
  ),
(
    21, 'normal', 'Professional Noise-Cancelling Wireless Earbuds', 'Premium wireless earbuds featuring advanced active noise cancellation, high-resolution audio, and customizable touch controls. Includes wireless charging case, multiple ear tip sizes, and up to 32 hours of battery life.', 'professional-noise-cancelling-earbuds', 149.99, 199.99, 25,
    4, 'electronics', 67, 'AudioPro',
    '["Active Noise Cancellation","Hi-Res Audio Certification","IPX7 Waterproof","Touch Controls","32-Hour Battery Life"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 423
  ),
(
    22, 'normal', 'Ultra HD 4K Smart TV (55-inch)', 'Premium 55-inch 4K Ultra HD Smart LED TV with QLED technology, HDR10+, and quantum processor for incredible picture quality. Features voice assistant compatibility, streaming apps, and wireless connectivity.', 'ultra-hd-4k-smart-tv-55inch', 799.99, 999.99, 20,
    4, 'electronics', 32, 'VisionTech',
    '["QLED Technology","4K UHD Resolution","HDR10+ Support","Voice Assistant Compatible","Bluetooth Audio"]', '{"courier":"3-5 days delivery, free shipping","local":"7-10 days delivery, $29.00","express":"2-day delivery, $49.99"}', '{}', 287
  ),
(
    23, 'normal', 'Waterproof Bluetooth Speaker with 360° Sound', 'Rugged waterproof Bluetooth speaker with 360° immersive sound, vibrant LED light show, and 24-hour battery life. Features built-in microphone for calls, USB-C fast charging, and floating design for pool use.', 'waterproof-bluetooth-speaker-360', 99.99, 129.99, 23,
    4, 'electronics', 108, 'SoundWave',
    '["IPX7 Waterproof","360° Sound Technology","24-Hour Battery","LED Light Show","USB-C Fast Charging"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 564
  ),
(
    24, 'normal', 'Professional Mirrorless Digital Camera with 4K Video', 'Advanced mirrorless camera featuring 24.2MP APS-C sensor, 4K/60fps video recording, and 5-axis image stabilization. Includes flip-out touchscreen, eye-tracking autofocus, and weather-sealed body.', 'professional-mirrorless-camera-4k', 1299.99, 1499.99, 13,
    4, 'electronics', 24, 'OptiView',
    '["24.2MP APS-C Sensor","4K/60fps Video","5-Axis Stabilization","Eye-Tracking AF","Weather-Sealed Body"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 189
  ),
(
    25, 'normal', 'Advanced Fitness Tracking Smartwatch', 'Premium fitness smartwatch with always-on AMOLED display, comprehensive health monitoring, and built-in GPS. Features 14-day battery life, 30+ sport modes, blood oxygen monitoring, and sleep tracking.', 'advanced-fitness-smartwatch', 229.99, 279.99, 18,
    4, 'electronics', 75, 'FitTech',
    '["AMOLED Display","14-Day Battery Life","24/7 Heart Monitoring","Blood Oxygen Sensor","Built-in GPS"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 378
  ),
(
    26, 'limited', 'STEM Educational Building Blocks Set (250 Pieces)', 'Comprehensive educational building set designed to develop STEM skills through creative play. Includes 250 interlocking pieces, instruction booklet with 20 project ideas, and storage container. Compatible with major building block brands.', 'stem-educational-building-blocks', 39.99, 49.99, 20,
    5, 'toys', 62, 'BrainBuilder',
    '["250 Interlocking Pieces","STEM Learning Focus","20 Project Ideas","Compatible with Major Brands","Includes Storage Container"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 345
  ),
(
    27, 'trending', 'High-Performance Remote Control Racing Car', 'Professional-grade RC racing car with brushless motor, 4WD system, and 40km/h top speed. Features proportional steering, adjustable suspension, and long-range 2.4GHz controller. Includes rechargeable battery and charger.', 'high-performance-rc-racing-car', 129.99, 159.99, 19,
    5, 'toys', 43, 'SpeedMaster',
    '["Brushless Motor","40km/h Top Speed","4WD System","Adjustable Suspension","2.4GHz Long-Range Control"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 267
  ),
(
    28, 'trending', 'Handcrafted Premium Teddy Bear (18-inch)', 'Luxurious handcrafted teddy bear made from ultra-soft premium plush material. Features safety-tested glass eyes, articulated limbs, and embroidered details. Comes with birth certificate and gift box.', 'handcrafted-premium-teddy-bear', 49.99, 64.99, 23,
    5, 'toys', 89, 'Cuddle Creations',
    '["Handcrafted Quality","Premium Plush Material","Safety Tested","Articulated Limbs","Gift Box Included"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 423
  ),
(
    29, 'trending', 'Family Strategy Board Game Collection', 'Comprehensive family board game featuring 5 different strategy games on a double-sided board. Includes chess, checkers, backgammon, ludo, and snakes & ladders. Made from sustainable wood with premium game pieces.', 'family-strategy-board-game', 54.99, 69.99, 21,
    5, 'toys', 37, 'GameMaster',
    '["5 Games in 1","Sustainable Wood Construction","Premium Game Pieces","Educational Value","Family-Friendly"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 312
  ),
(
    30, 'normal', 'Deluxe Art and Craft Supply Kit (150+ Pieces)', 'Comprehensive art and craft kit for children featuring 150+ pieces including colored pencils, markers, watercolors, oil pastels, and various craft supplies. Comes in a sturdy wooden storage case with organization compartments.', 'deluxe-art-craft-supply-kit', 59.99, 79.99, 25,
    5, 'toys', 54, 'Creative Kids',
    '["150+ Art Supplies","Wooden Storage Case","Non-Toxic Materials","Organization Compartments","Suitable for Ages 6+"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 287
  ),
(
    31, 'trending', 'Premium Single-Origin Coffee Beans (1lb)', 'Specialty-grade arabica coffee beans from the Ethiopian highlands, medium roasted to perfection. Features bright citrus notes, floral aroma, and a smooth chocolate finish. Ethically sourced and roasted in small batches.', 'premium-single-origin-coffee', 19.99, 24.99, 20,
    1, 'grocery', 87, 'Artisan Roasters',
    '["Single-Origin Ethiopian","Medium Roast","Specialty Grade","Fair Trade Certified","Small-Batch Roasted"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 456
  ),
(
    32, 'normal', 'Gourmet Belgian Chocolate Chip Cookies (Box of 12)', 'Handcrafted gourmet cookies made with premium Belgian chocolate chunks, organic butter, and Madagascar vanilla. Perfectly balanced crisp edge and chewy center. No preservatives or artificial ingredients.', 'gourmet-chocolate-chip-cookies', 14.99, 18.99, 21,
    1, 'grocery', 124, 'Sweet Indulgence',
    '["Premium Belgian Chocolate","Organic Ingredients","No Preservatives","Handcrafted Daily","Sustainable Packaging"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 378
  ),
(
    33, 'normal', 'Cold-Pressed Organic Orange Juice (32oz)', 'Fresh cold-pressed orange juice made from 100% organic, sun-ripened oranges. Unpasteurized to preserve maximum nutrients and flavor. Contains no added sugars, preservatives, or artificial ingredients.', 'cold-pressed-organic-orange-juice', 8.99, 10.99, 18,
    1, 'grocery', 92, 'Pure Harvest',
    '["Cold-Pressed","100% Organic","No Added Sugar","Rich in Vitamin C","Glass Bottle Packaging"]', '{"courier":"Same-day cold delivery, $3.99","local":"Next-day cold delivery, $6.00","pickup":"Free store pickup"}', '{}', 543
  ),
(
    34, 'trending', 'Advanced Anti-Aging Facial Serum (50ml)', 'Professional-grade facial serum featuring hyaluronic acid, vitamin C, retinol, and peptides. Targets fine lines, uneven skin tone, and loss of firmness. Dermatologist-tested and suitable for all skin types.', 'advanced-anti-aging-facial-serum', 79.99, 99.99, 20,
    1, 'grocery', 65, 'DermaPerfect',
    '["Hyaluronic Acid","Vitamin C Complex","Peptide Formula","Dermatologist Tested","Paraben-Free"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 298
  ),
(
    35, 'trending', 'Handcrafted Ceramic Minimalist Vase Set', 'Set of 3 contemporary handcrafted ceramic vases in complementary sizes and shapes. Each piece features a unique minimalist design with matte glaze finish. Perfect for modern home décor or as a thoughtful gift.', 'handcrafted-ceramic-vase-set', 79.99, 99.99, 20,
    2, 'household-appliances', 48, 'Modern Artisan',
    '["Handcrafted Ceramic","Set of 3 Vases","Unique Designs","Matte Glaze Finish","Fair Trade Certified"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 324
  ),
(
    36, 'trending', 'Premium Multisport GPS Smartwatch', 'Advanced multisport GPS smartwatch with 1.4" always-on display, titanium bezel, and sapphire crystal. Features comprehensive fitness tracking, advanced training metrics, onboard maps, and up to 21 days of battery life.', 'premium-multisport-gps-smartwatch', 549.99, 649.99, 15,
    4, 'electronics', 34, 'AthleteTime',
    '["GPS Navigation","Heart Rate Monitor","Pulse Oximeter","21-Day Battery Life","100m Water Resistance"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 267
  ),
(
    37, 'trending', 'Premium Stretch Selvedge Denim Jeans', 'Premium men''s jeans crafted from Japanese selvedge denim with 2% stretch for comfort. Features classic five-pocket design, copper rivets, and chain-stitched hem. Medium-dark wash with subtle fading for a timeless look.', 'premium-stretch-selvedge-jeans', 149.99, 189.99, 21,
    7, 'men', 57, 'Heritage Denim',
    '["Japanese Selvedge Denim","2% Stretch Comfort","Copper Hardware","Chain-Stitched Details","Medium-Dark Wash"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 376
  ),
(
    38, 'trending', 'Luxe Merino Wool Crewneck Sweater', 'Premium men''s sweater knitted from 100% extra-fine merino wool. Features classic crewneck design, ribbed cuffs and hem, and relaxed contemporary fit. Exceptionally soft, lightweight, and naturally temperature-regulating.', 'luxe-merino-wool-crewneck-sweater', 129.99, 159.99, 19,
    7, 'men', 73, 'Elevated Essentials',
    '["100% Extra-Fine Merino","Temperature Regulating","Odor-Resistant","Machine Washable","Sustainable Sourcing"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 289
  ),
(
    39, 'trending', 'Vintage-Inspired Wool Fedora Hat', 'Classic men''s fedora handcrafted from premium wool felt with grosgrain ribbon band and satin lining. Features pinched crown, structured brim, and moisture-wicking sweatband. Available in multiple sizes and colors.', 'vintage-inspired-wool-fedora', 89.99, 109.99, 18,
    7, 'men', 42, 'Heritage Headwear',
    '["Premium Wool Felt","Grosgrain Band","Satin Lining","Moisture-Wicking Sweatband","Handcrafted Quality"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 217
  ),
(
    40, 'trending', 'Premium Leather Oxford Dress Shoes', 'Classic men''s oxford shoes handcrafted from premium full-grain calfskin leather. Features Goodyear welt construction, leather soles, and hand-finished patina. Includes cedar shoe trees and luxury dust bags.', 'premium-leather-oxford-shoes', 299.99, 349.99, 14,
    7, 'men', 38, 'Heritage Footwear',
    '["Full-Grain Calfskin","Goodyear Welt Construction","Leather Soles","Hand-Finished Patina","Cedar Shoe Trees Included"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 198
  ),
(
    41, 'limited', 'Samsung Galaxy S24 Ultra (256GB, Titanium)', 'Flagship smartphone featuring a 6.8" Dynamic AMOLED 2X display, Snapdragon 8 Gen 3 processor, and 200MP camera system. Includes S Pen, 12GB RAM, all-day battery, and advanced AI features in a premium titanium frame.', 'samsung-galaxy-s24-ultra', 1299.99, 1399.99, 7,
    4, 'electronics', 25, 'Samsung',
    '[
      "200MP Camera System",
      "Snapdragon 8 Gen 3",
      "6.8\" AMOLED Display",
      "S Pen Included",
      "Titanium Frame"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    187
  ),
(
    42, 'limited', 'Professional Noise-Cancelling True Wireless Earbuds', 'Premium true wireless earbuds featuring adaptive noise cancellation, hi-fi sound with dual drivers, and 360° spatial audio. Includes wireless charging case, customizable touch controls, and multipoint connection.', 'premium-true-wireless-earbuds', 199.99, 249.99, 20,
    4, 'electronics', 43, 'SoundMaster',
    '[
      "Adaptive Noise Cancellation",
      "Dual-Driver System",
      "360° Spatial Audio",
      "36-Hour Battery Life",
      "IPX7 Waterproof"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    263
  ),
(
    43, 'limited', 'Premium Down-Insulated Winter Parka', 'Expedition-grade winter jacket featuring 800-fill power goose down insulation, waterproof breathable shell, and premium fur-lined hood. Includes reinforced shoulders, storm cuffs, and multiple functional pockets.', 'premium-down-insulated-winter-parka', 499.99, 599.99, 17,
    3, 'fashion', 32, 'Alpine Explorer',
    '["800-Fill Goose Down","Waterproof Breathable Shell","Fur-Lined Hood","Reinforced Shoulders","YKK Waterproof Zippers"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 176
  ),
(
    44, 'limited', 'New Balance 574 Legacy Premium Sneakers', 'Iconic sneakers crafted with premium suede and mesh uppers, ENCAP midsole cushioning, and durable rubber outsole. Features retro-inspired colorway, reflective details, and improved comfort footbed.', 'new-balance-574-legacy-premium', 119.99, 149.99, 20,
    3, 'fashion', 47, 'New Balance',
    '["Premium Suede/Mesh Upper","ENCAP Midsole Cushioning","Reflective Details","Improved Comfort Footbed","Durable Rubber Outsole"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 289
  ),
(
    45, 'limited', 'Professional Studio Wireless Headphones', 'Reference-quality studio headphones featuring premium 50mm beryllium drivers, adaptive EQ, and lossless audio support. Includes active noise cancellation, spatial audio, and up to 40 hours of battery life.', 'professional-studio-wireless-headphones', 349.99, 399.99, 13,
    4, 'electronics', 38, 'AudioElite',
    '["50mm Beryllium Drivers","Adaptive EQ Technology","Hi-Res Audio Certified","40-Hour Battery Life","Memory Foam Ear Cushions"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 217
  ),
(
    46, 'limited', 'Premium Smart Home Security Camera System (4-Pack)', 'Comprehensive home security system featuring 4 weather-resistant 2K HDR cameras with color night vision, facial recognition, and two-way audio. Includes AI-powered motion detection, encrypted local storage, and smartphone integration.', 'premium-smart-home-security-system', 499.99, 599.99, 17,
    4, 'electronics', 24, 'SecureHome',
    '["2K HDR Resolution","Color Night Vision","Facial Recognition","AI Motion Detection","Local Encrypted Storage"]', '{"courier":"2-3 days delivery, free shipping","local":"5-7 days delivery, $9.00","express":"Next-day delivery, $14.99"}', '{}', 156
  ),
(
    47, 'limited', 'Ultra-Slim Convertible Touchscreen Laptop', 'Premium 2-in-1 laptop featuring 14" 4K OLED touchscreen, Intel Core i7 processor, 16GB RAM, and 1TB SSD. Includes 360° hinge, precision touchpad, backlit keyboard, and all-day battery life in a sleek aluminum chassis.', 'ultra-slim-convertible-laptop', 1499.99, 1699.99, 12,
    6, 'laptops', 19, 'TechElite',
    '[
      "14\" 4K OLED Touchscreen",
      "Intel Core i7 Processor",
      "16GB RAM/1TB SSD",
      "360° Convertible Design",
      "Aluminum Unibody Chassis"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    134
  ),
(
    48, 'top-sale', 'Premium Ergonomic Computer Accessory Bundle', 'Complete ergonomic workspace solution featuring mechanical keyboard with palm rest, precision mouse, dual monitor arms, and premium headset. Designed for maximum comfort and productivity during extended work sessions.', 'premium-ergonomic-computer-accessory-bundle', 299.99, 399.99, 25,
    4, 'electronics', 43, 'ErgoTech',
    '[
      "Mechanical Keyboard",
      "Ergonomic Mouse",
      "Dual Monitor Arms",
      "Premium Headset",
      "Cable Management System"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    245
  ),
(
    49, 'top-sale', 'Mens Premium Business Casual Wardrobe Set', 'Curated mens wardrobe essentials including stretch chino pants, oxford shirt, merino wool sweater, and leather belt. Made from premium materials in versatile colors perfect for business casual environments.', 'mens-premium-business-casual-set', 329.99, 399.99, 18,
    7, 'men', 32, 'Executive Style',
    '[
      "Stretch Chino Pants",
      "Oxford Cotton Shirt",
      "Merino Wool Sweater",
      "Italian Leather Belt",
      "Travel-Friendly Fabrics"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    187
  ),
(
    50, 'top-sale', 'Cold-Pressed Pomegranate Juice (6-Pack)', 'Premium cold-pressed pomegranate juice made from hand-selected organic fruit. Rich in antioxidants with no added sugars, preservatives, or concentrates. Packaged in sustainable glass bottles, each containing 16oz of pure juice.', 'cold-pressed-pomegranate-juice', 59.99, 69.99, 14,
    1, 'grocery', 67, 'Orchard Fresh',
    '[
      "100% Pure Pomegranate",
      "Cold-Pressed",
      "No Added Sugar",
      "Rich in Antioxidants",
      "Glass Bottle Packaging"
    ]', 
    '{
      "courier": "Same-day cold delivery, $3.99",
      "local": "Next-day cold delivery, $6.00",
      "pickup": "Free store pickup"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    312
  ),
(
    51, 'top-sale', 'Premium Organic Dog Food Subscription Box', 'Customized monthly subscription of premium organic dog food formulated by veterinary nutritionists. Features human-grade ingredients, balanced nutrients, and options for special dietary needs. Available in various portion sizes.', 'premium-organic-dog-food-subscription', 79.99, 99.99, 20,
    1, 'grocery', 87, 'Healthy Paws',
    '[
      "Human-Grade Ingredients",
      "Vet-Formulated",
      "Organic Certified",
      "Customized Portions",
      "Monthly Delivery"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "col",
      "isComingSoon": false
    }', 
    276
  ),
(
    52, 'top-sale', 'Ultra HD Security Camera System with AI (8-Channel)', 'Professional-grade security system featuring 8-channel NVR with 4TB storage and 6 ultra HD cameras with color night vision. Includes AI-powered person/vehicle detection, smart alerts, and remote viewing capabilities.', 'ultra-hd-security-camera-system', 899.99, 1099.99, 18,
    4, 'electronics', 23, 'SecureVision',
    '[
      "8-Channel NVR",
      "6 Ultra HD Cameras",
      "Color Night Vision",
      "AI Person Detection",
      "4TB Storage Included"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $29.00",
      "express": "Next-day delivery, $49.99"
    }', 
    '{
      "setUpDesign": "double",
      "isComingSoon": false
    }', 
    165
  ),
(
    53, 'top-sale', 'Luxury Skincare Collection Gift Set', 'Comprehensive luxury skincare collection featuring cleanser, toner, serum, moisturizer, and eye cream. Formulated with premium ingredients including hyaluronic acid, peptides, and botanical extracts in elegant packaging.', 'luxury-skincare-collection-gift-set', 249.99, 299.99, 17,
    1, 'grocery', 41, 'DermaLuxe',
    '[
      "5-Step Skincare System",
      "Premium Formulations",
      "Anti-Aging Benefits",
      "Elegant Gift Packaging",
      "Suitable for All Skin Types"
    ]', 
    '{
      "courier": "2-3 days delivery, free shipping",
      "local": "5-7 days delivery, $9.00",
      "express": "Next-day delivery, $14.99"
    }', 
    '{
      "setUpDesign": "row",
      "isComingSoon": false
    }', 
    198
  );


-- Add product images
INSERT INTO product_images (id, product_id, url, is_primary, display_order) VALUES
-- Apple MacBook Pro (ID: 1)
(1, 1, '/images/products/electronics/macbook-pro-m1-gray-front.jpg', true, 0),
(2, 1, '/images/products/electronics/macbook-pro-m1-gray-angle.jpg', false, 1),
(3, 1, '/images/products/electronics/macbook-pro-m1-gray-open.jpg', false, 2),
(4, 1, '/images/products/electronics/macbook-pro-m1-gray-side.jpg', false, 3),
(5, 1, '/images/products/electronics/macbook-pro-m1-gray-keyboard.jpg', false, 4),
(6, 1, '/images/products/electronics/macbook-pro-m1-gray-ports.jpg', false, 5),
(7, 1, '/images/products/electronics/macbook-pro-m1-gray-display.jpg', false, 6),
(8, 1, '/images/products/electronics/macbook-pro-m1-gray-lifestyle.jpg', false, 7),

-- Organic Bananas (ID: 2)
(9, 2, '/images/products/fruits/organic-bananas-bunch.jpg', true, 0),
(10, 2, '/images/products/fruits/organic-bananas-single.jpg', false, 1),
(11, 2, '/images/products/fruits/organic-bananas-peeled.jpg', false, 2),
(12, 2, '/images/products/fruits/organic-bananas-bunch-angle.jpg', false, 3),
(13, 2, '/images/products/fruits/organic-bananas-lifestyle.jpg', false, 4),

-- Red Delicious Apples (ID: 3)
(14, 3, '/images/products/fruits/red-delicious-apples-pack.jpg', true, 0),
(15, 3, '/images/products/fruits/red-delicious-apples-single.jpg', false, 1),
(16, 3, '/images/products/fruits/red-delicious-apples-sliced.jpg', false, 2),
(17, 3, '/images/products/fruits/red-delicious-apples-basket.jpg', false, 3),
(18, 3, '/images/products/fruits/red-delicious-apples-tree.jpg', false, 4),
(19, 3, '/images/products/fruits/red-delicious-apples-lifestyle.jpg', false, 5),

-- Organic Berries Mix (ID: 4)
(20, 4, '/images/products/fruits/berries-mix-package.jpg', true, 0),
(21, 4, '/images/products/fruits/berries-mix-open.jpg', false, 1),
(22, 4, '/images/products/fruits/berries-mix-bowl.jpg', false, 2),
(23, 4, '/images/products/fruits/berries-mix-closeup.jpg', false, 3),
(24, 4, '/images/products/fruits/berries-mix-lifestyle.jpg', false, 4),

-- Artisan Sourdough Bread (ID: 5)
(25, 5, '/images/products/grocery/sourdough-bread-whole.jpg', true, 0),
(26, 5, '/images/products/grocery/sourdough-bread-sliced.jpg', false, 1),
(27, 5, '/images/products/grocery/sourdough-bread-crust.jpg', false, 2),
(28, 5, '/images/products/grocery/sourdough-bread-texture.jpg', false, 3),
(29, 5, '/images/products/grocery/sourdough-bread-lifestyle.jpg', false, 4),

-- Free-Range Organic Eggs (ID: 6)
(30, 6, '/images/products/grocery/organic-eggs-carton.jpg', true, 0),
(31, 6, '/images/products/grocery/organic-eggs-open.jpg', false, 1),
(32, 6, '/images/products/grocery/organic-eggs-single.jpg', false, 2),
(33, 6, '/images/products/grocery/organic-eggs-cracked.jpg', false, 3),
(34, 6, '/images/products/grocery/organic-eggs-farm.jpg', false, 4),

-- Grass-Fed Organic Milk (ID: 7)
(35, 7, '/images/products/grocery/organic-milk-bottle.jpg', true, 0),
(36, 7, '/images/products/grocery/organic-milk-pour.jpg', false, 1),
(37, 7, '/images/products/grocery/organic-milk-glass.jpg', false, 2),
(38, 7, '/images/products/grocery/organic-milk-cream-top.jpg', false, 3),
(39, 7, '/images/products/grocery/organic-milk-farm.jpg', false, 4),

-- Organic Hass Avocados (ID: 8)
(40, 8, '/images/products/fruits/hass-avocados-pack.jpg', true, 0),
(41, 8, '/images/products/fruits/hass-avocados-single.jpg', false, 1),
(42, 8, '/images/products/fruits/hass-avocados-halved.jpg', false, 2),
(43, 8, '/images/products/fruits/hass-avocados-sliced.jpg', false, 3),
(44, 8, '/images/products/fruits/hass-avocados-lifestyle.jpg', false, 4),

-- Tropical Mangoes (ID: 9)
(45, 9, '/images/products/fruits/tropical-mangoes-pack.jpg', true, 0),
(46, 9, '/images/products/fruits/tropical-mangoes-single.jpg', false, 1),
(47, 9, '/images/products/fruits/tropical-mangoes-sliced.jpg', false, 2),
(48, 9, '/images/products/fruits/tropical-mangoes-peeled.jpg', false, 3),
(49, 9, '/images/products/fruits/tropical-mangoes-lifestyle.jpg', false, 4),

-- Luxury Scented Candle Set (ID: 10)
(50, 10, '/images/products/household-appliances/scented-candle-set-package.jpg', true, 0),
(51, 12, '/images/products/household-appliances/egyptian-cotton-sheets-package.jpg', true, 0),
(52, 13, '/images/products/household-appliances/kitchen-utensil-set-complete.jpg', true, 0),
(53, 14, '/images/products/household-appliances/throw-pillow-collection-set.jpg', true, 0),
(54, 15, '/images/products/household-appliances/led-desk-lamp-front.jpg', true, 0),

-- Clothing & Fashion Products (ID: 16-20)
(55, 16, '/images/products/fashion/womens-denim-jacket-front.jpg', true, 0),
(56, 17, '/images/products/fashion/mens-graphic-tshirt-front.jpg', true, 0),
(57, 18, '/images/products/fashion/leather-sneakers-side.jpg', true, 0),
(58, 19, '/images/products/fashion/floral-maxi-dress-front.jpg', true, 0),
(59, 20, '/images/products/fashion/merino-wool-beanie-front.jpg', true, 0),

-- Electronics Products (ID: 21-25)
(60, 21, '/images/products/electronics/wireless-earbuds-case.jpg', true, 0),
(61, 22, '/images/products/electronics/4k-smart-tv-front.jpg', true, 0),
(62, 23, '/images/products/electronics/bluetooth-speaker-angle.jpg', true, 0),
(63, 24, '/images/products/electronics/mirrorless-camera-front.jpg', true, 0),
(64, 25, '/images/products/electronics/fitness-smartwatch-front.jpg', true, 0),

-- Toys & Games (ID: 26-30)
(65, 26, '/images/products/toys/building-blocks-set-complete.jpg', true, 0),
(66, 27, '/images/products/toys/rc-racing-car-angle.jpg', true, 0),
(67, 28, '/images/products/toys/teddy-bear-front.jpg', true, 0),
(68, 29, '/images/products/toys/strategy-board-game-package.jpg', true, 0),
(69, 30, '/images/products/toys/art-craft-kit-open.jpg', true, 0),

-- Food & Beverage Products (ID: 31-33)
(70, 31, '/images/products/grocery/coffee-beans-package.jpg', true, 0),
(71, 32, '/images/products/grocery/chocolate-chip-cookies-box.jpg', true, 0),
(72, 33, '/images/products/grocery/orange-juice-bottle.jpg', true, 0),

-- Beauty & Skincare (ID: 34)
(73, 34, '/images/products/beauty/facial-serum-bottle.jpg', true, 0),

-- Home Decor (ID: 35)
(74, 35, '/images/products/grocery/ceramic-vase-set-arranged.jpg', true, 0),

-- More Electronics (ID: 36)
(75, 36, '/images/products/electronics/gps-smartwatch-front.jpg', true, 0),

-- Men's Fashion (ID: 37-40)
(76, 37, '/images/products/fashion/selvedge-jeans-front.jpg', true, 0),
(77, 38, '/images/products/fashion/merino-wool-sweater-front.jpg', true, 0),
(78, 39, '/images/products/fashion/wool-fedora-hat-angle.jpg', true, 0),
(79, 40, '/images/products/fashion/oxford-dress-shoes-side.jpg', true, 0),

-- Premium Electronics (ID: 41-42, 45-47)
(80, 41, '/images/products/electronics/galaxy-s24-front.jpg', true, 0),
(81, 41, '/images/products/electronics/galaxy-s24-back.jpg', false, 1),
(82, 41, '/images/products/electronics/galaxy-s24-angle.jpg', false, 2),
(83, 41, '/images/products/electronics/galaxy-s24-lifestyle.jpg', false, 3),
(84, 41, '/images/products/electronics/galaxy-s24-camera.jpg', false, 4),
(85, 42, '/images/products/electronics/noise-cancelling-earbuds-case.jpg', true, 0),
(86, 42, '/images/products/electronics/noise-cancelling-earbuds-single.jpg', false, 1),
(87, 42, '/images/products/electronics/noise-cancelling-earbuds-wearing.jpg', false, 2),
(88, 42, '/images/products/electronics/noise-cancelling-earbuds-charging.jpg', false, 3),
(89, 42, '/images/products/electronics/noise-cancelling-earbuds-lifestyle.jpg', false, 4),
(90, 45, '/images/products/electronics/studio-headphones-front.jpg', true, 0),
(91, 45, '/images/products/electronics/studio-headphones-side.jpg', false, 1),
(92, 45, '/images/products/electronics/studio-headphones-folded.jpg', false, 2),
(93, 45, '/images/products/electronics/studio-headphones-detail.jpg', false, 3),
(94, 45, '/images/products/electronics/studio-headphones-lifestyle.jpg', false, 4),
(95, 46, '/images/products/electronics/security-camera-system-complete.jpg', true, 0),
(96, 46, '/images/products/electronics/security-camera-single.jpg', false, 1),
(97, 46, '/images/products/electronics/security-camera-nvr.jpg', false, 2),
(98, 46, '/images/products/electronics/security-camera-install.jpg', false, 3),
(99, 46, '/images/products/electronics/security-camera-app.jpg', false, 4),
(100, 47, '/images/products/electronics/convertible-laptop-front.jpg', true, 0),

-- Fashion Products (ID: 43-44)
(101, 43, '/images/products/fashion/winter-parka-front.jpg', true, 0),
(102, 43, '/images/products/fashion/winter-parka-back.jpg', false, 1),
(103, 43, '/images/products/fashion/winter-parka-hood.jpg', false, 2),
(104, 43, '/images/products/fashion/winter-parka-details.jpg', false, 3),
(105, 43, '/images/products/fashion/winter-parka-lifestyle.jpg', false, 4),
(106, 44, '/images/products/fashion/nb-574-sneakers-side.jpg', true, 0),

-- Top Selling Products (ID: 48-53)
(107, 48, '/images/products/electronics/computer-accessory-bundle-complete.jpg', true, 0),
(108, 49, '/images/products/fashion/business-casual-set-complete.jpg', true, 0),
(109, 50, '/images/products/grocery/pomegranate-juice-bottles.jpg', true, 0),
(110, 50, '/images/products/grocery/pomegranate-juice-single.jpg', false, 1),
(111, 50, '/images/products/grocery/pomegranate-juice-pour.jpg', false, 2),
(112, 50, '/images/products/grocery/pomegranate-juice-ingredients.jpg', false, 3),
(113, 51, '/images/products/pets/dog-food-subscription-box.jpg', true, 0),
(114, 51, '/images/products/pets/dog-food-package.jpg', false, 1),
(115, 51, '/images/products/pets/dog-food-serving.jpg', false, 2),
(116, 51, '/images/products/pets/dog-food-ingredients.jpg', false, 3),
(117, 52, '/images/products/electronics/security-system-complete.jpg', true, 0),
(118, 52, '/images/products/electronics/security-system-camera.jpg', false, 1),
(119, 52, '/images/products/electronics/security-system-monitor.jpg', false, 2),
(120, 52, '/images/products/electronics/security-system-app.jpg', false, 3),
(121, 53, '/images/products/beauty/skincare-collection-complete.jpg', true, 0),
(122, 53, '/images/products/beauty/skincare-collection-cleaner.jpg', false, 1),
(123, 53, '/images/products/beauty/skincare-collection-serum.jpg', false, 2),
(124, 53, '/images/products/beauty/skincare-collection-moisturizer.jpg', false, 3);



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
(10, 'household-appliances'),
(12, 'household-appliances'),
(13, 'household-appliances'),
(14, 'household-appliances'),
(15, 'household-appliances'),
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
(1, 1, '123e4567-e89b-12d3-a456-426614174000', 'John Smith', 5, 'Best MacBook I''ve ever had! The performance is amazing.', CURRENT_TIMESTAMP - INTERVAL '3 days'),
(2, 1, '123e4567-e89b-12d3-a456-426614174001', 'Emily Johnson', 4, 'Great laptop overall, but battery life could be better.', CURRENT_TIMESTAMP - INTERVAL '7 days'),
(3, 1, '123e4567-e89b-12d3-a456-426614174002', 'David Williams', 5, 'Switched from Windows and never looking back. Fast and reliable!', CURRENT_TIMESTAMP - INTERVAL '15 days'),
(4, 1, '123e4567-e89b-12d3-a456-426614174003', 'Sarah Davis', 4, 'Perfect for my design work. Love the display quality.', CURRENT_TIMESTAMP - INTERVAL '21 days'),
(5, 1, '123e4567-e89b-12d3-a456-426614174004', 'Michael Brown', 5, 'So lightweight yet powerful. Worth every penny.', CURRENT_TIMESTAMP - INTERVAL '30 days'),

-- Reviews for Samsung 4K TV
(6, 22, '123e4567-e89b-12d3-a456-426614174000', 'John Smith', 5, 'Picture quality is outstanding! The colors are vibrant and lifelike.', CURRENT_TIMESTAMP - INTERVAL '5 days'),
(7, 22, '123e4567-e89b-12d3-a456-426614174001', 'Emily Johnson', 4, 'Great TV for the price. Smart features work smoothly.', CURRENT_TIMESTAMP - INTERVAL '8 days'),
(8, 22, '123e4567-e89b-12d3-a456-426614174002', 'David Williams', 5, 'Perfect for gaming with low input lag.', CURRENT_TIMESTAMP - INTERVAL '12 days'),
(9, 22, '123e4567-e89b-12d3-a456-426614174003', 'Sarah Davis', 3, 'Good TV but the built-in speakers could be better.', CURRENT_TIMESTAMP - INTERVAL '18 days'),

-- Reviews for Wireless Earbuds
(10, 21, '123e4567-e89b-12d3-a456-426614174001', 'Emily Johnson', 5, 'The sound quality is amazing and battery life is impressive!', CURRENT_TIMESTAMP - INTERVAL '4 days'),
(11, 21, '123e4567-e89b-12d3-a456-426614174002', 'David Williams', 4, 'Comfortable fit, great for workouts.', CURRENT_TIMESTAMP - INTERVAL '10 days'),
(12, 21, '123e4567-e89b-12d3-a456-426614174004', 'Michael Brown', 3, 'Good for casual listening but not for audiophiles.', CURRENT_TIMESTAMP - INTERVAL '22 days'),

-- Reviews for Women's Denim Jacket
(13, 16, '123e4567-e89b-12d3-a456-426614174005', 'Jennifer Lee', 5, 'Perfect fit and looks exactly like the pictures. Very stylish!', CURRENT_TIMESTAMP - INTERVAL '2 days'),
(14, 16, '123e4567-e89b-12d3-a456-426614174006', 'Amanda Wilson', 5, 'The material quality is excellent. Highly recommend!', CURRENT_TIMESTAMP - INTERVAL '9 days'),
(15, 16, '123e4567-e89b-12d3-a456-426614174007', 'Sophia Martinez', 4, 'Runs slightly small but overall great quality.', CURRENT_TIMESTAMP - INTERVAL '14 days'),

-- Reviews for Smartwatch
(16, 25, '123e4567-e89b-12d3-a456-426614174003', 'Sarah Davis', 5, 'Love all the fitness tracking features. Battery lasts for days!', CURRENT_TIMESTAMP - INTERVAL '6 days'),
(17, 25, '123e4567-e89b-12d3-a456-426614174004', 'Michael Brown', 4, 'Great smartwatch with accurate heart rate monitoring.', CURRENT_TIMESTAMP - INTERVAL '16 days'),
(18, 25, '123e4567-e89b-12d3-a456-426614174005', 'Jennifer Lee', 5, 'The sleep tracking feature is a game changer for me.', CURRENT_TIMESTAMP - INTERVAL '25 days'),
(19, 25, '123e4567-e89b-12d3-a456-426614174006', 'Amanda Wilson', 3, 'Good functionality but the app could be more intuitive.', CURRENT_TIMESTAMP - INTERVAL '29 days'),

-- Reviews for Men's Casual Outfit
(20, 49, '123e4567-e89b-12d3-a456-426614174000', 'John Smith', 5, 'Very comfortable and looks great. Perfect for casual outings.', CURRENT_TIMESTAMP - INTERVAL '1 day'),
(21, 49, '123e4567-e89b-12d3-a456-426614174001', 'Emily Johnson', 4, 'Good quality fabric, slight issue with sizing.', CURRENT_TIMESTAMP - INTERVAL '11 days'),
(22, 49, '123e4567-e89b-12d3-a456-426614174007', 'Sophia Martinez', 5, 'Bought this for my husband and he loves it!', CURRENT_TIMESTAMP - INTERVAL '20 days'),

-- Reviews for Cotton Bed Sheets
(23, 12, '123e4567-e89b-12d3-a456-426614174005', 'Jennifer Lee', 5, 'So soft and comfortable! Love the quality.', CURRENT_TIMESTAMP - INTERVAL '5 days'),
(24, 12, '123e4567-e89b-12d3-a456-426614174006', 'Amanda Wilson', 5, 'The best sheets Ive ever purchased. Will buy again!', CURRENT_TIMESTAMP - INTERVAL '13 days'),
(25, 12, '123e4567-e89b-12d3-a456-426614174007', 'Sophia Martinez', 4, 'Great quality cotton, washes well without shrinking.', CURRENT_TIMESTAMP - INTERVAL '23 days'),

-- Reviews for Computer Accessories
(26, 48, '123e4567-e89b-12d3-a456-426614174002', 'David Williams', 5, 'This keyboard is amazing! Great tactile feedback and RGB lighting.', CURRENT_TIMESTAMP - INTERVAL '7 days'),
(27, 48, '123e4567-e89b-12d3-a456-426614174003', 'Sarah Davis', 4, 'Good gaming mouse with customizable buttons.', CURRENT_TIMESTAMP - INTERVAL '17 days'),
(28, 48, '123e4567-e89b-12d3-a456-426614174004', 'Michael Brown', 5, 'The mousepad is excellent quality and perfect size for gaming.', CURRENT_TIMESTAMP - INTERVAL '24 days'),
(29, 48, '123e4567-e89b-12d3-a456-426614174000', 'John Smith', 4, 'Solid build quality on all accessories. Recommended for gamers.', CURRENT_TIMESTAMP - INTERVAL '28 days');


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
    '/images/ads/banner/modern-furniture.jpg',
    '/category/electronics',
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
    '/images/ads/banner/worker-environment.jpg',
    '/category/grocery',
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
    '/images/ads/banner/tech-innovations.jpg',
    '/category/women',
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
    '/images/ads/banner/fashion-flash-sale.jpg',
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
    '/images/ads/banner/home-essentials.jpg',
    '/category/household-appliances',
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
    'Fresh Organic Produce', 
    'Farm to Table',
    'Enjoy freshly harvested organic fruits and vegetables',
    'Free Delivery on Orders $50+',
    'Organic Certified',
    '/images/ads/banner/organic-produce.jpg',
    '/category/fruits',
    'Shop Fresh',
    '#386641',
    '#f2e8cf',
    'fade',
    true,
    6,
    'hero',
    3, -- Red Delicious Apples product
    8 -- Fruits category
),
(
    'Premium Audio Experience', 
    'Hear Every Detail',
    'Immerse yourself in crystal-clear sound with our premium audio collection',
    'Up to 40% Off Select Items',
    'Limited Edition',
    '/images/ads/banner/premium-audio.jpg',
    '/category/electronics',
    'Discover More',
    '#14213d',
    '#e5e5e5',
    'slide',
    true,
    7,
    'promotional',
    23, -- Bluetooth Speaker product
    4 -- Electronics category
),
(
    'Seasonal Kitchen Essentials', 
    'Cook Like a Pro',
    'Upgrade your kitchen with premium tools and appliances',
    'Bundle Deals Available',
    'Chef Approved',
    '/images/ads/banner/kitchen-essentials.jpg',
    '/category/household-appliances',
    'Shop Collection',
    '#6d6875',
    '#ffcdb2',
    'zoom',
    true,
    8,
    'hero',
    13, -- Kitchen Utensil Set product
    2 -- Household Appliances category
),
(
    'Back to School Sale', 
    'Be Prepared',
    'Start the school year right with our collection of essentials',
    'Student Discount: 15% Off',
    'Limited Time',
    '/images/ads/banner/back-to-school.jpg',
    '/category/toys',
    'Shop Now',
    '#1b4965',
    '#cae9ff',
    'fade',
    true,
    9,
    'promotional',
    30, -- Art and Craft Supply Kit product
    5 -- Toys category
),
(
    'Winter Wardrobe Essentials', 
    'Stay Cozy & Stylish',
    'Discover our collection of premium winter clothing and accessories',
    'Up to 35% Off',
    'New Collection',
    '/images/ads/banner/winter-wardrobe.jpg',
    '/category/fashion',
    'Explore Now',
    '#354f52',
    '#cad2c5',
    'slide',
    true,
    10,
    'hero',
    43, -- Winter Parka product
    3 -- Fashion category
),
(
    'Smart Home Revolution', 
    'Connect Your Life',
    'Upgrade your home with the latest smart technology and devices',
    'Bundle & Save 25%',
    'Tech Innovation',
    '/images/ads/banner/smart-home.jpg',
    '/category/electronics',
    'Discover Smart Home',
    '#001219',
    '#94d2bd',
    'zoom',
    true,
    11,
    'promotional',
    46, -- Smart Home Security Camera System product
    4 -- Electronics category
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
    '/images/ads/banner/trousers_fashion.png',
    '{"has_more": false}',
    true
),
-- Banner 2: Watchmen fashion with gray Shop Now button
(
    'banner', 'category', 7, 2, -- Using Men category (ID=7)
    '/images/ads/banner/watchmen_fashion.png',
    '{"button_type": "gray", "button_text": "Shop Now", "has_more": false}',
    true
),
-- Banner 3: Denim fashion with "see more" option
(
    'banner', 'category', 10, 3, -- Using Women category (ID=10)
    '/images/ads/banner/denim_fashion.png',
    '{"has_more": true}',
    true
),
-- Banner 4: Domestic item with black Shop Now button
(
    'banner', 'category', 2, 4, -- Using Home category (ID=2)
    '/images/ads/banner/dometic.png',
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
    '/images/ads/display/be-winner.png',
    '{}',
    true
),
-- Display item 2: Redmi Y3 with gradient button
(
    'display', 'product', 21, 2, -- Using Wireless Earbuds (ID=21)
    '/images/ads/display/redmi-y3.png',
    '{"button_type": "gradient"}',
    true
),
-- Display item 3: Philips Ambilight TV with price and discount
(
    'display', 'product', 22, 3, -- Using 4K Smart TV (ID=22)
    '/images/ads/display/ambilighttv.png',
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
    'Headsets', '/images/ads/gaming/headsets.png',
    true
  ),
  -- Gaming item 2: Mouse
  (
    'gaming', 'category', 4, 2,
    'Mouse', '/images/ads/gaming/mouse.png',
    true
  ),
  -- Gaming item 3: Controller
  (
    'gaming', 'category', 4, 3,
    'Controller', '/images/ads/gaming/controller.png',
    true
  ),
  -- Gaming item 4: Chair
  (
    'gaming', 'category', 4, 4,
    'Chair', '/images/ads/gaming/chair.png',
    true
  );

