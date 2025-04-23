import React from "react";
import { getButtonClass } from "../../../utils/buttonUtils";
import { Product } from "@/types/product.model";

const MenCard: React.FC<Product> = ({ image_url, price, discount }) => {
  return (
    <div
      className="bg-white rounded-xl overflow-hidden relative w-[166px] h-[171px] 
  shadow-[0_-2px_4px_rgba(0,0,0,0.05),2px_0_5px_rgba(0,0,0,0.05),-2px_0_5px_rgba(0,0,0,0.05),0_4px_6px_rgba(0,0,0,0.2)]"
    >
      {(discount != 0) && (
        <div className="absolute top-4 left-4">
          <span className="font-medium text-[12px] p-[10px] z-10">
            {discount}% Off
          </span>
        </div>
      )}
      <a href="#" className="flex justify-center items-center w-full h-full">
        <img
          src={image_url}
          alt="Black pants"
          className="w-full h-full object-contain "
        />
      </a>
      <div className="absolute top-4 right-0 ">
        {(price != 0 )&& (
          <span className="text-ocean-green p-[10px] text-[12px] z-10">
            ${price}
          </span>
        )}
      </div>

      <div className="mt-2 absolute bottom-[7%] w-full">
        <div className="rounded-3xl w-full flex justify-center items-center ">
          <button
            className={`${getButtonClass(
              "secondary"
            )} flex justify-center items-center !text-[10px] w-[90px] h-[20px] `}
          >
            Buy Now
          </button>
        </div>
      </div>
    </div>
  );
};
export default MenCard;