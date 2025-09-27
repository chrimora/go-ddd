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

type TraceContext struct {
	RequestId string `json:"requestId"`
	UserId    string `json:"userId"`
}

var ErrBadContext = errors.New("bad trace context")

func NewTraceCtx(ctx context.Context) (TraceContext, error) {
	requestId, ok := ctx.Value(RequestIdKey).(string)
	if !ok || requestId == "" {
		return TraceContext{}, fmt.Errorf("%w: %s", ErrBadContext, RequestIdKey)
	}

	userId, ok := ctx.Value(UserIdKey).(string)
	if !ok || userId == "" {
		// return TraceContext{}, fmt.Errorf("%w: %s", ErrBadContext, UserIdKey)
	}

	return TraceContext{
		RequestId: requestId,
		UserId:    userId,
	}, nil
}
