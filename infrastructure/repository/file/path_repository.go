package file

import (
	"bexs/domain/exceptions"
	"bexs/domain/model"
	"bexs/interface/repository"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type filePathRepository struct {
	buffer io.ReadWriteSeeker

	mu sync.Mutex
}

func (fpr *filePathRepository) GetGraph() (model.Graph, error) {
	fpr.mu.Lock()
	defer fpr.mu.Unlock()

	fpr.buffer.Seek(0, io.SeekStart)
	r := csv.NewReader(fpr.buffer)

	graph := model.NewGraph()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return model.Graph{}, exceptions.ErrInvalidInputFile
		}

		price, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return model.Graph{}, exceptions.ErrInvalidInputFile
		}

		origin := model.VertexID(record[0])
		dest := model.VertexID(record[1])

		graph.AddRoute(model.Route{
			Origin: origin,
			Dest:   dest,
			Price:  price,
		})
	}

	return *graph, nil
}

func (fpr *filePathRepository) AddRoute(route model.Route) error {
	fpr.mu.Lock()
	defer fpr.mu.Unlock()

	pos, err := fpr.buffer.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	routeData := strings.Join([]string{string(route.Origin), string(route.Dest), strconv.FormatUint(route.Price, 10)}, ",")

	if pos > 0 {
		_, err = fpr.buffer.Write([]byte(fmt.Sprintf("\n%s", routeData)))
	} else {
		_, err = fpr.buffer.Write([]byte(routeData))
	}
	return err
}

func NewFilePathRepository(filepath string) (repository.PathRepository, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return NewBufferPathRepository(file)
}

func NewBufferPathRepository(buffer io.ReadWriteSeeker) (repository.PathRepository, error) {
	return &filePathRepository{
		buffer: buffer,
	}, nil
}
