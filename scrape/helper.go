package scrape

import (
	"reflect"
	"slices"
)

func checkNil(a any, varName string) error {
	if a == nil || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil()) {
		return GetNilErr(varName)
	}
	av := reflect.ValueOf(a)
	kinds := []reflect.Kind{reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice}
	if slices.Contains(kinds, av.Kind()) && av.IsNil() {
		return GetNilErr(varName)
	}
	return nil
}

func mapSlice[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
