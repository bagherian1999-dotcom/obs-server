package main

const dashboardHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GoSync Server Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 1200px;
            width: 100%;
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 { font-size: 2em; margin-bottom: 10px; }
        .header p { opacity: 0.9; }
        .content { padding: 30px; }
        .auth-section {
            max-width: 400px;
            margin: 0 auto;
        }
        .tabs {
            display: flex;
            margin-bottom: 30px;
            border-bottom: 2px solid #e0e0e0;
        }
        .tab {
            flex: 1;
            padding: 15px;
            text-align: center;
            cursor: pointer;
            border: none;
            background: none;
            font-size: 16px;
            color: #666;
            transition: all 0.3s;
        }
        .tab.active {
            color: #667eea;
            border-bottom: 3px solid #667eea;
            font-weight: bold;
        }
        .tab-content { display: none; }
        .tab-content.active { display: block; }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: #333;
            font-weight: 500;
        }
        .form-group input {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 8px;
            font-size: 14px;
            transition: border-color 0.3s;
        }
        .form-group input:focus {
            outline: none;
            border-color: #667eea;
        }
        .btn {
            width: 100%;
            padding: 14px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            transition: all 0.3s;
        }
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        .btn-danger {
            background: #e74c3c;
            color: white;
        }
        .alert {
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
        }
        .alert-error {
            background: #fee;
            color: #c33;
            border: 1px solid #fcc;
        }
        .alert-success {
            background: #efe;
            color: #3c3;
            border: 1px solid #cfc;
        }
        .dashboard-section {
            display: none;
        }
        .dashboard-section.active {
            display: block;
        }
        .user-info {
            background: #f5f5f5;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 30px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .user-info h3 {
            color: #667eea;
            margin-bottom: 5px;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 25px;
            border-radius: 12px;
            text-align: center;
        }
        .stat-card h4 {
            font-size: 2em;
            margin-bottom: 5px;
        }
        .stat-card p {
            opacity: 0.9;
        }
        .section {
            background: #f9f9f9;
            padding: 20px;
            border-radius: 10px;
            margin-bottom: 20px;
        }
        .section h3 {
            margin-bottom: 15px;
            color: #333;
        }
        .file-list {
            max-height: 400px;
            overflow-y: auto;
            background: white;
            border-radius: 8px;
            padding: 15px;
        }
        .file-item {
            padding: 10px;
            border-bottom: 1px solid #eee;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .file-item:last-child {
            border-bottom: none;
        }
        .log-entry {
            padding: 8px;
            margin-bottom: 5px;
            border-radius: 5px;
            font-family: monospace;
            font-size: 12px;
        }
        .log-INFO { background: #e3f2fd; color: #1565c0; }
        .log-ERROR { background: #ffebee; color: #c62828; }
        .log-CONNECT { background: #e8f5e9; color: #2e7d32; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔒 GoSync Server</h1>
            <p>Secure Obsidian Sync with User Authentication</p>
        </div>
        <div class="content">
            <!-- Auth Section -->
            <div id="authSection" class="auth-section">
                <div class="tabs">
                    <button class="tab active" onclick="switchTab('login')">Login</button>
                    <button class="tab" onclick="switchTab('register')">Register</button>
                </div>

                <div id="loginTab" class="tab-content active">
                    <div id="loginAlert" class="alert"></div>
                    <form onsubmit="handleLogin(event)">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" id="loginUsername" required>
                        </div>
                        <div class="form-group">
                            <label>Password</label>
                            <input type="password" id="loginPassword" required>
                        </div>
                        <button type="submit" class="btn btn-primary">Login</button>
                    </form>
                </div>

                <div id="registerTab" class="tab-content">
                    <div id="registerAlert" class="alert"></div>
                    <form onsubmit="handleRegister(event)">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" id="registerUsername" required minlength="3">
                        </div>
                        <div class="form-group">
                            <label>Email (optional)</label>
                            <input type="email" id="registerEmail">
                        </div>
                        <div class="form-group">
                            <label>Password</label>
                            <input type="password" id="registerPassword" required minlength="6">
                        </div>
                        <button type="submit" class="btn btn-primary">Register</button>
                    </form>
                </div>
            </div>

            <!-- Dashboard Section -->
            <div id="dashboardSection" class="dashboard-section">
                <div class="user-info">
                    <div>
                        <h3 id="userName">Welcome</h3>
                        <p id="userEmail"></p>
                    </div>
                    <button class="btn btn-danger" onclick="handleLogout()" style="width: auto; padding: 10px 20px;">Logout</button>
                </div>

                <div class="stats">
                    <div class="stat-card">
                        <h4 id="fileCount">0</h4>
                        <p>Files Synced</p>
                    </div>
                    <div class="stat-card">
                        <h4 id="clientCount">0</h4>
                        <p>Connected Clients</p>
                    </div>
                </div>

                <div class="section">
                    <h3>Your Files</h3>
                    <input type="text" id="fileSearch" placeholder="Search files..." 
                           style="width: 100%; padding: 10px; margin-bottom: 10px; border: 2px solid #e0e0e0; border-radius: 8px;">
                    <div class="file-list" id="fileList">
                        <p style="text-align: center; color: #999;">No files yet</p>
                    </div>
                </div>

                <div class="section">
                    <h3>Recent Activity</h3>
                    <div class="file-list" id="logList">
                        <p style="text-align: center; color: #999;">No activity</p>
                    </div>
                </div>

                <div class="section">
                    <h3>Danger Zone</h3>
                    <button class="btn btn-danger" onclick="handleReset()">Reset All Data</button>
                </div>
            </div>
        </div>
    </div>

    <script>
        let authToken = localStorage.getItem('authToken');
        let currentUser = null;

        // Check if already logged in
        if (authToken) {
            verifyToken();
        }

        function switchTab(tab) {
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));
            
            if (tab === 'login') {
                document.querySelectorAll('.tab')[0].classList.add('active');
                document.getElementById('loginTab').classList.add('active');
            } else {
                document.querySelectorAll('.tab')[1].classList.add('active');
                document.getElementById('registerTab').classList.add('active');
            }
        }

        async function handleLogin(e) {
            e.preventDefault();
            const username = document.getElementById('loginUsername').value;
            const password = document.getElementById('loginPassword').value;

            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, password })
                });

                if (response.ok) {
                    const data = await response.json();
                    authToken = data.token;
                    localStorage.setItem('authToken', authToken);
                    currentUser = data.user;
                    showDashboard();
                } else {
                    const error = await response.text();
                    showAlert('loginAlert', error, 'error');
                }
            } catch (err) {
                showAlert('loginAlert', 'Login failed: ' + err.message, 'error');
            }
        }

        async function handleRegister(e) {
            e.preventDefault();
            const username = document.getElementById('registerUsername').value;
            const email = document.getElementById('registerEmail').value;
            const password = document.getElementById('registerPassword').value;

            try {
                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, email, password })
                });

                if (response.ok) {
                    const data = await response.json();
                    authToken = data.token;
                    localStorage.setItem('authToken', authToken);
                    currentUser = data.user;
                    showDashboard();
                } else {
                    const error = await response.text();
                    showAlert('registerAlert', error, 'error');
                }
            } catch (err) {
                showAlert('registerAlert', 'Registration failed: ' + err.message, 'error');
            }
        }

        async function verifyToken() {
            try {
                const response = await fetch('/api/verify', {
                    headers: { 'Authorization': 'Bearer ' + authToken }
                });

                if (response.ok) {
                    const data = await response.json();
                    currentUser = data.user;
                    showDashboard();
                } else {
                    localStorage.removeItem('authToken');
                    authToken = null;
                }
            } catch (err) {
                localStorage.removeItem('authToken');
                authToken = null;
            }
        }

        function handleLogout() {
            localStorage.removeItem('authToken');
            authToken = null;
            currentUser = null;
            document.getElementById('authSection').style.display = 'block';
            document.getElementById('dashboardSection').classList.remove('active');
        }

        function showDashboard() {
            document.getElementById('authSection').style.display = 'none';
            document.getElementById('dashboardSection').classList.add('active');
            document.getElementById('userName').textContent = 'Welcome, ' + currentUser.username;
            document.getElementById('userEmail').textContent = currentUser.email || 'No email';
            fetchStatus();
            setInterval(fetchStatus, 3000);
        }

        async function fetchStatus() {
            if (!authToken) return;

            try {
                const response = await fetch('/api/status', {
                    headers: { 'Authorization': 'Bearer ' + authToken }
                });

                if (response.ok) {
                    const data = await response.json();
                    document.getElementById('fileCount').textContent = data.fileCount;
                    document.getElementById('clientCount').textContent = data.clients.length;

                    // Update logs
                    const logList = document.getElementById('logList');
                    if (data.logs && data.logs.length > 0) {
                        logList.innerHTML = data.logs.slice(-10).reverse().map(log => 
                            '<div class="log-entry log-' + log.level + '">' +
                            '[' + new Date(log.timestamp).toLocaleTimeString() + '] ' +
                            log.level + ': ' + log.message +
                            '</div>'
                        ).join('');
                    }
                } else if (response.status === 401) {
                    handleLogout();
                }
            } catch (err) {
                console.error('Failed to fetch status:', err);
            }
        }

        async function loadFiles(query = '') {
            if (!authToken) return;

            try {
                const url = query ? '/api/search?q=' + encodeURIComponent(query) : '/api/files';
                const response = await fetch(url, {
                    headers: { 'Authorization': 'Bearer ' + authToken }
                });

                if (response.ok) {
                    const files = await response.json();
                    const fileList = document.getElementById('fileList');
                    
                    if (files.length === 0) {
                        fileList.innerHTML = '<p style="text-align: center; color: #999;">No files found</p>';
                    } else {
                        fileList.innerHTML = files.map(file => 
                            '<div class="file-item">' +
                            '<span>' + file.path + '</span>' +
                            '<small>' + new Date(file.latest.timestamp * 1000).toLocaleString() + '</small>' +
                            '</div>'
                        ).join('');
                    }
                }
            } catch (err) {
                console.error('Failed to load files:', err);
            }
        }

        async function handleReset() {
            if (!confirm('Are you sure you want to delete ALL your files? This cannot be undone!')) {
                return;
            }

            try {
                const response = await fetch('/api/reset', {
                    method: 'POST',
                    headers: { 'Authorization': 'Bearer ' + authToken }
                });

                if (response.ok) {
                    alert('All data has been reset');
                    fetchStatus();
                    loadFiles();
                }
            } catch (err) {
                alert('Failed to reset: ' + err.message);
            }
        }

        function showAlert(id, message, type) {
            const alert = document.getElementById(id);
            alert.className = 'alert alert-' + type;
            alert.textContent = message;
            alert.style.display = 'block';
            setTimeout(() => alert.style.display = 'none', 5000);
        }

        // File search
        document.addEventListener('DOMContentLoaded', () => {
            const searchInput = document.getElementById('fileSearch');
            if (searchInput) {
                searchInput.addEventListener('input', (e) => {
                    loadFiles(e.target.value);
                });
            }
        });
    </script>
</body>
</html>
`
