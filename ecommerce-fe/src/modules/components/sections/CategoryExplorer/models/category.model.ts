
export interface CategoryItem {
  id?: string | number;
  name: string;
  image: string;
  link?: string;
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
  onViewAllClick?: () => void;
  itemWidth?: string;
  itemHeight?: string;
  className?: string;
  showNavigationArrow?: boolean;
  onItemClick?: (category: CategoryItem, index: number) => void;
}