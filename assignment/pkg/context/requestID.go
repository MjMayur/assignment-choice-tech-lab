package context

import (
	"context"
	"errors"

	"project/entity"
)

const (
	EMPTY_STRING               = ""
	EMPTY_REQUEST_ID           = "no value for the context key " + REQUEST_ID
	ERR_NO_VAL_SESSION         = "no value for the context key " + SESSION_DATA
	ERR_NO_VALID_VALUE_SESSION = "no valid value for the context key " + SESSION_DATA
)

func GetRequestIDFromContext(ctx context.Context) (string, error) {
	// log.Print("rubaru>pkg>context>requestID.go: GetRequestIDFromContext started")

	reqID := ctx.Value(ContextKey(REQUEST_ID))
	if reqID == nil {
		return EMPTY_STRING, errors.New(EMPTY_REQUEST_ID)
	}

	reqIDString, ok := reqID.(string)
	if !ok {
		return reqIDString, errors.New(EMPTY_REQUEST_ID)
	}

	// log.Print("rubaru>pkg>context>requestID.go: GetRequestIDFromContext completed")
	return reqIDString, nil
}

func GetTokenFromContext(ctx context.Context) (string, error) {

	token := ctx.Value(ContextKey(TOKEN))
	if token == nil {
		return "", errors.New(EMPTY_REQUEST_ID)
	}

	tokenStr, ok := token.(string)
	if !ok {
		return "", errors.New(EMPTY_REQUEST_ID)
	}

	return tokenStr, nil
}

func GetSessionFromContext(ctx context.Context) (entity.Session, error) {
	session := ctx.Value(ContextKey(SESSION_DATA))
	if session == nil {
		return entity.Session{}, errors.New(ERR_NO_VAL_SESSION)
	}

	sessionEntity, ok := session.(*entity.Session)
	if !ok {
		return entity.Session{}, errors.New(ERR_NO_VALID_VALUE_SESSION)
	}

	return *sessionEntity, nil
}
