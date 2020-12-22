package api

import (
	"bexs/infrastructure/service"
	"bexs/interface/interactor"
	"bexs/interface/presenter"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type restService struct {
	interactor interactor.PathInteractor
	presenter  presenter.PathPresenter

	server http.Server
	done   chan struct{}
}

func (rs *restService) Start(context context.Context) <-chan struct{} {
	done := make(chan struct{})
	rs.done = done

	go rs.run(context)

	return done
}

func (rs *restService) run(ctx context.Context) {
	router := http.NewServeMux()
	router.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var req interactor.FindPathRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			rs.presenter.ShowException(err, w)
			return
		}

		path, price, err := rs.interactor.FindPath(req)
		if err != nil {
			rs.presenter.ShowException(err, w)
			return
		}

		rs.presenter.ShowPath(path, price, w)
	})

	router.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var req interactor.AddRouteRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			rs.presenter.ShowException(err, w)
			return
		}

		rs.interactor.AddRoute(req)
		w.WriteHeader(200)
	})

	rs.server = http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-ctx.Done()

		tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rs.server.Shutdown(tctx)

		rs.done <- struct{}{}
	}()

	rs.server.ListenAndServe()
}

func NewRestService(interactor interactor.PathInteractor, presenter presenter.PathPresenter) service.Service {
	return &restService{
		interactor: interactor,
		presenter:  presenter,
	}
}
