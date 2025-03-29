// src/components/orders/modals/AddOrderModal.tsx

import React, { useState } from "react";
import { OrderStatus } from "../models/order.model";

interface Product {
  productId: string;
  productName: string;
  productImage: string;
  quantity: number;
  price: number;
}

export interface NewOrderData {
  customerId: string;
  customerName: string;
  products: Product[];
  totalAmount: number;
  paymentStatus: string;
}

interface AddOrderModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (orderData: NewOrderData) => void;
}

export const AddOrderModal: React.FC<AddOrderModalProps> = ({ 
  isOpen, 
  onClose, 
  onSave 
}) => {
  const [orderData, setOrderData] = useState<NewOrderData>({
    customerId: 'CUST001',
    customerName: 'John Doe',
    products: [
      {
        productId: 'PROD007',
        productName: "Casual Baseball Cap",
        productImage: '/images/products/cap.png',
        quantity: 1,
        price: 49.99,
      }
    ],
    totalAmount: 49.99,
    paymentStatus: 'Paid'
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave(orderData);
  };

  const handleQuantityChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newQuantity = parseInt(e.target.value) || 1;
    const newTotalAmount = orderData.products[0].price * newQuantity;
    
    setOrderData({
      ...orderData,
      products: [
        {
          ...orderData.products[0],
          quantity: newQuantity
        }
      ],
      totalAmount: newTotalAmount
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className="bg-white rounded-lg p-6 w-[500px] max-h-[90vh] overflow-y-auto">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Add New Order</h2>
          <button 
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700"
            aria-label="Close modal"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit}>
          {/* Customer Information */}
          <div className="mb-4">
            <h3 className="font-medium mb-2">Customer Information</h3>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700">Customer ID</label>
                <input 
                  type="text" 
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm"
                  value={orderData.customerId}
                  readOnly
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700">Customer Name</label>
                <input 
                  type="text" 
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm"
                  value={orderData.customerName}
                  readOnly
                />
              </div>
            </div>
          </div>

          {/* Product Information */}
          <div className="mb-4">
            <h3 className="font-medium mb-2">Product Information</h3>
            <div className="flex items-center p-4 border rounded-lg">
              <img 
                src={orderData.products[0].productImage} 
                alt={orderData.products[0].productName} 
                className="w-16 h-16 object-cover mr-4"
              />
              <div className="flex-grow">
                <p className="font-medium">{orderData.products[0].productName}</p>
                <p className="text-sm text-gray-500">Product ID: {orderData.products[0].productId}</p>
                <p className="text-sm font-medium">${orderData.products[0].price.toFixed(2)}</p>
              </div>
              <div className="w-20">
                <label className="block text-sm font-medium text-gray-700">Qty</label>
                <input 
                  type="number" 
                  min="1"
                  className="mt-1 block w-full px-2 py-1 border border-gray-300 rounded-md shadow-sm"
                  value={orderData.products[0].quantity}
                  onChange={handleQuantityChange}
                />
              </div>
            </div>
          </div>

          {/* Order Summary */}
          <div className="mb-4">
            <h3 className="font-medium mb-2">Order Summary</h3>
            <div className="bg-gray-50 p-4 rounded-lg">
              <div className="flex justify-between mb-2">
                <span>Subtotal:</span>
                <span>${orderData.totalAmount.toFixed(2)}</span>
              </div>
              <div className="flex justify-between font-medium">
                <span>Total:</span>
                <span>${orderData.totalAmount.toFixed(2)}</span>
              </div>
            </div>
          </div>

          {/* Payment Status */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-1">Payment Status</label>
            <select 
              className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm"
              value={orderData.paymentStatus}
              onChange={(e) => setOrderData({...orderData, paymentStatus: e.target.value})}
            >
              <option value="Paid">Paid</option>
              <option value="Unpaid">Unpaid</option>
              <option value="Refunded">Refunded</option>
            </select>
          </div>

          <div className="flex justify-end space-x-3">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-ocean-green hover:bg-green-700"
            >
              Create Order
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};