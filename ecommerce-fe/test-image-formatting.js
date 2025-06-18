// Test script to verify formatImageUrl logic
const BASE_API_URL = "http://localhost:58080/api/v1";
const API_VERSION = "v1";

const formatImageUrl = (url) => {
  // If the URL is empty, return empty string
  if (!url) return "";

  // Fix common typo in protocol (htpp:// -> http://)
  if (url.startsWith("htpp://")) {
    url = url.replace("htpp://", "http://");
  }

  // If URL is a presigned URL (contains X-Amz parameters), return it as is
  if (
    url.includes("X-Amz-") ||
    (url.includes("?") && url.includes("Signature=")) ||
    url.includes("minioapi.deploy.io.vn") ||
    url.includes("minio.deploy.io.vn")
  ) {
    console.log("Detected presigned URL:", url.split("?")[0]);
    return url;
  }

  // If URL has any domain with /images/ path, extract just the /images/ path
  const anyDomainImagePattern = /https?:\/\/[^/]+(\/images\/.*)/;
  const domainMatch = url.match(anyDomainImagePattern);
  if (domainMatch && domainMatch[1]) {
    return domainMatch[1];
  }

  // If URL already has http:// or https://, it's already a complete URL
  if (url.startsWith("http://") || url.startsWith("https://")) {
    return url;
  }

  // Special handling for different image types
  if (url.startsWith("/images/")) {
    console.log("Processing image URL:", url);

    // Get the base domain without /api/v1 if it exists
    let baseImageUrl = BASE_API_URL;
    if (baseImageUrl.endsWith(`/api/${API_VERSION}`)) {
      baseImageUrl = baseImageUrl.substring(
        0,
        baseImageUrl.length - `/api/${API_VERSION}`.length
      );
    }

    // Always add domain for ALL /images/ URLs for consistency
    const fullImageUrl = `${baseImageUrl}${url}`;

    // Check if this is a MinIO-stored image (special products-api folder)
    if (url.startsWith("/images/products-api/")) {
      console.log("MinIO image detected, requesting via broker:", fullImageUrl);
      return fullImageUrl;
    }

    // For regular product images (/images/products/) and other /images/ URLs
    // Return full domain URL like https://sapogo.deploy.io.vn/images/...
    console.log(
      "Regular image detected, returning full domain URL:",
      fullImageUrl
    );
    return fullImageUrl;
  }

  // Get the base domain without /api/v1 if it exists
  let baseImageUrl = BASE_API_URL;
  if (baseImageUrl.endsWith(`/api/${API_VERSION}`)) {
    baseImageUrl = baseImageUrl.substring(
      0,
      baseImageUrl.length - `/api/${API_VERSION}`.length
    );
  }

  // For all other URLs starting with /, add the domain without /api/v1
  if (url.startsWith("/")) {
    return `${baseImageUrl}${url}`;
  }

  // If URL doesn't start with slash, add one
  return `${baseImageUrl}/${url}`;
};

// Test cases
console.log("=== Testing formatImageUrl ===");

// Test 1: Regular product images (should keep relative)
console.log("\n1. Regular product images:");
console.log("Input: /images/products/laptop.jpg");
console.log("Output:", formatImageUrl("/images/products/laptop.jpg"));

// Test 2: MinIO images (should add domain)
console.log("\n2. MinIO images:");
console.log("Input: /images/products-api/phone.jpg");
console.log("Output:", formatImageUrl("/images/products-api/phone.jpg"));

// Test 3: Other images
console.log("\n3. Other images:");
console.log("Input: /images/banners/sale.jpg");
console.log("Output:", formatImageUrl("/images/banners/sale.jpg"));

// Test 4: Already complete URLs
console.log("\n4. Complete URLs:");
console.log("Input: https://example.com/images/products/test.jpg");
console.log(
  "Output:",
  formatImageUrl("https://example.com/images/products/test.jpg")
);

// Test 5: MinIO presigned URLs
console.log("\n5. MinIO presigned URLs:");
console.log(
  "Input: https://minioapi.deploy.io.vn/bucket/image.jpg?X-Amz-Signature=abc"
);
console.log(
  "Output:",
  formatImageUrl(
    "https://minioapi.deploy.io.vn/bucket/image.jpg?X-Amz-Signature=abc"
  )
);

console.log("\n=== Expected Results ===");
console.log(
  "1. Regular product images: http://localhost:58080/images/products/laptop.jpg (with domain)"
);
console.log(
  "2. MinIO images: http://localhost:58080/images/products-api/phone.jpg (with domain)"
);
console.log(
  "3. Other images: http://localhost:58080/images/banners/sale.jpg (with domain)"
);
console.log("4. Complete URLs: /images/products/test.jpg (extracted path)");
console.log("5. MinIO presigned: unchanged (keep presigned URL)");
