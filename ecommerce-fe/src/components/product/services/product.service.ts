// src/components/product/services/product.service.ts
import { Product } from "@/types/product.model";

// Generated Products Data
const products: Product[] = [
  // Grocery Products
  {
    type: "normal",
    id: "g1",
    name: "2020 Apple MacBook Pro with Apple M1 Chip (13-inch, 8GB RAM, 256GB SSD Storage) - Space Gray",
    slug: "organic-bananas",
    price: 100.25,
    originalPrice: 125.99,
    discount: 25,
    description: "The MacBook Air 13-inch is a slim and lightweight laptop from Apple with a sharp Retina display and impressive performance powered by the M1 chip. It's known for its portability, sleek design, and reliable performance.",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 5,
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    orders: 100,
    brand: "Apple",
    imgSlider: [
      "/images/products/grocery/02.png",
      "/images/products/grocery/laptop.png",
      "/images/products/grocery/03.png",
      "/images/products/grocery/04.png",
      "/images/products/grocery/05.png",
      "/images/products/grocery/06.png",
      "/images/products/grocery/bananas.png",

    ],
    tags: ["laptop", "apple", "macbook"],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g2",
    name: "Organic Bananas from sustainable farms",
    slug: "organic-bananas",
    price: 2.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Fresh organic bananas from sustainable farms",
    image_url: "/images/products/grocery/laptop.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 19,  
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/laptop.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g3",
    name: "Organic Bananas from sustainable farms",
    slug: "organic-bananas",
    price: 2.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Fresh organic bananas from sustainable farms",
    image_url: "/images/products/grocery/laptop.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 0,  
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g4",
    name: "Organic Bananas from sustainable farms",
    slug: "organic-bananas",
    price: 2.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Fresh organic bananas from sustainable farms",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 150,
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g5",
    name: "Whole Grain Bread",
    slug: "whole-grain-bread",
    price: 3.49,
    originalPrice: 3.99,
    discount: 25,
    description: "Freshly baked whole grain bread",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 75,
    reviews: {
      rating: 4.3,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g6",
    name: "Free Range Eggs (12pk)",
    slug: "free-range-eggs",
    price: 4.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Farm fresh free-range eggs",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 120,
    reviews: {
      rating: 4.7,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g7",
    name: "Organic Milk",
    slug: "organic-milk",
    price: 3.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Fresh organic whole milk",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 95,
    reviews: {
      rating: 4.2,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g8",
    name: "Avocados (3pk)",
    slug: "avocados",
    price: 5.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Ripe and ready to eat avocados",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 65,
    reviews: {
      rating: 4.6,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g9",
    name: "Avocados (3pk)",
    slug: "avocados",
    price: 5.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Ripe and ready to eat avocados",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "1",
    categorySlug: "grocery",
    stockQuantity: 65,
    reviews: {
      rating: 4.6,
      count: 10.878,
    },
    orders: 100,

    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },

  // Home Products
  {
    type: "normal",
    id: "h1",
    name: "Scented Candle Set",
    slug: "scented-candle-set",
    price: 24.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Set of 3 aromatic scented candles",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "2",
    categorySlug: "home",
    stockQuantity: 45,
    orders: 100,
    reviews: {
      rating: 4.4,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "h2",
    name: "Cotton Bed Sheets",
    slug: "cotton-bed-sheets",
    price: 49.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Soft 100% Egyptian cotton bed sheets",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "2",
    categorySlug: "home",
    reviews: {
      rating: 4.7,
      count: 10.878,
    },
    stockQuantity: 30,
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "h3",
    name: "Kitchen Utensil Set",
    slug: "kitchen-utensil-set",
    price: 35.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Complete set of silicone kitchen utensils",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "2",
    categorySlug: "home",
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    stockQuantity: 55,
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "h4",
    name: "Decorative Throw Pillows",
    slug: "decorative-throw-pillows",
    price: 22.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Set of 2 decorative throw pillows for sofa or bed",
    image_url: "/images/products/grocery/bananas.png",
    categoryId: "2",
    categorySlug: "home",
    stockQuantity: 0,
    reviews: {
      rating: 4.3,
      count: 10.878,
    },

    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "h5",
    name: "LED Table Lamp",
    slug: "led-table-lamp",
    price: 39.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Modern LED table lamp with adjustable brightness",
    image_url: "/images/products/home/lamp.jpg",
    categoryId: "2",
    categorySlug: "home",
    stockQuantity: 25,
    reviews: {
      rating: 4.3,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },

  // Fashion Products
  {
    type: "normal",
    id: "f1",
    name: "Women's Denim Jacket",
    slug: "womens-denim-jacket",
    price: 59.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Classic denim jacket for women",
    image_url: "/images/products/fashion/denim-jacket.jpg",
    categoryId: "3",
    reviews: {
      rating: 4.6,
      count: 10.878,
    },
    categorySlug: "fashion",
    stockQuantity: 40,
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "f2",
    name: "Men's Casual T-Shirt",
    slug: "mens-casual-tshirt",
    price: 19.99,
    originalPrice: 3.99,
    discount: 25,
    description: "100% cotton casual t-shirt for men",
    image_url: "/images/products/fashion/tshirt.jpg",
    categoryId: "3",
    categorySlug: "fashion",
    stockQuantity: 120,
    orders: 100,
    reviews: {
      rating: 4.3,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "f3",
    name: "Leather Sneakers",
    slug: "leather-sneakers",
    price: 89.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Comfortable leather sneakers for everyday use",
    image_url: "/images/products/fashion/sneakers.jpg",
    categoryId: "3",
    categorySlug: "fashion",
    stockQuantity: 35,
    orders: 100,
    reviews: {
      rating: 4.3,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "f4",
    name: "Summer Dress",
    slug: "summer-dress",
    price: 45.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Lightweight cotton summer dress",
    image_url: "/images/products/fashion/summer-dress.jpg",
    categoryId: "3",
    categorySlug: "fashion",
    stockQuantity: 0,
    reviews: {
      rating: 4.4,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "f5",
    name: "Unisex Beanie",
    slug: "unisex-beanie",
    price: 14.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Warm knitted beanie for all seasons",
    image_url: "/images/products/fashion/beanie.jpg",
    categoryId: "3",
    categorySlug: "fashion",
    stockQuantity: 60,
    orders: 100,
    reviews: {
      rating: 4.2,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },

  // Electronic Products
  {
    type: "normal",
    id: "e1",
    name: "Wireless Earbuds",
    slug: "wireless-earbuds",
    price: 129.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Bluetooth wireless earbuds with noise cancellation",
    image_url: "/images/products/electronics/earbuds.jpg",
    categoryId: "4",
    categorySlug: "electronic",
    stockQuantity: 45,
    orders: 100,
    reviews: {
      rating: 4.6,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "e2",
    name: "4K Smart TV",
    slug: "4k-smart-tv",
    price: 699.99,
    originalPrice: 3.99,
    discount: 25,
    description: "55-inch 4K Ultra HD Smart LED TV",
    image_url: "/images/products/electronics/smart-tv.jpg",
    categoryId: "4",
    categorySlug: "electronic",
    stockQuantity: 15,
    orders: 100,
    reviews: {
      rating: 4.8,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "e3",
    name: "Portable Bluetooth Speaker",
    slug: "portable-bluetooth-speaker",
    price: 89.99,
    originalPrice: 3.99,
    discount: 25,
    description:
      "Waterproof portable bluetooth speaker with 20-hour battery life",
    image_url: "/images/products/electronics/bluetooth-speaker.jpg",
    categoryId: "4",
    categorySlug: "electronic",
    stockQuantity: 30,
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "e4",
    name: "Digital Camera",
    slug: "digital-camera",
    price: 449.99,
    originalPrice: 3.99,
    discount: 25,
    description: "24MP digital camera with 4K video recording",
    image_url: "/images/products/electronics/digital-camera.jpg",
    categoryId: "4",
    categorySlug: "electronic",
    stockQuantity: 0,
    orders: 100,
    reviews: {
      rating: 4.7,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "e5",
    name: "Smartwatch",
    slug: "smartwatch",
    price: 199.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Fitness tracking smartwatch with heart rate monitor",
    image_url: "/images/products/electronics/smartwatch.jpg",
    categoryId: "4",
    categorySlug: "electronic",
    stockQuantity: 25,
    reviews: {
      rating: 4.4,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },

  // Toys Products
  {
    type: "normal",
    id: "t1",
    name: "Building Blocks Set",
    slug: "building-blocks-set",
    price: 29.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Educational building blocks set for children",
    image_url: "/images/products/toys/building-blocks.jpg",
    categoryId: "5",
    categorySlug: "toys",
    stockQuantity: 50,
    reviews: {
      rating: 4.7,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "t2",
    name: "Remote Control Car",
    slug: "remote-control-car",
    price: 49.99,
    originalPrice: 3.99,
    discount: 25,
    description: "High-speed remote control race car",
    image_url: "/images/products/toys/rc-car.jpg",
    categoryId: "5",
    categorySlug: "toys",
    stockQuantity: 35,
    reviews: {
      rating: 4.5,
      count: 10.878,
    },
    orders: 10,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "t3",
    name: "Stuffed Teddy Bear",
    slug: "stuffed-teddy-bear",
    price: 19.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Soft and huggable teddy bear for kids",
    image_url: "/images/products/toys/teddy-bear.jpg",
    categoryId: "5",
    categorySlug: "toys",
    stockQuantity: 80,
    reviews: {
      rating: 4.8,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "t4",
    name: "Educational Board Game",
    slug: "educational-board-game",
    price: 34.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Family-friendly educational board game",
    image_url: "/images/products/toys/board-game.jpg",
    categoryId: "5",
    categorySlug: "toys",
    stockQuantity: 0,
    orders: 100,
    reviews: {
      rating: 4.6,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "t5",
    name: "Art and Craft Kit",
    slug: "art-and-craft-kit",
    price: 24.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Creative art and craft kit for children",
    image_url: "/images/products/toys/art-kit.jpg",
    categoryId: "5",
    categorySlug: "toys",
    stockQuantity: 40,
    orders: 100,
      reviews: {
        rating: 4.4,
        count: 10.878,
      },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },

  // Additional Grocery Products (for grocery-2)
  {
    type: "normal",
    id: "g6",
    name: "Instant Coffee",
    slug: "instant-coffee",
    price: 7.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Premium instant coffee for quick preparation",
    image_url: "/images/products/grocery/coffee.jpg",
    categoryId: "6",
    categorySlug: "grocery-2",
    stockQuantity: 70,
    orders: 100,
    reviews: {
      rating: 4.4,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g7",
    name: "Chocolate Cookies",
    slug: "chocolate-cookies",
    price: 3.99,
    originalPrice: 3.99,
    discount: 25,
    description: "Crunchy chocolate chip cookies",
    image_url: "/images/products/grocery/cookies.jpg",
    categoryId: "6",
    categorySlug: "grocery-2",
    stockQuantity: 90,
    orders: 100,
    reviews: {
      rating: 4.4,
      count: 10.878,
    },
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
  {
    type: "normal",
    id: "g8",
    name: "Fresh Orange Juice",
    slug: "fresh-orange-juice",
    price: 4.49,
    originalPrice: 3.99,
    discount: 25,
    description: "100% freshly squeezed orange juice",
    image_url: "/images/products/grocery/orange-juice.jpg",
    categoryId: "6",
    categorySlug: "grocery-2",
    stockQuantity: 55,
    reviews: {
      rating: 4.4,
      count: 10.878,
    },
    orders: 100,
    imgSlider: [
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
      "/images/products/grocery/bananas.png",
    ],
    features: [
      {
        id: "f1",
        value: "Free Shipping",
      },
      {
        id: "f2",
        value: "Free 1 Year Warranty",
      },
      {
        id: "f3",
        value: "100% Money-back guarantee",
      },
      {
        id: "f4",
        value: "Secure payment method",
      },
      {
        id: "f5",
        value: "24/7 Customer support",
      },
    ],
    shippingInfo: {
      courier: "2-4 days, free shipping",
      local: "up to one week, $19.00",
      ups: "4-6 days, $29.00",
      global: "3-4 days, $39.00",
    },
  },
];

export const productService = {
  getProducts: () => {
    return Promise.resolve(products);
  },

  getProductsByCategory: (categoryId: string) => {
    return Promise.resolve(
      products.filter((product) => product.categoryId === categoryId)
    );
  },

  getProductsByCategorySlug: (categorySlug: string) => {
    return Promise.resolve(
      products.filter((product) => product.categorySlug === categorySlug)
    );
  },

  getProductById: (productId: string) => {
    return Promise.resolve(
      products.find((product) => product.id === productId)
    );
  },

  getProductBySlug: (productSlug: string) => {
    return Promise.resolve(
      products.find((product) => product.slug === productSlug)
    );
  },

  getProductByCategoryAndSlug: (categorySlug: string, productSlug: string) => {
    return Promise.resolve(
      products.find(
        (product) =>
          product.categorySlug === categorySlug && product.slug === productSlug
      )
    );
  },
};
