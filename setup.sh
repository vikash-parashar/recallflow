#!/bin/bash

# RecallFlow - Quick Setup Script
# This script helps you get RecallFlow running locally

set -e

echo "🚀 RecallFlow Setup"
echo "===================="
echo ""

# Check prerequisites
echo "📋 Checking prerequisites..."

command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed. Please install Go 1.22+"; exit 1; }
command -v node >/dev/null 2>&1 || { echo "❌ Node.js is not installed. Please install Node.js 18+"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "❌ Docker is not installed. Please install Docker"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ Docker Compose is not installed. Please install Docker Compose"; exit 1; }

echo "✅ All prerequisites satisfied"
echo ""

# Backend setup
echo "🔧 Setting up backend..."
cd backend

if [ ! -f .env ]; then
    echo "Creating .env file..."
    cp .env.example .env
    echo "⚠️  Please edit backend/.env with your actual credentials before running"
fi

echo "Starting PostgreSQL and Redis..."
docker-compose up -d

echo "Waiting for database to be ready..."
sleep 5

echo "Installing Go dependencies..."
go mod download

echo "✅ Backend setup complete"
echo ""

# Frontend setup
echo "🎨 Setting up frontend..."
cd ../frontend

if [ ! -f .env.local ]; then
    echo "Creating .env.local file..."
    cp .env.local.example .env.local
fi

echo "Installing npm dependencies..."
npm install

echo "✅ Frontend setup complete"
echo ""

# Summary
echo "🎉 Setup Complete!"
echo "=================="
echo ""
echo "Next steps:"
echo ""
echo "1. Edit backend/.env with your API credentials:"
echo "   - TWILIO_ACCOUNT_SID"
echo "   - TWILIO_AUTH_TOKEN"
echo "   - OPENAI_API_KEY"
echo "   - JWT_SECRET (generate a random string)"
echo ""
echo "2. Run database migrations:"
echo "   cd backend && make migrate"
echo ""
echo "3. Start the backend:"
echo "   cd backend && make dev"
echo "   (Runs on http://localhost:8080)"
echo ""
echo "4. Start the frontend (in a new terminal):"
echo "   cd frontend && npm run dev"
echo "   (Runs on http://localhost:3000)"
echo ""
echo "5. Visit http://localhost:3000 to see RecallFlow!"
echo ""
echo "📚 For more details, see docs/README.md"
