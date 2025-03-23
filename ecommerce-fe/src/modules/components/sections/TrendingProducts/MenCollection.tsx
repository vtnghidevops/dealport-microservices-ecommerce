import React from "react";
import MenCard from "./MenCard";
import { MenCollectionItem } from "./models/trendingProducts.model";

const MenCollection: React.FC<{products: MenCollectionItem[]}> = ({products}) => {
  return (
    <div className="w-[392px] min-h-[435px] mx-[0.5rem] p-6 rounded-[12px] border border-grep-300 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.1),0_4px_6px_-2px_rgba(0,0,0,0.05)] ">
      <h2 className="title mb-4 p-[10px]">Trend collection for men</h2>
      <div className="flex flex-wrap gap-5 mx-3">
        {products.map((product) => (
          <MenCard key={product.id} {...product}></MenCard>
        ))}
      </div>
    </div>
  );
};
export default MenCollection