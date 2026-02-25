package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social/internal/graph"
)

type ConnectionHandler struct {
	DB    *sql.DB
	Graph *graph.Graph
}

type CreateConnectionRequest struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
}

func (h *ConnectionHandler) CreateConnection(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FromUserID == 0 || req.ToUserID == 0 {
		http.Error(w, "from_user_id and to_user_id required", http.StatusBadRequest)
		return
	}

	var exists int

	err := h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id=$1", req.FromUserID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "From user not found", http.StatusNotFound)
		return
	}

	err = h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id=$1", req.ToUserID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists == 0 {
		http.Error(w, "To user not found", http.StatusNotFound)
		return
	}

	query := `
	INSERT INTO connections (from_user_id, to_user_id)
	VALUES ($1, $2)
	RETURNING id, created_at
	`

	var id int64
	var createdAt string

	err = h.DB.QueryRow(query, req.FromUserID, req.ToUserID).
		Scan(&id, &createdAt)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	h.Graph.AddEdge(int(req.FromUserID), int(req.ToUserID))

	resp := map[string]interface{}{
		"id":           id,
		"from_user_id": req.FromUserID,
		"to_user_id":   req.ToUserID,
		"created_at":   createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}