import React, {useEffect, useState} from "react";
import { getButtonClass } from "../../common/Button";
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
        {trendingProducts.map((product) => (
          <ProductCard
            key={product.id}
            title={product.title}
            description={product.description}
            price={product.price}
            originalPrice={product.originalPrice}
            discount={product.discount}
            reviews={product.review}
            imageUrl={product.imageUrl}
          />
        ))}
        {/* Collection for Men */}
        <div>
          <MenCollection products={menCollection}></MenCollection>
        </div>
      </div>
      <div className="absolute right-[5%] top-0">
        <button
          className={`${getButtonClass(
            ""
          )} w-[8rem] h-[3rem] rounded-3xl border border-black`}
        >
          View All
        </button>
      </div>
    </div>
  );
};

export default TrendingProducts;