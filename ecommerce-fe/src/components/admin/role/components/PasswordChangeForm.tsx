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
    error
  } = usePasswordChange();

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData(prev => ({ ...prev, [name]: value }));
  };

  return (
    <div className="bg-white rounded-lg shadow drop-shadow-sm filter p-[1rem] w-[360px] h-[480px] max-w-xl">
      <div className="flex justify-between items-center mb-12">
        <h3 className="text-[16px] font-bold text-cyprus">Change Password</h3>
        <button className="text-blue-600 hover:text-blue-800 text-sm">
          Need help?
        </button>
      </div>
      
      {error && (
        <div className="bg-red-100 border border-red-200 text-red-700 px-4 py-2 rounded mb-4">
          {error}
        </div>
      )}
      
      <div >
        <div className='mt-[1.5rem]'>
          <label className="block text-[15px] font-medium text-cyprus mb-8">Current Password</label>
          <PasswordInput
            type="password"
            placeholder="Enter password"
            value={passwordData.currentPassword}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: 'currentPassword' }})}
            showPassword={showCurrentPassword}
            toggleShow={() => setShowCurrentPassword(!showCurrentPassword)}
          />
          <div className="mt-1">
            <a href="#" className="text-primary 0 hover:text-blue-800 text-sm">
              Forgot Current Password? Click here
            </a>
          </div>
        </div>
        
        <div className='mt-[1.5rem]'>
          <label className="block text-[15px] font-medium text-cyprus mb-8">New Password</label>
          <PasswordInput
            type="password"
            placeholder="Enter password"
            value={passwordData.newPassword}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: 'newPassword' }})}
            showPassword={showNewPassword}
            toggleShow={() => setShowNewPassword(!showNewPassword)}
          />
        </div>
        
        <div className='mt-[1.5rem]'>
          <label className="block text-[15px] font-medium text-cyprus mb-8">Re-enter Password</label>
          <PasswordInput
            type="password"
            placeholder="Enter password"
            value={passwordData.confirmPassword}
            onChange={(e) => handleInputChange({ ...e, target: { ...e.target, name: 'confirmPassword' }})}
            showPassword={showConfirmPassword}
            toggleShow={() => setShowConfirmPassword(!showConfirmPassword)}
          />
        </div>
        
        <button 
          onClick={handlePasswordChange}
          className="w-full bg-ocean-green hover:bg-green-600 text-white py-[10px] px-[12px] mt-[1.5rem] rounded-md transition-colors font-medium"
        >
          Save Change
        </button>
      </div>
    </div>
  );
};