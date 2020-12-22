package interactor

import "bexs/domain/model"

type FindPathRequest struct {
	Origin string `json: "origin"`
	Dest   string `json: "dest"`
}

type AddRouteRequest struct {
	Origin string `json:"origin"`
	Dest   string `json:"dest"`
	Price  uint64 `json:"price"`
}

type PathInteractor interface {
	FindPath(FindPathRequest) ([]model.GraphVertex, int, error)
	AddRoute(AddRouteRequest)
}
