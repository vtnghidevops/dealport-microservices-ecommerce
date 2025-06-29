# 🚀 Performance Testing Suite - Complete Overview

## 📋 **PROJECT SUMMARY**

This performance testing suite provides comprehensive load testing for the ecommerce microservices platform using **k6** testing framework. The suite includes 3 test types with full authentication integration and covers all implemented API endpoints in the broker service.

---

## 🎯 **TEST SUITE ARCHITECTURE**

### **Environment Configuration**

- **Local:** `http://localhost:58080`
- **Staging:** `https://api-sapogo.deploy.io.vn`
- **Production:** `https://api-sapogo.deploy.io.vn`

### **Authentication System**

- **OTP Bypass:** Code `123456` enabled for staging environment
- **Token Management:** Automatic user registration and verification
- **Shared Token Pool:** Multiple authenticated users for distributed load testing
- **Graceful Fallback:** Public endpoints when authentication fails

---

## 🔥 **TEST TYPES & COVERAGE**

### **1. SMOKE TEST (`smoke-test.js`)**

**Purpose:** Basic functionality validation and health checks  
**Duration:** ~40 seconds  
**Load:** 1-2 virtual users  
**Environment:** All environments

**📍 PUBLIC ENDPOINTS TESTED:**

```
✅ GET /api/v1/categories                 - Product categories listing
✅ GET /api/v1/products?limit=5          - Product catalog browsing
✅ POST /api/v1/checkout/validate        - Checkout validation (guest)
```

**🔐 AUTHENTICATED ENDPOINTS TESTED:**

```
✅ GET /api/v1/users/me                  - User profile retrieval
✅ GET /api/v1/cart                      - Shopping cart access
✅ POST /api/v1/cart/items               - Add items to cart
✅ GET /api/v1/users/me/wishlist         - User wishlist access
✅ GET /api/v1/checkout/orders           - Order history retrieval
✅ POST /api/v1/checkout/orders          - Order creation
```

**🎯 SUCCESS CRITERIA:**

- All public endpoints return 200 status
- All authenticated endpoints work with valid tokens
- Basic e-commerce functionality operational
- Authentication system functional with OTP bypass

---

### **2. LOAD TEST (`load-test.js`)**

**Purpose:** Realistic user behavior simulation under normal load  
**Duration:** ~16 minutes  
**Load:** Up to 300 concurrent users (staging)  
**Environment:** Staging, Production

**👥 USER BEHAVIOR DISTRIBUTION:**

```
30% - Window Shopping          : Browse products & categories
30% - Search & Browse          : Product discovery patterns
20% - Cart Operations          : Browse → Add to cart → Profile
15% - Complete Purchase Journey: Full checkout flow with auth
5%  - Account Management       : Profile, orders, wishlist management
```

**📊 LOAD PATTERN:**

```
Phase 1: Ramp up to 20% users (2m)
Phase 2: Ramp up to 50% users (3m)
Phase 3: Ramp up to 100% users (2m)
Phase 4: Sustain full load (10m)
Phase 5: Ramp down to 50% (3m)
Phase 6: Ramp down to 0% (2m)
```

**🎯 PERFORMANCE THRESHOLDS:**

- **Response Time:** 95% of requests < 3 seconds
- **Error Rate:** < 5% request failures
- **Check Success:** > 95% of validations pass
- **Throughput:** Sustained load without degradation

---

### **3. STRESS TEST (`stress-test.js`)**

**Purpose:** Find system breaking points and resilience testing  
**Duration:** ~45 minutes  
**Load:** Up to 750 concurrent users (250% capacity)  
**Environment:** Staging, Production

**🔥 STRESS LEVELS:**

```
🟢 NORMAL   (375 users)  : Regular patterns with authentication focus
🟡 STRESS   (500 users)  : Sustained heavy load (150% capacity)
🔴 EXTREME  (625 users)  : Heavy operations + intensive checkout (200%)
💥 SPIKE    (750 users)  : Rapid-fire requests, double operations (250%)
```

**⚡ STRESS SCENARIOS:**

- **Spike Phase:** Rapid authenticated operations, double requests
- **Extreme Phase:** Heavy browsing + cart + checkout with authentication
- **Stress Phase:** Sustained load with multiple cart operations
- **Normal Phase:** Regular patterns with authentication focus

**🎯 STRESS THRESHOLDS:**

- **Response Time:** 95% of requests < 8 seconds (degraded performance expected)
- **Error Rate:** < 15% request failures
- **Check Success:** > 85% of validations pass
- **Recovery:** System recovers when load decreases

---

## 🛠 **API ENDPOINT COVERAGE**

### **✅ FULLY IMPLEMENTED & TESTED:**

**Authentication Endpoints:**

```
POST /api/v1/auth/register               - User registration
POST /api/v1/auth/verify-registration    - OTP verification (bypass: 123456)
POST /api/v1/auth/login                  - User login
GET  /api/v1/auth/validate               - Token validation
```

**Product & Category Endpoints:**

```
GET /api/v1/categories                   - List all categories
GET /api/v1/categories/{id}              - Get category by ID
GET /api/v1/categories/slug/{slug}       - Get category by slug
GET /api/v1/products                     - List all products
GET /api/v1/products/{id}                - Get product by ID
GET /api/v1/products/slug/{slug}         - Get product by slug
```

**Cart Management Endpoints:**

```
GET    /api/v1/cart                      - Get user cart
POST   /api/v1/cart/items                - Add item to cart
PUT    /api/v1/cart/items/{item_id}      - Update cart item
DELETE /api/v1/cart/items/{item_id}      - Remove cart item
DELETE /api/v1/cart                      - Clear cart
POST   /api/v1/cart/coupon               - Apply coupon
DELETE /api/v1/cart/coupon               - Remove coupon
```

**User Management Endpoints:**

```
GET /api/v1/users/me                     - Get user profile
PUT /api/v1/users/me                     - Update user profile
GET /api/v1/users/me/wishlist            - Get user wishlist
POST /api/v1/users/me/wishlist           - Add to wishlist
DELETE /api/v1/users/me/wishlist/{id}    - Remove from wishlist
```

**Checkout & Order Endpoints:**

```
POST /api/v1/checkout/validate           - Validate checkout (public)
POST /api/v1/checkout/orders             - Create order (auth required)
GET  /api/v1/checkout/orders             - List user orders (auth required)
GET  /api/v1/checkout/orders/{id}        - Get specific order (auth required)
```

**Coupon Endpoints:**

```
GET /api/v1/coupons                      - List available coupons (public)
GET /api/v1/coupons/{id}                 - Get coupon by ID (public)
GET /api/v1/coupons/code/{code}          - Get coupon by code (public)
```

### **❌ INTENTIONALLY EXCLUDED (NOT IMPLEMENTED IN BROKER):**

```
❌ GET /api/v1/products?search=           - Product search not implemented
❌ GET /api/v1/products?min_price=        - Price filtering not implemented
❌ GET /api/v1/checkout/shipping-options  - Shipping options endpoint missing
❌ GET /api/v1/checkout/payment-methods   - Payment methods endpoint missing
```

---

## 🔧 **TECHNICAL IMPLEMENTATION**

### **File Structure:**

```
scripts/performance-tests/
├── config/
│   └── test-config.js              # Environment & load configurations
├── utils/
│   ├── auth-utils.js               # Authentication & token management
│   └── test-scenarios.js           # Reusable test scenarios
├── tests/
│   ├── smoke-test.js               # Basic functionality validation
│   ├── load-test.js                # Realistic user behavior simulation
│   └── stress-test.js              # Breaking point analysis
└── PERFORMANCE_TESTING_OVERVIEW.md # This documentation
```

### **Key Technical Features:**

**Authentication Integration:**

- Automatic user registration with unique email generation
- OTP bypass code `123456` for staging environment
- Token pool management for distributed load testing
- Graceful fallback to public endpoints when auth fails

**Load Distribution:**

- Environment-specific configurations (local/staging/production)
- Realistic user behavior patterns based on e-commerce analytics
- Progressive load ramping with sustained periods
- Stress testing with multiple load levels

**Error Handling:**

- Comprehensive HTTP status code validation
- Graceful degradation under high load
- Detailed logging and error reporting
- Fallback mechanisms for failed operations

---

## 🚀 **EXECUTION INSTRUCTIONS**

### **Local Testing:**

```bash
cd scripts/performance-tests

# Smoke test
k6 run tests/smoke-test.js

# Load test
k6 run tests/load-test.js

# Stress test
k6 run tests/stress-test.js
```

### **CI/CD Pipeline:**

The tests are integrated into GitLab CI pipeline with:

- Environment variables: `ENVIRONMENT=staging`, `TARGET_URL=${STAGING_URL}`
- Automatic URL conversion from frontend to API URLs
- Performance results artifacts collection
- Threshold validation and pipeline failure on performance regression

### **Environment Variables:**

```bash
ENVIRONMENT=staging          # Test environment (local/staging/production)
TARGET_URL=https://api-...   # API base URL (optional, uses config if not set)
```

---

## 📊 **PERFORMANCE METRICS & THRESHOLDS**

### **Staging Environment Thresholds:**

```javascript
Smoke Test:
  - Response Time: 95% < 2s
  - Error Rate: < 10%
  - Check Success: > 90%

Load Test:
  - Response Time: 95% < 3s
  - Error Rate: < 5%
  - Check Success: > 95%

Stress Test:
  - Response Time: 95% < 8s
  - Error Rate: < 15%
  - Check Success: > 85%
```

### **Key Performance Indicators:**

- **Throughput:** Requests per second sustained
- **Response Time:** P95 latency across all endpoints
- **Error Rate:** Percentage of failed requests
- **Concurrent Users:** Maximum supported simultaneous users
- **Resource Utilization:** CPU, memory, database connections

---

## 🎯 **SUCCESS CRITERIA**

### **Functional Requirements:**

✅ All implemented API endpoints respond correctly  
✅ Authentication system works with OTP bypass  
✅ Cart operations function under load  
✅ Checkout process completes successfully  
✅ User profile management operates normally

### **Performance Requirements:**

✅ System handles expected user load (300 concurrent users)  
✅ Response times remain acceptable under normal load  
✅ System gracefully degrades under stress conditions  
✅ Recovery occurs when load returns to normal levels  
✅ No memory leaks or resource exhaustion detected

### **Reliability Requirements:**

✅ Error rates stay within acceptable limits  
✅ No complete system failures during testing  
✅ Database connections remain stable  
✅ Auto-scaling triggers appropriately in K8s  
✅ Circuit breakers activate under extreme load

---

## 🔍 **MONITORING & ANALYSIS**

### **Metrics to Monitor:**

- **Application Performance:** Response times, throughput, error rates
- **Infrastructure:** CPU, memory, disk I/O, network utilization
- **Database:** Connection pool usage, query performance, deadlocks
- **Kubernetes:** Pod scaling, resource limits, restarts
- **External Services:** Payment gateway, email service availability

### **Red Flags to Watch:**

🚨 **Critical Issues:**

- Complete system unresponsiveness
- Error rates above 25%
- Memory leaks or resource exhaustion
- Database deadlocks or connection exhaustion
- Cascading failures across services

⚠️ **Warning Signs:**

- Response times increasing significantly
- Error rates above threshold but below critical
- Resource utilization approaching limits
- Intermittent service timeouts
- Auto-scaling not triggering as expected

---

## 🛡 **SECURITY CONSIDERATIONS**

### **Authentication Security:**

- OTP bypass only enabled in staging environment
- Production uses full OTP verification workflow
- Token management with proper expiration
- Rate limiting on authentication endpoints

### **Test Data Security:**

- No real user data used in performance tests
- Generated test accounts with fake information
- Cleanup of test data after execution
- Secure handling of authentication tokens

---

## 📈 **CONTINUOUS IMPROVEMENT**

### **Regular Review Items:**

- Update user behavior patterns based on production analytics
- Adjust load thresholds as system capacity grows
- Add new endpoints as they are implemented
- Refine stress scenarios based on real-world incidents
- Update authentication mechanisms as they evolve

### **Performance Benchmarking:**

- Establish baseline performance metrics
- Track performance trends over time
- Compare results across environment types
- Identify performance regressions early
- Validate capacity planning assumptions

---

## 🔗 **RELATED DOCUMENTATION**

- **API Documentation:** Broker service routing configuration
- **Authentication Guide:** OTP bypass setup for staging
- **Deployment Guide:** K8s configuration and scaling policies
- **Monitoring Setup:** Performance metrics collection and alerting
- **Incident Response:** Performance degradation troubleshooting

---

## 📞 **SUPPORT & MAINTENANCE**

### **Test Maintenance:**

- Review and update quarterly
- Monitor for API changes and update endpoints
- Validate OTP bypass functionality in staging
- Update load patterns based on production usage
- Refresh test data and scenarios regularly

### **Issue Resolution:**

- Performance test failures: Check staging environment status
- Authentication issues: Verify OTP bypass configuration
- Load generation problems: Validate k6 infrastructure
- Threshold violations: Analyze system performance metrics
- Pipeline integration: Check GitLab CI configuration

---

**Last Updated:** January 2025  
**Version:** 1.0  
**Environment:** Staging/Production Ready  
**Test Coverage:** 100% of implemented broker service endpoints
