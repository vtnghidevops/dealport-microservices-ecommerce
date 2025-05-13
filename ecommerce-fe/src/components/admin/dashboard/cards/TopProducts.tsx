import React, { useState } from "react";
import { BestSellingProductStats } from "../services/product.service";

// interface Product {
//   id: number;
//   name: string;
//   category: string;
//   price: number;
//   imageUrl: string;
//   itemCode: string;
// }

interface TopProductsTableProps {
  products: BestSellingProductStats[];
}

const TopProductsTable: React.FC<TopProductsTableProps> = ({ products }) => {
  const [searchTerm, setSearchTerm] = useState("");

  const filteredProducts = products.filter((product) =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="bg-white rounded-lg filter drop-shadow-lg h-[420px] w-[360px] px-[16px] py-[20px]">
      <div className="flex justify-between items-center mb-3">
        <h2 className="text-[18px] font-bold text-gray-900">Top Products</h2>
        <a href="#" className="text-primary text-sm ">
          All
        </a>
      </div>

      <div className="relative mb-6 p-8 rounded-lg">
        <div className="absolute inset-y-0 left-0 flex items-center pl-[1.25rem] pointer-events-none">
          <svg
            className="w-[1rem] h-[1rem] text-gray-700"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            ></path>
          </svg>
        </div>
        <input
          type="search"
          className="w-full pl-10 py-2 bg-neutral-50 rounded-md text-sm text-gray-700 focus:outline-none"
          placeholder="Search"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
      </div>

      <div className="space-y-4">
        {filteredProducts.map((product) => (
          <div key={product.id} className="border-b border-gray-100 pb-4">
            <div className="flex items-center justify-between h-[64px] w-[328px]">
              <div className="flex items-center gap-[16px] w-full ">
                <div className="relative bg-gray-100 rounded-md overflow-hidden w-[56px] h-[56px]">
                  <img
                    src={product.imageSrc}
                    alt={product.name}
                    className="object-cover w-[56px] h-[56px]"
                  />
                </div>
                <div className="w-[180px] h-[46px] flex flex-col justify-center gap-[4px]">
                  <h3 className="text-[15px] font-medium text-cyprus">
                    {product.name}
                  </h3>
                  <p className="text-[12px] text-[#8B909A]">
                    Sold: {product.sold}
                  </p>
                </div>
                <div className="text-right ml-[0.5rem]">
                  <span className="text-[15px] font-bold text-gray-800">
                    ${product.price.toFixed(2)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TopProductsTable;
