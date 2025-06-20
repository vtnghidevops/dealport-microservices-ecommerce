/**
 * Authentication utilities for K6 performance tests
 * With OTP bypass enabled in staging, authentication is now straightforward
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { config } from "../config/test-config.js";

// Simple OTP bypass - only for staging
const STAGING_OTP_BYPASS = "123456";

/**
 * Check if current environment supports OTP bypass
 * @returns {boolean} True if bypass is available
 */
function isOTPBypassAvailable() {
  return config.environment.toLowerCase() === "staging";
}

/**
 * Get OTP code to use for verification
 * @returns {string} OTP bypass code for staging, or null for other environments
 */
function getOTPCode() {
  return isOTPBypassAvailable() ? STAGING_OTP_BYPASS : null;
}

/**
 * Register and verify user with OTP bypass
 * @param {Object} userData - User registration data
 * @returns {string|null} User token or null
 */
export function setupTestUser(userData = {}) {
  const baseUrl = config.baseUrls[config.environment];

  // Check if OTP bypass is available
  if (!isOTPBypassAvailable()) {
    console.log(
      `⚠️ OTP bypass not available in ${config.environment} environment`
    );
    console.log(
      "   Authentication will use real OTP workflow (not suitable for performance testing)"
    );
    return null;
  }

  // Generate unique user data
  const timestamp = Date.now();
  const randomId = Math.random().toString(36).substr(2, 9);

  const user = {
    email: userData.email || `perf-user-${timestamp}-${randomId}@test.com`,
    password: userData.password || "PerfTest123!",
    first_name: userData.firstName || "Performance",
    last_name: userData.lastName || "User",
    phone: userData.phone || "+1234567890",
  };

  try {
    console.log(`Setting up user: ${user.email} (${config.environment} env)`);

    // Step 1: Register user
    const registerResponse = http.post(
      `${baseUrl}/api/v1/auth/register`,
      JSON.stringify(user),
      {
        headers: { "Content-Type": "application/json" },
        tags: { scenario: "auth_setup", type: "register" },
      }
    );

    const registerSuccess = check(registerResponse, {
      "registration successful": (r) => r.status === 200 || r.status === 201,
    });

    if (!registerSuccess) {
      console.log(
        `Registration failed for ${user.email}: ${registerResponse.status}`
      );
      return null;
    }

    // Step 2: Verify with bypass OTP (staging only)
    const otpCode = getOTPCode();
    const verifyResponse = http.post(
      `${baseUrl}/api/v1/auth/verify-registration`,
      JSON.stringify({
        email: user.email,
        otp: otpCode,
      }),
      {
        headers: { "Content-Type": "application/json" },
        tags: { scenario: "auth_setup", type: "verify" },
      }
    );

    const verifySuccess = check(verifyResponse, {
      "verification successful": (r) => r.status === 200,
      "verification has token": (r) => {
        try {
          const data = JSON.parse(r.body);
          return data.data && data.data.access_token;
        } catch {
          return false;
        }
      },
    });

    if (verifySuccess && verifyResponse.status === 200) {
      const verifyData = JSON.parse(verifyResponse.body);
      console.log(`✅ User ${user.email} setup successful with OTP bypass`);
      return verifyData.data.access_token;
    } else {
      console.log(
        `Verification failed for ${user.email}: ${verifyResponse.status}`
      );
      return null;
    }
  } catch (error) {
    console.error(`Setup failed for ${user.email}:`, error.message);
    return null;
  }
}

/**
 * Setup multiple test users efficiently
 * @param {number} maxUsers - Maximum number of users to setup
 * @returns {Array} Array of user tokens
 */
export function setupMultipleTestUsers(maxUsers = 10) {
  console.log(`🔑 Setting up ${maxUsers} test users with OTP bypass...`);

  const userTokens = [];

  for (let i = 0; i < maxUsers; i++) {
    const token = setupTestUser({
      email: `perf-user-${i}-${Date.now()}@test.com`,
      firstName: `Perf${i}`,
      lastName: "User",
    });

    if (token) {
      userTokens.push(token);
      console.log(`   ✅ User ${i + 1}/${maxUsers} authenticated`);
    } else {
      console.log(`   ❌ User ${i + 1}/${maxUsers} failed`);
    }

    // Small delay between registrations
    sleep(0.2);
  }

  console.log(
    `🎯 Setup complete: ${userTokens.length}/${maxUsers} users ready`
  );
  return userTokens;
}

/**
 * Get authentication headers
 * @param {string} token - JWT token
 * @returns {Object} Headers object
 */
export function getAuthHeaders(token) {
  const baseHeaders = {
    "Content-Type": "application/json",
    Accept: "application/json",
  };

  if (token) {
    return Object.assign({}, baseHeaders, {
      Authorization: `Bearer ${token}`,
    });
  }

  return baseHeaders;
}

/**
 * Validate authentication token
 * @param {string} token - JWT token to validate
 * @returns {boolean} True if token is valid
 */
export function validateToken(token) {
  if (!token) return false;

  try {
    const response = http.get(
      `${config.baseUrls[config.environment]}/api/v1/user/profile`,
      {
        headers: getAuthHeaders(token),
        tags: { scenario: "token_validation" },
      }
    );

    return response.status === 200;
  } catch (error) {
    console.error("Token validation error:", error.message);
    return false;
  }
}

/**
 * Get a random user session from available sessions
 * @param {Array} userSessions - Array of user session objects/tokens
 * @returns {Object} User session object with token and metadata
 */
export function getRandomUserSession(userSessions = []) {
  if (!userSessions || userSessions.length === 0) {
    // Return guest session if no authenticated users available
    return {
      token: null,
      isGuest: true,
      sessionId: `guest-${Math.random().toString(36).substr(2, 9)}`,
    };
  }

  // Get random user session from available authenticated sessions
  const randomIndex = Math.floor(Math.random() * userSessions.length);
  const selectedToken = userSessions[randomIndex];

  // Handle both token strings and user session objects
  if (typeof selectedToken === "string") {
    return {
      token: selectedToken,
      isGuest: false,
      sessionId: `user-${randomIndex}-${Math.random()
        .toString(36)
        .substr(2, 6)}`,
    };
  }

  // If it's already a session object, return it with additional metadata
  return {
    ...selectedToken,
    sessionId:
      selectedToken.sessionId ||
      `user-${randomIndex}-${Math.random().toString(36).substr(2, 6)}`,
    isGuest: !selectedToken.token,
  };
}
