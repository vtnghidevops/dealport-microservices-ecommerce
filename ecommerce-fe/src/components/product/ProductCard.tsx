import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useWishlist } from '@/hooks/useWishList';
import { Product } from '@/types/product.model';

interface ProductCardProps {
  product: Product;
  onClick?: (e: React.MouseEvent) => void;
}

const ProductCard: React.FC<ProductCardProps> = ({ product, onClick }) => {
  const { addToWishlist, removeFromWishlist, isInWishlist } = useWishlist();
  const navigate = useNavigate();
  const isLiked = isInWishlist(product.id.toString()); // Check if the product is in the wishlist
  // const [liked, setLiked] = useState<boolean>(false); // State to manage like button
  // console.log("product in product card", product.name)
  // Calculate discount percentage if not provided directly
  const discountPercentage = product.discount ? product.discount : 
    product.originalPrice ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100) : 0;
  
  const handleLikeClick = (e: React.MouseEvent) => {
    e.preventDefault(); // Prevent event bubbling
    e.stopPropagation(); // Stop the event from propagating to parent elements
    // setLiked(!liked);
    if (isLiked) {
      removeFromWishlist(product.id.toString());
      
    } else {
      addToWishlist(product);
     
    }
  };

  const handleProductClick = (e: React.MouseEvent) => {
    e.preventDefault();
    
    // If parent provided onClick, use it
    if (onClick) {
      onClick(e);
    } else {
      // Otherwise handle navigation ourselves
      const productUrl = `/category/${product.categorySlug}/${product.slug}`;
      navigate(productUrl);
      
      // Scroll to top after navigation
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      });
    }
  };

  return (
    <div className="drop-shadow-sm filter p-3 min-h-[315px] min-w-[225px] group relative border rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow">
      {/* Special tags */}
      {(() => {
        if (product.stockQuantity === 0) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-error text-white text-xs font-medium px-2 py-1 rounded">
              OUT OF STOCK
            </div>
          );
        } else if (product.stockQuantity < 20) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-orange-500 text-white text-xs font-medium px-2 py-1 rounded">
              HOT
            </div>
          );
        } else if (product.price < 30) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-success text-white text-xs font-medium px-2 py-1 rounded">
              SALE
            </div>
          );
        }
        return null;
      })()}

      {/* Wishlist & Quick view buttons */}
      <div className="z-30 absolute top-2 right-2 flex flex-col space-y-2">
        <button
          onClick={handleLikeClick}
          className={`
            bg-white p-2 rounded-full shadow hover:bg-gray-100 
            transition-all duration-300 transform
            ${isLiked ? "scale-110" : "scale-100"}
          `}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill={isLiked ? "#FF69B4" : "none"}
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke={isLiked ? "#FF69B4" : "currentColor"}
            className="w-5 h-5 transition-colors"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z"
            />
          </svg>
        </button>
      </div>

      {/* Product image */}
      <div className="w-full flex items-center justify-center">
        <div onClick={handleProductClick} className="cursor-pointer">
          <div className="p-2 flex justify-center items-center relative h-[180px] w-[190px] overflow-hidden">
            <img
              src={product.imageUrl}
              alt={product.name}
              className="max-w-[150px] h-full object-contain transition-transform hover:scale-105"
            />
          </div>
        </div>
      </div>

      {/* Product info */}
      <div className="p-4 mt-3">
        <div className="cursor-pointer" onClick={handleProductClick}>
          <h3 className="text-[16px] font-bold text-cyprus hover:text-blue-600 transition-colors line-clamp-2">
            {product.name}
          </h3>
        </div>
        {/* Price section with original price and discount */}
        <div className="mt-2 flex items-center">
          <span className="text-lg font-semibold text-gray-900">
            ${product.price.toFixed(2)}
          </span>

          {product.originalPrice && product.originalPrice > 0 &&
            product.originalPrice > product.price && (
              <>
                <span className="ml-2 text-sm text-gray-500 line-through">
                  ${product.originalPrice.toFixed(2)}
                </span>
                <span className="ml-2 text-xs font-medium text-green-600">
                  {discountPercentage}% off
                </span>
              </>
            )}
        </div>

        {/* Rating */}
          <div className="mt-2 flex items-center">
            <div className="flex">
              {[...Array(5)].map((_, i) => (
                <span
                  key={i}
                  className={`text-[16px] ${
                    i < product.reviewsAvg.rating
                      ? "text-[#FF9017]"  
                      : "text-gray-300"
                  }`}
                >
                  ★
                </span>
              ))}
            </div>
            <span className="ml-1 mr-2 text-sm text-[#FF9017]">
              ({product.reviewsAvg.rating})
            </span>
            <div className="mr-1 w-1 h-1 rounded-full border border-neutral-400 bg-neutral-400"></div>
            <span className="ml-1 text-sm text-neutral-500">
              {product.orders} orders
            </span>
          </div>
      </div>
    </div>
  );
};

export default ProductCard;