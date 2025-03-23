import React from "react";
import {
  ProductCardProps,
  TopProductsProps,
} from "../../models/featuredProduct.models";
import { RxArrowTopRight } from "react-icons/rx";

const buttonClasses =
  "absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-white flex items-center justify-center";
const TopProducts: React.FC<TopProductsProps> = ({
  title = "Best selling product",
  products = [],
  viewAllLabel = "View All",
  onViewAllClick = () => {},
  onItemClick = () => {},
}) => {
  const defaultGridPositions = [
    "col-span-1 row-span-1", // Sub1
    "col-span-1 row-span-1", // Sub2
    "col-span-1 row-span-1", // Sub3
    "col-span-1 row-span-2", // Sub4
    "col-span-2 row-span-1 row-start-2", // Sub5
    "col-span-1 row-span-1 row-start-2", // Sub6
  ];

  return (
    <div className="w-full py-6 md:px-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-[2rem]">
        <h2 className="text-2xl font-bold">{title}</h2>
        <button
          onClick={onViewAllClick}
          className="w-[8rem] h-[3rem] rounded-3xl border border-black text-sm"
        >
          {viewAllLabel}
        </button>
      </div>

      {/* Grid layout */}
      <div className="grid grid-cols-4 gap-4">
        {products.map((product, index) => {
          // Ưu tiên gridSpan từ dữ liệu, nếu không có thì dùng mặc định
          const gridPosition =
            product.gridSpan?.col && product.gridSpan?.row
              ? `col-span-${product.gridSpan.col} row-span-${product.gridSpan.row} ${
                  index >= 4 ? "row-start-2" : ""
                }`
              : defaultGridPositions[index] || "col-span-1 row-span-1";

          return (
            <div key={product.id} className={`${gridPosition}`}>
              <ProductCard
                product={product}
                onClick={() => onItemClick(product, index)}
              />
            </div>
          );
        })}
      </div>
    </div>
  );

};

// Product Card Component
const ProductCard: React.FC<ProductCardProps> = ({ product, onClick }) => {
  const {
    name,
    price,
    image,
    actionLabel,
    setUpDesign,
    isCommingSoon,
    image_double,
  } = product;
  if (setUpDesign == "row") {
    if (isCommingSoon) {
      return (
        <div className="rounded-xl mr-[0.5rem] mt-[0.5rem] border border-gray-200">
          <a href="#">
            <img src={image} alt={name} className="rounded-xl h-[200px] w-[310px]"></img>
          </a>
        </div>
      );
    } else {
      return (
        <div className="rounded-xl relative mr-[0.5rem] mt-[0.5rem]  border border-gray-200">
          <a href="#">
            <img src={image} alt={name} className="rounded-xl h-[200px] w-[310px]"></img>
          </a>
          <span className="absolute bottom-0 p-2 text-white title font-bold right-5">
            ${price}
          </span>
          {actionLabel === "Visit store" ? (
            <button className="absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-white  flex items-center justify-center ">
              <span className="flex justify-center items-center ml-[-0.5rem]">
                {actionLabel}
              </span>
              <div className="bg-white border border-black w-[1.2rem] h-[1.2rem] rounded-[50%] absolute flex items-center justify-center right-2">
                <RxArrowTopRight></RxArrowTopRight>
              </div>
            </button>
          ) : (
            <button className="absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-white  flex items-center justify-center ">
              <span className="flex justify-center items-center">
                {actionLabel}
              </span>
            </button>
          )}
        </div>
      );
    }
  } else if (setUpDesign == "col") {
    return (
      <div className="rounded-xl mr-[0.5rem] relative mt-[0.5rem] ">
        <a href="#" className="block h-full">
          <img
            src={image}
            alt={name}
            className="rounded-xl w-[310px] h-[420px]"
          ></img>
        </a>
        <button className="text-white absolute bottom-[5%] left-1/2 -translate-x-1/2 text-[12px] w-[7rem] h-[2rem] rounded-3xl bg-ocean-blue flex items-center justify-center ">
          <span className="flex justify-center items-center ml-[-0.5rem]">
            {actionLabel}
          </span>
        </button>
      </div>
    );
  } else {
    return (
      <div className="flex rounded-xl mr-[0.5rem] items-center relative mt-[0.5rem] ">
        <a href="#" className="flex max-w-[38rem]">
          <img
            src={image}
            alt={name}
            className="max-w-[19.2rem] aspect-[16/9] h-[200px] w-[310px] rounded-l-xl"
          />
          <img
            src={image_double}
            alt={name}
            className="max-w-[19.2rem] aspect-[16/9] h-[200px] w-[310px] rounded-r-xl"
          />
        </a>
        {actionLabel === "Visit store" ? (
          <button className={buttonClasses}>
            <span className="flex justify-center items-center ml-[-0.5rem]">
              {actionLabel}
            </span>
            <div className="bg-white border border-black w-[1.2rem] h-[1.2rem] rounded-[50%] absolute flex items-center justify-center right-2">
              <RxArrowTopRight />
            </div>
          </button>
        ) : (
          <button className={buttonClasses}>
            <span className="flex justify-center items-center">
              {actionLabel}
            </span>
          </button>
        )}
      </div>
    );
  }
};

export default TopProducts;
