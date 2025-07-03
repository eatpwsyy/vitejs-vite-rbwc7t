#!/bin/bash

echo "🚀 Starting Complete HR Management System..."

# Function to kill processes on exit
cleanup() {
    echo "🛑 Shutting down HR System..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit
}

# Set up trap to cleanup on exit
trap cleanup SIGINT SIGTERM

# Start backend in background
echo "📊 Starting backend server..."
cd backend
go run main.go &
BACKEND_PID=$!

# Wait for backend to start
sleep 3

# Start frontend in background
echo "🎨 Starting frontend server..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!

echo ""
echo "✅ HR Management System is starting up!"
echo ""
echo "📊 Backend API: http://localhost:8080"
echo "🎨 Frontend UI: http://localhost:3000"
echo ""
echo "🔑 Default login credentials:"
echo "   Email: admin@hr.com"
echo "   Password: password"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for processes
wait