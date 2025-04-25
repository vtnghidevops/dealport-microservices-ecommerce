import React from "react";
import { useState, useEffect } from "react";
import { getButtonClass } from "@/utils/buttonUtils";
import { FaArrowLeft } from "react-icons/fa6";
import { FaArrowRight } from "react-icons/fa6";
import { BannerService } from "@/services/product.service";
import { SliderBannerItem } from "@/types/banner.model";
import { Link } from "react-router-dom";
import Loading from "@/components/shared/Loading";

const SliderBanner: React.FC = () => {

  const [sliderData, setSliderData] = useState<SliderBannerItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [currentIndex, setCurrentIndex] = useState<number>(0);

  // fetch banner data
  useEffect(() => {
    const fetchBanners = async () => {
      try {
        const data = await BannerService.getSliderBanners();
        setSliderData(data);
        setLoading(false);
      } catch (error) {
        console.error("Error fetching banners:", error);
        setLoading(false);
      }
    };

    fetchBanners();
  }, []);

  // Handle navigation
  const goToPrevious = () => {
    const isFirstSlide = currentIndex === 0;
    const newIndex = isFirstSlide ? sliderData.length - 1 : currentIndex - 1;
    setCurrentIndex(newIndex);
  };

  const goToNext = () => {
    const isLastSlide = currentIndex === sliderData.length - 1;
    const newIndex = isLastSlide ? 0 : currentIndex + 1;
    setCurrentIndex(newIndex);
  };

  // Auto-slide functionality
  useEffect(() => {
    if (sliderData.length <= 1) return; // Don't auto-slide if only one item

    const slideInterval = setInterval(() => {
      goToNext();
    }, 5000); // Change slide every 5 seconds

    return () => clearInterval(slideInterval);
  }, [currentIndex, sliderData.length]);

  const categories = [
    "Men",
    "Woman",
    "Baby",
    "Grocery & Essentials",
    "Streetwear",
    "Shoes",
    "Beauty",
    "Electronics",
    "Industrial equipment",
  ];

  if (loading) {
    return <Loading />
  }

  if (sliderData.length === 0) {
    return <div className="w-full h-[500px] flex items-center justify-center">No banner data available</div>;
  }

  // const currentSlide = sliderData[currentIndex];

  return (
    <div className="relative w-full text-white h-[500px] m-h-[500px]">
      <nav className="w-full p-2 z-10 relative">
        <ul className="mr-0 ml-[67px] w-full max-w-7xl mx-auto text-black flex items-center justify-between h-[48px] overflow-x-auto">
          {categories.map((category, index) => (
            <li key={index}>
              <a
                href="#"
                className="whitespace-nowrap px-3 py-2 hover:text-blue-600 transition-colors text-cyprus"
              >
                {category}
              </a>
            </li>
          ))}
          <li className="hover:text-blue-600 transition-colors">
            <a href="#" className="px-3 py-2 text-blue-500">
              See more
            </a>
          </li>
        </ul>
      </nav>

      {/* Slider container with transition effect */}
      <div className="w-full bg-cyprus h-full relative overflow-hidden">
        <div
          className="w-full h-full transition-transform duration-500 ease-in-out flex"
          style={{
            transform: `translateX(-${currentIndex * 100}%)`,
            width: `${sliderData.length * 100}%`
          }}
        >
           {/* !bg-cyprus" style={{
              backgroundColor: item.background_color || '#1e3a8a'
            }} */}
          {sliderData.map((item, index) => (
            <div key={index} className="w-full h-full flex-shrink-0 relative">
              <img
                src={item.imageUrl}
                alt={item.title}
                className="absolute right-[75%] h-full object-contain max-w-[951px] "
              />
              <div className="absolute left-10 -top-5 w-[40%] h-full flex flex-col justify-center pl-12">
                <div className="ml-[6rem] max-w-full z-10">
                  <h2 className="text-4xl font-bold mb-2 !text-white" style={{ color: item.textColor || '#ffffff' }}>
                    {item.title}
                  </h2>
                  <p className="text-5xl font-bold mb-6 italic !text-white" style={{ color: item.textColor || '#ffffff' }}>
                    {item.discount}
                  </p>
                </div>

                <div className="mt-5 z-10 ml-[6rem]">
                  <Link to={item.linkUrl || "#"}>
                    <button className={`${getButtonClass("secondary")} px-8 py-3`}>
                      {item.actionText || "Shop now"}
                    </button>
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Navigation dots */}
      {sliderData.length > 1 && (
        <div className="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex space-x-2">
          {sliderData.map((_, index) => (
            <button
              key={index}
              onClick={() => setCurrentIndex(index)}
              className={`w-3 h-3 rounded-full ${currentIndex === index ? 'bg-white' : 'bg-gray-400'}`}
              aria-label={`Go to slide ${index + 1}`}
            />
          ))}
        </div>
      )}

      {/* Navigation buttons - only show if more than one slide */}
      {sliderData.length > 1 && (
        <>
          <button
            onClick={goToPrevious}
            className="h-[3.5rem] text-black w-[3.5rem] absolute left-[50px] top-1/2 -translate-y-1/2 bg-white rounded-full p-4 flex items-center justify-center hover:bg-gray-100 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-300 z-10"
            aria-label="Previous slide"
          >
            <FaArrowLeft className="w-[20px] h-[20px]" />
          </button>

          <button
            onClick={goToNext}
            className="absolute text-black h-[3.5rem] w-[3.5rem] right-[50px] top-1/2 -translate-y-1/2 bg-white rounded-full p-4 flex items-center justify-center hover:bg-gray-100 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-300 z-10"
            aria-label="Next slide"
          >
            <FaArrowRight className="w-[20px] h-[20px]" />
          </button>
        </>
      )}
    </div>
  );
}

export default SliderBanner;