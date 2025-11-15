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

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine

WORKDIR /app

# Install serve globally for static file serving
RUN npm install -g serve

# Copy built static files from builder
COPY --from=builder /app/out ./out
COPY --from=builder /app/public ./public

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001

# Change ownership
RUN chown -R nextjs:nodejs /app

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Set environment to production
ENV NODE_ENV=production

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

# Start the application using serve for static files
CMD ["serve", "-s", "out", "-l", "3000"]
