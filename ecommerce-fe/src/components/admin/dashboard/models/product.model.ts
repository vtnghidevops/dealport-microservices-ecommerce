export interface Product {
  id: number;
  name: string;
  category: string;
  price: number;
  itemCode?: string;
  image_url?: string;
}

export interface BestSellingProduct {
  product: Product;
  totalOrder: number;
  status: "Stock" | "Stock out";
  price: number;
}

export interface ProductCategory {
  id: number;
  name: string;
  image_url: string;
}

export interface NewProduct {
  id: number;
  name: string;
  price: number;
  image_url?: string;
  itemCode?: string;
}
