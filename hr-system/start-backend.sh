#!/bin/bash

echo "🚀 Starting HR System Backend..."

cd backend

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Run the backend server
echo "🔥 Starting backend server on http://localhost:8080"
go run main.go