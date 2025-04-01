import React, { useState, useEffect } from 'react';
import { IoFilterOutline } from "react-icons/io5";
import { PiArrowsDownUp } from "react-icons/pi";
import { HiDotsVertical } from "react-icons/hi";
import { MdOutlineSearch } from "react-icons/md";
import { orderService } from "../services/order.service";

interface OrderFilterProps {
  onSearch: (searchTerm: string) => void;
  onFilterChange: (filter: string) => void;
  counts?: {
    all: number;
    completed: number;
    pending: number;
    shipped: number;
    cancelled: number;
  };
  loading?: boolean;
}

export const OrderFilter: React.FC<OrderFilterProps> = ({ 
  onSearch, 
  onFilterChange, 
  counts = { all: 0, completed: 0, pending: 0, shipped: 0, cancelled: 0 },
  loading = false
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [activeFilter, setActiveFilter] = useState('All order');
  const [orderCounts, setOrderCounts] = useState(counts);
  const [isLoading, setIsLoading] = useState(loading);
  
  // Fetch order counts when component mounts
  useEffect(() => {
    const fetchOrderCounts = async () => {
      setIsLoading(true);
      try {
        // Get all orders to calculate counts
        const { total: all } = await orderService.fetchOrders({ page: 1, limit: 1 });
        const { total: completed } = await orderService.fetchOrders({ status: 'Delivered', page: 1, limit: 1 });
        const { total: pending } = await orderService.fetchOrders({ status: 'Pending', page: 1, limit: 1 });
        const { total: shipped } = await orderService.fetchOrders({ status: 'Shipped', page: 1, limit: 1 });
        const { total: cancelled } = await orderService.fetchOrders({ status: 'Cancelled', page: 1, limit: 1 });
        
        setOrderCounts({
          all,
          completed,
          pending,
          shipped,
          cancelled
        });
      } catch (error) {
        console.error("Failed to fetch order counts:", error);
      } finally {
        setIsLoading(false);
      }
    };

    // Only fetch if counts are not provided as props
    if (counts.all === 0) {
      fetchOrderCounts();
    } else {
      setOrderCounts(counts);
    }
  }, [counts]);
  
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(searchTerm);
  };
  
  const handleFilterClick = (filter: string) => {
    setActiveFilter(filter);
    onFilterChange(filter);
  };
  
  return (
    <div className="flex flex-col md:flex-row justify-between mb-5 ">
      <div className="gap-[14px] p-4 rounded-lg flex items-center space-x-2 mb-3 md:mb-0 overflow-x-auto bg-aqua-spring w-[540px] h-[40px]">
        <button
          onClick={() => handleFilterClick("All order")}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === "All order"
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          All order <span className="text-xs ml-1">({isLoading ? '...' : orderCounts.all})</span>
        </button>
        <button
          onClick={() => handleFilterClick("Completed")}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === "Completed"
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          Completed <span className="text-xs ml-1">({isLoading ? '...' : orderCounts.completed})</span>
        </button>
        <button
          onClick={() => handleFilterClick("Pending")}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === "Pending"
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          Pending <span className="text-xs ml-1">({isLoading ? '...' : orderCounts.pending})</span>
        </button>
        <button
          onClick={() => handleFilterClick("Cancelled")}
          className={`px-[12px] py-[6px] h-[32px] rounded-md whitespace-nowrap ${
            activeFilter === "Cancelled"
              ? "bg-white text-black"
              : "bg-aqua-spring text-neutral-600"
          }`}
        >
          Cancelled <span className="text-xs ml-1">({isLoading ? '...' : orderCounts.cancelled})</span>
        </button>
      </div>

      <form onSubmit={handleSearch} className="flex">
        <div className="w-[420px] h-[40px] flex items-center gap-12">
          <div className='w-[264px] h-full relative flex items-center'>
            <input
              type="text"
              placeholder="Search order report"
              className="w-[264px] py-6 px-8 pl-[1rem] bg-neutral-50 h-[40px] rounded-md focus:outline-none text-neutral-500"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
            <MdOutlineSearch className='text-[25px] text-neutral-500 absolute right-[5%] top-1/2 -translate-y-1/2'></MdOutlineSearch>
          </div>
          <button className='w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring'>
            <IoFilterOutline></IoFilterOutline>
          </button>
          <button className='w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring'>
            <PiArrowsDownUp></PiArrowsDownUp>
          </button>
          <button className='w-[40px] h-[40px] flex items-center justify-center border border-neutral-200 rounded-lg hover:bg-aqua-spring'>
            <HiDotsVertical></HiDotsVertical>
          </button>
        </div>
      </form>
    </div>
  );
};
// Updated OrderFilter component to use the reusable components
// import React from 'react';
// import TabFilter from '../../../common/TabFilter';
// import SearchFilter from '../../../common/SearchFilter';

// interface OrderFilterProps {
//   onSearch: (searchTerm: string) => void;
//   onFilterChange: (filter: string) => void;
//   counts?: {
//     all: number;
//     completed: number;
//     pending: number;
//     shipped: number;
//     cancelled: number;
//   };
//   activeFilter: string;
//   loading?: boolean;
// }

// export const OrderFilter: React.FC<OrderFilterProps> = ({ 
//   onSearch, 
//   onFilterChange, 
//   counts = { all: 0, completed: 0, pending: 0, shipped: 0, cancelled: 0 },
//   activeFilter = "All order",
//   loading = false
// }) => {
//   // Filter tabs configuration
//   const filterTabs = [
//     { id: "All order", label: "All order", count: counts.all },
//     { id: "Completed", label: "Completed", count: counts.completed },
//     { id: "Pending", label: "Pending", count: counts.pending },
//     { id: "Cancelled", label: "Cancelled", count: counts.cancelled },
//   ];
  
//   return (
//     <div className="flex flex-col md:flex-row justify-between mb-5">
//       <TabFilter
//         tabs={filterTabs}
//         activeTab={activeFilter}
//         onTabChange={onFilterChange}
//         loading={loading}
//         containerClassName="bg-aqua-spring w-[540px] h-[40px]"
//       />

//       <SearchFilter
//         onSearch={onSearch}
//         placeholder="Search order report"
//         showFilterButtons={true}
//       />
//     </div>
//   );
// };