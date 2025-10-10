package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	RequestIdKey = "requestId"
	UserIdKey    = "userId"
)

type TraceContext struct {
	RequestId string `json:"request_id"`
	UserId    string `json:"user_id"`
}

var ErrBadContext = errors.New("bad trace context")

func NewTraceCtx(ctx context.Context) (tc TraceContext) {
	requestId, ok := ctx.Value(RequestIdKey).(string)
	if ok {
		tc.RequestId = requestId
	}
	userId, ok := ctx.Value(UserIdKey).(string)
	if ok {
		tc.UserId = userId
	}
	return tc
}
func (tc *TraceContext) IsComplete() error {
	if tc.RequestId == "" {
		return fmt.Errorf("%w: %s", ErrBadContext, RequestIdKey)
	}
	if tc.UserId == "" {
		// return fmt.Errorf("%w: %s", ErrBadContext, UserIdKey)
	}
	return nil
}

func JsonToCtx(jsonCtx []byte, ctx context.Context) context.Context {
	var trace TraceContext
	err := json.Unmarshal(jsonCtx, &trace)
	if err != nil {
		// TODO; shouldnt happen
		return ctx
	}

	ctx = context.WithValue(ctx, RequestIdKey, trace.RequestId)
	ctx = context.WithValue(ctx, UserIdKey, trace.UserId)
	return ctx
}
