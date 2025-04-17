// Global banner model definitions

export interface Banner {
  id: string;
  title: string;
  subtitle?: string;
  description?: string;
  image_url: string;
  link_url?: string;
  action_text?: string;
  position?: string;
  isActive: boolean;
  startDate?: string;
  endDate?: string;
  priority?: number;
  type: BannerType;
  createdAt?: string;
  updatedAt?: string;
}

export enum BannerType {
  HERO = 'hero',
  PROMOTION = 'promotion',
  CATEGORY = 'category',
  COLLECTION = 'collection',
  SEASONAL = 'seasonal',
  FEATURED = 'featured'
}

export interface SliderBannerItem extends Omit<Banner, 'type'> {
  discount?: number | string;
  highlight_text?: string;
  background_color?: string;
  text_color?: string;
  animation_type?: 'fade' | 'slide' | 'zoom';
}

export interface BannerResponse {
  banners: Banner[];
  total: number;
  page?: number;
  limit?: number;
}

export interface BannerFilter {
  type?: BannerType;
  isActive?: boolean;
  search?: string;
  startDate?: string;
  endDate?: string;
  sortBy?: 'createdAt' | 'priority';
  sortOrder?: 'asc' | 'desc';
  page?: number;
  limit?: number;
} 