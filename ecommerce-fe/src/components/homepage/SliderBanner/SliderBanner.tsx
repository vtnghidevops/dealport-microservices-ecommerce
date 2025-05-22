import React from "react";
import { useState, useEffect, useRef } from "react";
import { getButtonClass } from "@/utils/buttonUtils";
import { FaArrowLeft } from "react-icons/fa6";
import { FaArrowRight } from "react-icons/fa6";
import { BannerService } from "@/services/product/product.service";
import { SliderBannerItem } from "@/types/banner.model";
import { Link } from "react-router-dom";
import { SliderSkeleton } from "@/components/ui/skeletons";
import { handleCreateSlug } from "@/utils/helpers";
import { useNavigate } from "react-router-dom";

const SliderBanner: React.FC = () => {
  const [sliderData, setSliderData] = useState<SliderBannerItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [currentIndex, setCurrentIndex] = useState<number>(0);
  const [direction, setDirection] = useState<'next' | 'prev'>('next');
  const [isPaused, setIsPaused] = useState<boolean>(false);
  const pauseTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const navigate = useNavigate();

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

  // Function to handle pausing and resuming the slider
  const handlePauseAndResume = () => {
    setIsPaused(true); // Pause the slider immediately

    // Clear any existing timeout
    if (pauseTimeoutRef.current) {
      clearTimeout(pauseTimeoutRef.current);
    }

    // Set a new timeout to resume the slider after inactivity
    pauseTimeoutRef.current = setTimeout(() => {
      setIsPaused(false); // Resume auto-sliding
    }, 8000); // Resume after 8 seconds of inactivity
  };

  // Handle navigation
  const goToPrevious = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent banner click when clicking navigation
    const isFirstSlide = currentIndex === 0;
    const newIndex = isFirstSlide ? sliderData.length - 1 : currentIndex - 1;
    setDirection('prev');
    setCurrentIndex(newIndex);
    handlePauseAndResume(); // Pause auto-sliding when user interacts
  };

  const goToNext = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent banner click when clicking navigation
    const isLastSlide = currentIndex === sliderData.length - 1;
    const newIndex = isLastSlide ? 0 : currentIndex + 1;
    setDirection('next');
    setCurrentIndex(newIndex);
    handlePauseAndResume(); // Pause auto-sliding when user interacts
  };

  // Handle banner click
  const handleBannerClick = () => {
    handlePauseAndResume(); // Pause auto-sliding when user interacts

    if (sliderData[currentIndex]?.linkUrl) {
      navigate(sliderData[currentIndex].linkUrl);
    } else {
      navigate("/products");
    }
  };

  // Handle mouse enter/leave for slider
  const handleMouseEnter = () => {
    setIsPaused(true); // Pause on mouse enter
  };

  const handleMouseLeave = () => {
    setIsPaused(false); // Resume on mouse leave
  };

  // Auto-slide functionality
  useEffect(() => {
    if (sliderData.length <= 1 || isPaused) return; // Don't auto-slide if only one item or paused

    const slideInterval = setInterval(() => {
      setDirection('next');
      setCurrentIndex((prevIndex) => {
        const isLastSlide = prevIndex === sliderData.length - 1;
        return isLastSlide ? 0 : prevIndex + 1;
      });
    }, 5000); // Change slide every 5 seconds

    return () => clearInterval(slideInterval);
  }, [sliderData.length, isPaused]);

  // Clean up timeout on unmount
  useEffect(() => {
    return () => {
      if (pauseTimeoutRef.current) {
        clearTimeout(pauseTimeoutRef.current);
      }
    };
  }, []);

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
    return <SliderSkeleton />
  }

  // Function to get slide transition classes
  const getSlideClasses = (index: number) => {
    if (index === currentIndex) {
      return 'opacity-100 z-1 scale-100 transform-none';
    } else {
      const isPrevious = (index === currentIndex - 1) || (currentIndex === 0 && index === sliderData.length - 1);
      const isNext = (index === currentIndex + 1) || (currentIndex === sliderData.length - 1 && index === 0);

      if (direction === 'next' && isPrevious) {
        return 'opacity-0 z-0 -translate-x-1/4 scale-95';
      } else if (direction === 'prev' && isNext) {
        return 'opacity-0 z-0 translate-x-1/4 scale-95';
      } else {
        return 'opacity-0 z-0 scale-105';
      }
    }
  };

  return (
    <div className="w-full mb-8 relative z-0">
      {/* Categories Navigation */}
      <div className="w-full bg-white shadow-sm border-b border-gray-200 relative z-10">
        <div className="max-w-7xl mx-auto">
          <ul className="flex items-center gap-4 h-[48px] overflow-x-auto px-4 md:px-6">
            {categories.map((category, index) => (
              <li key={index} className="flex-shrink-0">
                <Link
                  to={`/category/${handleCreateSlug(category)}`}
                  className="whitespace-nowrap px-3 py-2 hover:text-blue-600 transition-colors text-cyprus"
                >
                  {category}
                </Link>
              </li>
            ))}
            <li className="flex-shrink-0 hover:text-blue-600 transition-colors ml-auto">
              <Link to="/products" className="px-3 py-2 text-blue-500">
                See more
              </Link>
            </li>
          </ul>
        </div>
      </div>

      {/* Slider container */}
      <div
        className="w-full h-[400px] sm:h-[450px] md:h-[500px] relative overflow-hidden border border-gray-200 cursor-pointer rounded-xl"
        onClick={handleBannerClick}
        onMouseEnter={handleMouseEnter}
        onMouseLeave={handleMouseLeave}
      >
        {sliderData.map((item, index) => (
          <div
            key={index}
            className={`absolute top-0 left-0 w-full h-full transition-all duration-700 ease-in-out ${getSlideClasses(index)
              }`}
          >
            {/* Full width image */}
            <img
              src={item.imageUrl}
              alt={item.title}
              className="w-full h-full object-cover object-center select-none"
              loading="eager"
              style={{
                objectFit: "cover",
                imageRendering: "crisp-edges",
                backfaceVisibility: "hidden",
                transform: "translateZ(0)",
              }}
            />
          </div>
        ))}

        {/* Navigation dots - moved higher up */}
        {sliderData.length > 1 && (
          <div
            className="absolute bottom-10 md:bottom-14 left-1/2 transform -translate-x-1/2 flex space-x-3 z-10 bg-black/20 px-4 py-2 rounded-full"
            onClick={(e) => e.stopPropagation()}
          >
            {sliderData.map((_, index) => (
              <button
                key={index}
                onClick={(e) => {
                  e.stopPropagation();
                  setDirection(index < currentIndex ? 'prev' : 'next');
                  setCurrentIndex(index);
                  handlePauseAndResume(); // Pause auto-sliding when clicking dots
                }}
                className={`w-3 h-3 rounded-full transition-all duration-300 ${currentIndex === index
                  ? 'bg-white w-5 h-3'
                  : 'bg-gray-400/70 hover:bg-gray-200'
                  }`}
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
              className="h-10 w-10 md:h-[3rem] md:w-[3rem] absolute left-2 md:left-4 top-1/2 -translate-y-1/2 bg-white/90 rounded-full flex items-center justify-center hover:bg-white transition-colors focus:outline-none focus:ring-2 focus:ring-blue-300 z-10 shadow-md"
              aria-label="Previous slide"
            >
              <FaArrowLeft className="w-4 h-4 md:w-[16px] md:h-[16px]" />
            </button>

            <button
              onClick={goToNext}
              className="h-10 w-10 md:h-[3rem] md:w-[3rem] absolute right-2 md:right-4 top-1/2 -translate-y-1/2 bg-white/90 rounded-full flex items-center justify-center hover:bg-white transition-colors focus:outline-none focus:ring-2 focus:ring-blue-300 z-10 shadow-md"
              aria-label="Next slide"
            >
              <FaArrowRight className="w-4 h-4 md:w-[16px] md:h-[16px]" />
            </button>
          </>
        )}
      </div>
    </div>
  );
}

export default SliderBanner;