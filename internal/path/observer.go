package path

import "bexs/domain/model"

//PathObserver is a listener for graph changes
type PathObserver interface {
	OnNewRoute(model.Route)
}

//PathObservable is the struct who notifies graph changes
type PathObservable struct {
	observers []PathObserver
}

//Add adds new observable to list
func (po *PathObservable) Add(observer PathObserver) {
	po.observers = append(po.observers, observer)
}

//Notify notifies all observers
func (po *PathObservable) Notify(route model.Route) {
	for _, observer := range po.observers {
		observer.OnNewRoute(route)
	}
}
