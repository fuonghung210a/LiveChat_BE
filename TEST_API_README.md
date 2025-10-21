# Auth API Testing Guide

## Overview
This document explains how to test the Authentication API and addresses issues found during testing.

## Test Script
A comprehensive bash test script has been created: `test_auth_api.sh`

### Features
- 12 comprehensive test cases covering register, login, and profile endpoints
- Color-coded output for easy result identification
- Automatic token management between tests
- Detailed request/response logging

### Usage
```bash
./test_auth_api.sh
```

## Test Cases

### Registration Tests (1-5)
1. **Register new user** - Valid registration with all required fields
2. **Duplicate email** - Ensures duplicate emails are rejected
3. **Invalid email format** - Validates email format
4. **Short password** - Ensures password minimum length (6 chars)
5. **Missing fields** - Validates required fields

### Login Tests (6-9)
6. **Valid credentials** - Successful login
7. **Wrong password** - Invalid password rejection
8. **Non-existent email** - User not found
9. **Invalid email format** - Email validation on login

### Profile Tests (10-12)
10. **Valid token** - Get profile with authentication
11. **No token** - Unauthorized access rejection
12. **Invalid token** - Malformed token rejection

## Issues Found & Solutions

### Issue 1: JWT Token Generation Fails (HTTP 500)
**Error:** `{"error":"failed to generate token"}`

**Root Cause:** In `internal/util/jwt.go:37`, the `token.SignedString()` method expects `[]byte` but receives a `string`.

**Fix:**
```go
// Change line 37 from:
tokenString, err := token.SignedString(JWTSecretKey)

// To:
tokenString, err := token.SignedString([]byte(JWTSecretKey))
```

### Issue 2: Missing Environment Variables
**Error:** JWT secret not configured

**Solution:** Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

Then update the JWT_SECRET:
```env
JWT_SECRET=your-super-secret-key-change-in-production
```

### Issue 3: AuthMiddleware Missing for Profile Endpoint
The `/auth/profile` endpoint expects `user_id` in the context (set by middleware), but the middleware isn't applied in the router.

**Fix in `internal/router/router.go`:**
```go
// Add middleware import
import "go_starter/internal/middleware"

// Apply auth middleware to protected routes
authGroup := api.Group("/auth")
{
    authGroup.POST("/register", authHandler.Register)
    authGroup.POST("/login", authHandler.Login)
    // Protected route - requires authentication
    authGroup.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
}
```

## Required Fixes Summary

### 1. Fix JWT Signing (CRITICAL)
File: `internal/util/jwt.go:37`
```go
tokenString, err := token.SignedString([]byte(JWTSecretKey))
```

### 2. Create .env File
```bash
cp .env.example .env
# Edit .env and set JWT_SECRET to a secure value
```

### 3. Add Auth Middleware
File: `internal/router/router.go`

Create middleware file if it doesn't exist: `internal/middleware/auth.go`
```go
package middleware

import (
    "go_starter/internal/util"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        // Extract token from "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        token := parts[1]
        claims, err := util.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        // Set user information in context
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)
        c.Next()
    }
}
```

Then apply it in the router:
```go
authGroup.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
```

## Running the Tests

### Prerequisites
1. Ensure your Go application is running:
   ```bash
   go run cmd/main.go
   ```

2. Make sure the database is accessible

3. (Optional) Install `jq` for formatted JSON output:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install jq

   # macOS
   brew install jq
   ```

### Execute Tests
```bash
chmod +x test_auth_api.sh
./test_auth_api.sh
```

## Expected Results
After implementing the fixes above, all 12 tests should pass:
- ✓ Tests 3, 4, 5, 9: Validation tests (already passing)
- ✓ Tests 7, 8, 11, 12: Rejection tests (already passing)
- ✓ Tests 1, 6: Registration and login (will pass after JWT fix)
- ✓ Test 2: Duplicate email (will pass after JWT fix for first user creation)
- ✓ Test 10: Profile retrieval (will pass after middleware fix)

## API Endpoints Reference

### POST /api/auth/register
Register a new user account.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201):**
```json
{
  "message": "user registered successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### POST /api/auth/login
Authenticate and receive JWT token.

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

### GET /api/auth/profile
Get current user profile (requires authentication).

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200):**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

## Troubleshooting

### Connection Refused
- Ensure the server is running on port 8080
- Check `BASE_URL` in the test script matches your server configuration

### Database Errors
- Verify database connection in `.env`
- Ensure MySQL is running
- Check database migrations have been applied

### JWT Token Errors
- Verify `JWT_SECRET` is set in `.env`
- Ensure the JWT signing fix has been applied
- Check token isn't expired (default 24h)

## Notes
- The test script generates unique email addresses using timestamps to avoid conflicts
- Temporary files are automatically cleaned up on script exit
- The script requires `curl` which is typically pre-installed on most systems
