package buffer

import (
	"bexs/domain/model"
	"bexs/interface/presenter"
	"fmt"
	"io"
	"strings"
)

type bufferPathPresenter struct {
}

func (bp *bufferPathPresenter) ShowPath(path []model.GraphVertex, price int, writer io.Writer) error {
	var vertexes []string = make([]string, len(path))

	for i := range path {
		vertexes[i] = string(path[i].ID)
	}

	_, err := writer.Write([]byte(fmt.Sprintf("best route: %s > $%d\n", strings.Join(vertexes, " - "), price)))
	return err
}

func (bp *bufferPathPresenter) ShowException(err error, writer io.Writer) error {
	_, err = writer.Write([]byte(err.Error() + "\n"))
	return err
}

func NewBufferPathPresenter() presenter.PathPresenter {
	return &bufferPathPresenter{}
}
