export interface FeaturedProductProps {
  id?: string | number;
  name: string;
  price?: number | string;
  image: string;
  setUpDesign?: string;   // set up container with row or col or double - row, double, col
  actionLabel?: string;     // "Buy Now", "Visit Store", etc.
  actionLink?: string;
  isCommingSoon?: boolean;
  gridSpan?: {
    col?: number;  // Số cột mà sản phẩm chiếm
    row?: number;  // Số hàng mà sản phẩm chiếm
  };
  isImageDouble?: boolean
  image_double?: string
}

export interface ProductCardProps {
  product: FeaturedProductProps;
  onClick?: () => void;
}

export interface TopProductsProps {
  title?: string;
  products: FeaturedProductProps[];
  viewAllLabel?: string;
  onViewAllClick?: () => void;
  showNavigationArrow?: boolean;
  onItemClick?: (product: FeaturedProductProps, index: number) => void;
}