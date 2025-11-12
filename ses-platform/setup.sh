#!/bin/bash
# SES Platform - Complete Setup Script
# This script will create the entire project structure and files

set -e  # Exit on error

echo "🚀 SES Platform - Repository Setup"
echo "===================================="
echo ""

# Create base directory
PROJECT_NAME="ses-platform"
read -p "Enter project directory name [$PROJECT_NAME]: " input
PROJECT_NAME=${input:-$PROJECT_NAME}

if [ -d "$PROJECT_NAME" ]; then
    read -p "⚠️  Directory $PROJECT_NAME already exists. Continue? (y/n): " confirm
    if [ "$confirm" != "y" ]; then
        echo "Aborted."
        exit 1
    fi
else
    mkdir -p "$PROJECT_NAME"
fi

cd "$PROJECT_NAME"

echo "📁 Creating directory structure..."

# Create all directories
mkdir -p .github/workflows
mkdir -p backend
mkdir -p frontend/{public,src}
mkdir -p database/migrations
mkdir -p capabilities
mkdir -p enablers
mkdir -p monitoring/grafana/{dashboards,datasources}
mkdir -p docs
mkdir -p scripts
mkdir -p uploads
mkdir -p logs
mkdir -p backups

echo "✅ Directory structure created"
echo ""
echo "📝 Now you need to:"
echo ""
echo "1. Copy the following files from Claude's artifacts:"
echo "   - Copy 'SES Platform Backend API' → backend/main.go"
echo "   - Copy 'SES Platform Database Schema' → database/schema.sql"
echo "   - Copy 'SES Platform React App' → frontend/src/App.jsx"
echo "   - Copy 'Docker Compose Setup' → docker-compose.yml"
echo "   - Copy 'Setup & Documentation' → README.md"
echo "   - Copy 'Complete GitHub Repository Structure' → SETUP_GUIDE.md"
echo ""
echo "2. Copy your existing capability and enabler markdown files:"
echo "   - Copy capabilities/*.md → capabilities/"
echo "   - Copy enablers/*.md → enablers/"
echo ""
echo "3. Create configuration files using templates from SETUP_GUIDE.md:"
echo ""

# Create .gitignore
cat > .gitignore << 'EOF'
# Environment variables
.env
.env.local
.env.*.local

# Dependencies
node_modules/
vendor/

# Build outputs
build/
dist/
*.exe
*.dll
*.so
*.dylib

# IDEs
.idea/
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# Database
*.db
*.sqlite
pgdata/

# Uploads
uploads/
tmp/

# Test coverage
coverage/
*.cover
*.coverage

# Docker
docker-compose.override.yml
EOF

echo "   ✅ Created .gitignore"

# Create .env.example
cat > .env.example << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ses_platform
DB_USER=ses_user
DB_PASSWORD=ses_password
DB_SSLMODE=disable

# Backend Configuration
SERVER_PORT=8080
GIN_MODE=release
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=*

# Frontend Configuration
REACT_APP_API_URL=http://localhost:8080/api/v1

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# Monitoring
PROMETHEUS_PORT=9090
GRAFANA_PORT=3001
GRAFANA_ADMIN_PASSWORD=admin

# Security (CHANGE IN PRODUCTION!)
JWT_SECRET=your-secret-key-change-in-production
SESSION_SECRET=your-session-secret-change-in-production
EOF

echo "   ✅ Created .env.example"

# Create LICENSE
cat > LICENSE << 'EOF'
MIT License

Copyright (c) 2025 Your Organization

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

echo "   ✅ Created LICENSE"

# Create backend .gitignore
cat > backend/.gitignore << 'EOF'
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
main

# Test binary
*.test

# Output
*.out

# Go workspace
go.work

# Vendor
vendor/

# Environment
.env
.env.local
EOF

echo "   ✅ Created backend/.gitignore"

# Create backend Dockerfile
cat > backend/Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates wget tzdata
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./main"]
EOF

echo "   ✅ Created backend/Dockerfile"

# Create backend go.mod
cat > backend/go.mod << 'EOF'
module ses-platform

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/driver/postgres v1.5.4
    gorm.io/gorm v1.25.5
)
EOF

echo "   ✅ Created backend/go.mod"

# Create frontend .gitignore
cat > frontend/.gitignore << 'EOF'
# dependencies
/node_modules
/.pnp
.pnp.js

# testing
/coverage

# production
/build

# misc
.DS_Store
.env.local
.env.development.local
.env.test.local
.env.production.local

npm-debug.log*
yarn-debug.log*
yarn-error.log*
EOF

echo "   ✅ Created frontend/.gitignore"

# Create frontend Dockerfile
cat > frontend/Dockerfile << 'EOF'
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=$REACT_APP_API_URL
RUN npm run build

FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/build /usr/share/nginx/html
EXPOSE 80
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost/ || exit 1
CMD ["nginx", "-g", "daemon off;"]
EOF

echo "   ✅ Created frontend/Dockerfile"

# Create nginx.conf
cat > frontend/nginx.conf << 'EOF'
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://backend:8080/api/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    error_page 404 /index.html;
}
EOF

echo "   ✅ Created frontend/nginx.conf"

# Create package.json
cat > frontend/package.json << 'EOF'
{
  "name": "ses-platform-frontend",
  "version": "1.0.0",
  "description": "Simulation Environment Specification Platform - Frontend",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "lucide-react": "^0.263.1",
    "recharts": "^2.10.0"
  },
  "devDependencies": {
    "react-scripts": "5.0.1",
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "typescript": "^4.9.5",
    "tailwindcss": "^3.3.0",
    "autoprefixer": "^10.4.14",
    "postcss": "^8.4.24"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test",
    "eject": "react-scripts eject"
  },
  "eslintConfig": {
    "extends": ["react-app"]
  },
  "browserslist": {
    "production": [">0.2%", "not dead", "not op_mini all"],
    "development": ["last 1 chrome version", "last 1 firefox version", "last 1 safari version"]
  }
}
EOF

echo "   ✅ Created frontend/package.json"

# Create index.html
cat > frontend/public/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="%PUBLIC_URL%/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#000000" />
    <meta name="description" content="Simulation Environment Specification Platform" />
    <title>SES Platform</title>
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>
  </body>
</html>
EOF

echo "   ✅ Created frontend/public/index.html"

# Create index.js
cat > frontend/src/index.js << 'EOF'
import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
EOF

echo "   ✅ Created frontend/src/index.js"

# Create index.css
cat > frontend/src/index.css << 'EOF'
@tailwind base;
@tailwind components;
@tailwind utilities;

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

code {
  font-family: source-code-pro, Menlo, Monaco, Consolas, 'Courier New', monospace;
}
EOF

echo "   ✅ Created frontend/src/index.css"

# Create tailwind.config.js
cat > frontend/tailwind.config.js << 'EOF'
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {},
  },
  plugins: [],
}
EOF

echo "   ✅ Created frontend/tailwind.config.js"

# Create postcss.config.js
cat > frontend/postcss.config.js << 'EOF'
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
EOF

echo "   ✅ Created frontend/postcss.config.js"

# Create prometheus.yml
cat > monitoring/prometheus.yml << 'EOF'
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'ses-backend'
    static_configs:
      - targets: ['backend:8080']
    metrics_path: '/metrics'

  - job_name: 'postgres'
    static_configs:
      - targets: ['postgres:5432']

  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
EOF

echo "   ✅ Created monitoring/prometheus.yml"

# Create setup script
cat > scripts/setup.sh << 'EOF'
#!/bin/bash
echo "🚀 Setting up SES Platform..."
command -v docker >/dev/null 2>&1 || { echo "❌ Docker required" >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose required" >&2; exit 1; }
[ ! -f .env ] && cp .env.example .env && echo "📝 Created .env file"
mkdir -p uploads logs backups
docker-compose pull
docker-compose build
echo "✅ Setup complete! Run: docker-compose up -d"
EOF

chmod +x scripts/setup.sh
echo "   ✅ Created scripts/setup.sh"

echo ""
echo "======================================"
echo "✅ Base structure created successfully!"
echo "======================================"
echo ""
echo "📋 Next steps:"
echo ""
echo "1. Copy these artifacts from Claude into the files:"
echo "   • backend/main.go (from 'SES Platform - Go Backend API')"
echo "   • database/schema.sql (from 'SES Platform - PostgreSQL Database Schema')"
echo "   • frontend/src/App.jsx (from 'Simulation Environment Specification Platform')"
echo "   • docker-compose.yml (from 'SES Platform - Docker Compose Setup')"
echo "   • README.md (from 'SES Platform - Setup & Documentation')"
echo ""
echo "2. Copy your capability and enabler files:"
echo "   • capabilities/*.md"
echo "   • enablers/*.md"
echo ""
echo "3. Initialize git and push to GitHub:"
echo ""
echo "   cd $PROJECT_NAME"
echo "   git init"
echo "   git add ."
echo "   git commit -m 'Initial commit: Complete SES Platform implementation'"
echo "   git remote add origin https://github.com/YOUR_USERNAME/ses-platform.git"
echo "   git branch -M main"
echo "   git push -u origin main"
echo ""
echo "4. After pushing, run the setup:"
echo ""
echo "   ./scripts/setup.sh"
echo "   docker-compose up -d"
echo ""
echo "🎉 Then access your platform at http://localhost:3000"
echo ""
EOF

chmod +x "$0"
