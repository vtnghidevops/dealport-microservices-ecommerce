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
  discoverCategories,
} from "./services/category.service";
import AdminHeader from "../layout/AdminHeader";
import { LuCirclePlus } from "react-icons/lu";
import { BsThreeDotsVertical } from "react-icons/bs";
import { IoIosArrowForward } from "react-icons/io";
/**
 * Main component for category management
 * Contains all functionality for managing product categories
 */
export const CategoryManagement: React.FC = () => {
  // State for categories and filter
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);
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
      } finally {
        setLoading(false);
      }
    };

    fetchCategories();
  }, [filter, activeFilterTab]);

  // Handle editing a category
  const handleEditCategory = (category: Category) => {
    console.log("Edit category:", category);
    // Implement edit functionality here
  };

  // Handle deleting a category
  const handleDeleteCategory = (id: number) => {
    console.log("Delete category with ID:", id);
    // Implement delete functionality here
  };

  // Handle adding a product
  const handleAddProduct = () => {
    console.log("Add product clicked");
    // Implement add product functionality
  };

  // Handle filter changes
  // const handleFilterChange = (newFilter: FilterType) => {
  //   setFilter(newFilter);
  // };

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
                {discoverCategories.map((category, index) => (
                  <CategoryCard
                    key={index}
                    name={category.name}
                    imageUrl={category.imageUrl}
                  />
                ))}
                <div className="absolute right-[2%] top-1/2 -translate-y-1/2 ">
                  <button className="hover:bg-neutral-50 w-[48px] h-[48px] rounded-full bg-white flex justify-center items-center drop-shadow-md">
                    <IoIosArrowForward />
                  </button>
                </div>
              </div>
            </div>

            {/* Category Filter and Table */}
            {/* <div className="w-[1116px] h-[979px] bg-white rounded-lg shadow p-[1rem] mb-6 mt-[1rem] drop-shadow filter">
                  <OrderFilter
                    onSearch={handleSearch}
                    onFilterChange={handleFilterChange}
                    counts={filterCounts}
                    loading={loading}
                  /> */}
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
                // <CategoryTable
                //   orders={orders}
                //   onStatusChange={handleStatusChange}
                //  onViewDetails={handleViewDetails}
                //                   />
              )}
              <div className="mt-[3rem]">
                {/* {renderPagination()} */}
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
