// components/admin/role/hooks/usePasswordChange.ts
import { useState } from 'react';
import { PasswordChangeData } from "../models/adminRole.model"
import { changePassword } from "../services/roleAdmin.service"

export const usePasswordChange = () => {
  const [passwordData, setPasswordData] = useState<PasswordChangeData>({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  
  const [showCurrentPassword, setShowCurrentPassword] = useState<boolean>(false);
  const [showNewPassword, setShowNewPassword] = useState<boolean>(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const resetPasswordForm = () => {
    setPasswordData({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    });
  };

  const handlePasswordChange = async () => {
    // Validation
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setError("New passwords don't match");
      return false;
    }

    if (passwordData.newPassword.length < 6) {
      setError("New password must be at least 6 characters");
      return false;
    }

    try {
      setIsLoading(true);
      setError(null);
      const result = await changePassword(
        passwordData.currentPassword,
        passwordData.newPassword
      );
      if (result) {
        resetPasswordForm();
      }
      return result;
    } catch (err) {
      setError('Failed to change password');
      console.error('Failed to change password', err);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    passwordData,
    setPasswordData,
    showCurrentPassword,
    setShowCurrentPassword,
    showNewPassword,
    setShowNewPassword,
    showConfirmPassword, 
    setShowConfirmPassword,
    handlePasswordChange,
    resetPasswordForm,
    error,
    isLoading
  };
};