/**
 * Load Test - Simulates normal expected load
 *
 * Purpose: Verify system performance under typical usage patterns
 * Load: Gradually increase to 50 users, maintain for 5 minutes, then ramp down
 * Thresholds: Stricter response time and error rate requirements
 */

import { sleep } from "k6";
import { config } from "../config/test-config.js";
import { setupTestUser } from "../utils/auth-utils.js";
import {
  browseProducts,
  cartOperations,
  checkoutProcess,
  userProfileOperations,
  searchAndFilter,
  completeUserJourney,
} from "../utils/test-scenarios.js";

// Test configuration for load test
export const options = {
  stages: [
    { duration: "2m", target: 10 }, // Ramp up to 10 users over 2 minutes
    { duration: "3m", target: 25 }, // Ramp up to 25 users over 3 minutes
    { duration: "2m", target: 50 }, // Ramp up to 50 users over 2 minutes
    { duration: "5m", target: 50 }, // Stay at 50 users for 5 minutes
    { duration: "2m", target: 25 }, // Ramp down to 25 users over 2 minutes
    { duration: "2m", target: 0 }, // Ramp down to 0 users over 2 minutes
  ],
  thresholds: config.thresholds.load,
  tags: {
    test_type: "load",
    environment: config.environment,
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("📊 Starting Load Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);
  console.log("Expected load: Up to 50 concurrent users for 16 minutes total");

  // Setup multiple test users for different scenarios
  const users = [];

  for (let i = 0; i < 5; i++) {
    const userToken = setupTestUser({
      email: `load-customer-${i}-${Date.now()}@test.com`,
      password: "testpassword123",
      firstName: `Load${i}`,
      lastName: "Customer",
    });

    if (userToken) {
      users.push(userToken);
    }
  }

  console.log(`Setup completed with ${users.length} test users`);

  return {
    userTokens: users,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { userTokens } = data;

  // Randomly select a user token for this iteration
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[Math.floor(Math.random() * userTokens.length)]
      : null;

  // Simulate different user behavior patterns
  const userBehavior = Math.random();

  if (userBehavior < 0.3) {
    // 30% - Browsing only users (window shoppers)
    console.log("Simulating browsing-only user...");
    browseProducts(userToken);
    searchAndFilter();
    sleep(Math.random() * 3 + 2); // 2-5 seconds thinking time
  } else if (userBehavior < 0.6) {
    // 30% - Users who browse and add to cart but don't checkout
    console.log("Simulating browse-and-cart user...");
    browseProducts(userToken);

    if (userToken) {
      userProfileOperations(userToken);
      cartOperations(userToken);
    }

    searchAndFilter();
    sleep(Math.random() * 2 + 1); // 1-3 seconds thinking time
  } else if (userBehavior < 0.85) {
    // 25% - Complete user journey (browse, cart, checkout)
    console.log("Simulating complete user journey...");

    if (userToken) {
      completeUserJourney(userToken);
    } else {
      // Fallback to browsing if no token
      browseProducts();
      searchAndFilter();
    }

    sleep(Math.random() * 2 + 1); // 1-3 seconds thinking time
  } else {
    // 15% - Quick searchers (users who know what they want)
    console.log("Simulating quick search user...");
    searchAndFilter();

    // Quick product view
    browseProducts(userToken);

    if (userToken && Math.random() < 0.7) {
      // 70% chance to add to cart
      cartOperations(userToken);
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
