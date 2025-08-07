#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Build the frontend application
echo "Building frontend..."
(
  cd frontend
  npm install
  npm run build
)

# Compile the backend application
echo "Building backend..."
(
  cd backend
  go mod tidy
  go build -o ../che-url-shortener-server .
)

# Launch the compiled backend server as a background process
echo "Starting backend server..."
./che-url-shortener-server &

echo "Server started successfully in the background."
echo "Access the application at http://localhost:8080"
