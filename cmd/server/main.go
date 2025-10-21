package main

import (
	"log"
	"net/http"

	"github.com/aiplaybookin/tiffin-go/internal/server"
)

func main() {
	srv := server.NewServer()
	srv.Start()

	// API routes
	http.HandleFunc("/api/create", srv.HandleCreateGame)
	http.HandleFunc("/api/join", srv.HandleJoinGame)
	http.HandleFunc("/ws", srv.HandleWebSocket)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
