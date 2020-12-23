package repository

import (
	"bexs/domain/model"
	"bexs/internal/path"
)

type PathRepository interface {
	GetGraph() (model.Graph, error)
	AddRoute(route model.Route) error
}

type PathCacheRepository interface {
	path.PathObserver

	//AddPath adds calculated path in cache
	AddPath(model.Path)
	//GetPath check if path already cached
	GetPath(model.VertexID, model.VertexID) *model.Path

	//AddPathNotFound stores information that inexists a path between origin and destination
	AddPathNotFound(model.VertexID, model.VertexID)
	//InPathNotFound check if no path exists
	InPathNotFound(model.VertexID, model.VertexID) bool
}
