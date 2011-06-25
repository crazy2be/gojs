package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSValueRef.h>
import "C"
import "os"
import "fmt"
import "unsafe"

type Error struct {
	Name    string
	Message string
	Context *Context
	Value   *Value
}

type Exception struct {
	msg string // Code error value, string.
	val *Value // Javascript error value, could be any type
	ctx *Context
}

// Used for reporting errors in javascipt code to go code
func (ctx *Context) NewException() *Exception {
	err := new(Exception)
	err.ctx = ctx
	return err
}

// Attempts to convert the error to a string. Pretty-prints with %#v if unable to.
func (e *Exception) String() string {
	str, err := e.ctx.ToString(e.val)
	if err != nil {
		return fmt.Sprintf("%#v (failed to convert to string) %s", e, e.msg)
	}
	return fmt.Sprintf("%#v (string representation: %s %s", e, str, e.msg)
}

func newPanicError(ctx *Context, value *Value) *Error {
	typ := ctx.ValueType(value)

	if typ == TypeString || typ == TypeNumber || typ == TypeBoolean {
		var exception C.JSValueRef
		ret := C.JSValueToStringCopy(ctx.ref, C.JSValueRef(unsafe.Pointer(value)), &exception)
		if exception != nil {
			// An error occurred during extraction of string
			// Let's not go to far down the rabbit hole
			panic(os.ENOMEM)
		}
		defer C.JSStringRelease(ret)

		return &Error{"", (*String)(unsafe.Pointer(ret)).String(), ctx, value}
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
	return &Error{"Unknown", "Unknown error", ctx, value}
}
