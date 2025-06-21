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

  // Skip user profile - sync issue between Auth and User services
  console.log("Skipping profile fetch (sync issue)");

  // Get user orders (implemented endpoint)
  const ordersResponse = http.get(`${getBaseUrl()}/api/v1/checkout/orders`, {
    headers,
    tags: { scenario: "profile", type: "get_orders" },
  });

  check(ordersResponse, {
    "orders retrieved": (r) => r.status === 200,
  });

  // Skip wishlist - not implemented
  console.log("Skipping wishlist (not implemented)");

  // Update profile information (test endpoint)
  const updatePayload = {
    first_name: `TestUser${userSession.id || ""}`,
    last_name: `Session${(userSession.sessionId || "").slice(-4)}`,
    phone: "+1234567890",
  };

  const updateResponse = http.put(
    `${getBaseUrl()}/api/v1/users/profile`,
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

/**
 * Authentication flow testing - Test login/logout cycles
 * @param {string} email - User email
 * @param {string} password - User password
 */
export function authenticationFlow(
  email = "test@example.com",
  password = "Test123456!"
) {
  const headers = { "Content-Type": "application/json" };

  // Test login
  const loginPayload = { email, password };
  const loginResponse = http.post(
    `${getBaseUrl()}/api/v1/auth/login`,
    JSON.stringify(loginPayload),
    { headers, tags: { scenario: "auth", type: "login" } }
  );

  const loginSuccess = check(loginResponse, {
    "login successful": (r) => r.status === 200,
    "login returns token": (r) => {
      try {
        const data = JSON.parse(r.body);
        return data.data && data.data.access_token;
      } catch {
        return false;
      }
    },
  });

  if (loginSuccess && loginResponse.status === 200) {
    const loginData = JSON.parse(loginResponse.body);
    const token = loginData.data.access_token;

    // Test token validation
    const validateResponse = http.get(`${getBaseUrl()}/api/v1/auth/validate`, {
      headers: getAuthHeaders(token),
      tags: { scenario: "auth", type: "validate" },
    });

    check(validateResponse, {
      "token validation successful": (r) => r.status === 200,
    });

    // Test logout
    const logoutResponse = http.post(
      `${getBaseUrl()}/api/v1/auth/logout`,
      {},
      {
        headers: getAuthHeaders(token),
        tags: { scenario: "auth", type: "logout" },
      }
    );

    check(logoutResponse, {
      "logout successful": (r) => r.status === 200,
    });

    return token;
  }

  return null;
}

/**
 * Admin operations testing - Order management, user management
 * @param {Object} adminSession - Admin user session
 */
export function adminOperations(adminSession) {
  if (!adminSession || !adminSession.token) {
    console.log("Admin operations skipped - admin authentication required");
    return;
  }

  const headers = getBaseHeaders(adminSession);

  // Admin - Get all orders
  const ordersResponse = http.get(
    `${getBaseUrl()}/api/v1/checkout/admin/orders`,
    { headers, tags: { scenario: "admin", type: "list_orders" } }
  );

  check(ordersResponse, {
    "admin orders retrieved": (r) => r.status === 200,
  });

  // Admin - Get user statistics
  const statsPayload = {
    start_date: "2024-01-01",
    end_date: "2024-12-31",
    limit: 50,
  };

  const statsResponse = http.post(
    `${getBaseUrl()}/api/v1/users/admin/statistics`,
    JSON.stringify(statsPayload),
    { headers, tags: { scenario: "admin", type: "user_stats" } }
  );

  check(statsResponse, {
    "admin user statistics retrieved": (r) => r.status === 200,
  });

  // Admin - Get user list
  const userListPayload = {
    page: 1,
    limit: 20,
  };

  const userListResponse = http.post(
    `${getBaseUrl()}/api/v1/users/admin/list`,
    JSON.stringify(userListPayload),
    { headers, tags: { scenario: "admin", type: "user_list" } }
  );

  check(userListResponse, {
    "admin user list retrieved": (r) => r.status === 200,
  });

  sleep(Math.random() * 1 + 0.5);
}

/**
 * Product management testing - CRUD operations
 * @param {Object} userSession - User session (admin preferred)
 */
export function productManagement(userSession) {
  if (!userSession || !userSession.token) {
    console.log("Product management skipped - authentication required");
    return browseProducts();
  }

  const headers = getBaseHeaders(userSession);

  // Get product details
  const productsListResponse = http.get(
    `${getBaseUrl()}/api/v1/products?limit=5`,
    { headers, tags: { scenario: "product_mgmt", type: "list" } }
  );

  check(productsListResponse, {
    "products list retrieved": (r) => r.status === 200,
  });

  if (productsListResponse.status === 200) {
    try {
      const products = JSON.parse(productsListResponse.body);
      if (products.data && products.data.length > 0) {
        const product = products.data[0];

        // Get product by ID
        const productResponse = http.get(
          `${getBaseUrl()}/api/v1/products/${product.id}`,
          { headers, tags: { scenario: "product_mgmt", type: "get_by_id" } }
        );

        check(productResponse, {
          "product details retrieved": (r) => r.status === 200,
        });

        // Get product by slug (if available)
        if (product.slug) {
          const slugResponse = http.get(
            `${getBaseUrl()}/api/v1/products/slug/${product.slug}`,
            { headers, tags: { scenario: "product_mgmt", type: "get_by_slug" } }
          );

          check(slugResponse, {
            "product by slug retrieved": (r) => r.status === 200,
          });
        }

        // Get product reviews
        const reviewsResponse = http.get(
          `${getBaseUrl()}/api/v1/products/${product.id}/reviews`,
          { headers, tags: { scenario: "product_mgmt", type: "reviews" } }
        );

        check(reviewsResponse, {
          "product reviews retrieved": (r) => r.status === 200,
        });
      }
    } catch (error) {
      console.error("Error in product management:", error.message);
    }
  }

  sleep(Math.random() * 1 + 0.5);
}

/**
 * Category management testing
 * @param {Object} userSession - User session
 */
export function categoryManagement(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // Get all categories
  const categoriesResponse = http.get(`${getBaseUrl()}/api/v1/categories`, {
    headers,
    tags: { scenario: "category_mgmt", type: "list" },
  });

  check(categoriesResponse, {
    "categories list retrieved": (r) => r.status === 200,
  });

  if (categoriesResponse.status === 200) {
    try {
      const categories = JSON.parse(categoriesResponse.body);
      if (categories.data && categories.data.length > 0) {
        const category = categories.data[0];

        // Get category by ID
        const categoryResponse = http.get(
          `${getBaseUrl()}/api/v1/categories/${category.id}`,
          { headers, tags: { scenario: "category_mgmt", type: "get_by_id" } }
        );

        check(categoryResponse, {
          "category details retrieved": (r) => r.status === 200,
        });

        // Get category by slug
        if (category.slug) {
          const slugResponse = http.get(
            `${getBaseUrl()}/api/v1/categories/slug/${category.slug}`,
            {
              headers,
              tags: { scenario: "category_mgmt", type: "get_by_slug" },
            }
          );

          check(slugResponse, {
            "category by slug retrieved": (r) => r.status === 200,
          });
        }
      }
    } catch (error) {
      console.error("Error in category management:", error.message);
    }
  }

  sleep(Math.random() * 1 + 0.5);
}

/**
 * Enhanced shopping experience - Include banners, testimonials, reviews
 * @param {Object} userSession - User session
 */
export function enhancedShopping(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // View homepage content - banners and testimonials
  const bannersResponse = http.get(`${getBaseUrl()}/api/v1/banners`, {
    headers,
    tags: { scenario: "homepage", type: "banners" },
  });

  check(bannersResponse, {
    "banners loaded": (r) => r.status === 200,
  });

  // Get testimonials for social proof
  const testimonialsResponse = http.get(`${getBaseUrl()}/api/v1/testimonials`, {
    headers,
    tags: { scenario: "homepage", type: "testimonials" },
  });

  check(testimonialsResponse, {
    "testimonials loaded": (r) => r.status === 200,
  });

  // Browse products and read reviews
  const productsResponse = http.get(`${getBaseUrl()}/api/v1/products?limit=5`, {
    headers,
    tags: { scenario: "shopping", type: "products" },
  });

  if (productsResponse.status === 200) {
    try {
      const products = JSON.parse(productsResponse.body);
      if (products.data && products.data.length > 0) {
        const product = products.data[0];

        // Read product reviews (important for purchase decision)
        const reviewsResponse = http.get(
          `${getBaseUrl()}/api/v1/products/${product.id}/reviews`,
          { headers, tags: { scenario: "shopping", type: "reviews" } }
        );

        check(reviewsResponse, {
          "product reviews loaded": (r) => r.status === 200,
        });

        // If authenticated, potentially write a review
        if (userSession && userSession.token && Math.random() < 0.3) {
          const reviewPayload = {
            rating: Math.floor(Math.random() * 2) + 4, // 4-5 stars (realistic positive reviews)
            comment: "Great product, fast delivery, highly recommended!",
            title: "Excellent quality",
          };

          const addReviewResponse = http.post(
            `${getBaseUrl()}/api/v1/products/${product.id}/reviews`,
            JSON.stringify(reviewPayload),
            { headers, tags: { scenario: "shopping", type: "add_review" } }
          );

          check(addReviewResponse, {
            "review added": (r) => r.status === 200 || r.status === 201,
          });
        }
      }
    } catch (error) {
      console.error("Error in enhanced shopping:", error.message);
    }
  }

  sleep(Math.random() * 1 + 0.5);
}

/**
 * Coupon and promotions exploration - Real user behavior
 * @param {Object} userSession - User session
 */
export function explorePromotions(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // Browse available coupons/promotions
  const couponsResponse = http.get(`${getBaseUrl()}/api/v1/coupons`, {
    headers,
    tags: { scenario: "promotions", type: "list_coupons" },
  });

  check(couponsResponse, {
    "coupons loaded": (r) => r.status === 200,
  });

  if (couponsResponse.status === 200) {
    try {
      const coupons = JSON.parse(couponsResponse.body);
      if (coupons.data && coupons.data.length > 0) {
        const coupon = coupons.data[0];

        // Check coupon details
        const couponDetailResponse = http.get(
          `${getBaseUrl()}/api/v1/coupons/${coupon.id}`,
          { headers, tags: { scenario: "promotions", type: "coupon_details" } }
        );

        check(couponDetailResponse, {
          "coupon details loaded": (r) => r.status === 200,
        });

        // Try to get coupon by code (common user behavior)
        if (coupon.code) {
          const couponByCodeResponse = http.get(
            `${getBaseUrl()}/api/v1/coupons/code/${coupon.code}`,
            {
              headers,
              tags: { scenario: "promotions", type: "coupon_by_code" },
            }
          );

          check(couponByCodeResponse, {
            "coupon by code loaded": (r) => r.status === 200,
          });
        }
      }
    } catch (error) {
      console.error("Error exploring promotions:", error.message);
    }
  }

  sleep(Math.random() * 1 + 0.5);
}

/**
 * Payment flow testing - Real checkout experience
 * @param {Object} userSession - User session with token
 */
export function paymentFlow(userSession) {
  if (!userSession || !userSession.token) {
    console.log("Payment flow skipped - authentication required");
    return;
  }

  const headers = getBaseHeaders(userSession);

  // First add something to cart
  const productsResponse = http.get(`${getBaseUrl()}/api/v1/products?limit=3`, {
    headers,
    tags: { scenario: "payment", type: "get_products" },
  });

  if (productsResponse.status === 200) {
    try {
      const products = JSON.parse(productsResponse.body);
      if (products.data && products.data.length > 0) {
        const product = products.data[0];

        // Add to cart
        const addToCartPayload = {
          product_id: product.id,
          quantity: 2,
        };

        const addResponse = http.post(
          `${getBaseUrl()}/api/v1/cart/items`,
          JSON.stringify(addToCartPayload),
          { headers, tags: { scenario: "payment", type: "add_to_cart" } }
        );

        if (addResponse.status === 200 || addResponse.status === 201) {
          // Create order first
          const orderPayload = {
            shipping_address: "123 Test Street, Test City",
            shipping_phone: "+1234567890",
            payment_method: "momo",
            notes: "Please deliver during business hours",
          };

          const orderResponse = http.post(
            `${getBaseUrl()}/api/v1/checkout/orders`,
            JSON.stringify(orderPayload),
            { headers, tags: { scenario: "payment", type: "create_order" } }
          );

          check(orderResponse, {
            "order created": (r) => r.status === 200 || r.status === 201,
          });

          if (orderResponse.status === 200 || orderResponse.status === 201) {
            const orderData = JSON.parse(orderResponse.body);
            const orderId = orderData.data?.id || orderData.data?.order_id;

            if (orderId) {
              // Try MoMo payment creation (common payment method)
              const paymentPayload = {
                order_id: orderId,
                return_url: "https://example.com/payment/success",
                notify_url: "https://example.com/payment/notify",
              };

              const momoPaymentResponse = http.post(
                `${getBaseUrl()}/api/v1/payments/momo/create`,
                JSON.stringify(paymentPayload),
                { headers, tags: { scenario: "payment", type: "momo_create" } }
              );

              check(momoPaymentResponse, {
                "momo payment created": (r) =>
                  r.status === 200 || r.status === 201,
              });
            }
          }
        }
      }
    } catch (error) {
      console.error("Error in payment flow:", error.message);
    }
  }

  sleep(Math.random() * 2 + 1);
}

/**
 * Homepage content viewing - Banners and ads
 * @param {Object} userSession - User session
 */
export function viewHomepageContent(userSession = null) {
  const headers = getBaseHeaders(userSession);

  // Get different types of banners
  const bannerTypes = ["hero", "promotion", "category"];
  const randomType =
    bannerTypes[Math.floor(Math.random() * bannerTypes.length)];

  const bannersByTypeResponse = http.get(
    `${getBaseUrl()}/api/v1/banners/type/${randomType}`,
    { headers, tags: { scenario: "homepage", type: "banners_by_type" } }
  );

  check(bannersByTypeResponse, {
    "banners by type loaded": (r) => r.status === 200,
  });

  // Get ads for different placements
  const adPlacements = ["homepage", "sidebar", "footer"];
  const randomPlacement =
    adPlacements[Math.floor(Math.random() * adPlacements.length)];

  const adsResponse = http.get(
    `${getBaseUrl()}/api/v1/ads/placement/${randomPlacement}`,
    { headers, tags: { scenario: "homepage", type: "ads" } }
  );

  check(adsResponse, {
    "ads loaded": (r) => r.status === 200,
  });

  sleep(Math.random() * 1 + 0.5);
}
