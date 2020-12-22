package service

import "context"

type Service interface {
	Start(context.Context) <-chan struct{}
}
