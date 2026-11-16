
# Build React application
FROM node:22-alpine AS builder

WORKDIR /app

COPY application/package*.json ./
RUN npm install

COPY application/ .
RUN npm run build

# Serve React using unprivileged NGINX
FROM nginxinc/nginx-unprivileged:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 8080 