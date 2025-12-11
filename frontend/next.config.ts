import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  webpack: (config) => {
    config.resolve.alias['@'] = process.cwd();
    return config;
  },
};

export default nextConfig;
