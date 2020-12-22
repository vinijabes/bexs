package file

import (
	"bexs/domain/model"
	"bexs/interface/repository"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type filePathRepository struct {
	filepath string
	file     *os.File
}

func (fpr *filePathRepository) GetGraph() (model.Graph, error) {
	file, err := os.Open(fpr.filepath)
	if err != nil {
		return model.Graph{}, err
	}

	r := csv.NewReader(file)
	graph := model.NewGraph()

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return model.Graph{}, err
		}

		price, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return model.Graph{}, err
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
	file, err := os.OpenFile(fpr.filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("\n%s", strings.Join([]string{string(route.Origin), string(route.Dest), strconv.FormatUint(route.Price, 10)}, ",")))
	return err
}

func NewFilePathRepository(filepath string) (repository.PathRepository, error) {
	return &filePathRepository{
		filepath: filepath,
	}, nil
}
