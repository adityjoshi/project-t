#!/bin/bash

# Project Synapse - Quick Start Script

set -e

echo "🧠 Project Synapse - Quick Start"
echo "================================"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo ""
    echo "Please create a .env file with your OpenAI API key:"
    echo "  echo 'OPENAI_API_KEY=sk-your-key-here' > .env"
    echo ""
    echo "Get your API key from: https://platform.openai.com/api-keys"
    exit 1
fi

# Check if OPENAI_API_KEY is set
if ! grep -q "OPENAI_API_KEY=sk-" .env 2>/dev/null; then
    echo "⚠️  Warning: OPENAI_API_KEY not found in .env file"
    echo "   Make sure your .env file contains: OPENAI_API_KEY=sk-..."
    echo ""
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running!"
    echo "   Please start Docker Desktop and try again"
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Build and start services
echo "🚀 Starting all services..."
echo ""

docker-compose up -d --build

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check service health
echo ""
echo "📊 Checking service status..."
docker-compose ps

echo ""
echo "✅ Services are starting!"
echo ""
echo "🌐 Access the application:"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8080"
echo ""
echo "📝 View logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Stop services:"
echo "   docker-compose down"
echo ""

