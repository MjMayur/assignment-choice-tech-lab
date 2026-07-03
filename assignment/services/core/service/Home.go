package service

import (
	"context"

	"project/pkg/errors"
)

type HomeServiceImpl struct{}

func NewHomeServiceImpl() (*HomeServiceImpl, error) {
	return &HomeServiceImpl{}, nil
}

func (s *HomeServiceImpl) Home(ctx context.Context) errors.Response {

	return nil
}
