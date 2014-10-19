package gojs

// #include <stdlib.h>
// #cgo pkg-config: javascriptcoregtk-3.0
// #include <JavaScriptCore/JSBase.h>
import "C"
import (
	"log"
	"unsafe"
)

// EvaluateScript evaluates the JavaScript code in script.
func (ctx *Context) EvaluateScript(script string, thisObject *Object, sourceURL string, startingLineNumber int) (*Value, error) {
	scriptRef := NewString(script)
	defer scriptRef.Release()

	var sourceRef *String
	if sourceURL != "" {
		sourceRef = NewString(sourceURL)
		defer sourceRef.Release()
	}

	if thisObject == nil {
		thisObject = ctx.NewEmptyObject()
	}

	log.Println("About to evaluate script:", script, thisObject, sourceURL, startingLineNumber)

	errVal := ctx.newErrorValue()
	ret := C.JSEvaluateScript(ctx.ref,
		C.JSStringRef(unsafe.Pointer(scriptRef)), thisObject.ref,
		C.JSStringRef(unsafe.Pointer(sourceRef)), C.int(startingLineNumber), &errVal.ref)
	if ret == nil {
		// An error occurred
		// Error information should be stored in exception
		return nil, errVal
	}

	// Successful evaluation
	return ctx.newValue(ret), nil
}

// CheckScriptSyntax checks the JavaScript syntax of script.
func (ctx *Context) CheckScriptSyntax(script string, sourceURL string, startingLineNumber int) error {
	scriptRef := NewString(script)
	defer scriptRef.Release()

	var sourceRef *String
	if sourceURL != "" {
		sourceRef = NewString(sourceURL)
		defer sourceRef.Release()
	}

	errVal := ctx.newErrorValue()
	ret := C.JSCheckScriptSyntax(ctx.ref,
		C.JSStringRef(unsafe.Pointer(scriptRef)), C.JSStringRef(unsafe.Pointer(sourceRef)),
		C.int(startingLineNumber), &errVal.ref)
	if !ret {
		// A syntax error was found
		// exception should be non-nil
		return errVal
	}

	// exception should be nil
	return nil
}

// GarbageCollect performs a JavaScript garbage collection.
func (ctx *Context) GarbageCollect() {
	C.JSGarbageCollect(ctx.ref)
}
