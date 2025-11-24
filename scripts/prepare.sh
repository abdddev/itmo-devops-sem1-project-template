#!/bin/bash
echo "Building Docker images"

cd "$(dirname "$0")/.."

docker compose build

echo "Docker images built successfully"