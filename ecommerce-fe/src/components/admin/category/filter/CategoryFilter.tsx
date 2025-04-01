import React, { useState, useEffect } from 'react';
import { IoFilterOutline } from "react-icons/io5";
import { PiArrowsDownUp } from "react-icons/pi";
import { HiDotsVertical } from "react-icons/hi";
import { MdOutlineSearch } from "react-icons/md";
import { CategoryFilterCounts } from '../models/category.model';

interface CategoryFilterProps {
  onSearch?: (searchTerm: string) => void;
  onFilterChange: (filter: 'all' | 'featured' | 'onSale' | 'outOfStock') => void;
  activeFilter?: 'all' | 'featured' | 'onSale' | 'outOfStock';
  counts?: CategoryFilterCounts;
  loading?: boolean;
}

export const CategoryFilter: React.FC<CategoryFilterProps> = ({ 
  onSearch, 
  onFilterChange, 
  activeFilter = 'all',
  counts = { all: 0, featured: 0, onSale: 0, outOfStock: 0 },
  loading = false
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSearch) {
      onSearch(searchTerm);
    }
  };
  
  // Map tab names to filter values
  const filterLabels = {
    'all': 'All Product',
    'featured': 'Featured Products',
    'onSale': 'On Sale',
    'outOfStock': 'Out of Stock'
  };
  
  return (
    <div className="flex flex-col md:flex-row justify-between mb-5 ">
      <div className="gap-[14px] p-4 rounded-lg flex items-center space-x-2 mb-3 md:mb-0 overflow-x-auto bg-aqua-spring w-[640px] h-[40px]">
        <button
          onClick={() => onFilterChange('all')}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === 'all'
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          All Product <span className="text-xs ml-1">({loading ? '...' : counts.all})</span>
        </button>
        <button
          onClick={() => onFilterChange('featured')}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === 'featured'
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          Featured Products <span className="text-xs ml-1">({loading ? '...' : counts.featured})</span>
        </button>
        <button
          onClick={() => onFilterChange('onSale')}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === 'onSale'
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          On Sale <span className="text-xs ml-1">({loading ? '...' : counts.onSale})</span>
        </button>
        <button
          onClick={() => onFilterChange('outOfStock')}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === 'outOfStock'
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          Out of Stock <span className="text-xs ml-1">({loading ? '...' : counts.outOfStock})</span>
        </button>
      </div>

      <form onSubmit={handleSearch} className="flex">
        <div className="w-[420px] h-[40px] flex items-center gap-12">
          <div className='w-[264px] h-full relative flex items-center'>
            <input
              type="text"
              placeholder="Search your product"
              className="w-[264px] py-6 px-8 pl-[1rem] bg-neutral-50 h-[40px] rounded-md focus:outline-none text-neutral-500"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            <MdOutlineSearch className='text-[25px] text-neutral-500 absolute right-[5%] top-1/2 -translate-y-1/2'></MdOutlineSearch>
          </div>
          <button type="button" className='w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring'>
            <IoFilterOutline></IoFilterOutline>
          </button>
          <button type="button" className='w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring'>
            <PiArrowsDownUp></PiArrowsDownUp>
          </button>
          <button type="button" className='w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring'>
            <HiDotsVertical></HiDotsVertical>
          </button>
        </div>
      </form>
    </div>
  );
};

export default CategoryFilter;