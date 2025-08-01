# SocialPay Deployment

## Quick Deploy to Server

Deploy both V1 and V2 backends with frontend to the server:

```bash
./deploy_to_server.sh
```

## Local Development

For local development with all services:

```bash
./deploy.sh --production
```

## Service URLs

### Local Development
- Frontend: http://localhost:3000
- Backend V1: http://localhost:8004
- Backend V2: http://localhost:8082
- Swagger: http://localhost:8082/swagger/index.html

### Server (196.190.251.194)
- Frontend: http://196.190.251.194:3000
- Backend V1: http://196.190.251.194:8004
- Backend V2: http://196.190.251.194:8082
- Swagger: http://196.190.251.194:8082/swagger/index.html
- Nginx: http://196.190.251.194

## What's Deployed

- ✅ Frontend (Next.js)
- ✅ Backend V1 (Mux)
- ✅ Backend V2 (Gin) with Swagger
- ✅ Nginx (Reverse Proxy) 