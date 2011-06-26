package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
import "C"
import "unsafe"

func NewContext() *Context {
	const c_nil = unsafe.Pointer(uintptr(0))

	ctx := new(Context)
	
	ctx.ref = C.JSContextRef(C.JSGlobalContextCreate((*[0]uint8)(c_nil)))
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
	return ctx.NewObject(ret)
}
