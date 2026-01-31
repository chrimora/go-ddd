package commontest

import (
	"context"
	"fmt"
	commoninfrastructure "goddd/internal/common/infrastructure"
	"reflect"

	"github.com/google/uuid"
)

func TestContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, commoninfrastructure.RequestIdKey, uuid.NewString())
	ctx = context.WithValue(ctx, commoninfrastructure.UserIdKey, uuid.NewString())
	return ctx
}

// obj needs to be pointer
func Merge(obj any, overrides []map[string]any) {
	st := reflect.ValueOf(obj).Elem()
	for _, values := range overrides {
		for k, v := range values {
			f := st.FieldByName(k)
			if !f.IsValid() || !f.CanSet() {
				panic(fmt.Sprintf("Cannot set field: %s in %T", k, obj))
			}
			v := reflect.ValueOf(v)
			f.Set(v)
		}
	}
}
