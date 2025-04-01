// components/admin/role/AdminRoleManagement.tsx
import React from "react";
import {
  ProfileCard,
  ProfileUpdateForm,
  PasswordChangeForm,
} from "./components";
import { useAdminRole } from "./hooks/useAdminRole";
import AdminHeader from "../layout/AdminHeader";

const AdminRoleManagement: React.FC = () => {
  const { adminData, isLoading, error, handleProfileUpdate, setAdminData } =
    useAdminRole();

  if (isLoading) {
    return (
      <div className="p-6 flex justify-center items-center">
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error || !adminData) {
    return (
      <div className="p-6">
        <div className="bg-red-100 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error || "Failed to load admin data"}
        </div>
      </div>
    );
  }

  return (
    <div className="flex bg-neutral-50">
      <div className="flex-1 overflow-auto">
        <AdminHeader title="Admin role" />
        <main className="p-[1rem]">
          <div className="w-[1116px] mx-auto">
            <h2 className="mb-[1.5rem] text-[22px] font-bold text-cyprus">
              About section
            </h2>
            <div className="flex">
              <div className="flex flex-col gap-[1.5rem]">
                <ProfileCard adminData={adminData} />
                {/* Change Password Section */}
                <PasswordChangeForm />
              </div>

              <ProfileUpdateForm
                adminData={adminData}
                onUpdateProfile={handleProfileUpdate}
                setAdminData={setAdminData}
              />
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default AdminRoleManagement;
