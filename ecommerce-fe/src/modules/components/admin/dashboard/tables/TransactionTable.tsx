import React from 'react';
import { Transaction } from '../../../sections/AdminDashboard/models/transaction.model';

interface TransactionTableProps {
  transactions: Transaction[];
}

const TransactionTable: React.FC<TransactionTableProps> = ({ transactions }) => {
  return (
    <div className="bg-white rounded-lg p-6  mt-[1.5rem]">      
      <div className="overflow-x-auto">
        <table className="min-w-full table-fixed">
          <thead>
            <tr className="border-b border-[#D1D1D1] w-[808px]">
              <th className="py-3 text-left text-sm font-medium text-gray-500 w-[18%]">
                No
              </th>
              <th className="py-3 text-left text-sm font-medium text-gray-500 w-[20%]">
                Id Customer
              </th>
              <th className="py-3 text-left text-sm font-medium text-gray-500 w-[30%]">
                Order Date
              </th>
              <th className="py-3 text-left text-sm font-medium text-gray-500 w-[20%]">
                Status
              </th>
              <th className="py-3 text-left text-sm font-medium text-gray-500 w-[20%]">
                Amount
              </th>
            </tr>
          </thead>
          <tbody>
            {transactions.map((transaction, index) => (
              <tr key={transaction.id} className="text-[14px] font-bold border-b border-white hover:bg-gray-50 h-[41px]">
                <td className="py-4 text-sm text-gray-500">
                  {index + 1}.
                </td>
                <td className="py-4">
                  <div className="text-[14px] font-medium text-gray-900">{transaction.customerId}</div>
                </td>
                <td className="py-4">
                  <div className="text-[14px] font-medium text-gray-700">{transaction.orderDate}</div>
                </td>
                <td className="py-4">
                  {transaction.status === 'Paid' ? (
                    <div className="flex items-center">
                      <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                      <span className="text-sm font-medium text-gray-700">Paid</span>
                    </div>
                  ) : transaction.status === 'Pending' ? (
                    <div className="flex items-center">
                      <span className="w-2 h-2 bg-yellow-500 rounded-full mr-2"></span>
                      <span className="text-sm font-medium text-gray-700">Pending</span>
                    </div>
                  ) : (
                    <div className="flex items-center">
                      <span className="w-2 h-2 bg-red-500 rounded-full mr-2"></span>
                      <span className="font-medium text-sm text-gray-700">{transaction.status}</span>
                    </div>
                  )}
                </td>
                <td className="py-4 text-sm text-gray-700">
                  ${transaction.amount}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default TransactionTable;