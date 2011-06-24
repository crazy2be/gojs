package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
import "C"
import "unsafe"

func NewContext() *Context {
	const c_nil = unsafe.Pointer(uintptr(0))

	ctx := C.JSGlobalContextCreate((*[0]uint8)(c_nil))
	return (*Context)(unsafe.Pointer(ctx))
}

func (ctx *Context) Retain() {
	C.JSGlobalContextRetain(C.JSContextRef(unsafe.Pointer(ctx)))
}

func (ctx *Context) Release() {
	C.JSGlobalContextRelease(C.JSContextRef(unsafe.Pointer(ctx)))
}

func (ctx *Context) GlobalObject() *Object {
	ret := C.JSContextGetGlobalObject(C.JSContextRef(unsafe.Pointer(ctx)))
	return (*Object)(unsafe.Pointer(ret))
}
