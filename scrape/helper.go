package scrape

import (
	"reflect"
	"slices"
)

// ValidateNotNil checks if the given o variable is nil or is a pointer
// that points to nil and returns [NillErr] in the true case. varName
// is a string the specifies the name of the variable that is nill.
func ValidateNotNil(o any, varName string) error {
	err := NilErr{Var: varName}
	if o == nil || (reflect.ValueOf(o).Kind() == reflect.Ptr && reflect.ValueOf(o).IsNil()) {
		return err
	}
	av := reflect.ValueOf(o)
	kinds := []reflect.Kind{reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice}
	if slices.Contains(kinds, av.Kind()) && av.IsNil() {
		return err
	}
	return nil
}
