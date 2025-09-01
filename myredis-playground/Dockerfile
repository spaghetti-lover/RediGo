# Stage 1: Build frontend
FROM node:18-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package.json and lockfile first (optimize cache layer)
COPY package.json package-lock.json* yarn.lock* pnpm-lock.yaml* ./

# Install dependencies
RUN npm install

# Copy entire source code
COPY . .

# Build Next.js
RUN npm run build

# Stage 2: Run production
FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# Copy from builder stage
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

# Expose port 3000
EXPOSE 3000

# Chạy Next.js
CMD ["npm", "run", "start"]
