export interface Product {
  id: string;
  name: string;
  slug: string;
  price: number;
  description: string;
  imageUrl: string;
  categoryId: string;
  categorySlug: string; // add categorySlug link to category page
  stock: number;
  rating?: number;
}