package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

type CommandRequest struct {
	Cmd string `json:"cmd"`
}

type CommandResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

func main() {
	http.HandleFunc("/command", handleCommand)
	fmt.Println("HTTP Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

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

// Gửi command tới Redis server qua TCP (port 6379)
func sendToRedis(cmd string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:6379") // đổi thành host:port Redis của bạn
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Gửi lệnh theo RESP (ở mức đơn giản: tách theo space, wrap lại)
	parts := strings.Split(cmd, " ")
	resp := fmt.Sprintf("*%d\r\n", len(parts))
	for _, p := range parts {
		resp += fmt.Sprintf("$%d\r\n%s\r\n", len(p), p)
	}

	_, err = conn.Write([]byte(resp))
	if err != nil {
		return "", err
	}

	// Đọc trả lời
	buf, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}
