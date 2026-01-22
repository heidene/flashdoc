---
title: Configuration
description: Learn how to configure your application
---

# Configuration Guide

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum ac diam sit amet quam vehicula elementum sed sit amet dui.

## Configuration Files

Your application can be configured through multiple files:

### config.json

```json
{
  "appName": "My Application",
  "port": 3000,
  "logLevel": "info"
}
```

Proin eget tortor risus. Curabitur non nulla sit amet nisl tempus convallis quis ac lectus.

### Environment Variables

Create a `.env` file with the following variables:

```env
NODE_ENV=production
API_TIMEOUT=5000
CACHE_ENABLED=true
```

Curabitur arcu erat, accumsan id imperdiet et, porttitor at sem. Vivamus magna justo, lacinia eget consectetur sed.

## Configuration Options

### General Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `appName` | string | "App" | Vestibulum ac diam |
| `port` | number | 3000 | Mauris blandit aliquet |
| `logLevel` | string | "info" | Nulla porttitor accumsan |

### Advanced Options

#### Caching

Lorem ipsum dolor sit amet:

- `CACHE_ENABLED`: Boolean to enable/disable caching
- `CACHE_TTL`: Time to live in seconds (default: 3600)
- `CACHE_STORAGE`: Storage backend ("memory" or "redis")

#### Security

Cras ultricies ligula sed magna dictum porta:

- `JWT_SECRET`: Secret key for JWT tokens
- `SESSION_TIMEOUT`: Session timeout in milliseconds
- `CORS_ORIGINS`: Allowed CORS origins (comma-separated)

## Environment-Specific Configuration

### Development

```json
{
  "logLevel": "debug",
  "cache": {
    "enabled": false
  }
}
```

### Production

```json
{
  "logLevel": "error",
  "cache": {
    "enabled": true,
    "ttl": 7200
  }
}
```

Pellentesque in ipsum id orci porta dapibus. Donec sollicitudin molestie malesuada.

## Loading Configuration

The application loads configuration in the following order:

1. Default configuration
2. Environment-specific configuration
3. Environment variables
4. Command-line arguments

Later sources override earlier ones.

## Best Practices

- ✅ Use environment variables for secrets
- ✅ Keep configuration files in version control
- ✅ Document all configuration options
- ❌ Never commit secrets to the repository
- ❌ Don't use production credentials in development

Sed porttitor lectus nibh. Curabitur aliquet quam id dui posuere blandit.
