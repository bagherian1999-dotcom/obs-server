# GoSync Server - Setup Complete! 🎉

Your GoSync server is now ready for easy deployment with IP/domain configuration and SSL support!

## What's Been Added

### 1. Configuration System
- ✅ `server/config.go` - Configuration loader
- ✅ `server/config.json.example` - Example configuration file
- ✅ Updated `server/main.go` - Now uses configuration system

### 2. Deployment Scripts
- ✅ `deploy.sh` - Interactive deployment script
- ✅ `install-as-service.sh` - Install as systemd service
- ✅ `setup-ssl.sh` - SSL certificate setup helper

### 3. Docker Support
- ✅ `Dockerfile` - Multi-stage Docker build
- ✅ `docker-compose.yml` - Docker Compose configuration

### 4. Documentation
- ✅ `DEPLOYMENT.md` - Comprehensive deployment guide
- ✅ `QUICKSTART.md` - Quick start guide
- ✅ Updated `.gitignore` - Excludes config and data files

## Quick Start

### For Local Testing (No SSL)
```bash
./deploy.sh
```
- Host: `0.0.0.0`
- Port: `8080`
- SSL: `n`

### For Production (With SSL)
```bash
./deploy.sh
./setup-ssl.sh
```

### As System Service
```bash
sudo ./install-as-service.sh
```

## Configuration File Structure

The server now uses `server/config.json`:

```json
{
  "host": "0.0.0.0",           // Bind address (0.0.0.0 = all interfaces)
  "port": "8080",              // Port number
  "domain": "sync.example.com", // Your domain (optional)
  "enable_ssl": true,          // Enable HTTPS
  "ssl_cert_path": "/path/to/cert.pem",
  "ssl_key_path": "/path/to/key.pem",
  "data_dir": "./data"         // Data directory
}
```

## Features Added

### ✅ Configurable IP/Domain
- Bind to specific IP or all interfaces
- Custom port configuration
- Domain name support for display

### ✅ SSL/HTTPS Support
- Native SSL support in the server
- Let's Encrypt integration helper
- Self-signed certificate generator
- Automatic URL scheme detection (http/https)

### ✅ Easy Deployment
- One-command interactive setup
- Systemd service installer
- Docker/Docker Compose support
- Background process management

### ✅ Production Ready
- Service auto-restart on failure
- Proper logging
- Security best practices
- Firewall configuration guidance

## File Structure

```
GoSync/
├── deploy.sh                    # Main deployment script
├── install-as-service.sh        # Service installer
├── setup-ssl.sh                 # SSL setup helper
├── Dockerfile                   # Docker image
├── docker-compose.yml           # Docker Compose config
├── DEPLOYMENT.md                # Detailed deployment guide
├── QUICKSTART.md                # Quick start guide
├── README.md                    # Original README
├── server/
│   ├── main.go                  # Updated main server (with config support)
│   ├── config.go                # Configuration loader
│   ├── config.json.example      # Example configuration
│   └── go.mod                   # Go dependencies
└── ssl/                         # SSL certificates (created by setup-ssl.sh)
```

## Next Steps

1. **Choose Your Deployment Method:**
   - Quick test: Run `./deploy.sh`
   - Production: Run `./deploy.sh` then `./setup-ssl.sh`
   - System service: Run `sudo ./install-as-service.sh`
   - Docker: Use `docker-compose up -d`

2. **Configure Firewall:**
   ```bash
   sudo ufw allow 8080/tcp    # For HTTP on port 8080
   sudo ufw allow 443/tcp     # For HTTPS on port 443
   ```

3. **Test the Server:**
   ```bash
   curl http://localhost:8080/api/status
   ```

4. **Access Dashboard:**
   Open `http://YOUR_SERVER_IP:8080` in browser

5. **Configure Obsidian Plugin:**
   - Server URL: Your server address
   - Device Name: Friendly name for this device
   - Enable Sync: Toggle ON

## Deployment Examples

### Example 1: Home Server (Local Network)
```bash
./deploy.sh
# Host: 0.0.0.0, Port: 8080, SSL: No
# Access: http://192.168.1.100:8080
```

### Example 2: VPS with Domain
```bash
./deploy.sh
# Host: 0.0.0.0, Port: 443, Domain: sync.example.com, SSL: Yes

./setup-ssl.sh
# Choose Let's Encrypt

sudo ./install-as-service.sh
# Install as service for auto-start

# Access: https://sync.example.com
```

### Example 3: Docker Deployment
```bash
# Create config.json
cp server/config.json.example server/config.json
# Edit server/config.json with your settings

# Run with Docker Compose
docker-compose up -d

# Access: http://localhost:8080
```

## Troubleshooting

### "Port already in use"
```bash
sudo lsof -i :8080
# Kill the process or choose different port
```

### "Permission denied" on SSL files
```bash
sudo chown $USER:$USER /path/to/cert.pem /path/to/key.pem
```

### Can't connect from other devices
- Check firewall: `sudo ufw status`
- Verify binding to `0.0.0.0` in config
- Check router port forwarding

### SSL certificate errors
```bash
# Verify certificates
openssl x509 -in cert.pem -text -noout
```

## Security Checklist

- [ ] SSL enabled for production
- [ ] Firewall configured
- [ ] Strong certificates (Let's Encrypt or CA-signed)
- [ ] Regular backups of `data/` directory
- [ ] Server kept updated
- [ ] Logs monitored regularly

## Support & Documentation

- **Quick Start:** See [QUICKSTART.md](QUICKSTART.md)
- **Detailed Guide:** See [DEPLOYMENT.md](DEPLOYMENT.md)
- **Original README:** See [README.md](README.md)

## Important Notes

⚠️ **Before Production Deployment:**
1. Always use SSL/HTTPS for production
2. Use Let's Encrypt for proper certificates
3. Set up automatic certificate renewal
4. Configure proper firewall rules
5. Regular backup of data directory

✅ **Configuration is Complete:**
- Server supports configurable IP/domain
- SSL/HTTPS fully supported
- Easy deployment scripts ready
- Production-ready setup available

🚀 **You're Ready to Deploy!**

Start with `./deploy.sh` and follow the prompts!
