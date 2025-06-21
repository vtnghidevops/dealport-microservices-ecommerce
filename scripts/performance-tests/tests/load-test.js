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
  authenticationFlow,
  productManagement,
  categoryManagement,
  enhancedShopping,
  explorePromotions,
  paymentFlow,
  viewHomepageContent,
} from "../utils/test-scenarios.js";
import { getAuthHeaders } from "../utils/auth-utils.js";
import { http, check } from "k6";

// Get environment-specific configuration
const loadPattern = getLoadPattern();
const thresholds = getThresholds();

// Test configuration for load test
export const options = {
  // Increase setup timeout for large user setups
  setupTimeout: "10m", // Allow 10 minutes for 300 user setup

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
  // Target: ${config.baseUrls[config.environment]}

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
  // Need enough users to cover most VUs for realistic authenticated testing
  const userCount = Math.min(
    loadPattern.maxUsers, // RESTORED - match max VUs for full coverage
    250 // RESTORED high cap but not too high to avoid overwhelming setup
  );
  console.log(`Attempting to authenticate ${userCount} shared accounts...`);

  const userTokens = setupMultipleTestUsers(userCount, {
    batchSize: 25, // OPTIMIZED batch size for parallel processing
    setupDelay: 0.3, // OPTIMIZED delay - faster than before but not too fast
  });

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

  // Get authenticated user token for this VU (same approach as smoke test)
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[(__VU - 1) % userTokens.length] // Use modulo like smoke test
      : null;

  const currentVUs = __VU;
  const elapsedMinutes = (Date.now() - startTime) / (1000 * 60);

  // Simple logging for debugging
  if (userToken) {
    console.log(
      `VU${currentVUs}: Using token: ${userToken.substring(0, 20)}...`
    );
  } else {
    console.log(`VU${currentVUs}: Running as guest user`);
  }

  // Distribute user behavior based on realistic e-commerce patterns
  const userBehavior = Math.random();

  try {
    if (userBehavior < 0.25) {
      // 25% - Window shopping with homepage content
      console.log(`VU${currentVUs}: Window shopping pattern`);

      // Start with homepage content (realistic user journey)
      viewHomepageContent(userToken);
      sleep(Math.random() * 1 + 0.5);

      // Browse products with enhanced features
      enhancedShopping(userToken);
      sleep(Math.random() * 1 + 0.5);

      // Explore promotions (realistic bargain hunting)
      explorePromotions(userToken);
    } else if (userBehavior < 0.5) {
      // 25% - Search and enhanced browsing
      console.log(`VU${currentVUs}: Enhanced browsing pattern`);

      // Start with homepage to see what's featured
      viewHomepageContent(userToken);
      sleep(Math.random() * 1 + 0.5);

      // Enhanced shopping with reviews
      enhancedShopping(userToken);
      sleep(Math.random() * 1 + 0.5);

      // Traditional search/browse
      searchAndFilter(userToken);
    } else if (userBehavior < 0.8) {
      // 30% - Authenticated operations (increased from 20%)
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
        console.log(`VU${currentVUs}: Guest user - comprehensive browsing`);
        browseProducts(); // Guest can browse
        searchAndFilter(); // Guest can search
        categoryManagement(); // Guest can view categories
        checkoutProcess(); // Guest can validate checkout (without actual order)
      }
    } else if (userBehavior < 0.9) {
      // 10% - Complete purchase journey with payment
      if (userToken) {
        console.log(`VU${currentVUs}: Complete purchase journey with payment`);

        // Start with promotions (users often look for deals before buying)
        explorePromotions(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Enhanced shopping experience
        enhancedShopping(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Complete purchase with payment flow
        paymentFlow(userToken);
      } else {
        console.log(`VU${currentVUs}: Guest user - complete shopping journey`);
        viewHomepageContent(); // Start with homepage
        enhancedShopping(); // Enhanced browsing
        explorePromotions(); // Check for deals
        checkoutProcess(); // Guest checkout validation
      }
    } else if (userBehavior < 0.95) {
      // 5% - Product and category management
      if (userToken) {
        console.log(`VU${currentVUs}: Product management testing`);

        // Test product management features
        productManagement(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Test category management
        categoryManagement(userToken);
      } else {
        console.log(`VU${currentVUs}: Guest user - category browsing`);
        categoryManagement();
        browseProducts();
      }
    } else {
      // 5% - Account management and admin operations
      if (userToken) {
        console.log(`VU${currentVUs}: Account management pattern`);

        // Focus on profile operations
        userProfileOperations(userToken);
        sleep(Math.random() * 1 + 0.5);

        // Check orders using same approach as smoke test
        const authHeaders = {
          "Content-Type": "application/json",
          Authorization: `Bearer ${userToken}`,
        };

        const ordersResponse = http.get(
          `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
          {
            headers: authHeaders,
            tags: { scenario: "account", type: "orders" },
          }
        );
        check(ordersResponse, {
          "orders retrieved": (r) => r.status === 200,
        });

        // Test authentication flow cycling
        if (Math.random() < 0.3) {
          authenticationFlow();
        }
      } else {
        console.log(`VU${currentVUs}: Guest user - account exploration`);
        browseProducts(); // Browse as guest
        productManagement(); // View product details (no auth needed for viewing)
        categoryManagement(); // Explore all categories
        authenticationFlow(); // Guest trying to login (realistic scenario)
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
