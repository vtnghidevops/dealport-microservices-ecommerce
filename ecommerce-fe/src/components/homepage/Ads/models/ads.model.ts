export interface BannerShowCaseItem {
  id: number,
  imageUrl: string,
  buttonType?: string,
  buttonText?: string,
  hasMore: boolean,
  type: string,     // Type of item (product or category)
  categorySlug: string,      // Slug for navigation
};

export interface DisplayItem {
  id: number,
  imageUrl: string,
  buttonType?: string,
  price?: string | number,
  title?: string,
  discountImg?: string,
  type: string,     // Type of item (product or category)
  productSlug: string,      // Slug for navigation
  categorySlug: string // product type
}


export interface GamingItem {
  id: number,
  subtitle: string,
  imageUrl: string,
  type: string,     // Type of item (product or category)
  categorySlug: string // product type
}

export interface NewFashionItem {
  id: number,
  title: string,
  imageUrl: string,
  buttonType: string,
  type: string,     // Type of item (product or category)
  productSlug: string,      // Slug for 
  categorySlug: string // product type
}