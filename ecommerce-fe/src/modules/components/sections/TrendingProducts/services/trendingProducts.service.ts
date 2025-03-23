import { TrendingProductItem } from "../models/trendingProducts.model";

const urlImgTrending = "images/trending/";
// Mock data - trong thực tế sẽ fetch từ API
const trendingProductData: TrendingProductItem[] = [
  {
    id: 1,
    title: "Radiant Glow Hydrating Serum",
    description:
      "Gentle yet effective, our Radiant Boosting Foaming our Radiant Boosting Foaming our Radiant Boosting Foaming",
    price: 29.99,
    originalPrice: 39.99,
    discount: 20,
    review: {
      rating: 4.8,
      count: 345,
    },
    imageUrl: `${urlImgTrending}serum.png`,
  },
  {
    id: 2,
    title: "Modern Minimalist Vase",
    description:
      "Track your workouts, heart rate, sleep quality and receive notifications. Water resistant up to 50m with 7-day battery life.",
    price: 40.99,
    originalPrice: "",
    discount: "",
    review: {
      rating: 4.6,
      count: 842,
    },
    imageUrl: `${urlImgTrending}vase.png`,
  },
  {
    id: 3,
    title: "FitPro 3000 Smart Watch",
    description:
      "Fast-charging power bank with dual USB ports and USB-C compatibility. Charge multiple devices simultaneously on the go.",
    price: 119.99,
    originalPrice: "",
    discount: "",
    review: {
      rating: 4.0,
      count: 2105,
    },
    imageUrl: `${urlImgTrending}smartwatch.png`,
  },
];



// resolve data trendingProduct
export const TrendingPorductService = {
  getDataTredings: (): Promise<TrendingProductItem[]> => {
    return Promise.resolve(trendingProductData);
  }
};