import React from "react";
import { Link } from "react-router-dom";
import { ProductCard } from "./models/topProducts.model";
// import { RxArrowTopRight } from "react-icons/rx"; 


// const buttonClasses =
//   "flex justify-center bg-white items-center absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl hover:bg-aqua-spring hover:shadow-md transition-all duration-300 ease-out";

// Product Card Component
const ProductCardItem: React.FC<ProductCard> = ({ product, onClick }) => {
  const {
    name,
    price,
    imageUrl,
    // actionLabel,
    uiMetadata,
    categorySlug,
    slug
  } = product;
  if (uiMetadata.setUpDesign == "row") {
    if (uiMetadata.isCommingSoon) {
      return (
        <div className="rounded-xl mr-[0.5rem] mt-[0.5rem] border border-gray-200">
          <Link to={`/category/${categorySlug}/${slug}`}>
            <img src={imageUrl} alt={name} className="rounded-xl h-[200px] w-[310px]"></img>
          </Link>
        </div>
      );
    } else {
      return (
        <div className="rounded-xl relative mr-[0.5rem] mt-[0.5rem]  border border-gray-200">
          <Link to={`/category/${categorySlug}/${slug}`}>
            <img src={imageUrl} alt={name} className="rounded-xl h-[200px] w-[310px]"></img>
          </Link>
          <span className="absolute bottom-0 p-2 text-white title font-bold right-5">
            ${price}
          </span>
        </div>
      );
    }
  } else if (uiMetadata.setUpDesign == "col") {
    return (
      <div className="rounded-xl mr-[0.5rem] relative mt-[0.5rem] ">
        <Link to={`/category/${categorySlug}/${slug}`} className="block h-full">
          <img
            src={imageUrl}
            alt={name}
            className="rounded-xl w-[310px] h-[420px]"
          ></img>
        </Link>
      </div>
    );
  } else {
    return (
      <div className="flex rounded-xl mr-[0.5rem] items-center relative mt-[0.5rem] ">
        <Link to={`/category/${categorySlug}/${slug}`} className="flex max-w-[39rem]">
          <img
            src={imageUrl}
            alt={name}
            className="aspect-[16/9] h-[200px] w-[630px] rounded-xl"
          />
        </Link>
      
      </div>
    );
  }
};

export default ProductCardItem;

