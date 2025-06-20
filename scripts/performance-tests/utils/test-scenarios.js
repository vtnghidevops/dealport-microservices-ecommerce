import http from "k6/http";
import { check, sleep } from "k6";
import { getBaseUrl, headers, config } from "../config/test-config.js";
import { getAuthHeaders } from "./auth-utils.js";

/**
 * Browse products scenario - simulates user browsing product catalog
 * @param {string} token - Optional auth token for personalized experience
 */
export function browseProducts(token = null) {
  const requestHeaders = token ? getAuthHeaders(token) : headers;

  // Get all products
  const productsResponse = http.get(
    `${getBaseUrl()}/api/v1/products?page=1&limit=20`,
    { headers: requestHeaders }
  );

  check(productsResponse, {
    "products list loaded": (r) => r.status === 200,
    "products response time OK": (r) => r.timings.duration < 2000,
  });

  sleep(1); // Simulate user reading time

  // Get categories
  const categoriesResponse = http.get(`${getBaseUrl()}/api/v1/categories`, {
    headers: requestHeaders,
  });

  check(categoriesResponse, {
    "categories loaded": (r) => r.status === 200,
    "categories response time OK": (r) => r.timings.duration < 1000,
  });

  sleep(0.5);

  // Get a specific product (simulate clicking on product)
  const productDetailResponse = http.get(`${getBaseUrl()}/api/v1/products/1`, {
    headers: requestHeaders,
  });

  check(productDetailResponse, {
    "product detail loaded": (r) => r.status === 200,
    "product detail response time OK": (r) => r.timings.duration < 1500,
  });

  sleep(2); // Simulate user viewing product details
}

/**
 * Cart operations scenario - add, update, remove items from cart
 * @param {string} token - Auth token (required for cart operations)
 */
export function cartOperations(token) {
  if (!token) {
    console.error("Cart operations require authentication token");
    return;
  }

  const authHeaders = getAuthHeaders(token);

  // Get current cart
  const getCartResponse = http.get(`${getBaseUrl()}/api/v1/cart`, {
    headers: authHeaders,
  });

  check(getCartResponse, {
    "get cart successful": (r) => r.status === 200,
    "get cart response time OK": (r) => r.timings.duration < 1000,
  });

  sleep(0.5);

  // Add item to cart
  const addItemPayload = {
    product_id: config.testProducts[0].id,
    quantity: config.testProducts[0].quantity,
  };

  const addItemResponse = http.post(
    `${getBaseUrl()}/api/v1/cart/items`,
    JSON.stringify(addItemPayload),
    { headers: authHeaders }
  );

  check(addItemResponse, {
    "add item to cart successful": (r) => r.status === 200 || r.status === 201,
    "add item response time OK": (r) => r.timings.duration < 1500,
  });

  sleep(1);

  // Update cart item quantity
  const updateItemPayload = {
    quantity: config.testProducts[0].quantity + 1,
  };

  const updateItemResponse = http.put(
    `${getBaseUrl()}/api/v1/cart/items/1`, // Assuming item ID 1
    JSON.stringify(updateItemPayload),
    { headers: authHeaders }
  );

  check(updateItemResponse, {
    "update cart item successful": (r) => r.status === 200,
    "update item response time OK": (r) => r.timings.duration < 1000,
  });

  sleep(0.5);
}

/**
 * Checkout scenario - validate checkout and create order
 * @param {string} token - Auth token (required for checkout)
 */
export function checkoutProcess(token) {
  if (!token) {
    console.error("Checkout process requires authentication token");
    return;
  }

  const authHeaders = getAuthHeaders(token);

  // Validate checkout
  const validatePayload = {
    billing_info: {
      first_name: "Test",
      last_name: "User",
      email: "test@example.com",
      phone: "+1234567890",
      address: "123 Test St",
      city: "Test City",
      state: "Test State",
      zip_code: "12345",
      country: "US",
    },
    shipping_info: {
      first_name: "Test",
      last_name: "User",
      address: "123 Test St",
      city: "Test City",
      state: "Test State",
      zip_code: "12345",
      country: "US",
    },
  };

  const validateResponse = http.post(
    `${getBaseUrl()}/api/v1/checkout/validate`,
    JSON.stringify(validatePayload),
    { headers: authHeaders }
  );

  check(validateResponse, {
    "checkout validation successful": (r) => r.status === 200,
    "checkout validation response time OK": (r) => r.timings.duration < 2000,
  });

  sleep(1);

  // Create order
  const orderPayload = Object.assign({}, validatePayload, {
    payment_method: "credit_card",
    notes: "Performance test order",
  });

  const createOrderResponse = http.post(
    `${getBaseUrl()}/api/v1/checkout/orders`,
    JSON.stringify(orderPayload),
    { headers: authHeaders }
  );

  check(createOrderResponse, {
    "order creation successful": (r) => r.status === 200 || r.status === 201,
    "order creation response time OK": (r) => r.timings.duration < 3000,
  });

  sleep(2);
}

/**
 * User profile operations scenario
 * @param {string} token - Auth token (required)
 */
export function userProfileOperations(token) {
  if (!token) {
    console.error("User profile operations require authentication token");
    return;
  }

  const authHeaders = getAuthHeaders(token);

  // Get user profile
  const profileResponse = http.get(`${getBaseUrl()}/api/v1/users/me`, {
    headers: authHeaders,
  });

  check(profileResponse, {
    "get profile successful": (r) => r.status === 200,
    "get profile response time OK": (r) => r.timings.duration < 1000,
  });

  sleep(0.5);

  // Get wishlist
  const wishlistResponse = http.get(
    `${getBaseUrl()}/api/v1/users/me/wishlist`,
    { headers: authHeaders }
  );

  check(wishlistResponse, {
    "get wishlist successful": (r) => r.status === 200,
    "get wishlist response time OK": (r) => r.timings.duration < 1000,
  });

  sleep(0.5);
}

/**
 * Search and filter scenario
 */
export function searchAndFilter() {
  // Search products
  const searchResponse = http.get(
    `${getBaseUrl()}/api/v1/products?search=test&page=1&limit=10`,
    { headers }
  );

  check(searchResponse, {
    "product search successful": (r) => r.status === 200,
    "search response time OK": (r) => r.timings.duration < 2000,
  });

  sleep(1);

  // Filter by category
  const filterResponse = http.get(
    `${getBaseUrl()}/api/v1/products?category_id=1&page=1&limit=10`,
    { headers }
  );

  check(filterResponse, {
    "product filter successful": (r) => r.status === 200,
    "filter response time OK": (r) => r.timings.duration < 2000,
  });

  sleep(0.5);
}

/**
 * Complete user journey - from browsing to checkout
 * @param {string} token - Auth token
 */
export function completeUserJourney(token) {
  console.log("Starting complete user journey...");

  // 1. Browse products
  browseProducts(token);

  // 2. Search and filter
  searchAndFilter();

  // 3. User profile operations (if authenticated)
  if (token) {
    userProfileOperations(token);
  }

  // 4. Cart operations (if authenticated)
  if (token) {
    cartOperations(token);
  }

  // 5. Checkout process (if authenticated)
  if (token) {
    checkoutProcess(token);
  }

  console.log("Completed user journey");
}
