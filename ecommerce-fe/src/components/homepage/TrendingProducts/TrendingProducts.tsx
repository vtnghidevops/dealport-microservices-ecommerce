import React, { useEffect, useState } from "react";
import ProductCard from "../../common/Card";
import MenCollection from "./MenCollection";
import ProductService from "@/services/product/product.service";
import { Product } from "@/types/product.model.ts";
import { useNavigate } from "react-router-dom";
import { TopProductSkeleton } from "@/components/ui/skeletons";

const TrendingProducts: React.FC = () => {

  const [trendingProducts, setTrendingProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [menCollection, setMenCollection] = useState<Product[]>([])
  const navigate = useNavigate();

  // Maximum number of products to display
  const MAX_TRENDING_PRODUCTS = 3;
  const MAX_MEN_COLLECTION_PRODUCTS = 4;

  // Category slug for navigation
  const TRENDING_CATEGORY_SLUG = "trending";

  // fetch data trending Products
  useEffect(() => {
    const fetchTrendingProducts = async () => {
      try {
        // Get trending products and limit to MAX_TRENDING_PRODUCTS
        const data = await ProductService.getTrendingProducts(1, MAX_TRENDING_PRODUCTS);
        setTrendingProducts(data.slice(0, MAX_TRENDING_PRODUCTS));
        setLoading(false);
      } catch (error) {
        console.error("Error fetching trending products:", error);
        setLoading(false);
      }
    };

    fetchTrendingProducts();
  }, []);

  // fetch data collection
  useEffect(() => {
    const fetchMenCollection = async () => {
      try {
        // Get men collection products and limit to MAX_MEN_COLLECTION_PRODUCTS
        const data = await ProductService.getMenCollection(1, MAX_MEN_COLLECTION_PRODUCTS);
        setMenCollection(data.slice(0, MAX_MEN_COLLECTION_PRODUCTS));
        setLoading(false);
      } catch (error) {
        console.error("Error fetching men collection:", error);
        setLoading(false);
      }
    };
    fetchMenCollection();
  }, []);

  // Handle View All click
  const handleViewAllClick = () => {
    // Navigate to trending products page
    navigate(`/${TRENDING_CATEGORY_SLUG}`);
  };

  if (loading) {
    return <TopProductSkeleton count={5} />;
  }
  return (
    <div className="h-full relative">
      <span className="header-2 font-bold p-[10px] block mb-3">
        Trending Products
      </span>
      <div className="h-full w-full flex justify-start items-start">
        {trendingProducts.map((product, index) => (
          <ProductCard
            key={index}
            product={product}
          />
        ))}
        {/* Collection for Men */}
        <div>
          <MenCollection products={menCollection}></MenCollection>
        </div>
      </div>

      <div className="absolute right-[5%] top-[2%]">
        <button
          className="w-[8rem] h-[3rem] rounded-3xl border border-black bg-white hover:bg-black hover:text-white transition-all duration-300 ease-in-out"
          onClick={handleViewAllClick}
        >
          View All
        </button>
      </div>
    </div>
  );
};

export default TrendingProducts;