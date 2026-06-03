# GoSync - Server Deployment Guide

## 🚀 Ready for Production Deployment!

Your GoSync server now supports:
- ✅ **Configurable IP/Domain** - Easy network setup
- ✅ **SSL/HTTPS Support** - Secure connections
- ✅ **One-Command Deployment** - Interactive setup script
- ✅ **System Service** - Auto-start on boot
- ✅ **Docker Support** - Containerized deployment

---

## Quick Deploy (3 Steps)

### 1. Run Deployment Script
```bash
./deploy.sh
```

### 2. (Optional) Setup SSL
```bash
./setup-ssl.sh
```

### 3. Access Your Server
```
http://YOUR_SERVER_IP:8080
```

---

## What's New

### Configuration System
The server now uses `server/config.json` for all settings:
- Host/IP binding
- Port configuration  
- Domain name
- SSL certificate paths
- Data directory location

### Deployment Scripts
- **`deploy.sh`** - Interactive setup wizard
- **`setup-ssl.sh`** - SSL certificate helper (Let's Encrypt or self-signed)
- **`install-as-service.sh`** - Install as systemd service

### Files Added
```
server/config.go              # Configuration loader
server/config.json.example    # Configuration template
server/main.go               # Updated with config support
Dockerfile                   # Docker image
docker-compose.yml           # Docker Compose config
DEPLOYMENT.md               # Detailed deployment guide
QUICKSTART.md              # Quick start guide
```

---

## Deployment Options

### Option 1: Standard Deployment
```bash
./deploy.sh
cd server && ./gosync-server
```

### Option 2: Background Service
```bash
./deploy.sh
sudo ./install-as-service.sh
sudo systemctl start gosync
```

### Option 3: Docker
```bash
docker-compose up -d
```

---

## Configuration Examples

### Local Network (No SSL)
```json
{
  "host": "0.0.0.0",
  "port": "8080",
  "domain": "",
  "enable_ssl": false,
  "data_dir": "./data"
}
```

### Production with SSL
```json
{
  "host": "0.0.0.0",
  "port": "443",
  "domain": "sync.yourdomain.com",
  "enable_ssl": true,
  "ssl_cert_path": "/etc/letsencrypt/live/sync.yourdomain.com/fullchain.pem",
  "ssl_key_path": "/etc/letsencrypt/live/sync.yourdomain.com/privkey.pem",
  "data_dir": "./data"
}
```

---

## Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - Get started in 5 minutes
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Comprehensive deployment guide
- **[SETUP_COMPLETE.md](SETUP_COMPLETE.md)** - Summary of changes

---

## Support

For detailed instructions, troubleshooting, and best practices, see the comprehensive guides in:
- DEPLOYMENT.md
- QUICKSTART.md

---

**Ready to deploy? Start with:**
```bash
./deploy.sh
```
