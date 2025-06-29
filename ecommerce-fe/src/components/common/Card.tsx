import { FaRegHeart, FaHeart } from "react-icons/fa";
import { useCart } from '@/hooks/useCart';
import { useWishlist } from "@/hooks/useWishList";
import { Product } from "@/types/product.model";
import { Link } from "react-router-dom";

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
  const rating = product.reviewsAvg?.rating || 0;
  const reviewCount = product.reviewsAvg?.count || 0;
  const { addToCart } = useCart();
  const { addToWishlist, removeFromWishlist, isInWishlist } = useWishlist();
  const isLiked = isInWishlist(product.id.toString());

  const handleLikeClick = (e: React.MouseEvent) => {
    e.preventDefault(); // Prevent event bubbling
    e.stopPropagation(); // Stop the event from propagating to parent elements
    if (isLiked) {
      removeFromWishlist(product.id.toString());
    } else {
      addToWishlist(product);
    }
  };

  const handleAddToCart = () => {
    addToCart({
      productId: Number(product.id),
      name: product.name,
      price: product.price,
      originalPrice: product.originalPrice || undefined,
      imageUrl: product.imageUrl
    }, 1);
  };

  return (
    <div className="mx-2 p-[10px] rounded-xl max-w-[272px] min-w-[272px] border border-gray-200 shadow-sm hover:shadow-md hover:border-aqua-spring transition-all duration-300 hover:scale-[1.02]">
      <div className="relative overflow-hidden rounded-xl">
        <Link to={`/category/${product.categorySlug || 'uncategorized'}/${product.slug}`} className="block overflow-hidden">
          <img
            src={product.imageUrl}
            alt={product.name}
            className="min-w-[248px] min-h-[180px] max-w-[248px] max-h-[180px] object-contain rounded-xl transition-transform duration-500 hover:scale-105"
          />
        </Link>
        <button
          onClick={handleLikeClick}
          className="absolute top-2 right-[5%] flex items-center justify-center rounded-full h-[1.5rem] w-[1.5rem] bg-white hover:bg-gray-100 transition-colors duration-300 shadow-sm"
        >
          {isLiked ? (
            <FaHeart className="text-red-500" />
          ) : (
            <FaRegHeart className="text-gray-500" />
          )}
        </button>

        {/* Price tag */}
        {product.originalPrice !== 0 && product.originalPrice > product.price && (
          <span className="absolute top-2 left-2 px-2 py-1 bg-red-500 text-white font-bold text-xs rounded-md">
            {Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)}% OFF
          </span>
        )}
      </div>

      <div className="mt-4">
        <h3 className="title font-bold text-cyprus line-clamp-1 hover:text-ocean-green transition-colors duration-300">{product.name}</h3>
        <p className="text-sm mt-1 line-clamp-2 text-gray-600">{product.description}</p>

        <div className="flex items-center mt-2">
          <StarRating rating={rating} />
          <span className="text-xs ml-2 text-gray-500">({reviewCount || 0})</span>
        </div>

        <div className="flex items-center mt-2">
          <span className="text-lg font-bold text-ocean-green">${product.price}</span>
          {product.originalPrice !== 0 && (
            <span className="text-sm line-through text-gray-500 ml-2">
              ${product.originalPrice}
            </span>
          )}
        </div>

        <div className="flex justify-between items-center mt-[15px]">
          <Link to={`/category/${product.categorySlug || 'uncategorized'}/${product.slug}`} className="text-primary text-sm hover:text-ocean-green transition-all duration-300 font-medium">
            View Details
          </Link>
          <button
            onClick={handleAddToCart}
            className="bg-aqua-spring text-cyprus font-medium text-sm py-1.5 px-3 rounded-md hover:shadow-md hover:bg-green-500 hover:text-white transition-all duration-300 ease-out"
          >
            Add to Cart
          </button>
        </div>
      </div>
    </div>
  );
}

export default ProductCard;