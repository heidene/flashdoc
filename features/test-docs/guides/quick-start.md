---
title: Quick Start
---

# Quick Start Guide

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Get started with your first application in just 5 minutes!

## Your First Application

### Create a New Project

```bash
npx create-app my-app
cd my-app
```

Vivamus suscipit tortor eget felis porttitor volutpat. Vestibulum ante ipsum primis in faucibus orci luctus.

### Run the Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser. Donec rutrum congue leo eget malesuada.

## Basic Example

Here's a simple "Hello World" example:

```javascript
import { App } from '@example/core';

const app = new App({
  port: 3000,
  title: 'My First App'
});

app.get('/', (req, res) => {
  res.send('Hello, World!');
});

app.start();
```

Praesent sapien massa, convallis a pellentesque nec, egestas non nisi. Curabitur non nulla sit amet nisl tempus convallis.

## Project Structure

Your project will have the following structure:

```
my-app/
├── src/
│   ├── index.js      # Main entry point
│   ├── routes/       # Route handlers
│   └── utils/        # Utility functions
├── public/           # Static files
├── tests/            # Test files
└── package.json      # Dependencies
```

## Add Your First Route

Create `src/routes/hello.js`:

```javascript
export function helloRoute(app) {
  app.get('/hello/:name', (req, res) => {
    const { name } = req.params;
    res.json({ 
      message: `Hello, ${name}!`,
      timestamp: new Date()
    });
  });
}
```

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Import and register the route:

```javascript
import { helloRoute } from './routes/hello.js';

helloRoute(app);
```

## Next Steps

Now that you have a basic application running:

1. **Explore the [API Reference](../api/overview.md)** - Quisque velit nisi, pretium ut lacinia
2. **Learn about [Configuration](./configuration.md)** - Vestibulum ac diam sit amet quam
3. **Check out [Examples](./examples.md)** - Sed porttitor lectus nibh

## Common Tasks

### Adding Dependencies

```bash
npm install package-name
```

### Running Tests

```bash
npm test
```

### Building for Production

```bash
npm run build
```

Pellentesque in ipsum id orci porta dapibus. Nulla quis lorem ut libero malesuada feugiat.

## Getting Help

- 📖 [Documentation](../README.md)
- 💬 [Community Forum](https://forum.example.com)
- 🐛 [Issue Tracker](https://github.com/example/project/issues)

Curabitur aliquet quam id dui posuere blandit. Vivamus magna justo, lacinia eget consectetur sed.
