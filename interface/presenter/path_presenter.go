package presenter

import (
	"bexs/domain/model"
	"io"
)

type PathPresenter interface {
	ShowPath([]model.GraphVertex, int, io.Writer) error
	ShowException(error, io.Writer) error
}
