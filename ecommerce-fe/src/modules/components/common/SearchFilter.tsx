import React, { useState, useEffect } from "react";
import { MdOutlineSearch } from "react-icons/md";
import { IoFilterOutline } from "react-icons/io5";
import { PiArrowsDownUp } from "react-icons/pi";
import { HiDotsVertical } from "react-icons/hi";

interface SearchFilterProps {
  onSearch: (searchTerm: string) => void;
  placeholder?: string;
  initialValue?: string;
  showFilterButtons?: boolean;
  onFilterClick?: () => void;
  onSortClick?: () => void;
  onMoreClick?: () => void;
}

/**
 * A reusable search filter component
 * Can be used for any list or table that needs search functionality
 */
const SearchFilter: React.FC<SearchFilterProps> = ({
  onSearch,
  placeholder = "Search",
  initialValue = "",
  showFilterButtons = true,
  onFilterClick,
  onSortClick,
  onMoreClick,
}) => {
  const [searchTerm, setSearchTerm] = useState(initialValue);

  useEffect(() => {
    setSearchTerm(initialValue);
  }, [initialValue]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchTerm);
  };

  return (
    <form onSubmit={handleSubmit} className="flex">
      <div className="w-full md:w-[420px] h-[40px] flex items-center gap-4 md:gap-12">
        <div className="w-full md:w-[264px] h-full relative flex items-center">
          <input
            type="text"
            placeholder={placeholder}
            className="w-full md:w-[264px] py-6 px-8 pl-[1rem] bg-neutral-50 h-[40px] rounded-md focus:outline-none text-neutral-500"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
          <button type="submit" className="absolute right-[5%] top-1/2 -translate-y-1/2">
            <MdOutlineSearch className="text-[25px] text-neutral-500" />
          </button>
        </div>
        
        {showFilterButtons && (
          <>
            <button
              type="button"
              onClick={onFilterClick}
              className="w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring"
            >
              <IoFilterOutline />
            </button>
            <button
              type="button"
              onClick={onSortClick}
              className="w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring"
            >
              <PiArrowsDownUp />
            </button>
            <button
              type="button"
              onClick={onMoreClick}
              className="w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring"
            >
              <HiDotsVertical />
            </button>
          </>
        )}
      </div>
    </form>
  );
};

export default SearchFilter;