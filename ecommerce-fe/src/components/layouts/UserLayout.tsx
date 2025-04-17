// components/layouts/UserLayout.tsx
import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { FiUser, FiPackage, FiMapPin, FiHeart, FiRefreshCw, FiCreditCard, FiSettings, FiLogOut } from 'react-icons/fi';

interface UserLayoutProps {
  children: React.ReactNode;
}

const UserLayout: React.FC<UserLayoutProps> = ({ children }) => {
  const { pathname } = useLocation();
  const { authState } = useAuth();

  const navigationItems = [
    { path: '/user/dashboard', label: 'Dashboard', icon: FiUser },
    { path: '/user/order-history', label: 'Order History', icon: FiPackage },
    { path: '/user/track-order', label: 'Track Order', icon: FiMapPin },
    { path: '/user/orders', label: 'Shopping Cart', icon: FiPackage },
    { path: '/user/wishlist', label: 'Wishlist', icon: FiHeart },
    { path: '/user/addresses', label: 'Cards & Address', icon: FiCreditCard },
    { path: '/user/profile', label: 'Setting', icon: FiSettings },
    { path: '/logout', label: 'Log-out', icon: FiLogOut },
  ];

  return (
    <div className="min-h-screen py-5">
      <div className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row gap-8 shadow-sm filter">
          {/* Sidebar navigation */}
          <div className="w-full h-fit md:w-[260px] flex-shrink-0 border border-neutral-50 rounded-lg">
            <div className="bg-white p-4 rounded-lg shadow-sm">
              <nav className="space-y-1">
                {navigationItems.map(item => {
                  const Icon = item.icon;
                  const isActive = pathname === item.path;

                  return (
                    <Link
                      key={item.path}
                      to={item.path}
                      className={`flex items-center px-3 py-3 text-[15px] font-medium ${isActive
                        ? 'bg-orange-500 text-white'
                        : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
                        }`}
                    >
                      <Icon
                        size={20}
                        className={`mr-3 ${isActive ? 'text-white' : 'text-gray-500'}`}
                      />
                      {item.label}
                    </Link>
                  );
                })}
              </nav>
            </div>
          </div>

          {/* Main content */}
          <div className="flex-1">
            {children}
          </div>
        </div>
      </div>
    </div>
  );
};

export default UserLayout;