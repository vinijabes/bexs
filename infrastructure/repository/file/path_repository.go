package file

import (
	"bexs/domain/model"
	"bexs/interface/repository"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type filePathRepository struct {
	filepath string
	file     *os.File
}

func (fpr *filePathRepository) GetGraph() model.Graph {
	file, err := os.Open(fpr.filepath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)

	var vertexes []model.Vertex = []model.Vertex{}
	var vertexesReference map[model.VertexID]int = make(map[model.VertexID]int)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		price, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		origin := model.VertexID(record[0])
		dest := model.VertexID(record[1])

		if _, exists := vertexesReference[origin]; !exists {
			vertexes = append(vertexes, model.Vertex{
				ID: model.VertexID(origin),
			})

			vertexesReference[origin] = len(vertexes) - 1
		}

		if _, exists := vertexesReference[dest]; !exists {
			vertexes = append(vertexes, model.Vertex{
				ID: model.VertexID(dest),
			})

			vertexesReference[dest] = len(vertexes) - 1
		}

		originVertexIndex := vertexesReference[origin]

		edge := model.Edge{
			Origin: originVertexIndex,
			Dest:   vertexesReference[dest],
			Price:  price,
		}

		vertexes[originVertexIndex].Edges = append(vertexes[originVertexIndex].Edges, edge)
	}

	return model.Graph{Vertexes: vertexes, VertexReference: vertexesReference}
}

func (fpr *filePathRepository) AddEdge(edge model.Edge) error {
	file, err := os.OpenFile(fpr.filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("\n%s", strings.Join([]string{string(edge.Origin), string(edge.Dest), strconv.FormatUint(edge.Price, 10)}, ",")))
	return err
}

func NewFilePathRepository(filepath string) (repository.PathRepository, error) {
	return &filePathRepository{
		filepath: filepath,
	}, nil
}
