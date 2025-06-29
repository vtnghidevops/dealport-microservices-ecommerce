import React, { useState } from 'react';
import { IoFilterOutline } from "react-icons/io5";
import { PiArrowsDownUp } from "react-icons/pi";
import { HiDotsVertical } from "react-icons/hi";
import { MdOutlineSearch } from "react-icons/md";


interface OrderFilterProps {
  onSearch: (searchTerm: string) => void;
  onFilterChange: (filter: string) => void;
  counts: {
    all: number;
    paid: number;
    pending: number;
    shipped: number;
    cancelled: number;
    processing?: number;
    delivered?: number;
    refunded?: number;
  };
  loading?: boolean;
}

export const OrderFilter: React.FC<OrderFilterProps> = ({
  onSearch,
  onFilterChange,
  counts = { all: 0, paid: 0, pending: 0, shipped: 0, cancelled: 0 },
  loading = false,
}) => {
  const [activeFilter, setActiveFilter] = useState<string>("All");
  const [searchTerm, setSearchTerm] = useState("");

  // Ensure counts are valid numbers
  const safeCounts = {
    all: typeof counts.all === 'number' ? counts.all : 0,
    paid: typeof counts.paid === 'number' ? counts.paid : 0,
    pending: typeof counts.pending === 'number' ? counts.pending : 0,
    shipped: typeof counts.shipped === 'number' ? counts.shipped : 0,
    cancelled: typeof counts.cancelled === 'number' ? counts.cancelled : 0,
    processing: typeof counts.processing === 'number' ? counts.processing : 0,
    delivered: typeof counts.delivered === 'number' ? counts.delivered : 0,
    refunded: typeof counts.refunded === 'number' ? counts.refunded : 0,
  };

  const handleFilterChange = (filter: string) => {
    setActiveFilter(filter);
    onFilterChange(filter);
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchTerm);
  };

  // Handle Enter key press in search input
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      onSearch(searchTerm);
    }
  };

  const filterButtons = [
    { name: "All", count: safeCounts.all },
    { name: "Pending", count: safeCounts.pending },
    { name: "Processing", count: safeCounts.processing },
    { name: "Paid", count: safeCounts.paid },
    { name: "Shipped", count: safeCounts.shipped },
    { name: "Delivered", count: safeCounts.delivered },
    { name: "Cancelled", count: safeCounts.cancelled },
    { name: "Refunded", count: safeCounts.refunded },
  ];

  return (
    <>
      <div className="mb-4 flex justify-between">
        <form onSubmit={handleSearch} className="relative flex items-center w-[326px] h-[48px]">
          <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
            <MdOutlineSearch className="w-[20px] h-[20px] text-gray-500" />
          </div>
          <div className='w-full border border-gray-300 rounded-lg '>
            <input
              type="search"
              className="w-[80%] h-[48px] p-4 pl-10 text-sm text-gray-900 outline-none"
              placeholder="Search"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              onKeyDown={handleKeyDown}
              disabled={loading}
            />
            <button
              type="submit"
              className="text-white absolute end-2.5 !px-5 bottom-1.5 top-1.5 bg-ocean-green hover:bg-green-700 font-medium rounded text-sm py-2"
              disabled={loading}
            >
              Search
            </button>
          </div>

        </form>

        <div className="flex items-center">
          <button className="w-[48px] h-[48px] rounded-lg border border-neutral-200 flex justify-center items-center mr-2 bg-white">
            <IoFilterOutline className="w-[20px] h-[20px] text-gray-500" />
          </button>
          <button className="w-[48px] h-[48px] rounded-lg border border-neutral-200 flex justify-center items-center mr-2 bg-white">
            <PiArrowsDownUp className="w-[20px] h-[20px] text-gray-500" />
          </button>
          <button className="w-[48px] h-[48px] rounded-lg border border-neutral-200 flex justify-center items-center bg-white">
            <HiDotsVertical className="w-[20px] h-[20px] text-gray-500" />
          </button>
        </div>
      </div>

      <div className="mb-6 mt-5 flex flex-wrap">
        {filterButtons.map((button) => (
          <button
            key={button.name}
            className={`mr-6 mb-2 px-5 py-2 rounded-lg text-[15px] font-medium ${activeFilter === button.name
              ? "bg-ocean-green text-white"
              : "bg-gray-100 text-gray-500 hover:bg-gray-200"
              }`}
            onClick={() => handleFilterChange(button.name)}
            disabled={loading}
          >
            {button.name} ({button.count})
          </button>
        ))}
      </div>
    </>
  );
};