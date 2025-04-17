// import { TrendingProductItem } from "../models/trendingProducts.model";
import { Product } from "@/types/product.model";
const urlImgTrending = "/images/trending/";
// Mock data - trong thực tế sẽ fetch từ API
const trendingProductData: Product[] = [
  {
    id: 't1',
    name: "Radiant Glow Hydrating Serum",
    description:
      "Gentle yet effective, our Radiant Boosting Foaming our Radiant Boosting Foaming our Radiant Boosting Foaming",
    price: 29.99,
    originalPrice: 39.99,
    discount: 20,
    type: "trending",
    reviews: {
      rating: 4.8,
      count: 345,
    },
    image_url: `${urlImgTrending}serum.png`,
    categoryId: "1",
    categorySlug: "grocery",
    slug: "radiant-glow-hydrating-serum",
    stockQuantity: 5,
  },
  {
    type: "trending",
    id: 't2',
    name: "Modern Minimalist Vase",
    description:
      "Track your workouts, heart rate, sleep quality and receive notifications. Water resistant up to 50m with 7-day battery life.",
    price: 40.99,
    originalPrice: 49.99,
    discount: 18,
    reviews: {
      rating: 4.6,
      count: 842,
    },
    image_url: `${urlImgTrending}vase.png`,
    categoryId: "1",
    categorySlug: "grocery",
    slug: "modern-minimalist-vase",
    stockQuantity: 5,
  },
  {
    type: "trending",
    id: 't3',
    name: "FitPro 3000 Smart Watch",
    description:
      "Fast-charging power bank with dual USB ports and USB-C compatibility. Charge multiple devices simultaneously on the go.",
    price: 119.99,
    originalPrice: 0,
    discount: 0,
    reviews: {
      rating: 4.0,
      count: 2105,
    },
    image_url: `${urlImgTrending}smartwatch.png`,
    categoryId: "1",
    categorySlug: "grocery",
    slug: "fitpro-3000-smart-watch",
    stockQuantity: 5,
  },
];



// resolve data trendingProduct
export const TrendingPorductService = {
  getDataTredings: (): Promise<Product[]> => {
    return Promise.resolve(trendingProductData);
  }
};