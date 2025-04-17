import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const LogoutPage = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    // Execute logout
    logout();

    // Redirect to login page
    navigate('/login');
  }, [logout, navigate]);

  // This component doesn't render anything as it immediately redirects
  return null;
};

export default LogoutPage; 