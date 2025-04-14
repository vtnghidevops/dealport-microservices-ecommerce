import { LimitedDealItem } from "../models/limitedDeal.model";

// Mock data - trong thực tế sẽ fetch từ API
//   id: number;
//   title: string;
//   description: string;
//   price: number | string;
//   originalPrice: number | string;
//   discount: number | string;
//   review: {
//     rating: number;
//     count: number;
//   };
//   imageUrl: string;
const urlLimited = "images/limited/" ;
const limitedDealData: LimitedDealItem[] = [
  {
    id: 'l1',
    title: "Samsung Galaxy S24",
    description: "Gentle yet effective, our hydrating serum infuses skin with essential moisture while brightening and plumping for a radiant complexion.",
    price: 29.99,
    originalPrice: 39.99,
    discount: 25,
    review: {
      rating: 4.8,
      count: 345,
    },
    image_url: `${urlLimited}samsung_s24.png`,
  },
  {
    id: 'l2',
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
  },
  {
    id: 'l3',
    title: "Winter fashion jacket",
    description: "Delivers extreme hydration and helps strengthen the skin barrier. Perfect for both daytime wear and overnight rejuvenation. [[2]]",
    price: 39.99,
    originalPrice: 49.99,
    discount: 20,
    review: {
      rating: 4.9,
      count: 412,
    },
    image_url: `${urlLimited}winter_jacket.png`,
  },
  {
    id: 'l4',
    title: "New Balance 574 Senekers",
    description: "Illuminating serum infused with glow-boosting intelligent botanicals that leave the skin visibly radiant while reducing redness. [[4]]",
    price: 44.99,
    originalPrice: 59.99,
    discount: 25,
    review: {
      rating: 4.6,
      count: 278,
    },
    image_url: `${urlLimited}senekers.png`,
  },
  {
    id: 'l5',
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
  },
  {
    id: 'l6',
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
  },
  {
    id: 'l7',
    title: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    review: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
  }
]



// resolve data trendingProduct
export const LimitedDealService = {
  getLitmitedData: (): Promise<LimitedDealItem[]> => {
    return Promise.resolve(limitedDealData);
  }
};