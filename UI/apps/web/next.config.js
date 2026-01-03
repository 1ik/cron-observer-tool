/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  transpilePackages: ['@cron-observer/ui', '@cron-observer/lib'],
}

module.exports = nextConfig

