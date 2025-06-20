/**
 * Smoke Test - Verifies basic functionality with minimal load
 *
 * Purpose: Ensure the system is working and all critical endpoints are accessible
 * Load: 1-2 virtual users for 30 seconds
 * Thresholds: Basic response time and error rate checks
 */

import { sleep } from "k6";
import { config } from "../config/test-config.js";
import { loginUser, setupTestUser } from "../utils/auth-utils.js";
import {
  browseProducts,
  cartOperations,
  userProfileOperations,
} from "../utils/test-scenarios.js";

// Test configuration for smoke test
export const options = {
  stages: [
    { duration: "10s", target: 1 }, // Ramp up to 1 user
    { duration: "20s", target: 2 }, // Stay at 2 users
    { duration: "10s", target: 0 }, // Ramp down to 0 users
  ],
  thresholds: config.thresholds.smoke,
  tags: {
    test_type: "smoke",
    environment: config.environment,
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("Starting Smoke Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);

  // Always return data to allow test to continue even if setup fails
  let customerToken = null;

  try {
    // Setup test users - Note: May fail in staging due to OTP/email verification
    customerToken = setupTestUser({
      email: `smoke-customer-${Date.now()}@test.com`,
      password: "testpassword123",
      firstName: "Smoke",
      lastName: "Customer",
    });

    if (customerToken) {
      console.log("Authentication successful - full test coverage available");
    } else {
      console.log(
        "ℹAuthentication not available - testing public endpoints only"
      );
      console.log(
        "   This is expected in staging environments with OTP/email verification"
      );
    }
  } catch (error) {
    console.log(
      "ℹAuthentication setup error (expected in staging):",
      error.message
    );
    console.log("   Continuing with public endpoint tests...");
  }

  return {
    customerToken: customerToken,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { customerToken } = data;

  // Test 1: Basic product browsing (no auth required)
  console.log("Testing product browsing...");
  browseProducts();

  sleep(1);

  // Test 2: Authenticated user operations
  if (customerToken) {
    console.log("Testing authenticated operations...");

    // Test user profile operations
    userProfileOperations(customerToken);

    sleep(1);

    // Test cart operations
    cartOperations(customerToken);

    sleep(1);
  } else {
    console.warn("No customer token available, skipping authenticated tests");
  }

  // Random sleep to simulate real user behavior
  sleep(Math.random() * 2 + 1);
}

// Teardown function - runs once after the test
export function teardown(data) {
  console.log("🏁 Smoke Test completed");
  console.log(
    "Check the results above to ensure all critical paths are working"
  );
}
