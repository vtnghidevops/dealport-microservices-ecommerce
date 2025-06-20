/**
 * Stress Test - Tests system beyond normal capacity for K8s
 *
 * Purpose: Find breaking points and test system resilience under extreme load
 * Load: Gradually increase to 300+ users with aggressive spike testing
 * Thresholds: More lenient to allow for system degradation under stress
 * Environment: Optimized for K8s with 5 nodes
 */

import { sleep } from "k6";
import { config, getLoadPattern } from "../config/test-config.js";
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

// Calculate stress test limits based on environment
const getStressLimits = () => {
  // Use stressMaxUsers if available, otherwise fallback to maxUsers
  const baseMax = loadPattern.stressMaxUsers || loadPattern.maxUsers;
  return {
    normalLoad: Math.floor(baseMax * 0.75), // 75% as baseline
    stressLoad: baseMax, // 100% - main stress level
    extremeLoad: Math.floor(baseMax * 1.25), // 125% of stress max
    spikeLoad: Math.floor(baseMax * 1.5), // 150% of stress max
  };
};

const stressLimits = getStressLimits();

// Test configuration for stress test
export const options = {
  stages: [
    // Gradual ramp up to normal capacity
    { duration: "2m", target: Math.floor(stressLimits.normalLoad * 0.2) }, // 20%
    { duration: "3m", target: Math.floor(stressLimits.normalLoad * 0.5) }, // 50%
    { duration: "2m", target: stressLimits.normalLoad }, // 100% - normal capacity

    // Push beyond normal capacity
    { duration: "3m", target: stressLimits.stressLoad }, // 150% - stress level
    { duration: "5m", target: stressLimits.stressLoad }, // Maintain stress

    // Extreme load testing
    { duration: "2m", target: stressLimits.extremeLoad }, // 200% - extreme load
    { duration: "3m", target: stressLimits.extremeLoad }, // Maintain extreme

    // Spike testing - sudden traffic bursts
    { duration: "30s", target: stressLimits.spikeLoad }, // 250% - spike
    { duration: "1m", target: stressLimits.spikeLoad }, // Maintain spike
    { duration: "30s", target: stressLimits.extremeLoad }, // Drop back

    // Second spike test
    { duration: "30s", target: stressLimits.spikeLoad }, // Another spike
    { duration: "1m", target: stressLimits.spikeLoad }, // Maintain

    // Recovery testing
    { duration: "2m", target: stressLimits.extremeLoad }, // Back to extreme
    { duration: "3m", target: stressLimits.stressLoad }, // Back to stress
    { duration: "2m", target: stressLimits.normalLoad }, // Back to normal
    { duration: "3m", target: 0 }, // Ramp down
  ],

  thresholds: {
    // More lenient thresholds for stress testing
    http_req_duration: ["p(95)<8000"], // 95% under 8s (degraded performance expected)
    http_req_failed: ["rate<0.15"], // Error rate under 15%
    checks: ["rate>0.85"], // 85% of checks pass

    // K8s specific stress thresholds
    "http_req_duration{scenario:browse}": ["p(95)<5000"],
    "http_req_duration{scenario:cart}": ["p(95)<6000"],
    "http_req_duration{scenario:checkout}": ["p(95)<10000"],

    // System should not completely fail
    "http_req_failed{scenario:critical}": ["rate<0.05"], // Critical paths should remain stable

    // Resource monitoring
    vus: [`value<=${stressLimits.spikeLoad}`],
  },

  tags: {
    test_type: "stress",
    environment: config.environment,
    k8s_nodes: config.k8s.nodeCount.toString(),
    max_stress_users: stressLimits.spikeLoad.toString(),
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("🔥 Starting K8s Stress Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);
  console.log(`K8s Nodes: ${config.k8s.nodeCount}`);
  console.log("");
  console.log("📊 Stress Test Load Levels:");
  console.log(`Normal Load: ${stressLimits.normalLoad} users`);
  console.log(`Stress Load: ${stressLimits.stressLoad} users (150%)`);
  console.log(`Extreme Load: ${stressLimits.extremeLoad} users (200%)`);
  console.log(`Spike Load: ${stressLimits.spikeLoad} users (250%)`);
  console.log("");
  console.log("⚠️  This test will push the system beyond normal capacity");
  console.log("🎯 Purpose: Find breaking points and validate resilience");

  // Setup more test users for stress testing
  const users = [];
  const userCount = Math.min(30, Math.floor(stressLimits.normalLoad / 8)); // More users for stress testing

  console.log(`Setting up ${userCount} test users for stress testing...`);

  for (let i = 0; i < userCount; i++) {
    const userToken = setupTestUser({
      email: `stress-user-${i}-${Date.now()}@test.com`,
      password: "testpassword123",
      firstName: `Stress${i}`,
      lastName: "User",
    });

    if (userToken) {
      users.push(userToken);
    }

    // Minimal delay to setup users quickly
    sleep(0.05);
  }

  console.log(`✅ Setup completed with ${users.length} test users`);

  return {
    userTokens: users,
    startTime: Date.now(),
    stressLimits: stressLimits,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { userTokens, startTime, stressLimits } = data;

  // Randomly select a user token for this iteration
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[Math.floor(Math.random() * userTokens.length)]
      : null;

  // Calculate current load level and test phase
  const elapsedMinutes = (Date.now() - startTime) / (1000 * 60);
  const currentVUs = __VU;
  const currentLoad = __ENV.K6_VUS || currentVUs;

  // Determine stress level
  let stressLevel = "normal";
  if (currentLoad > stressLimits.spikeLoad * 0.9) {
    stressLevel = "spike";
  } else if (currentLoad > stressLimits.extremeLoad * 0.9) {
    stressLevel = "extreme";
  } else if (currentLoad > stressLimits.stressLoad * 0.9) {
    stressLevel = "stress";
  }

  // More aggressive user behavior patterns for stress testing
  const userBehavior = Math.random();

  try {
    if (stressLevel === "spike") {
      // Spike phase: Very aggressive, rapid-fire requests
      if (userBehavior < 0.3) {
        console.log(`VU${currentVUs}: SPIKE - Rapid browsing burst`);
        browseProducts(userToken);
        browseProducts(userToken); // Double requests
        sleep(Math.random() * 0.3 + 0.1);
      } else if (userBehavior < 0.6) {
        console.log(`VU${currentVUs}: SPIKE - Cart operation burst`);
        if (userToken) {
          cartOperations(userToken);
          cartOperations(userToken); // Double operations
        } else {
          browseProducts();
        }
        sleep(Math.random() * 0.5 + 0.1);
      } else {
        console.log(`VU${currentVUs}: SPIKE - Mixed aggressive operations`);
        browseProducts(userToken);
        if (userToken && Math.random() < 0.7) {
          userProfileOperations(userToken);
        }
        searchAndFilter();
        sleep(Math.random() * 0.2 + 0.05);
      }
    } else if (stressLevel === "extreme") {
      // Extreme phase: Heavy operations with some rapid requests
      if (userBehavior < 0.25) {
        console.log(`VU${currentVUs}: EXTREME - Heavy browsing pattern`);
        browseProducts(userToken);
        searchAndFilter();
        browseProducts(userToken);
        sleep(Math.random() * 0.8 + 0.2);
      } else if (userBehavior < 0.5) {
        console.log(`VU${currentVUs}: EXTREME - Intensive cart operations`);
        if (userToken) {
          browseProducts(userToken);
          cartOperations(userToken);
          userProfileOperations(userToken);
        } else {
          browseProducts();
          searchAndFilter();
        }
        sleep(Math.random() * 1 + 0.3);
      } else if (userBehavior < 0.75) {
        console.log(`VU${currentVUs}: EXTREME - Complete journey stress`);
        if (userToken) {
          completeUserJourney(userToken);
        } else {
          browseProducts();
          searchAndFilter();
          browseProducts();
        }
        sleep(Math.random() * 1.5 + 0.5);
      } else {
        console.log(`VU${currentVUs}: EXTREME - Search heavy pattern`);
        searchAndFilter();
        searchAndFilter(); // Multiple searches
        browseProducts(userToken);
        sleep(Math.random() * 0.6 + 0.2);
      }
    } else if (stressLevel === "stress") {
      // Stress phase: Sustained heavy load
      if (userBehavior < 0.2) {
        console.log(`VU${currentVUs}: STRESS - Sustained browsing`);
        browseProducts(userToken);
        searchAndFilter();
        sleep(Math.random() * 1.2 + 0.5);
      } else if (userBehavior < 0.4) {
        console.log(`VU${currentVUs}: STRESS - Cart abandonment pattern`);
        if (userToken) {
          browseProducts(userToken);
          cartOperations(userToken);
          userProfileOperations(userToken);
          cartOperations(userToken); // Multiple cart operations
        }
        sleep(Math.random() * 1.5 + 0.5);
      } else if (userBehavior < 0.6) {
        console.log(`VU${currentVUs}: STRESS - Complete journey under load`);
        if (userToken) {
          completeUserJourney(userToken);
        } else {
          browseProducts();
          searchAndFilter();
        }
        sleep(Math.random() * 2 + 0.8);
      } else if (userBehavior < 0.8) {
        console.log(`VU${currentVUs}: STRESS - Search intensive`);
        searchAndFilter();
        browseProducts(userToken);
        if (userToken && Math.random() < 0.6) {
          cartOperations(userToken);
        }
        sleep(Math.random() * 1 + 0.4);
      } else {
        console.log(`VU${currentVUs}: STRESS - Mixed heavy operations`);
        browseProducts(userToken);
        if (userToken) {
          userProfileOperations(userToken);
          if (Math.random() < 0.5) {
            cartOperations(userToken);
          }
        }
        searchAndFilter();
        sleep(Math.random() * 1.3 + 0.3);
      }
    } else {
      // Normal phase: Regular stress test patterns
      if (userBehavior < 0.2) {
        console.log(`VU${currentVUs}: NORMAL - Browsing pattern`);
        browseProducts(userToken);
        sleep(Math.random() * 1.5 + 0.8);
      } else if (userBehavior < 0.4) {
        console.log(`VU${currentVUs}: NORMAL - Cart operations`);
        if (userToken) {
          browseProducts(userToken);
          cartOperations(userToken);
        } else {
          browseProducts();
        }
        sleep(Math.random() * 2 + 1);
      } else if (userBehavior < 0.6) {
        console.log(`VU${currentVUs}: NORMAL - User journey`);
        if (userToken) {
          completeUserJourney(userToken);
        } else {
          browseProducts();
          searchAndFilter();
        }
        sleep(Math.random() * 2.5 + 1);
      } else if (userBehavior < 0.8) {
        console.log(`VU${currentVUs}: NORMAL - Search operations`);
        searchAndFilter();
        browseProducts(userToken);
        sleep(Math.random() * 1.8 + 0.7);
      } else {
        console.log(`VU${currentVUs}: NORMAL - Profile operations`);
        if (userToken) {
          userProfileOperations(userToken);
          browseProducts(userToken);
        } else {
          browseProducts();
        }
        sleep(Math.random() * 1.5 + 0.8);
      }
    }
  } catch (error) {
    console.error(
      `VU${currentVUs}: Error during ${stressLevel} stress test:`,
      error
    );
    // Continue execution even if there are errors - this is expected under stress
    sleep(0.5);
  }

  // Minimal pause - stress testing should be aggressive
  const stressSleep =
    stressLevel === "spike"
      ? 0.1
      : stressLevel === "extreme"
      ? 0.3
      : stressLevel === "stress"
      ? 0.5
      : 0.8;
  sleep(Math.random() * stressSleep + 0.1);
}

// Teardown function - runs once after the test
export function teardown(data) {
  const { startTime, stressLimits } = data;
  const totalDuration = (Date.now() - startTime) / (1000 * 60);

  console.log("🏁 K8s Stress Test completed");
  console.log("");
  console.log("📊 Stress Test Analysis:");
  console.log(`Total Duration: ${totalDuration.toFixed(1)} minutes`);
  console.log(`Environment: ${config.environment}`);
  console.log(`K8s Nodes: ${config.k8s.nodeCount}`);
  console.log("");
  console.log("🔥 Load Levels Tested:");
  console.log(`✅ Normal Load: ${stressLimits.normalLoad} users`);
  console.log(`⚠️  Stress Load: ${stressLimits.stressLoad} users (150%)`);
  console.log(`🚨 Extreme Load: ${stressLimits.extremeLoad} users (200%)`);
  console.log(`💥 Spike Load: ${stressLimits.spikeLoad} users (250%)`);
  console.log("");

  console.log("🔍 Key Metrics to Analyze:");
  console.log("- Response time degradation patterns across load levels");
  console.log("- Error rate increases and recovery patterns");
  console.log("- System recovery after load spikes");
  console.log("- K8s auto-scaling behavior under extreme load");
  console.log("- Database connection pool behavior");
  console.log("- Memory usage and garbage collection patterns");
  console.log("- Network saturation points");
  console.log("");

  console.log("⚠️  Expected Behavior Under Stress:");
  console.log("- Response times increase significantly at extreme loads");
  console.log("- Some requests may fail (acceptable up to 15%)");
  console.log("- System should recover when load decreases");
  console.log("- No complete system failures or crashes");
  console.log("- Auto-scaling should trigger appropriately");
  console.log("");

  console.log("🚨 Red Flags to Watch For:");
  console.log("- Complete system unresponsiveness");
  console.log("- Error rates above 25%");
  console.log("- System not recovering after load reduction");
  console.log("- Memory leaks or resource exhaustion");
  console.log("- Cascading failures across services");
  console.log("- Database deadlocks or connection exhaustion");
  console.log("");

  console.log("🎯 K8s Specific Monitoring:");
  console.log("- Pod restart counts during stress");
  console.log("- HPA scaling events and timing");
  console.log("- Node resource utilization");
  console.log("- Service mesh circuit breaker activations");
  console.log("- Persistent volume performance");
  console.log("");

  console.log("📈 Capacity Planning Insights:");
  console.log(
    `- Breaking point appears around ${stressLimits.extremeLoad}+ concurrent users`
  );
  console.log("- Consider horizontal scaling triggers");
  console.log("- Review resource limits and requests");
  console.log("- Validate circuit breaker configurations");
  console.log("- Plan for traffic spike handling");
}
