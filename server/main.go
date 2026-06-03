package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MaxLogs    = 100
	ConfigFile = "./config.json"
)

// --- Data Structures ---

type FileVersion struct {
	Hash      string `json:"hash"`
	Timestamp int64  `json:"timestamp"`
	Device    string `json:"device"`
}

type FileMetadata struct {
	Path    string        `json:"path"`
	History []FileVersion `json:"history"`
	Latest  FileVersion   `json:"latest"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

type ClientInfo struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	IP          string    `json:"ip"`
	Pn          string    `json:"pn"`
	UserID      string    `json:"user_id"`
	ConnectedAt time.Time `json:"connectedAt"`
}

type ServerState struct {
	mu      sync.RWMutex
	Clients map[*websocket.Conn]*ClientInfo
	Logs    []LogEntry
}

// User-specific file state (per user)
type UserFileState struct {
	mu    sync.RWMutex
	Files map[string]FileMetadata
}

// --- Globals ---

var state = ServerState{
	Clients: make(map[*websocket.Conn]*ClientInfo),
	Logs:    make([]LogEntry, 0),
}

var userFileStates = make(map[string]*UserFileState) // key: userID
var userFileStatesMu sync.RWMutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var broadcast = make(chan interface{})

// --- Main ---

func main() {
	// Load configuration
	if err := loadConfig(ConfigFile); err != nil {
		log.Printf("Warning: Could not load config file (%v), using defaults", err)
	}

	// Ensure data directory exists
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Initialize user store
	if err := initUserStore(); err != nil {
		log.Fatal("Failed to initialize user store:", err)
	}

	addLog("INFO", "Server started with authentication enabled")

	go handleMessages()

	// Public endpoints
	http.HandleFunc("/api/register", enableCors(handleRegister))
	http.HandleFunc("/api/login", enableCors(handleLogin))
	http.HandleFunc("/", handleDashboard)

	// Protected endpoints (require authentication)
	http.HandleFunc("/api/profile", enableCorsWithAuth(requireAuth(handleProfile)))
	http.HandleFunc("/api/verify", enableCorsWithAuth(requireAuth(handleVerifyToken)))
	http.HandleFunc("/api/files", enableCorsWithAuth(requireAuth(handleListFiles)))
	http.HandleFunc("/api/file", enableCorsWithAuth(requireAuth(handleFileOperations)))
	http.HandleFunc("/api/files/delete", enableCorsWithAuth(requireAuth(handleBulkDelete)))
	http.HandleFunc("/api/search", enableCorsWithAuth(requireAuth(handleSearch)))
	http.HandleFunc("/api/cleanup", enableCorsWithAuth(requireAuth(handleCleanup)))
	http.HandleFunc("/api/reset", enableCorsWithAuth(requireAuth(handleReset)))
	http.HandleFunc("/api/status", enableCorsWithAuth(requireAuth(handleServerStatus)))

	// WebSocket (requires authentication via query param)
	http.HandleFunc("/ws", handleConnections)

	serverAddr := config.GetAddress()
	publicURL := config.GetPublicURL()

	addLog("INFO", fmt.Sprintf("GoSync Server listening on %s", serverAddr))
	fmt.Printf("GoSync Server started at %s\n", publicURL)
	fmt.Printf("Authentication enabled - users must register/login\n")

	if config.EnableSSL {
		addLog("INFO", "SSL enabled")
		log.Fatal(http.ListenAndServeTLS(serverAddr, config.SSLCertPath, config.SSLKeyPath, nil))
	} else {
		addLog("INFO", "Running without SSL (HTTP only)")
		log.Fatal(http.ListenAndServe(serverAddr, nil))
	}
}

// --- User File State Management ---

func getUserFileState(userID string) *UserFileState {
	userFileStatesMu.RLock()
	ufs, exists := userFileStates[userID]
	userFileStatesMu.RUnlock()

	if exists {
		return ufs
	}

	// Create new state
	userFileStatesMu.Lock()
	defer userFileStatesMu.Unlock()

	// Double-check after lock
	if ufs, exists := userFileStates[userID]; exists {
		return ufs
	}

	ufs = &UserFileState{
		Files: make(map[string]FileMetadata),
	}
	userFileStates[userID] = ufs

	// Load metadata for this user
	loadUserMetadata(userID, ufs)

	return ufs
}

func loadUserMetadata(userID string, ufs *UserFileState) {
	metadataFile := getUserMetadataFile(userID)

	file, err := os.Open(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // No metadata yet
		}
		addLog("ERROR", fmt.Sprintf("Failed to load metadata for user %s: %v", userID, err))
		return
	}
	defer file.Close()

	var files map[string]FileMetadata
	if err := json.NewDecoder(file).Decode(&files); err != nil {
		addLog("ERROR", fmt.Sprintf("Failed to decode metadata for user %s: %v", userID, err))
		return
	}

	ufs.mu.Lock()
	ufs.Files = files
	ufs.mu.Unlock()

	addLog("INFO", fmt.Sprintf("Loaded %d files for user %s", len(files), userID))
}

func saveUserMetadata(userID string, ufs *UserFileState) error {
	metadataFile := getUserMetadataFile(userID)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(metadataFile), 0755); err != nil {
		return err
	}

	ufs.mu.RLock()
	defer ufs.mu.RUnlock()

	file, err := os.Create(metadataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ufs.Files)
}

// --- Logging & State Helpers ---

func addLog(level, message string) {
	state.mu.Lock()
	defer state.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   message,
		Level:     level,
	}

	state.Logs = append(state.Logs, entry)
	if len(state.Logs) > MaxLogs {
		state.Logs = state.Logs[1:]
	}

	fmt.Printf("[%s] %s: %s\n", entry.Timestamp.Format("15:04:05"), level, message)
}

// --- HTTP Handlers ---

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, dashboardHTML)
}

func handleServerStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	state.mu.RLock()
	defer state.mu.RUnlock()

	// Get user's file state
	ufs := getUserFileState(user.ID)
	ufs.mu.RLock()
	fileCount := len(ufs.Files)
	ufs.mu.RUnlock()

	// Convert Clients map to list for this user
	clientList := make([]*ClientInfo, 0)
	for _, c := range state.Clients {
		if c.UserID == user.ID {
			clientList = append(clientList, c)
		}
	}

	response := map[string]interface{}{
		"clients":   clientList,
		"logs":      state.Logs,
		"fileCount": fileCount,
		"user":      user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter for WebSocket
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	// Verify token
	claims, err := verifyToken(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	user, err := getUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		addLog("ERROR", fmt.Sprintf("WebSocket Upgrade error: %v", err))
		return
	}
	defer ws.Close()

	clientIP := r.RemoteAddr
	deviceName := user.Username

	info := &ClientInfo{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:        deviceName,
		IP:          clientIP,
		UserID:      user.ID,
		ConnectedAt: time.Now(),
	}

	state.mu.Lock()
	state.Clients[ws] = info
	state.mu.Unlock()

	addLog("CONNECT", fmt.Sprintf("Client connected: %s (User: %s)", clientIP, user.Username))

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				addLog("ERROR", fmt.Sprintf("Client %s error: %v", deviceName, err))
			} else {
				addLog("CONNECT", fmt.Sprintf("Client disconnected: %s", deviceName))
			}
			break
		}

		if msgType, ok := msg["type"].(string); ok {
			if msgType == "identify" {
				if name, ok := msg["deviceName"].(string); ok {
					deviceName = name
					state.mu.Lock()
					info.Name = deviceName
					if pn, ok := msg["pluginName"].(string); ok {
						info.Pn = pn
					}
					state.mu.Unlock()
					addLog("INFO", fmt.Sprintf("Client identified: %s (User: %s)", deviceName, user.Username))
				}
			}
		}

		broadcast <- msg
	}

	state.mu.Lock()
	delete(state.Clients, ws)
	state.mu.Unlock()
}

func handleMessages() {
	for {
		msg := <-broadcast
		state.mu.RLock()
		for client := range state.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(state.Clients, client)
			}
		}
		state.mu.RUnlock()
	}
}

func handleListFiles(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ufs := getUserFileState(user.ID)
	ufs.mu.RLock()
	defer ufs.mu.RUnlock()

	files := make([]FileMetadata, 0, len(ufs.Files))
	for _, meta := range ufs.Files {
		files = append(files, meta)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func handleFileOperations(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Missing path parameter", http.StatusBadRequest)
		return
	}

	ufs := getUserFileState(user.ID)

	switch r.Method {
	case "GET":
		handleGetFile(w, r, user.ID, path, ufs)
	case "PUT":
		handlePutFile(w, r, user.ID, path, ufs)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetFile(w http.ResponseWriter, r *http.Request, userID, path string, ufs *UserFileState) {
	ufs.mu.RLock()
	meta, exists := ufs.Files[path]
	ufs.mu.RUnlock()

	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	blobPath := filepath.Join(getUserDataDir(userID), meta.Latest.Hash)
	data, err := os.ReadFile(blobPath)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

func handlePutFile(w http.ResponseWriter, r *http.Request, userID, path string, ufs *UserFileState) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Calculate hash
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// Save blob
	blobPath := filepath.Join(getUserDataDir(userID), hashStr)
	if err := os.WriteFile(blobPath, data, 0644); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	device := r.URL.Query().Get("device")
	if device == "" {
		device = "unknown"
	}

	version := FileVersion{
		Hash:      hashStr,
		Timestamp: time.Now().Unix(),
		Device:    device,
	}

	ufs.mu.Lock()
	meta, exists := ufs.Files[path]
	if !exists {
		meta = FileMetadata{
			Path:    path,
			History: []FileVersion{},
		}
	}
	meta.History = append(meta.History, version)
	meta.Latest = version
	ufs.Files[path] = meta
	ufs.mu.Unlock()

	saveUserMetadata(userID, ufs)

	addLog("INFO", fmt.Sprintf("File uploaded: %s (User: %s)", path, userID))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"hash":   hashStr,
	})

	broadcast <- map[string]interface{}{
		"type": "file_updated",
		"path": path,
		"hash": hashStr,
	}
}

func handleBulkDelete(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	patterns := r.URL.Query().Get("patterns")
	if patterns == "" {
		http.Error(w, "Missing patterns parameter", http.StatusBadRequest)
		return
	}

	patternList := strings.Split(patterns, ",")
	ufs := getUserFileState(user.ID)

	ufs.mu.Lock()
	deleted := 0
	for path := range ufs.Files {
		for _, pattern := range patternList {
			pattern = strings.TrimSpace(pattern)
			if strings.Contains(path, pattern) {
				delete(ufs.Files, path)
				deleted++
				break
			}
		}
	}
	ufs.mu.Unlock()

	saveUserMetadata(user.ID, ufs)
	addLog("INFO", fmt.Sprintf("Bulk delete: %d files (User: %s)", deleted, user.Username))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deleted": deleted,
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")
	ufs := getUserFileState(user.ID)

	ufs.mu.RLock()
	defer ufs.mu.RUnlock()

	results := make([]FileMetadata, 0)
	for _, meta := range ufs.Files {
		if strings.Contains(meta.Path, query) {
			results = append(results, meta)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleCleanup(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ufs := getUserFileState(user.ID)
	ufs.mu.RLock()
	validHashes := make(map[string]bool)
	for _, meta := range ufs.Files {
		for _, version := range meta.History {
			validHashes[version.Hash] = true
		}
	}
	ufs.mu.RUnlock()

	userDir := getUserDataDir(user.ID)
	entries, err := os.ReadDir(userDir)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	deletedBlobs := 0
	freedBytes := int64(0)

	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "metadata.json" {
			continue
		}

		if !validHashes[entry.Name()] {
			path := filepath.Join(userDir, entry.Name())
			if info, err := os.Stat(path); err == nil {
				freedBytes += info.Size()
			}
			if err := os.Remove(path); err == nil {
				deletedBlobs++
			}
		}
	}

	addLog("INFO", fmt.Sprintf("Cleanup: removed %d blobs, freed %d bytes (User: %s)", deletedBlobs, freedBytes, user.Username))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deletedBlobs": deletedBlobs,
		"freedBytes":   freedBytes,
	})
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userDir := getUserDataDir(user.ID)
	if err := os.RemoveAll(userDir); err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset: %v", err), http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(userDir, 0755); err != nil {
		http.Error(w, "Failed to recreate directory", http.StatusInternalServerError)
		return
	}

	ufs := getUserFileState(user.ID)
	ufs.mu.Lock()
	ufs.Files = make(map[string]FileMetadata)
	ufs.mu.Unlock()

	saveUserMetadata(user.ID, ufs)

	addLog("INFO", fmt.Sprintf("Server reset by user: %s", user.Username))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "reset_complete",
	})
}

