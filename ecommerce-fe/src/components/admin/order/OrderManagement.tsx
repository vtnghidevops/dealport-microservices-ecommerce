import React, { useState, useEffect, useCallback, useRef } from "react";
import AdminHeader from "../layout/AdminHeader";
import { OrderSummaryCard } from "./cards/OrderSummaryCard";
import { OrderTable } from "./tables/OrderTable";
import { OrderFilter } from "./filters/OrderFilter";
import { FiPlusCircle } from "react-icons/fi";

// Import from services
import { Order, OrderStatus } from '@/services/user/order.service';
import { adminOrderService, OrderSummary } from '@/services/admin/order.service';

import { AddOrderModal, NewOrderData } from "./modals/AddOrderModal";
import { useToast } from "@/hooks/use-toast";
import Pagination from "../../common/Pagination";

export const OrderManagement: React.FC = () => {
  const { toast } = useToast();

  // State variables
  const [orders, setOrders] = useState<Order[]>([]);
  const [allOrders, setAllOrders] = useState<Order[]>([]);
  const [orderSummary, setOrderSummary] = useState<OrderSummary | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [initialLoading, setInitialLoading] = useState<boolean>(true);

  // Pagination state
  const [currentPage, setCurrentPage] = useState<number>(1);
  const itemsPerPage = 10;
  const [totalItems, setTotalItems] = useState<number>(0);

  // Filter state
  const [activeFilter, setActiveFilter] = useState<string>("All");
  const [searchTerm, setSearchTerm] = useState<string>("");

  // Filter results
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

  // Use ref to track if data has been loaded
  const dataHasBeenLoaded = useRef(false);

  // Calculate counts for each status - memoized
  const updateStatusCounts = useCallback((orders: Order[]) => {
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
  }, []);

  // Fetch orders - not memoized to avoid dependencies in useEffect
  const fetchOrders = async () => {
    console.log("Fetching admin orders...");
    try {
      // Fetch all orders (we'll filter client-side)
      const { orders: fetchedOrders } = await adminOrderService.fetchOrders({
        page: 1,
        limit: 1000 // Get a large batch to handle locally
      });

      // Filter out invalid orders
      const validOrders = Array.isArray(fetchedOrders)
        ? fetchedOrders.filter(order => order && typeof order === 'object' && order.id)
        : [];

      console.log(`Fetched ${validOrders.length} valid orders`);

      // Save all orders for client-side filtering
      setAllOrders(validOrders);
      setFilteredOrders(validOrders);
      setTotalItems(validOrders.length);

      // Calculate counts for each status
      updateStatusCounts(validOrders);

      return validOrders;
    } catch (error) {
      console.error("Failed to load order data:", error);
      toast({
        title: "Failed to load order data",
        description: error instanceof Error ? error.message : "Unknown error",
        variant: "destructive",
      });
      return [];
    }
  };

  // Load order summary data
  const fetchOrderSummary = async () => {
    try {
      console.log("Fetching order summary...");
      const summaryData = await adminOrderService.fetchOrderSummary();
      setOrderSummary(summaryData);
      return summaryData;
    } catch (error) {
      console.error("Failed to load order summary:", error);
      toast({
        title: "Failed to load order summary",
        description: error instanceof Error ? error.message : "Unknown error",
        variant: "destructive",
      });
      return null;
    }
  };

  // Load all orders and summary on component mount ONLY
  useEffect(() => {
    const loadAllData = async () => {
      if (dataHasBeenLoaded.current) {
        console.log("Data already loaded, skipping fetch");
        setInitialLoading(false);
        return;
      }

      console.log("Loading order data for the first time");
      setInitialLoading(true);

      try {
        dataHasBeenLoaded.current = true;

        // Parallel loading of summary and orders
        await Promise.all([
          fetchOrderSummary(),
          fetchOrders()
        ]);
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
    // No dependencies to avoid re-execution
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Update filtered orders when filter changes - only run when allOrders or filters change
  useEffect(() => {
    if (!allOrders.length) return; // Skip if no orders loaded

    console.log("Applying filters to orders");
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
        // Search in basic order information
        if ((order.id && order.id.toLowerCase().includes(search)) ||
          (order.orderNumber && order.orderNumber.toLowerCase().includes(search)) ||
          (order.notes && order.notes.toLowerCase().includes(search)) ||
          (order.paymentMethod && order.paymentMethod.toLowerCase().includes(search))) {
          return true;
        }

        // Search in billing info
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

        // Search in shipping info
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

        // Search in product names
        if (order.items && Array.isArray(order.items)) {
          return order.items.some((item: any) =>
            (item.name && item.name.toLowerCase().includes(search)) ||
            (item.productId && item.productId.toLowerCase().includes(search))
          );
        }

        return false;
      });
    }

    console.log(`Filter applied: ${filtered.length} orders match criteria`);
    setFilteredOrders(filtered);
    setTotalItems(filtered.length);
    setCurrentPage(1); // Reset to first page when filter changes
    setLoading(false);
  }, [activeFilter, searchTerm, allOrders]);

  // Update displayed orders when page changes
  useEffect(() => {
    if (!filteredOrders.length) return;

    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = Math.min(startIndex + itemsPerPage, filteredOrders.length);

    if (startIndex >= filteredOrders.length) {
      setOrders([]);
    } else {
      const paginatedOrders = filteredOrders.slice(startIndex, endIndex);
      console.log(`Displaying orders ${startIndex + 1}-${endIndex} of ${filteredOrders.length}`);
      setOrders(paginatedOrders);
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

  // Handler for creating order
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

      console.log("New order created:", newOrder.id);

      // Add the new order to our data
      const updatedOrders = [newOrder, ...allOrders];
      setAllOrders(updatedOrders);

      // Update counts
      updateStatusCounts(updatedOrders);

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

  // Handlers
  const handleStatusChange = async (orderId: string, status: string) => {
    const statusValue = status as OrderStatus;
    setLoading(true);

    try {
      const success = await adminOrderService.updateOrderStatus(orderId, statusValue);
      if (success) {
        // Update both allOrders and filteredOrders
        const updatedAllOrders = allOrders.map((order: Order) => {
          if (order.id === orderId) {
            const updatedOrder = { ...order, status: statusValue };

            // Automatically update payment status if order status is paid
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

        // Refresh the orders list
        fetchOrders();

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
