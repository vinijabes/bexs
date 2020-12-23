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
	Vertexes []string `json:"connections"`
	Price    int      `json:"price"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (jp *jsonPathPresenter) ShowPath(path model.Path, writer io.Writer) error {
	resp := PathResponse{}

	var vertexes []string = make([]string, len(path.Connections))

	for i := range path.Connections {
		vertexes[i] = string(path.Connections[i])
	}

	resp.Vertexes = vertexes
	resp.Price = int(path.Dist)

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
