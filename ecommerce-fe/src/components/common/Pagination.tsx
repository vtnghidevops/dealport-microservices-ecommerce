// components/common/Pagination.tsx
import React from 'react';
import { IoMdArrowRoundBack, IoMdArrowRoundForward } from "react-icons/io";

interface PaginationProps {
  currentPage: number;
  totalItems: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  showNavigation?: boolean;
}

/**
 * Reusable pagination component for tables and lists
 */
const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalItems,
  pageSize,
  onPageChange,
  showNavigation = true,
}) => {
  const totalPages = Math.ceil(totalItems / pageSize);
  const pages: React.ReactNode[] = [];

  // Previous button
  if (showNavigation) {
    pages.push(
      <button
        key="previous"
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className={`absolute left-0 px-4 py-2 w-[120px] h-[45px] justify-center flex items-center border rounded-xl transition-all duration-200 ${
          currentPage === 1
            ? "bg-gray-100 text-gray-400 cursor-not-allowed"
            : "bg-white text-[16px] font-medium shadow-md hover:bg-gray-50 hover:border-gray-300 active:bg-gray-100"
        }`}
      >
        <IoMdArrowRoundBack className="mr-2 h-[20px] w-[20px]" />
        Previous
      </button>
    );
  }

  // Page numbers
  for (let i = 1; i <= Math.min(5, totalPages); i++) {
    pages.push(
      <button
        key={i}
        className={`px-3 py-1 rounded-md w-[36px] h-[36px] transition-all duration-200 ${
          currentPage === i 
            ? "bg-[#C1E6BA] text-black font-medium border-2 border-[#A3D89A] hover:bg-[#B1E0A8]" 
            : "border-2 border-aqua-spring hover:border-[#A3D89A] hover:bg-[#F0F9EE] active:bg-[#E1F4DD]"
        }`}
        onClick={() => onPageChange(i)}
      >
        {i}
      </button>
    );
  }

  // Ellipsis and last page for pagination with many pages
  if (totalPages > 5) {
    pages.push(<span key="ellipsis">...</span>);
    pages.push(
      <button
        key={totalPages}
        className={`px-3 py-1 rounded-md w-[36px] h-[36px] transition-all duration-200 ${
          currentPage === totalPages 
            ? "bg-[#C1E6BA] text-black font-medium border-2 border-[#A3D89A] hover:bg-[#B1E0A8]" 
            : "border-2 border-aqua-spring hover:border-[#A3D89A] hover:bg-[#F0F9EE] active:bg-[#E1F4DD]"
        }`}
        onClick={() => onPageChange(totalPages)}
      >
        {totalPages}
      </button>
    );
  }

  // Next button
  if(showNavigation){
    pages.push(
      <button
        key="next"
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className={`absolute right-0 px-4 py-2 w-[120px] h-[45px] justify-center flex items-center border rounded-xl transition-all duration-200 ${
          currentPage === totalPages
            ? "bg-gray-100 text-gray-400 cursor-not-allowed"
            : "bg-white text-[16px] font-medium shadow-md hover:bg-gray-50 hover:border-gray-300 active:bg-gray-100"
        }`}
      >
        Next
        <IoMdArrowRoundForward className="ml-2 h-[24px] w-[24px]" />
      </button>
    );
  }
  
  return (
    <div className="relative flex justify-center items-center space-x-2 mt-5">
      {pages}
    </div>
  );
};

export default Pagination;