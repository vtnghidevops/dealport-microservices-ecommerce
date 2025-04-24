export interface Product {
  type: string; // normal, trending, top-sale, new, limited
  id: number;
  name: string;
  description: string;
  slug: string;
  price: number;
  image_url: string;
  categoryId: string;
  categorySlug: string;
  stockQuantity: number;
  originalPrice: number;
  discount: number;
  brand: string;
  tags: string[];
  reviewsAvg: ProductRating;
  orders: number;
  imgSlider: string[];
  features: string[];
  shippingInfo: ShippingInfo;
  uiMetadata?: any;
  images?: ProductImage[];
}

export interface ProductImage {
  id?: number;
  product_id?: number;
  url: string;
  is_primary: boolean;
  display_order?: number;
}

export interface ProductReview {
  id?: number;
  product_id: number;
  user_id: number;
  user_name?: string;
  rating: number;
  comment: string;
  created_at?: string;
}

export interface ShippingInfo {
  courier: string;
  local: string;
  ups: string;
  global: string;
}

export interface ProductRating {
  rating: number;
  count: number;
}

