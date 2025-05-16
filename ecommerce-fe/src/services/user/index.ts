import axios from 'axios';

// Create axios instance with default configuration
const api = axios.create({
  baseURL: import.meta.env.VITE_PUBLIC_BROKER_API_URL || 'https://testapi.deploy.io.vn/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to attach auth token
api.interceptors.request.use(
  (config) => {
    // Get the token from localStorage if it exists
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add a response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Handle unauthorized errors (401)
    if (error.response && error.response.status === 401) {
      // Clear local storage and redirect to login
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api; 