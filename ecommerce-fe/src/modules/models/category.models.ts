export interface Category {
  id?: string | number;
  name: string;
  image: string;
  link?: string;
}

export interface CategoryExplorerProps {
  title?: string;
  categories: Category[];
  viewAllLabel?: string;
  onViewAllClick?: () => void;
  itemWidth?: string;
  itemHeight?: string;
  className?: string;
  showNavigationArrow?: boolean;
  onItemClick?: (category: Category, index: number) => void;
}