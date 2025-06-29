import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Ensures a URL has the correct protocol (http:// or https://)
 * Particularly useful for image URLs that might be missing the protocol
 */
export function ensureImageProtocol(url: string | undefined | null): string {
  if (!url) return "";

  // If URL already has a protocol, return as is
  if (url.startsWith("http://") || url.startsWith("https://")) {
    return url;
  }

  // Add protocol based on environment - use https in production, http in development
  const protocol =
    process.env.NODE_ENV === "production" ? "https://" : "http://";

  return protocol + url;
}
