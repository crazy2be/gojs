package gojs

// #include <stdlib.h>
// #cgo pkg-config: webkit-1.0
// #include <JavaScriptCore/JSBase.h>
import "C"
import (
	"unsafe"
	"log"
)

//=========================================================
// ContextRef
//

type ContextGroup struct {

}

type Context struct {
	ref C.JSContextRef
}

type GlobalContext Context

func (ctx *Context) EvaluateScript(script string, obj *Object, source_url string, startingLineNumber int) (*Value, *Exception) {
	scriptRef := NewString(script)
	defer scriptRef.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = NewString(source_url)
		defer sourceRef.Release()
	}
	
	if obj == nil {
		obj = ctx.NewEmptyObject()
	}

	exception := ctx.NewException()

	log.Println("About to evaluate script:", script, obj, source_url, startingLineNumber)
	
	ret := C.JSEvaluateScript(ctx.ref,
		C.JSStringRef(unsafe.Pointer(scriptRef)), obj.ref,
		C.JSStringRef(unsafe.Pointer(sourceRef)), C.int(startingLineNumber), &exception.val.ref)
	if ret == nil {
		// An error occurred
		// Error information should be stored in exception
		return nil, exception
	}

	// Successful evaluation
	return ctx.newValue(ret), nil
}

func (ctx *Context) CheckScriptSyntax(script string, source_url string, startingLineNumber int) *Exception {
	scriptRef := NewString(script)
	defer scriptRef.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = NewString(source_url)
		defer sourceRef.Release()
	}

	var exception = ctx.NewException()

	ret := C.JSCheckScriptSyntax(ctx.ref,
		C.JSStringRef(unsafe.Pointer(scriptRef)), C.JSStringRef(unsafe.Pointer(sourceRef)),
		C.int(startingLineNumber), &exception.val.ref)
	if !ret {
		// A syntax error was found
		// exception should be non-nil
		return exception
	}

	// exception should be nil
	return nil
}

func (ctx *Context) GarbageCollect() {
	C.JSGarbageCollect(ctx.ref)
}
