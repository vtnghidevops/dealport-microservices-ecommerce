import React, { useState, useEffect, JSX } from "react";
import AdminHeader from "../layout/AdminHeader";
import { OrderSummaryCard } from "./cards/OrderSummaryCard";
import { OrderTable } from "./tables/OrderTable";
import { OrderFilter } from "./filters/OrderFilter";
import { FiPlusCircle } from "react-icons/fi";
import { IoMdArrowRoundBack, IoMdArrowRoundForward } from "react-icons/io";
import {
  Order,
  OrderStatus,
  OrderFilterParams,
  OrderSummary,
} from "./models/order.model";
import { orderService } from "./services/order.service";
import { AddOrderModal, NewOrderData } from "./modals/AddOrderModal";
import { showSuccess, showError } from "../../../utils/notifications";
import Pagination from "../../common/Pagination";

export const OrderManagement: React.FC = () => {
  // State variables
  const [orders, setOrders] = useState<Order[]>([]);
  const [orderSummary, setOrderSummary] = useState<OrderSummary | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [filterParams, setFilterParams] = useState<OrderFilterParams>({
    page: 1,
    limit: 10,
  });
  const [totalOrders, setTotalOrders] = useState<number>(0);
  const [activeStatus, setActiveStatus] = useState<OrderStatus | "All">("All");
  const [filterCounts, setFilterCounts] = useState({
    all: 0,
    completed: 0,
    pending: 0,
    shipped: 0,
    cancelled: 0,
  });
  const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);

  // Caching state variables
  const [allOrdersCache, setAllOrdersCache] = useState<Order[]>([]);
  const [displayedOrders, setDisplayedOrders] = useState<Order[]>([]);
  const [isDataLoaded, setIsDataLoaded] = useState<boolean>(false);

  // Fetching functions
  const fetchOrders = async (params: OrderFilterParams) => {
    const result = await orderService.fetchOrders(params);
    return result;
  };

  const fetchOrderSummary = async () => {
    const summary = await orderService.fetchOrderSummary();
    return summary;
  };

  const updateOrderStatus = async (orderId: string, status: OrderStatus) => {
    return await orderService.updateOrderStatus(orderId, status);
  };

  // Handlers
  const handleCreateOrder = async (orderData: NewOrderData) => {
    setLoading(true);
    try {
      const newOrder = await orderService.createOrder({
        customerId: orderData.customerId,
        customerName: orderData.customerName,
        products: orderData.products,
        totalAmount: orderData.totalAmount,
        paymentStatus: orderData.paymentStatus,
      });

      const updatedOrders = [newOrder, ...allOrdersCache];
      setAllOrdersCache(updatedOrders);
      setTotalOrders(totalOrders + 1);

      if (filterParams.page === 1) {
        const ordersToShow = updatedOrders.slice(0, filterParams.limit);
        setOrders(ordersToShow);
        setDisplayedOrders(ordersToShow);
      }

      setIsAddModalOpen(false);
      showSuccess("Order created successfully!");
    } catch (error) {
      console.error("Error creating order:", error);
      showError("Failed to create order");
    } finally {
      setLoading(false);
    }
  };

  const handleStatusChange = async (orderId: string, status: OrderStatus) => {
    const success = await updateOrderStatus(orderId, status);
    if (success) {
      const updatedOrders = orders.map((order) =>
        order.id === orderId ? { ...order, status } : order
      );
      setOrders(updatedOrders);

      setAllOrdersCache(
        allOrdersCache.map((order) =>
          order.id === orderId ? { ...order, status } : order
        )
      );
    }
  };

  const handleViewDetails = (orderId: string) => {
    console.log(`View details for order: ${orderId}`);
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
    let status: OrderStatus | undefined;

    switch (filter) {
      case "Completed":
        status = "Delivered";
        break;
      case "Pending":
        status = "Pending";
        break;
      case "Cancelled":
        status = "Cancelled";
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
    const totalPages = Math.ceil(totalOrders / filterParams.limit);
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
        const summaryData = await fetchOrderSummary();
        setOrderSummary(summaryData);

        const { orders: allOrders, total } = await fetchOrders({
          ...filterParams,
          page: 1,
          limit: 1000,
          status: filterParams.status,
        });

        setAllOrdersCache(allOrders);
        setTotalOrders(total);
        setIsDataLoaded(true);

        const { total: all } = await orderService.fetchOrders({
          page: 1,
          limit: 1,
        });
        const { total: completed } = await orderService.fetchOrders({
          status: "Delivered",
          page: 1,
          limit: 1,
        });
        const { total: pending } = await orderService.fetchOrders({
          status: "Pending",
          page: 1,
          limit: 1,
        });
        const { total: shipped } = await orderService.fetchOrders({
          status: "Shipped",
          page: 1,
          limit: 1,
        });
        const { total: cancelled } = await orderService.fetchOrders({
          status: "Cancelled",
          page: 1,
          limit: 1,
        });

        setFilterCounts({
          all,
          completed,
          pending,
          shipped,
          cancelled,
        });
      } catch (error) {
        console.error("Failed to load order data:", error);
      } finally {
        setLoading(false);
      }
    };

    if (!isDataLoaded || filterParams.status || filterParams.searchTerm) {
      loadData();
    }
  }, [filterParams.status, filterParams.searchTerm, isDataLoaded]);

  useEffect(() => {
    if (allOrdersCache.length > 0) {
      const startIndex = (filterParams.page - 1) * filterParams.limit;
      const endIndex = startIndex + filterParams.limit;

      const ordersForCurrentPage = allOrdersCache.slice(startIndex, endIndex);
      setOrders(ordersForCurrentPage);
      setDisplayedOrders(ordersForCurrentPage);
    }
  }, [filterParams.page, filterParams.limit, allOrdersCache]);


  // Component rendering
  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Order Management" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            <div className="flex justify-between items-center mb-6">
              <h1 className="text-2xl font-semibold">Order List</h1>
              <button
                onClick={() => setIsAddModalOpen(true)}
                className="rounded-lg h-[48px] w-[142px] py-6 px-8 bg-ocean-green hover:bg-green-700"
              >
                <span className="flex items-center justify-center text-white w-full text-[15px] font-medium ">
                  <FiPlusCircle className="mr-[0.5rem] text-[20px]" />
                  Add Order
                </span>
              </button>
            </div>

            {/* Order Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6 mt-[3rem]">
              {orderSummary ? (
                <>
                  <OrderSummaryCard
                    title="Total Orders"
                    value={orderSummary.totalOrders}
                    growthRate={orderSummary.growthRate.total}
                    period={orderSummary.lastUpdated}
                    color="primary"
                  />
                  <OrderSummaryCard
                    title="New Orders"
                    value={orderSummary.newOrders}
                    growthRate={orderSummary.growthRate.new}
                    period={orderSummary.lastUpdated}
                    color="warning"
                  />
                  <OrderSummaryCard
                    title="Completed Orders"
                    value={orderSummary.completedOrders}
                    growthRate={orderSummary.growthRate.completed}
                    period={orderSummary.lastUpdated}
                    color="success"
                  />
                  <OrderSummaryCard
                    title="Canceled Orders"
                    value={orderSummary.cancelledOrders}
                    growthRate={orderSummary.growthRate.cancelled}
                    period={orderSummary.lastUpdated}
                    color="error"
                  />
                </>
              ) : (
                <div className="col-span-4 flex justify-center py-10">
                  <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
                </div>
              )}
            </div>

            {/* Order Table and Filters */}
            <div className="w-[1116px] h-[979px] bg-white rounded-lg shadow p-[1rem] mb-6 mt-[1rem] drop-shadow filter">
              <OrderFilter
                onSearch={handleSearch}
                onFilterChange={handleFilterChange}
                counts={filterCounts}
                loading={loading}
              />
              {/* <OrderFilter
                onSearch={handleSearch}
                onFilterChange={handleFilterChange}
                counts={filterCounts}
                activeFilter={activeStatus}
                loading={loading}
              /> */}

              {loading ? (
                <div className="flex justify-center py-10">
                  <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-500"></div>
                </div>
              ) : (
                <>
                  <OrderTable
                    orders={orders}
                    onStatusChange={handleStatusChange}
                    onViewDetails={handleViewDetails}
                  />
                  {/* <OrderTable
                    orders={orders}
                    onStatusChange={handleStatusChange}
                    onViewDetails={handleViewDetails}
                    totalItems={totalOrders}
                    currentPage={filterParams.page}
                    pageSize={filterParams.limit}
                    onPageChange={handlePageChange}
                    loading={loading}
                  /> */}
                  <div className="mt-[3rem]">
                    {/* {renderPagination()} */}
                    <Pagination
                        currentPage={filterParams.page}
                        totalItems={totalOrders}
                        pageSize={filterParams.limit}
                        onPageChange={handlePageChange}
                      />
                  </div>
                </>
              )}
            </div>
          </div>
          <AddOrderModal
            isOpen={isAddModalOpen}
            onClose={() => setIsAddModalOpen(false)}
            onSave={handleCreateOrder}
          />
        </main>
      </div>
    </div>
  );
};
