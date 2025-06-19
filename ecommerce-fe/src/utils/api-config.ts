/**
 * API configuration utility to centralize URL handling
 * and ensure consistency across services
 */

// Base URL without /api/v1 path
export const BASE_API_URL =
  import.meta.env.VITE_PUBLIC_BROKER_API_URL || "http://localhost:58080/api/v1";

// Default API version path
export const API_VERSION = "v1";

// Admin API prefix for administrative endpoints
export const ADMIN_PREFIX = "admin";

/**
 * Format image URL - special handling for image paths
 * @param url The image URL to format
 * @returns Formatted image URL
 */
export const formatImageUrl = (url: string): string => {
  // If the URL is empty, return empty string
  if (!url) return "";
  
  // If URL is a presigned URL (contains X-Amz parameters), return it as is
  if (
    url.includes("X-Amz-") ||
    (url.includes("?") && url.includes("Signature=")) ||
    url.includes("minioapi.deploy.io.vn") ||
    url.includes("minio.deploy.io.vn")
  ) {
    console.log("Detected presigned URL:", url.split("?")[0]);
    return url;
  }

  // If URL has any domain with /images/ path, extract just the /images/ path
  // This works for any domain, not just localhost
  // const anyDomainImagePattern = /https?:\/\/[^/]+(\/images\/.*)/;
  // const domainMatch = url.match(anyDomainImagePattern);
  // if (domainMatch && domainMatch[1]) {
  //   return domainMatch[1];
  // }

  // If URL already has http:// or https://, it's already a complete URL
  // if (url.startsWith("http://") || url.startsWith("https://")) {
  //   return url;
  // }

  return url;
};

/**
 * Generate a complete API URL with proper path handling
 * @param path The API endpoint path
 * @param isAdmin Whether to use admin prefix in the URL
 * @returns Complete API URL string
 */
export const getApiUrl = (path: string, isAdmin = false): string => {
  // Fix common typo in protocol (htpp:// -> http://)
  if (path.startsWith("htpp://")) {
    path = path.replace("htpp://", "http://");
  }
  // If path has any domain with /images/ path, extract just the /images/ path
  const anyDomainImagePattern = /https?:\/\/[^/]+(\/images\/.*)/;
  const domainMatch = path.match(anyDomainImagePattern);
  if (domainMatch && domainMatch[1]) {
    return domainMatch[1];
  }

  // If path already has a protocol, return as is
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }

  // If path starts with /images/, keep it as is without adding any domain
  if (path.startsWith("/images/")) {
    return path;
  }

  // Special handling for API paths that include /images/ but not at the start
  // For example /api/products/images/... should get domain added

  // Remove leading slash from path if it exists
  const cleanPath = path.startsWith("/") ? path.substring(1) : path;

  // Check if path already includes an admin prefix to avoid duplication
  const hasAdminPrefix = cleanPath.startsWith(ADMIN_PREFIX);

  // If admin prefix is requested and not already in path, prepend it
  const adminPath =
    isAdmin && !hasAdminPrefix ? `${ADMIN_PREFIX}/${cleanPath}` : cleanPath;

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
  const token = localStorage.getItem("token");
  if (!token) {
    console.warn("No token found for authenticated request");
    return {};
  }
  return { Authorization: `Bearer ${token}` };
};

/**
 * Get common headers including Content-Type and authorization
 * @param contentType Content type for the request
 * @returns Headers object with Content-Type and authorization
 */
export const getCommonHeaders = (contentType = "application/json") => {
  return {
    "Content-Type": contentType,
    ...getAuthHeader(),
  };
};

/**
 * Get complete admin headers including admin-specific headers
 * @returns Headers object with admin-specific headers
 */
export const getAdminHeaders = () => {
  return {
    ...getCommonHeaders(),
    "X-Admin-Access": "true",
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
  formatImageUrl,
};
