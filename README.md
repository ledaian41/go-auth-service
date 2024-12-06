# Authentication Service

## Overview
The **Authentication Service** is a microservice designed to handle user authentication, provide secure JWT tokens for access, and manage user identities across multiple sites. This service can be integrated by third-party applications for user login, signup, and token validation.

## Features
- **User Authentication**:
    - Provides login and signup functionality with secure password hashing (bcrypt).
- **JWT Token Generation**:
    - Issues **JWT access tokens** to authenticated users for secure access.
    - **Refresh Token**: Stored in a cookie and used to re-generate the access token when it expires.
    - **Access Token**: Used in the `Authorization` header for verifying access to protected routes.
- **Token Validation**:
    - Verifies the authenticity of JWT tokens to ensure only authenticated users can access protected resources.
- **Multiple Site Support**:
    - Supports authentication for users across different sites identified by `siteId`.
    - Different JWT tokens are generated for each site to provide isolation of user data.
- **Scalability**:
    - Designed to handle a growing user base efficiently with optimized token management.
- **Security**:
    - Implements secure practices to protect user data and JWT tokens, including password hashing, secure token storage (via cookies), and token expiration management.

## Technologies
- **Go**: The backend programming language used to build the service.
- **Gin**: A fast and lightweight web framework for building the API.
- **JWT**: JSON Web Tokens (JWT) for secure user authentication and authorization.
- **bcrypt**: Used for secure password hashing.

