---
title: Authentication API
---

# Authentication API

Lorem ipsum dolor sit amet, consectetur adipiscing elit. This page documents all authentication-related endpoints.

## Endpoints

### POST /auth/login

Authenticate a user and receive an access token.

**Request:**

```json
{
  "username": "johndoe",
  "password": "securepassword123"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600,
    "tokenType": "Bearer"
  }
}
```

Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia.

### POST /auth/logout

Invalidate the current session token.

**Request:**

```bash
POST /auth/logout
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "message": "Successfully logged out"
}
```

### POST /auth/refresh

Refresh an expired access token using a refresh token.

**Request:**

```json
{
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 3600
  }
}
```

Praesent sapien massa, convallis a pellentesque nec, egestas non nisi.

### POST /auth/register

Register a new user account.

**Request:**

```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepass123",
  "firstName": "John",
  "lastName": "Doe"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "userId": "user_789",
    "username": "newuser",
    "email": "newuser@example.com"
  }
}
```

### POST /auth/password-reset

Request a password reset email.

**Request:**

```json
{
  "email": "user@example.com"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Password reset email sent"
}
```

Curabitur non nulla sit amet nisl tempus convallis quis ac lectus.

### POST /auth/password-reset/confirm

Confirm password reset with token.

**Request:**

```json
{
  "token": "reset_token_here",
  "newPassword": "newsecurepassword123"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Password successfully reset"
}
```

## Token Structure

JWT tokens contain the following claims:

```json
{
  "sub": "user_123",
  "username": "johndoe",
  "email": "john@example.com",
  "roles": ["user", "admin"],
  "iat": 1706875200,
  "exp": 1706878800
}
```

Donec sollicitudin molestie malesuada. Nulla quis lorem ut libero malesuada feugiat.

## Security Best Practices

### Token Storage

- ✅ Store tokens in httpOnly cookies
- ✅ Use secure flag in production
- ❌ Never store tokens in localStorage
- ❌ Don't expose tokens in URLs

### Password Requirements

Passwords must meet these requirements:

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

### Rate Limiting

Login attempts are rate limited:

- **5 failed attempts** → 5 minute lockout
- **10 failed attempts** → 30 minute lockout
- **20 failed attempts** → Account locked (requires email verification)

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| `AUTH_001` | Invalid credentials | Username or password incorrect |
| `AUTH_002` | Token expired | Access token has expired |
| `AUTH_003` | Invalid token | Token is malformed or invalid |
| `AUTH_004` | Account locked | Too many failed login attempts |
| `AUTH_005` | Email not verified | User must verify email first |

Vivamus suscipit tortor eget felis porttitor volutpat.

## Example Usage

### JavaScript

```javascript
import { AuthClient } from '@example/sdk';

const auth = new AuthClient({
  baseURL: 'https://api.example.com'
});

// Login
const { accessToken } = await auth.login({
  username: 'johndoe',
  password: 'password123'
});

// Use token for authenticated requests
auth.setToken(accessToken);

// Get current user
const user = await auth.getCurrentUser();
```

### Python

```python
from example_sdk import AuthClient

auth = AuthClient(base_url='https://api.example.com')

# Login
response = auth.login(
    username='johndoe',
    password='password123'
)

access_token = response['accessToken']

# Use token
auth.set_token(access_token)
user = auth.get_current_user()
```

### cURL

```bash
# Login
curl -X POST https://api.example.com/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"johndoe","password":"password123"}'

# Use token
curl https://api.example.com/api/users/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

## OAuth 2.0 Support

The API supports OAuth 2.0 authentication with the following providers:

- **Google** - `GET /auth/oauth/google`
- **GitHub** - `GET /auth/oauth/github`
- **Facebook** - `GET /auth/oauth/facebook`

### OAuth Flow

1. Redirect user to `/auth/oauth/{provider}`
2. User authenticates with provider
3. Redirect back to your callback URL with authorization code
4. Exchange code for access token

Cras ultricies ligula sed magna dictum porta. Pellentesque in ipsum id orci porta dapibus.

## Multi-Factor Authentication

MFA can be enabled for additional security. Mauris blandit aliquet elit:

### Enable MFA

```bash
POST /auth/mfa/enable
Authorization: Bearer <token>
```

Returns QR code and backup codes.

### Verify MFA

```bash
POST /auth/mfa/verify
{
  "code": "123456"
}
```

Sed porttitor lectus nibh. Vivamus magna justo, lacinia eget consectetur sed.
