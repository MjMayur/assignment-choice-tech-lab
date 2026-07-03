package repo

import (
	"context"

	projectContext "project/pkg/context"
	"project/pkg/errors"
	"project/pkg/log"
)

type HomeRepoImpl struct{}

func NewHomeRepoImpl() (*HomeRepoImpl, error) {
	return &HomeRepoImpl{}, nil
}

func (r *HomeRepoImpl) Home(ctx context.Context) errors.Response {
	reqID, _ := projectContext.GetRequestIDFromContext(ctx)
	log.Info("core>repo>login: GetUserDetailsForLogin started", reqID)

	log.Info("core>repo>login: GetUserDetailsForLogin completed", reqID)
	return nil
}
