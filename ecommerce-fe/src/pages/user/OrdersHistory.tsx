// pages/user/Orders.tsx
import React, { useState, useEffect } from 'react';
import { useAuth } from '../../hooks/useAuth';
import UserLayout from '../../components/layouts/UserLayout';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../components/ui/tabs';
import OrderCard from '../../components/user/OrderCard';
import { Button } from '../../components/ui/button';
import { useNavigate } from 'react-router-dom';
import Loading from '@/components/shared/Loading';

// Mock order data
const mockOrders = [
  {
    id: 'ord-001',
    date: '2023-05-15',
    status: 'delivered',
    total: 129.99,
    items: [
      { id: 'prod-1', name: 'Wireless Headphones', price: 79.99, quantity: 1, image: '/images/products/headphones.png' },
      { id: 'prod-2', name: 'Phone Case', price: 24.99, quantity: 2, image: '/images/products/case.png' },
    ]
  },
  {
    id: 'ord-002',
    date: '2023-06-22',
    status: 'processing',
    total: 56.97,
    items: [
      { id: 'prod-3', name: 'Smart Watch Band', price: 18.99, quantity: 3, image: '/images/products/watch-band.png' },
    ]
  }
];

const OrdersHistory: React.FC = () => {
  const { authState } = useAuth();
  const [orders, setOrders] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    // In a real app, fetch orders from API
    const fetchOrders = async () => {
      try {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 800));
        setOrders(mockOrders as any);
      } catch (error) {
        console.error('Error fetching orders:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchOrders();
  }, [authState.user?.id]);

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
              value="delivered"
            >
              Delivered
            </TabsTrigger>
          </TabsList>

          <TabsContent value="all" className="space-y-4">
            {orders.length > 0 ? (
              orders.map((order: any) => (
                <OrderCard key={order.id} order={order} />
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
            {orders.filter((o: any) => o.status === 'processing').map((order: any) => (
              <OrderCard key={order.id} order={order} />
            ))}
          </TabsContent>

          <TabsContent value="shipped" className="space-y-4">
            {orders.filter((o: any) => o.status === 'shipped').map((order: any) => (
              <OrderCard key={order.id} order={order} />
            ))}
          </TabsContent>

          <TabsContent value="delivered" className="space-y-4">
            {orders.filter((o: any) => o.status === 'delivered').map((order: any) => (
              <OrderCard key={order.id} order={order} />
            ))}
          </TabsContent>
        </Tabs>
      </div>
    </UserLayout>
  );
};

export default OrdersHistory;