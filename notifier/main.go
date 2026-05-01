package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/notify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, _ := io.ReadAll(r.Body)
			fmt.Printf("notification: received message from Task Manager: %s\n", string(body))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("notification successfully received"))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("method not allowed"))
		}
	})

	fmt.Println("notification service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
