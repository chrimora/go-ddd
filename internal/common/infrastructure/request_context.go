package commoninfrastructure

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type requestContextKeyType struct{}

var RequestContextKey = requestContextKeyType{}

const (
	requestIdMetaKey = "requestId"
	userIdMetaKey    = "userId"
)

type RequestContext struct {
	RequestId uuid.UUID `json:"request_id"`
	UserId    uuid.UUID `json:"user_id"`
}

func MustGetRequestCtx(ctx context.Context) RequestContext {
	rc, ok := ctx.Value(RequestContextKey).(RequestContext)
	if !ok {
		panic("Request context not found.")
	}
	return rc
}

func NewRequestCtxFromJson(data []byte) (rc RequestContext, err error) {
	return rc, json.Unmarshal(data, &rc)
}

func NewRequestCtxFromMessage(metadata message.Metadata) RequestContext {
	return RequestContext{
		RequestId: uuid.MustParse(metadata.Get(requestIdMetaKey)),
		UserId:    uuid.MustParse(metadata.Get(userIdMetaKey)),
	}
}

func (rc RequestContext) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestContextKey, rc)
}

func (rc RequestContext) ToMessageMetadata(metadata message.Metadata) {
	metadata.Set(requestIdMetaKey, rc.RequestId.String())
	metadata.Set(userIdMetaKey, rc.UserId.String())
}
