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
import http from "k6/http";
import { check } from "k6";

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
    // === AUTHENTICATION FLOW TESTING (if no token, test auth process) ===
    if (!userToken && Math.random() < 0.1) {
      // 10% of users without tokens test auth flow
      console.log(`VU${currentVUs}: Testing authentication flow...`);

      // Test registration + verification
      const timestamp = Date.now();
      const randomId = Math.random().toString(36).substr(2, 9);
      const testUser = {
        email: `load-test-${timestamp}-${randomId}@test.com`,
        password: "LoadTest123!",
        first_name: "Load",
        last_name: "Test",
        phone: "+1234567890",
      };

      const registerResponse = http.post(
        `${config.baseUrls[config.environment]}/api/v1/auth/register`,
        JSON.stringify(testUser),
        {
          headers: { "Content-Type": "application/json" },
          tags: { scenario: "load", type: "auth_register" },
        }
      );
      check(registerResponse, {
        "registration successful": (r) => r.status === 200 || r.status === 201,
      });

      if (registerResponse.status === 200 || registerResponse.status === 201) {
        sleep(0.2); // Wait for OTP system

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
            tags: { scenario: "load", type: "auth_verify" },
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
          sleep(0.5); // Wait for user to be fully created

          const loginResponse = http.post(
            `${config.baseUrls[config.environment]}/api/v1/auth/login`,
            JSON.stringify({
              email: testUser.email,
              password: testUser.password,
            }),
            {
              headers: { "Content-Type": "application/json" },
              tags: { scenario: "load", type: "auth_login" },
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
                tags: { scenario: "load", type: "auth_validate" },
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
    console.log(`VU${currentVUs}: Testing core endpoints...`);

    // 1. Categories endpoint test
    const categoriesResponse = http.get(
      `${config.baseUrls[config.environment]}/api/v1/categories`,
      { tags: { scenario: "load", type: "categories" } }
    );
    check(categoriesResponse, {
      "categories endpoint works": (r) => r.status === 200,
    });

    // 2. Products endpoint test
    const productsResponse = http.get(
      `${config.baseUrls[config.environment]}/api/v1/products?limit=10`,
      { tags: { scenario: "load", type: "products" } }
    );
    check(productsResponse, {
      "products endpoint works": (r) => r.status === 200,
    });

    // 3. Checkout validation test
    const checkoutValidateResponse = http.post(
      `${config.baseUrls[config.environment]}/api/v1/checkout/validate`,
      JSON.stringify({
        items: [{ product_id: 1, quantity: 2, price: 99.99 }],
        shipping_address: {
          street: "123 Load Test St",
          city: "Load City",
          postal_code: "12345",
          country: "US",
        },
      }),
      {
        headers: { "Content-Type": "application/json" },
        tags: { scenario: "load", type: "checkout_validate" },
      }
    );
    check(checkoutValidateResponse, {
      "checkout validation works": (r) =>
        r.status === 200 || r.status === 400 || r.status === 422,
    });

    // 4. AUTHENTICATED ENDPOINTS (if token available)
    if (userToken) {
      console.log(`VU${currentVUs}: Testing authenticated endpoints...`);

      const authHeaders = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${userToken}`,
      };

      // User profile test
      const profileResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/users/profile`,
        { headers: authHeaders, tags: { scenario: "load", type: "profile" } }
      );
      check(profileResponse, {
        "user profile works": (r) => r.status === 200,
      });

      // Cart access test
      const cartResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/cart`,
        { headers: authHeaders, tags: { scenario: "load", type: "cart_get" } }
      );
      check(cartResponse, {
        "cart access works": (r) => r.status === 200,
      });

      // Add to cart test
      const addToCartResponse = http.post(
        `${config.baseUrls[config.environment]}/api/v1/cart/items`,
        JSON.stringify({
          product_id: Math.floor(Math.random() * 10) + 1,
          quantity: Math.floor(Math.random() * 3) + 1,
        }),
        { headers: authHeaders, tags: { scenario: "load", type: "cart_add" } }
      );
      check(addToCartResponse, {
        "add to cart works": (r) => r.status === 200 || r.status === 201,
      });

      // Orders access test
      const ordersResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
        { headers: authHeaders, tags: { scenario: "load", type: "orders" } }
      );
      check(ordersResponse, {
        "orders access works": (r) => r.status === 200,
      });

      // Order creation test (realistic load test scenario)
      if (Math.random() < 0.1) {
        // 10% of users try to create orders
        const orderCreateResponse = http.post(
          `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
          JSON.stringify({
            items: [
              {
                product_id: Math.floor(Math.random() * 10) + 1,
                quantity: 1,
                price: 99.99,
              },
            ],
            shipping_address: {
              street: "123 Load Test St",
              city: "Load City",
              postal_code: "12345",
              country: "US",
            },
            payment_method: "stripe",
          }),
          {
            headers: authHeaders,
            tags: { scenario: "load", type: "order_create" },
          }
        );
        check(orderCreateResponse, {
          "order creation works": (r) =>
            r.status === 200 || r.status === 201 || r.status === 422,
        });
      }
    }

    // === ENHANCED LOAD TEST SCENARIOS ===
    // Create user session object for scenario functions
    const userSession = userToken ? { token: userToken, id: currentVUs } : null;

    // Add realistic load test behavior patterns on top of core endpoint testing
    if (userBehavior < 0.25) {
      // 25% - Enhanced browsing with multiple product views
      console.log(`VU${currentVUs}: Enhanced browsing pattern`);
      enhancedShopping(userSession);
      sleep(Math.random() * 2 + 1);

      // View homepage content
      viewHomepageContent(userSession);
    } else if (userBehavior < 0.5) {
      // 25% - Search and filter intensive usage
      console.log(`VU${currentVUs}: Search intensive pattern`);
      searchAndFilter(userSession);
      sleep(Math.random() * 1 + 0.5);

      // Explore promotions
      explorePromotions(userSession);
    } else if (userBehavior < 0.8) {
      // 30% - Cart and checkout operations
      if (userSession) {
        console.log(`VU${currentVUs}: Cart operations pattern`);
        cartOperations(userSession);
        sleep(Math.random() * 1 + 0.5);

        // User profile operations
        userProfileOperations(userSession);
      } else {
        console.log(`VU${currentVUs}: Guest browsing pattern`);
        browseProducts();
        sleep(Math.random() * 1 + 0.5);
        windowShopping();
      }
    } else {
      // 20% - Complete purchase journey
      if (userSession) {
        console.log(`VU${currentVUs}: Complete purchase journey`);
        completeUserJourney(userSession);
        sleep(Math.random() * 2 + 1);

        // Payment flow
        paymentFlow(userSession);
      } else {
        console.log(`VU${currentVUs}: Guest complete journey`);
        completeUserJourney();
        sleep(Math.random() * 1 + 0.5);
        checkoutProcess();
      }
    }

    // Random sleep to simulate realistic user behavior
    sleep(Math.random() * 2 + 1);
  } catch (error) {
    console.error(`VU${currentVUs}: Error in load test: ${error.message}`);
  }
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
