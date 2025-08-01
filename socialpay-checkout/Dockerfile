FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

COPY package.json ./

RUN yarn

# Copy the rest of the app
COPY . .

# Build the Next.js app
RUN yarn build

# -------- Stage 2: Production --------
FROM node:22-alpine AS runner

# Set working directory
WORKDIR /app

# Copy necessary files from build stage
COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

# Expose the default port
EXPOSE 3000

# Start the app
CMD ["node", "server.js"]

