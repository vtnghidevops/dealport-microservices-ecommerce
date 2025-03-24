import React, { useState } from "react";
import { FaRegHeart, FaHeart } from "react-icons/fa";

const StarRating = ({ rating, maxRating = 5 }) => {
  return (
    <div className="flex">
      {[...Array(maxRating)].map((_, index) => {
        const starValue = index + 1;
        return (
          <span
            key={index}
            className={`text-lg ${
              starValue <= rating
                ? "text-yellow-400" // Full star
                : starValue <= rating + 0.5
                ? "text-yellow-400" // Half star
                : "text-gray-300" // Empty star
            }`}
          >
            {starValue <= rating ? "★" : starValue <= rating + 0.5 ? "★" : "☆"}
          </span>
        );
      })}
    </div>
  );
};
export default function ProductCard({
  title,
  description,
  price,
  originalPrice,
  discount,
  reviews,
  imageUrl,
}) {
  const [isFavorite, setIsFavorite] = useState(false);
  // Assuming reviews contains a rating property
  const rating = reviews?.rating || 0;
  const reviewCount = reviews?.count || 0;
  
  const toggleFavorite = () => {
    setIsFavorite(!isFavorite);
  };
  
  return (
    <div className="mx-2 p-[10px] rounded-xl w-[272px] border border-grep-300 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.1),0_4px_6px_-2px_rgba(0,0,0,0.05)] ">
      <div className="relative overflow-hidden rounded-xl">
        <a href="#" className="block overflow-hidden">
          <img
            src={imageUrl}
            alt={title}
            className="w-[248px] h-[180px] object-cover rounded-xl transition-transform duration-700"
          />
        </a>
        <button 
          onClick={toggleFavorite}
          className="absolute top-2 right-[5%] flex items-center justify-center rounded-full h-[1.5rem] w-[1.5rem] bg-white hover:bg-gray-100 transition-colors duration-300"
        >
          {isFavorite ? (
            <FaHeart className="text-red-500" />
          ) : (
            <FaRegHeart className="text-gray-500" />
          )}
        </button>
      </div>

      <div className="mt-4">
        <h3 className="title font-bold">{title}</h3>
        <p className="text-sm mt-1 line-clamp-2">{description}</p>

        <div className="flex items-center mt-2">
          <StarRating rating={rating} />
          <span className="text-xs ml-2">({reviewCount} reviews)</span>
        </div>

        <div className="flex items-center mt-2">
          <span className="text-lg font-bold text-ocean-green">${price}</span>
          { 
            (originalPrice !== 0) && (originalPrice !== "") && (
              <span className="text-sm line-through text-grey ml-2">
                (${originalPrice})
              </span>
            )}
          {(discount !== 0) && (discount !== "")  && (
            <span className="text-sm ml-2 text-ocean-green">
              {discount}% Off
            </span>
          )}
        </div>

        <div className="flex justify-between items-center mt-[15px]">
          <a href="#" className="text-primary text-sm hover:text-ocean-green transition-all duration-300">
            View Details
          </a>
          <button className="btn-primary text-sm py-1 h-[39px] w-[120px] border-2 border-transparent hover:bg-green-500 hover:shadow-md transition-all duration-300 ease-out">
            Add to Cart
          </button>
        </div>
      </div>
    </div>
  );
}