#!/usr/bin/env python3

# Read the original main.go
with open('server/main.go', 'r') as f:
    lines = f.readlines()

new_lines = []
i = 0
while i < len(lines):
    line = lines[i]
    
    # Update constants section
    if line.strip().startswith('const ('):
        new_lines.append(line)
        i += 1
        # Skip old constants and add new ones
        while i < len(lines) and not lines[i].strip().startswith(')'):
            i += 1
        new_lines.append('\tMetadataFile = "./metadata.json"\n')
        new_lines.append('\tMaxLogs      = 100\n')
        new_lines.append('\tConfigFile   = "./config.json"\n')
        new_lines.append(')\n')
        i += 1
        continue
    
    # Update main function start
    if line.strip() == 'func main() {':
        new_lines.append(line)
        new_lines.append('\t// Load configuration\n')
        new_lines.append('\tif err := loadConfig(ConfigFile); err != nil {\n')
        new_lines.append('\t\tlog.Printf("Warning: Could not load config file (%v), using defaults", err)\n')
        new_lines.append('\t}\n')
        new_lines.append('\n')
        new_lines.append('\tDataDir := config.DataDir\n')
        new_lines.append('\t\n')
        i += 1
        continue
    
    # Update MkdirAll call
    if 'os.MkdirAll(DataDir,' in line:
        new_lines.append(line.replace('DataDir', 'config.DataDir'))
        i += 1
        continue
    
    # Update server start section
    if 'addLog("INFO", fmt.Sprintf("GoSync Server listening on %s", Port))' in line:
        new_lines.append('\tserverAddr := config.GetAddress()\n')
        new_lines.append('\tpublicURL := config.GetPublicURL()\n')
        new_lines.append('\t\n')
        new_lines.append('\taddLog("INFO", fmt.Sprintf("GoSync Server listening on %s", serverAddr))\n')
        new_lines.append('\tfmt.Printf("GoSync Server started at %s\\n", publicURL)\n')
        new_lines.append('\t\n')
        new_lines.append('\tif config.EnableSSL {\n')
        new_lines.append('\t\taddLog("INFO", "SSL enabled")\n')
        new_lines.append('\t\tlog.Fatal(http.ListenAndServeTLS(serverAddr, config.SSLCertPath, config.SSLKeyPath, nil))\n')
        new_lines.append('\t} else {\n')
        new_lines.append('\t\taddLog("INFO", "Running without SSL (HTTP only)")\n')
        new_lines.append('\t\tlog.Fatal(http.ListenAndServe(serverAddr, nil))\n')
        new_lines.append('\t}\n')
        # Skip the next 3 lines (old Printf and ListenAndServe)
        i += 1
        while i < len(lines) and 'log.Fatal(http.ListenAndServe(Port, nil))' not in lines[i]:
            i += 1
        i += 1
        continue
    
    new_lines.append(line)
    i += 1

# Write updated content
with open('server/main.go', 'w') as f:
    f.writelines(new_lines)

print("main.go updated successfully!")
