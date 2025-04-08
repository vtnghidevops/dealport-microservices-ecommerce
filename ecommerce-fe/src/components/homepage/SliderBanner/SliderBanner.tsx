import {useState, useEffect} from "react";
import { getButtonClass } from "@/utils/buttonUtils";
import { FaArrowLeft } from "react-icons/fa6";
import { FaArrowRight } from "react-icons/fa6";
import { SliderBannerItem } from "./models/sliderBanner.model";
import { SliderBannerService } from "./services/sliderBanner.service";

export default function HeroBanner() {
  
    const [sliderData, setSliderData] = useState<SliderBannerItem[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    
      // fetch data trending Products
      useEffect(() => {
        const fetchTestimonials = async () => {
          try {
            const data = await SliderBannerService.getDataSlider();
            setSliderData(data);
            setLoading(false);
          } catch (error) {
            console.error("Error fetching testimonials:", error);
            setLoading(false);
          }
        };
    
        fetchTestimonials();
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

  return (
    <div className="relative w-full  text-white h-[500px] m-h-[500px]">
      <nav className="w-full p-2 ">
        <ul className="mr-0 ml-[67px] w-full max-w-7xl mx-auto text-black flex items-center justify-between h-[48px]">
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
      {sliderData.map((item, index) => {
        return (
          <div key={index} className="w-full bg-cyprus h-full relative">
            <img
              src={item.image_url}
              alt="Clothing rack"
              className="absolute right-0 h-full w-[70%]"
            />
            <div className="absolute inset-0 flex flex-col justify-center px-16 mb-[5rem]">
              <div>
                <h2 className="header-2 ml-[10.5rem]">{item.title}</h2>
                <p className="header-1 mb-6 ml-[10rem] italic">
                  {item.discount}
                </p>
              </div>

              <div className="ml-[10rem] mt-[1.5rem]">
                <div className="w-[12rem] h-[4rem] rounded-3xl ">
                  <button className={`${getButtonClass("secondary")}`}>
                    Show now
                  </button>
                </div>
              </div>
            </div>
          </div>
        );
      })}

      <button className="h-[3.5rem] text-black w-[3.5rem] absolute left-[50px] top-1/2 -translate-y-1/2 bg-white rounded-full p-4">
        <FaArrowLeft className="w-[49px] h-[25px]"></FaArrowLeft>
      </button>
      <button className="absolute text-black h-[3.5rem] w-[3.5rem] right-[50px]  top-1/2 -translate-y-1/2 bg-white rounded-full p-4">
        <FaArrowRight className="w-[49px] h-[25px]"></FaArrowRight>
      </button>
    </div>
  );
}
