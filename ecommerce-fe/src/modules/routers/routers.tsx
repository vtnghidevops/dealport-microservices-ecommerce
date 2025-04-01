import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';

// Import các components trang
import Home from '../pages/HomePage/Home';
import DashboardAdmin from '../pages/admin/dashboard/Dashboard';
import OrderPage from '../pages/admin/order'; // Import Order Management page
import AdminLayout from '../components/admin/layout/AdminLayout';
import CustomerPage from '../pages/admin/customer';
import CategoryPage from '../pages/admin/category';
import ProductPage from '../pages/admin/product';
import AdminRolePage from '../pages/admin/role';
const AppRouters: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />

        <Route path="/admin" element={<AdminLayout />}>
          <Route index element={<DashboardAdmin />} />
          <Route path="dashboard" element={<DashboardAdmin />} />
          <Route path="orders" element={<OrderPage />} />
          <Route path="customers" element={<CustomerPage />} />
          <Route path="categories" element={<CategoryPage />} />
          <Route path="products" element={<ProductPage />} />
          <Route path="role" element={<AdminRolePage />} />
        </Route>

        {/* Redirect if route not found */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
};

export default AppRouters;