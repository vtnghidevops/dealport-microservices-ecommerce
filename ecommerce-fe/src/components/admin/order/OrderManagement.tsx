import React, { useState, useEffect } from "react";
import AdminHeader from "../layout/AdminHeader";
import { OrderSummaryCard } from "./cards/OrderSummaryCard";
import { OrderTable } from "./tables/OrderTable";
import { OrderFilter } from "./filters/OrderFilter";
import { FiPlusCircle } from "react-icons/fi";

// Import từ services toàn cục
import { Order, OrderStatus } from '@/services/user/order.service';
import { adminOrderService, OrderFilterParams, OrderSummary } from '@/services/admin/order.service';

import { AddOrderModal, NewOrderData } from "./modals/AddOrderModal";
import { useToast } from "@/hooks/use-toast";
import Pagination from "../../common/Pagination";

export const OrderManagement: React.FC = () => {
  const { toast } = useToast();

  // State variables
  const [orders, setOrders] = useState<Order[]>([]);
  const [orderSummary, setOrderSummary] = useState<OrderSummary | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [initialLoading, setInitialLoading] = useState<boolean>(true);

  // Pagination state
  const [currentPage, setCurrentPage] = useState<number>(1);
  const [itemsPerPage, setItemsPerPage] = useState<number>(10);
  const [totalItems, setTotalItems] = useState<number>(0);

  // Filter state
  const [activeFilter, setActiveFilter] = useState<string>("All");
  const [searchTerm, setSearchTerm] = useState<string>("");

  // Data caching
  const [allOrders, setAllOrders] = useState<Order[]>([]);
  const [filteredOrders, setFilteredOrders] = useState<Order[]>([]);
  const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);
  const [statusCounts, setStatusCounts] = useState({
    all: 0,
    paid: 0,
    pending: 0,
    shipped: 0,
    cancelled: 0,
    processing: 0,
    delivered: 0,
    refunded: 0,
  });

  // Load all orders from the API once
  useEffect(() => {
    const loadAllData = async () => {
      setInitialLoading(true);
      try {
        // Fetch summary data
        const summaryData = await adminOrderService.fetchOrderSummary();
        setOrderSummary(summaryData);

        // Fetch all orders (we'll filter client-side)
        const { orders: fetchedOrders, total } = await adminOrderService.fetchOrders({
          page: 1,
          limit: 1000, // Get a large batch to handle locally
        });

        // Filter out invalid orders
        const validOrders = Array.isArray(fetchedOrders)
          ? fetchedOrders.filter(order => order && typeof order === 'object' && order.id)
          : [];

        // Save all orders for client-side filtering
        setAllOrders(validOrders);
        setFilteredOrders(validOrders);
        setTotalItems(validOrders.length);

        // Calculate counts for each status
        updateStatusCounts(validOrders);
      } catch (error) {
        console.error("Failed to load order data:", error);
        toast({
          title: "Failed to load order data",
          description: error instanceof Error ? error.message : "Unknown error",
          variant: "destructive",
        });
      } finally {
        setInitialLoading(false);
      }
    };

    loadAllData();
  }, []);

  // Calculate counts for each status
  const updateStatusCounts = (orders: Order[]) => {
    if (!Array.isArray(orders)) return;

    const counts = {
      all: orders.length,
      paid: orders.filter(order => order.status === OrderStatus.Paid).length,
      pending: orders.filter(order => order.status === OrderStatus.Pending).length,
      shipped: orders.filter(order => order.status === OrderStatus.Shipped).length,
      cancelled: orders.filter(order => order.status === OrderStatus.Cancelled).length,
      processing: orders.filter(order => order.status === OrderStatus.Processing).length,
      delivered: orders.filter(order => order.status === OrderStatus.Delivered).length,
      refunded: orders.filter(order => order.status === OrderStatus.Refunded).length,
    };

    setStatusCounts(counts);
  };

  // Update filtered orders when filter changes
  useEffect(() => {
    setLoading(true);

    // Apply filters (status and search term)
    let filtered = [...allOrders];

    // Apply status filter
    if (activeFilter !== "All") {
      const statusValue = getStatusFromFilter(activeFilter);
      if (statusValue) {
        filtered = filtered.filter(order => order.status === statusValue);
      }
    }

    // Apply search filter
    if (searchTerm) {
      const search = searchTerm.toLowerCase().trim();
      filtered = filtered.filter(order => {
        // Tìm kiếm trong thông tin đơn hàng cơ bản
        if ((order.id && order.id.toLowerCase().includes(search)) ||
          (order.orderNumber && order.orderNumber.toLowerCase().includes(search)) ||
          (order.notes && order.notes.toLowerCase().includes(search)) ||
          (order.paymentMethod && order.paymentMethod.toLowerCase().includes(search))) {
          return true;
        }

        // Tìm kiếm trong thông tin billing
        if (order.billingInfo) {
          if ((order.billingInfo.firstName && order.billingInfo.firstName.toLowerCase().includes(search)) ||
            (order.billingInfo.lastName && order.billingInfo.lastName.toLowerCase().includes(search)) ||
            (order.billingInfo.email && order.billingInfo.email.toLowerCase().includes(search)) ||
            (order.billingInfo.phone && order.billingInfo.phone.includes(search)) ||
            (order.billingInfo.address && order.billingInfo.address.toLowerCase().includes(search)) ||
            (order.billingInfo.city && order.billingInfo.city.toLowerCase().includes(search)) ||
            (order.billingInfo.country && order.billingInfo.country.toLowerCase().includes(search))) {
            return true;
          }
        }

        // Tìm kiếm trong thông tin shipping
        if (order.shippingInfo) {
          if ((order.shippingInfo.firstName && order.shippingInfo.firstName.toLowerCase().includes(search)) ||
            (order.shippingInfo.lastName && order.shippingInfo.lastName.toLowerCase().includes(search)) ||
            (order.shippingInfo.address && order.shippingInfo.address.toLowerCase().includes(search)) ||
            (order.shippingInfo.city && order.shippingInfo.city.toLowerCase().includes(search)) ||
            (order.shippingInfo.country && order.shippingInfo.country.toLowerCase().includes(search)) ||
            (order.shippingInfo.shippingMethod && order.shippingInfo.shippingMethod.toLowerCase().includes(search))) {
            return true;
          }
        }

        // Tìm kiếm trong tên sản phẩm
        if (order.items && Array.isArray(order.items)) {
          return order.items.some(item =>
            (item.name && item.name.toLowerCase().includes(search)) ||
            (item.productId && item.productId.toLowerCase().includes(search))
          );
        }

        return false;
      });
    }

    setFilteredOrders(filtered);
    setTotalItems(filtered.length);
    setCurrentPage(1); // Reset to first page when filter changes
    setLoading(false);
  }, [activeFilter, searchTerm, allOrders]);

  // Update displayed orders when page changes
  useEffect(() => {
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = Math.min(startIndex + itemsPerPage, filteredOrders.length);

    if (startIndex >= filteredOrders.length) {
      setOrders([]);
    } else {
      setOrders(filteredOrders.slice(startIndex, endIndex));
    }
  }, [currentPage, itemsPerPage, filteredOrders]);

  // Convert filter name to OrderStatus
  const getStatusFromFilter = (filter: string): OrderStatus | undefined => {
    switch (filter) {
      case "Paid":
        return OrderStatus.Paid;
      case "Pending":
        return OrderStatus.Pending;
      case "Processing":
        return OrderStatus.Processing;
      case "Shipped":
        return OrderStatus.Shipped;
      case "Delivered":
        return OrderStatus.Delivered;
      case "Cancelled":
        return OrderStatus.Cancelled;
      case "Refunded":
        return OrderStatus.Refunded;
      default:
        return undefined;
    }
  };

  // Handlers
  const handleCreateOrder = async (orderData: NewOrderData) => {
    setLoading(true);
    try {
      const newOrder = await adminOrderService.createOrder({
        customerId: orderData.customerId,
        customerName: orderData.customerName,
        products: orderData.products,
        totalAmount: orderData.totalAmount,
        paymentStatus: orderData.paymentStatus,
      });

      // Add the new order to our cached orders
      const updatedAllOrders = [newOrder, ...allOrders];
      setAllOrders(updatedAllOrders);

      // Update filtered orders if the new order matches current filter
      if (activeFilter === "All" || newOrder.status === getStatusFromFilter(activeFilter)) {
        const updatedFiltered = [newOrder, ...filteredOrders];
        setFilteredOrders(updatedFiltered);
        setTotalItems(updatedFiltered.length);
      }

      // Update status counts
      updateStatusCounts(updatedAllOrders);

      setIsAddModalOpen(false);
      toast({
        title: "Order created successfully!",
        variant: "success",
      });
    } catch (error) {
      console.error("Error creating order:", error);
      toast({
        title: "Failed to create order",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleStatusChange = async (orderId: string, status: string) => {
    const statusValue = status as OrderStatus;
    setLoading(true);

    try {
      const success = await adminOrderService.updateOrderStatus(orderId, statusValue);
      if (success) {
        // Update both allOrders and filteredOrders
        const updatedAllOrders = allOrders.map(order => {
          if (order.id === orderId) {
            const updatedOrder = { ...order, status: statusValue };

            // Tự động cập nhật payment status nếu order status là paid
            if (statusValue === OrderStatus.Paid) {
              if (!updatedOrder.paymentStatus ||
                updatedOrder.paymentStatus === 'pending' ||
                updatedOrder.paymentStatus === 'processing') {
                updatedOrder.paymentStatus = 'paid';
              }

              if (updatedOrder.paymentInfo) {
                updatedOrder.paymentInfo = {
                  ...updatedOrder.paymentInfo,
                  status: 'paid'
                };
              }
            }

            return updatedOrder;
          }
          return order;
        });

        setAllOrders(updatedAllOrders);

        // If the order no longer matches the current filter, remove it from filtered orders
        if (activeFilter !== "All" && getStatusFromFilter(activeFilter) !== statusValue) {
          const updatedFiltered = filteredOrders.filter(order => order.id !== orderId);
          setFilteredOrders(updatedFiltered);
          setTotalItems(updatedFiltered.length);
        } else {
          // Otherwise update it in filtered orders
          const updatedFiltered = filteredOrders.map(order => {
            if (order.id === orderId) {
              const updatedOrder = { ...order, status: statusValue };

              // Tự động cập nhật payment status nếu order status là paid
              if (statusValue === OrderStatus.Paid) {
                if (!updatedOrder.paymentStatus ||
                  updatedOrder.paymentStatus === 'pending' ||
                  updatedOrder.paymentStatus === 'processing') {
                  updatedOrder.paymentStatus = 'paid';
                }

                if (updatedOrder.paymentInfo) {
                  updatedOrder.paymentInfo = {
                    ...updatedOrder.paymentInfo,
                    status: 'paid'
                  };
                }
              }

              return updatedOrder;
            }
            return order;
          });
          setFilteredOrders(updatedFiltered);
        }

        // Update status counts
        updateStatusCounts(updatedAllOrders);

        toast({
          title: "Order status updated",
          description: `Order status changed to ${status}`,
          variant: "success",
        });
      } else {
        toast({
          title: "Failed to update order status",
          variant: "destructive",
        });
      }
    } catch (error) {
      console.error("Error updating order status:", error);
      toast({
        title: "Error updating order status",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleViewDetails = (orderId: string) => {
    const order = orders.find(o => o.id === orderId);
    const orderIdentifier = order?.orderNumber || orderId.substring(0, 8);
    console.log(`View details for order: ${orderIdentifier}`);
    // Implement view details functionality
  };

  const handleSearch = (term: string) => {
    setSearchTerm(term);
  };

  const handleFilterChange = (filter: string) => {
    setActiveFilter(filter);
  };

  const handlePageChange = (page: number) => {
    const totalPages = Math.ceil(totalItems / itemsPerPage);
    const validPage = Math.max(1, Math.min(page, totalPages));
    setCurrentPage(validPage);
  };

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
              {orderSummary && !initialLoading ? (
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
                counts={statusCounts}
                loading={loading || initialLoading}
              />

              {loading || initialLoading ? (
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
                  <div className="mt-[3rem]">
                    <Pagination
                      currentPage={currentPage}
                      totalItems={totalItems}
                      pageSize={itemsPerPage}
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
