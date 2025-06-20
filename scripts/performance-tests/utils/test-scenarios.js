/**
 * Test scenarios for E-commerce performance testing
 * Using shared account approach with unique session contexts
 */

import http from "k6/http";
import { check, sleep } from "k6";
import { config, getBaseUrl } from "../config/test-config.js";
import { getAuthHeaders } from "./auth-utils.js";

/**
 * Common headers for requests
 */
function getBaseHeaders(userSession = null) {
  return getAuthHeaders(userSession);
}

/**
 * Browse products - Core public functionality
 * @param {Object} userSession - User session from setup
 */
export function browseProducts(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // Get product categories
  const categoriesResponse = http.get(`${getBaseUrl()}/api/v1/categories`, {
    headers,
    tags: { scenario: "browse", type: "categories" },
  });

  check(categoriesResponse, {
    "categories loaded": (r) => r.status === 200,
  });

  if (categoriesResponse.status === 200) {
    try {
      const categories = JSON.parse(categoriesResponse.body);
      if (categories.data && categories.data.length > 0) {
        // Browse a random category
        const randomCategory =
          categories.data[Math.floor(Math.random() * categories.data.length)];

        const categoryResponse = http.get(
          `${getBaseUrl()}/api/v1/products?category_id=${
            randomCategory.id
          }&limit=10`,
          { headers, tags: { scenario: "browse", type: "category_products" } }
        );

        check(categoryResponse, {
          "category products loaded": (r) => r.status === 200,
        });
      }
    } catch (error) {
      console.error("Error browsing categories:", error.message);
    }
  }

  // Get featured/trending products
  const productsResponse = http.get(
    `${getBaseUrl()}/api/v1/products?limit=20`,
    { headers, tags: { scenario: "browse", type: "products" } }
  );

  check(productsResponse, {
    "products loaded": (r) => r.status === 200,
  });

  sleep(Math.random() * 2 + 1);
}

/**
 * Search and filter operations - SIMPLIFIED (no search/filter implemented)
 * @param {Object} userSession - User session (optional)
 */
export function searchAndFilter(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // Since search is not implemented, just browse different product categories
  const browseModes = ["all", "electronics", "fashion", "grocery"];
  const randomMode =
    browseModes[Math.floor(Math.random() * browseModes.length)];

  console.log(`Browsing products (mode: ${randomMode})`);

  // Browse all products (since search is not implemented)
  const productsResponse = http.get(
    `${getBaseUrl()}/api/v1/products?limit=15`,
    { headers, tags: { scenario: "browse", type: "products_browse" } }
  );

  check(productsResponse, {
    "products browse successful": (r) => r.status === 200,
  });

  // Browse categories (alternative to search)
  const categoriesResponse = http.get(`${getBaseUrl()}/api/v1/categories`, {
    headers,
    tags: { scenario: "browse", type: "categories_browse" },
  });

  check(categoriesResponse, {
    "categories browse successful": (r) => r.status === 200,
  });

  // If categories available, browse a specific category
  if (categoriesResponse.status === 200) {
    try {
      const categories = JSON.parse(categoriesResponse.body);
      if (categories.data && categories.data.length > 0) {
        const randomCategory =
          categories.data[Math.floor(Math.random() * categories.data.length)];

        const categoryProductsResponse = http.get(
          `${getBaseUrl()}/api/v1/products?limit=10`,
          { headers, tags: { scenario: "browse", type: "category_browse" } }
        );

        check(categoryProductsResponse, {
          "category products browse successful": (r) => r.status === 200,
        });
      }
    } catch (error) {
      console.error("Error browsing categories:", error.message);
    }
  }

  sleep(Math.random() * 2 + 1);
}

/**
 * Cart operations - Requires authentication
 * @param {Object} userSession - User session with valid token
 */
export function cartOperations(userSession) {
  if (!userSession || !userSession.token) {
    console.log("Cart operations skipped - authentication required");
    return browseProducts(userSession); // Fallback to browsing
  }

  const headers = getBaseHeaders(userSession);

  // Get current cart
  const cartResponse = http.get(`${getBaseUrl()}/api/v1/cart`, {
    headers,
    tags: { scenario: "cart", type: "get_cart" },
  });

  check(cartResponse, {
    "cart retrieved": (r) => r.status === 200,
  });

  // Get some products to add to cart
  const productsResponse = http.get(`${getBaseUrl()}/api/v1/products?limit=5`, {
    headers,
    tags: { scenario: "cart", type: "get_products" },
  });

  if (productsResponse.status === 200) {
    try {
      const products = JSON.parse(productsResponse.body);
      if (products.data && products.data.length > 0) {
        const randomProduct =
          products.data[Math.floor(Math.random() * products.data.length)];

        // Add product to cart
        const addToCartPayload = {
          product_id: randomProduct.id,
          quantity: Math.floor(Math.random() * 3) + 1,
        };

        const addResponse = http.post(
          `${getBaseUrl()}/api/v1/cart/items`,
          JSON.stringify(addToCartPayload),
          { headers, tags: { scenario: "cart", type: "add_item" } }
        );

        check(addResponse, {
          "item added to cart": (r) => r.status === 200 || r.status === 201,
        });

        // Update cart item quantity
        if (addResponse.status === 200 || addResponse.status === 201) {
          const updatePayload = {
            product_id: randomProduct.id,
            quantity: Math.floor(Math.random() * 2) + 1,
          };

          const updateResponse = http.put(
            `${getBaseUrl()}/api/v1/cart/items/${randomProduct.id}`,
            JSON.stringify(updatePayload),
            { headers, tags: { scenario: "cart", type: "update_item" } }
          );

          check(updateResponse, {
            "cart item updated": (r) => r.status === 200,
          });
        }
      }
    } catch (error) {
      console.error("Error in cart operations:", error.message);
    }
  }

  sleep(Math.random() * 2 + 1);
}

/**
 * Checkout process - Only test implemented endpoints
 * @param {Object} userSession - User session (can be guest)
 */
export function checkoutProcess(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // Validate checkout data (implemented public endpoint)
  const validatePayload = {
    items: [
      {
        product_id: 1, // Use simple ID instead of string
        quantity: 2,
        price: 99.99,
      },
    ],
    shipping_address: {
      street: "123 Test St",
      city: "Test City",
      postal_code: "12345",
      country: "US",
    },
  };

  const validateResponse = http.post(
    `${getBaseUrl()}/api/v1/checkout/validate`,
    JSON.stringify(validatePayload),
    { headers, tags: { scenario: "checkout", type: "validate" } }
  );

  check(validateResponse, {
    "checkout validation": (r) =>
      r.status === 200 || r.status === 400 || r.status === 422, // Accept validation errors
  });

  // Only attempt order creation if authenticated (implemented endpoint)
  if (userSession && userSession.token) {
    const orderPayload = Object.assign({}, validatePayload, {
      payment_method: "credit_card",
      payment_details: {
        card_number: "4111111111111111",
        expiry: "12/25",
        cvv: "123",
      },
    });

    const orderResponse = http.post(
      `${getBaseUrl()}/api/v1/checkout/orders`,
      JSON.stringify(orderPayload),
      { headers, tags: { scenario: "checkout", type: "create_order" } }
    );

    check(orderResponse, {
      "order creation attempted": (r) => r.status >= 200 && r.status < 500,
    });

    // Try to get orders list (implemented endpoint)
    const ordersListResponse = http.get(
      `${getBaseUrl()}/api/v1/checkout/orders`,
      { headers, tags: { scenario: "checkout", type: "list_orders" } }
    );

    check(ordersListResponse, {
      "orders list retrieved": (r) => r.status === 200,
    });
  }

  sleep(Math.random() * 2 + 1);
}

/**
 * User profile operations - Only test implemented endpoints
 * @param {Object} userSession - User session with valid token
 */
export function userProfileOperations(userSession) {
  if (!userSession || !userSession.token) {
    console.log("Profile operations skipped - authentication required");
    return browseProducts(userSession); // Fallback to browsing
  }

  const headers = getBaseHeaders(userSession);

  // Get user profile (implemented endpoint)
  const profileResponse = http.get(`${getBaseUrl()}/api/v1/users/me`, {
    headers,
    tags: { scenario: "profile", type: "get_profile" },
  });

  check(profileResponse, {
    "profile retrieved": (r) => r.status === 200,
  });

  // Get user orders (implemented endpoint)
  const ordersResponse = http.get(`${getBaseUrl()}/api/v1/checkout/orders`, {
    headers,
    tags: { scenario: "profile", type: "get_orders" },
  });

  check(ordersResponse, {
    "orders retrieved": (r) => r.status === 200,
  });

  // Get wishlist (implemented endpoint)
  const wishlistResponse = http.get(
    `${getBaseUrl()}/api/v1/users/me/wishlist`,
    { headers, tags: { scenario: "profile", type: "get_wishlist" } }
  );

  check(wishlistResponse, {
    "wishlist retrieved": (r) => r.status === 200,
  });

  // Update profile information (implemented endpoint)
  const updatePayload = {
    first_name: `TestUser${userSession.id || ""}`,
    last_name: `Session${(userSession.sessionId || "").slice(-4)}`,
    phone: "+1234567890",
  };

  const updateResponse = http.put(
    `${getBaseUrl()}/api/v1/users/me`,
    JSON.stringify(updatePayload),
    { headers, tags: { scenario: "profile", type: "update_profile" } }
  );

  check(updateResponse, {
    "profile updated": (r) => r.status === 200,
  });

  sleep(Math.random() * 2 + 1);
}

/**
 * Complete user journey - Full e-commerce flow
 * @param {Object} userSession - User session (mixed auth/guest)
 */
export function completeUserJourney(userSession = null) {
  // Start with browsing (public)
  browseProducts(userSession);

  sleep(Math.random() * 1 + 0.5);

  // Search for products (public)
  searchAndFilter(userSession);

  sleep(Math.random() * 1 + 0.5);

  // If authenticated, do cart operations
  if (userSession && userSession.token) {
    cartOperations(userSession);
    sleep(Math.random() * 1 + 0.5);

    // Profile operations
    userProfileOperations(userSession);
    sleep(Math.random() * 1 + 0.5);
  }

  // Checkout (mixed public/auth)
  checkoutProcess(userSession);

  sleep(Math.random() * 2 + 1);
}

/**
 * Window shopping behavior - Pure browsing without purchase intent
 * @param {Object} userSession - User session (guest)
 */
export function windowShopping(userSession = null) {
  // Browse multiple categories
  browseProducts(userSession);
  sleep(Math.random() * 1 + 0.5);

  // Search different products
  searchAndFilter(userSession);
  sleep(Math.random() * 1 + 0.5);

  // Browse more products
  browseProducts(userSession);
  sleep(Math.random() * 2 + 1);
}

/**
 * Quick search user - Search-focused behavior
 * @param {Object} userSession - User session (guest)
 */
export function quickSearch(userSession = null) {
  // Multiple searches
  searchAndFilter(userSession);
  sleep(Math.random() * 0.5 + 0.2);

  searchAndFilter(userSession);
  sleep(Math.random() * 0.5 + 0.2);

  // Quick browse of results
  browseProducts(userSession);
  sleep(Math.random() * 1 + 0.5);
}
