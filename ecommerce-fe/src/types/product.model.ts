export interface Product {
  type: string; // normal, trending, top-sale, new, limited
  id: string;
  name: string;
  description: string;
  slug: string;
  price: number;
  image_url: string;
  categoryId: string;
  categorySlug: string;
  stockQuantity: number;


  originalPrice?: number;
  discount?: number;
  images?: string[];
  categories?: string[];

  brand?: string;
  tags?: string[];
  reviews?: {
    rating?: number;
    count?: number;
  };
  orders?: number;
  imgSlider?: string[];

  features?: {
    id: string;
    value: string;
  }[];
  
  shippingInfo?: {
    courier?: string;
    local?: string;
    ups?: string;
    global?: string;
  };

  // // Các trường từ dashboard product
  // category?: string;
  // itemCode?: string;

  // // Các trường từ admin product
  // saleAmount?: number;
  // taxIncluded?: boolean;
  // expirationStart?: string;
  // expirationEnd?: string;
  // highlighted?: boolean;
  // color?: string;
}


