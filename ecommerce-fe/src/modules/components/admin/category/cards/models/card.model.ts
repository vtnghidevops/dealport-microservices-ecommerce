// Models for the category summary cards
export interface CategorySummaryData {
  totalCategories: number;
  activeCategories: number;
  featuredCategories: number;
  popularCategory: {
    name: string;
    productCount: number;
  };
}