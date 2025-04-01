import React from 'react';
import { CustomerStatus } from '../models/customer.model';

interface StatusFilterProps {
  currentStatus: CustomerStatus | null;
  onStatusChange: (status: CustomerStatus | null) => void;
}

const StatusFilter: React.FC<StatusFilterProps> = ({ currentStatus, onStatusChange }) => {
  return (
    <div className="flex space-x-2 mb-4">
      <button
        className={`px-3 py-1 text-sm rounded-full ${!currentStatus ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'}`}
        onClick={() => onStatusChange(null)}
      >
        All
      </button>
      <button
        className={`px-3 py-1 text-sm rounded-full ${currentStatus === CustomerStatus.ACTIVE ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'}`}
        onClick={() => onStatusChange(CustomerStatus.ACTIVE)}
      >
        Active
      </button>
      <button
        className={`px-3 py-1 text-sm rounded-full ${currentStatus === CustomerStatus.INACTIVE ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'}`}
        onClick={() => onStatusChange(CustomerStatus.INACTIVE)}
      >
        Inactive
      </button>
      <button
        className={`px-3 py-1 text-sm rounded-full ${currentStatus === CustomerStatus.VIP ? 'bg-green-500 text-white' : 'bg-gray-100 text-gray-700'}`}
        onClick={() => onStatusChange(CustomerStatus.VIP)}
      >
        VIP
      </button>
    </div>
  );
};

export default StatusFilter;