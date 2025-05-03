import React, { useState, useEffect } from "react";
import CategoryFilter from "./filter/CategoryFilter";
import CategoryTable from "./tables/CategoryTable";
import CategoryCard from "./discover/CategoryCard";
import Pagination from "../../common/Pagination";
import {
  Category,
  CategoryFilter as FilterType,
  CategoryFilterCounts,
} from "./models/category.model";
import {
  CategoryService,
} from "./services/category.service";
import AdminHeader from "../layout/AdminHeader";
import { LuCirclePlus } from "react-icons/lu";
import { BsThreeDotsVertical } from "react-icons/bs";
import { IoIosArrowForward } from "react-icons/io";
import { useToast } from "@/hooks/use-toast";
/**
 * Main component for category management
 * Contains all functionality for managing product categories
 */
export const CategoryManagement: React.FC = () => {
  // State for categories and filter
  const [categories, setCategories] = useState<Category[]>([]);
  const [discoverCategories, setDiscoverCategories] = useState<Array<{ name: string, imageUrl: string }>>([]);
  const [loading, setLoading] = useState(true);
  const [loadingDiscover, setLoadingDiscover] = useState(true);
  const [filter, setFilter] = useState<FilterType>({
    search: "",
    status: "all",
    page: 1,
    limit: 10,
  });
  const [totalItems, setTotalItems] = useState(0);
  const [activeFilterTab, setActiveFilterTab] = useState<
    "all" | "featured" | "onSale" | "outOfStock"
  >("all");
  const [filterCounts, setFilterCounts] = useState<CategoryFilterCounts>({
    all: 0,
    featured: 0,
    onSale: 0,
    outOfStock: 0,
  });
  const { toast } = useToast();

  // Fetch filter counts on component mount
  useEffect(() => {
    const fetchFilterCounts = async () => {
      try {
        const counts = await CategoryService.getCategoryFilterCounts();
        setFilterCounts(counts);
      } catch (error) {
        console.error("Error fetching filter counts:", error);
      }
    };

    fetchFilterCounts();
  }, []);

  // Load discover categories from real API data
  useEffect(() => {
    const fetchDiscoverCategories = async () => {
      setLoadingDiscover(true);
      try {
        // Get all categories for discover section (no pagination)
        const response = await CategoryService.getCategories({
          search: "",
          status: "active",
          limit: 8 // Limit to 8 categories for the discover section
        });

        // Map to the format needed for CategoryCard component
        const mappedCategories = response.categories.map(cat => ({
          name: cat.name,
          imageUrl: cat.imageUrl && cat.imageUrl.startsWith('http')
            ? cat.imageUrl
            : getDefaultCategoryImage(cat.name)
        }));

        setDiscoverCategories(mappedCategories);
      } catch (error) {
        console.error("Error fetching discover categories:", error);
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to load discover categories."
        });

        // Fallback to default categories if API fails
        setDiscoverCategories([
          { name: 'Electronics', imageUrl: '/images/exploring/electronic.png' },
          { name: 'Fashion', imageUrl: '/images/exploring/fashion.png' },
          { name: 'Home & Kitchen', imageUrl: '/images/exploring/home.png' },
          { name: 'Sports', imageUrl: '/images/exploring/grocery.png' },
        ]);
      } finally {
        setLoadingDiscover(false);
      }
    };

    fetchDiscoverCategories();
  }, [toast]);

  // Helper function to get default image based on category name
  const getDefaultCategoryImage = (name: string): string => {
    const nameLower = name.toLowerCase();
    if (nameLower.includes('electronic')) return '/images/exploring/electronic.png';
    if (nameLower.includes('fashion') || nameLower.includes('cloth')) return '/images/exploring/fashion.png';
    if (nameLower.includes('home') || nameLower.includes('kitchen')) return '/images/exploring/home.png';
    if (nameLower.includes('sport')) return '/images/exploring/grocery.png';
    if (nameLower.includes('toy') || nameLower.includes('game')) return '/images/exploring/toys.png';
    if (nameLower.includes('book')) return '/images/exploring/books.png';
    return '/images/exploring/category.png'; // Default image
  };

  // Fetch categories on component mount and when filter changes
  useEffect(() => {
    const fetchCategories = async () => {
      setLoading(true);
      try {
        const response = await CategoryService.getCategories({
          ...filter,
          productFilter:
            activeFilterTab !== "all" ? activeFilterTab : undefined,
        });
        setCategories(response.categories);
        setTotalItems(response.total);
      } catch (error) {
        console.error("Error fetching categories:", error);
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to load categories. Please try again."
        });
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, [filter, activeFilterTab, toast]);

  // Handle editing a category
  const handleEditCategory = (category: Category) => {
    console.log("Edit category:", category);
    // Implement edit functionality here
  };

  // Handle deleting a category
  const handleDeleteCategory = async (id: string) => {
    try {
      setLoading(true);
      const success = await CategoryService.deleteCategory(id);

      if (success) {
        // Remove the deleted category from the state
        setCategories(categories.filter(category => category.id !== id));

        toast({
          variant: "success",
          title: "Success",
          description: "Category deleted successfully"
        });
      } else {
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to delete category"
        });
      }
    } catch (error) {
      console.error("Error deleting category:", error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to delete category. Please try again."
      });
    } finally {
      setLoading(false);
    }
  };

  // Handle adding a product
  const handleAddProduct = () => {
    console.log("Add product clicked");
    // Implement add product functionality
  };

  // Handle page change
  const handlePageChange = (page: number) => {
    const totalPages = Math.ceil(totalItems / filter.limit!);
    const validPage = Math.max(1, Math.min(page, totalPages));

    setFilter({
      ...filter,
      page: validPage,
    });
  };
  // Handle tab filter change
  const handleTabFilterChange = (
    tab: "all" | "featured" | "onSale" | "outOfStock"
  ) => {
    setActiveFilterTab(tab);
    setFilter({
      ...filter,
      page: 1,
    });
  };

  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Categories" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            {/* Discover Section */}
            <div className="mb-8">
              <div className="flex justify-between items-center mb-[25px]">
                <h2 className="text-[22px] font-bold text-cyprus">Discover</h2>
                <div className="flex gap-[12px] h-[48px]">
                  <button
                    onClick={handleAddProduct}
                    className="gap-4 flex justify-center items-center w-[142px] py-6 px-12  bg-ocean-green text-white rounded-lg hover:bg-green-600 transition-colors"
                  >
                    <LuCirclePlus className=""></LuCirclePlus>
                    <span className="text-[15px] font-bold">Add Product</span>
                  </button>
                  <button className="gap-4 flex justify-center items-center w-[142px] py-6 px-12  bg-white border text-cyprus rounded-lg hover:bg-gray-50 transition-colors">
                    <span className="text-[15px] font-bold">More Action</span>
                    <BsThreeDotsVertical className=""></BsThreeDotsVertical>
                  </button>
                </div>
              </div>

              {/* Category Cards Grid */}
              <div className="flex flex-wrap gap-16 mb-8 relative">
                {loadingDiscover ? (
                  // Loading skeleton for discover categories
                  <>
                    {[1, 2, 3, 4].map((index) => (
                      <div key={index} className="animate-pulse">
                        <div className="bg-gray-200 rounded-lg w-[120px] h-[120px] mb-2"></div>
                        <div className="bg-gray-200 h-5 w-24 rounded"></div>
                      </div>
                    ))}
                  </>
                ) : (
                  // Actual discover categories from API
                  discoverCategories.map((category, index) => (
                    <CategoryCard
                      key={index}
                      name={category.name}
                      imageUrl={category.imageUrl}
                    />
                  ))
                )}

                <div className="absolute right-[2%] top-1/2 -translate-y-1/2 ">
                  <button className="hover:bg-neutral-50 w-[48px] h-[48px] rounded-full bg-white flex justify-center items-center drop-shadow-md">
                    <IoIosArrowForward />
                  </button>
                </div>
              </div>
            </div>

            {/* Category Filter and Table */}
            <div className="w-[1116px] h-[979px] bg-white rounded-lg shadow p-[1rem] mb-6 mt-[1rem] drop-shadow filter">
              {/* Product Filter Tabs */}
              <CategoryFilter
                counts={filterCounts}
                onFilterChange={(filter) => handleTabFilterChange(filter)}
                activeFilter={activeFilterTab}
                onSearch={(term) => {
                  setFilter({
                    ...filter,
                    search: term,
                    page: 1,
                  });
                }}
                loading={loading}
              />

              {/* Categories Table */}
              {loading ? (
                <div className="animate-pulse">
                  <div className="h-8 bg-gray-200 rounded mb-4"></div>
                  <div className="h-40 bg-gray-200 rounded"></div>
                </div>
              ) : (
                <CategoryTable
                  categories={categories}
                  onEdit={handleEditCategory}
                  onDelete={handleDeleteCategory}
                />
              )}
              <div className="mt-[3rem]">
                <Pagination
                  currentPage={filter.page || 1}
                  totalItems={totalItems}
                  pageSize={filter.limit || 10}
                  onPageChange={handlePageChange}
                />
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};
