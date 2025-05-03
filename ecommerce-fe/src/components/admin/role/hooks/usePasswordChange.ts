// components/admin/role/hooks/usePasswordChange.ts
import { useState } from 'react';
import { PasswordChangeData } from "../models/adminRole.model"
import { changePassword } from "../services/roleAdmin.service"
import { useToast } from '@/hooks/use-toast';

export const usePasswordChange = () => {
  const { toast } = useToast();
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
    if (!passwordData.currentPassword) {
      setError("Current password is required");
      return false;
    }

    if (!passwordData.newPassword) {
      setError("New password is required");
      return false;
    }

    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setError("New passwords don't match");
      return false;
    }

    if (passwordData.newPassword.length < 8) {
      setError("New password must be at least 8 characters");
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
        toast({
          variant: 'success',
          title: 'Success',
          description: 'Password has been changed successfully'
        });
      }
      return result;
    } catch (err: any) {
      // Extract error message from API response if available
      const errorMessage = err.response?.data?.message ||
        err.response?.data?.error ||
        err.message ||
        'Failed to change password';

      setError(errorMessage);

      toast({
        variant: 'destructive',
        title: 'Password Change Failed',
        description: errorMessage
      });

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