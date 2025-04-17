import { useState } from "react";
import { FaRegHeart, FaHeart } from "react-icons/fa";
import { useCart } from '@/hooks/useCart';
import { useWishlist } from "@/hooks/useWishList";
import { Product } from "@/types/product.model";

interface StarRatingProps {
  rating: number;
  maxRating?: number;
}


interface ProductCardProps {
  product: Product;
}

const StarRating: React.FC<StarRatingProps> = ({ rating, maxRating = 5 }) => {
  return (
    <div className="flex">
      {[...Array(maxRating)].map((_, index) => {
        const starValue = index + 1;
        return (
          <span
            key={index}
            className={`text-lg ${starValue <= rating
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

const ProductCard: React.FC<ProductCardProps> = ({
  product
}) => {
  // type for cart item
  const productForCart = {
    id: product.id,
    name: product.name,
    price: product.price,
    originalPrice: product.originalPrice || undefined,
    image: product.image_url
  };
  //console.log(product);                    

  const rating = product.reviews?.rating || 0;
  const reviewCount = product.reviews?.count || 0;
  const { addToCart } = useCart();
  const { addToWishlist, removeFromWishlist, isInWishlist } = useWishlist();
  const isLiked = isInWishlist(product.id);

  const handleLikeClick = (e: React.MouseEvent) => {
    e.preventDefault(); // Prevent event bubbling
    e.stopPropagation(); // Stop the event from propagating to parent elements
    // setLiked(!liked);
    if (isLiked) {
      removeFromWishlist(product.id);

    } else {
      addToWishlist(product);

    }
  };


  return (
    <div className="mx-2 p-[10px] rounded-xl w-[272px] border border-grep-300 shadow-[0_2px_15px_-3px_rgba(0,0,0,0.1),0_4px_6px_-2px_rgba(0,0,0,0.05)] ">
      <div className="relative overflow-hidden rounded-xl">
        <a href={`/${product.categorySlug}/${product.slug}`} className="block overflow-hidden">
          <img
            src={product.image_url}
            alt={product.name}
            className="w-[248px] h-[180px] object-cover rounded-xl transition-transform duration-700"
          />
        </a>
        <button
          onClick={handleLikeClick}
          className="absolute top-2 right-[5%] flex items-center justify-center rounded-full h-[1.5rem] w-[1.5rem] bg-white hover:bg-gray-100 transition-colors duration-300"
        >
          {isLiked ? (
            <FaHeart className="text-red-500" />
          ) : (
            <FaRegHeart className="text-gray-500" />
          )}
        </button>
      </div>

      <div className="mt-4">
        <h3 className="title font-bold">{product.name}</h3>
        <p className="text-sm mt-1 line-clamp-2">{product.description}</p>

        <div className="flex items-center mt-2">
          <StarRating rating={rating} />
          <span className="text-xs ml-2">({reviewCount} reviews)</span>
        </div>

        <div className="flex items-center mt-2">
          <span className="text-lg font-bold text-ocean-green">${product.price}</span>
          {product.originalPrice !== 0 && (
            <span className="text-sm line-through text-grey ml-2">
              (${product.originalPrice})
            </span>
          )}
          {product.discount !== 0 && (
            <span className="text-sm ml-2 text-ocean-green">
              {product.discount}% Off
            </span>
          )}
        </div>

        <div className="flex justify-between items-center mt-[15px]">
          <a href={`/${product.categorySlug}/${product.slug}`} className="text-primary text-sm hover:text-ocean-green transition-all duration-300">
            View Details
          </a>
          {/* onClick={() => addToCart(productForCart)} */}
          <button onClick={() => addToCart(productForCart)} className="btn-primary text-sm py-1 h-[39px] w-[120px] border-2 border-transparent hover:bg-green-500 hover:shadow-md transition-all duration-300 ease-out">
            Add to Cart
          </button>
        </div>
      </div>
    </div>
  );
}

export default ProductCard;