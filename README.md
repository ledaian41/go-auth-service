# Core Authentication Service

## Overview
The **Core Authentication Service** is a microservice designed to handle user authentication and provide JWT tokens for secure access. Third-party applications can integrate with this service to authenticate users and validate their identities.

## Features
- **User Authentication**: Login and signup functionality with secure password hashing.
- **JWT Token Generation**: Issue JWT tokens to authenticated users for secure access.
- **Token Validation**: Verify the authenticity of JWT tokens.
- **Scalability**: Designed to handle a growing user base efficiently.
- **Security**: Implements secure practices to protect user data and tokens.

## Technologies
- **Go**: Backend programming language.
- **Gin**: Lightweight web framework.
- **MongoDB**: Database for storing user data.
- **JWT**: JSON Web Token for secure user authentication.
