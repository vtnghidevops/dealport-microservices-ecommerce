import React, { useState, useEffect } from "react";
import {
  TopProductItem,
  TopProductsProps,
} from './models/topProducts.model'
import ProductService from "@/services/product.service";
import { handleViewAll } from "../../../utils/helpers";
import { handleProductItemClick } from "../../../utils/helpers";
import ProductCardItem from "./ProductsCard";
import Loading from '@/components/shared/Loading'
import { useNavigate } from "react-router-dom";
const TopProducts: React.FC  = () => {
  // fetch data
  const [topProductsData, setTopProductsData] = useState<TopProductItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
        const fetchTopProductsData = async () => {
          try {
            const data = await ProductService.getTopSaleProducts();
            setTopProductsData(data);
            setLoading(false);
          } catch (error) {
            console.error("Error fetching testimonials:", error);
            setLoading(false);
          }
        };
    
        fetchTopProductsData();
      }, []);


    
  var topProducts: TopProductsProps = {
    title: "Best selling product",
    products: topProductsData,
    viewAllLabel: "View All",
    onViewAllClick: () => handleViewAll(useNavigate()), // navigate to /products
    onItemClick: (product: TopProductItem, index: number) => handleProductItemClick(useNavigate(), product, index), // Sử dụng hàm đã định nghĩa ở trên
  };


  const defaultGridPositions = [
    "col-span-1 row-span-1", // Sub1
    "col-span-1 row-span-1", // Sub2
    "col-span-1 row-span-1", // Sub3
    "col-span-1 row-span-2", // Sub4
    "col-span-2 row-span-1 row-start-2", // Sub5
    "col-span-1 row-span-1 row-start-2", // Sub6
  ];

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="w-full py-6 md:px-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-[2rem]">
        <h2 className="header-2 font-bold">{topProducts.title}</h2>
        <button
          onClick={topProducts.onViewAllClick}
          className="w-[8rem] h-[3rem] rounded-3xl border border-black bg-white hover:bg-black hover:text-white transition-all duration-300 ease-in-out"
        >
          {topProducts.viewAllLabel}
        </button>
        
      </div>

      {/* Grid layout */}
      
      <div className="grid grid-cols-4 gap-4">
        {topProducts.products.map((product, index) => {
          // Ưu tiên gridSpan từ dữ liệu, nếu không có thì dùng mặc định
          const gridPosition = defaultGridPositions[index] || "col-span-1 row-span-1";
          // console.log("product in top", product)
          return (
            <div key={product.id} className={`${gridPosition}`}>
              <ProductCardItem
                product={product}
                onClick={() => topProducts.onItemClick?.(product, index)}
              />
            </div>
          );
        })}
      </div>
    </div>
  );

};

export default TopProducts;
