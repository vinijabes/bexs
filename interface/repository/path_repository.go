package repository

import "bexs/domain/model"

type PathRepository interface {
	GetGraph() (model.Graph, error)
	AddRoute(route model.Route) error
}
