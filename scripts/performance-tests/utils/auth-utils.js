/**
 * Authentication utilities for K6 performance tests
 * With OTP bypass enabled in staging, authentication is now straightforward
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { config } from "../config/test-config.js";

/**
 * Check if OTP bypass is available in current environment
 * @returns {boolean} True if bypass is available
 */
function isOTPBypassAvailable() {
  // OTP bypass is only available in staging environment
  return config.environment === "staging" || __ENV.ENVIRONMENT === "staging";
}

/**
 * Get OTP code for verification
 * @returns {string} OTP code (bypass code for staging, random for others)
 */
function getOTPCode() {
  if (isOTPBypassAvailable()) {
    return "123456"; // Bypass code for staging
  }
  // Generate random 6-digit code for other environments (will likely fail)
  return Math.floor(100000 + Math.random() * 900000).toString();
}

/**
 * Register and verify user with OTP bypass
 * @param {Object} userData - User registration data
 * @returns {string|null} User token or null
 */
export function setupTestUser(userData = {}, options = {}) {
  const { setupMode = false, shortDelay = 0.3 } = options;

  const baseUrl = config.baseUrls[config.environment];

  if (!baseUrl) {
    console.error(
      `No base URL configured for environment: ${config.environment}`
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
    if (!setupMode) {
      console.log(`Setting up user: ${user.email} (${config.environment} env)`);
    }

    // Step 1: Register user
    const registerResponse = http.post(
      `${baseUrl}/api/v1/auth/register`,
      JSON.stringify(user),
      {
        headers: { "Content-Type": "application/json" },
        tags: { scenario: "auth_setup", type: "register" },
        timeout: "30s", // Keep timeout for reliability
      }
    );

    const registerSuccess = check(registerResponse, {
      "registration successful": (r) => r.status === 200 || r.status === 201,
    });

    if (!registerSuccess) {
      if (!setupMode) {
        console.log(
          `Registration failed for ${user.email}: ${
            registerResponse.status
          } - ${
            registerResponse.body
              ? registerResponse.body.substring(0, 100)
              : "No response body"
          }`
        );
      }
      return null;
    }

    // Essential delay for OTP bypass to work in staging
    sleep(0.2); // MINIMAL delay just for OTP system

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
        timeout: "30s", // Keep timeout
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

      if (!setupMode) {
        console.log(`✅ User ${user.email} setup successful with OTP bypass`);
        // Add delay for microservices sync only during runtime
        console.log(`⏳ Waiting 3s for user sync across microservices...`);
        sleep(3.0); // Keep this for runtime reliability
      } else {
        // During setup, minimal delay
        sleep(shortDelay);
      }

      return verifyData.data.access_token;
    } else {
      if (!setupMode) {
        console.log(
          `Verification failed for ${user.email}: ${verifyResponse.status} - ${
            verifyResponse.body
              ? verifyResponse.body.substring(0, 100)
              : "No response body"
          }`
        );
      }
      return null;
    }
  } catch (error) {
    if (!setupMode) {
      console.error(`Setup failed for ${user.email}:`, error.message);
    }
    return null;
  }
}

/**
 * Setup multiple test users efficiently with parallel processing
 * @param {number} maxUsers - Maximum number of users to setup
 * @param {Object} options - Setup options
 * @returns {Array} Array of user tokens
 */
export function setupMultipleTestUsers(maxUsers = 10, options = {}) {
  const { batchSize = 25, setupDelay = 0.3 } = options; // OPTIMIZED for high load

  console.log(
    `🔑 Setting up ${maxUsers} test users with optimized parallel processing...`
  );
  console.log(
    `📦 Using batches of ${batchSize} users with ${setupDelay}s delay`
  );

  const userTokens = [];
  const totalBatches = Math.ceil(maxUsers / batchSize);

  for (let batch = 0; batch < totalBatches; batch++) {
    const batchStart = batch * batchSize;
    const batchEnd = Math.min(batchStart + batchSize, maxUsers);
    const batchSize_actual = batchEnd - batchStart;

    console.log(
      `🔄 Processing batch ${
        batch + 1
      }/${totalBatches} (${batchSize_actual} users)...`
    );

    const batchTokens = [];
    const batchStartTime = Date.now();

    // Process batch efficiently
    for (let i = batchStart; i < batchEnd; i++) {
      const token = setupTestUser(
        {
          email: `perf-user-${i}-${Date.now()}@test.com`,
          firstName: `Perf${i}`,
          lastName: "User",
        },
        {
          setupMode: true, // Skip the 3s runtime delay
          shortDelay: setupDelay,
        }
      );

      if (token) {
        batchTokens.push(token);
      }

      // Minimal delay between users - just enough for staging OTP bypass
      sleep(0.1); // OPTIMIZED - very short delay
    }

    userTokens.push(...batchTokens);

    const batchTime = (Date.now() - batchStartTime) / 1000;
    console.log(
      `   ✅ Batch ${batch + 1} completed: ${
        batchTokens.length
      }/${batchSize_actual} users in ${batchTime.toFixed(1)}s`
    );

    // Brief delay between batches - not too long
    if (batch < totalBatches - 1) {
      console.log(`   ⏳ Waiting 1s between batches...`);
      sleep(1.0); // OPTIMIZED - shorter delay
    }
  }

  const successRate = ((userTokens.length / maxUsers) * 100).toFixed(1);
  console.log(
    `🎯 Setup complete: ${userTokens.length}/${maxUsers} users ready (${successRate}% success)`
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
      `${config.baseUrls[config.environment]}/api/v1/users/profile`,
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
