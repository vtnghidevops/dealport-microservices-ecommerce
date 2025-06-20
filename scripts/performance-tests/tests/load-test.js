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
  guestCheckoutFlow,
  windowShoppingFlow,
} from "../utils/test-scenarios.js";

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
  const { userTokens, startTime, maxUsers } = data;

  // Randomly select a user token for this iteration
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[Math.floor(Math.random() * userTokens.length)]
      : null;

  // Realistic user behavior patterns based on e-commerce analytics
  const userBehavior = Math.random();

  if (userBehavior < 0.4) {
    // 40% - Window shoppers (just browsing, no purchase intent)
    console.log("Simulating window shopping user...");
    windowShoppingFlow();
    sleep(Math.random() * 3 + 2); // 2-5 seconds thinking time
  } else if (userBehavior < 0.7) {
    // 30% - Guest users who validate checkout but can't complete
    console.log("Simulating guest checkout flow...");
    guestCheckoutFlow();
    sleep(Math.random() * 2 + 1); // 1-3 seconds thinking time
  } else if (userBehavior < 0.85) {
    // 15% - Authenticated users with complete journey
    console.log("Simulating authenticated complete journey...");
    if (userToken) {
      completeUserJourney(userToken);
    } else {
      // Fallback to guest flow if no auth
      console.log("No auth token, falling back to guest checkout...");
      guestCheckoutFlow();
    }
    sleep(Math.random() * 2 + 1); // 1-3 seconds thinking time
  } else if (userBehavior < 0.95) {
    // 10% - Authenticated cart users (add to cart but don't checkout)
    console.log("Simulating authenticated cart operations...");
    if (userToken) {
      browseProducts(userToken);
      userProfileOperations(userToken);
      cartOperations(userToken);
    } else {
      // Fallback to window shopping if no auth
      console.log("No auth token, falling back to browsing...");
      windowShoppingFlow();
    }
    sleep(Math.random() * 2 + 1); // 1-3 seconds thinking time
  } else {
    // 5% - Quick searchers (users who know what they want)
    console.log("Simulating quick search user...");
    searchAndFilter();
    browseProducts(userToken);

    // Try guest checkout validation
    if (Math.random() < 0.5) {
      guestCheckoutFlow();
    }

    sleep(Math.random() * 1 + 0.5); // 0.5-1.5 seconds thinking time
  }

  // Random pause to simulate real user behavior
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
