package cache

import (
	"bexs/domain/model"
	"bexs/interface/repository"
	"fmt"
	"sync"
)

type cachePathRepository struct {
	cache    map[string]model.Path
	notFound map[string]bool

	mu sync.RWMutex
}

func (cp *cachePathRepository) getKey(origin model.VertexID, dest model.VertexID) string {
	return fmt.Sprintf("%s-%s", origin, dest)
}

func (cp *cachePathRepository) AddPath(path model.Path) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	initial := path.Connections[0]
	for i := 1; i < len(path.Connections); i++ {
		cp.cache[cp.getKey(initial, path.Connections[i])] = model.Path{
			Connections: path.Connections[0 : i+1],
			Dist:        path.SegmentDist[i],
		}
	}

}

func (cp *cachePathRepository) GetPath(origin model.VertexID, dest model.VertexID) *model.Path {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	var path model.Path
	var ok bool
	path, ok = cp.cache[cp.getKey(origin, dest)]

	if !ok {
		return nil
	}

	return &path
}

func (cp *cachePathRepository) AddPathNotFound(origin model.VertexID, dest model.VertexID) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.notFound[cp.getKey(origin, dest)] = true
}
func (cp *cachePathRepository) InPathNotFound(origin model.VertexID, dest model.VertexID) bool {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	val, exists := cp.notFound[cp.getKey(origin, dest)]
	return val && exists
}

func (cp *cachePathRepository) OnNewRoute(route model.Route) {
	cp.notFound = make(map[string]bool)

	for key, value := range cp.cache {
		for _, connection := range value.Connections {
			if connection == route.Origin || connection == route.Dest {
				delete(cp.cache, key)
			}
		}
	}
}

func NewPathCacheRepository() repository.PathCacheRepository {
	return &cachePathRepository{
		cache:    make(map[string]model.Path),
		notFound: make(map[string]bool),
	}
}
