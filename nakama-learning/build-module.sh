#!/bin/bash

# Build Nakama Go module
echo "🔨 Building Nakama Go module..."

# Build the module inside Docker (Nakama runs on Linux)
docker run --rm \
  -v $(pwd)/modules:/workspace \
  -w /workspace \
  golang:1.21 \
  bash -c "go mod tidy && go build -buildmode=plugin -o main.so main.go"

# Copy to Nakama's module directory
cp modules/main.so data/modules/

echo "✅ Module built and copied to data/modules/"
echo "🔄 Restarting Nakama to load the module..."

# Restart Nakama to load the new module
docker compose restart nakama

echo "✅ Done! Check logs with: docker compose logs -f nakama"
