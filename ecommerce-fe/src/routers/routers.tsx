import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from '../pages/homepage/Home';
import DashboardAdmin from '../pages/admin/dashboard/Dashboard';
import OrderPage from '../pages/admin/order'; // Import Order Management page
import AdminLayout from '../components/admin/layout/AdminLayout';
import CustomerPage from '../pages/admin/customer';
import CategoryPage from '../pages/admin/category';
import ProductPage from '../pages/admin/product';
import AdminRolePage from '../pages/admin/role';


import MainLayout from '@/components/layouts/MainLayout';
import ProductList from '../pages/products/ProductList' 
import ProductDetail from '../pages/products/ProductDetails'; 
import NotFound from '../pages/system/NotFound';

const AppRouters: React.FC = () => {
  return (
    <Router>
      <Routes >
        {/* Public routes with MainLayout */}
        <Route path="/" element={<MainLayout />}>
        <Route path="/" element={<Home />} />
          <Route path="/categories" element={<ProductList />} />
          
          {/* Route cho các URL theo slug */}
          <Route path="/:categorySlug" element={<ProductList />} />
          <Route path="/:categorySlug/:productSlug" element={<ProductDetail />} />
          
          {/* Routes cho giỏ hàng và thanh toán */}
          <Route path="/cart" element={<ProductList />} />
          <Route path="/checkout" element={<ProductList />} />
        </Route>

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
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Router>
  );
};

export default AppRouters;