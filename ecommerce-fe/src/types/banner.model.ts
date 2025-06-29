// Global banner model definitions

export interface Banner {
  id: number;
  title: string;
  subtitle?: string;
  description?: string;
  imageUrl: string;
  linkUrl?: string;
  actionText?: string;
  isActive: boolean;
  priority?: number;
  type: BannerType;
  productId?: number;
  categoryId?: number;
  discount?: string;
  highlightText?: string;
  backgroundColor?: string;
  textColor?: string;
  animationType?: 'fade' | 'slide' | 'zoom';
  createdAt?: string;
  updatedAt?: string;
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
  imageUrl: string;
  linkUrl?: string;
  actionText?: string;
  isActive: boolean;
  priority?: number;
  discount?: string;
  highlightText?: string;
  backgroundColor?: string;
  textColor?: string;
  animationType?: 'fade' | 'slide' | 'zoom';
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
  isActive?: boolean;
  productId?: number;
  categoryId?: number;
  page?: number;
  pageSize?: number;
  orderBy?: string;
  orderDir?: 'ASC' | 'DESC';
} 