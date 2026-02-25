package graph

import (
	"sort"
	"sync"
)

type Graph struct {
	AdjList map[int][]int
	mu      sync.RWMutex
}

type InfluenceMetrics struct {
	UserID              int     `json:"user_id"`
	InDegree            int     `json:"in_degree"`
	OutDegree           int     `json:"out_degree"`
	NormalizedInDegree  float64 `json:"normalized_in_degree"`
}

func NewGraph() *Graph {
	return &Graph{
		AdjList: make(map[int][]int),
	}
}

func (g *Graph) AddNode(id int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.AdjList[id]; !exists {
		g.AdjList[id] = []int{}
	}
}

func (g *Graph) AddEdge(from, to int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.AdjList[from]; !exists {
		g.AdjList[from] = []int{}
	}
	if _, exists := g.AdjList[to]; !exists {
		g.AdjList[to] = []int{}
	}

	g.AdjList[from] = append(g.AdjList[from], to)
}

func (g *Graph) HasNode(id int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, exists := g.AdjList[id]
	return exists
}

func (g *Graph) NodeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.AdjList)
}

func (g *Graph) EdgeCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	count := 0
	for _, neighbors := range g.AdjList {
		count += len(neighbors)
	}
	return count
}

func (g *Graph) OutDegree(userID int) int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.AdjList[userID])
}

func (g *Graph) inDegreeUnsafe(userID int) int {
	count := 0
	for _, neighbors := range g.AdjList {
		for _, neighbor := range neighbors {
			if neighbor == userID {
				count++
			}
		}
	}
	return count
}

func (g *Graph) Influence(userID int) InfluenceMetrics {
	g.mu.RLock()
	defer g.mu.RUnlock()

	in := g.inDegreeUnsafe(userID)
	out := len(g.AdjList[userID])
	totalNodes := len(g.AdjList)

	var normalized float64
	if totalNodes > 1 {
		normalized = float64(in) / float64(totalNodes-1)
	}

	return InfluenceMetrics{
		UserID:             userID,
		InDegree:           in,
		OutDegree:          out,
		NormalizedInDegree: normalized,
	}
}

func (g *Graph) ShortestPath(src, dest int) []int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if src == dest {
		return []int{src}
	}

	visited := make(map[int]bool)
	parent := make(map[int]int)
	queue := []int{src}
	visited[src] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.AdjList[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = current
				queue = append(queue, neighbor)

				if neighbor == dest {
					path := []int{}
					for at := dest; ; {
						path = append([]int{at}, path...)
						if at == src {
							break
						}
						at = parent[at]
					}
					return path
				}
			}
		}
	}
	return nil
}

func (g *Graph) ConnectionsWithinDepth(userID int, depth int) []int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	type nodeDepth struct {
		node  int
		depth int
	}

	visited := make(map[int]bool)
	queue := []nodeDepth{{userID, 0}}
	visited[userID] = true

	result := []int{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.depth == depth {
			continue
		}

		for _, neighbor := range g.AdjList[current.node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				result = append(result, neighbor)
				queue = append(queue, nodeDepth{neighbor, current.depth + 1})
			}
		}
	}

	return result
}

func (g *Graph) TopInfluencers(limit int) []int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	inDegreeMap := make(map[int]int)

	// Initialize
	for id := range g.AdjList {
		inDegreeMap[id] = 0
	}

	// Count inbound edges
	for _, neighbors := range g.AdjList {
		for _, neighbor := range neighbors {
			inDegreeMap[neighbor]++
		}
	}

	type userScore struct {
		id    int
		score int
	}

	var scores []userScore
	for id, score := range inDegreeMap {
		scores = append(scores, userScore{id, score})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	if limit > len(scores) {
		limit = len(scores)
	}

	result := make([]int, limit)
	for i := 0; i < limit; i++ {
		result[i] = scores[i].id
	}

	return result
}

func (g *Graph) dfs(node int, visited map[int]bool) {
	visited[node] = true
	for _, neighbor := range g.AdjList[node] {
		if !visited[neighbor] {
			g.dfs(neighbor, visited)
		}
	}
}

func (g *Graph) ConnectedComponents() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[int]bool)
	components := 0

	for node := range g.AdjList {
		if !visited[node] {
			g.dfs(node, visited)
			components++
		}
	}

	return components
}