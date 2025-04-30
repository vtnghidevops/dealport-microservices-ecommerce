// pages/user/Orders.tsx
import React, { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';
import UserLayout from '../../components/layouts/UserLayout';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../components/ui/tabs';
import OrderCard from '../../components/user/OrderCard';
import { Button } from '../../components/ui/button';
import { useNavigate } from 'react-router-dom';
import Loading from '@/components/shared/Loading';
import orderService, { Order } from '@/services/user/order.service';
import { useToast } from '@/hooks/use-toast';

// Define a possible response type that includes a data property
interface OrdersResponse {
  data?: Order[];
  [key: string]: any;
}

const OrdersHistory: React.FC = () => {
  const { authState } = useAuth();
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();
  const { toast } = useToast();

  useEffect(() => {
    // Fetch orders from the backend API
    const fetchOrders = async () => {
      if (!authState.isAuthenticated) {
        // console.log("OrdersHistory: User not authenticated, skipping order fetch");
        setIsLoading(false);
        return;
      }

      // console.log("OrdersHistory: Starting to fetch orders for authenticated user");
      try {
        setIsLoading(true);
        // console.log("OrdersHistory: Calling orderService.getOrders()");
        const fetchedOrders = await orderService.getOrders();
        // console.log("OrdersHistory: Received response from orderService.getOrders()");
        // console.log("OrdersHistory: Response type:", typeof fetchedOrders);
        // console.log("OrdersHistory: Is array?", Array.isArray(fetchedOrders));
        // console.log("OrdersHistory: Response value:", fetchedOrders);

        // Ensure fetchedOrders is always an array
        if (Array.isArray(fetchedOrders)) {
          // console.log(`OrdersHistory: Setting ${fetchedOrders.length} orders to state`);
          setOrders(fetchedOrders);
        } else {
          // console.error('OrdersHistory: Expected orders array but got:', fetchedOrders);
          // If it's an object with a data property that's an array, use that instead
          const ordersResponse = fetchedOrders as OrdersResponse;
          if (ordersResponse && typeof ordersResponse === 'object' && Array.isArray(ordersResponse.data)) {
            // console.log(`OrdersHistory: Found orders in data property, setting ${ordersResponse.data.length} orders to state`);
            setOrders(ordersResponse.data);
          } else {
            // Otherwise, set to empty array
            // console.error('OrdersHistory: Could not extract orders from response, setting empty array');
            setOrders([]);
            toast({
              title: 'Error loading orders',
              description: 'The orders data format was invalid',
              variant: 'destructive'
            });
          }
        }
      } catch (error) {
        // console.error('OrdersHistory: Error fetching orders:', error);
        setOrders([]);
        toast({
          title: 'Error fetching orders',
          description: error instanceof Error ? error.message : 'An unknown error occurred',
          variant: 'destructive'
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchOrders();
  }, [authState.isAuthenticated, toast]);

  // Function to handle retrying a payment
  const handleRetryPayment = (orderId: string) => {
    // console.log(`OrdersHistory: Retrying payment for order ${orderId}`);

    // Find the order to get payment details
    const order = orders.find(o => o.id === orderId);
    if (!order) {
      toast({
        title: 'Error',
        description: 'Order not found',
        variant: 'destructive'
      });
      return;
    }

    // Navigate to payment page with order details
    navigate(`/checkout/payment/${orderId}`, {
      state: {
        orderId: orderId,
        amount: order.total,
        paymentMethod: order.paymentMethod
      }
    });
  };

  // Function to safely filter orders
  const filterOrders = (status: string) => {
    // Debug filter operations
    // console.log(`OrdersHistory: Filtering orders for status "${status}"`);
    // console.log("OrdersHistory: Current orders:", orders);
    // console.log("OrdersHistory: Orders is array?", Array.isArray(orders));

    // Ensure orders is an array before filtering
    if (!Array.isArray(orders)) {
      // console.error('OrdersHistory: Orders is not an array:', orders);
      return [];
    }

    // Special case for completed tab - show both delivered and paid orders
    if (status === 'completed') {
      const filtered = orders.filter((o) => o.status === 'delivered' || o.status === 'paid');
      // console.log(`OrdersHistory: Found ${filtered.length} completed orders (delivered/paid)`);
      return filtered;
    }

    // Normal filter by status
    const filtered = orders.filter((o) => o.status === status);
    // console.log(`OrdersHistory: Found ${filtered.length} orders with status "${status}"`);
    return filtered;
  };

  if (isLoading) {
    return (
      <UserLayout>
        <Loading />
      </UserLayout>
    );
  }

  // Debug orders before rendering
  // console.log("OrdersHistory: About to render with orders:", orders);
  // console.log("OrdersHistory: Orders is array?", Array.isArray(orders));
  // console.log("OrdersHistory: Orders length:", orders?.length);

  return (
    <UserLayout>
      <div className="max-w-4xl mx-auto p-6">
        <h1 className="text-[20px] font-medium font-sans mb-5">ORDER HISTORY</h1>

        <Tabs defaultValue="all">
          <TabsList className="mb-6 gap-8">
            <TabsTrigger
              className="min-w-[100px] h-[40px] data-[state=active]:bg-[#0496FF]  data-[state=inactive]:hover:text-gray-900 data-[state=inactive]:hover:bg-blue-300 data-[state=active]:text-white data-[state=inactive]:bg-neutral-200 data-[state=inactive]:text-gray-700 rounded-lg"
              value="all"
            >
              All Orders
            </TabsTrigger>
            <TabsTrigger
              className="min-w-[100px] h-[40px] data-[state=active]:bg-[#0496FF] data-[state=inactive]:hover:text-gray-900 data-[state=inactive]:hover:bg-blue-300 data-[state=active]:text-white data-[state=inactive]:bg-neutral-200 data-[state=inactive]:text-gray-700 rounded-lg"
              value="processing"
            >
              Processing
            </TabsTrigger>
            <TabsTrigger
              className="min-w-[100px] h-[40px]  data-[state=active]:bg-[#0496FF] data-[state=inactive]:hover:text-gray-900 data-[state=inactive]:hover:bg-blue-300 data-[state=active]:text-white data-[state=inactive]:bg-neutral-200 data-[state=inactive]:text-gray-700 rounded-lg"
              value="shipped"
            >
              Shipped
            </TabsTrigger>
            <TabsTrigger
              className="min-w-[100px] h-[40px]  data-[state=active]:text-white data-[state=active]:bg-[#0496FF] data-[state=inactive]:hover:text-gray-900 data-[state=inactive]:hover:bg-blue-300 data-[state=inactive]:bg-neutral-200 data-[state=inactive]:text-gray-700 rounded-lg"
              value="completed"
            >
              Completed
            </TabsTrigger>
            <TabsTrigger
              className="min-w-[100px] h-[40px]  data-[state=active]:text-white data-[state=active]:bg-[#0496FF] data-[state=inactive]:hover:text-gray-900 data-[state=inactive]:hover:bg-blue-300 data-[state=inactive]:bg-neutral-200 data-[state=inactive]:text-gray-700 rounded-lg"
              value="pending"
            >
              Pending
            </TabsTrigger>
          </TabsList>

          <TabsContent value="all" className="space-y-4">
            {orders.length > 0 ? (
              orders.map((order) => (
                <OrderCard
                  key={order.id}
                  order={order}
                  onRetryPayment={handleRetryPayment}
                />
              ))
            ) : (
              <div className="text-center p-8 bg-white rounded-lg shadow-sm">
                <h3 className="text-lg font-medium mb-2">No orders yet</h3>
                <p className="text-gray-500 mb-4">You haven't placed any orders with us yet.</p>
                <Button onClick={() => navigate('/products')}>
                  Start Shopping
                </Button>
              </div>
            )}
          </TabsContent>

          <TabsContent value="processing" className="space-y-4">
            {filterOrders('processing').length > 0 ? (
              filterOrders('processing').map((order) => (
                <OrderCard
                  key={order.id}
                  order={order}
                  onRetryPayment={handleRetryPayment}
                />
              ))
            ) : (
              <div className="text-center p-8 bg-white rounded-lg shadow-sm">
                <p className="text-gray-500">No processing orders</p>
              </div>
            )}
          </TabsContent>

          <TabsContent value="shipped" className="space-y-4">
            {filterOrders('shipped').length > 0 ? (
              filterOrders('shipped').map((order) => (
                <OrderCard
                  key={order.id}
                  order={order}
                  onRetryPayment={handleRetryPayment}
                />
              ))
            ) : (
              <div className="text-center p-8 bg-white rounded-lg shadow-sm">
                <p className="text-gray-500">No shipped orders</p>
              </div>
            )}
          </TabsContent>

          <TabsContent value="completed" className="space-y-4">
            {(filterOrders('completed').length > 0) ? (
              filterOrders('completed').map((order) => (
                <OrderCard
                  key={order.id}
                  order={order}
                  onRetryPayment={handleRetryPayment}
                />
              ))
            ) : (
              <div className="text-center p-8 bg-white rounded-lg shadow-sm">
                <p className="text-gray-500">No completed orders</p>
              </div>
            )}
          </TabsContent>

          <TabsContent value="pending" className="space-y-4">
            {filterOrders('pending').length > 0 ? (
              filterOrders('pending').map((order) => (
                <OrderCard
                  key={order.id}
                  order={order}
                  onRetryPayment={handleRetryPayment}
                />
              ))
            ) : (
              <div className="text-center p-8 bg-white rounded-lg shadow-sm">
                <p className="text-gray-500">No pending orders</p>
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>
    </UserLayout>
  );
};

export default OrdersHistory;