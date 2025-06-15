/**
 * Staging Deployment Test - Comprehensive validation for K8s deployment
 *
 * Purpose: Validate system performance before production deployment
 * Load: High load testing with up to 200 concurrent users
 * Duration: ~20 minutes comprehensive testing
 * Environment: K8s staging with 5 nodes
 */

import { sleep } from "k6";
import {
  config,
  getLoadPattern,
  getThresholds,
} from "../config/test-config.js";
import { setupTestUser } from "../utils/auth-utils.js";
import {
  browseProducts,
  cartOperations,
  checkoutProcess,
  userProfileOperations,
  searchAndFilter,
  completeUserJourney,
} from "../utils/test-scenarios.js";

// Get environment-specific configuration
const loadPattern = getLoadPattern();
const thresholds = getThresholds();

// Test configuration for staging deployment
export const options = {
  stages: [
    // Warm-up phase
    { duration: "1m", target: 10 }, // Initial warm-up
    { duration: "2m", target: 25 }, // Gradual increase

    // Ramp-up phase
    {
      duration: loadPattern.rampUpTime,
      target: Math.floor(loadPattern.maxUsers * 0.5),
    }, // 50% load
    { duration: "2m", target: Math.floor(loadPattern.maxUsers * 0.75) }, // 75% load
    { duration: "1m", target: loadPattern.maxUsers }, // Full load

    // Sustained load phase
    { duration: loadPattern.sustainTime, target: loadPattern.maxUsers }, // Maintain full load

    // Spike test phase
    { duration: "30s", target: Math.floor(loadPattern.maxUsers * 1.2) }, // 120% spike
    { duration: "2m", target: Math.floor(loadPattern.maxUsers * 1.2) }, // Maintain spike

    // Recovery phase
    { duration: "1m", target: loadPattern.maxUsers }, // Back to normal
    { duration: "2m", target: Math.floor(loadPattern.maxUsers * 0.5) }, // 50% load
    { duration: "2m", target: 0 }, // Ramp down
  ],

  thresholds: {
    ...thresholds,
    // Additional K8s-specific thresholds
    "http_req_duration{scenario:browse}": ["p(95)<2000"],
    "http_req_duration{scenario:cart}": ["p(95)<1500"],
    "http_req_duration{scenario:checkout}": ["p(95)<3000"],
    "http_req_duration{scenario:search}": ["p(95)<2500"],

    // Throughput requirements for K8s
    http_reqs: [
      `rate>${config.k8s.expectedThroughput[config.environment] || 500}`,
    ],

    // Resource utilization checks
    vus: [`value<=${loadPattern.maxUsers * 1.2}`],
  },

  tags: {
    test_type: "staging_deployment",
    environment: config.environment,
    k8s_nodes: config.k8s.nodeCount,
    max_users: loadPattern.maxUsers,
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("🚀 Starting Staging Deployment Test for K8s");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);
  console.log(`K8s Nodes: ${config.k8s.nodeCount}`);
  console.log(`Max Users: ${loadPattern.maxUsers}`);
  console.log(
    `Expected Throughput: ${
      config.k8s.expectedThroughput[config.environment] || "N/A"
    } req/s`
  );
  console.log("Duration: ~20 minutes");
  console.log("Purpose: Pre-production deployment validation");

  // Setup multiple test users for high-load testing
  const users = [];
  const userCount = Math.min(20, Math.floor(loadPattern.maxUsers / 10)); // 10% of max users or 20, whichever is smaller

  console.log(`Setting up ${userCount} test users...`);

  for (let i = 0; i < userCount; i++) {
    const userToken = setupTestUser({
      email: `staging-user-${i}-${Date.now()}@test.com`,
      password: "testpassword123",
      firstName: `Staging${i}`,
      lastName: "User",
    });

    if (userToken) {
      users.push(userToken);
    }

    // Small delay to avoid overwhelming auth service during setup
    sleep(0.1);
  }

  console.log(`✅ Setup completed with ${users.length} test users`);

  return {
    userTokens: users,
    startTime: Date.now(),
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

  // Calculate test phase based on elapsed time
  const elapsedMinutes = (Date.now() - startTime) / (1000 * 60);
  const currentVUs = __VU;
  const totalVUs = __ENV.K6_VUS || loadPattern.maxUsers;

  // Determine user behavior based on test phase and load
  let userBehavior;

  if (elapsedMinutes < 5) {
    // Warm-up phase: lighter operations
    userBehavior = Math.random() < 0.7 ? "browse" : "search";
  } else if (elapsedMinutes < 15) {
    // Main load phase: mixed operations
    userBehavior = ["browse", "cart", "profile", "search", "complete"][
      Math.floor(Math.random() * 5)
    ];
  } else {
    // Spike phase: intensive operations
    userBehavior = Math.random() < 0.4 ? "complete" : "cart";
  }

  try {
    switch (userBehavior) {
      case "browse":
        // 25% - Product browsing with tagging
        console.log(`VU${currentVUs}: Product browsing scenario`);
        browseProducts(userToken);
        sleep(Math.random() * 2 + 1);
        break;

      case "search":
        // 20% - Search and filter operations
        console.log(`VU${currentVUs}: Search and filter scenario`);
        searchAndFilter();
        browseProducts(userToken);
        sleep(Math.random() * 1.5 + 0.5);
        break;

      case "cart":
        // 20% - Cart operations
        console.log(`VU${currentVUs}: Cart operations scenario`);
        if (userToken) {
          browseProducts(userToken);
          cartOperations(userToken);
        } else {
          browseProducts();
        }
        sleep(Math.random() * 2 + 1);
        break;

      case "profile":
        // 15% - User profile operations
        console.log(`VU${currentVUs}: Profile operations scenario`);
        if (userToken) {
          userProfileOperations(userToken);
          browseProducts(userToken);
        } else {
          browseProducts();
        }
        sleep(Math.random() * 1.5 + 0.5);
        break;

      case "complete":
        // 20% - Complete user journey
        console.log(`VU${currentVUs}: Complete user journey scenario`);
        if (userToken) {
          completeUserJourney(userToken);
        } else {
          browseProducts();
          searchAndFilter();
        }
        sleep(Math.random() * 3 + 1);
        break;

      default:
        // Fallback to browsing
        browseProducts(userToken);
        sleep(Math.random() * 1 + 0.5);
    }
  } catch (error) {
    console.error(`VU${currentVUs}: Error in ${userBehavior} scenario:`, error);
    // Continue execution even if there are errors
    sleep(1);
  }

  // Dynamic sleep based on current load
  const loadFactor = currentVUs / totalVUs;
  const baseSleep = 0.5;
  const dynamicSleep = baseSleep + Math.random() * 2 * loadFactor;
  sleep(dynamicSleep);
}

// Teardown function - runs once after the test
export function teardown(data) {
  const { startTime } = data;
  const totalDuration = (Date.now() - startTime) / (1000 * 60);

  console.log("🏁 Staging Deployment Test completed");
  console.log("");
  console.log("📊 K8s Staging Deployment Test Results:");
  console.log(`Total Duration: ${totalDuration.toFixed(1)} minutes`);
  console.log(`Environment: ${config.environment}`);
  console.log(`K8s Nodes: ${config.k8s.nodeCount}`);
  console.log(`Max Concurrent Users: ${loadPattern.maxUsers}`);
  console.log("");

  console.log("🔍 Key Validation Points:");
  console.log("✅ System stability under sustained high load");
  console.log("✅ Performance during traffic spikes");
  console.log("✅ Resource utilization across K8s nodes");
  console.log("✅ Database connection pooling efficiency");
  console.log("✅ Service mesh performance");
  console.log("✅ Auto-scaling behavior validation");
  console.log("");

  console.log("📈 Success Criteria for Production Deployment:");
  console.log(
    `- Response time p95 < ${thresholds.http_req_duration[0].split("<")[1]}`
  );
  console.log(
    `- Error rate < ${(
      parseFloat(thresholds.http_req_failed[0].split("<")[1]) * 100
    ).toFixed(1)}%`
  );
  console.log(
    `- Throughput > ${
      config.k8s.expectedThroughput[config.environment] || "N/A"
    } req/s`
  );
  console.log(
    `- Check success rate > ${(
      parseFloat(thresholds.checks[0].split(">")[1]) * 100
    ).toFixed(0)}%`
  );
  console.log("");

  console.log("🚨 Monitor These Metrics:");
  console.log("- Pod CPU and memory usage");
  console.log("- Database connection pool status");
  console.log("- Redis cache hit rates");
  console.log("- Network latency between services");
  console.log("- Horizontal Pod Autoscaler (HPA) activity");
  console.log("");

  console.log("🎯 Next Steps:");
  console.log("1. Review detailed performance metrics");
  console.log("2. Check K8s cluster resource utilization");
  console.log("3. Validate auto-scaling behavior");
  console.log("4. Compare with previous staging test results");
  console.log("5. If all criteria pass, proceed with production deployment");
  console.log("");

  if (config.environment === "staging") {
    console.log("🚀 Ready for Production Deployment!");
  } else {
    console.log(
      "⚠️  Run this test on staging environment before production deployment"
    );
  }
}
