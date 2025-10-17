package commontest

import (
	"fmt"
	"reflect"
)

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
