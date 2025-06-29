import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const LogoutPage = () => {
  const { logout } = useAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  useEffect(() => {
    // Check if the "all" parameter is present to determine if we should logout from all devices
    const logoutAll = searchParams.get('all') === 'true';

    // Execute logout with the appropriate parameter
    logout(logoutAll);

    // Redirect to login page
    navigate('/login');
  }, [logout, navigate, searchParams]);

  // This component doesn't render anything as it immediately redirects
  return null;
};

export default LogoutPage; 