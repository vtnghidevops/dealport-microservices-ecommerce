import React, {useEffect, useState} from "react";
import ProductCard from "../../common/Card";
import MenCollection from "./MenCollection";
import { TrendingProductItem } from "./models/trendingProducts.model.ts";
import { TrendingPorductService } from "./services/trendingProducts.service.ts";
import { MenCollectionItem } from "./models/trendingProducts.model.ts";
import { MenCollectionService } from "./services/menCollection.service.ts";

const TrendingProducts: React.FC = () => {
  
  const [trendingProducts, setTrendingProducts] = useState<TrendingProductItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [menCollection, setMenCollection] = useState<MenCollectionItem[]>([])
  
    // fetch data trending Products
    useEffect(() => {
      const fetchTestimonials = async () => {
        try {
          const data = await TrendingPorductService.getDataTredings();
          setTrendingProducts(data);
          setLoading(false);
        } catch (error) {
          console.error("Error fetching testimonials:", error);
          setLoading(false);
        }
      };
  
      fetchTestimonials();
    }, []);
    // fetch data collection
    useEffect(() => {
      const fetchMenCollection = async () => {
        try {
          const data = await MenCollectionService.getDataMen();
          setMenCollection(data);
          setLoading(false);
        } catch (error) {
          console.error("Error fetching testimonials:", error);
          setLoading(false);
        }
      };
      fetchMenCollection();
    }, []);
    


  return (
    <div className="h-full">
      <span className="header-2 font-bold p-[10px] block mb-3">
        Trending Products
      </span>
      <div className="h-full w-full flex justify-start items-start">
        {trendingProducts.map((product, index) => (
          <ProductCard
            key={index}
            id={product.id}
            title={product.title}
            description={product.description}
            price={Number(product.price)}
            originalPrice={Number(product.originalPrice) || 0}
            discount={Number(product.discount) || 0}
            reviews={product.review}
            image_url={product.image_url}
          />
        ))}
        {/* Collection for Men */}
        <div>
          <MenCollection products={menCollection}></MenCollection>
        </div>
      </div>

      <div className="absolute right-[5%] top-[2%]">
        <button className="w-[8rem] h-[3rem] rounded-3xl border border-black bg-white hover:bg-black hover:text-white transition-all duration-300 ease-in-out">
          View All
        </button>
      </div>
    </div>
  );
};

export default TrendingProducts;