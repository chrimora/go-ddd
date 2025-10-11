package commontest

import (
	"reflect"
)

func Merge(obj any, overrides []map[string]any) {
	// obj needs to be pointer
	st := reflect.ValueOf(obj).Elem()
	for _, values := range overrides {
		for k, v := range values {
			f := st.FieldByName(k)
			v := reflect.ValueOf(v)
			f.Set(v)
		}
	}
}
