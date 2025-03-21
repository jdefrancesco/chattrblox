#!/bin/zsh

set -e

echo "[+] Starting Chatrblox Dev Environment..."

# --- CONFIG ---
REDIS_PORT=6379
GO_SERVER_PORT=8080
REDIS_CONTAINER_NAME=chatrblox-redis
GO_MAIN="./cmd/main.go"
LOG_DIR="./logs"
mkdir -p "$LOG_DIR"

# --- REDIS ---
if ! docker ps --format '{{.Names}}' | grep -q "^${REDIS_CONTAINER_NAME}$"; then
  echo "[+] Starting Redis in Docker..."
  docker run -d \
    --name $REDIS_CONTAINER_NAME \
    -p $REDIS_PORT:6379 \
    redis:7 > /dev/null
else
  echo "[!] Redis already running in Docker"
fi

# --- GO BACKEND ---
echo "[+] Starting Go backend on :${GO_SERVER_PORT}..."
go run "$GO_MAIN" 2>&1 | tee "$LOG_DIR/server.log"