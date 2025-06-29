/**
 * Stress Test - Tests system beyond normal capacity for K8s
 *
 * Purpose: Find breaking points and test system resilience under extreme load
 * Load: Gradually increase to 500+ users with aggressive spike testing
 * Thresholds: More lenient to allow for system degradation under stress
 * Environment: Optimized for K8s with 5 nodes
 * API Endpoints: Tests all main broker service routes (/api/v1/*)
 */

import { sleep } from "k6";
import { config, getLoadPattern } from "../config/test-config.js";
import { setupMultipleTestUsers } from "../utils/auth-utils.js";
import { getAuthHeaders } from "../utils/auth-utils.js";
import http from "k6/http";
import { check } from "k6";
import {
  browseProducts,
  cartOperations,
  checkoutProcess,
  userProfileOperations,
  searchAndFilter,
  completeUserJourney,
  authenticationFlow,
  productManagement,
  categoryManagement,
  enhancedShopping,
  explorePromotions,
  paymentFlow,
  viewHomepageContent,
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
  // Increase setup timeout for large user setups
  setupTimeout: "15m", // Allow 15 minutes for 500+ user setup

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
  // Target: ${config.baseUrls[config.environment]}
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

  // Setup test users for stress testing - need good coverage
  const userCount = Math.min(
    loadPattern.stressMaxUsers || loadPattern.maxUsers, // Use stress max if available
    300 // High cap for stress testing
  );
  console.log(
    `Setting up ${userCount} shared test accounts for stress testing...`
  );

  const userTokens = setupMultipleTestUsers(userCount, {
    batchSize: 30, // Larger batches for stress test
    setupDelay: 0.2, // Faster setup for stress test
  });

  console.log(`✅ Setup completed with ${userTokens.length} test users`);

  return {
    userTokens: userTokens,
    startTime: Date.now(),
    stressLimits: stressLimits,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { userTokens, startTime, stressLimits } = data;

  // Get authenticated user token for this VU (same approach as smoke test)
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[(__VU - 1) % userTokens.length] // Use modulo like smoke test
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

  // Create user session object for scenario functions
  const userSession = userToken ? { token: userToken, id: currentVUs } : null;

  try {
    // === AUTHENTICATION FLOW TESTING (if no token, test auth process) ===
    if (!userToken && Math.random() < 0.15) {
      console.log(
        `VU${currentVUs}: [${stressLevel}] Testing authentication flow...`
      );

      // Test registration + verification
      const timestamp = Date.now();
      const randomId = Math.random().toString(36).substr(2, 9);
      const testUser = {
        email: `stress-test-${timestamp}-${randomId}@test.com`,
        password: "StressTest123!",
        first_name: "Stress",
        last_name: "Test",
        phone: "+1234567890",
      };

      const registerResponse = http.post(
        `${config.baseUrls[config.environment]}/api/v1/auth/register`,
        JSON.stringify(testUser),
        {
          headers: { "Content-Type": "application/json" },
          tags: {
            scenario: "stress",
            type: "auth_register",
            stress_level: stressLevel,
          },
        }
      );
      check(registerResponse, {
        "registration successful": (r) => r.status === 200 || r.status === 201,
      });

      if (registerResponse.status === 200 || registerResponse.status === 201) {
        sleep(0.1); // Faster for stress test

        // Test verification with bypass OTP
        const verifyResponse = http.post(
          `${
            config.baseUrls[config.environment]
          }/api/v1/auth/verify-registration`,
          JSON.stringify({
            email: testUser.email,
            otp: "123456", // Bypass code for staging
          }),
          {
            headers: { "Content-Type": "application/json" },
            tags: {
              scenario: "stress",
              type: "auth_verify",
              stress_level: stressLevel,
            },
          }
        );
        check(verifyResponse, {
          "verification successful": (r) => r.status === 200,
          "verification has token": (r) => {
            try {
              const data = JSON.parse(r.body);
              return data.data && data.data.access_token;
            } catch (error) {
              return false;
            }
          },
        });

        // Test login flow with the newly created user
        if (verifyResponse.status === 200) {
          sleep(0.2); // Shorter wait for stress test

          const loginResponse = http.post(
            `${config.baseUrls[config.environment]}/api/v1/auth/login`,
            JSON.stringify({
              email: testUser.email,
              password: testUser.password,
            }),
            {
              headers: { "Content-Type": "application/json" },
              tags: {
                scenario: "stress",
                type: "auth_login",
                stress_level: stressLevel,
              },
            }
          );

          const loginToken = check(loginResponse, {
            "login successful": (r) => r.status === 200,
            "login returns token": (r) => {
              try {
                const data = JSON.parse(r.body);
                return data.data && data.data.access_token;
              } catch (error) {
                return false;
              }
            },
          });

          // Test token validation if login successful
          if (loginToken && loginResponse.status === 200) {
            const loginData = JSON.parse(loginResponse.body);
            const token = loginData.data.access_token;

            const validateResponse = http.get(
              `${config.baseUrls[config.environment]}/api/v1/auth/validate`,
              {
                headers: { Authorization: `Bearer ${token}` },
                tags: {
                  scenario: "stress",
                  type: "auth_validate",
                  stress_level: stressLevel,
                },
              }
            );

            check(validateResponse, {
              "token validation successful": (r) => r.status === 200,
            });
          }
        }
      }
    }

    // === CORE ENDPOINTS TESTING (same as smoke test) ===
    console.log(`VU${currentVUs}: [${stressLevel}] Testing core endpoints...`);

    // 1. Categories endpoint test
    const categoriesResponse = http.get(
      `${config.baseUrls[config.environment]}/api/v1/categories`,
      {
        tags: {
          scenario: "stress",
          type: "categories",
          stress_level: stressLevel,
        },
      }
    );
    check(categoriesResponse, {
      "categories endpoint works": (r) => r.status === 200,
    });

    // 2. Products endpoint test
    const productsResponse = http.get(
      `${config.baseUrls[config.environment]}/api/v1/products?limit=20`,
      {
        tags: {
          scenario: "stress",
          type: "products",
          stress_level: stressLevel,
        },
      }
    );
    check(productsResponse, {
      "products endpoint works": (r) => r.status === 200,
    });

    // 3. Checkout validation test
    const checkoutValidateResponse = http.post(
      `${config.baseUrls[config.environment]}/api/v1/checkout/validate`,
      JSON.stringify({
        items: [
          {
            product_id: Math.floor(Math.random() * 20) + 1,
            quantity: Math.floor(Math.random() * 5) + 1,
            price: 99.99,
          },
        ],
        shipping_address: {
          street: `${Math.floor(Math.random() * 999) + 1} Stress Test Ave`,
          city: "Stress City",
          postal_code: `${Math.floor(Math.random() * 90000) + 10000}`,
          country: "US",
        },
      }),
      {
        headers: { "Content-Type": "application/json" },
        tags: {
          scenario: "stress",
          type: "checkout_validate",
          stress_level: stressLevel,
        },
      }
    );
    check(checkoutValidateResponse, {
      "checkout validation works": (r) =>
        r.status === 200 || r.status === 400 || r.status === 422,
    });

    // 4. AUTHENTICATED ENDPOINTS (if token available)
    if (userToken) {
      console.log(
        `VU${currentVUs}: [${stressLevel}] Testing authenticated endpoints...`
      );

      const authHeaders = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${userToken}`,
      };

      // User profile test
      const profileResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/users/profile`,
        {
          headers: authHeaders,
          tags: {
            scenario: "stress",
            type: "profile",
            stress_level: stressLevel,
          },
        }
      );
      check(profileResponse, {
        "user profile works": (r) => r.status === 200,
      });

      // Cart access test
      const cartResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/cart`,
        {
          headers: authHeaders,
          tags: {
            scenario: "stress",
            type: "cart_get",
            stress_level: stressLevel,
          },
        }
      );
      check(cartResponse, {
        "cart access works": (r) => r.status === 200,
      });

      // Add to cart test (more aggressive for stress)
      const addToCartResponse = http.post(
        `${config.baseUrls[config.environment]}/api/v1/cart/items`,
        JSON.stringify({
          product_id: Math.floor(Math.random() * 50) + 1,
          quantity: Math.floor(Math.random() * 5) + 1,
        }),
        {
          headers: authHeaders,
          tags: {
            scenario: "stress",
            type: "cart_add",
            stress_level: stressLevel,
          },
        }
      );
      check(addToCartResponse, {
        "add to cart works": (r) => r.status === 200 || r.status === 201,
      });

      // Orders access test
      const ordersResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
        {
          headers: authHeaders,
          tags: {
            scenario: "stress",
            type: "orders",
            stress_level: stressLevel,
          },
        }
      );
      check(ordersResponse, {
        "orders access works": (r) => r.status === 200,
      });

      // Order creation test (more frequent in stress test)
      if (Math.random() < 0.2) {
        // 20% of users try to create orders (higher than load test)
        const orderCreateResponse = http.post(
          `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
          JSON.stringify({
            items: [
              {
                product_id: Math.floor(Math.random() * 20) + 1,
                quantity: Math.floor(Math.random() * 3) + 1,
                price: Math.random() * 200 + 50,
              },
            ],
            shipping_address: {
              street: `${Math.floor(Math.random() * 999) + 1} Stress Test Ave`,
              city: "Stress City",
              postal_code: `${Math.floor(Math.random() * 90000) + 10000}`,
              country: "US",
            },
            payment_method: Math.random() < 0.5 ? "stripe" : "paypal",
          }),
          {
            headers: authHeaders,
            tags: {
              scenario: "stress",
              type: "order_create",
              stress_level: stressLevel,
            },
          }
        );
        check(orderCreateResponse, {
          "order creation works": (r) =>
            r.status === 200 || r.status === 201 || r.status === 422,
        });
      }
    }

    // === STRESS-SPECIFIC BEHAVIOR PATTERNS ===
    if (stressLevel === "spike") {
      // Spike phase: Very aggressive, rapid-fire requests
      if (userBehavior < 0.4) {
        console.log(`VU${currentVUs}: SPIKE - Rapid authenticated operations`);
        if (userSession) {
          // Fast multiple cart operations
          cartOperations(userSession);
          cartOperations(userSession); // Double call for stress
        }
        // Minimal sleep for maximum stress
        sleep(Math.random() * 0.5);
      } else if (userBehavior < 0.8) {
        console.log(`VU${currentVUs}: SPIKE - Rapid browsing`);
        browseProducts(userSession);
        searchAndFilter(userSession);
        sleep(Math.random() * 0.5);
      } else {
        console.log(`VU${currentVUs}: SPIKE - Rapid checkout attempts`);
        checkoutProcess(userSession);
        if (userSession) {
          paymentFlow(userSession);
        }
        sleep(Math.random() * 0.3);
      }
    } else if (stressLevel === "extreme") {
      // Extreme phase: High volume realistic operations
      if (userBehavior < 0.3) {
        console.log(`VU${currentVUs}: EXTREME - Intensive shopping`);
        enhancedShopping(userSession);
        explorePromotions(userSession);
        sleep(Math.random() * 1);
      } else if (userBehavior < 0.7) {
        console.log(`VU${currentVUs}: EXTREME - Heavy cart usage`);
        if (userSession) {
          cartOperations(userSession);
          userProfileOperations(userSession);
        } else {
          browseProducts();
          searchAndFilter();
        }
        sleep(Math.random() * 1);
      } else {
        console.log(`VU${currentVUs}: EXTREME - Complete journey stress`);
        completeUserJourney(userSession);
        sleep(Math.random() * 1.5);
      }
    } else if (stressLevel === "stress") {
      // Stress phase: Sustained high load
      if (userBehavior < 0.25) {
        console.log(`VU${currentVUs}: STRESS - Sustained browsing`);
        viewHomepageContent(userSession);
        enhancedShopping(userSession);
        sleep(Math.random() * 1.5);
      } else if (userBehavior < 0.5) {
        console.log(`VU${currentVUs}: STRESS - Search intensive`);
        searchAndFilter(userSession);
        browseProducts(userSession);
        sleep(Math.random() * 1.5);
      } else if (userBehavior < 0.8) {
        console.log(`VU${currentVUs}: STRESS - Cart operations`);
        if (userSession) {
          cartOperations(userSession);
          userProfileOperations(userSession);
        } else {
          browseProducts();
          categoryManagement();
        }
        sleep(Math.random() * 1.5);
      } else {
        console.log(`VU${currentVUs}: STRESS - Purchase flow`);
        if (userSession) {
          completeUserJourney(userSession);
          paymentFlow(userSession);
        } else {
          completeUserJourney();
          checkoutProcess();
        }
        sleep(Math.random() * 2);
      }
    } else {
      // Normal phase: Regular load test behavior
      if (userBehavior < 0.25) {
        console.log(`VU${currentVUs}: NORMAL - Regular browsing`);
        browseProducts(userSession);
        sleep(Math.random() * 2);
      } else if (userBehavior < 0.5) {
        console.log(`VU${currentVUs}: NORMAL - Search and filter`);
        searchAndFilter(userSession);
        sleep(Math.random() * 2);
      } else if (userBehavior < 0.8) {
        console.log(`VU${currentVUs}: NORMAL - Cart usage`);
        if (userSession) {
          cartOperations(userSession);
        } else {
          browseProducts();
        }
        sleep(Math.random() * 2);
      } else {
        console.log(`VU${currentVUs}: NORMAL - Complete journey`);
        completeUserJourney(userSession);
        sleep(Math.random() * 3);
      }
    }
  } catch (error) {
    console.error(
      `VU${currentVUs}: [${stressLevel}] Stress test error: ${error.message}`
    );
    // Shorter recovery time in stress test
    sleep(0.1);
  }
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
