

export interface LimitedDealItem {
  id: number;
  title: string;
  description: string;
  price: number | string;
  originalPrice: number | string;
  discount: number | string;
  review: {
    rating: number;
    count: number;
  };
  imageUrl: string;
}