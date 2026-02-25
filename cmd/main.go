package main

import (
	"log"
	"net/http"

	"social/internal/db"
	"social/internal/graph"
	"social/internal/handlers"
)

func main() {

	database := db.Connect()
	db.InitSchema(database)

	inMemoryGraph := graph.LoadGraph(database)

	userHandler := &handlers.UserHandler{
		DB:    database,
		Graph: inMemoryGraph,
	}

	connectionHandler := &handlers.ConnectionHandler{
		DB:    database,
		Graph: inMemoryGraph,
	}

	graphHandler := &handlers.GraphHandler{
		Graph: inMemoryGraph,
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/api/v1/bulk-import", userHandler.BulkImport)
	http.HandleFunc("/api/v1/users", userHandler.CreateUser)
	http.HandleFunc("/api/v1/connections", connectionHandler.CreateConnection)

	http.HandleFunc("/api/v1/graph/stats", graphHandler.Stats)
	http.HandleFunc("/api/v1/graph/degree/", graphHandler.Degree)
	http.HandleFunc("/api/v1/graph/reach/", graphHandler.Reach)
	http.HandleFunc("/api/v1/graph/influence/", graphHandler.Influence)

	http.HandleFunc("/api/v1/graph/shortest-path", graphHandler.ShortestPath)
	http.HandleFunc("/api/v1/graph/components", graphHandler.Components)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}