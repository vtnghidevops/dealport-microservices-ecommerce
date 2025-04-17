import React, { useRef, useState, useEffect } from "react";
import ProductCard from "../../common/Card";
import { MdChevronRight, MdChevronLeft } from "react-icons/md";
import { Product } from "@/types/product.model";
import { LimitedDealService } from "./services/litmitedDeal.service";
import { handleViewAll } from "@/utils/helpers";
import { useNavigate } from "react-router-dom";
const LimitedDeal: React.FC = () => {
  // fetch data
  const [limitedData, setLimitedData] = useState<Product[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const navigate = useNavigate();
  useEffect(() => {
    const fetchLimitedData = async () => {
      try {
        const data = await LimitedDealService.getLitmitedData();
        setLimitedData(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching testimonials:", error);
        setLoading(false);
      }
    };

    fetchLimitedData();
  }, []);

  // event slider
  const sliderRef = useRef<HTMLDivElement | null>(null);
  const [showLeftArrow, setShowLeftArrow] = useState(false);
  const [showRightArrow, setShowRightArrow] = useState(true);

  const checkArrowVisibility = () => {
    if (!sliderRef.current) return;

    // Check if we're at the beginning of the slider
    setShowLeftArrow(sliderRef.current.scrollLeft > 0);

    // Check if we're at the end of the slider
    const isAtEnd =
      Math.ceil(sliderRef.current.scrollLeft + sliderRef.current.clientWidth) >=
      sliderRef.current.scrollWidth;

    setShowRightArrow(!isAtEnd);
  };

  // Check arrow visibility on initial load and when products change
  useEffect(() => {
    checkArrowVisibility();
  }, [limitedData]);

  const scroll = (direction: "left" | "right") => {
    if (sliderRef.current) {
      const scrollAmount =
        direction === "left"
          ? -sliderRef.current.clientWidth / 2
          : sliderRef.current.clientWidth / 2;

      sliderRef.current.scrollBy({
        left: scrollAmount,
        behavior: "smooth",
      });

      // Set a timeout to check arrow visibility after scrolling completes
      setTimeout(checkArrowVisibility, 500);
    }
  };

  return (
    <div className="h-full relative">
      <span className="header-2 font-bold p-[10px] !px-0 block mb-3">
        Limited-Time Deal
      </span>

      <div className="relative">
        {/* Left gradient and button */}
        {showLeftArrow && (
          <>
            <div className="-ml-[1rem] absolute top-0 left-0 h-full w-[5rem] bg-gradient-to-r from-white to-transparent pointer-events-none z-10"></div>
            <button
              onClick={() => scroll("left")}
              className="flex justify-center items-center h-10 w-10 absolute left-3 top-1/2 -translate-y-1/2 z-20 bg-white bg-opacity-75 rounded-full p-2 shadow-md transition-opacity hover:bg-opacity-100"
              aria-label="Scroll left"
            >
              <MdChevronLeft className="h-5 w-5 text-gray-500" />
            </button>
          </>
        )}

        {/* Product slider */}
        <div
          ref={sliderRef}
          className="-ml-[1rem] h-full w-full flex overflow-x-auto scroll-smooth py-4 px-2 "
          style={{
            scrollbarWidth: "none",
            msOverflowStyle: "none",
          }}
          onScroll={checkArrowVisibility}
        >
          {limitedData.map((product, index) => (
            <div className="min-w-[280px]">
              <ProductCard
                key={index} 
                product={product}
              />
            </div>
          ))}
        </div>

        {/* Right gradient and button */}
        {showRightArrow && (
          <>
            <div className="absolute top-0 right-0 h-full w-[5rem] bg-gradient-to-l from-white to-transparent pointer-events-none z-10"></div>
            <button
              onClick={() => scroll("right")}
              className="flex justify-center items-center h-10 w-10 absolute right-3 top-1/2 -translate-y-1/2 z-20 bg-white bg-opacity-75 rounded-full p-2 shadow-md transition-opacity hover:bg-opacity-100"
              aria-label="Scroll right"
            >
              <MdChevronRight className="h-5 w-5 text-gray-500" />
            </button>
          </>
        )}
      </div>
      <div className="absolute right-[1%] top-[2%]">
        <button onClick={() => handleViewAll(navigate)} className="w-[8rem] h-[3rem] rounded-3xl border border-black bg-white hover:bg-black hover:text-white transition-all duration-300 ease-in-out">
          View All
        </button>
      </div>
    </div>
    
  );
};

export default LimitedDeal;
