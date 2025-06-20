/**
 * Load Test - Simulates normal expected load
 *
 * Purpose: Verify system performance under typical usage patterns
 * Load: Gradually increase to 50 users, maintain for 5 minutes, then ramp down
 * Thresholds: Stricter response time and error rate requirements
 */

import { sleep } from "k6";
import {
  config,
  getLoadPattern,
  getThresholds,
} from "../config/test-config.js";
import { setupMultipleTestUsers } from "../utils/auth-utils.js";
import {
  browseProducts,
  cartOperations,
  checkoutProcess,
  userProfileOperations,
  searchAndFilter,
  completeUserJourney,
  windowShopping,
  quickSearch,
} from "../utils/test-scenarios.js";
import { getAuthHeaders } from "../utils/auth-utils.js";
import { http, check } from "k6";

// Get environment-specific configuration
const loadPattern = getLoadPattern();
const thresholds = getThresholds();

// Test configuration for load test
export const options = {
  stages: [
    { duration: "2m", target: Math.floor(loadPattern.maxUsers * 0.2) }, // 20% ramp up
    {
      duration: loadPattern.rampUpTime,
      target: Math.floor(loadPattern.maxUsers * 0.5),
    }, // 50% ramp up
    { duration: "2m", target: loadPattern.maxUsers }, // Full load
    { duration: loadPattern.sustainTime, target: loadPattern.maxUsers }, // Sustain load
    { duration: "3m", target: Math.floor(loadPattern.maxUsers * 0.5) }, // Ramp down to 50%
    { duration: "2m", target: 0 }, // Ramp down to 0
  ],
  thresholds: thresholds,
  tags: {
    test_type: "load",
    environment: config.environment,
    max_users: loadPattern.maxUsers.toString(),
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("📊 Starting Load Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);

  // Calculate total duration
  const totalDuration =
    parseFloat(loadPattern.rampUpTime) +
    parseFloat(loadPattern.sustainTime) +
    9; // +9 for other stages
  console.log(
    `Expected load: Up to ${loadPattern.maxUsers} concurrent users for ~${totalDuration} minutes total`
  );
  console.log(
    `Ramp up time: ${loadPattern.rampUpTime}, Sustain time: ${loadPattern.sustainTime}`
  );

  // Setup multiple shared test accounts for better distribution
  const userCount = Math.min(
    10,
    config.testUsers.performanceTestAccounts.length
  );
  console.log(`Attempting to authenticate ${userCount} shared accounts...`);

  const userTokens = setupMultipleTestUsers(userCount);

  console.log(
    `✅ Setup completed with ${userTokens.length}/${userCount} authenticated accounts`
  );
  console.log(`Starting load test with max ${loadPattern.maxUsers} users...`);

  return {
    userTokens: userTokens,
    startTime: Date.now(),
    maxUsers: loadPattern.maxUsers,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { userTokens, startTime } = data;

  // Randomly select a user token for this iteration
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[Math.floor(Math.random() * userTokens.length)]
      : null;

  const currentVUs = __VU;
  const elapsedMinutes = (Date.now() - startTime) / (1000 * 60);

  // Distribute user behavior based on realistic e-commerce patterns
  const userBehavior = Math.random();

  try {
    if (userBehavior < 0.3) {
      // 30% - Window shopping (browsing without purchase intent)
      console.log(`VU${currentVUs}: Window shopping pattern`);

      // Browse categories and products
      browseProducts(userToken);
      sleep(Math.random() * 1 + 0.5);

      // Browse more products
      browseProducts(userToken);
    } else if (userBehavior < 0.6) {
      // 30% - Search and browse behavior
      console.log(`VU${currentVUs}: Search and browse pattern`);

      // Search/browse products
      searchAndFilter(userToken);
      sleep(Math.random() * 1 + 0.5);

      // Follow up with more browsing
      browseProducts(userToken);
    } else if (userBehavior < 0.8) {
      // 20% - Authenticated cart operations
      if (userToken) {
        console.log(`VU${currentVUs}: Authenticated cart operations`);

        // Browse products first
        browseProducts(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Add items to cart
        cartOperations(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Check profile/orders
        userProfileOperations(userToken);
      } else {
        console.log(`VU${currentVUs}: No auth - fallback to browsing`);
        browseProducts();
        searchAndFilter();
      }
    } else if (userBehavior < 0.95) {
      // 15% - Complete purchase journey
      if (userToken) {
        console.log(`VU${currentVUs}: Complete purchase journey`);

        // Full user journey
        completeUserJourney(userToken);
      } else {
        console.log(`VU${currentVUs}: No auth - guest checkout attempt`);
        browseProducts();
        checkoutProcess(); // Guest checkout validation
      }
    } else {
      // 5% - Profile and account management
      if (userToken) {
        console.log(`VU${currentVUs}: Account management pattern`);

        // Focus on profile operations
        userProfileOperations(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Check orders and wishlist
        const headers = getAuthHeaders(userToken);

        // Get orders
        const ordersResponse = http.get(
          `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
          { headers, tags: { scenario: "account", type: "orders" } }
        );
        check(ordersResponse, {
          "orders retrieved": (r) => r.status === 200,
        });

        // Get wishlist
        const wishlistResponse = http.get(
          `${config.baseUrls[config.environment]}/api/v1/users/me/wishlist`,
          { headers, tags: { scenario: "account", type: "wishlist" } }
        );
        check(wishlistResponse, {
          "wishlist retrieved": (r) => r.status === 200,
        });
      } else {
        console.log(`VU${currentVUs}: No auth - fallback to browsing`);
        browseProducts();
      }
    }
  } catch (error) {
    console.error(`VU${currentVUs}: Load test error:`, error.message);
    // Continue execution even if there are errors
    sleep(0.5);
  }

  // Think time between user actions
  sleep(Math.random() * 2 + 1);
}

// Teardown function - runs once after the test
export function teardown(data) {
  console.log("🏁 Load Test completed");
  console.log(
    "Review the performance metrics to ensure the system can handle expected load"
  );
  console.log("Key metrics to check:");
  console.log("- Average response time should be under 3 seconds");
  console.log("- 95th percentile response time should be under 5 seconds");
  console.log("- Error rate should be under 5%");
  console.log("- All checks should pass at least 95% of the time");
}
