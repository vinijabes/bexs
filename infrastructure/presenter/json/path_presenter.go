package json

import (
	"bexs/domain/model"
	"bexs/interface/presenter"
	"encoding/json"
	"io"
)

type jsonPathPresenter struct {
}

type PathResponse struct {
	Vertexes []string `json:"vertexes"`
	Price    int      `json:"price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (jp *jsonPathPresenter) ShowPath(path []model.GraphVertex, price int, writer io.Writer) error {
	resp := PathResponse{}

	var vertexes []string = make([]string, len(path))

	for i := range path {
		vertexes[i] = string(path[i].ID)
	}

	resp.Vertexes = vertexes
	resp.Price = price

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (jp *jsonPathPresenter) ShowException(err error, writer io.Writer) error {
	resp := ErrorResponse{
		Message: err.Error(),
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func NewJsonPathPresenter() presenter.PathPresenter {
	return &jsonPathPresenter{}
}
