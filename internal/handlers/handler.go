package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"social/internal/graph"
)

type GraphHandler struct {
	Graph *graph.Graph
}

func (h *GraphHandler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"total_nodes": h.Graph.NodeCount(),
		"total_edges": h.Graph.EdgeCount(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *GraphHandler) Influence(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	idStr = strings.TrimSuffix(idStr, "/influence")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if !h.Graph.HasNode(userID) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	metrics := h.Graph.Influence(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

func (h *GraphHandler) ShortestPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	from, err1 := strconv.Atoi(fromStr)
	to, err2 := strconv.Atoi(toStr)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	if !h.Graph.HasNode(from) || !h.Graph.HasNode(to) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	path := h.Graph.ShortestPath(from, to)
	if path == nil {
		http.Error(w, "No path found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"from": from,
		"to":   to,
		"path": path,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *GraphHandler) TopInfluencers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit", http.StatusBadRequest)
		return
	}

	top := h.Graph.TopInfluencers(limit)

	response := map[string]interface{}{
		"top_influencers": top,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *GraphHandler) Connections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/users/")
	idStr = strings.TrimSuffix(idStr, "/connections")

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if !h.Graph.HasNode(userID) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	depthStr := r.URL.Query().Get("depth")
	depth, err := strconv.Atoi(depthStr)
	if err != nil || depth < 1 {
		http.Error(w, "Invalid depth", http.StatusBadRequest)
		return
	}

	connections := h.Graph.ConnectionsWithinDepth(userID, depth)

	response := map[string]interface{}{
		"user_id":     userID,
		"connections": connections,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}