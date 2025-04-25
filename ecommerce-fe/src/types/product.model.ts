export interface Product {
  type: string; // normal, trending, top-sale, new, limited
  id: number;
  name: string;
  description: string;
  slug: string;
  price: number;
  imageUrl: string;
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
  productId?: number;
  url: string;
  isPrimary: boolean;
  displayOrder?: number;
}

export interface ProductReview {
  id?: number;
  productId: number;
  userId: number;
  userName?: string;
  rating: number;
  comment: string;
  createdAt?: string;
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

