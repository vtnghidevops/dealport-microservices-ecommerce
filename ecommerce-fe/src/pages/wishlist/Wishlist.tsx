import React from 'react';
import WishlistTable from '../../components/wishlist/WishlistTable';
import EmptyWishlist from '../system/EmptyWishList'
import { useWishlist } from '@/context/WishlistContext';

const Wishlist: React.FC = () => {
  const { wishlistItems } = useWishlist();

  if (wishlistItems.length === 0) {
    return <EmptyWishlist />;
  }

  return (
    <div className="container mx-auto py-[1.5rem] px-[5rem]">
      <h3 className="text-[20px] font-bold mb-6">Wishlist</h3>
      <WishlistTable wishlistItems={wishlistItems} />
    </div>
  );
};

export default Wishlist;