# Getting Started

This guide will help you set up and run the RuoYi Web application.

## Prerequisites

- Node.js (version 16 or higher)
- pnpm (recommended) or npm

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   ```

2. Navigate to the project directory:
   ```bash
   cd ruoyi-web
   ```

3. Install dependencies:
   ```bash
   pnpm install
   ```

## Running the Application

### Development Server

To start the development server:

```bash
pnpm dev
```

The application will be available at `http://localhost:3000` (or another port if 3000 is in use).

### Documentation

To start the VitePress documentation server:

```bash
pnpm docs:dev
```

The documentation will be available at `http://localhost:4000`.

## Building for Production

To build the application for production:

```bash
pnpm build
```

To build the documentation:

```bash
pnpm docs:build
```