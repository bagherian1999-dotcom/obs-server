# Authentication System for GoSync Server

## Overview

This implementation adds complete user authentication and authorization to the OBS sync server. Each user has isolated storage, ensuring complete data privacy and security.

## Features

### 🔐 Authentication
- **User Registration**: New users can create accounts with username, email, and password
- **User Login**: Existing users authenticate with username and password
- **JWT Tokens**: Secure token-based authentication (24-hour expiration)
- **Password Security**: Passwords are hashed using bcrypt (cost factor 10)

### 👤 User Isolation
- **Private Data Storage**: Each user has a separate data directory (`data/users/{userID}/`)
- **Isolated File Metadata**: Per-user metadata.json files
- **Private File Access**: Users can only access their own files
- **Separate WebSocket Connections**: Client connections are user-specific

### 🛡️ Security Features
- All API endpoints (except `/api/register` and `/api/login`) require authentication
- JWT tokens must be provided in the `Authorization: Bearer {token}` header
- WebSocket connections require token in query parameter: `/ws?token={token}`
- CORS enabled for cross-origin requests
- Session management with token verification

## API Endpoints

### Public Endpoints (No Authentication Required)

#### POST /api/register
Register a new user account.

**Request:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "a1b2c3d4e5f6...",
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

**Validation:**
- Username: minimum 3 characters
- Password: minimum 6 characters
- Email: optional

#### POST /api/login
Authenticate an existing user.

**Request:**
```json
{
  "username": "john_doe",
  "password": "secure123"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "a1b2c3d4e5f6...",
    "username": "john_doe",
    "email": "john@example.com"
  }
}
```

### Protected Endpoints (Authentication Required)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer {your_jwt_token}
```

#### GET /api/profile
Get current user's profile information.

**Response:**
```json
{
  "id": "a1b2c3d4e5f6...",
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2025-06-03T10:30:00Z",
  "last_login": "2025-06-03T15:45:00Z"
}
```

#### GET /api/verify
Verify if the current token is valid.

**Response:**
```json
{
  "valid": true,
  "user": {
    "id": "a1b2c3d4e5f6...",
    "username": "john_doe"
  }
}
```

#### GET /api/files
List all files for the authenticated user.

**Response:**
```json
[
  {
    "path": "Notes/Daily/2025-06-03.md",
    "history": [...],
    "latest": {
      "hash": "abc123...",
      "timestamp": 1717423456,
      "device": "laptop"
    }
  }
]
```

#### GET /api/file?path={filepath}
Download a specific file.

**Query Parameters:**
- `path`: File path to retrieve

**Response:** File content (binary)

#### PUT /api/file?path={filepath}&device={device_name}
Upload or update a file.

**Query Parameters:**
- `path`: File path to save
- `device`: Device name (optional, defaults to "unknown")

**Request Body:** File content (binary)

**Response:**
```json
{
  "status": "ok",
  "hash": "abc123..."
}
```

#### POST /api/files/delete?patterns={pattern1,pattern2}
Delete files matching patterns.

**Query Parameters:**
- `patterns`: Comma-separated list of patterns to match

**Response:**
```json
{
  "deleted": 5
}
```

#### GET /api/search?q={query}
Search files by path.

**Query Parameters:**
- `q`: Search query

**Response:**
```json
[
  {
    "path": "matching/file.md",
    "history": [...],
    "latest": {...}
  }
]
```

#### POST /api/cleanup
Clean up orphaned file versions (files deleted from metadata but blobs still exist).

**Response:**
```json
{
  "deletedBlobs": 10,
  "freedBytes": 524288
}
```

#### POST /api/reset
Delete all user data (DANGEROUS - cannot be undone).

**Response:**
```json
{
  "status": "reset_complete"
}
```

#### GET /api/status
Get server status and user statistics.

**Response:**
```json
{
  "clients": [
    {
      "id": "123456789",
      "name": "laptop",
      "ip": "192.168.1.100:54321",
      "user_id": "a1b2c3d4...",
      "connectedAt": "2025-06-03T15:00:00Z"
    }
  ],
  "logs": [
    {
      "timestamp": "2025-06-03T15:30:00Z",
      "message": "File uploaded: Notes/test.md",
      "level": "INFO"
    }
  ],
  "fileCount": 42,
  "user": "john_doe"
}
```

### WebSocket Connection

WebSocket endpoint for real-time updates.

**URL:** `/ws?token={jwt_token}`

The token must be provided as a query parameter since WebSocket headers are limited.

**Messages:**
```json
{
  "type": "identify",
  "deviceName": "laptop",
  "pluginName": "obsidian-sync-plugin"
}
```

**Broadcasts:**
```json
{
  "type": "file_updated",
  "path": "Notes/test.md",
  "hash": "abc123..."
}
```

## Data Structure

### Directory Layout
```
.
├── config.json           # Server configuration
├── users.json           # User database
└── data/
    └── users/
        ├── {userID1}/
        │   ├── metadata.json
        │   ├── abc123...     # File blob (hash as filename)
        │   └── def456...
        └── {userID2}/
            ├── metadata.json
            └── ...
```

### users.json Format
```json
{
  "john_doe": {
    "id": "a1b2c3d4e5f6...",
    "username": "john_doe",
    "email": "john@example.com",
    "password_hash": "$2a$10$...",
    "created_at": "2025-06-03T10:00:00Z",
    "last_login": "2025-06-03T15:45:00Z"
  }
}
```

### User metadata.json Format
```json
{
  "Notes/test.md": {
    "path": "Notes/test.md",
    "history": [
      {
        "hash": "abc123...",
        "timestamp": 1717423456,
        "device": "laptop"
      }
    ],
    "latest": {
      "hash": "abc123...",
      "timestamp": 1717423456,
      "device": "laptop"
    }
  }
}
```

## Configuration

### Environment Variables (Recommended)

Create a `.env` file or set environment variables:

```bash
# JWT Secret (CHANGE THIS!)
export JWT_SECRET="your-very-secret-key-here"

# Server configuration
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="8080"
export DATA_DIR="./data"

# SSL (optional)
export ENABLE_SSL="false"
export SSL_CERT_PATH="/path/to/cert.pem"
export SSL_KEY_PATH="/path/to/key.pem"
```

### config.json

```json
{
  "host": "0.0.0.0",
  "port": "8080",
  "domain": "",
  "enable_ssl": false,
  "ssl_cert_path": "",
  "ssl_key_path": "",
  "data_dir": "./data"
}
```

## Security Considerations

### ⚠️ Important Security Notes

1. **Change JWT Secret**: The default `JWTSecret` in `auth.go` is **NOT SECURE**. Change it to a strong random value:
   ```go
   const JWTSecret = "your-very-long-random-secret-key-at-least-32-characters"
   ```
   Or better yet, load it from environment variables.

2. **Use HTTPS**: For production, enable SSL/TLS to encrypt all traffic:
   ```json
   {
     "enable_ssl": true,
     "ssl_cert_path": "/path/to/cert.pem",
     "ssl_key_path": "/path/to/key.pem"
   }
   ```

3. **Password Requirements**: Current minimum is 6 characters. Consider increasing for production.

4. **Rate Limiting**: Not implemented. Consider adding rate limiting for login/register endpoints.

5. **Token Refresh**: Tokens expire after 24 hours. Clients need to re-authenticate.

6. **File Permissions**: Data directory should have restricted permissions (0700 or 0755).

## Building and Running

### Prerequisites
- Go 1.21 or later

### Install Dependencies
```bash
cd server
go mod tidy
```

### Build
```bash
go build -o gosync-server
```

### Run
```bash
./gosync-server
```

The server will start on `http://localhost:8080` (or configured address).

### First-Time Setup

1. Open browser to `http://localhost:8080`
2. Click "Register" tab
3. Create a new account
4. You'll be automatically logged in

## Client Integration

### Example: Login and File Upload (JavaScript)

```javascript
// Login
const loginResponse = await fetch('http://localhost:8080/api/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    username: 'john_doe',
    password: 'secure123'
  })
});

const { token } = await loginResponse.json();

// Upload a file
const fileContent = 'Hello World!';
await fetch('http://localhost:8080/api/file?path=test.md&device=browser', {
  method: 'PUT',
  headers: {
    'Authorization': `Bearer ${token}`
  },
  body: fileContent
});

// List files
const filesResponse = await fetch('http://localhost:8080/api/files', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const files = await filesResponse.json();
console.log(files);
```

### WebSocket Connection

```javascript
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

ws.onopen = () => {
  // Identify the client
  ws.send(JSON.stringify({
    type: 'identify',
    deviceName: 'my-laptop',
    pluginName: 'obsidian-sync-v1.0'
  }));
};

ws.onmessage = (event) => {
  const msg = JSON.parse(event.data);
  if (msg.type === 'file_updated') {
    console.log('File updated:', msg.path);
  }
};
```

## Updating Existing Obsidian Plugin

The Obsidian plugin needs to be updated to support authentication:

1. Add login screen to collect username/password
2. Store JWT token in plugin settings
3. Include `Authorization: Bearer {token}` header in all API requests
4. Add `?token={token}` to WebSocket connection URL
5. Handle 401 responses (re-authenticate)

## Troubleshooting

### "Invalid or expired token"
- Token may have expired (24h lifetime)
- Re-authenticate to get a new token

### "User not found"
- User account may have been deleted
- Register a new account

### WebSocket connection fails
- Ensure token is included in URL: `/ws?token={your_token}`
- Check that token is valid

### Files not appearing
- Ensure you're authenticated as the correct user
- Check that files were uploaded with proper authentication

## Migration from Non-Authenticated Version

If you have an existing server without authentication:

1. **Backup existing data**: `cp -r data data.backup`
2. **Deploy new version** with authentication
3. **Register an account** for each previous user
4. **Manually migrate data** by copying files to user directories:
   ```bash
   mkdir -p data/users/{userID}
   cp data.backup/* data/users/{userID}/
   ```

## License

Same as the parent project.
