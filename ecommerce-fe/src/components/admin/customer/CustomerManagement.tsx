import React, { useState, useEffect } from "react";
import AdminHeader from "../layout/AdminHeader";
import { CustomerSummaryCard } from "./cards/CustomerSummaryCard";
import CustomerTable from "./tables/CustomerTable";
import CustomerActivityChart from "./charts/CustomerActivityChart";
import CustomerSidebar from "./detail/CustomerSidebar";
import Pagination from '../../common/Pagination';
import {
  Customer,
  CustomerStatus,
  CustomerFilterParams,
  CustomerOverview,
  CustomerChartData,
} from "./models/customer.model";
import { CustomerService } from "./services/customer.service";
import { useToast } from "@/hooks/use-toast";

const CustomerManagement: React.FC = () => {
  // State variables
  const [customers, setCustomers] = useState<Customer[]>([]);
  const [overview, setOverview] = useState<CustomerOverview | null>(null);
  const [chartData, setChartData] = useState<CustomerChartData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [filterParams, setFilterParams] = useState<CustomerFilterParams>({
    page: 1,
    limit: 10,
  });
  const [totalCustomers, setTotalCustomers] = useState<number>(0);
  const [activeStatus, setActiveStatus] = useState<CustomerStatus | "All">("All");
  const [filterCounts, setFilterCounts] = useState({
    all: 0,
    active: 0,
    inactive: 0,
    vip: 0,
  });
  const [selectedCustomer, setSelectedCustomer] = useState<Customer | null>(null);
  const [showSidebar, setShowSidebar] = useState<boolean>(false);
  const { toast } = useToast();

  // Handlers
  const handleViewCustomer = (customer: Customer) => {
    // If the same customer is clicked again, toggle the sidebar
    if (selectedCustomer && selectedCustomer.id === customer.id) {
      setShowSidebar(!showSidebar);
    } else {
      // If a different customer is clicked, select them and ensure sidebar is shown
      setSelectedCustomer(customer);
      setShowSidebar(true);
    }
  };

  const handleSearch = (searchTerm: string) => {
    setFilterParams({
      ...filterParams,
      searchTerm,
      page: 1,
    });
  };

  const handleFilterChange = (filter: string) => {
    let status: CustomerStatus | undefined;

    switch (filter) {
      case "Active":
        status = CustomerStatus.ACTIVE;
        break;
      case "Inactive":
        status = CustomerStatus.INACTIVE;
        break;
      case "VIP":
        status = CustomerStatus.VIP;
        break;
      default:
        status = undefined;
    }

    setActiveStatus(status || "All");
    setFilterParams({
      ...filterParams,
      status,
      page: 1,
    });
  };

  const handlePageChange = (page: number) => {
    const totalPages = Math.ceil(totalCustomers / filterParams.limit);
    const validPage = Math.max(1, Math.min(page, totalPages));

    setFilterParams({
      ...filterParams,
      page: validPage,
    });
  };

  const handleStatusChange = async (customerId: string, newStatus: CustomerStatus) => {
    try {
      const success = await CustomerService.updateCustomerStatus(customerId, newStatus);

      if (success) {
        // Update the customer in the local state
        setCustomers(prevCustomers =>
          prevCustomers.map(customer =>
            customer.id === customerId
              ? { ...customer, status: newStatus }
              : customer
          )
        );

        // If we're updating the status of the selected customer, update that too
        if (selectedCustomer && selectedCustomer.id === customerId) {
          setSelectedCustomer({
            ...selectedCustomer,
            status: newStatus
          });
        }

        toast({
          variant: "success",
          title: "Status Updated",
          description: `Customer status has been updated to ${newStatus}`
        });

        // Refresh the filter counts
        refreshFilterCounts();
      } else {
        toast({
          variant: "destructive",
          title: "Update Failed",
          description: "Could not update customer status"
        });
      }
    } catch (error) {
      console.error("Error updating customer status:", error);
      toast({
        variant: "destructive",
        title: "Update Failed",
        description: "An error occurred while updating customer status"
      });
    }
  };

  const handleDeleteCustomer = async (customerId: string) => {
    try {
      const success = await CustomerService.deleteCustomer(customerId);

      if (success) {
        // Remove the customer from the local state
        setCustomers(prevCustomers => prevCustomers.filter(customer => customer.id !== customerId));

        // If the deleted customer was selected, hide the sidebar
        if (selectedCustomer && selectedCustomer.id === customerId) {
          setSelectedCustomer(null);
          setShowSidebar(false);
        }

        // Decrease the total count
        setTotalCustomers(prevTotal => Math.max(0, prevTotal - 1));

        toast({
          variant: "success",
          title: "Customer Deleted",
          description: "Customer has been successfully deleted"
        });

        // Refresh the filter counts
        refreshFilterCounts();
      } else {
        toast({
          variant: "destructive",
          title: "Delete Failed",
          description: "Could not delete the customer"
        });
      }
    } catch (error) {
      console.error("Error deleting customer:", error);
      toast({
        variant: "destructive",
        title: "Delete Failed",
        description: "An error occurred while deleting the customer"
      });
    }
  };

  // Helper function to refresh filter counts
  const refreshFilterCounts = async () => {
    try {
      const { total: all } = await CustomerService.getCustomers({
        page: 1,
        limit: 1,
      });
      const { total: active } = await CustomerService.getCustomers({
        status: CustomerStatus.ACTIVE,
        page: 1,
        limit: 1,
      });
      const { total: inactive } = await CustomerService.getCustomers({
        status: CustomerStatus.INACTIVE,
        page: 1,
        limit: 1,
      });
      const { total: vip } = await CustomerService.getCustomers({
        status: CustomerStatus.VIP,
        page: 1,
        limit: 1,
      });

      setFilterCounts({
        all,
        active,
        inactive,
        vip,
      });
    } catch (error) {
      console.error("Error fetching filter counts:", error);
    }
  };

  // Effects
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        // Fetch overview data
        const overviewData = await CustomerService.getCustomerOverview();
        setOverview(overviewData);

        // Fetch chart data
        const chartData = await CustomerService.getCustomerChartData();
        setChartData(chartData);

        // Fetch customers with filters
        const { customers, total } = await CustomerService.getCustomers(filterParams);

        setCustomers(customers);
        setTotalCustomers(total);

        // Get counts for each status filter
        await refreshFilterCounts();
      } catch (error) {
        console.error("Failed to load customer data:", error);
        toast({
          variant: "destructive",
          title: "Data Loading Error",
          description: "Failed to load customer data. Please try again later."
        });
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [filterParams, toast]);

  // Component rendering
  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Customer Management" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            {/* Customer Summary Cards */}
            <div className="grid grid-cols-4 grid-rows-3 gap-x-5 gap-y-5">
              {overview ? (
                <>
                  <div className="col-span-1 row-span-1">
                    <CustomerSummaryCard
                      title="Total Customers"
                      value={overview.totalCustomers.count}
                      growthRate={overview.totalCustomers.growth}
                      period={overview.totalCustomers.period}
                      color="primary"
                    />
                  </div>
                  <div className="col-span-1 row-span-1 row-start-2">
                    <CustomerSummaryCard
                      title="New Customers"
                      value={overview.newCustomers.count}
                      growthRate={overview.newCustomers.growth}
                      period={overview.newCustomers.period}
                      color="warning"
                    />
                  </div>
                  <div className="col-span-1 row-span-1 row-start-3">
                    <CustomerSummaryCard
                      title="Visitor"
                      value={overview.visitors.count}
                      growthRate={overview.visitors.growth}
                      period={overview.visitors.period}
                      color="warning"
                    />
                  </div>

                  {/* Customer Activity Chart */}
                  <div className="col-span-3 row-span-3 row-start-1 filter drop-shadow-lg">
                    {loading ? (
                      <div className="flex justify-center py-10">
                        <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-green-500"></div>
                      </div>
                    ) : (
                      <CustomerActivityChart
                        chartData={chartData}
                        overview={overview}
                      />
                    )}
                  </div>
                </>
              ) : (
                <div className="col-span-3 flex justify-center py-10">
                  <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
                </div>
              )}
            </div>

            {/* Customer Table */}
            <div className="w-[1116px] bg-white rounded-lg shadow p-[1rem] mb-6 mt-[1rem] drop-shadow filter">
              <div className="flex ">
                <div
                  className={`transition-all duration-500 ease-in-out ${showSidebar ? "w-[788px] mr-5" : "w-full"
                    }`}
                >
                  {loading ? (
                    <div className="flex justify-center py-10">
                      <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
                    </div>
                  ) : (
                    <CustomerTable
                      customers={customers}
                      onViewCustomer={handleViewCustomer}
                      selectedCustomerId={selectedCustomer?.id}
                      onStatusChange={handleStatusChange}
                      onDeleteCustomer={handleDeleteCustomer}
                      onSearch={handleSearch}
                      onFilterChange={handleFilterChange}
                      counts={filterCounts}
                      activeFilter={activeStatus}
                      loading={loading}
                    />
                  )}
                  {!loading && (
                    <div className="mt-[3rem]">
                      <Pagination
                        currentPage={filterParams.page}
                        totalItems={totalCustomers}
                        pageSize={filterParams.limit}
                        onPageChange={handlePageChange}
                      />
                    </div>
                  )}
                </div>

                {showSidebar && selectedCustomer && (
                  <div className="filter drop-shadow-lg w-[306px] h-[552px] border rounded-lg transition-all duration-800 ease-linear">
                    <CustomerSidebar customer={selectedCustomer} />
                  </div>
                )}
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default CustomerManagement;