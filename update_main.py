#!/usr/bin/env python3

import re

# Read the original main.go
with open('server/main.go', 'r') as f:
    content = f.read()

# Update constants section
content = re.sub(
    r'const \(\s*Port\s*=\s*"[^"]*"\s*DataDir\s*=\s*"[^"]*"\s*MetadataFile\s*=\s*"[^"]*"\s*MaxLogs\s*=\s*\d+\s*\)',
    '''const (
\tMetadataFile = "./metadata.json"
\tMaxLogs      = 100
\tConfigFile   = "./config.json"
)''',
    content
)

# Update main function
old_main = r'func main\(\) \{\s*// Ensure data directory exists\s*if err := os\.MkdirAll\(DataDir, 0755\); err != nil \{'
new_main = '''func main() {
\t// Load configuration
\tif err := loadConfig(ConfigFile); err != nil {
\t\tlog.Printf("Warning: Could not load config file (%v), using defaults", err)
\t}

\tDataDir := config.DataDir
\t
\t// Ensure data directory exists
\tif err := os.MkdirAll(config.DataDir, 0755); err != nil {'''

content = re.sub(old_main, new_main, content, flags=re.DOTALL)

# Update server start section
old_start = r'\taddLog\("INFO", fmt\.Sprintf\("GoSync Server listening on %s", Port\)\)\s*fmt\.Printf\("GoSync Server started at http://localhost%s\\n", Port\)\s*\s*log\.Fatal\(http\.ListenAndServe\(Port, nil\)\)'
new_start = '''\tserverAddr := config.GetAddress()
\tpublicURL := config.GetPublicURL()
\t
\taddLog("INFO", fmt.Sprintf("GoSync Server listening on %s", serverAddr))
\tfmt.Printf("GoSync Server started at %s\\n", publicURL)
\t
\tif config.EnableSSL {
\t\taddLog("INFO", "SSL enabled")
\t\tlog.Fatal(http.ListenAndServeTLS(serverAddr, config.SSLCertPath, config.SSLKeyPath, nil))
\t} else {
\t\taddLog("INFO", "Running without SSL (HTTP only)")
\t\tlog.Fatal(http.ListenAndServe(serverAddr, nil))
\t}'''

content = re.sub(old_start, new_start, content)

# Write updated content
with open('server/main.go', 'w') as f:
    f.write(content)

print("main.go updated successfully!")
