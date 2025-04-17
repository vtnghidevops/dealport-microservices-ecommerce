// context/AuthContext.tsx
import React, { createContext, useState, ReactNode } from 'react';
import { User, AuthState } from '@/types/user.model';

// Initial state
const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: false,
  error: null
};

// Create the context
export const AuthContext = createContext<{
  authState: AuthState;
  login: (email: string, password: string) => Promise<void>;
  register: (userData: any) => Promise<void>;
  logout: () => void;
  updateProfile: (userData: Partial<User>) => Promise<void>;
}>({
  authState: initialState,
  login: async () => {},
  register: async () => {},
  logout: () => {},
  updateProfile: async () => {}
});

// Create provider
export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [authState, setAuthState] = useState<AuthState>(() => {
    // Check local storage for existing session
    const token = localStorage.getItem('token');
    const user = localStorage.getItem('user');
    
    if (token && user) {
      return {
        user: JSON.parse(user),
        token,
        isAuthenticated: true,
        isLoading: false,
        error: null
      };
    }
    
    return initialState;
  });

  // Login function
  const login = async (email: string, password: string) => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // API call would go here
      // Example:
      // const response = await fetch('/api/auth/login', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ email, password })
      // });
      // const data = await response.json();
      
      // For now, mock the response
      const mockUser = { id: '1', email, firstName: 'John', lastName: 'Doe' };
      const mockToken = 'mock-token-xxx';
      
      // Save to localStorage
      localStorage.setItem('token', mockToken);
      localStorage.setItem('user', JSON.stringify(mockUser));
      setAuthState({
        user: {
          ...mockUser,
          profile: null,
          role: 'USER',
          status: 'ACTIVE',
          addresses: [],
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        } as unknown as User,
        token: mockToken,
        isAuthenticated: true,
        isLoading: false,
        error: null
      });
    } catch (error) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: 'Authentication failed. Please check your credentials.'
      }));
    }
  };

  // Register function
  const register = async (userData: any) => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // API call would go here
      
      // Mock registration response
      const mockUser = { 
        id: '1', 
        email: userData.email, 
        firstName: userData.firstName, 
        lastName: userData.lastName 
      };
      const mockToken = 'mock-token-xxx';
      
      localStorage.setItem('token', mockToken);
      localStorage.setItem('user', JSON.stringify(mockUser));
      setAuthState({
        user: {
          ...mockUser,
          profile: null,
          role: 'USER',
          status: 'ACTIVE', 
          addresses: [],
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString()
        } as unknown as User,
        token: mockToken,
        isAuthenticated: true,
        isLoading: false,
        error: null
      });
    } catch (error) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: 'Registration failed. Please try again.'
      }));
    }
  };

  // Logout function
  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    
    setAuthState(initialState);
  };

  // Update profile function
  const updateProfile = async (userData: Partial<User>) => {
    try {
      setAuthState(prev => ({ ...prev, isLoading: true, error: null }));
      
      // API call would go here
      
      // Update local storage
      const updatedUser = { ...authState.user, ...userData };
      localStorage.setItem('user', JSON.stringify(updatedUser));
      
      setAuthState(prev => ({
        ...prev,
        user: updatedUser as User, // Cast to User type to ensure type safety (đảm bảo - ensure)
        isLoading: false
      }));
    } catch (error) {
      setAuthState(prev => ({
        ...prev,
        isLoading: false,
        error: 'Profile update failed.'
      }));
    }
  };

  return (
    <AuthContext.Provider value={{ authState, login, register, logout, updateProfile }}>
      {children}
    </AuthContext.Provider>
  );
};