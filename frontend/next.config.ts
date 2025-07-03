import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: 'standalone',
  // Disable telemetry in production for better performance
  experimental: {
    optimizePackageImports: ['@heroicons/react', 'lucide-react']
  }
};

export default nextConfig;
