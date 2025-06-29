import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useWishlist } from '@/hooks/useWishList';
import { Product } from '@/types/product.model';
import { formatImageUrl } from '@/utils/api-config';

interface ProductCardProps {
  product: Product;
  onClick?: (e: React.MouseEvent) => void;
}

const ProductCard: React.FC<ProductCardProps> = ({ product, onClick }) => {
  const { addToWishlist, removeFromWishlist, isInWishlist } = useWishlist();
  const navigate = useNavigate();
  const isLiked = isInWishlist(product.id.toString());

  // Calculate discount percentage if not provided directly
  const discountPercentage = product.discount ? product.discount :
    product.originalPrice ? Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100) : 0;

  const handleLikeClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

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
    <div className="h-full w-full p-3 flex flex-col group relative border border-gray-200 rounded-lg overflow-hidden shadow-sm hover:shadow-md hover:border-aqua-spring transition-all duration-300 hover:scale-[1.02]">
      {/* Special tags */}
      {(() => {
        if (product.stockQuantity === 0) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-error text-white text-xs font-medium px-2 py-1 rounded-md shadow-sm">
              OUT OF STOCK
            </div>
          );
        } else if (product.stockQuantity < 20) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-orange-500 text-white text-xs font-medium px-2 py-1 rounded-md shadow-sm">
              HOT
            </div>
          );
        } else if (discountPercentage > 0) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-red-500 text-white text-xs font-medium px-2 py-1 rounded-md shadow-sm">
              {discountPercentage}% OFF
            </div>
          );
        } else if (product.price < 30) {
          return (
            <div className="z-20 absolute top-2 left-2 bg-success text-white text-xs font-medium px-2 py-1 rounded-md shadow-sm">
              SALE
            </div>
          );
        }
        return null;
      })()}

      {/* Wishlist button */}
      <div className="z-30 absolute top-2 right-2 flex flex-col space-y-2">
        <button
          onClick={handleLikeClick}
          className={`
            bg-white p-2 rounded-full shadow-sm hover:bg-gray-100 
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
      <div className="w-full flex-none h-[180px] flex items-center justify-center overflow-hidden">
        <div onClick={handleProductClick} className="cursor-pointer w-full h-full">
          <div className="p-2 flex justify-center items-center relative h-full w-full overflow-hidden">
            <img
              src={formatImageUrl(product.imageUrl)}
              alt={product.name}
              className="max-h-[150px] max-w-[150px] object-contain transition-transform duration-500 hover:scale-110"
            />
          </div>
        </div>
      </div>

      {/* Product info */}
      <div className="p-2 flex-grow flex flex-col">
        <div className="cursor-pointer mb-2" onClick={handleProductClick}>
          <h3 className="text-[16px] font-bold text-cyprus hover:text-ocean-green transition-colors duration-300 line-clamp-2 h-[48px]">
            {product.name}
          </h3>
        </div>

        {/* Price section with original price and discount */}
        <div className="mt-auto mb-2 flex items-center">
          <span className="text-lg font-bold text-ocean-green">
            ${product.price.toFixed(2)}
          </span>

          {product.originalPrice && product.originalPrice > 0 &&
            product.originalPrice > product.price && (
              <span className="ml-2 text-sm text-gray-500 line-through">
                ${product.originalPrice.toFixed(2)}
              </span>
            )}
        </div>

        {/* Rating */}
        <div className="flex items-center">
          <div className="flex">
            {[...Array(5)].map((_, i) => (
              <span
                key={i}
                className={`text-[16px] ${i < product.reviewsAvg.rating
                  ? "text-[#FF9017]"
                  : "text-gray-300"
                  }`}
              >
                ★
              </span>
            ))}
          </div>
          <span className="ml-1 mr-2 text-sm text-[#FF9017]">
            {product.reviewsAvg?.rating?.toFixed(1)}
          </span>
          <div className="mr-1 w-1 h-1 rounded-full bg-gray-400"></div>
          <span className="ml-1 text-sm text-gray-500">
            {product.orders} orders
          </span>
        </div>

        {/* Add view details button for consistency */}
        <div className="mt-3 flex justify-end">
          <button
            onClick={handleProductClick}
            className="bg-aqua-spring text-cyprus font-medium text-xs py-1.5 px-3 rounded-md hover:shadow-md hover:bg-green-500 hover:text-white transition-all duration-300"
          >
            View Details
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;