import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';

// Import các components trang
import Home from '../pages/HomePage/Home';
import DashboardAdmin from '../pages/admin/dashboard/Dashboard';
import Sidebar from '../components/admin/layout/Sidebar';

const AppRouters: React.FC = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />

        {/* Route cho trang admin dashboard */}
        <Route path="/admin" element={<DashboardAdmin/>} />

        {/* Redirect nếu không tìm thấy route */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
};

export default AppRouters;
