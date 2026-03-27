package commoninfrastructure

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
)

type requestContextKeyType struct{}

var RequestContextKey = requestContextKeyType{}

const (
	requestIdMetaKey = "requestId"
	userIdMetaKey    = "userId"
)

type RequestContext struct {
	RequestId string `json:"request_id"`
	UserId    string `json:"user_id"`
}

func MustGetRequestCtx(ctx context.Context) RequestContext {
	rc, ok := ctx.Value(RequestContextKey).(RequestContext)
	if !ok {
		panic("request context not found in context — ensure middleware is configured")
	}
	return rc
}

func NewRequestCtxFromJson(data []byte) (rc RequestContext, err error) {
	return rc, json.Unmarshal(data, &rc)
}

func NewRequestCtxFromMessage(metadata message.Metadata) RequestContext {
	return RequestContext{
		RequestId: metadata.Get(requestIdMetaKey),
		UserId:    metadata.Get(userIdMetaKey),
	}
}

func (rc RequestContext) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestContextKey, rc)
}

func (rc RequestContext) ToMessageMetadata(metadata message.Metadata) {
	metadata.Set(requestIdMetaKey, rc.RequestId)
	metadata.Set(userIdMetaKey, rc.UserId)
}
