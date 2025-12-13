# Project Architecture & Solution Summary

## 1. Architecture Overview
This project is an **Authentication Service** built with Go, designed to support both internal microservices (via **gRPC**) and external clients (via **REST API**).

-   **Language**: Go
-   **Web Framework**: Gin (HTTP)
-   **RPC Framework**: gRPC
-   **Persistence**: SQLite (via GORM)
-   **Caching/State**: Redis (Token Versioning)

## 2. Key Components

### 2.1. Service Layer
The application is structured into modular services:
-   **Auth Service**: Central orchestrator for Signup, Login, and Token verification.
-   **User Service**: Manages user identity, CRUD operations, and persistence.
-   **Token Service**: Manages session validity using a JTI (JWT ID) Allowlist.
-   **Site Service**: Handles multi-tenancy configuration (Site IDs and Secret Keys).

### 2.2. Interfaces
-   **REST API**: Exposes endpoints for public access (`/login`, `/signup`, `/refresh`).
-   **gRPC API**: Exposes internal methods for other services to verify tokens or exchange credentials.

## 3. Core Solutions

### 3.1. Authentication Strategy (Dual Token)
The system uses a **Dual Token System** (Access Token + Refresh Token) to balance security and user experience.

-   **Access Token**:
    -   **Format**: JWT.
    -   **Lifespan**: Short (e.g., 15-30 mins).
    -   **Signing**: Signed with the **Site's Secret Key**. This ensures isolation between different tenants/sites.
    -   **Usage**: Used for API authorization.

-   **Refresh Token**:
    -   **Format**: JWT.
    -   **Lifespan**: Long (e.g., 7 days).
    -   **Signing**: Signed with the **Global Secret Key**.
    -   **Content**: Contains a `jti` (Session ID).
    -   **Usage**: Used to obtain new Access Tokens.

### 3.2. Secure Session Management (JTI Allowlist)
To secure Refresh Tokens without storing sensitive data, we implemented a **JTI Allowlist** mechanism.

1.  **Usage**: When a user logs in, a random unique Session ID (`jti`) is generated.
2.  **Storage**: The `jti` is stored in the database (`user_tokens` table). The actual Refresh Token string is **never stored**.
3.  **Token Issuance**: The `jti` is embedded into the Refresh Token's claims.
4.  **Validation**: When a Refresh Token is presented, the system extracts the `jti` and checks if it exists in the database.
5.  **Revocation**:
    -   **Logout**: Deletes the `jti` from the database.
    -   **Effect**: The Refresh Token becomes immediately invalid because its `jti` is no longer in the allowlist.

### 3.3. User Persistence
-   **Old Approach**: Users were stored in memory (lost on restart).
-   **Current Approach**: Users are persisted in **SQLite** using **GORM**. This ensures distinct user accounts per Site and data durability.

### 3.4. Token Versioning (Redis)
-   To handle cases like "Change Password" (where all existing tokens must be invalidated), the system uses a `token_version` integer stored in **Redis**.
-   Tokens include a `token_version` claim.
-   If the Redis version is higher than the token's version, the token is rejected.
