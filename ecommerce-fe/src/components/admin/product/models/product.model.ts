// components/admin/product/models/product.model.ts
export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  discountedPrice: number;
  saleAmount: number;
  taxIncluded: boolean;
  expirationStart: string;
  expirationEnd: string;
  stockQuantity: string;
  stockStatus: string;
  highlighted: boolean;
  images: string[];
  categories: string[];
  tags: string[];
  color: string;
}