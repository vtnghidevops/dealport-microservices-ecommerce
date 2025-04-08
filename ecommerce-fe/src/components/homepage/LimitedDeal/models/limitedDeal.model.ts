

export interface LimitedDealItem {
  id: string;
  title: string;
  description: string;
  price?: number | string;
  originalPrice?: number | string;
  discount?: number | string;
  review: {
    rating: number;
    count: number;
  };
  image_url: string;
}