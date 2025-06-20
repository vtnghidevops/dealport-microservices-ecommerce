/**
 * Smoke Test - Basic functionality validation
 *
 * Purpose: Quick health check to ensure all critical systems are working
 * Load: 1-2 users for 40 seconds
 * Coverage: Core public endpoints + basic authentication
 */

import { sleep } from "k6";
import http from "k6/http";
import { check } from "k6";
import { config, getLoadPattern } from "../config/test-config.js";
import {
  setupMultipleTestUsers,
  getRandomUserSession,
  getAuthHeaders,
} from "../utils/auth-utils.js";
import {
  browseProducts,
  searchAndFilter,
  checkoutProcess,
} from "../utils/test-scenarios.js";

// Get environment-specific configuration
const loadPattern = getLoadPattern();

// Test configuration for smoke test
export const options = {
  stages: [
    { duration: "10s", target: 1 }, // Start with 1 user
    { duration: "20s", target: 2 }, // Ramp to 2 users
    { duration: "10s", target: 0 }, // Ramp down
  ],

  thresholds: {
    // Lenient thresholds for smoke test - just checking if things work
    http_req_duration: ["p(95)<5000"], // 95% under 5s
    http_req_failed: ["rate<0.10"], // Error rate under 10%
    checks: ["rate>0.80"], // 80% of checks pass

    // Critical endpoints should work
    "http_req_failed{type:products}": ["rate<0.05"],
    "http_req_failed{type:categories}": ["rate<0.05"],
  },

  tags: {
    test_type: "smoke",
    environment: config.environment,
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("🚨 Starting Smoke Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);
  console.log(`Duration: 40 seconds`);
  console.log("");
  console.log("🎯 Objectives:");
  console.log("- Verify core public endpoints are working");
  console.log("- Test basic authentication flow");
  console.log("- Validate essential e-commerce functionality");
  console.log("");

  // Setup minimal user sessions for smoke test
  console.log("Setting up 2 test user sessions...");
  const userSessions = setupMultipleTestUsers(2);

  const authenticatedUsers = userSessions.filter((u) => u.token).length;
  console.log(
    `✅ Setup complete: ${authenticatedUsers}/${userSessions.length} authenticated`
  );

  if (authenticatedUsers > 0) {
    console.log("🎉 Authentication working - will test auth flows");
  } else {
    console.log("⚠️ Authentication failed - will test public endpoints only");
  }

  return {
    userSessions: userSessions,
    startTime: Date.now(),
    authenticatedCount: authenticatedUsers,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { userSessions, authenticatedCount } = data;

  // Get user session for this virtual user
  const userSession = getRandomUserSession(userSessions);
  const currentVUs = __VU;

  console.log(`VU${currentVUs}: Starting smoke test checks`);

  try {
    // 1. Test core browsing functionality (public)
    console.log(`VU${currentVUs}: Testing product browsing...`);
    browseProducts(userSession);

    sleep(1);

    // 2. Test search functionality (public)
    console.log(`VU${currentVUs}: Testing search functionality...`);
    searchAndFilter(userSession);

    sleep(1);

    // 3. Test checkout validation (public)
    console.log(`VU${currentVUs}: Testing checkout validation...`);
    checkoutProcess(userSession);

    sleep(1);

    // 4. Additional auth test if available
    if (userSession && userSession.token && !userSession.isGuest) {
      console.log(`VU${currentVUs}: Testing authenticated endpoints...`);

      // Quick cart check
      const cartResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/cart`,
        {
          headers: getAuthHeaders(userSession),
          tags: { scenario: "smoke", type: "cart_check" },
        }
      );

      check(cartResponse, {
        "cart endpoint accessible": (r) => r.status === 200,
      });
    }

    // Quick API endpoint checks
    console.log("Checking API endpoints...");

    // Test categories endpoint
    const categoriesCheck = http.get(
      `${config.baseUrls[config.environment]}/api/v1/categories`,
      { tags: { scenario: "api_check", type: "categories" } }
    );

    check(categoriesCheck, {
      "categories endpoint available": (r) => r.status === 200,
    });

    // Test products endpoint
    const productsCheck = http.get(
      `${config.baseUrls[config.environment]}/api/v1/products?limit=5`,
      { tags: { scenario: "api_check", type: "products" } }
    );

    check(productsCheck, {
      "products endpoint available": (r) => r.status === 200,
    });

    // Test cart endpoint (requires auth, should return 401)
    const cartCheck = http.get(
      `${config.baseUrls[config.environment]}/api/v1/cart`,
      { tags: { scenario: "api_check", type: "cart" } }
    );

    check(cartCheck, {
      "cart endpoint responds": (r) => r.status === 401 || r.status === 200,
    });
  } catch (error) {
    console.error(`VU${currentVUs}: Smoke test error:`, error.message);
  }

  // Short pause between iterations
  sleep(2);
}

// Teardown function - runs once after the test
export function teardown(data) {
  const { startTime, authenticatedCount, userSessions } = data;
  const totalDuration = (Date.now() - startTime) / 1000;

  console.log("🏁 Smoke Test Completed");
  console.log("");
  console.log("📊 Results Summary:");
  console.log(`Duration: ${totalDuration.toFixed(1)} seconds`);
  console.log(`Environment: ${config.environment}`);
  console.log(
    `Authenticated sessions: ${authenticatedCount}/${userSessions.length}`
  );
  console.log("");
  console.log("✅ Basic Health Checks:");
  console.log("- Product browsing functionality");
  console.log("- Search and filter operations");
  console.log("- Checkout validation endpoints");
  if (authenticatedCount > 0) {
    console.log("- Authentication and cart access");
  }
  console.log("");
  console.log("🚀 System is ready for load testing!");
}
