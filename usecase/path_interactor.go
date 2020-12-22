package usecase

import (
	"bexs/domain/model"
	"bexs/interface/interactor"
	"bexs/interface/repository"
)

type pathInteractor struct {
	graph model.Graph

	repository repository.PathRepository
}

func (pi *pathInteractor) FindPath(req interactor.FindPathRequest) ([]model.GraphVertex, int, error) {
	path, dist, err := pi.graph.CalculatePath(model.VertexID(req.Origin), model.VertexID(req.Dest))

	if err != nil {
		return nil, 0, err
	}

	return path, dist, nil
}

func (pi *pathInteractor) AddRoute(req interactor.AddRouteRequest) {
	route := model.Route{
		Origin: model.VertexID(req.Origin),
		Dest:   model.VertexID(req.Dest),
		Price:  req.Price,
	}

	pi.repository.AddRoute(route)
	pi.graph.AddRoute(route)
}

func NewPathInteractor(repository repository.PathRepository) (interactor.PathInteractor, error) {
	graph, err := repository.GetGraph()

	if err != nil {
		return nil, err
	}

	return &pathInteractor{repository: repository, graph: graph}, nil
}
