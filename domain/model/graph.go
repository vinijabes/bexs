package model

import (
	"bexs/domain/exceptions"
)

type VertexID string

type Route struct {
	Origin VertexID
	Dest   VertexID
	Price  uint64
}

type GraphEdge struct {
	Origin int
	Dest   int
	Price  uint64
}

type GraphVertex struct {
	ID    VertexID
	Edges []GraphEdge
}

type Graph struct {
	Vertexes        []GraphVertex
	VertexReference map[VertexID]int
}

const maxInt = int(^uint(0) >> 1)

func (g Graph) minDistance(dist []int, set []bool) int {
	var min int = maxInt
	var minIndex int

	for v := range dist {
		if set[v] == false && dist[v] < min {
			min = dist[v]
			minIndex = v
		}
	}

	return minIndex
}

func (g *Graph) AddRoute(edge Route) {
	if _, exists := g.VertexReference[edge.Origin]; !exists {
		g.Vertexes = append(g.Vertexes, GraphVertex{
			ID: edge.Origin,
		})

		g.VertexReference[edge.Origin] = len(g.Vertexes) - 1
	}

	if _, exists := g.VertexReference[edge.Dest]; !exists {
		g.Vertexes = append(g.Vertexes, GraphVertex{
			ID: edge.Dest,
		})

		g.VertexReference[edge.Dest] = len(g.Vertexes) - 1
	}

	originVertexIndex := g.VertexReference[edge.Origin]

	graphEdge := GraphEdge{
		Origin: originVertexIndex,
		Dest:   g.VertexReference[edge.Dest],
		Price:  edge.Price,
	}

	g.Vertexes[originVertexIndex].Edges = append(g.Vertexes[originVertexIndex].Edges, graphEdge)
}

func (g Graph) CalculatePath(start VertexID, end VertexID) ([]GraphVertex, int, error) {
	var vertexes []GraphVertex = []GraphVertex{}

	var vertexCount = len(g.Vertexes)
	var invalidIndex = int(vertexCount + 1)

	var startIndex, endIndex int
	var ok bool

	startIndex, ok = g.VertexReference[start]
	if !ok {
		return nil, 0, exceptions.ErrVertexNotFound
	}

	endIndex, ok = g.VertexReference[end]
	if !ok {
		return nil, 0, exceptions.ErrVertexNotFound
	}

	var set []bool = make([]bool, vertexCount)
	var dist []int = make([]int, vertexCount)
	var prev []int = make([]int, vertexCount)

	for i := range g.Vertexes {
		dist[i] = maxInt
		prev[i] = invalidIndex
		set[i] = false
	}

	dist[startIndex] = 0

	for i := 0; i < vertexCount-1; i++ {
		var currentVertexIndex int = g.minDistance(dist, set)
		set[currentVertexIndex] = true

		if currentVertexIndex == endIndex {
			break
		}

		for _, edge := range g.Vertexes[currentVertexIndex].Edges {
			temporaryDist := dist[currentVertexIndex] + int(edge.Price)
			if !set[edge.Dest] && temporaryDist > 0 && temporaryDist < dist[edge.Dest] {
				dist[edge.Dest] = temporaryDist
				prev[edge.Dest] = currentVertexIndex
			}
		}
	}

	var index int = endIndex
	for index != invalidIndex {
		vertexes = append(vertexes, g.Vertexes[index])
		index = prev[index]
	}

	if vertexes[len(vertexes)-1].ID != start {
		return nil, 0, exceptions.ErrPathNotFound
	}

	for i, j := 0, len(vertexes)-1; i < j; i, j = i+1, j-1 {
		vertexes[i], vertexes[j] = vertexes[j], vertexes[i]
	}

	return vertexes, dist[endIndex], nil
}

func NewGraph() *Graph {
	return &Graph{
		Vertexes:        []GraphVertex{},
		VertexReference: make(map[VertexID]int),
	}
}
