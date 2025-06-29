import { Product } from "@/types/product.model";

export interface TopProductItem extends Product {
  uiMetadata: {
    setUpDesign: string;
    isCommingSoon: boolean;
  }
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