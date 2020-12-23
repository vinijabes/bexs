package main

import (
	"bexs/infrastructure/presenter/buffer"
	"bexs/infrastructure/presenter/json"
	"bexs/infrastructure/repository/cache"
	"bexs/infrastructure/repository/file"
	"bexs/infrastructure/service/api"
	"bexs/infrastructure/service/console"
	"bexs/interface/interactor"
	"bexs/internal/path"
	"bexs/usecase"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing input file")
		return
	}

	var interactor interactor.PathInteractor

	useCache, useCachePresent := os.LookupEnv("USE_CACHE")

	repo, err := file.NewFilePathRepository(os.Args[1])
	if err != nil {
		panic(err)
	}

	if useCachePresent && useCache != "false" {
		observable := path.PathObservable{}

		cache := cache.NewPathCacheRepository()
		if err != nil {
			panic(err)
		}

		observable.Add(cache)

		interactor, err = usecase.NewPathInteractorWithCache(repo, cache, observable)
	} else {
		interactor, err = usecase.NewPathInteractor(repo)
	}

	if err != nil {
		panic(err)
	}

	consolePresenter := buffer.NewBufferPathPresenter()
	jsonPresenter := json.NewJsonPathPresenter()

	ctx, cancel := context.WithCancel(context.Background())

	console := console.NewConsoleService(interactor, consolePresenter)
	api := api.NewRestService(interactor, jsonPresenter)

	apiDone := api.Start(ctx)
	consoleDone := console.Start(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	cancel()

	<-consoleDone
	<-apiDone

	fmt.Println("Finished")
}
