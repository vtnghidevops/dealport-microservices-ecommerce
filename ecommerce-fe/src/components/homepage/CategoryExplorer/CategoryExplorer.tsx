import React, { useRef, useState, useEffect } from "react";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { CategoryExplorerProps } from "./models/category.model";
import CategoryCard from "./CategoryCard";
import { CategoryItem } from "./models/category.model";
import { CategoryExplorerService } from "./services/categoryExplorer.service";
import { handleViewAll } from "../../../utils/helpers";
import { handleCategoryClick } from "../../../utils/helpers";

const CategoryExplorer: React.FC = () => {
  // fetch data
  const [categories, setCategories] = useState<CategoryItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  useEffect(() => {
    const fetchCategoryExplorer = async () => {
      try {
        const data = await CategoryExplorerService.getCategoryData();
        setCategories(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching testimonials:", error);
        setLoading(false);
      }
    };

    fetchCategoryExplorer();
  }, []);

  var categoryExplore: CategoryExplorerProps = {
    title: "Start exploring now",
    categories: categories,
    viewAllLabel: "View All",
    onViewAllClick: handleViewAll,
    itemWidth: "w-[180px]",
    itemHeight: "h-[220px]",
    showNavigationArrow: true,
    onItemClick: handleCategoryClick,
  };

  // scroll event
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const [showLeftArrow, setShowLeftArrow] = useState(false);
  const [showRightArrow, setShowRightArrow] = useState(true);

  // Check if arrows should be shown
  const checkForArrows = () => {
    const container = scrollContainerRef.current;
    if (!container) return;

    // Show left arrow if scrolled to the right
    setShowLeftArrow(container.scrollLeft > 20);

    // Show right arrow if there's more content to scroll
    const hasMoreToScroll =
      container.scrollWidth > container.clientWidth + container.scrollLeft + 20;
    setShowRightArrow(hasMoreToScroll);
  };

  // Initial check and add scroll listener
  useEffect(() => {
    checkForArrows();
    const container = scrollContainerRef.current;

    if (container) {
      container.addEventListener("scroll", checkForArrows);
      // Check on window resize as well
      window.addEventListener("resize", checkForArrows);
    }

    return () => {
      if (container) {
        container.removeEventListener("scroll", checkForArrows);
      }
      window.removeEventListener("resize", checkForArrows);
    };
  }, [categories]);

  const scrollRight = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: 300,
        behavior: "smooth",
      });
    }
  };

  const scrollLeft = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: -300,
        behavior: "smooth",
      });
    }
  };

  return (
    <div className="w-full py-6 md:px-6">
      {/* Header with title and view all button */}
      <div className="flex items-center justify-between mb-4 p-2 !pl-0 relative">
        <h2 className="header-2 font-bold text-gray-800">
          {categoryExplore.title}
        </h2>
        <div className="absolute right-0 top-[2%]">
          <button className="w-[8rem] h-[3rem] rounded-3xl border border-black bg-white hover:bg-black hover:text-white transition-all duration-300 ease-in-out">
            View All
          </button>
        </div>
        <button
          onClick={categoryExplore.onViewAllClick}
          className="hidden px-4 py-2 border border-gray-300 rounded-full text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          {categoryExplore.viewAllLabel}
        </button>
      </div>

      {/* Category cards container */}
      <div className="relative">
        <div
          ref={scrollContainerRef}
          className="flex overflow-x-auto gap-4 pb-4 -mx-4 px-4 scrollbar-hide"
        >
          {categories.map((category) => (
            <CategoryCard
              category={category}
              itemWidth={categoryExplore.itemWidth}
              itemHeight={categoryExplore.itemHeight}
              onItemClick={categoryExplore.onItemClick}
            ></CategoryCard>
          ))}
        </div>

        {/* Left navigation arrow */}
        {categoryExplore.showNavigationArrow && showLeftArrow && (
          <button
            className="absolute left-0 top-1/2 -translate-y-1/2 bg-white rounded-full shadow-md h-10 w-10 flex items-center justify-center z-10"
            onClick={scrollLeft}
          >
            <MdChevronLeft className="h-5 w-5 text-gray-500" />
          </button>
        )}

        {/* Right navigation arrow */}
        {categoryExplore.showNavigationArrow && showRightArrow && (
          <button
            className="absolute right-0 top-1/2 -translate-y-1/2 bg-white rounded-full shadow-md h-10 w-10 flex items-center justify-center z-10"
            onClick={scrollRight}
          >
            <MdChevronRight className="h-5 w-5 text-gray-500" />
          </button>
        )}

        {/* Left gradient effect */}
        {showLeftArrow && (
          <div className="absolute top-0 -left-1 h-full w-[5rem] bg-gradient-to-r from-white to-transparent pointer-events-none"></div>
        )}

        {/* Right gradient effect */}
        {showRightArrow && (
          <div className="absolute top-0 -right-1 h-full w-[5rem] bg-gradient-to-l from-white to-transparent pointer-events-none"></div>
        )}
      </div>
    </div>
  );
};

export default CategoryExplorer;
