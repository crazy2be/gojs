package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSValueRef.h>
import "C"
import "os"
import "unsafe"

type Error struct {
	Name    string
	Message string
	Context *Context
	Value   *Value
}

func (e *Error) String() string {
	return e.Name + ": " + e.Message
}

func newPanicError(ctx *Context, value *Value) *Error {
	typ := ctx.ValueType(value)

	if typ == TypeString || typ == TypeNumber || typ == TypeBoolean {
		var exception C.JSValueRef
		ret := C.JSValueToStringCopy(C.JSContextRef(unsafe.Pointer(ctx)), C.JSValueRef(unsafe.Pointer(value)), &exception)
		if exception != nil {
			// An error occurred during extraction of string
			// Let's not go to far down the rabbit hole
			panic(os.ENOMEM)
		}
		defer C.JSStringRelease(ret)

		return &Error{"Error", (*String)(unsafe.Pointer(ret)).String(), ctx, value}
	}

	if typ == TypeObject {
		obj := (*Object)(unsafe.Pointer(value))

		name := ""
		prop, _ := ctx.GetProperty(obj, "name")
		if prop != nil {
			name = ctx.ToStringOrDie(prop)
		} else {
			name = "Error"
		}

		msg := ""
		prop, _ = ctx.GetProperty(obj, "message")
		if prop != nil {
			msg = ctx.ToStringOrDie(prop)
		} else {
			msg = "Unknown error"
		}

		return &Error{name, msg, ctx, value}
	}

	// Not certain what else to make of the error
	return &Error{"Error", "Unknown error", ctx, value}
}
