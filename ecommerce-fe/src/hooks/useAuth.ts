// hooks/useAuth.tsx
import { useContext } from 'react';
import { AuthContext } from '@/context/AuthContext';
// import authService from '@/services/auth/auth.service';

/**
 * Hook để truy cập AuthContext một cách dễ dàng
 * @returns AuthContext
 */
export const useAuth = () => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export default useAuth;