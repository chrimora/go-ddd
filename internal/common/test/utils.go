package commontest

import (
	"context"
	"fmt"
	commondomain "goddd/internal/common/domain"
	commoninfrastructure "goddd/internal/common/infrastructure"
	"reflect"
)

func TestContext() context.Context {
	rc := commoninfrastructure.RequestContext{
		RequestId: commondomain.NewUUID().String(),
		UserId:    commondomain.NewUUID().String(),
	}
	return rc.ToCtx(context.Background())
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
