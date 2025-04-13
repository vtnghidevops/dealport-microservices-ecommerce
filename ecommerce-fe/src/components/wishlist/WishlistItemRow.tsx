import React from 'react';
import { FiTrash2 } from 'react-icons/fi';
import { Product } from '../product/models/product.model';
import { useWishlist } from '@/context/WishlistContext';
import { useNavigate } from 'react-router-dom';
// import { useCart } from '@/context/CartContext';
interface WishlistItemRowProps {
  item: Product;
}

const WishlistItemRow: React.FC<WishlistItemRowProps> = ({ item }) => {
  const { removeFromWishlist } = useWishlist();
  // const { addToCart } = useCart();
  const navigate = useNavigate();
  // const handleAddToCart = async () => {
  //   try {
  //     await addToCart(item);
  //     enqueueSnackbar('Added to cart successfully!', { variant: 'success' });
  //     // Optionally remove from wishlist after adding to cart
  //     // removeFromWishlist(item.id);
  //   } catch (error) {
  //     enqueueSnackbar('Failed to add to cart', { variant: 'error' });
  //   }
  // };

  const formatPrice = (price: number) => {
    return `$${price.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`;
  };

  const handleProductClick = () => {
    navigate(`/${item.categorySlug}/${item.slug}`);
  };

  return (
    <tr className="h-[60px] hover:bg-neutral-50">
      <td className="px-5 py-4 w-1/2">
        <div className="flex items-center cursor-pointer" onClick={handleProductClick}>
          <div className="h-[30px] w-[30px] flex-shrink-0">
            <img
              className="h-full w-full object-contain"
              src={item.image_url}
              alt={item.name}
            />
          </div>
          <div className="ml-5">
            <div className="max-w-[400px] text-sm font-medium text-neutral-700">
              {item.name}
            </div>
          </div>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap w-1/6">
        <div className="flex items-center">
          {item.originalPrice && (
            <span className="mr-2 text-xs text-gray-400 line-through">
              {formatPrice(item.originalPrice)}
            </span>
          )}
          <span className="text-sm font-bold text-neutral-600">
            {formatPrice(item.price)}
          </span>
        </div>
      </td>
      <td className="px-6 py-4 whitespace-nowrap w-[14%]">
        <span
          className={`min-w-5 px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${
            item.stock
              ? "bg-green-100 text-success"
              : "bg-red-100 text-error"
          }`}
        >
          {item.stock ? "IN STOCK" : "OUT OF STOCK"}
        </span>
      </td>
      <td className="px-6 py-4 whitespace-nowrap">
        <div className="flex items-center ">
          <button
            // onClick={handleAddToCart}
            disabled={!item.stock}
            className={`h-[40px] w-[150px] text-[14px] flex justify-center items-center font-medium px-2 rounded-full  text-white py-3 mr-3 transition-all duration-300 transform text-sm ${
              item.stock
                ? "bg-[#0496FF] hover:bg-blue-500 text-white"
                : "bg-gray-300 text-gray-500 cursor-not-allowed"
            }`}
          >
            ADD TO CART
          </button>
          <button
            onClick={() => removeFromWishlist(item.id)}
            className="text-gray-400 hover:text-red-500"
          >
            <FiTrash2 size={18} />
          </button>
        </div>
      </td>
    </tr>
  );
};

export default WishlistItemRow;
