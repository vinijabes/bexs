package interactor

import "bexs/domain/model"

//FindPathRequest ...
type FindPathRequest struct {
	Origin string `json:"origin"`
	Dest   string `json:"destination"`
}

//AddRouteRequest ...
type AddRouteRequest struct {
	Origin string `json:"origin"`
	Dest   string `json:"destination"`
	Price  int    `json:"price"`
}

//PathInteractor contains the usecases methods
type PathInteractor interface {
	//FindPath calculates path between two points
	FindPath(FindPathRequest) (model.Path, error)
	//AddRoute adds new route in the graph
	AddRoute(AddRouteRequest) error
}
