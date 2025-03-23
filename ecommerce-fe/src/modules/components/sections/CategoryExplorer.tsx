import React, { useRef, useState, useEffect } from 'react';
import { MdChevronRight, MdChevronLeft } from 'react-icons/md';
import { CategoryExplorerProps } from '../../models/Category';

const CategoryExplorer: React.FC<CategoryExplorerProps> = ({
  title = "Start exploring now",
  categories = [],
  viewAllLabel = "View All",
  onViewAllClick = () => {},
  itemWidth = "w-[180px]",
  itemHeight = "h-[220px]",
  showNavigationArrow = true,
  onItemClick = () => {},
}) => {
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
    const hasMoreToScroll = container.scrollWidth > container.clientWidth + container.scrollLeft + 20;
    setShowRightArrow(hasMoreToScroll);
  };

  // Initial check and add scroll listener
  useEffect(() => {
    checkForArrows();
    const container = scrollContainerRef.current;
    
    if (container) {
      container.addEventListener('scroll', checkForArrows);
      // Check on window resize as well
      window.addEventListener('resize', checkForArrows);
    }
    
    return () => {
      if (container) {
        container.removeEventListener('scroll', checkForArrows);
      }
      window.removeEventListener('resize', checkForArrows);
    };
  }, [categories]);

  const scrollRight = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: 300,
        behavior: 'smooth'
      });
    }
  };

  const scrollLeft = () => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: -300,
        behavior: 'smooth'
      });
    }
  };

  return (
    <div className="w-full py-6 md:px-6">
      {/* Header with title and view all button */}
      <div className="flex items-center justify-between mb-4 p-2 !pl-0 relative">
        <h2 className="header-2 font-bold text-gray-800">{title}</h2>
        <div className="absolute right-0 top-0">
          <button
            onClick={onViewAllClick}
            className="buttonText w-[8rem] h-[3rem] rounded-3xl border border-black"
          >
            View All
          </button>
        </div>
        <button
          onClick={onViewAllClick}
          className="hidden px-4 py-2 border border-gray-300 rounded-full text-sm font-medium text-gray-700 hover:bg-gray-50"
        >
          {viewAllLabel}
        </button>
      </div>

      {/* Category cards container */}
      <div className="relative">
        <div
          ref={scrollContainerRef}
          className="flex overflow-x-auto gap-4 pb-4 -mx-4 px-4 scrollbar-hide"
        >
          {categories.map((category, index) => (
            <div
              key={`category-${index}-${category.id || category.name}`}
              className={`flex-shrink-0 ${itemWidth} cursor-pointer h-[220px] w-[180px mr-[3rem] rounded-xl`}
              onClick={() => onItemClick(category, index)}
            >
              <div className="relative rounded-lg overflow-hidden shadow-sm border border-gray-300 h-full ">
                <div
                  className={`${itemHeight} w-full flex items-center justify-center`}
                >
                  <img
                    src={category.image}
                    alt={category.name}
                    className="w-[148px] h-[140px] object-cover"
                  />
                </div>
                <div className="p-2 text-center absolute bottom-0 flex justify-center w-full">
                  <h3 className="text-sm font-medium text-gray-800">
                    {category.name}
                  </h3>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Left navigation arrow */}
        {showNavigationArrow && showLeftArrow && (
          <button
            className="absolute left-0 top-1/2 -translate-y-1/2 bg-white rounded-full shadow-md h-10 w-10 flex items-center justify-center z-10"
            onClick={scrollLeft}
          >
            <MdChevronLeft className="h-5 w-5 text-gray-500" />
          </button>
        )}

        {/* Right navigation arrow */}
        {showNavigationArrow && showRightArrow && (
          <button
            className="absolute right-0 top-1/2 -translate-y-1/2 bg-white rounded-full shadow-md h-10 w-10 flex items-center justify-center z-10"
            onClick={scrollRight}
          >
            <MdChevronRight className="h-5 w-5 text-gray-500" />
          </button>
        )}
        
        {/* Left gradient effect */}
        {showLeftArrow && (
          <div className="absolute top-0 -left-1 h-full w-[5rem] bg-gradient-to-r from-white to-transparent pointer-events-none">
          </div>
        )}
        
        {/* Right gradient effect */}
        {showRightArrow && (
          <div className="absolute top-0 -right-1 h-full w-[5rem] bg-gradient-to-l from-white to-transparent pointer-events-none">
          </div>
        )}
      </div>
    </div>
  );
};

export default CategoryExplorer;