#!/bin/bash

echo "🚀 Starting HR System Frontend..."

cd frontend

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or higher."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm."
    exit 1
fi

# Install dependencies
echo "📦 Installing npm dependencies..."
npm install

# Start the development server
echo "🔥 Starting frontend server on http://localhost:3000"
npm run dev