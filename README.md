# Authentication Service

## Overview
The **Authentication Service** is a microservice designed to handle user authentication following the **OAuth2 standard**. It provides secure JWT token issuance, session management, and token validation. Third-party applications can integrate this service for user login, signup, and authentication.

## Features
- **OAuth2-Based Authentication**:
  - Implements OAuth2 flows for authentication and token management.
  - Supports **password grant**, **client credentials**, and **refresh token** flows.
- **JWT Token Issuance & Management**:
  - Issues **access tokens** (short-lived) for API authentication.
  - Issues **refresh tokens** (long-lived, stored in cookies) for renewing access tokens.
- **User Session Management**:
  - Each login generates a **refresh token tied to a session ID**.
  - Redis **blacklist** is used to immediately revoke access tokens.
  - SQLite **stores revoked refresh tokens** to prevent reuse.
- **Token Validation**:
  - Ensures that only valid tokens grant access to protected resources.
  - API users must include the **access token** in the `Authorization` header.
  - When the access token expires, the **refresh token** is used to generate a new one.
- **Multiple Site Support**:
  - Supports authentication across different `siteId` values.
  - Separate JWT tokens are issued per site to maintain data isolation.
- **Security**:
  - Secure password hashing using `bcrypt`.
  - Tokens stored securely with proper expiration and revocation policies.
  - Protection against **token replay attacks** with session-based validation.

## Technologies
- **Go**: The backend programming language for the service.
- **Gin**: A fast and lightweight web framework for the API.
- **OAuth2**: Industry-standard authentication framework.
- **JWT**: JSON Web Tokens for secure authentication.
- **bcrypt**: Secure password hashing.
- **Redis**: Manages access token blacklist for immediate revocation.
- **SQLite**: Stores user session data and revoked refresh tokens.

## Token Flow
1. **User Login**
  - User logs in with credentials via `/:siteId/login`.
  - The system generates:
    - **Access token** (short-lived, used for authentication)
    - **Refresh token** (long-lived, stored in a cookie, tied to a session ID)
2. **Authenticated API Calls**
  - The user includes the **access token** in the `Authorization` header when calling APIs.
  - Example: `GET /:siteId/jwt` to retrieve user info.
3. **Access Token Expire**
  - If the access token expires, the client sends the **refresh token** to `/:siteId/refresh` to obtain a new access token.
4. **Logout & Token Revocation**
  - The user logs out via `/:siteId/signout`:
    - The **access token** is immediately blacklisted in Redis.
    - The **refresh token** is revoked in SQLite.
    - Future API requests with a blacklisted access token are denied.
    - Future refresh attempts using a revoked refresh token are denied.

## API Endpoints
- `POST /:siteId/signup` → Register a new user.
- `POST /:siteId/login` → Authenticate user, return **access & refresh tokens**.
- `GET /:siteId/jwt` → Validate **access token**, return user information.
- `GET /:siteId/refresh` → Use refresh token to obtain a new access token.
- `GET /:siteId/signout` → Revoke tokens and terminate the session.

## Environment Variables
- `SECRET_KEY=1` → Secret key for signing JWT tokens.
- `REDIS_HOST=localhost:6379` → Redis connection for managing blacklisted tokens.
- `CACHE_PATH=db.sqlite` → SQLite database path for storing revoked refresh tokens.

This service ensures secure authentication while adhering to OAuth2 best practices.

