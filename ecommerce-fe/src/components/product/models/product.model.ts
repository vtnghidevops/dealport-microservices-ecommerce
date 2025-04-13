export interface Product {
  id: string;
  name: string;
  slug: string;
  price: number;
  originalPrice: number;
  discount: number;
  description: string;
  image_url: string;
  categoryId: string;
  categorySlug: string; // add categorySlug link to category page
  stock: number;
  rating?: number;
  brand?: string;
  tags?: string[];
  reviews: number;
  orders: number;
  imgSlider: string[];
  features: {
    id: string;
    value: string;
  }[];
  shippingInfo: {
    courier: string;
    local: string;
    ups: string;
    global: string;
  };
}