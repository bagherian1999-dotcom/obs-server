package main

const dashboardHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>GoSync - Secure Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        :root {
            --primary: #667eea;
            --primary-dark: #764ba2;
            --success: #10b981;
            --danger: #ef4444;
            --dark: #1f2937;
            --light: #f9fafb;
            --border: #e5e7eb;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
            min-height: 100vh;
            color: var(--dark);
        }
        .login-page {
            display: flex;
            min-height: 100vh;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .login-box {
            background: white;
            border-radius: 16px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 420px;
            width: 100%;
            overflow: hidden;
        }
        .login-header {
            background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
            color: white;
            padding: 50px 40px;
            text-align: center;
        }
        .login-header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 700;
        }
        .login-header p {
            opacity: 0.9;
            font-size: 0.95em;
        }
        .login-content {
            padding: 40px;
        }
        .tabs {
            display: flex;
            margin-bottom: 30px;
            border-bottom: 2px solid var(--border);
        }
        .tab {
            flex: 1;
            padding: 15px;
            text-align: center;
            cursor: pointer;
            border: none;
            background: none;
            font-size: 15px;
            font-weight: 600;
            color: #666;
            transition: all 0.3s;
            position: relative;
        }
        .tab.active {
            color: var(--primary);
        }
        .tab.active::after {
            content: '';
            position: absolute;
            bottom: -2px;
            left: 0;
            right: 0;
            height: 3px;
            background: var(--primary);
        }
        .tab-content { display: none; }
        .tab-content.active { display: block; }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            color: var(--dark);
            font-weight: 600;
            font-size: 14px;
        }
        .form-group input {
            width: 100%;
            padding: 12px 14px;
            border: 2px solid var(--border);
            border-radius: 8px;
            font-size: 14px;
            background: var(--light);
            transition: all 0.3s;
        }
        .form-group input:focus {
            outline: none;
            border-color: var(--primary);
            background: white;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }
        .btn {
            width: 100%;
            padding: 12px;
            border: none;
            border-radius: 8px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }
        .btn-primary {
            background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
            color: white;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102, 126, 234, 0.3);
        }
        .alert {
            padding: 14px 16px;
            border-radius: 8px;
            margin-bottom: 20px;
            display: none;
            font-size: 14px;
            font-weight: 500;
        }
        .alert-error {
            background: #fee2e2;
            color: #991b1b;
            border-left: 4px solid var(--danger);
        }
        .alert-success {
            background: #dcfce7;
            color: #166534;
            border-left: 4px solid var(--success);
        }
        /* DASHBOARD */
        .dashboard-page {
            display: none;
            min-height: 100vh;
            padding: 20px;
        }
        .dashboard-page.active {
            display: block;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        .dashboard-header {
            background: white;
            border-radius: 12px;
            padding: 30px;
            margin-bottom: 30px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            display: flex;
            justify-content: space-between;
            align-items: center;
            flex-wrap: wrap;
            gap: 20px;
        }
        .user-profile {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        .user-avatar {
            width: 60px;
            height: 60px;
            border-radius: 50%;
            background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
            color: white;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 28px;
            font-weight: 700;
        }
        .user-details h2 {
            font-size: 1.4em;
            margin-bottom: 4px;
            color: var(--dark);
        }
        .user-details p {
            color: #6b7280;
            font-size: 0.95em;
        }
        .dashboard-actions {
            display: flex;
            gap: 10px;
        }
        .btn-sm {
            padding: 10px 20px;
            font-size: 14px;
            width: auto;
        }
        .btn-danger {
            background: var(--danger);
            color: white;
        }
        .btn-danger:hover {
            background: #dc2626;
            transform: translateY(-2px);
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: white;
            padding: 25px;
            border-radius: 12px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            text-align: center;
            border-left: 4px solid var(--primary);
        }
        .stat-card h4 {
            font-size: 2.5em;
            color: var(--primary);
            margin-bottom: 8px;
        }
        .stat-card p {
            color: #6b7280;
            font-size: 0.95em;
        }
        .grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .card {
            background: white;
            border-radius: 12px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .card-header {
            background: linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%);
            color: white;
            padding: 20px;
            font-weight: 600;
            font-size: 1.1em;
        }
        .card-content {
            padding: 20px;
        }
        .file-list {
            max-height: 400px;
            overflow-y: auto;
        }
        .file-item {
            padding: 12px;
            border-bottom: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
            font-size: 0.95em;
        }
        .file-item:last-child {
            border-bottom: none;
        }
        .file-name {
            font-weight: 500;
            color: var(--dark);
            flex: 1;
            word-break: break-word;
        }
        .file-time {
            color: #6b7280;
            font-size: 0.85em;
            margin-left: 10px;
            white-space: nowrap;
        }
        .log-entry {
            padding: 10px 12px;
            margin-bottom: 8px;
            border-radius: 6px;
            font-family: monospace;
            font-size: 0.85em;
            display: flex;
            gap: 10px;
        }
        .log-time {
            color: #6b7280;
            min-width: 130px;
        }
        .log-level {
            font-weight: 600;
            min-width: 60px;
        }
        .log-INFO { background: #dbeafe; color: #1e40af; }
        .log-ERROR { background: #fee2e2; color: #991b1b; }
        .log-CONNECT { background: #dcfce7; color: #166534; }
        .log-WARN { background: #fef3c7; color: #92400e; }
        .search-input {
            width: 100%;
            padding: 12px 14px;
            border: 2px solid var(--border);
            border-radius: 8px;
            font-size: 14px;
            margin-bottom: 15px;
            background: var(--light);
        }
        .search-input:focus {
            outline: none;
            border-color: var(--primary);
            background: white;
        }
        .empty-state {
            text-align: center;
            color: #6b7280;
            padding: 40px 20px;
        }
        .empty-state-icon {
            font-size: 3em;
            margin-bottom: 10px;
            opacity: 0.5;
        }
        @media (max-width: 768px) {
            .dashboard-header {
                flex-direction: column;
                text-align: center;
            }
            .user-profile {
                flex-direction: column;
                width: 100%;
            }
            .dashboard-actions {
                width: 100%;
            }
            .grid {
                grid-template-columns: 1fr;
            }
            .login-header h1 {
                font-size: 2em;
            }
        }
    </style>
</head>
<body>
    <!-- LOGIN PAGE -->
    <div id="loginPage" class="login-page">
        <div class="login-box">
            <div class="login-header">
                <h1>🔐 GoSync</h1>
                <p>Secure Obsidian Sync Server</p>
            </div>
            <div class="login-content">
                <div class="tabs">
                    <button class="tab active" onclick="switchTab('login')">Login</button>
                    <button class="tab" onclick="switchTab('register')">Register</button>
                </div>

                <div id="loginTab" class="tab-content active">
                    <div id="loginAlert" class="alert"></div>
                    <form onsubmit="handleLogin(event)">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" id="loginUsername" required autocomplete="username">
                        </div>
                        <div class="form-group">
                            <label>Password</label>
                            <input type="password" id="loginPassword" required autocomplete="current-password">
                        </div>
                        <button type="submit" class="btn btn-primary">Login</button>
                    </form>
                </div>

                <div id="registerTab" class="tab-content">
                    <div id="registerAlert" class="alert"></div>
                    <form onsubmit="handleRegister(event)">
                        <div class="form-group">
                            <label>Username</label>
                            <input type="text" id="registerUsername" required minlength="3" autocomplete="username">
                        </div>
                        <div class="form-group">
                            <label>Email (optional)</label>
                            <input type="email" id="registerEmail" autocomplete="email">
                        </div>
                        <div class="form-group">
                            <label>Password (min. 6 characters)</label>
                            <input type="password" id="registerPassword" required minlength="6" autocomplete="new-password">
                        </div>
                        <button type="submit" class="btn btn-primary">Create Account</button>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <!-- DASHBOARD PAGE -->
    <div id="dashboardPage" class="dashboard-page">
        <div class="container">
            <div class="dashboard-header">
                <div class="user-profile">
                    <div class="user-avatar" id="userAvatar">U</div>
                    <div class="user-details">
                        <h2 id="userName">Welcome</h2>
                        <p id="userEmail">user@example.com</p>
                    </div>
                </div>
                <div class="dashboard-actions">
                    <button class="btn btn-sm btn-danger" onclick="handleLogout()">Logout</button>
                </div>
            </div>

            <div class="stats">
                <div class="stat-card">
                    <h4 id="fileCount">0</h4>
                    <p>Files Synced</p>
                </div>
                <div class="stat-card">
                    <h4 id="clientCount">0</h4>
                    <p>Connected Devices</p>
                </div>
            </div>

            <div class="grid">
                <div class="card">
                    <div class="card-header">📁 Your Files</div>
                    <div class="card-content">
                        <input type="text" id="fileSearch" placeholder="🔍 Search files..." class="search-input">
                        <div class="file-list" id="fileList">
                            <div class="empty-state">
                                <div class="empty-state-icon">📭</div>
                                <p>No files yet</p>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="card">
                    <div class="card-header">📊 Recent Activity</div>
                    <div class="card-content">
                        <div class="file-list" id="logList">
                            <div class="empty-state">
                                <div class="empty-state-icon">🌫️</div>
                                <p>No activity</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div class="card">
                <div class="card-header" style="background: linear-gradient(135deg, var(--danger) 0%, #b91c1c 100%);">⚠️ Danger Zone</div>
                <div class="card-content">
                    <p style="margin-bottom: 15px; color: #6b7280;">This action will delete all your synced files permanently. This cannot be undone.</p>
                    <button class="btn btn-sm btn-danger" onclick="handleReset()">Reset All Data</button>
                </div>
            </div>
        </div>
    </div>

    <script>
        let authToken = localStorage.getItem('authToken');
        let currentUser = null;
        let statusInterval = null;

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
                    document.getElementById('loginUsername').value = '';
                    document.getElementById('loginPassword').value = '';
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
                    document.getElementById('registerUsername').value = '';
                    document.getElementById('registerEmail').value = '';
                    document.getElementById('registerPassword').value = '';
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
            if (statusInterval) clearInterval(statusInterval);
            document.getElementById('loginPage').style.display = 'flex';
            document.getElementById('dashboardPage').classList.remove('active');
            switchTab('login');
        }

        function showDashboard() {
            document.getElementById('loginPage').style.display = 'none';
            document.getElementById('dashboardPage').classList.add('active');
            
            const avatar = currentUser.username.charAt(0).toUpperCase();
            document.getElementById('userAvatar').textContent = avatar;
            document.getElementById('userName').textContent = 'Welcome, ' + currentUser.username;
            document.getElementById('userEmail').textContent = currentUser.email || 'No email set';
            
            fetchStatus();
            if (statusInterval) clearInterval(statusInterval);
            statusInterval = setInterval(fetchStatus, 3000);
            loadFiles();
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

                    const logList = document.getElementById('logList');
                    if (data.logs && data.logs.length > 0) {
                        logList.innerHTML = data.logs.slice(-15).reverse().map(log => {
                            const date = new Date(log.timestamp);
                            const time = date.toLocaleTimeString();
                            return '<div class="log-entry log-' + log.level + '">' +
                                '<span class="log-time">' + time + '</span>' +
                                '<span class="log-level">' + log.level + '</span>' +
                                '<span>' + log.message + '</span>' +
                                '</div>';
                        }).join('');
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
                    
                    if (!files || files.length === 0) {
                        fileList.innerHTML = '<div class="empty-state"><div class="empty-state-icon">📭</div><p>No files found</p></div>';
                    } else {
                        fileList.innerHTML = files.map(file => {
                            const date = new Date(file.latest.timestamp * 1000);
                            const timeStr = date.toLocaleString();
                            return '<div class="file-item">' +
                                '<span class="file-name">' + escapeHtml(file.path) + '</span>' +
                                '<span class="file-time">' + timeStr + '</span>' +
                                '</div>';
                        }).join('');
                    }
                }
            } catch (err) {
                console.error('Failed to load files:', err);
            }
        }

        async function handleReset() {
            if (!confirm('⚠️ Are you sure? This will delete ALL your files permanently and cannot be undone!')) {
                return;
            }
            if (!confirm('This is your last chance! Click OK again to confirm.')) {
                return;
            }

            try {
                const response = await fetch('/api/reset', {
                    method: 'POST',
                    headers: { 'Authorization': 'Bearer ' + authToken }
                });

                if (response.ok) {
                    showAlert('You can dismiss this', 'All data has been reset', 'success');
                    fetchStatus();
                    loadFiles();
                }
            } catch (err) {
                alert('Failed to reset: ' + err.message);
            }
        }

        function showAlert(id, message, type) {
            let alertElement = null;
            if (id === 'loginAlert') {
                alertElement = document.getElementById('loginAlert');
            } else if (id === 'registerAlert') {
                alertElement = document.getElementById('registerAlert');
            }
            
            if (alertElement) {
                alertElement.className = 'alert alert-' + type;
                alertElement.textContent = message;
                alertElement.style.display = 'block';
                setTimeout(() => alertElement.style.display = 'none', 5000);
            }
        }

        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
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
