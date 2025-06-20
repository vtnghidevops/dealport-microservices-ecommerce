// Test configuration for different environments
export const config = {
  // Base URLs for different environments
  baseUrls: {
    local: "http://localhost:58080",
    staging: "https://sapogo.deploy.io.vn",
    production: "https://sapogo.deploy.io.vn",
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
      http_req_duration: ["p(95)<2500"], // 95% of requests under 2.5s
      http_req_failed: ["rate<0.03"], // Error rate under 3%
      checks: ["rate>0.97"], // 97% of checks pass
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
      maxUsers: 200, // Increased for K8s
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
};

// Get base URL for current environment
export function getBaseUrl() {
  // Priority 1: Use TARGET_URL from environment (for CI/CD pipelines)
  if (__ENV.TARGET_URL) {
    return __ENV.TARGET_URL;
  }

  // Priority 2: Use configured URL for environment
  return config.baseUrls[config.environment] || config.baseUrls.local;
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
