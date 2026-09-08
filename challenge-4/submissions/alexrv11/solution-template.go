package main

import (
    "fmt"
    "sync"
    )

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	// TODO: Implement concurrency-based BFS for multiple queries.
	// Return an empty map so the code compiles but fails tests if unchanged.
	results := map[int][]int{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	if len(queries) == 0 {
	    return results
	}
	
	if numWorkers < 1 {
	    return results
	}
	
	searchNodes := make(chan int, len(queries))
	
	for w := 0; w < numWorkers; w++ {
	    wg.Add(1)
	    go func() {
	       defer wg.Done()
	       for startNode := range searchNodes {
	           order := Bfs(graph, startNode)
	           mu.Lock()
	           results[startNode] = order
	           mu.Unlock()
	       }
	    }()
	}
	
	for _, query := range queries {
	    searchNodes <- query
	}
	
	close(searchNodes)
	wg.Wait()
	
	
	
	return results
}

func Bfs(graph map[int][]int, startNode int) []int {
    visited := map[int]bool{ startNode: true}
    queue := []int{ startNode }
    order := make([]int, 0, len(graph))
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        order = append(order, node)
        for _, nb := range graph[node] {
            if !visited[nb] {
                visited[nb] = true;
                queue = append(queue, nb)
            }
        }
    }
    
    return order
}

func main() {
	graph := map[int][]int{
        0: {1, 2},
        1: {2, 3},
        2: {3},
        3: {4},
        4: {},
    }
    queries := []int{0, 1, 2}
    numWorkers := 2

    results := ConcurrentBFSQueries(graph, queries, numWorkers)
    fmt.Print(results)
}
