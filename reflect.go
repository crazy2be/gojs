package gojs

import "reflect"

// NewValue returns a JavaScript value corresponding to a Go value.
func (ctx *Context) NewValue(goValue interface{}) *Value {
	// Handle simple case right off
	if goValue == nil {
		return ctx.NewNullValue()
	}

	return ctx.reflectToJSValue(reflect.ValueOf(goValue))
}
