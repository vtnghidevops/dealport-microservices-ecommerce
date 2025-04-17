import { Product } from "@/types/product.model";

export interface TopProductItem extends Product {
  setUpDesign?: string;   // set up container with row or col or double - row, double, col
  actionLabel?: string;     // "Buy Now", "Visit Store", etc.
  actionLink?: string;
  isCommingSoon?: boolean;
  gridSpan: {
    col: number;  
    row: number;  
  };
  isImageDouble?: boolean
  image_double_url?: string
}

export interface ProductCard {
  product: TopProductItem;
  onClick?: () => void;
}

export interface TopProductsProps {
  title?: string;
  products: TopProductItem[];
  viewAllLabel?: string;
  onViewAllClick?: () => void;
  showNavigationArrow?: boolean;
  onItemClick?: (product: TopProductItem, index: number) => void;
}