package gojs

// #include <stdlib.h>
import "C"
import "reflect"

func (ctx *Context) NewValue(value interface{}) *Value {
	// Handle simple case right off
	if value == nil {
		return ctx.NewNullValue()
	}

	ret := ctx.reflectToJSValue(reflect.ValueOf(value))
	return ret
}
