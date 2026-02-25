package graph

import (
	"database/sql"
)

func LoadGraph(db *sql.DB) *Graph {

	g := NewGraph()

	rows, _ := db.Query("SELECT id FROM users")
	defer rows.Close()

	for rows.Next() {
		var id int
		rows.Scan(&id)
		g.AdjList[id] = []int{}
	}

	connRows, _ := db.Query("SELECT from_user_id, to_user_id FROM connections")
	defer connRows.Close()

	for connRows.Next() {
		var from, to int
		connRows.Scan(&from, &to)
		g.AdjList[from] = append(g.AdjList[from], to)
	}

	return g
}

func ReloadGraph(g *Graph, db *sql.DB) {

	g.mu.Lock()
	defer g.mu.Unlock()

	g.AdjList = make(map[int][]int)

	rows, _ := db.Query("SELECT id FROM users")
	defer rows.Close()

	for rows.Next() {
		var id int
		rows.Scan(&id)
		g.AdjList[id] = []int{}
	}

	connRows, _ := db.Query("SELECT from_user_id, to_user_id FROM connections")
	defer connRows.Close()

	for connRows.Next() {
		var from, to int
		connRows.Scan(&from, &to)
		g.AdjList[from] = append(g.AdjList[from], to)
	}
}