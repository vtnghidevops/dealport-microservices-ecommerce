import { Product } from "@/types/product.model";

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
const urlLimited = "/images/limited/" ;
const limitedDealData: Product[] = [
  {
    id: 'l1',
    name: "Samsung Galaxy S24",
    description: "Gentle yet effective, our hydrating serum infuses skin with essential moisture while brightening and plumping for a radiant complexion.",
    price: 29.99,
    originalPrice: 39.99,
    discount: 25,
    reviews: {
      rating: 4.8,
      count: 345,
    },
    image_url: `${urlLimited}samsung_s24.png`,
    type: "limited-deal",
    slug: "samsung-galaxy-s24",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  },
  {
    id: 'l2',
    name: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    reviews: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
    type: "limited-deal",
    slug: "ui-tws-7002-earbud",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  },
  {
    id: 'l3',
    name: "Winter fashion jacket",
    description: "Delivers extreme hydration and helps strengthen the skin barrier. Perfect for both daytime wear and overnight rejuvenation. [[2]]",
    price: 39.99,
    originalPrice: 49.99,
    discount: 20,
    reviews: {
      rating: 4.9,
      count: 412,
    },
    image_url: `${urlLimited}winter_jacket.png`,
    type: "limited-deal",
    slug: "winter-fashion-jacket",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  },
  {
    id: 'l4',
    name: "New Balance 574 Senekers",
    description: "Illuminating serum infused with glow-boosting intelligent botanicals that leave the skin visibly radiant while reducing redness. [[4]]",
    price: 44.99,
    originalPrice: 59.99,
    discount: 25,
    reviews: {
      rating: 4.6,
      count: 278,
    },
    image_url: `${urlLimited}senekers.png`,
    type: "limited-deal",
    slug: "new-balance-574-senekers",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  },
  {
    id: 'l5',
    name: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    reviews: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
    type: "limited-deal",
    slug: "ui-tws-7002-earbud",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  },
  {
    id: 'l6',
    name: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    reviews: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
    type: "limited-deal",
    slug: "ui-tws-7002-earbud",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  },
  {
    id: 'l7',
    name: "Ui TWS 7002 Earbud",
    description: "Exceptional hydration with antioxidants, hyaluronic acid and collagen for visibly smooth and glowy skin. Fragrance-free for all skin types. [[1]] [[7]]",
    price: 34.99,
    originalPrice: 42.99,
    discount: 18,
    reviews: {
      rating: 4.7,
      count: 287,
    },
    image_url: `${urlLimited}earbud.png`,
    type: "limited-deal",
    slug: "ui-tws-7002-earbud",
    categoryId: "1",
    categorySlug: "category",
    stockQuantity: 0,
  }
]



// resolve data trendingProduct
export const LimitedDealService = {
  getLitmitedData: (): Promise<Product[]> => {
    return Promise.resolve(limitedDealData);
  }
};