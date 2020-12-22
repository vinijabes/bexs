package presenter

import "io"

type PathPresenter interface {
	ShowPath(io.Writer) error
}
