package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
import "C"
import "unsafe"

// Context wraps a JavaScriptCore JSContextRef.
type Context struct {
	ref C.JSContextRef
}

// GlobalContext wraps a JavaScriptCore JSGlobalContextRef.
type GlobalContext Context

func NewContext() *Context {
	const c_nil = unsafe.Pointer(uintptr(0))

	ctx := new(Context)

	ctx.ref = C.JSContextRef(C.JSGlobalContextCreate((C.JSClassRef)(c_nil)))
	return ctx
}

type RawContext C.JSContextRef

type RawGlobalContext C.JSGlobalContextRef

func NewContextFrom(raw RawContext) *Context {
	ctx := new(Context)
	ctx.ref = C.JSContextRef(raw)
	return ctx
}

func NewGlobalContextFrom(raw RawGlobalContext) *GlobalContext {
	ctx := new(GlobalContext)
	ctx.ref = C.JSContextRef(raw)
	return ctx
}

func (ctx *Context) Retain() {
	C.JSGlobalContextRetain(ctx.ref)
}

func (ctx *Context) Release() {
	C.JSGlobalContextRelease(ctx.ref)
}

func (ctx *Context) GlobalObject() *Object {
	ret := C.JSContextGetGlobalObject(ctx.ref)
	return ctx.newObject(ret)
}
