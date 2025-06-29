import React from 'react';
import { OrderProduct } from '../models/order.model';

interface OrderDetailsTableProps {
  products: OrderProduct[];
}

export const OrderDetailsTable: React.FC<OrderDetailsTableProps> = ({ products }) => {
  return (
    <div className="overflow-x-auto">
      <table className="min-w-full bg-white">
        <thead>
          <tr className="bg-gray-50 text-gray-600 text-sm leading-normal">
            <th className="py-3 px-4 text-left">Product</th>
            <th className="py-3 px-4 text-left">Quantity</th>
            <th className="py-3 px-4 text-left">Price</th>
            <th className="py-3 px-4 text-left">Total</th>
          </tr>
        </thead>
        <tbody className="text-gray-600 text-sm">
          {products.map((product, index) => (
            <tr key={index} className="border-b border-gray-100 hover:bg-gray-50">
              <td className="py-3 px-4">
                <div className="flex items-center">
                  <img 
                    src={product.productImage} 
                    alt={product.productName} 
                    className="w-10 h-10 mr-3 object-cover rounded"
                  />
                  <span>{product.productName}</span>
                </div>
              </td>
              <td className="py-3 px-4">{product.quantity}</td>
              <td className="py-3 px-4">${product.price.toFixed(2)}</td>
              <td className="py-3 px-4">${(product.price * product.quantity).toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};