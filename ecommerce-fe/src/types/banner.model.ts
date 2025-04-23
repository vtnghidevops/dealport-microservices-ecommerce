// Global banner model definitions

export interface Banner {
  id: number;
  title: string;
  subtitle?: string;
  description?: string;
  image_url: string;
  link_url?: string;
  action_text?: string;
  is_active: boolean;
  priority?: number;
  type: BannerType;
  product_id?: number;
  category_id?: number;
  discount?: string;
  highlight_text?: string;
  background_color?: string;
  text_color?: string;
  animation_type?: 'fade' | 'slide' | 'zoom';
  created_at?: string;
  updated_at?: string;
}

export enum BannerType {
  HERO = 'hero',
  PROMOTION = 'promotional',
  CATEGORY = 'category',
  PRODUCT = 'product',
  SEASONAL = 'seasonal',
  FEATURED = 'featured'
}

// Frontend-specific model for slider banners
export interface SliderBannerItem {
  id: number;
  title: string;
  subtitle?: string;
  description?: string;
  image_url: string;
  link_url?: string;
  action_text?: string;
  is_active: boolean;
  priority?: number;
  discount?: string;
  highlight_text?: string;
  background_color?: string;
  text_color?: string;
  animation_type?: 'fade' | 'slide' | 'zoom';
}

export interface BannerResponse {
  status: number;
  data: Banner[];
  meta?: {
    current_page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
  };
}

export interface BannerFilter {
  type?: BannerType;
  is_active?: boolean;
  product_id?: number;
  category_id?: number;
  page?: number;
  page_size?: number;
  order_by?: string;
  order_dir?: 'ASC' | 'DESC';
} 