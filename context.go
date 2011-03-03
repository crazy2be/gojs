package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
// #include "callback.h"
import "C"
import "unsafe"

func NewContext() *Context {
	const c_nil = unsafe.Pointer( uintptr(0) )
	ctx := C.JSGlobalContextCreate( (*[0]uint8)(c_nil) );
	return &Context{ ctx }
}

func (ctx *Context) Retain() {
	C.JSGlobalContextRetain( ctx.value )
}

func (ctx *Context) Release() {
	C.JSGlobalContextRelease( ctx.value )
}

func (ctx *Context) GlobalObject() *Object {
	ret := C.JSContextGetGlobalObject( ctx.value )
	return (*Object)( unsafe.Pointer( ret ) )
}

