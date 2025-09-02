/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "export",
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
  devIndicators: false,
  eslint: {
    ignoreDuringBuilds: true,
  },
};

module.exports = nextConfig;
