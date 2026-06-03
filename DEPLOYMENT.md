# GoSync Server - Deployment Guide

This guide will help you deploy GoSync server on your production server with IP/domain configuration and SSL support.

## Quick Start

### Option 1: Interactive Setup (Recommended)

Run the deployment script which will guide you through the setup:

```bash
./deploy.sh
```

This script will:
- Ask for your server IP/domain configuration
- Configure SSL if needed
- Install dependencies
- Build the server binary
- Create configuration file
- Optionally start the server

### Option 2: Manual Setup

1. **Create Configuration File**

Copy the example configuration:
```bash
cp server/config.json.example server/config.json
```

Edit `server/config.json`:
```json
{
  "host": "0.0.0.0",
  "port": "8080",
  "domain": "sync.yourdomain.com",
  "enable_ssl": true,
  "ssl_cert_path": "/etc/ssl/certs/gosync.crt",
  "ssl_key_path": "/etc/ssl/private/gosync.key",
  "data_dir": "./data"
}
```

2. **Build the Server**

```bash
cd server
go mod tidy
go build -o gosync-server main.go config.go
```

3. **Run the Server**

```bash
./gosync-server
```

## Configuration Options

### config.json Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `host` | IP address to bind to (use 0.0.0.0 for all interfaces) | `"0.0.0.0"` |
| `port` | Port number to listen on | `"8080"` or `"443"` |
| `domain` | Your domain name (optional, for display purposes) | `"sync.example.com"` |
| `enable_ssl` | Enable HTTPS/SSL | `true` or `false` |
| `ssl_cert_path` | Path to SSL certificate file | `"/etc/ssl/certs/cert.pem"` |
| `ssl_key_path` | Path to SSL private key file | `"/etc/ssl/private/key.pem"` |
| `data_dir` | Directory to store synchronized files | `"./data"` |

## SSL Setup

### Option 1: Let's Encrypt (Recommended for Production)

Use the SSL setup script:

```bash
./setup-ssl.sh
```

Choose option 1 for Let's Encrypt and follow the prompts.

**Manual Let's Encrypt Setup:**

```bash
# Install certbot
sudo apt install certbot  # Ubuntu/Debian
# or
sudo yum install certbot  # CentOS/RHEL

# Obtain certificate
sudo certbot certonly --standalone -d sync.yourdomain.com

# Certificates will be at:
# /etc/letsencrypt/live/sync.yourdomain.com/fullchain.pem
# /etc/letsencrypt/live/sync.yourdomain.com/privkey.pem
```

Update your `config.json`:
```json
{
  "enable_ssl": true,
  "ssl_cert_path": "/etc/letsencrypt/live/sync.yourdomain.com/fullchain.pem",
  "ssl_key_path": "/etc/letsencrypt/live/sync.yourdomain.com/privkey.pem"
}
```

**Note:** Let's Encrypt certificates expire in 90 days. Set up auto-renewal:
```bash
sudo certbot renew --dry-run
```

### Option 2: Self-Signed Certificate (Testing/Internal Use)

```bash
./setup-ssl.sh
```

Choose option 2 for self-signed certificate.

**Manual Self-Signed Certificate:**

```bash
mkdir -p ssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout ssl/key.pem \
  -out ssl/cert.pem \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=yourdomain.com"
```

Update your `config.json`:
```json
{
  "enable_ssl": true,
  "ssl_cert_path": "./ssl/cert.pem",
  "ssl_key_path": "./ssl/key.pem"
}
```

**Warning:** Self-signed certificates will show security warnings in browsers and require clients to accept the certificate.

## Running as a System Service

### Linux (systemd)

Use the service installation script:

```bash
sudo ./install-as-service.sh
```

This will:
- Install GoSync to `/opt/gosync` (or your chosen directory)
- Create a systemd service
- Set up auto-start on boot (optional)

**Manual Service Setup:**

Create `/etc/systemd/system/gosync.service`:

```ini
[Unit]
Description=GoSync Server
After=network.target

[Service]
Type=simple
User=yourusername
WorkingDirectory=/opt/gosync
ExecStart=/opt/gosync/gosync-server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable gosync
sudo systemctl start gosync
sudo systemctl status gosync
```

**Service Commands:**
```bash
sudo systemctl start gosync    # Start service
sudo systemctl stop gosync     # Stop service
sudo systemctl restart gosync  # Restart service
sudo systemctl status gosync   # Check status
sudo journalctl -u gosync -f   # View logs
```

## Docker Deployment

### Using Docker Compose

1. **Setup Configuration:**

Create `server/config.json` with your settings.

2. **Run with Docker Compose:**

```bash
docker-compose up -d
```

3. **View Logs:**

```bash
docker-compose logs -f
```

4. **Stop:**

```bash
docker-compose down
```

### Manual Docker Build

```bash
# Build image
docker build -t gosync-server .

# Run container
docker run -d \
  --name gosync \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/server/config.json:/app/config.json \
  gosync-server
```

## Network Configuration

### Firewall Rules

**Allow HTTP (port 80):**
```bash
sudo ufw allow 80/tcp
```

**Allow HTTPS (port 443):**
```bash
sudo ufw allow 443/tcp
```

**Allow custom port (e.g., 8080):**
```bash
sudo ufw allow 8080/tcp
```

### Reverse Proxy (Nginx)

If you want to run GoSync behind Nginx:

```nginx
server {
    listen 80;
    server_name sync.yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name sync.yourdomain.com;
    
    ssl_certificate /etc/letsencrypt/live/sync.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/sync.yourdomain.com/privkey.pem;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

In this case, set your `config.json` to:
```json
{
  "host": "127.0.0.1",
  "port": "8080",
  "enable_ssl": false
}
```

## Accessing the Server

### From Web Browser

- **HTTP:** `http://your-server-ip:8080`
- **HTTPS:** `https://your-domain.com:443`
- **Custom Port:** `http://your-domain.com:8080`

### From Obsidian Plugin

Configure the Obsidian plugin with:
- **Server URL:** `https://sync.yourdomain.com` (or your configured URL)
- **Device Name:** A friendly name for your device

## Troubleshooting

### Port Already in Use

Check what's using the port:
```bash
sudo lsof -i :8080
```

Kill the process or choose a different port in `config.json`.

### SSL Certificate Errors

**Verify certificate files exist:**
```bash
ls -la /path/to/cert.pem
ls -la /path/to/key.pem
```

**Check certificate validity:**
```bash
openssl x509 -in /path/to/cert.pem -text -noout
```

**Check certificate and key match:**
```bash
openssl x509 -noout -modulus -in /path/to/cert.pem | openssl md5
openssl rsa -noout -modulus -in /path/to/key.pem | openssl md5
```

### Permission Denied

Ensure the user running the server has:
- Read access to SSL certificate files
- Write access to the data directory
- Permission to bind to the port (ports < 1024 require root/sudo)

### Cannot Connect from Other Devices

1. Check firewall rules
2. Verify server is binding to `0.0.0.0` not `127.0.0.1`
3. Ensure port forwarding is configured on router (if needed)
4. Check that domain DNS points to your server IP

## Security Best Practices

1. **Use SSL in Production:** Always enable SSL for production deployments
2. **Firewall:** Only open necessary ports
3. **Regular Updates:** Keep Go and dependencies updated
4. **Strong Certificates:** Use Let's Encrypt or proper CA-signed certificates
5. **Data Backup:** Regularly backup the `data` directory
6. **Monitor Logs:** Check logs regularly for suspicious activity
7. **Limit Access:** Use firewall rules to restrict access if only needed from specific IPs

## Updating the Server

1. **Backup your data:**
```bash
tar -czf backup-$(date +%Y%m%d).tar.gz data/ server/config.json
```

2. **Pull latest changes:**
```bash
git pull origin main
```

3. **Rebuild:**
```bash
cd server
go build -o gosync-server main.go config.go
```

4. **Restart service:**
```bash
sudo systemctl restart gosync
```

## Support

For issues and questions:
- Check the [README.md](README.md) for general information
- Review server logs: `sudo journalctl -u gosync -f`
- Check the web dashboard at `http://your-server/`

## Examples

### Example 1: Local Network (No SSL)

```json
{
  "host": "0.0.0.0",
  "port": "8080",
  "domain": "",
  "enable_ssl": false,
  "data_dir": "./data"
}
```

Access: `http://192.168.1.100:8080`

### Example 2: Production with Let's Encrypt

```json
{
  "host": "0.0.0.0",
  "port": "443",
  "domain": "sync.example.com",
  "enable_ssl": true,
  "ssl_cert_path": "/etc/letsencrypt/live/sync.example.com/fullchain.pem",
  "ssl_key_path": "/etc/letsencrypt/live/sync.example.com/privkey.pem",
  "data_dir": "/var/lib/gosync/data"
}
```

Access: `https://sync.example.com`

### Example 3: Behind Nginx Reverse Proxy

```json
{
  "host": "127.0.0.1",
  "port": "8080",
  "domain": "sync.example.com",
  "enable_ssl": false,
  "data_dir": "./data"
}
```

Nginx handles SSL, forwards to `http://127.0.0.1:8080`

Access: `https://sync.example.com` (via Nginx)
