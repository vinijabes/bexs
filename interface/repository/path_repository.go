package repository

import "bexs/domain/model"

type PathRepository interface {
	GetGraph() model.Graph
	AddEdge(edge model.Edge) error
}
