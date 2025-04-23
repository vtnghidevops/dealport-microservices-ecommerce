export interface BannerShowCaseItem {
  id: number,
  image_url: string,
  buttonType?: string,
  buttonText?: string,
  hasMore: boolean,
  type: string,     // Type of item (product or category)
  categorySlug: string,      // Slug for navigation
};

export interface DisplayItem {
  id: number,
  image_url: string,
  buttonType?: string,
  price?: string | number,
  title?: string,
  discount_img?: string,
  type: string,     // Type of item (product or category)
  productSlug: string,      // Slug for navigation
  categorySlug: string // product type
}


export interface GamingItem {
  id: number,
  subtitle: string,
  image_url: string,
  type: string,     // Type of item (product or category)
  categorySlug: string // product type
}

export interface NewFashionItem {
  id: number,
  title: string,
  image_url: string,
  buttonType: string,
  type: string,     // Type of item (product or category)
  productSlug: string,      // Slug for 
  categorySlug: string // product type
}