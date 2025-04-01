import React from 'react';
import { Customer, CustomerStatus } from '../models/customer.model';

interface CustomerTableProps {
  customers: Customer[];
  onViewCustomer: (customer: Customer) => void;
  selectedCustomerId?: string;
}

const CustomerTable: React.FC<CustomerTableProps> = ({ customers, onViewCustomer, selectedCustomerId }) => {
  return (
    <div className="overflow-x-auto bg-white rounded-lg shadow">
      <table className="min-w-full">
        <thead>
          <tr className="bg-aqua-spring h-[56px] border-b border-gray-200">
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Customer Id
            </th>
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Name
            </th>
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Phone
            </th>
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Order Count
            </th>
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Total Spend
            </th>
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Status
            </th>
            <th className="text-[15px] font-medium px-6 py-3 text-center text-cyprus  tracking-wider">
              Action
            </th>
          </tr>
        </thead>
        <tbody className="bg-white divide-y divide-gray-200">
          {customers.map((customer, index) => (
            <tr 
              key={index} 
              className={`hover:bg-gray-50 cursor-pointer ${customer.id === selectedCustomerId ? 'bg-blue-50' : ''}`}
              onClick={() => onViewCustomer(customer)}
            >
              <td className="flex items-center justify-center h-[64px] px-6 pl-0 py-4  text-center whitespace-nowrap text-[15px]  text-black">
                #{customer.id}
              </td>
              <td className="text-center h-[64px] px-6 py-4 whitespace-nowrap text-[15px]  text-black">
                {customer.name}
              </td>
              <td className="text-center h-[64px] px-6 py-4 whitespace-nowrap text-[15px]  text-black">
                {customer.phone}
              </td>
              <td className="text-center h-[64px] px-6 py-4 whitespace-nowrap text-[15px]  text-black">
                {customer.orderCount}
              </td>
              <td className="text-center h-[64px] px-6 py-4 whitespace-nowrap text-[15px]  text-black">
                ${customer.totalSpend.toFixed(2)}
              </td>
              <td className="text-center h-[64px] px-6 py-4 whitespace-nowrap">
                <span className={`!text-[15px] px-2 inline-flex text-xs leading-5 rounded-full 
                  ${customer.status === CustomerStatus.ACTIVE ? ' text-success' : ''} 
                  ${customer.status === CustomerStatus.INACTIVE ? ' text-error' : ''}
                  ${customer.status === CustomerStatus.VIP ? ' text-pending' : ''}`}>
                  • {customer.status}
                </span>
              </td>
              <td className="text-center h-[64px] px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button 
                  className="text-gray-500 hover:text-gray-700 mr-3"
                  onClick={(e) => {
                    e.stopPropagation();
                    onViewCustomer(customer);
                  }}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                    <path fillRule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clipRule="evenodd" />
                  </svg>
                </button>
                <button 
                  className="text-gray-500 hover:text-gray-700 mr-3"
                  onClick={(e) => e.stopPropagation()}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M13.586 3.586a2 2 0 112.828 2.828l-.793.793-2.828-2.828.793-.793zM11.379 5.793L3 14.172V17h2.828l8.38-8.379-2.83-2.828z" />
                  </svg>
                </button>
                <button 
                  className="text-gray-500 hover:text-gray-700"
                  onClick={(e) => e.stopPropagation()}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                  </svg>
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      
    </div>
  );
};

export default CustomerTable;