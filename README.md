# Social Network Influence Analyzer

A Go-based RESTful web service that ingests a social network dataset,
stores it in PostgreSQL, models relationships using an in-memory
adjacency list graph, and exposes APIs to compute graph metrics such as
shortest path, connected components, and centrality.

------------------------------------------------------------------------

##  Features

-   Adjacency List Graph Representation
-   BFS (Shortest Path)
-   DFS (Connected Components Detection)
-   Degree Centrality (In-degree & Out-degree)
-   Influence Score Calculation
-   Thread-safe Graph using sync.RWMutex
-   PostgreSQL Persistence
-   RESTful API
-   Proper HTTP Status Codes (200, 201, 400, 404, 500)

------------------------------------------------------------------------

##  Tech Stack

-   Go (net/http)
-   PostgreSQL
-   sync.RWMutex (Concurrency Control)
-   JSON-based REST APIs

------------------------------------------------------------------------

##  Architecture Overview

Client → Handlers → In-Memory Graph → PostgreSQL

-   Database is the source of truth.
-   Graph is loaded into memory for fast computation.
-   Handlers bridge HTTP requests and business logic.

------------------------------------------------------------------------

##  Project Structure

social/ │ ├── internal/ │ ├── db/\
│ ├── graph/\
│ ├── handlers/\
│ ├── main.go\
├── go.mod └── README.md

------------------------------------------------------------------------

##  Setup Instructions

###  Install PostgreSQL

Create database:

CREATE DATABASE social;

###  Configure Database Connection

Update connection settings inside internal/db/db.go

Example:

connStr := "user=postgres password=yourpassword dbname=social
sslmode=disable"

###  Install Dependencies

go mod tidy

###  Run the Server

go run main.go

Server runs at:

http://localhost:8080

------------------------------------------------------------------------

##  API Endpoints

GET /health

POST /api/v1/users

POST /api/v1/connections

POST /api/v1/bulk-import

GET /api/v1/graph/stats

GET /api/v1/graph/degree/{id}

GET /api/v1/graph/reach/{id}

GET /api/v1/graph/shortest-path?src=1&dest=4

GET /api/v1/graph/components

GET /api/v1/graph/influence/{id}

------------------------------------------------------------------------

##  Example Curl

Create User:

curl -X POST http://localhost:8080/api/v1/users -H "Content-Type:
application/json" -d
'{"name":"Alice","email":"alice@example.com","bio":"Engineer"}'

Shortest Path:

curl "http://localhost:8080/api/v1/graph/shortest-path?src=1&dest=4"

------------------------------------------------------------------------

##  Algorithms Implemented

-   BFS (Shortest Path)
-   DFS (Connected Components)
-   Degree Centrality (In & Out)
-   Influence Score

Time Complexity:

-   BFS: O(N + E)
-   DFS: O(N + E)

------------------------------------------------------------------------

##  Concurrency Handling

Graph protected using sync.RWMutex.

-   Multiple readers allowed
-   Single writer allowed

------------------------------------------------------------------------

##  Author

Vyankatesh Kulkarni
