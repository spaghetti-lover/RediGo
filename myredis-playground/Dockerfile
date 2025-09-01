# Stage 1: Build frontend
FROM node:18-alpine AS builder

# Đặt thư mục làm việc
WORKDIR /app

# Copy package.json và lockfile trước (tối ưu cache layer)
COPY package.json package-lock.json* yarn.lock* pnpm-lock.yaml* ./

# Cài dependencies
RUN npm install

# Copy toàn bộ source code
COPY . .

# Build Next.js
RUN npm run build

# Stage 2: Run production
FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production

# Copy từ stage builder
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

# Expose port 3000
EXPOSE 3000

# Chạy Next.js
CMD ["npm", "run", "start"]
