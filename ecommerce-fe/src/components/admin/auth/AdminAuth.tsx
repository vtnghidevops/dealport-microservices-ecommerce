import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useToast } from '@/hooks/use-toast';
import { FiEye, FiEyeOff } from "react-icons/fi";
import { extractErrorMessage } from '@/utils/error-handler';

// Simple logo component
const Logo: React.FC = () => {
  return (
    <div className="flex flex-col items-center">
      <Link to="/" className="flex items-center">
        <img src="/images/common/logo.png" alt="logo" className="h-10" />
        <span className="ml-2 font-semibold text-xl text-gray-800">
          EcomSphere
        </span>
      </Link>
    </div>
  );
};

const AdminAuth: React.FC = () => {
  const { authState, login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const { toast } = useToast();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Get the redirect path from location state if it exists
  const from = location.state?.from?.pathname || "/admin/dashboard";

  // Check if already logged in as admin
  useEffect(() => {
    if (authState.isAuthenticated && authState.user?.role === 'admin') {
      // Already logged in as admin, ask if they want to continue to dashboard
      const confirmed = window.confirm('You are already logged in as admin. Go to admin dashboard?');
      if (confirmed) {
        navigate(from);
      }
    } else if (authState.isAuthenticated && authState.user?.role !== 'admin') {
      // Logged in but not as admin - show warning
      setError('Your account does not have admin privileges. Please log in with an admin account.');
    }
  }, [authState.isAuthenticated, authState.user?.role, navigate, from]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const result = await login(email, password);

      if (result.success) {
        // Explicitly check if role is 'admin' (case sensitive)
        if (authState.user?.role === 'admin') {
          toast({
            variant: 'success',
            title: 'Login Successful',
            description: 'Welcome to admin dashboard'
          });
          navigate(from);
        } else {
          setError('This account does not have admin privileges.');
          toast({
            variant: 'destructive',
            title: 'Access Denied',
            description: 'You do not have admin privileges'
          });
          // Log out since they're not an admin
          setTimeout(() => {
            window.location.reload();
          }, 1500);
        }
      } else {
        const friendlyError = extractErrorMessage(result.error);
        setError(friendlyError);
        toast({
          variant: 'destructive',
          title: 'Login Failed',
          description: friendlyError
        });
      }
    } catch (error: any) {
      const errorMessage = extractErrorMessage(error);
      setError(errorMessage);
      toast({
        variant: 'destructive',
        title: 'Login Error',
        description: errorMessage
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen bg-gray-50">
      <div className="flex flex-col justify-center items-center flex-1 px-4 sm:px-6 lg:flex-none lg:px-20 xl:px-24">
        <div className="w-full max-w-sm lg:w-96">
          <div className="text-center mb-8">
            <Logo />
            <h2 className="mt-6 text-3xl font-extrabold text-gray-900">Admin Login</h2>
            <p className="mt-2 text-sm text-gray-600">
              Enter your credentials to access the admin dashboard
            </p>
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-600 rounded-md">
              {error}
            </div>
          )}

          <form className="space-y-6" onSubmit={handleLogin}>
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                Email address
              </label>
              <div className="mt-1">
                <Input
                  id="email"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="h-[44px]"
                  placeholder="admin@example.com"
                />
              </div>
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                Password
              </label>
              <div className="relative mt-1">
                <Input
                  id="password"
                  name="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="current-password"
                  required
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="h-[44px] pr-10"
                />
                <button
                  type="button"
                  className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-700"
                  tabIndex={-1}
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? <FiEyeOff /> : <FiEye />}
                </button>
              </div>
            </div>

            <div>
              <Button
                type="submit"
                className="w-full h-[44px] flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-orange-500 hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500"
                disabled={loading}
              >
                {loading ? (
                  <>
                    <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Logging in...
                  </>
                ) : (
                  'Sign in to Admin'
                )}
              </Button>
            </div>
          </form>

          <div className="mt-6">
            <Button
              variant="ghost"
              className="w-full text-sm text-indigo-600 hover:text-indigo-500"
              onClick={() => navigate('/')}
            >
              Back to Home
            </Button>
          </div>
        </div>
      </div>
      <div className="hidden lg:block relative w-0 flex-1">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-indigo-600 flex justify-center items-center">
          <div className="px-12 text-white">
            <h1 className="text-4xl font-bold mb-6">Admin Dashboard</h1>
            <p className="text-xl">Manage your e-commerce platform with powerful admin tools.</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminAuth; 