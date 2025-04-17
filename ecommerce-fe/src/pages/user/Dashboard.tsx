import React, { useEffect, useState } from 'react';
import { FiBell, FiEye, FiEdit } from 'react-icons/fi';
import UserLayout from '@/components/layouts/UserLayout';
import { Link } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

// Types for user dashboard data
interface UserData {
  id: string;
  name: string;
  email: string;
  phone: string;
  location: string;
}

interface BillingAddress {
  address: string;
  city: string;
  zipCode: string;
  country: string;
}

interface OrderStats {
  total: number;
  pending: number;
  completed: number;
}

interface PaymentCard {
  id: string;
  last4: string;
  cardHolder: string;
  type: 'visa' | 'mastercard';
  balance?: number;
}

interface Order {
  id: string;
  status: 'IN PROGRESS' | 'COMPLETED' | 'CANCELED';
  date: string;
  total: number;
  productCount: number;
}

const UserDashboard: React.FC = () => {
  const { authState } = useAuth();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [userData, setUserData] = useState<UserData | null>(null);
  const [billingAddress, setBillingAddress] = useState<BillingAddress | null>(null);
  const [orderStats, setOrderStats] = useState<OrderStats>({ total: 0, pending: 0, completed: 0 });
  const [paymentCards, setPaymentCards] = useState<PaymentCard[]>([]);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      setIsLoading(true);
      try {
        // For demo purposes, we'll simulate API calls with setTimeout
        // In a real app, you would fetch this data from your backend API

        // Simulate API response delay
        await new Promise(resolve => setTimeout(resolve, 500));

        // Set demo data - in a real app, these would come from API responses
        setUserData({
          id: '1',
          name: authState?.user?.email?.split('@')[0] || 'Kevin Gilbert',
          email: authState?.user?.email || 'kevin.gilbert@gmail.com',
          phone: '+1-202-555-0118',
          location: 'Dhaka-1207, Bangladesh'
        });

        setBillingAddress({
          address: 'East Tejturi Bazar, Word No. 04, Road No. 13/x, House no. 1320/C, Flat No. 5D',
          city: 'Dhaka',
          zipCode: '1200',
          country: 'Bangladesh'
        });

        setOrderStats({
          total: 154,
          pending: 5,
          completed: 149
        });

        setPaymentCards([
          {
            id: '1',
            last4: '3814',
            cardHolder: userData?.name || 'Kevin Gilbert',
            type: 'visa',
            balance: 95400.00
          },
          {
            id: '2',
            last4: '1761',
            cardHolder: userData?.name || 'Kevin Gilbert',
            type: 'mastercard',
            balance: 87583.00
          }
        ]);

        setRecentOrders([
          {
            id: '96459781',
            status: 'IN PROGRESS',
            date: 'Dec 30, 2019 05:18',
            total: 1500,
            productCount: 5
          },
          {
            id: '71667167',
            status: 'COMPLETED',
            date: 'Feb 2, 2019 19:28',
            total: 80,
            productCount: 1
          },
          {
            id: '95214362',
            status: 'CANCELED',
            date: 'Mar 20, 2019 23:14',
            total: 160,
            productCount: 3
          },
          {
            id: '51746385',
            status: 'COMPLETED',
            date: 'Feb 2, 2019 19:28',
            total: 2300,
            productCount: 2
          }
        ]);
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboardData();
  }, [authState]);

  if (isLoading) {
    return (
      <UserLayout>
        <div className="bg-white p-6 rounded-lg shadow-sm flex justify-center items-center h-96">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-orange-500"></div>
        </div>
      </UserLayout>
    );
  }

  return (
    <UserLayout>
      <div className="bg-white p-5 rounded-lg shadow-sm">
        {/* Welcome Section */}
        <div className="mb-6">
          <h1 className="text-2xl font-semibold mb-2">Hello, {userData?.name || 'User'}</h1>
          <p className="text-gray-600 mb-2">
            From your account dashboard, you can easily check & view your{' '}
            <Link to="/user/order-history" className="text-blue-500 hover:text-blue-700">
              Recent Orders
            </Link>
            , manage your{' '}
            <Link to="/user/addresses" className="text-blue-500 hover:text-blue-700">
              Shipping and Billing Addresses
            </Link>{' '}
            and edit your{' '}
            <Link to="/user/profile" className="text-blue-500 hover:text-blue-700">
              Password
            </Link>{' '}
            and{' '}
            <Link to="/user/profile" className="text-blue-500 hover:text-blue-700">
              Account Details
            </Link>
            .
          </p>
        </div>

        {/* Account Info Section */}
        <div className="mb-8">
          <div className="flex flex-col md:flex-row gap-[15px]">
            <div className="!min-w-1/3 w-1/3 border border-neutral-200 flex flex-col justify-start px-5 bg-white rounded-lg py-5 space-x-4">
              <span className="text-[16px] font-medium font-sans mb-5">ACCOUNT INFO</span>
              <div className="flex gap-16 items-center">
                <img
                  src="/images/system/default-avatars.png"
                  alt="User Avatar"
                  className="w-[40px] h-[40px] rounded-full object-cover"
                  onError={(e) => {
                    const target = e.target as HTMLImageElement;
                    target.src = `https://ui-avatars.com/api/?name=${encodeURIComponent(userData?.name || 'User')}&background=random`;
                  }}
                />
                <div>
                  <h3 className="font-bold text-[16px]">{userData?.name}</h3>
                  <p className="text-gray-500 text-sm">{userData?.location}</p>
                </div>


              </div>
              <div className="flex flex-col gap-4 mt-5">
                <span className="text-sm font-bold">Email: <span className="text-gray-600">{userData?.email}</span></span>
                <span className="text-sm font-bold">Phone: <span className="text-gray-600">{userData?.phone}</span></span>
              </div>

              <div className="flex justify-start mt-5">
                <button className="text-[#0496FF] px-5 hover:bg-blue-50 text-sm font-medium border border-[#0496FF] rounded-md py-2">
                  EDIT ACCOUNT
                </button>
              </div>

            </div>
            <div className="bg-white w-1/3 min-w-1/3 rounded-lg p-5 border border-neutral-200 flex-1">
              <div>
                <h3 className="font-medium font-sans text-[16px] mb-3">BILLING ADDRESS</h3>
                <span className="font-bold text-[16px] mb-3 block">{userData?.name}</span>
                <p className="text-gray-600 text-sm font-medium">
                  {billingAddress?.address},
                  <br />{billingAddress?.city}-{billingAddress?.zipCode}, {billingAddress?.country}
                </p>
              </div>
              <div className="flex flex-col gap-4 mt-3">
                <span className="text-sm font-bold">Email: <span className="text-gray-600">{userData?.email}</span></span>
                <span className="text-sm font-bold">Phone: <span className="text-gray-600">{userData?.phone}</span></span>
              </div>
              <div className="flex justify-start mt-5">
                <button className="text-[#0496FF] px-5 hover:bg-blue-50 text-sm font-medium border border-[#0496FF] rounded-md py-2">
                  EDIT ACCOUNT
                </button>
              </div>
            </div>

            {/* Order Statistics */}
            <div className="mb-8 w-1/3">
              <div className="flex flex-col gap-5">
                <div className="bg-blue-50 rounded-lg flex items-center justify-start h-[90px] p-5 gap-5">
                  <div className="bg-white rounded-full p-3 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-[40px] w-[40px] text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
                    </svg>
                  </div>
                  <div className='w-1/2'>
                    <p className="text-2xl font-bold">{orderStats.total}</p>
                    <h3 className="text-sm text-neutral-600 font-medium">Total Orders</h3>
                  </div>
                </div>

                <div className="bg-orange-100 rounded-lg flex items-center justify-start h-[90px] p-5 gap-5">
                  <div className="bg-white rounded-full p-3 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-[40px] w-[40px] text-pending" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>

                  <div className='w-1/2'>
                    <p className="text-2xl font-bold">{orderStats.pending.toString().padStart(2, '0')}</p>
                    <h3 className="text-sm text-neutral-600 font-medium">Pending Orders</h3>
                  </div>

                </div>
                <div className="bg-aqua-spring rounded-lg flex items-center justify-start h-[90px] p-5 gap-5">
                  <div className="bg-white rounded-full p-3 shadow-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-[40px] w-[40px] text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                  <div>
                    <p className="text-2xl font-bold">{orderStats.completed}</p>
                    <h3 className="text-sm text-neutral-600 font-medium">Completed Orders</h3>
                  </div>

                </div>
              </div>
            </div>

          </div>

        </div>

        {/* Recent Orders */}
        <div className='mt-5 px-5 border border-neutral-200 rounded-lg'>
          <div className="flex justify-between items-center mb-4 mt-5">
            <h2 className="text-lg font-sans font-medium mb-4">RECENT ORDER</h2>
            <Link to="/user/order-history" className="text-sm text-orange-500 font-medium font-sans flex items-center">
              View All
            </Link>
          </div>
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th scope="col" className="font-medium font-sans pl-5 px-6 py-3 text-left text-xs text-gray-500 uppercase tracking-wider">
                    ORDER ID
                  </th>
                  <th scope="col" className="font-medium font-sans px-6 py-3 text-left text-xs text-gray-500 uppercase tracking-wider">
                    STATUS
                  </th>
                  <th scope="col" className="font-medium font-sans px-6 py-3 text-left text-xs text-gray-500 uppercase tracking-wider">
                    DATE
                  </th>
                  <th scope="col" className="font-medium font-sans px-6 py-3 text-left text-xs text-gray-500 uppercase tracking-wider">
                    TOTAL
                  </th>
                  <th scope="col" className="font-medium font-sans px-6 py-3 text-left text-xs text-gray-500 uppercase tracking-wider">
                    ACTION
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {recentOrders.map(order => (
                  <tr key={order.id} className='font-sans h-[50px]'>
                    <td className="pl-5 px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      #{order.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full 
                        ${order.status === 'COMPLETED' ? 'bg-green-100 text-success' :
                          order.status === 'IN PROGRESS' ? 'bg-yellow-100 text-yellow-800' :
                            'bg-red-100 text-red-800'}`}>
                        {order.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {order.date}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      ${order.total.toLocaleString()} ({order.productCount} {order.productCount === 1 ? 'Product' : 'Products'})
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <Link to={`/user/order-history/${order.id}`} className="text-blue-600 hover:text-blue-900 flex items-center">
                        <FiEye className="mr-1" /> View Details
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </UserLayout>
  );
};

export default UserDashboard; 