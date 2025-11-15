/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Removed 'output: export' to enable SSR and dynamic features
  output: 'standalone', // Use standalone build for Docker deployment
  images: {
    unoptimized: true, // Keep unoptimized for Docker compatibility
  },
  env: {
    NEXT_PUBLIC_GRAPHQL_URL: process.env.NEXT_PUBLIC_GRAPHQL_URL,
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL,
  },
  typescript: {
    ignoreBuildErrors: false,
  },
  eslint: {
    ignoreDuringBuilds: false,
  },
  // Disable trailing slashes
  trailingSlash: false,
};

module.exports = nextConfig;
