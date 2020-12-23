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

func (bp *bufferPathPresenter) ShowPath(path model.Path, writer io.Writer) error {
	var vertexes []string = make([]string, len(path.Connections))

	for i := range path.Connections {
		vertexes[i] = string(path.Connections[i])
	}

	_, err := writer.Write([]byte(fmt.Sprintf("best route: %s > $%d\n", strings.Join(vertexes, " - "), path.Dist)))
	return err
}

func (bp *bufferPathPresenter) ShowException(err error, writer io.Writer) error {
	_, err = writer.Write([]byte(err.Error() + "\n"))
	return err
}

func NewBufferPathPresenter() presenter.PathPresenter {
	return &bufferPathPresenter{}
}
