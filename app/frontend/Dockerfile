# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install --legacy-peer-deps

# Copy source files
COPY . .

# Build arguments for Next.js public env vars
ARG NEXT_PUBLIC_GRAPHQL_URL=http://localhost:4000/graphql
ARG NEXT_PUBLIC_SOCKETIO_URL=http://localhost:3002
ARG NEXT_PUBLIC_ANALYTICS_ENABLED=true
ARG NEXT_PUBLIC_FEATURE_FLAGS_ENABLED=true

# Set environment variables for build
ENV NEXT_PUBLIC_GRAPHQL_URL=$NEXT_PUBLIC_GRAPHQL_URL
ENV NEXT_PUBLIC_SOCKETIO_URL=$NEXT_PUBLIC_SOCKETIO_URL
ENV NEXT_PUBLIC_ANALYTICS_ENABLED=$NEXT_PUBLIC_ANALYTICS_ENABLED
ENV NEXT_PUBLIC_FEATURE_FLAGS_ENABLED=$NEXT_PUBLIC_FEATURE_FLAGS_ENABLED

# Build the application with standalone output
RUN npm run build

# Production stage
FROM node:18-alpine AS runner

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001

# Set environment to production
ENV NODE_ENV=production

# Copy standalone build output
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
COPY --from=builder --chown=nextjs:nodejs /app/public ./public

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Health check - use the Next.js health endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

# Start Next.js server
CMD ["node", "server.js"]
