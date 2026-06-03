# 🚀 GoSync Server - Start Here

## Your Server is Ready for Deployment!

This GoSync server has been configured for easy deployment with IP/domain configuration and SSL support.

---

## ⚡ Quick Start (Choose One)

### Option 1: Interactive Setup (Recommended)
```bash
./deploy.sh
```
Follow the prompts to configure your server.

### Option 2: Docker
```bash
cp server/config.json.example server/config.json
# Edit server/config.json with your settings
docker-compose up -d
```

### Option 3: Manual
```bash
cd server
cp config.json.example config.json
# Edit config.json
go mod tidy
go build -o gosync-server main.go config.go
./gosync-server
```

---

## 📚 Documentation

- **[QUICKSTART.md](QUICKSTART.md)** - 5-minute deployment guide
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Comprehensive guide with SSL, systemd, Docker
- **[SETUP_COMPLETE.md](SETUP_COMPLETE.md)** - What's new and how to use it

---

## ✨ What's Included

### Configuration System
- JSON-based configuration (`server/config.json`)
- Configurable IP/domain, port, SSL, data directory
- Example config included

### SSL/HTTPS Support
- Native SSL in Go server
- Let's Encrypt helper script
- Self-signed certificate generator

### Deployment Tools
- **`deploy.sh`** - Interactive deployment wizard
- **`setup-ssl.sh`** - SSL certificate setup
- **`install-as-service.sh`** - Install as systemd service

### Docker Support
- Multi-stage Dockerfile
- Docker Compose configuration
- Volume management for data

---

## 🔧 Configuration Example

Create `server/config.json`:
```json
{
  "host": "0.0.0.0",
  "port": "8080",
  "domain": "sync.yourdomain.com",
  "enable_ssl": false,
  "ssl_cert_path": "",
  "ssl_key_path": "",
  "data_dir": "./data"
}
```

---

## 🌐 Common Scenarios

### Local Network (No SSL)
```bash
./deploy.sh
# Host: 0.0.0.0, Port: 8080, SSL: no
```
Access: `http://192.168.x.x:8080`

### Production with SSL
```bash
./deploy.sh
# Host: 0.0.0.0, Port: 443, Domain: your.domain, SSL: yes
./setup-ssl.sh
# Choose Let's Encrypt
```
Access: `https://your.domain`

### System Service
```bash
./deploy.sh
sudo ./install-as-service.sh
sudo systemctl enable --now gosync
```

---

## 📦 Files Created

```
GoSync/
├── deploy.sh                    # Main deployment script
├── setup-ssl.sh                 # SSL setup helper
├── install-as-service.sh        # Service installer
├── Dockerfile                   # Docker image
├── docker-compose.yml           # Docker Compose
├── QUICKSTART.md                # Quick guide
├── DEPLOYMENT.md                # Full guide
├── SETUP_COMPLETE.md            # Summary
└── server/
    ├── config.go                # Config loader
    ├── config.json.example      # Config template
    └── main.go                  # Updated server
```

---

## 🔒 Security Notes

**For Production:**
- ✅ Enable SSL/HTTPS
- ✅ Use Let's Encrypt or CA-signed certificates
- ✅ Configure firewall (ports 80/443)
- ✅ Regular backups of data directory
- ✅ Set up certificate auto-renewal

---

## 🆘 Quick Help

**Port already in use:**
```bash
sudo lsof -i :8080
```

**Check if server is running:**
```bash
curl http://localhost:8080/api/status
```

**View service logs:**
```bash
sudo journalctl -u gosync -f
```

---

## 🎯 Next Steps

1. **Deploy:** Run `./deploy.sh`
2. **Configure SSL:** Run `./setup-ssl.sh` (if needed)
3. **Install Service:** Run `sudo ./install-as-service.sh` (optional)
4. **Access Dashboard:** Open `http://YOUR_SERVER_IP:8080`
5. **Configure Plugin:** Set server URL in Obsidian plugin settings

---

**Need detailed instructions?** Read [QUICKSTART.md](QUICKSTART.md) or [DEPLOYMENT.md](DEPLOYMENT.md)

**Ready to deploy?**
```bash
./deploy.sh
```
