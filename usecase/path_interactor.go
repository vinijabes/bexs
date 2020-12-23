package usecase

import (
	"bexs/domain/exceptions"
	"bexs/domain/model"
	"bexs/interface/interactor"
	"bexs/interface/repository"
	"bexs/internal/path"
	"sync"
)

type pathInteractor struct {
	graph model.Graph

	observable path.PathObservable
	repository repository.PathRepository
	cache      repository.PathCacheRepository

	mu sync.RWMutex
}

func (pi *pathInteractor) FindPath(req interactor.FindPathRequest) (model.Path, error) {
	pi.mu.RLock()
	defer pi.mu.RUnlock()

	origin, dest := model.VertexID(req.Origin), model.VertexID(req.Dest)

	if pi.cache != nil {
		cachedPath := pi.cache.GetPath(origin, dest)
		if cachedPath != nil {
			return *cachedPath, nil
		}

		if pi.cache.InPathNotFound(origin, dest) {
			return model.Path{}, exceptions.ErrPathNotFound
		}
	}

	path, err := pi.graph.CalculatePath(origin, dest)

	if err != nil {
		pi.cache.AddPathNotFound(origin, dest)
		return model.Path{}, err
	}

	if pi.cache != nil {
		pi.cache.AddPath(*path)
	}

	return *path, nil
}

func (pi *pathInteractor) AddRoute(req interactor.AddRouteRequest) error {
	pi.mu.Lock()
	defer pi.mu.Unlock()

	if req.Origin == "" || req.Dest == "" {
		return exceptions.ErrBadParameters
	}

	if req.Price < 0 {
		return exceptions.ErrBadParameters
	}

	route := model.Route{
		Origin: model.VertexID(req.Origin),
		Dest:   model.VertexID(req.Dest),
		Price:  uint64(req.Price),
	}

	if pi.graph.AddRoute(route) {
		pi.repository.AddRoute(route)

		if pi.cache != nil {
			pi.observable.Notify(route)
		}
	} else {
		return exceptions.ErrRouteAlreadyExists
	}

	return nil
}

//NewPathInteractor returns a instance of path interactor
func NewPathInteractor(repository repository.PathRepository) (interactor.PathInteractor, error) {
	graph, err := repository.GetGraph()

	if err != nil {
		return nil, err
	}

	return &pathInteractor{repository: repository, graph: graph}, nil
}

//NewPathInteractorWithCache return a instance of path interactor with cache enabled
func NewPathInteractorWithCache(repository repository.PathRepository, cache repository.PathCacheRepository, observable path.PathObservable) (interactor.PathInteractor, error) {
	graph, err := repository.GetGraph()

	if err != nil {
		return nil, err
	}

	return &pathInteractor{repository: repository, cache: cache, observable: observable, graph: graph}, nil
}
