package common

import (
	"context"
	"errors"
	"fmt"
)

const (
	RequestIdKey = "requestId"
	UserIdKey    = "userId"
)

type ServiceContext struct {
	RequestId string `json:"requestId"`
	UserId    string `json:"userId"`
}

var ErrBadContext = errors.New("bad context")

func NewServiceCtx(ctx context.Context) (ServiceContext, error) {
	requestId, ok := ctx.Value(RequestIdKey).(string)
	if !ok || requestId == "" {
		return ServiceContext{}, fmt.Errorf("%w: %s", ErrBadContext, RequestIdKey)
	}

	userId, ok := ctx.Value(UserIdKey).(string)
	if !ok || userId == "" {
		// return ServiceContext{}, fmt.Errorf("%w: %s", ErrBadContext, UserIdKey)
	}

	return ServiceContext{
		RequestId: requestId,
		UserId:    userId,
	}, nil
}
