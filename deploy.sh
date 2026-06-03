#!/bin/bash

# GoSync Server Deployment Script
# This script helps you deploy GoSync server with SSL support

set -e

echo "=== GoSync Server Deployment Setup ==="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go from https://go.dev/"
    exit 1
fi

# Configuration
read -p "Enter server IP or domain (default: 0.0.0.0): " HOST
HOST=${HOST:-0.0.0.0}

read -p "Enter server port (default: 8080): " PORT
PORT=${PORT:-8080}

read -p "Enter domain name (leave empty for IP-based access): " DOMAIN

read -p "Enable SSL? (y/n, default: n): " ENABLE_SSL
ENABLE_SSL=${ENABLE_SSL:-n}

if [[ "$ENABLE_SSL" == "y" || "$ENABLE_SSL" == "Y" ]]; then
    ENABLE_SSL_BOOL="true"
    
    read -p "SSL certificate path (default: /etc/ssl/certs/gosync.crt): " SSL_CERT
    SSL_CERT=${SSL_CERT:-/etc/ssl/certs/gosync.crt}
    
    read -p "SSL key path (default: /etc/ssl/private/gosync.key): " SSL_KEY
    SSL_KEY=${SSL_KEY:-/etc/ssl/private/gosync.key}
else
    ENABLE_SSL_BOOL="false"
    SSL_CERT="/etc/ssl/certs/gosync.crt"
    SSL_KEY="/etc/ssl/private/gosync.key"
fi

read -p "Data directory (default: ./data): " DATA_DIR
DATA_DIR=${DATA_DIR:-./data}

# Navigate to server directory
cd server

echo ""
echo "=== Installing Dependencies ==="
go mod tidy

echo ""
echo "=== Creating Configuration File ==="
cat > config.json << CONFIGEOF
{
  "host": "$HOST",
  "port": "$PORT",
  "domain": "$DOMAIN",
  "enable_ssl": $ENABLE_SSL_BOOL,
  "ssl_cert_path": "$SSL_CERT",
  "ssl_key_path": "$SSL_KEY",
  "data_dir": "$DATA_DIR"
}
CONFIGEOF

echo "Configuration saved to server/config.json"

echo ""
echo "=== Building Server Binary ==="
go build -o gosync-server main.go config.go

echo ""
echo "=== Setup Complete! ==="
echo ""
echo "Configuration:"
echo "  Host: $HOST"
echo "  Port: $PORT"
echo "  Domain: $DOMAIN"
echo "  SSL Enabled: $ENABLE_SSL_BOOL"
if [[ "$ENABLE_SSL_BOOL" == "true" ]]; then
    echo "  SSL Certificate: $SSL_CERT"
    echo "  SSL Key: $SSL_KEY"
fi
echo "  Data Directory: $DATA_DIR"
echo ""

if [[ "$ENABLE_SSL_BOOL" == "true" ]]; then
    if [[ ! -f "$SSL_CERT" ]] || [[ ! -f "$SSL_KEY" ]]; then
        echo "WARNING: SSL is enabled but certificate files not found!"
        echo "You need to obtain SSL certificates. Options:"
        echo "  1. Use Let's Encrypt (recommended for production)"
        echo "  2. Generate self-signed certificates (for testing)"
        echo ""
        read -p "Generate self-signed certificates now? (y/n): " GEN_CERT
        if [[ "$GEN_CERT" == "y" || "$GEN_CERT" == "Y" ]]; then
            echo "Generating self-signed certificates..."
            mkdir -p $(dirname "$SSL_CERT") $(dirname "$SSL_KEY")
            
            CERT_DOMAIN=${DOMAIN:-localhost}
            openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
                -keyout "$SSL_KEY" \
                -out "$SSL_CERT" \
                -subj "/C=US/ST=State/L=City/O=Organization/CN=$CERT_DOMAIN"
            
            echo "Self-signed certificates generated!"
            echo "NOTE: Self-signed certificates will show browser warnings."
        fi
    fi
fi

echo ""
echo "To start the server:"
echo "  cd server"
echo "  ./gosync-server"
echo ""
echo "Or run in background:"
echo "  cd server"
echo "  nohup ./gosync-server > server.log 2>&1 &"
echo ""

if [[ "$ENABLE_SSL_BOOL" == "true" ]]; then
    SCHEME="https"
else
    SCHEME="http"
fi

if [[ -n "$DOMAIN" ]]; then
    echo "Access your server at: $SCHEME://$DOMAIN:$PORT"
else
    echo "Access your server at: $SCHEME://localhost:$PORT"
    echo "Or from other devices: $SCHEME://YOUR_SERVER_IP:$PORT"
fi

echo ""
read -p "Start the server now? (y/n): " START_NOW
if [[ "$START_NOW" == "y" || "$START_NOW" == "Y" ]]; then
    echo "Starting GoSync Server..."
    ./gosync-server
fi
