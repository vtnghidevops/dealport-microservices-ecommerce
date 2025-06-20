import http from "k6/http";
import { check, sleep } from "k6";
import { getBaseUrl, headers, config } from "../config/test-config.js";

/**
 * Login user and return authentication token
 * @param {Object} user - User credentials {email, password}
 * @returns {string} - JWT token
 */
export function loginUser(user) {
  const loginPayload = {
    email: user.email,
    password: user.password,
  };

  const response = http.post(
    `${getBaseUrl()}/api/v1/auth/login`,
    JSON.stringify(loginPayload),
    { headers }
  );

  const loginSuccess = check(response, {
    "login successful": (r) => r.status === 200,
    "login response has token": (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.access_token;
      } catch (e) {
        return false;
      }
    },
  });

  if (loginSuccess && response.status === 200) {
    try {
      const responseBody = JSON.parse(response.body);
      return responseBody.data.access_token;
    } catch (e) {
      console.error("Failed to parse login response:", e);
      return null;
    }
  }

  console.error("Login failed:", response.status, response.body);
  return null;
}

/**
 * Register a new user for testing
 * @param {Object} userData - User registration data
 * @returns {boolean} - Registration success
 */
export function registerUser(userData) {
  const registerPayload = {
    email: userData.email,
    password: userData.password,
    first_name: userData.firstName || "Test",
    last_name: userData.lastName || "User",
    phone: userData.phone || "+1234567890",
  };

  const response = http.post(
    `${getBaseUrl()}/api/v1/auth/register`,
    JSON.stringify(registerPayload),
    { headers }
  );

  return check(response, {
    "registration successful": (r) => r.status === 201 || r.status === 200,
  });
}

/**
 * Get authenticated headers with JWT token
 * @param {string} token - JWT token
 * @returns {Object} - Headers with authorization
 */
export function getAuthHeaders(token) {
  return Object.assign({}, headers, {
    Authorization: `Bearer ${token}`,
  });
}

/**
 * Validate token
 * @param {string} token - JWT token to validate
 * @returns {boolean} - Token validity
 */
export function validateToken(token) {
  const response = http.get(`${getBaseUrl()}/api/v1/auth/validate`, {
    headers: getAuthHeaders(token),
  });

  return check(response, {
    "token is valid": (r) => r.status === 200,
  });
}

/**
 * Setup test user (register if needed, then login)
 * @param {Object} userData - User data
 * @returns {string} - JWT token
 */
export function setupTestUser(userData) {
  // Try to login first
  let token = loginUser(userData);

  if (!token) {
    // If login fails, try to register then login
    console.log("Login failed, attempting to register user...");
    const registered = registerUser(userData);

    if (registered) {
      // Wait a bit for registration to complete
      sleep(1);
      token = loginUser(userData);
    }
  }

  return token;
}
