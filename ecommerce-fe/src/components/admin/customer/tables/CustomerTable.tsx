import React, { useState } from 'react';
import { Customer, CustomerStatus } from '../models/customer.model';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';

interface CustomerTableProps {
  customers: Customer[];
  onViewCustomer: (customer: Customer) => void;
  selectedCustomerId?: string;
  onStatusChange?: (customerId: string, newStatus: CustomerStatus) => void;
  onSearch?: (searchTerm: string) => void;
  onFilterChange?: (filter: string) => void;
  onDeleteCustomer?: (customerId: string) => void;
  counts?: {
    all: number;
    active: number;
    inactive: number;
    vip: number;
  };
  activeFilter?: CustomerStatus | "All";
  loading?: boolean;
}

const CustomerTable: React.FC<CustomerTableProps> = ({
  customers,
  onViewCustomer,
  selectedCustomerId,
  onStatusChange,
  onSearch,
  onFilterChange,
  onDeleteCustomer,
  counts = { all: 0, active: 0, inactive: 0, vip: 0 },
  activeFilter = "All",
  loading = false
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusDropdownOpen, setStatusDropdownOpen] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [customerToDelete, setCustomerToDelete] = useState<Customer | null>(null);

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSearch) {
      onSearch(searchTerm);
    }
  };

  const handleFilterClick = (filter: string) => {
    if (onFilterChange) {
      onFilterChange(filter);
    }
  };

  const handleStatusChange = (e: React.MouseEvent, customerId: string, newStatus: CustomerStatus) => {
    e.stopPropagation();
    if (onStatusChange) {
      onStatusChange(customerId, newStatus);
    }
    setStatusDropdownOpen(null);
  };

  const toggleStatusDropdown = (e: React.MouseEvent, customerId: string) => {
    e.stopPropagation();
    setStatusDropdownOpen(statusDropdownOpen === customerId ? null : customerId);
  };

  const handleDeleteClick = (e: React.MouseEvent, customer: Customer) => {
    e.stopPropagation();
    setCustomerToDelete(customer);
    setDeleteDialogOpen(true);
  };

  const confirmDelete = () => {
    if (customerToDelete && onDeleteCustomer) {
      onDeleteCustomer(customerToDelete.id);
      setDeleteDialogOpen(false);
      setCustomerToDelete(null);
    }
  };

  return (
    <div className="overflow-hidden">
      {/* Search and Filter Bar */}
      <div className="mb-5 p-4 bg-white rounded-t-lg ">
        <div className="flex flex-wrap items-center justify-between">
          <div className="w-full md:w-1/3 mb-5 md:mb-0">
            <form onSubmit={handleSearchSubmit} className="relative">
              <input
                type="text"
                className="w-full rounded-md border-2 border-neutral-200 shadow-sm pl-10 pr-4 py-2 focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
                placeholder="Search customers..."
                value={searchTerm}
                onChange={handleSearchChange}
              />
              <div className="absolute left-3 top-2.5 text-gray-400">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <button type="submit" className="hidden">Search</button>
            </form>
          </div>

          <div className="w-full md:w-2/3 flex justify-end">
            <div className="flex space-x-2">
              <button
                onClick={() => handleFilterClick("All")}
                className={`px-5 py-2 text-sm font-medium rounded-md ${activeFilter === "All"
                  ? "bg-blue-100 text-blue-700"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                  }`}
              >
                All ({counts.all})
              </button>
              <button
                onClick={() => handleFilterClick("Active")}
                className={`px-5 py-2 text-sm font-medium rounded-md ${activeFilter === CustomerStatus.ACTIVE
                  ? "bg-green-100 text-green-700"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                  }`}
              >
                Active ({counts.active})
              </button>
              <button
                onClick={() => handleFilterClick("Inactive")}
                className={`px-5 py-2 text-sm font-medium rounded-md ${activeFilter === CustomerStatus.INACTIVE
                  ? "bg-red-100 text-red-700"
                  : "bg-gray-100 text-gray-700 hover:bg-gray-200"
                  }`}
              >
                Inactive ({counts.inactive})
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Table */}
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
            {loading ? (
              <tr>
                <td colSpan={7} className="text-center py-8">
                  <div className="flex justify-center">
                    <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
                  </div>
                </td>
              </tr>
            ) : customers.length === 0 ? (
              <tr>
                <td colSpan={7} className="text-center py-8 text-gray-500">
                  No customers found
                </td>
              </tr>
            ) : (
              customers.map((customer, index) => (
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
                      title="View details"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                        <path fillRule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clipRule="evenodd" />
                      </svg>
                    </button>

                    {/* Status Change Dropdown */}
                    <div className="relative inline-block text-left mr-3">
                      <button
                        className="text-gray-500 hover:text-gray-700 focus:outline-none"
                        onClick={(e) => toggleStatusDropdown(e, customer.id)}
                        title="Change status"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                      {statusDropdownOpen === customer.id && (
                        <div className="absolute right-0 mt-2 w-[8rem] flex justify-center items-center rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-10">
                          <div className="py-1" role="menu" aria-orientation="vertical">
                            <button
                              className="block w-[8rem] text-center py-2 text-sm text-gray-700 hover:bg-gray-100"
                              onClick={(e) => handleStatusChange(e, customer.id, CustomerStatus.ACTIVE)}
                            >
                              Set Active
                            </button>
                            <button
                              className="block text-center w-[8rem] py-2 text-sm text-gray-700 hover:bg-gray-100"
                              onClick={(e) => handleStatusChange(e, customer.id, CustomerStatus.INACTIVE)}
                            >
                              Set Inactive
                            </button>
                            <button
                              className="block  text-center w-[8rem] py-2 text-sm text-gray-700 hover:bg-gray-100"
                              onClick={(e) => handleStatusChange(e, customer.id, CustomerStatus.VIP)}
                            >
                              Set VIP
                            </button>
                          </div>
                        </div>
                      )}
                    </div>

                    <button
                      className="text-gray-500 hover:text-red-700"
                      onClick={(e) => handleDeleteClick(e, customer)}
                      title="Delete Customer"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                      </svg>
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Delete Confirmation Dialog */}
      <Dialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Confirm Customer Deletion</DialogTitle>
            <DialogDescription>
              Are you sure you want to delete customer {customerToDelete?.name}? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex space-x-2 justify-end">
            <Button
              variant="outline"
              onClick={() => setDeleteDialogOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={confirmDelete}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default CustomerTable;