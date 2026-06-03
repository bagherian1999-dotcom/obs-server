#!/bin/bash

# SSL Certificate Setup Helper for GoSync
# Supports both Let's Encrypt and self-signed certificates

set -e

echo "=== GoSync SSL Certificate Setup ==="
echo ""
echo "Choose SSL certificate method:"
echo "1) Let's Encrypt (recommended for production with domain)"
echo "2) Self-signed certificate (for testing/internal use)"
echo ""
read -p "Enter choice (1 or 2): " CERT_METHOD

read -p "Enter your domain name: " DOMAIN

if [[ -z "$DOMAIN" ]]; then
    echo "Error: Domain name is required for SSL setup"
    exit 1
fi

SSL_DIR="./ssl"
mkdir -p "$SSL_DIR"

if [[ "$CERT_METHOD" == "1" ]]; then
    echo ""
    echo "=== Let's Encrypt Setup ==="
    echo ""
    
    # Check if certbot is installed
    if ! command -v certbot &> /dev/null; then
        echo "Certbot is not installed."
        echo "Install it with:"
        echo "  Ubuntu/Debian: sudo apt install certbot"
        echo "  CentOS/RHEL: sudo yum install certbot"
        echo "  macOS: brew install certbot"
        exit 1
    fi
    
    read -p "Enter email for Let's Encrypt notifications: " EMAIL
    
    if [[ -z "$EMAIL" ]]; then
        echo "Error: Email is required for Let's Encrypt"
        exit 1
    fi
    
    echo ""
    echo "Running certbot in standalone mode..."
    echo "Make sure port 80 is open and not in use!"
    echo ""
    
    sudo certbot certonly --standalone \
        --preferred-challenges http \
        --email "$EMAIL" \
        --agree-tos \
        --no-eff-email \
        -d "$DOMAIN"
    
    # Copy certificates to ssl directory
    sudo cp "/etc/letsencrypt/live/$DOMAIN/fullchain.pem" "$SSL_DIR/cert.pem"
    sudo cp "/etc/letsencrypt/live/$DOMAIN/privkey.pem" "$SSL_DIR/key.pem"
    sudo chown $USER:$USER "$SSL_DIR/cert.pem" "$SSL_DIR/key.pem"
    
    echo ""
    echo "Let's Encrypt certificates installed successfully!"
    echo "Certificate: $SSL_DIR/cert.pem"
    echo "Key: $SSL_DIR/key.pem"
    echo ""
    echo "NOTE: Let's Encrypt certificates expire in 90 days."
    echo "Set up auto-renewal with: sudo certbot renew --dry-run"
    
elif [[ "$CERT_METHOD" == "2" ]]; then
    echo ""
    echo "=== Self-Signed Certificate Setup ==="
    echo ""
    
    CERT_FILE="$SSL_DIR/cert.pem"
    KEY_FILE="$SSL_DIR/key.pem"
    
    read -p "Certificate validity in days (default: 365): " DAYS
    DAYS=${DAYS:-365}
    
    echo "Generating self-signed certificate..."
    openssl req -x509 -nodes -days "$DAYS" -newkey rsa:2048 \
        -keyout "$KEY_FILE" \
        -out "$CERT_FILE" \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=$DOMAIN"
    
    echo ""
    echo "Self-signed certificate generated successfully!"
    echo "Certificate: $CERT_FILE"
    echo "Key: $KEY_FILE"
    echo ""
    echo "WARNING: Self-signed certificates will show security warnings in browsers."
    echo "Use Let's Encrypt for production environments."
    
else
    echo "Invalid choice"
    exit 1
fi

echo ""
echo "=== Update Configuration ==="
echo ""
echo "Update your server/config.json with:"
echo "  \"enable_ssl\": true,"
echo "  \"ssl_cert_path\": \"$SSL_DIR/cert.pem\","
echo "  \"ssl_key_path\": \"$SSL_DIR/key.pem\""
echo ""

read -p "Update config.json automatically? (y/n): " UPDATE_CONFIG
if [[ "$UPDATE_CONFIG" == "y" || "$UPDATE_CONFIG" == "Y" ]]; then
    if [[ -f "server/config.json" ]]; then
        # Backup existing config
        cp server/config.json server/config.json.backup
        
        # Update config using jq if available, otherwise manual
        if command -v jq &> /dev/null; then
            jq ".enable_ssl = true | .ssl_cert_path = \"$SSL_DIR/cert.pem\" | .ssl_key_path = \"$SSL_DIR/key.pem\" | .domain = \"$DOMAIN\"" \
                server/config.json > server/config.json.tmp && mv server/config.json.tmp server/config.json
            echo "Configuration updated successfully!"
        else
            echo "jq not found. Please update config.json manually."
        fi
    else
        echo "config.json not found. Please run deploy.sh first."
    fi
fi

echo ""
echo "SSL setup complete!"
