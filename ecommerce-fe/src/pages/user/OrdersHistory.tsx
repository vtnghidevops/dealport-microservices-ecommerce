// pages/user/Orders.tsx
import React, { useState, useEffect, useRef } from 'react';
import UserLayout from '../../components/layouts/UserLayout';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../components/ui/tabs';
import OrderCard from '../../components/user/OrderCard';
import { Button } from '../../components/ui/button';
import { useNavigate } from 'react-router-dom';
import Loading from '@/components/shared/Loading';
import { useCheckout } from '@/hooks/useCheckout';
import { CheckoutOrder } from '@/types/checkout.model';
import { Order, OrderStatus, OrderItem, BillingInfo as OrderBillingInfo, ShippingInfo as OrderShippingInfo, PaymentInfo as OrderPaymentInfo } from '@/services/user/order.service';

// Adapter function to convert CheckoutOrder to Order
const adaptCheckoutOrderToOrder = (checkoutOrder: CheckoutOrder): Order => {
  // Convert cart items to order items
  const orderItems: OrderItem[] = (checkoutOrder.items || []).map(item => ({
    id: typeof item.id === 'string' ? item.id : '',
    productId: typeof item.productId === 'number' ? String(item.productId) : '',
    name: item.name || '',
    price: item.price || 0,
    quantity: item.quantity || 1,
    imageUrl: item.imageUrl
  }));

  // Convert billing info - note: some fields might be missing depending on the API format
  const billingInfo: OrderBillingInfo = {
    address: checkoutOrder.billingInfo?.address || '',
    city: checkoutOrder.billingInfo?.city || '',
    country: checkoutOrder.billingInfo?.country || '',
    email: checkoutOrder.billingInfo?.email || '',
    firstName: checkoutOrder.billingInfo?.firstName || '',
    lastName: checkoutOrder.billingInfo?.lastName || '',
    phone: checkoutOrder.billingInfo?.phone || '',
    region: checkoutOrder.billingInfo?.state || '',
    zipCode: checkoutOrder.billingInfo?.zipCode || '',
    companyName: ''  // This field might be missing in some API responses
  };

  // Convert shipping info
  const shippingInfo: OrderShippingInfo = {
    address: checkoutOrder.shippingInfo?.address || '',
    city: checkoutOrder.shippingInfo?.city || '',
    companyName: checkoutOrder.shippingInfo?.companyName || '',
    country: checkoutOrder.shippingInfo?.country || '',
    firstName: checkoutOrder.shippingInfo?.firstName || '',
    lastName: checkoutOrder.shippingInfo?.lastName || '',
    region: checkoutOrder.shippingInfo?.region || '',
    shipToDifferentAddress: checkoutOrder.shippingInfo?.shipToDifferentAddress || false,
    shippingCost: typeof checkoutOrder.shippingInfo?.shippingCost === 'number' ? checkoutOrder.shippingInfo.shippingCost : 0,
    shippingMethod: checkoutOrder.shippingInfo?.shippingMethod || '',
    zipCode: checkoutOrder.shippingInfo?.zipCode || ''
  };

  // Convert payment info
  const paymentInfo: OrderPaymentInfo = {
    amount: checkoutOrder.paymentInfo?.amount || 0,
    currency: checkoutOrder.paymentInfo?.currency || 'USD',
    paymentDate: checkoutOrder.paymentInfo?.paymentDate || new Date().toISOString(),
    paymentMethod: checkoutOrder.paymentInfo?.paymentMethod || '',
    status: checkoutOrder.paymentInfo?.status || 'pending',
    transactionId: checkoutOrder.paymentInfo?.transactionId || ''
  };

  // Convert shipping from string to number if needed
  let shipping: number | string = 0;
  if (checkoutOrder.totals?.shipping !== undefined) {
    if (typeof checkoutOrder.totals.shipping === 'number') {
      shipping = checkoutOrder.totals.shipping;
    } else if (typeof checkoutOrder.totals.shipping === 'string') {
      if (checkoutOrder.totals.shipping.toLowerCase() === 'free') {
        shipping = 'Free';
      } else {
        const numValue = parseFloat(checkoutOrder.totals.shipping);
        shipping = isNaN(numValue) ? 0 : numValue;
      }
    }
  }

  return {
    id: checkoutOrder.id,
    userId: checkoutOrder.userId,
    status: checkoutOrder.status as OrderStatus,
    items: orderItems,
    total: checkoutOrder.totals?.total || 0,
    subtotal: checkoutOrder.totals?.subtotal || 0,
    tax: checkoutOrder.totals?.tax || 0,
    shipping: shipping,
    discount: checkoutOrder.totals?.discount || 0,
    createdAt: checkoutOrder.createdAt,
    updatedAt: checkoutOrder.updatedAt,
    paymentMethod: checkoutOrder.paymentInfo?.paymentMethod || '',
    paymentStatus: checkoutOrder.paymentInfo?.status || 'pending',
    transactionId: checkoutOrder.paymentInfo?.transactionId,
    orderNumber: checkoutOrder.orderNumber,
    billingInfo: billingInfo,
    shippingInfo: shippingInfo,
    paymentInfo: paymentInfo,
    notes: checkoutOrder.notes
  };
};

const OrdersHistory: React.FC = () => {
  const [checkoutOrders, setCheckoutOrders] = useState<CheckoutOrder[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();
  const { listOrders } = useCheckout();

  // Use ref to track if orders have been fetched already
  const hasOrdersBeenFetched = useRef(false);

  // Fetch orders function
  const fetchOrders = async () => {
    try {
      console.log('Fetching orders...');
      setIsLoading(true);
      const result = await listOrders();
      console.log('Orders fetched:', result.orders.length);
      setCheckoutOrders(result.orders);
    } catch (error) {
      console.error('Error fetching orders:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // Fetch orders on component mount ONLY
  useEffect(() => {
    // Only fetch if we haven't fetched before
    if (!hasOrdersBeenFetched.current) {
      console.log('Initiating first order fetch');
      hasOrdersBeenFetched.current = true;
      fetchOrders();
    } else {
      console.log('Orders already fetched, skipping fetch');
      setIsLoading(false);
    }
    // No dependencies to prevent re-fetch
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Function to handle retrying a payment
  const handleRetryPayment = (orderId: string) => {
    // Find the order to get payment details
    const order = checkoutOrders.find(o => o.id === orderId);
    if (!order) {
      return;
    }

    navigate(`/checkout/payment/${orderId}`, {
      state: {
        orderId: orderId,
        amount: order.totals?.total || 0,
        paymentMethod: order.paymentInfo?.paymentMethod
      }
    });
  };

  // Filter orders by status
  const filterOrdersByStatus = (status: string) => {
    // Special case for completed tab - show both delivered and paid orders
    if (status === 'completed') {
      return checkoutOrders.filter(o => o.status === 'delivered' || o.status === 'paid');
    }

    // Normal filter by status
    return checkoutOrders.filter(o => o.status === status);
  };

  if (isLoading) {
    return (
      <UserLayout>
        <Loading />
      </UserLayout>
    );
  }

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
            {checkoutOrders.length > 0 ? (
              checkoutOrders.map((checkoutOrder) => (
                <OrderCard
                  key={checkoutOrder.id}
                  order={adaptCheckoutOrderToOrder(checkoutOrder)}
                  onRetryPayment={handleRetryPayment}
                />
              ))
            ) : (
              <div className="mt-5 text-center p-8 bg-white rounded-lg shadow-sm">
                <h3 className="text-lg font-medium mb-2">No orders yet</h3>
                <p className="text-gray-500 mb-4">You haven't placed any orders with us yet.</p>
                <Button className="mt-4 bg-[#0496FF] text-white p-5 text-sm hover:bg-blue-500" onClick={() => navigate('/category/products')}>
                  Start Shopping
                </Button>
              </div>
            )}
          </TabsContent>

          <TabsContent value="processing" className="space-y-4">
            {filterOrdersByStatus('processing').length > 0 ? (
              filterOrdersByStatus('processing').map((checkoutOrder) => (
                <OrderCard
                  key={checkoutOrder.id}
                  order={adaptCheckoutOrderToOrder(checkoutOrder)}
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
            {filterOrdersByStatus('shipped').length > 0 ? (
              filterOrdersByStatus('shipped').map((checkoutOrder) => (
                <OrderCard
                  key={checkoutOrder.id}
                  order={adaptCheckoutOrderToOrder(checkoutOrder)}
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
            {(filterOrdersByStatus('completed').length > 0) ? (
              filterOrdersByStatus('completed').map((checkoutOrder) => (
                <OrderCard
                  key={checkoutOrder.id}
                  order={adaptCheckoutOrderToOrder(checkoutOrder)}
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
            {filterOrdersByStatus('pending').length > 0 ? (
              filterOrdersByStatus('pending').map((checkoutOrder) => (
                <OrderCard
                  key={checkoutOrder.id}
                  order={adaptCheckoutOrderToOrder(checkoutOrder)}
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