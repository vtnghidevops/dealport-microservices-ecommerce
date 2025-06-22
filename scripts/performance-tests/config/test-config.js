// Test configuration for different environments
export const config = {
  // Base URLs for different environments
  baseUrls: {
    local: "http://localhost:58080",
    staging: "https://staging-api-sapogo.deploy.io.vn",
    production: "https://api-sapogo.deploy.io.vn",
  },

  // Get current environment (default to local)
  environment: __ENV.ENVIRONMENT || "local",

  // Test data
  testUsers: {
    admin: {
      email: "admin@test.com",
      password: "admin123456",
    },
    customer: {
      email: "customer@test.com",
      password: "customer123456",
    },
    // Shared test account pool for performance testing
    performanceTestAccounts: [
      {
        email: "perf-user-1@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User1",
      },
      {
        email: "perf-user-2@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User2",
      },
      {
        email: "perf-user-3@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User3",
      },
      {
        email: "perf-user-4@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User4",
      },
      {
        email: "perf-user-5@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User5",
      },
      {
        email: "perf-user-6@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User6",
      },
      {
        email: "perf-user-7@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User7",
      },
      {
        email: "perf-user-8@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User8",
      },
      {
        email: "perf-user-9@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User9",
      },
      {
        email: "perf-user-10@sapogo.test",
        password: "PerfTest123!",
        firstName: "Performance",
        lastName: "User10",
      },
    ],
  },

  // Test products for cart operations
  testProducts: [
    { id: 1, quantity: 2 },
    { id: 2, quantity: 1 },
    { id: 3, quantity: 3 },
    { id: 4, quantity: 5 },
    { id: 5, quantity: 1 },
  ],

  // Performance thresholds - updated for K8s deployment
  thresholds: {
    smoke: {
      http_req_duration: ["p(95)<2000"], // 95% of requests under 2s
      http_req_failed: ["rate<0.1"], // Error rate under 10%
      checks: ["rate>0.9"], // 90% of checks pass
    },
    load: {
      http_req_duration: ["p(95)<3000"], // 95% of requests under 3s
      http_req_failed: ["rate<0.05"], // Error rate under 5%
      checks: ["rate>0.95"], // 95% of checks pass
    },
    stress: {
      http_req_duration: ["p(95)<5000"], // 95% of requests under 5s
      http_req_failed: ["rate<0.1"], // Error rate under 10%
      checks: ["rate>0.9"], // 90% of checks pass
    },
    // New: High-load thresholds for K8s staging
    staging: {
      http_req_duration: ["p(95)<5000"], // More reasonable 5s threshold
      http_req_failed: ["rate<0.10"], // 10% error rate - still reasonable for load testing
      checks: ["rate>0.85"], // 85% success rate - achievable
    },
    // New: Production-ready thresholds
    production: {
      http_req_duration: ["p(95)<2000"], // 95% of requests under 2s
      http_req_failed: ["rate<0.02"], // Error rate under 2%
      checks: ["rate>0.98"], // 98% of checks pass
    },
  },

  // Load patterns for different environments
  loadPatterns: {
    local: {
      maxUsers: 50,
      rampUpTime: "2m",
      sustainTime: "5m",
    },
    staging: {
      maxUsers: 300,
      stressMaxUsers: 500,
      rampUpTime: "3m",
      sustainTime: "10m",
    },
    production: {
      maxUsers: 500, // Production-level load
      rampUpTime: "5m",
      sustainTime: "15m",
    },
  },

  // K8s specific configuration
  k8s: {
    nodeCount: 5,
    expectedThroughput: {
      staging: 1000, // requests per second
      production: 2000,
    },
    resourceLimits: {
      cpu: "2000m", // 2 CPU cores per pod
      memory: "4Gi", // 4GB RAM per pod
    },
  },

  // Validation and monitoring configuration
  validation: {
    // Minimum expected metrics for test validity
    minimumMetrics: {
      smoke: {
        duration: 30, // seconds
        totalRequests: 10,
        virtualUsers: 1,
      },
      load: {
        duration: 900, // 15 minutes
        totalRequests: 3000,
        virtualUsers: 200,
      },
      stress: {
        duration: 1500, // 25 minutes
        totalRequests: 8000,
        virtualUsers: 400,
      },
    },

    // Health check endpoints for validation
    healthChecks: {
      api: "/api/v1/health",
      auth: "/api/v1/auth/health",
      products: "/api/v1/products/health",
      cart: "/api/v1/cart/health",
    },

    // Expected status codes for different operations
    expectedStatusCodes: {
      browse: [200],
      search: [200],
      auth: [200, 201],
      cart: [200, 201],
      checkout: [200, 422], // 422 for validation errors is acceptable
    },

    // Test execution monitoring
    monitoring: {
      logInterval: 30, // seconds
      metricsCollection: true,
      realTimeValidation: true,
      failureThreshold: 0.25, // 25% failure rate stops test
    },
  },

  // Environment-specific validation rules
  environmentValidation: {
    staging: {
      requiredFeatures: ["otp_bypass", "performance_testing"],
      maxResponseTime: 5000, // ms
      minThroughput: 100, // req/s
    },
    production: {
      requiredFeatures: ["security", "monitoring"],
      maxResponseTime: 2000, // ms
      minThroughput: 500, // req/s
    },
  },
};

// Get base URL for current environment
export function getBaseUrl() {
  // Always use API URLs, ignore TARGET_URL from pipeline
  const envFromK6 = __ENV.ENVIRONMENT || config.environment;
  const baseUrl = config.baseUrls[envFromK6] || config.baseUrls.local;

  // Environment: ${envFromK6}, Base URL: ${baseUrl}
  return baseUrl;
}

// Get load pattern for current environment
export function getLoadPattern() {
  return config.loadPatterns[config.environment] || config.loadPatterns.local;
}

// Get thresholds for current environment
export function getThresholds() {
  return config.thresholds[config.environment] || config.thresholds.local;
}

// Common headers
export const headers = {
  "Content-Type": "application/json",
  Accept: "application/json",
};
