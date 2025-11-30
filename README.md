# Ecommerce Microservice GRPC

    Ecommerce with microservice and GRPC have gateway is graphQL manage and handle request via resolver

GRAPHQL gateway

Incoming Request
      │
      ▼
[ HTTP Middleware ]  <-- 1. EXTRACTS TOKEN (Sets ctx value)
      │
      ▼
[ GraphQL Server  ]  <-- 2. PARSES QUERY
      │
      ▼
   [ Directive ]     <-- 3. CHECKS PERMISSION (Reads ctx value)
         │
         ▼
    [ Resolver ]     <-- 4. EXECUTES LOGIC (Database calls, etc.)

- middleware validate role before access resolver

All service - metadata of user from centralized authorization at gateway

- account service
- authentication service
handle login and refresh token
- catalog service
- order service
