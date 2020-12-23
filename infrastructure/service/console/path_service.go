package console

import (
	"bexs/domain/exceptions"
	"bexs/infrastructure/service"
	"bexs/interface/interactor"
	"bexs/interface/presenter"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

type consoleService struct {
	interactor interactor.PathInteractor
	presenter  presenter.PathPresenter

	done  chan struct{}
	input chan string
}

func (cs *consoleService) Start(context context.Context) <-chan struct{} {
	done := make(chan struct{})
	cs.done = done

	go cs.run(context)

	return done
}

func (cs *consoleService) readInput() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("please enter the route: ")
	input, err := reader.ReadString('\n')
	if err == io.EOF {
		close(cs.input)
		return
	}

	cs.input <- input
}

func (cs *consoleService) run(context context.Context) {
	defer func() {
		cs.done <- struct{}{}
	}()

	for {
		go cs.readInput()

		select {
		case <-context.Done():
			return
		case input := <-cs.input:
			input = strings.Replace(input, "\n", "", -1)
			input = strings.Replace(input, "\r", "", -1)
			segments := strings.Split(input, "-")

			if len(segments) != 2 {
				cs.presenter.ShowException(exceptions.ErrBadParameters, os.Stdout)
				continue
			}

			origin, dest := segments[0], segments[1]

			path, err := cs.interactor.FindPath(interactor.FindPathRequest{
				Origin: origin,
				Dest:   dest,
			})

			if err != nil {
				cs.presenter.ShowException(err, os.Stdout)
				continue
			}

			cs.presenter.ShowPath(path, os.Stdout)
		}
	}
}

func NewConsoleService(interactor interactor.PathInteractor, presenter presenter.PathPresenter) service.Service {
	return &consoleService{
		interactor: interactor,
		presenter:  presenter,
		input:      make(chan string),
	}
}
