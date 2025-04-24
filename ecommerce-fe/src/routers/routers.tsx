import React from "react";
import { Route, Routes } from "react-router-dom";
import Home from "../pages/homepage/Home";
import DashboardAdmin from "../pages/admin/dashboard/Dashboard";
import OrderPage from "../pages/admin/order"; // Import Order Management page
import AdminLayout from "../components/admin/layout/AdminLayout";
import CustomerPage from "../pages/admin/customer";
import CategoryPage from "../pages/admin/category";
// import ProductPage from "../pages/admin/product";
import AdminRolePage from "../pages/admin/role";

import MainLayout from "@/components/layouts/MainLayout";
import ProductList from "../pages/products/ProductList";
import ProductDetail from "../pages/products/ProductDetails";
import NotFound from "../pages/system/NotFound";
import ScrollToTop from "@/components/common/ScrollToTop";
import Login from "@/pages/system/Login";
import ForgotPassword from "@/pages/system/ForgotPassword";
import VerifyEmail from "@/pages/system/VerifyEmail";
import Register from "@/pages/system/Register";
import ResetPassword from "@/pages/system/ResetPassword";
import Wishlist from "@/pages/wishlist/Wishlist";
import Cart from "@/pages/cart/Cart";
import Checkout from "@/pages/checkout/Checkout";
import SuccessfulPayment from "../pages/checkout/SuccessfulPayment";

import Profile from "../pages/user/Profile";
import OrdersHistory from "../pages/user/OrdersHistory";
import Addresses from "../pages/user/Addresses";
import SecuritySettings from "../pages/user/SecuritySettings";
import UserDashboard from "../pages/user/Dashboard";
import ProductPage from "../pages/admin/product";
const AppRouters: React.FC = () => {
  return (
    <>
      <ScrollToTop />
      <Routes>
        {/* Public routes with MainLayout */}
        <Route path="/" element={<MainLayout />}>
          {/* Home route */}
          <Route index element={<Home />} />

          {/* Authentication routes */}
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/verify-email" element={<VerifyEmail />} />
          <Route path="/forgot-password" element={<ForgotPassword />} />
          <Route path="/reset-password" element={<ResetPassword />} />

          {/* Product/Category routes - grouped for easier management */}
          <Route path="/products" element={<ProductList />} />
          <Route path="/categories" element={<ProductList />} />
          <Route path="/category/:categorySlug" element={<ProductList />} />
          <Route path="/category/:categorySlug/:productSlug" element={<ProductDetail />} />

          {/* Checkout flow routes */}


          {/* User account routes */}
          <Route path="/user">
            <Route path="dashboard" element={<UserDashboard />} />
            <Route path="wishlist" element={<Wishlist />} />
            <Route path="profile" element={<Profile />} />
            <Route path="order-history" element={<OrdersHistory />} />
            <Route path="addresses" element={<Addresses />} />
            <Route path="security-settings" element={<SecuritySettings />} />
            <Route path="orders" element={<Cart />} />
            <Route path="checkout" element={<Checkout />} />
            <Route path="checkout/success" element={<SuccessfulPayment />} />
          </Route>
        </Route>

        {/* Admin routes */}
        <Route path="/admin" element={<AdminLayout />}>
          <Route index element={<DashboardAdmin />} />
          <Route path="dashboard" element={<DashboardAdmin />} />
          <Route path="orders" element={<OrderPage />} />
          <Route path="customers" element={<CustomerPage />} />
          <Route path="categories" element={<CategoryPage />} />
          {/* <Route path="products" element={<ProductPage />} /> */}
          <Route path="role" element={<AdminRolePage />} />
          <Route path="add-product" element={<ProductPage />} />
        </Route>

        {/* Not found route */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </>
  );
};

export default AppRouters;
