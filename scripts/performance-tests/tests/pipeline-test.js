/**
 * Pipeline Test - Quick performance validation for CI/CD
 *
 * Purpose: Fast performance check suitable for CI/CD pipelines
 * Load: Light load with essential functionality checks
 * Duration: Under 5 minutes total
 * Output: JSON results for automated analysis
 */

import { sleep } from "k6";
import { config } from "../config/test-config.js";
import { setupTestUser } from "../utils/auth-utils.js";
import {
  browseProducts,
  cartOperations,
  userProfileOperations,
  searchAndFilter,
} from "../utils/test-scenarios.js";

// Test configuration for pipeline test
export const options = {
  stages: [
    { duration: "30s", target: 5 }, // Quick ramp up to 5 users
    { duration: "2m", target: 10 }, // Increase to 10 users
    { duration: "1m", target: 10 }, // Maintain 10 users
    { duration: "30s", target: 0 }, // Quick ramp down
  ],
  thresholds: {
    // Stricter thresholds for pipeline - must pass for deployment
    http_req_duration: ["p(95)<2500"], // 95% under 2.5s
    http_req_failed: ["rate<0.02"], // Error rate under 2%
    checks: ["rate>0.98"], // 98% of checks must pass

    // Specific endpoint thresholds
    "http_req_duration{endpoint:products}": ["p(95)<2000"],
    "http_req_duration{endpoint:auth}": ["p(95)<1500"],
    "http_req_duration{endpoint:cart}": ["p(95)<1000"],
  },
  tags: {
    test_type: "pipeline",
    environment: config.environment,
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("🚀 Starting Pipeline Performance Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);
  console.log("Duration: ~4 minutes");
  console.log("Purpose: CI/CD performance validation");

  // Setup a single test user for pipeline testing
  const customerToken = setupTestUser({
    email: `pipeline-customer-${Date.now()}@test.com`,
    password: "testpassword123",
    firstName: "Pipeline",
    lastName: "Customer",
  });

  if (!customerToken) {
    console.error(
      "❌ Failed to setup test user - pipeline test may be incomplete"
    );
  }

  return {
    customerToken: customerToken,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { customerToken } = data;

  // Essential functionality tests for pipeline validation

  // Test 1: Product browsing (critical path)
  browseProducts(customerToken);
  sleep(0.5);

  // Test 2: Search functionality (critical path)
  searchAndFilter();
  sleep(0.5);

  // Test 3: Authenticated operations (if token available)
  if (customerToken) {
    // User profile operations
    userProfileOperations(customerToken);
    sleep(0.5);

    // Cart operations
    cartOperations(customerToken);
    sleep(0.5);
  }

  // Minimal sleep to simulate user behavior
  sleep(Math.random() * 1 + 0.5);
}

// Teardown function - runs once after the test
export function teardown(data) {
  console.log("🏁 Pipeline Performance Test completed");
  console.log("");
  console.log("📊 Pipeline Test Results Summary:");
  console.log(
    "This test validates essential performance criteria for deployment"
  );
  console.log("");
  console.log("✅ Success Criteria:");
  console.log("- All critical endpoints respond within acceptable time limits");
  console.log("- Error rate is below 2%");
  console.log("- 98% of functional checks pass");
  console.log("");
  console.log("❌ Failure Criteria (blocks deployment):");
  console.log("- Any endpoint exceeds response time thresholds");
  console.log("- Error rate exceeds 2%");
  console.log("- Less than 98% of checks pass");
  console.log("");
  console.log(
    "📈 Check the JSON output for detailed metrics and automated analysis"
  );
}
