package model

type VertexID string

type Edge struct {
	Origin int
	Dest   int
	Price  uint64
}

type Vertex struct {
	ID    VertexID
	Edges []Edge
}

type Graph struct {
	Vertexes        []Vertex
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

func (g Graph) CalculatePath(start VertexID, end VertexID) ([]Vertex, int, error) {
	var vertexes []Vertex = []Vertex{}

	var vertexCount = len(g.Vertexes)
	var invalidIndex = int(vertexCount + 1)

	var startIndex = g.VertexReference[start]
	var endIndex = g.VertexReference[end]

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

	for i, j := 0, len(vertexes)-1; i < j; i, j = i+1, j-1 {
		vertexes[i], vertexes[j] = vertexes[j], vertexes[i]
	}

	return vertexes, dist[endIndex], nil
}
