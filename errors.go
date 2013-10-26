package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include <JavaScriptCore/JSValueRef.h>
// #include "callback.h"
import "C"

// NewError constructs a new JavaScript Error object with message.
func (ctx *Context) NewError(message string) (*Object, error) {
	errVal := ctx.newErrorValue()
	msg := ctx.NewStringValue(message)
	ret := C.JSObjectMakeError(ctx.ref, C.size_t(1), &msg.ref, &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) newErrorObjectOrValue(err error) C.JSValueRef {
	errObj, err := ctx.NewError(err.Error())
	if err == nil {
		// Any JSObjectRef can be safely cast to a JSValueRef.
		// https://lists.webkit.org/pipermail/webkit-dev/2009-May/007530.html
		return C.JSValueRef(errObj.ref)
	}
	// If we failed to construct a new Error, fall back to just using a
	// string and hope it works. We might be out of memory.
	return ctx.NewStringValue(err.Error()).ref
}

type errorValue struct {
	ctx *Context
	ref C.JSValueRef
}

func (ctx *Context) newErrorValue() *errorValue {
	return &errorValue{ctx, nil}
}

// Error returns a string describing the exception. If r.ref is nil, it panics.
func (r errorValue) Error() string {
	if r.ref == nil {
		panic("errorValue.ref is nil")
	}
	v := r.ctx.newValue(r.ref)
	if r.ctx.IsString(v) {
		return r.ctx.ToStringOrDie(v)
	}
	return "obj"
}
