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

  // Caching state variables
  const [allCustomersCache, setAllCustomersCache] = useState<Customer[]>([]);
  const [displayedCustomers, setDisplayedCustomers] = useState<Customer[]>([]);
  const [isDataLoaded, setIsDataLoaded] = useState<boolean>(false);

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
    setIsDataLoaded(false);
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
    setIsDataLoaded(false);
  };

  const handlePageChange = (page: number) => {
    const totalPages = Math.ceil(totalCustomers / filterParams.limit);
    const validPage = Math.max(1, Math.min(page, totalPages));

    setFilterParams({
      ...filterParams,
      page: validPage,
    });
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

        // Fetch all customers for caching
        const { customers: allCustomers, total } = await CustomerService.getCustomers({
          ...filterParams,
          page: 1,
          limit: 1000,
          status: filterParams.status,
        });

        setAllCustomersCache(allCustomers);
        setTotalCustomers(total);
        setIsDataLoaded(true);

        // Get counts for each status filter
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
        console.error("Failed to load customer data:", error);
      } finally {
        setLoading(false);
      }
    };

    if (!isDataLoaded || filterParams.status || filterParams.searchTerm) {
      loadData();
    }
  }, [filterParams.status, filterParams.searchTerm, isDataLoaded]);

  useEffect(() => {
    if (allCustomersCache.length > 0) {
      const startIndex = (filterParams.page - 1) * filterParams.limit;
      const endIndex = startIndex + filterParams.limit;

      const customersForCurrentPage = allCustomersCache.slice(startIndex, endIndex);
      setCustomers(customersForCurrentPage);
      setDisplayedCustomers(customersForCurrentPage);
    }
  }, [filterParams.page, filterParams.limit, allCustomersCache]);

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
                  className={`transition-all duration-500 ease-in-out ${
                    showSidebar ? "w-[788px] mr-5" : "w-full"
                  }`}
                >
                  {loading ? (
                    <div className="flex justify-center py-10">
                      <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
                    </div>
                  ) : (
                    <CustomerTable
                      customers={displayedCustomers}
                      onViewCustomer={handleViewCustomer}
                      selectedCustomerId={selectedCustomer?.id}
                    />
                  )}
                  {!loading && (
                    <div className="mt-[3rem]">
                      {/* {renderPagination()} */}
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