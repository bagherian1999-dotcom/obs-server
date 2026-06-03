# GoSync - Quick Start Guide for Server Deployment

## Prerequisites

- Go 1.16 or higher installed
- A Linux/Unix server (VPS, cloud instance, or local machine)
- (Optional) A domain name pointing to your server
- (Optional) Port 80/443 access for SSL certificates

## 5-Minute Setup

### Step 1: Clone and Navigate

```bash
git clone <your-repo-url>
cd GoSync
```

### Step 2: Run Deployment Script

```bash
chmod +x deploy.sh
./deploy.sh
```

Answer the prompts:
- Server IP: `0.0.0.0` (to accept connections from all interfaces)
- Port: `8080` (or `443` if using SSL)
- Domain: Your domain name or leave empty for IP access
- Enable SSL: `y` or `n`

The script will:
✅ Install dependencies
✅ Create configuration file
✅ Build the server binary
✅ Optionally generate SSL certificates

### Step 3: Access Your Server

Open in browser:
- **Without SSL:** `http://YOUR_SERVER_IP:8080`
- **With SSL:** `https://YOUR_DOMAIN:443`

## Common Deployment Scenarios

### Scenario A: Local Network (Home/Office)

```bash
./deploy.sh
```
- Host: `0.0.0.0`
- Port: `8080`
- Domain: (leave empty)
- SSL: `n`

Access from devices on same network: `http://192.168.x.x:8080`

### Scenario B: Cloud Server with Domain

```bash
./deploy.sh
```
- Host: `0.0.0.0`
- Port: `443`
- Domain: `sync.yourdomain.com`
- SSL: `y`

Then run SSL setup:
```bash
./setup-ssl.sh
```
Choose Let's Encrypt option.

Access: `https://sync.yourdomain.com`

### Scenario C: Run as Background Service

After deployment:

```bash
sudo ./install-as-service.sh
```

This sets up systemd service that:
- Starts automatically on boot
- Restarts on failure
- Runs in background

Control with:
```bash
sudo systemctl start gosync
sudo systemctl stop gosync
sudo systemctl status gosync
```

## Configure Obsidian Plugin

1. Open Obsidian Settings → Community Plugins
2. Find "GoSync" plugin
3. Configure:
   - **Server URL:** `https://sync.yourdomain.com` or `http://your-ip:8080`
   - **Device Name:** A friendly name (e.g., "My Laptop")
   - **Enable Sync:** Toggle ON

## Quick Commands Reference

```bash
# Start server manually
cd server && ./gosync-server

# Start in background
cd server && nohup ./gosync-server > server.log 2>&1 &

# Check if running
ps aux | grep gosync-server

# Stop background server
pkill gosync-server

# View logs (if using systemd)
sudo journalctl -u gosync -f

# Check configuration
cat server/config.json

# Test server access
curl http://localhost:8080/api/status
```

## Troubleshooting

**Server won't start:**
```bash
# Check if port is in use
sudo lsof -i :8080

# Check logs
cat server/server.log
```

**Can't connect from other devices:**
- Verify firewall: `sudo ufw allow 8080`
- Verify binding to `0.0.0.0` not `127.0.0.1`
- Check router port forwarding (if applicable)

**SSL certificate errors:**
```bash
# Verify certificates exist
ls -la /etc/letsencrypt/live/yourdomain.com/
```

## Next Steps

- Read [DEPLOYMENT.md](DEPLOYMENT.md) for detailed configuration options
- Configure automatic backups of your `data/` directory
- Set up SSL certificate auto-renewal for Let's Encrypt
- Configure reverse proxy (Nginx/Apache) for advanced setups

## Security Notes

⚠️ **For production use:**
- Always enable SSL
- Use strong, CA-signed certificates (Let's Encrypt)
- Keep the server updated
- Regularly backup the `data/` directory
- Use firewall rules to restrict access
