package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type NotificationPayload struct {
	Message string `json:"message"`
}

type NotificationLog struct {
	ReceivedAt time.Time `json:"received_at"`
	Message    string    `json:"message"`
}

var logs []NotificationLog

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var payload NotificationPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			payload.Message = string(body)
		}

		entry := NotificationLog{ReceivedAt: time.Now(), Message: payload.Message}
		logs = append(logs, entry)

		log.Printf("📣 [NOTIFICATION] %s — %s", entry.ReceivedAt.Format(time.RFC3339), entry.Message)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "received"})
	})

	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "notifier"})
	})

	fmt.Println("🔔 Notifier service running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
