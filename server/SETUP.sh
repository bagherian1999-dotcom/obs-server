#!/bin/bash

echo "======================================"
echo "GoSync Server - Authentication Setup"
echo "======================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "Please install Go 1.21 or later from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
cd "$(dirname "$0")"
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ Failed to install dependencies"
    exit 1
fi

echo "✓ Dependencies installed"
echo ""

# Build the server
echo "🔨 Building server..."
go build -o gosync-server

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✓ Server built successfully"
echo ""

# Create data directory
mkdir -p data/users
echo "✓ Created data directory"
echo ""

echo "======================================"
echo "✅ Setup Complete!"
echo "======================================"
echo ""
echo "To start the server:"
echo "  ./gosync-server"
echo ""
echo "⚠️  IMPORTANT SECURITY NOTES:"
echo "1. Change the JWT secret in auth.go (line 19)"
echo "2. Use HTTPS in production (configure SSL in config.json)"
echo "3. Set strong password requirements"
echo ""
echo "First time users:"
echo "1. Open http://localhost:8080 in your browser"
echo "2. Click 'Register' to create an account"
echo "3. Start syncing!"
echo ""
