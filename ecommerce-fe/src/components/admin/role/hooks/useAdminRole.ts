// components/admin/role/hooks/useAdminRole.ts
import { useState, useEffect } from 'react';
import { AdminRole } from '../models/adminRole.model';
import { fetchAdminRole, updateAdminRole } from '../services/roleAdmin.service';

export const useAdminRole = () => {
  const [adminData, setAdminData] = useState<AdminRole | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadAdminData = async () => {
      try {
        setIsLoading(true);
        const data = await fetchAdminRole();
        setAdminData(data);
        setError(null);
      } catch (err) {
        setError('Failed to load admin data');
        console.error('Failed to load admin data', err);
      } finally {
        setIsLoading(false);
      }
    };
    
    loadAdminData();
  }, []);
  
  const handleProfileUpdate = async (updatedData: Partial<AdminRole>) => {
    if (!adminData) return false;
    
    try {
      setIsLoading(true);
      const updated = await updateAdminRole({ ...adminData, ...updatedData });
      setAdminData(updated);
      return true;
    } catch (err) {
      setError('Failed to update profile');
      console.error('Failed to update profile', err);
      return false;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    adminData,
    isLoading,
    error,
    handleProfileUpdate,
    setAdminData
  };
};