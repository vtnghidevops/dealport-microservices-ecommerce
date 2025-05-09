import React, { useState, useEffect } from "react";
import {
  TopProductItem,
  TopProductsProps,
} from './models/topProducts.model'
import ProductService from "@/services/product/product.service";
import { handleViewAll } from "../../../utils/helpers";
import { handleProductItemClick } from "../../../utils/helpers";
import ProductCardItem from "./ProductsCard";
import { useNavigate } from "react-router-dom";
import { TopProductSkeleton } from "@/components/ui/skeletons";

const TopProducts: React.FC<TopProductsProps> = () => {
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
  // console.log("topProductsData", topProductsData[0].uiMetadata)


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
    return <TopProductSkeleton />
  }

  // Rearrange products to prioritize specific designs for positions 4 and 5
  const arrangeProductsByDesign = (products: TopProductItem[]): TopProductItem[] => {
    // Make a copy of the products array to avoid modifying the original
    let sortedProducts = [...products];

    // First, handle specific design priorities

    // Find products with specific designs
    const colProducts = sortedProducts.filter(
      (product) => product.uiMetadata.setUpDesign === "col"
    );

    const doubleProducts = sortedProducts.filter(
      (product) => product.uiMetadata.setUpDesign === "double"
    );

    // Remove the col and double products from the original array
    sortedProducts = sortedProducts.filter(
      (product) =>
        product.uiMetadata.setUpDesign !== "col" &&
        product.uiMetadata.setUpDesign !== "double"
    );

    // Create a new array with prioritized positions
    const result: TopProductItem[] = [];

    // Fill positions 1-3 with regular products
    for (let i = 0; i < Math.min(3, sortedProducts.length); i++) {
      result.push(sortedProducts[i]);
    }

    // Position 4 should be col design if available
    if (colProducts.length > 0) {
      result.push(colProducts[0]);
    } else if (sortedProducts.length > 3) {
      // Use next regular product if no col design
      result.push(sortedProducts[3]);
    }

    // Position 5 should be double design if available
    if (doubleProducts.length > 0) {
      result.push(doubleProducts[0]);
    } else if (sortedProducts.length > 4) {
      // Use next regular product if no double design
      result.push(sortedProducts[4]);
    }

    // Add a 6th product if available and we don't have 6 yet
    if (result.length < 6) {
      const remainingProducts = sortedProducts.filter(
        (_, index) => index > (result.length - (colProducts.length + doubleProducts.length) - 1)
      );

      if (remainingProducts.length > 0) {
        result.push(remainingProducts[0]);
      }
    }

    // Ensure we limit to exactly 6 products maximum
    return result.slice(0, 6);
  };

  // Apply the arrangement
  const arrangedProducts = arrangeProductsByDesign(topProducts.products);

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
        {arrangedProducts.map((product, index) => {
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
