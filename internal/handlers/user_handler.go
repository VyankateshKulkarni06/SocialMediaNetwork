package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social/internal/graph"
)

type UserHandler struct {
	DB    *sql.DB
	Graph *graph.Graph
}


type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and Email required", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO users (name, email, bio)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
	`

	var id int64
	var createdAt string

	err := h.DB.QueryRow(query, req.Name, req.Email, req.Bio).
		Scan(&id, &createdAt)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
	return
	}

	h.Graph.AddNode(int(id))

	resp := map[string]interface{}{
		"id":         id,
		"name":       req.Name,
		"email":      req.Email,
		"bio":        req.Bio,
		"created_at": createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}



type BulkUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Bio   string `json:"bio"`
}

type BulkConnection struct {
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
}

type BulkImportRequest struct {
	Users       []BulkUser       `json:"users"`
	Connections []BulkConnection `json:"connections"`
}

func (h *UserHandler) BulkImport(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req BulkImportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	for _, user := range req.Users {

		if user.Name == "" || user.Email == "" {
			tx.Rollback()
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}

		_, err := tx.Exec(`
			INSERT INTO users (id, name, email, bio)
			VALUES ($1, $2, $3, $4)
		`, user.ID, user.Name, user.Email, user.Bio)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, conn := range req.Connections {

		_, err := tx.Exec(`
			INSERT INTO connections (from_user_id, to_user_id)
			VALUES ($1, $2)
		`, conn.FromUserID, conn.ToUserID)

		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Transaction commit failed", http.StatusInternalServerError)
		return
	}

	graph.ReloadGraph(h.Graph, h.DB)

	response := map[string]interface{}{
		"message":           "Bulk import successful",
		"users_inserted":    len(req.Users),
		"connections_added": len(req.Connections),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}