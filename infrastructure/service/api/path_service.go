package api

import (
	"bexs/infrastructure/service"
	"bexs/interface/interactor"
	"bexs/interface/presenter"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	router.HandleFunc("/v1/path", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var req interactor.FindPathRequest

		values := r.URL.Query()

		req.Origin = values.Get("origin")
		req.Dest = values.Get("destination")

		path, err := rs.interactor.FindPath(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rs.presenter.ShowException(err, w)
			return
		}

		rs.presenter.ShowPath(path, w)
	})

	router.HandleFunc("/v1/route", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var req interactor.AddRouteRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rs.presenter.ShowException(err, w)
			return
		}

		err = rs.interactor.AddRoute(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			rs.presenter.ShowException(err, w)
			return
		}

		w.WriteHeader(200)
	})

	port, portExists := os.LookupEnv("PORT")
	if !portExists {
		port = ":8080"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	rs.server = http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		<-ctx.Done()

		tctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rs.server.Shutdown(tctx)

		rs.done <- struct{}{}
	}()

	fmt.Printf("REST API listening on %s\n", port)
	err := rs.server.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func NewRestService(interactor interactor.PathInteractor, presenter presenter.PathPresenter) service.Service {
	return &restService{
		interactor: interactor,
		presenter:  presenter,
	}
}
