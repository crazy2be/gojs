package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSBase.h>
import "C"
import "unsafe"

//=========================================================
// ContextRef
//

type ContextGroup struct {
}

type Context struct {
	value C.JSGlobalContextRef
}

type GlobalContext Context

func (ctx *Context) EvaluateScript( script string, obj *Object, source_url string, startingLineNumber int ) (*Value, *Value) {
	scriptRef := ctx.NewString( script )
	defer scriptRef.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = ctx.NewString( source_url )
		defer sourceRef.Release()
	}

	var exception C.JSValueRef

	ret := C.JSEvaluateScript( ctx.value, C.JSStringRef(unsafe.Pointer(scriptRef)), C.JSObjectRef(unsafe.Pointer(obj)), 
		C.JSStringRef(unsafe.Pointer(sourceRef)), C.int(startingLineNumber), &exception )
	if ret == nil {
		// An error occurred
		// Error information should be stored in exception
		return nil, (*Value)(unsafe.Pointer( exception ))
	}

	// Successful evaluation
	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) CheckScriptSyntax( script string, source_url string, startingLineNumber int ) *Value {
	scriptRef := ctx.NewString( script )
	defer scriptRef.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = ctx.NewString( source_url )
		defer sourceRef.Release()
	} 

	var exception C.JSValueRef

	ret := C.JSCheckScriptSyntax( ctx.value, C.JSStringRef(unsafe.Pointer(scriptRef)), C.JSStringRef(unsafe.Pointer(sourceRef)), 
		C.int(startingLineNumber), &exception )
	if !ret {
		// A syntax error was found
		// exception should be non-nil
		return (*Value)(unsafe.Pointer(exception))
	}

	// exception should be nil
	return nil
}

func (ctx *Context) GarbageCollect() {
	C.JSGarbageCollect( ctx.value )
}

