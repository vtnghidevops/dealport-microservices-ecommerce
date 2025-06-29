import React, { useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import AdminLayout from '../layout/AdminLayout';

const ProtectedAdminRoute: React.FC = () => {
  const { authState } = useAuth();
  const location = useLocation();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    // Short delay to allow auth state to be properly loaded/updated
    const timer = setTimeout(() => {
      setIsChecking(false);
    }, 500);

    return () => clearTimeout(timer);
  }, []);

  // Show loading while checking authentication
  if (isChecking) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
          <p className="mt-4 text-lg font-semibold">Checking admin access...</p>
        </div>
      </div>
    );
  }

  // If user is not authenticated, redirect to admin login
  if (!authState.isAuthenticated) {
    return <Navigate to="/admin/login" state={{ from: location }} replace />;
  }

  // If user is authenticated but not an admin, redirect to home
  if (authState.user?.role !== 'admin') {
    return <Navigate to="/" replace />;
  }

  // User is authenticated and has admin role, render admin layout with outlet
  return <AdminLayout />;
};

export default ProtectedAdminRoute; 