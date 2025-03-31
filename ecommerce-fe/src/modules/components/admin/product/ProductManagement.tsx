// components/admin/product/ProductManagement.tsx
import React, { useState } from 'react';
import { AddProduct } from './add';
import AdminHeader from '../layout/AdminHeader';

export const ProductManagement: React.FC = () => {
  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Add Product" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            <AddProduct />
          </div>
        </main>
      </div>
    </div>
  );
};