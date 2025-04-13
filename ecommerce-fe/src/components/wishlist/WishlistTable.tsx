import React from 'react';
import { Product } from '../product/models/product.model';
import WishlistItemRow from './WishlistItemRow';

interface WishlistTableProps {
  wishlistItems: Product[];
}

const WishlistTable: React.FC<WishlistTableProps> = ({ wishlistItems }) => {
  return (
    <div className="bg-white rounded-lg shadow overflow-hidden">
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-neutral-100">
          <tr>
            <th scope="col" className="px-7 py-3 text-left text-xs font-medium text-neutral-600 uppercase tracking-wider">
              Products
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-neutral-600  uppercase tracking-wider">
              Price
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-neutral-600  uppercase tracking-wider">
              Stock Status
            </th>
            <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-neutral-600  uppercase tracking-wider">
              Actions
            </th>
            <th scope="col" className="relative px-6 py-3">
              <span className="sr-only">Remove</span>
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {wishlistItems.map((item) => (
            <WishlistItemRow key={item.id} item={item} />
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default WishlistTable;