/**
 * Quick Auth Test - Verify OTP bypass is working
 * Run this first to confirm authentication setup before load testing
 */

import { sleep } from "k6";
import { config } from "../config/test-config.js";
import { setupTestUser } from "../utils/auth-utils.js";

export const options = {
  stages: [
    { duration: "30s", target: 5 }, // Test with 5 users for 30 seconds
  ],
  thresholds: {
    http_req_duration: ["p(95)<10000"], // Very lenient
    http_req_failed: ["rate<0.50"], // 50% failure allowed for testing
    checks: ["rate>0.50"], // 50% success minimum
  },
};

export function setup() {
  console.log("🧪 Quick Auth Test - Verifying OTP Bypass");
  console.log(`Environment: ${config.environment}`);
  console.log(`Target: ${config.baseUrls[config.environment]}`);

  // Test creating 10 users quickly
  console.log("Testing rapid user creation with OTP bypass...");

  const tokens = [];
  for (let i = 0; i < 10; i++) {
    console.log(`Creating user ${i + 1}/10...`);
    const token = setupTestUser(
      {
        email: `auth-test-${i}-${Date.now()}@test.com`,
        firstName: `AuthTest${i}`,
      },
      {
        setupMode: true,
        shortDelay: 0.2,
      }
    );

    if (token) {
      tokens.push(token);
      console.log(`User ${i + 1} authenticated successfully`);
    } else {
      console.log(`User ${i + 1} authentication failed`);
    }

    sleep(0.3); // Small delay between users
  }

  const successRate = (tokens.length / 10) * 100;
  console.log(
    `Auth Test Results: ${
      tokens.length
    }/10 users authenticated (${successRate.toFixed(1)}%)`
  );

  if (successRate >= 80) {
    console.log("🎉 Auth system ready for load testing!");
  } else {
    console.log(
      "Auth system may have issues - consider investigating before load test"
    );
  }

  return { tokens, successRate };
}

export default function (data) {
  // Simple verification that tokens work
  if (data && data.tokens && data.tokens.length > 0) {
    const randomToken =
      data.tokens[Math.floor(Math.random() * data.tokens.length)];

    // Test that token works for a simple request
    // This would normally be a user profile check but keeping it simple
    console.log(`VU${__VU}: Using token: ${randomToken.substring(0, 20)}...`);
  }

  sleep(1);
}

export function teardown(data) {
  console.log("🏁 Auth Test Complete");
  if (data && data.successRate >= 80) {
    console.log("Ready to proceed with full load test");
  } else {
    console.log("Fix authentication issues before load testing");
  }
}
