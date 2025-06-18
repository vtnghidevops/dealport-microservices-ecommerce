// Test file to simulate image URL formatting logic
// This helps debug the image URL issue without running the full application

console.log("=== Testing Image URL Formatting Logic ===\n");

// Mock environment variables
const mockEnv = {
  VITE_PUBLIC_BROKER_API_URL: "http://localhost:58080/api/v1",
};

// Simulate the formatImageUrl function from api-config.ts
function formatImageUrl(url) {
  console.log(`\n🔍 Testing URL: "${url}"`);

  // If the URL is empty, return empty string
  if (!url) {
    console.log("❌ Empty URL, returning empty string");
    return "";
  }

  // Fix common typo in protocol (htpp:// -> http://)
  if (url.startsWith("htpp://")) {
    url = url.replace("htpp://", "http://");
    console.log(`🔧 Fixed typo: ${url}`);
  }

  // If URL is a presigned URL (contains X-Amz parameters), return it as is
  if (
    url.includes("X-Amz-") ||
    (url.includes("?") && url.includes("Signature=")) ||
    url.includes("minioapi.deploy.io.vn") ||
    url.includes("minio.deploy.io.vn")
  ) {
    console.log("✅ Detected presigned URL, returning as-is");
    return url;
  }

  // If URL has any domain with /images/ path, extract just the /images/ path
  const anyDomainImagePattern = /https?:\/\/[^/]+(\/images\/.*)/;
  const domainMatch = url.match(anyDomainImagePattern);
  if (domainMatch && domainMatch[1]) {
    console.log(`🔧 Extracted relative path: ${domainMatch[1]}`);
    return domainMatch[1];
  }

  // If URL already has http:// or https://, it's already a complete URL
  if (url.startsWith("http://") || url.startsWith("https://")) {
    console.log("✅ Already complete URL, returning as-is");
    return url;
  }

  // Special handling for different image types
  if (url.startsWith("/images/")) {
    console.log("📂 Processing image URL...");

    // Check if this is a MinIO-stored image (special products-api folder)
    if (url.startsWith("/images/products-api/")) {
      // Get the base domain without /api/v1 if it exists
      let baseImageUrl = mockEnv.VITE_PUBLIC_BROKER_API_URL;
      if (baseImageUrl.endsWith("/api/v1")) {
        baseImageUrl = baseImageUrl.substring(
          0,
          baseImageUrl.length - "/api/v1".length
        );
      }

      const fullImageUrl = `${baseImageUrl}${url}`;
      console.log(
        `🗄️  MinIO image detected, requesting via broker: ${fullImageUrl}`
      );
      return fullImageUrl;
    }

    // For regular product images (/images/products/) and other /images/ URLs
    // Return relative URL without domain - browser will serve them directly
    console.log("🖼️  Regular image detected, returning relative URL");
    return url;
  }

  // Get the base domain without /api/v1 if it exists
  let baseImageUrl = mockEnv.VITE_PUBLIC_BROKER_API_URL;
  if (baseImageUrl.endsWith("/api/v1")) {
    baseImageUrl = baseImageUrl.substring(
      0,
      baseImageUrl.length - "/api/v1".length
    );
  }

  // For all other URLs starting with /, add the domain without /api/v1
  if (url.startsWith("/")) {
    const result = `${baseImageUrl}${url}`;
    console.log(`🔗 Adding domain to relative URL: ${result}`);
    return result;
  }

  // If URL doesn't start with slash, add one
  const result = `${baseImageUrl}/${url}`;
  console.log(`🔗 Adding domain and slash: ${result}`);
  return result;
}

// Simulate admin product service getAbsoluteUrl logic (the problematic one)
function getAbsoluteUrl_AdminService(url) {
  console.log(`\n🔍 [ADMIN SERVICE] Testing URL: "${url}"`);

  if (!url) return "";

  const brokerBaseUrl =
    mockEnv.VITE_PUBLIC_BROKER_API_URL?.replace("/api/v1", "") ||
    "http://localhost:8080";

  // Check if this is a MinIO URL
  const isMinioUrl =
    url.includes("minioapi.deploy.io.vn") ||
    url.includes("minio.deploy.io.vn") ||
    url.includes("X-Amz-");

  // If it's already an absolute URL
  if (url.startsWith("http://") || url.startsWith("https://")) {
    if (isMinioUrl) {
      console.log("✅ [ADMIN] Detected presigned URL, not modifying");
      return url;
    }

    // Extract the path portion if it contains /images/
    const anyDomainImagePattern = /https?:\/\/[^/]+(\/images\/.*)/;
    const domainMatch = url.match(anyDomainImagePattern);
    if (domainMatch && domainMatch[1]) {
      console.log(`🔧 [ADMIN] Extracted relative path: ${domainMatch[1]}`);
      return domainMatch[1];
    }

    return url;
  }

  // If URL is a blob URL
  if (url.startsWith("blob:")) {
    console.log("✅ [ADMIN] Blob URL, returning as-is");
    return url;
  }

  // If URL includes "localhost" with port
  if (url.includes("localhost:")) {
    const localhostNoProtocolPattern = /localhost:\d+(\/images\/.*)/;
    const localhostMatch = url.match(localhostNoProtocolPattern);
    if (localhostMatch && localhostMatch[1]) {
      console.log(
        `🔧 [ADMIN] Extracted from localhost URL: ${localhostMatch[1]}`
      );
      return localhostMatch[1];
    }

    return url.startsWith("//") ? `http:${url}` : `http://${url}`;
  }

  // Special case: If URL starts with /images/
  if (url.startsWith("/images/")) {
    if (url.startsWith("/images/products-api/")) {
      console.log("🗄️  [ADMIN] MinIO image path, preserving as-is");
      return url;
    }

    console.log("🖼️  [ADMIN] Local storage path, preserving as-is");
    return url;
  }

  // If it's a relative URL starting with '/', check if it's an image path first
  if (url.startsWith("/")) {
    // Don't add domain to image paths - they should be served relatively
    if (url.includes("/images/")) {
      console.log("📂 [ADMIN] Detected image path in relative URL, preserving");
      return url;
    }

    // For other relative URLs (API endpoints), add the broker service base URL
    const result = `${brokerBaseUrl}${url}`;
    console.log(
      `🔗 [ADMIN] Converting relative URL with broker base: ${result}`
    );
    return result;
  }

  // Otherwise, add the full path
  const result = `${brokerBaseUrl}/${url}`;
  console.log(`🔗 [ADMIN] Converting partial path: ${result}`);
  return result;
}

// Test cases
const testUrls = [
  // The problematic case reported by user
  "/images/products/electronics/computer-accessory-bundle-complete.jpg",

  // MinIO URLs
  "/images/products-api/minio-stored-image.jpg",

  // Already full URLs that should be converted to relative
  "http://localhost:58080/images/products/toys/rc-racing-car-angle.jpg",
  "http://localhost:58080/images/products-api/minio-file.jpg",

  // Presigned URLs (should be kept as-is)
  "https://minioapi.deploy.io.vn/bucket/image.jpg?X-Amz-Algorithm=...",

  // Other paths
  "/images/banners/sale.jpg",
  "/api/v1/products",
  "partial-path",

  // Edge cases
  "",
  "blob:http://localhost:3000/abc-123",
  "localhost:58080/images/products/test.jpg",
];

console.log("🧪 Testing with main api-config.ts formatImageUrl logic:");
console.log("=".repeat(80));

testUrls.forEach((url) => {
  const result = formatImageUrl(url);
  console.log(`➡️  Result: "${result}"`);
  console.log("-".repeat(60));
});

console.log("\n🧪 Testing with admin service getAbsoluteUrl logic:");
console.log("=".repeat(80));

testUrls.forEach((url) => {
  const result = getAbsoluteUrl_AdminService(url);
  console.log(`➡️  Result: "${result}"`);
  console.log("-".repeat(60));
});

console.log("\n📋 SUMMARY:");
console.log("=".repeat(40));
console.log("✅ Expected behavior:");
console.log("   - /images/products/* → /images/products/* (relative)");
console.log(
  "   - /images/products-api/* → http://localhost:58080/images/products-api/* (proxy to broker)"
);
console.log(
  "   - http://localhost:58080/images/* → /images/* (extract relative path)"
);
console.log("\n❌ Problematic behavior:");
console.log(
  "   - /images/products/* → http://localhost:58080/images/products/* (should be relative!)"
);

console.log("\n🎯 KEY FINDINGS:");
if (
  formatImageUrl(
    "/images/products/electronics/computer-accessory-bundle-complete.jpg"
  ).startsWith("http://")
) {
  console.log(
    "❌ Main formatImageUrl is adding domain to local images - THIS IS THE BUG!"
  );
} else {
  console.log("✅ Main formatImageUrl correctly keeps local images relative");
}

if (
  getAbsoluteUrl_AdminService(
    "/images/products/electronics/computer-accessory-bundle-complete.jpg"
  ).startsWith("http://")
) {
  console.log(
    "❌ Admin getAbsoluteUrl is adding domain to local images - THIS IS THE BUG!"
  );
} else {
  console.log("✅ Admin getAbsoluteUrl correctly keeps local images relative");
}
