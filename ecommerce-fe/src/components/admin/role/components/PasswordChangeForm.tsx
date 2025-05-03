// components/admin/role/components/PasswordChangeForm.tsx
import React from 'react';
import { PasswordInput } from './PasswordInput';
import { usePasswordChange } from '../hooks/usePasswordChange'

export const PasswordChangeForm: React.FC = () => {
  const {
    passwordData,
    setPasswordData,
    showCurrentPassword,
    setShowCurrentPassword,
    showNewPassword,
    setShowNewPassword,
    showConfirmPassword,
    setShowConfirmPassword,
    handlePasswordChange,
    error,
    isLoading
  } = usePasswordChange();

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData(prev => ({ ...prev, [name]: value }));
  };

  return (
    <div className="bg-white rounded-lg shadow drop-shadow-sm filter p-[1rem] w-[360px] h-[480px] max-w-xl">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-[16px] font-bold text-cyprus">Change Password</h3>
        <button className="text-blue-600 hover:text-blue-800 text-sm">
          Need help?
        </button>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-200 text-red-700 px-4 py-2 rounded mb-4 text-sm">
          {error}
        </div>
      )}

      <div className="mt-4">
        <div className='mb-4'>
          <label className="block text-[15px] font-medium text-cyprus mb-2">Current Password</label>
          <PasswordInput
            type="password"
            placeholder="Enter current password"
            value={passwordData.currentPassword}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: 'currentPassword' } })}
            showPassword={showCurrentPassword}
            toggleShow={() => setShowCurrentPassword(!showCurrentPassword)}
            disabled={isLoading}
          />
          <div className="mt-1">
            <a href="/forgot-password" className="text-primary hover:text-blue-800 text-sm">
              Forgot Current Password?
            </a>
          </div>
        </div>

        <div className='mb-4'>
          <label className="block text-[15px] font-medium text-cyprus mb-2">New Password</label>
          <PasswordInput
            type="password"
            placeholder="Enter new password"
            value={passwordData.newPassword}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: 'newPassword' } })}
            showPassword={showNewPassword}
            toggleShow={() => setShowNewPassword(!showNewPassword)}
            disabled={isLoading}
          />
          <div className="mt-1 text-xs text-gray-500">
            Password must be at least 8 characters
          </div>
        </div>

        <div className='mb-5'>
          <label className="block text-[15px] font-medium text-cyprus mb-2">Re-enter Password</label>
          <PasswordInput
            type="password"
            placeholder="Confirm new password"
            value={passwordData.confirmPassword}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: 'confirmPassword' } })}
            showPassword={showConfirmPassword}
            toggleShow={() => setShowConfirmPassword(!showConfirmPassword)}
            disabled={isLoading}
          />
        </div>

        <button
          onClick={handlePasswordChange}
          disabled={isLoading}
          className={`w-full bg-ocean-green hover:bg-green-600 text-white py-[10px] px-[12px] mt-4 rounded-md transition-colors font-medium flex justify-center items-center
            ${isLoading ? 'opacity-70 cursor-not-allowed' : ''}`}
        >
          {isLoading ? (
            <>
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Saving...
            </>
          ) : (
            'Save Change'
          )}
        </button>
      </div>
    </div>
  );
};