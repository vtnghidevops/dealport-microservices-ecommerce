import React from "react";
import { Route, Routes } from "react-router-dom";
import Home from "../pages/homepage/Home";
import DashboardAdmin from "../pages/admin/dashboard/Dashboard";
import OrderPage from "../pages/admin/order"; // Import Order Management page
import AdminLayout from "../components/admin/layout/AdminLayout";
import CustomerPage from "../pages/admin/customer";
import CategoryPage from "../pages/admin/category";
import ProductPage from "../pages/admin/product";
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

const AppRouters: React.FC = () => {
  return (
    <>
      <ScrollToTop />
      <Routes>
        {/* Public routes with MainLayout */}
        <Route path="/" element={<MainLayout />}>
          <Route path="/" element={<Home />} />
          <Route path="/categories" element={<ProductList />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route path="/verify-email" element={<VerifyEmail />} />
          <Route path="/forgot-password" element={<ForgotPassword />} />
          <Route path="/reset-password" element={<ResetPassword />} />
          <Route path="/wishlist" element={<Wishlist />} />
          <Route path="/cart" element={<Cart />} />



          {/* Route cho các URL theo slug */}
          <Route path="/:categorySlug" element={<ProductList />} />
          <Route
            path="/:categorySlug/:productSlug"
            element={<ProductDetail />}
          />

          {/* Routes cho giỏ hàng và thanh toán */}
          <Route path="/products" element={<ProductList />} />
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
    </>
  );
};

export default AppRouters;
