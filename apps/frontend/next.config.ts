import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Disable ESLint during builds for now
  eslint: {
    ignoreDuringBuilds: true,
  },
  
  // Disable type checking for faster builds during development
  typescript: {
    ignoreBuildErrors: true,
  },
  
  // Add any other configuration here
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: '**',
      },
    ],
  },
};

export default nextConfig;
