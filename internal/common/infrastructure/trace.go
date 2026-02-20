package commoninfrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
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

func NewTraceCtxFromCtx(ctx context.Context) (TraceContext, error) {
	var tc TraceContext

	requestId, ok := ctx.Value(RequestIdKey).(string)
	if ok {
		tc.RequestId = requestId
	} else {
		return tc, fmt.Errorf("%w: %s", ErrBadContext, RequestIdKey)
	}

	userId, ok := ctx.Value(UserIdKey).(string)
	if ok {
		tc.UserId = userId
	} else {
		return tc, fmt.Errorf("%w: %s", ErrBadContext, UserIdKey)
	}

	return tc, nil
}

func NewTraceCtxFromJson(ctx []byte) (tc TraceContext, err error) {
	return tc, json.Unmarshal(ctx, &tc)
}

func NewTraceCtxFromMessage(metadata message.Metadata) TraceContext {
	return TraceContext{
		RequestId: metadata.Get(RequestIdKey),
		UserId:    metadata.Get(UserIdKey),
	}
}

func (tc *TraceContext) ToCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, RequestIdKey, tc.RequestId)
	ctx = context.WithValue(ctx, UserIdKey, tc.UserId)
	return ctx
}
