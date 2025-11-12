#!/bin/bash
echo "🚀 Setting up SES Platform..."
command -v docker >/dev/null 2>&1 || { echo "❌ Docker required" >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose required" >&2; exit 1; }
[ ! -f .env ] && cp .env.example .env && echo "📝 Created .env file"
mkdir -p uploads logs backups
docker-compose pull
docker-compose build
echo "✅ Setup complete! Run: docker-compose up -d"
