
export interface CategoryItem {
  id: string;
  name: string;
  image_url: string;
  link?: string;
  slug: string;
}

export interface CategoryExplorerProps {
  title?: string;
  categories: CategoryItem[];
  viewAllLabel?: string;
  onViewAllClick?: () => void;
  itemWidth?: string;
  itemHeight?: string;
  className?: string;
  showNavigationArrow?: boolean;
  onItemClick?: (category: CategoryItem, index: number) => void;
}

export interface CategoryExplorerCard {
  title?: string;
  category: CategoryItem;
  viewAllLabel?: string;
  itemWidth?: string;
  itemHeight?: string;
  className?: string;
  showNavigationArrow?: boolean;
  onViewAllClick?: () => void;
  onItemClick?: (category: CategoryItem, index: number) => void;
}