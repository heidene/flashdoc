---
title: Examples
---

# Code Examples

Lorem ipsum dolor sit amet, consectetur adipiscing elit. This page provides practical examples for common use cases.

## Basic Examples

### Hello World

The simplest possible application:

```javascript
import { App } from '@example/core';

const app = new App();

app.get('/', (req, res) => {
  res.send('Hello, World!');
});

app.start();
```

Curabitur arcu erat, accumsan id imperdiet et, porttitor at sem.

### REST API

Create a simple REST API for managing users:

```javascript
import { App, Router } from '@example/core';
import { Database } from '@example/db';

const app = new App();
const router = new Router();
const db = new Database();

// Get all users
router.get('/users', async (req, res) => {
  const users = await db.query('SELECT * FROM users');
  res.json(users);
});

// Get user by ID
router.get('/users/:id', async (req, res) => {
  const user = await db.query(
    'SELECT * FROM users WHERE id = ?',
    [req.params.id]
  );
  res.json(user);
});

// Create user
router.post('/users', async (req, res) => {
  const { username, email } = req.body;
  const result = await db.insert('users', {
    username,
    email,
    createdAt: new Date()
  });
  res.status(201).json(result);
});

// Update user
router.put('/users/:id', async (req, res) => {
  const { username, email } = req.body;
  await db.update('users', req.params.id, {
    username,
    email,
    updatedAt: new Date()
  });
  res.json({ message: 'User updated' });
});

// Delete user
router.delete('/users/:id', async (req, res) => {
  await db.delete('users', req.params.id);
  res.status(204).send();
});

app.use('/api', router);
app.start();
```

Vestibulum ac diam sit amet quam vehicula elementum sed sit amet dui.

## Intermediate Examples

### Authentication Middleware

Implement JWT-based authentication:

```javascript
import jwt from 'jsonwebtoken';

const SECRET = process.env.JWT_SECRET;

// Middleware to verify JWT
function authenticate(req, res, next) {
  const token = req.headers.authorization?.replace('Bearer ', '');
  
  if (!token) {
    return res.status(401).json({ 
      error: 'No token provided' 
    });
  }
  
  try {
    const decoded = jwt.verify(token, SECRET);
    req.user = decoded;
    next();
  } catch (err) {
    return res.status(401).json({ 
      error: 'Invalid token' 
    });
  }
}

// Login endpoint
router.post('/auth/login', async (req, res) => {
  const { username, password } = req.body;
  
  // Verify credentials (simplified)
  const user = await db.findUserByUsername(username);
  if (!user || !await verifyPassword(password, user.hashedPassword)) {
    return res.status(401).json({ 
      error: 'Invalid credentials' 
    });
  }
  
  // Generate token
  const token = jwt.sign(
    { 
      userId: user.id, 
      username: user.username 
    },
    SECRET,
    { expiresIn: '1h' }
  );
  
  res.json({ token });
});

// Protected route
router.get('/api/profile', authenticate, (req, res) => {
  res.json({ user: req.user });
});
```

Praesent sapien massa, convallis a pellentesque nec, egestas non nisi.

### File Upload

Handle file uploads with validation:

```javascript
import multer from 'multer';
import path from 'path';

// Configure storage
const storage = multer.diskStorage({
  destination: './uploads',
  filename: (req, file, cb) => {
    const uniqueName = `${Date.now()}-${file.originalname}`;
    cb(null, uniqueName);
  }
});

// File filter
const fileFilter = (req, file, cb) => {
  const allowedTypes = /jpeg|jpg|png|gif/;
  const extname = allowedTypes.test(
    path.extname(file.originalname).toLowerCase()
  );
  const mimetype = allowedTypes.test(file.mimetype);
  
  if (extname && mimetype) {
    cb(null, true);
  } else {
    cb(new Error('Only images are allowed'));
  }
};

const upload = multer({
  storage,
  fileFilter,
  limits: { fileSize: 5 * 1024 * 1024 } // 5MB
});

// Upload endpoint
router.post('/upload', upload.single('file'), (req, res) => {
  if (!req.file) {
    return res.status(400).json({ 
      error: 'No file uploaded' 
    });
  }
  
  res.json({
    message: 'File uploaded successfully',
    filename: req.file.filename,
    path: req.file.path,
    size: req.file.size
  });
});
```

Curabitur non nulla sit amet nisl tempus convallis quis ac lectus.

## Advanced Examples

### WebSocket Real-Time Chat

Build a real-time chat application:

```javascript
import { WebSocketServer } from 'ws';

const wss = new WebSocketServer({ port: 8080 });

const clients = new Map();
const rooms = new Map();

wss.on('connection', (ws) => {
  const clientId = generateId();
  clients.set(clientId, ws);
  
  ws.on('message', (data) => {
    const message = JSON.parse(data);
    
    switch (message.type) {
      case 'join':
        joinRoom(clientId, message.room);
        break;
        
      case 'message':
        broadcastToRoom(message.room, {
          type: 'message',
          from: clientId,
          text: message.text,
          timestamp: Date.now()
        });
        break;
        
      case 'leave':
        leaveRoom(clientId, message.room);
        break;
    }
  });
  
  ws.on('close', () => {
    clients.delete(clientId);
    removeFromAllRooms(clientId);
  });
});

function joinRoom(clientId, roomId) {
  if (!rooms.has(roomId)) {
    rooms.set(roomId, new Set());
  }
  rooms.get(roomId).add(clientId);
  
  broadcastToRoom(roomId, {
    type: 'user_joined',
    userId: clientId
  });
}

function broadcastToRoom(roomId, message) {
  const room = rooms.get(roomId);
  if (!room) return;
  
  room.forEach(clientId => {
    const client = clients.get(clientId);
    if (client?.readyState === WebSocket.OPEN) {
      client.send(JSON.stringify(message));
    }
  });
}
```

Donec sollicitudin molestie malesuada. Vivamus suscipit tortor eget felis porttitor volutpat.

### Database Transactions

Handle complex transactions with rollback:

```javascript
async function transferFunds(fromAccountId, toAccountId, amount) {
  const transaction = await db.beginTransaction();
  
  try {
    // Check sender balance
    const fromAccount = await transaction.query(
      'SELECT balance FROM accounts WHERE id = ? FOR UPDATE',
      [fromAccountId]
    );
    
    if (fromAccount.balance < amount) {
      throw new Error('Insufficient funds');
    }
    
    // Debit sender
    await transaction.update('accounts', fromAccountId, {
      balance: fromAccount.balance - amount
    });
    
    // Credit receiver
    const toAccount = await transaction.query(
      'SELECT balance FROM accounts WHERE id = ? FOR UPDATE',
      [toAccountId]
    );
    
    await transaction.update('accounts', toAccountId, {
      balance: toAccount.balance + amount
    });
    
    // Log transaction
    await transaction.insert('transactions', {
      fromAccount: fromAccountId,
      toAccount: toAccountId,
      amount,
      timestamp: new Date()
    });
    
    await transaction.commit();
    return { success: true };
    
  } catch (error) {
    await transaction.rollback();
    throw error;
  }
}

// Usage
router.post('/transfer', async (req, res) => {
  try {
    const { from, to, amount } = req.body;
    await transferFunds(from, to, amount);
    res.json({ message: 'Transfer successful' });
  } catch (error) {
    res.status(400).json({ error: error.message });
  }
});
```

### Caching Layer

Implement a caching layer with Redis:

```javascript
import Redis from 'ioredis';

const redis = new Redis({
  host: 'localhost',
  port: 6379
});

// Cache middleware
function cache(duration = 60) {
  return async (req, res, next) => {
    const key = `cache:${req.originalUrl}`;
    
    try {
      const cached = await redis.get(key);
      if (cached) {
        return res.json(JSON.parse(cached));
      }
      
      // Override res.json to cache the response
      const originalJson = res.json.bind(res);
      res.json = (data) => {
        redis.setex(key, duration, JSON.stringify(data));
        return originalJson(data);
      };
      
      next();
    } catch (error) {
      console.error('Cache error:', error);
      next();
    }
  };
}

// Usage
router.get('/api/products', cache(300), async (req, res) => {
  const products = await db.query('SELECT * FROM products');
  res.json(products);
});

// Invalidate cache
router.post('/api/products', async (req, res) => {
  const product = await db.insert('products', req.body);
  
  // Invalidate related caches
  await redis.del('cache:/api/products');
  
  res.status(201).json(product);
});
```

Pellentesque in ipsum id orci porta dapibus. Cras ultricies ligula sed magna dictum porta.

## Testing Examples

### Unit Tests

```javascript
import { describe, it, expect, beforeEach } from 'vitest';
import { UserService } from './user-service.js';

describe('UserService', () => {
  let service;
  
  beforeEach(() => {
    service = new UserService();
  });
  
  it('should create a new user', async () => {
    const user = await service.create({
      username: 'testuser',
      email: 'test@example.com'
    });
    
    expect(user).toHaveProperty('id');
    expect(user.username).toBe('testuser');
  });
  
  it('should validate email format', async () => {
    await expect(
      service.create({
        username: 'testuser',
        email: 'invalid-email'
      })
    ).rejects.toThrow('Invalid email format');
  });
});
```

### Integration Tests

```javascript
import { describe, it, expect } from 'vitest';
import request from 'supertest';
import { app } from './app.js';

describe('API Integration Tests', () => {
  it('should get all users', async () => {
    const response = await request(app)
      .get('/api/users')
      .expect(200);
    
    expect(response.body).toBeInstanceOf(Array);
  });
  
  it('should create a user', async () => {
    const response = await request(app)
      .post('/api/users')
      .send({
        username: 'newuser',
        email: 'new@example.com'
      })
      .expect(201);
    
    expect(response.body).toHaveProperty('id');
  });
});
```

Sed porttitor lectus nibh. Mauris blandit aliquet elit, eget tincidunt nibh pulvinar a.

## More Resources

- [API Reference](../api/overview.md)
- [Configuration Guide](./configuration.md)
- [Best Practices](./best-practices.md)

Vivamus magna justo, lacinia eget consectetur sed, convallis at tellus.
