import React from "react";
import { FaRegHeart } from "react-icons/fa";
// import { FaHeart } from "react-icons/fa";

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
  // Assuming reviews contains a rating property
  const rating = reviews?.rating || 0;
  const reviewCount = reviews?.count || 0;
  return (
    <div className="mx-2 p-[10px] rounded-xl w-[272px] border border-grep-300 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.1),0_4px_6px_-2px_rgba(0,0,0,0.05)]">
      <div className="relative">
        <a href="#">
          <img
            src={imageUrl}
            alt={title}
            className="w-[248px] h-[180px] object-cover rounded-xl"
          />
        </a>
        <button className="absolute top-2 right-[5%] flex items-center justify-center rounded-full h-[1.5rem] w-[1.5rem] bg-white">
          {/* Heart icon could be replaced with an actual icon component */}
          <FaRegHeart></FaRegHeart>
        </button>
      </div>

      <div className="mt-4">
        <h3 className="title font-bold">{title}</h3>
        <p className="text-sm mt-1 line-clamp-2">{description}</p>

        <div className="flex items-center mt-2">
          {/* Star ratings - you might want to replace with an actual rating component */}
          <StarRating rating={rating} />
          <span className="text-xs ml-2">({reviewCount} reviews)</span>
        </div>

        <div className="flex items-center mt-2">
          <span className="text-lg font-bold text-ocean-green">${price}</span>
          {originalPrice !== undefined &&
            originalPrice !== null &&
            originalPrice !== 0 && (
              <span className="text-sm line-through text-grey ml-2">
                (${originalPrice})
              </span>
            )}
          {discount !== undefined && discount !== null && discount !== 0 && (
            <span className="text-sm ml-2 text-ocean-green">
              {discount}% Off
            </span>
          )}
        </div>

        <div className="flex justify-between items-center mt-[15px]">
          <a href="#" className="text-primary text-sm">
            View Details
          </a>
          <button className="btn-primary text-sm py-1 h-[39px] w-[120px]">
            Add to cart
          </button>
        </div>
      </div>
    </div>
  );
}
