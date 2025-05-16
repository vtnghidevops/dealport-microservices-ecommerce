/**
 * API configuration utility to centralize URL handling 
 * and ensure consistency across services
 */

// Base URL without /api/v1 path
export const BASE_API_URL = import.meta.env.VITE_PUBLIC_BROKER_API_URL || 'http://localhost:58080/api/v1';

// Default API version path
export const API_VERSION = 'v1';

// Admin API prefix for administrative endpoints
export const ADMIN_PREFIX = 'admin';

/**
 * Format image URL - special handling for image paths
 * @param url The image URL to format
 * @returns Formatted image URL
 */
export const formatImageUrl = (url: string): string => {
  // If the URL is empty, return empty string
  if (!url) return '';

  // If URL already has http:// or https://, it's already a complete URL
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url;
  }

  // Get the base domain without /api/v1 if it exists
  let baseImageUrl = BASE_API_URL;
  if (baseImageUrl.endsWith(`/api/${API_VERSION}`)) {
    baseImageUrl = baseImageUrl.substring(0, baseImageUrl.length - `/api/${API_VERSION}`.length);
  }

  // For all image URLs, just add the domain without /api/v1
  if (url.startsWith('/')) {
    return `${baseImageUrl}${url}`;
  }

  // If URL doesn't start with slash, add one
  return `${baseImageUrl}/${url}`;
};

/**
 * Generate a complete API URL with proper path handling
 * @param path The API endpoint path
 * @param isAdmin Whether to use admin prefix in the URL
 * @returns Complete API URL string
 */
export const getApiUrl = (path: string, isAdmin = false): string => {
  // If path starts with /images, keep it as is and just append to base URL
  // This special case is for image URLs that should not have /api/v1 added
  if (path.startsWith('/images')) {
    return `${BASE_API_URL}${path}`;
  }

  // Remove leading slash from path if it exists
  const cleanPath = path.startsWith('/') ? path.substring(1) : path;

  // Check if path already includes an admin prefix to avoid duplication
  const hasAdminPrefix = cleanPath.startsWith(ADMIN_PREFIX);

  // If admin prefix is requested and not already in path, prepend it
  const adminPath = isAdmin && !hasAdminPrefix ? `${ADMIN_PREFIX}/${cleanPath}` : cleanPath;

  // Handle different URL configurations
  // 1. If BASE_API_URL already ends with /api/v1, just append the path
  if (BASE_API_URL.endsWith(`/api/${API_VERSION}`)) {
    return `${BASE_API_URL}/${adminPath}`;
  }

  // 2. Otherwise, add /api/v1 before the path
  return `${BASE_API_URL}/api/${API_VERSION}/${adminPath}`;
};

/**
 * Get authentication headers for API requests
 * @returns Authentication headers object
 */
export const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  if (!token) {
    console.warn('No token found for authenticated request');
    return {};
  }
  return { Authorization: `Bearer ${token}` };
};

/**
 * Get common headers including Content-Type and authorization
 * @param contentType Content type for the request
 * @returns Headers object with Content-Type and authorization
 */
export const getCommonHeaders = (contentType = 'application/json') => {
  return {
    'Content-Type': contentType,
    ...getAuthHeader()
  };
};

/**
 * Get complete admin headers including admin-specific headers
 * @returns Headers object with admin-specific headers
 */
export const getAdminHeaders = () => {
  return {
    ...getCommonHeaders(),
    'X-Admin-Access': 'true'
  };
};

export default {
  BASE_API_URL,
  API_VERSION,
  ADMIN_PREFIX,
  getApiUrl,
  getAuthHeader,
  getCommonHeaders,
  getAdminHeaders,
  formatImageUrl
}; 