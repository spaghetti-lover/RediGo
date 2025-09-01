package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/spaghetti-lover/multithread-redis/internal/config"
	"github.com/spaghetti-lover/multithread-redis/internal/core"
)

type CommandRequest struct {
	Cmd string `json:"cmd"`
}

type CommandResponse struct {
	Output interface{} `json:"output"`
	Error  string      `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/command", corsMiddleware(handleCommand))
	http.HandleFunc("/stats", corsMiddleware(handleStats))

	fmt.Printf("HTTP Gateway listening on port %s", config.HTTPPort)
	log.Printf("Will connect to Redis at: %s", config.Port)

	server := &http.Server{
		Addr:         config.HTTPPort,
		ReadTimeout:  time.Duration(config.HTTPReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.HTTPWriteTimeout) * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"keys":    10,    // ví dụ: hàm trả số keys hiện tại
		"memory":  1024,  // MB
		"uptime":  10000, // giây
		"clients": 10,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	output, err := sendToRedis(req.Cmd)
	resp := CommandResponse{Output: output}
	if err != nil {
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Send command to Redis server thoruh TCP (port 6379).
func sendToRedis(cmd string) (interface{}, error) {
	redisAddr := config.RedisAddr

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	conn, err := net.DialTimeout("tcp", redisAddr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Set read/write timeout
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// Parse command to array of strings.
	parts := strings.Fields(strings.TrimSpace(cmd))
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	// Convert to []interface{} to Encode
	cmdArray := make([]interface{}, len(parts))
	for i, part := range parts {
		cmdArray[i] = part
	}

	// Encode command to RESP
	respData := core.Encode(cmdArray, false)

	// Log for debug
	log.Printf("Sending RESP: %q", string(respData))

	// Send command
	_, err = conn.Write(respData)
	if err != nil {
		return nil, err
	}

	// Read response
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("connection closed by server: %v", err)
		}
		return nil, err
	}

	rawResponse := buf[:n]
	log.Printf("Received raw response: %q", string(rawResponse))

	// Decode RESP command
	decodedResponse, err := core.Decode(rawResponse)
	if err != nil {
		log.Printf("Decode error: %v, raw data: %q", err, string(rawResponse))
		return string(rawResponse), nil
	}

	log.Printf("Decoded response: %+v", decodedResponse)
	return decodedResponse, nil
}
