import React, { useEffect } from 'react';
import AdminAuth from '@/components/admin/auth/AdminAuth';

const AdminLogin: React.FC = () => {
  useEffect(() => {
    document.title = 'Admin Login | EcomSphere';
  }, []);

  return <AdminAuth />;
};

export default AdminLogin; 