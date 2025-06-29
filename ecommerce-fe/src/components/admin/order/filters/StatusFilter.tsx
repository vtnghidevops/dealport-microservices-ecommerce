import React from 'react';
import { OrderStatus } from '../models/order.model';

interface StatusFilterProps {
  onStatusChange: (status: OrderStatus | 'All') => void;
  currentStatus: OrderStatus | 'All';
}

export const StatusFilter: React.FC<StatusFilterProps> = ({ onStatusChange, currentStatus }) => {
  return (
    <div className="flex space-x-2 mb-4 overflow-x-auto pb-2">
      <button
        className={`px-4 py-2 rounded-md whitespace-nowrap ${
          currentStatus === 'All' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-800'
        }`}
        onClick={() => onStatusChange('All')}
      >
        All
      </button>
      <button
        className={`px-4 py-2 rounded-md whitespace-nowrap ${
          currentStatus === 'Delivered' ? 'bg-green-600 text-white' : 'bg-gray-100 text-gray-800'
        }`}
        onClick={() => onStatusChange('Delivered')}
      >
        Delivered
      </button>
      <button
        className={`px-4 py-2 rounded-md whitespace-nowrap ${
          currentStatus === 'Pending' ? 'bg-yellow-600 text-white' : 'bg-gray-100 text-gray-800'
        }`}
        onClick={() => onStatusChange('Pending')}
      >
        Pending
      </button>
      <button
        className={`px-4 py-2 rounded-md whitespace-nowrap ${
          currentStatus === 'Shipped' ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-800'
        }`}
        onClick={() => onStatusChange('Shipped')}
      >
        Shipped
      </button>
      <button
        className={`px-4 py-2 rounded-md whitespace-nowrap ${
          currentStatus === 'Cancelled' ? 'bg-red-600 text-white' : 'bg-gray-100 text-gray-800'
        }`}
        onClick={() => onStatusChange('Cancelled')}
      >
        Cancelled
      </button>
    </div>
  );
};