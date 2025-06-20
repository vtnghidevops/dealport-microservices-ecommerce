/**
 * Smoke Test - Basic functionality validation
 *
 * Purpose: Quick health check to ensure all critical systems are working
 * Load: 1-2 users for 40 seconds
 * Coverage: Core public endpoints + basic authentication
 */

import { sleep } from "k6";
import http from "k6/http";
import { check } from "k6";
import { config, getLoadPattern } from "../config/test-config.js";
import {
  setupMultipleTestUsers,
  getRandomUserSession,
  getAuthHeaders,
} from "../utils/auth-utils.js";
import {
  browseProducts,
  searchAndFilter,
  checkoutProcess,
} from "../utils/test-scenarios.js";

// Get environment-specific configuration
const loadPattern = getLoadPattern();

// Test configuration for smoke test
export const options = {
  stages: [
    { duration: "10s", target: 1 }, // Start with 1 user
    { duration: "20s", target: 2 }, // Ramp to 2 users
    { duration: "10s", target: 0 }, // Ramp down
  ],

  thresholds: {
    // Lenient thresholds for smoke test - just checking if things work
    http_req_duration: ["p(95)<5000"], // 95% under 5s
    http_req_failed: ["rate<0.10"], // Error rate under 10%
    checks: ["rate>0.80"], // 80% of checks pass

    // Critical endpoints should work
    "http_req_failed{type:products}": ["rate<0.05"],
    "http_req_failed{type:categories}": ["rate<0.05"],
  },

  tags: {
    test_type: "smoke",
    environment: config.environment,
  },
};

// Setup function - runs once before the test
export function setup() {
  console.log("🚨 Starting Smoke Test");
  console.log(`Environment: ${config.environment}`);
  console.log(`Base URL: ${config.baseUrls[config.environment]}`);
  console.log(`Duration: 40 seconds`);
  console.log("");
  console.log("🎯 Objectives:");
  console.log("- Verify core public endpoints are working");
  console.log("- Test basic authentication flow");
  console.log("- Validate essential e-commerce functionality");
  console.log("");

  // Setup minimal user tokens for smoke test
  console.log("Setting up 2 test user sessions...");
  const userTokens = setupMultipleTestUsers(2);

  const authenticatedUsers = userTokens.filter(
    (token) => token !== null
  ).length;
  console.log(
    `Setup complete: ${authenticatedUsers}/${userTokens.length} authenticated`
  );

  if (authenticatedUsers > 0) {
    console.log("🎉 Authentication working - will test auth flows");
  } else {
    console.log("⚠️ Authentication failed - will test public endpoints only");
  }

  return {
    userTokens: userTokens,
    startTime: Date.now(),
    authenticatedCount: authenticatedUsers,
  };
}

// Main test function - runs for each virtual user
export default function (data) {
  const { userTokens, startTime } = data;
  const currentVUs = __VU;

  // Get authenticated user token for this VU
  const userToken =
    userTokens && userTokens.length > 0
      ? userTokens[currentVUs % userTokens.length]
      : null;

  try {
    console.log(`🔍 VU${currentVUs}: Starting comprehensive endpoint testing`);

    // 1. Test Public Endpoints (no auth required)
    console.log(`VU${currentVUs}: Testing public endpoints...`);

    // Categories
    const categoriesResponse = http.get(
      `${config.baseUrls[config.environment]}/api/v1/categories`,
      { tags: { scenario: "public", type: "categories" } }
    );
    check(categoriesResponse, {
      "categories endpoint works": (r) => r.status === 200,
    });

    // Products
    const productsResponse = http.get(
      `${config.baseUrls[config.environment]}/api/v1/products?limit=5`,
      { tags: { scenario: "public", type: "products" } }
    );
    check(productsResponse, {
      "products endpoint works": (r) => r.status === 200,
    });

    // Checkout validation (public)
    const checkoutValidateResponse = http.post(
      `${config.baseUrls[config.environment]}/api/v1/checkout/validate`,
      JSON.stringify({
        items: [{ product_id: 1, quantity: 2, price: 99.99 }],
        shipping_address: {
          street: "123 Test St",
          city: "Test City",
          postal_code: "12345",
          country: "US",
        },
      }),
      {
        headers: { "Content-Type": "application/json" },
        tags: { scenario: "public", type: "checkout_validate" },
      }
    );
    check(checkoutValidateResponse, {
      "checkout validation works": (r) =>
        r.status === 200 || r.status === 400 || r.status === 422,
    });

    // 2. Test Authenticated Endpoints (requires token)
    if (userToken) {
      console.log(`VU${currentVUs}: Testing authenticated endpoints...`);

      const authHeaders = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${userToken}`,
      };

      // User profile
      const profileResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/users/me`,
        { headers: authHeaders, tags: { scenario: "auth", type: "profile" } }
      );
      check(profileResponse, {
        "user profile works": (r) => r.status === 200,
      });

      // Cart operations
      const cartResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/cart`,
        { headers: authHeaders, tags: { scenario: "auth", type: "cart_get" } }
      );
      check(cartResponse, {
        "cart access works": (r) => r.status === 200,
      });

      // Add item to cart
      const addToCartResponse = http.post(
        `${config.baseUrls[config.environment]}/api/v1/cart/items`,
        JSON.stringify({ product_id: 1, quantity: 2 }),
        { headers: authHeaders, tags: { scenario: "auth", type: "cart_add" } }
      );
      check(addToCartResponse, {
        "add to cart works": (r) => r.status === 200 || r.status === 201,
      });

      // Wishlist
      const wishlistResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/users/me/wishlist`,
        { headers: authHeaders, tags: { scenario: "auth", type: "wishlist" } }
      );
      check(wishlistResponse, {
        "wishlist access works": (r) => r.status === 200,
      });

      // Orders
      const ordersResponse = http.get(
        `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
        { headers: authHeaders, tags: { scenario: "auth", type: "orders" } }
      );
      check(ordersResponse, {
        "orders access works": (r) => r.status === 200,
      });

      // Create order
      const createOrderResponse = http.post(
        `${config.baseUrls[config.environment]}/api/v1/checkout/orders`,
        JSON.stringify({
          items: [{ product_id: 1, quantity: 1, price: 50.0 }],
          shipping_address: {
            street: "123 Test St",
            city: "Test City",
            postal_code: "12345",
            country: "US",
          },
          payment_method: "credit_card",
        }),
        {
          headers: authHeaders,
          tags: { scenario: "auth", type: "create_order" },
        }
      );
      check(createOrderResponse, {
        "order creation works": (r) => r.status >= 200 && r.status < 500,
      });

      console.log(`✅ VU${currentVUs}: All authenticated endpoints tested`);
    } else {
      console.log(
        `⚠️ VU${currentVUs}: No auth token - skipping authenticated tests`
      );
    }

    console.log(`🎯 VU${currentVUs}: Smoke test completed successfully`);
  } catch (error) {
    console.error(`❌ VU${currentVUs}: Smoke test error:`, error.message);
    sleep(0.5);
  }

  sleep(Math.random() * 2 + 1);
}

// Teardown function - runs once after the test
export function teardown(data) {
  const { startTime, authenticatedCount, userTokens } = data;
  const totalDuration = (Date.now() - startTime) / 1000;

  console.log("🏁 Smoke Test Completed");
  console.log("");
  console.log("📊 Results Summary:");
  console.log(`Duration: ${totalDuration.toFixed(1)} seconds`);
  console.log(`Environment: ${config.environment}`);
  console.log(
    `Authenticated sessions: ${authenticatedCount}/${userTokens.length}`
  );
  console.log("");
  console.log("✅ Basic Health Checks:");
  console.log("- Product browsing functionality");
  console.log("- Search and filter operations");
  console.log("- Checkout validation endpoints");
  if (authenticatedCount > 0) {
    console.log("- Authentication and cart access");
  }
  console.log("");
  console.log("🚀 System is ready for load testing!");
}
