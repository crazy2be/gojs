package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "unsafe"

func release_jsstringref_array( refs []C.JSStringRef ) {
	for i:=0; i<len(refs); i++ {
		if refs[i] != nil {
			C.JSStringRelease( refs[i] )
		}
	}
}

func (ctx *Context) MakeFunction(name string, parameters []string, body string, source_url string, starting_line_number int ) (*Object,*Value) {
	Cname := ctx.NewString( name )
	defer Cname.Release()

	Cparameters := make( []C.JSStringRef, len(parameters) )
	defer release_jsstringref_array( Cparameters )
	for i:=0; i<len(parameters); i++ {
		Cparameters[i] = (C.JSStringRef)(unsafe.Pointer(ctx.NewString( parameters[i] ) ) )
	}

	Cbody := ctx.NewString( body )
	defer Cbody.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = ctx.NewString( source_url )
		defer sourceRef.Release()
	}

	var exception C.JSValueRef

	ret := C.JSObjectMakeFunction( ctx.value, 
		(C.JSStringRef)( unsafe.Pointer( Cname ) ),
		C.unsigned(len(Cparameters)), &Cparameters[0], 
		(C.JSStringRef)( unsafe.Pointer( Cbody ) ), 
		(C.JSStringRef)( unsafe.Pointer( sourceRef ) ), 
		C.int(starting_line_number), &exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) GetPrototype(obj *Object) *Value {
	ret := C.JSObjectGetPrototype( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)) )
	return (*Value)(unsafe.Pointer( ret ))
}

type PropertyNameArray struct {
}

func (ctx *Context) ObjectHasProperty(obj *Object, name string) bool {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	ret := C.JSObjectHasProperty( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)) )
	return bool(ret)
}

func (ctx *Context) ObjectGetProperty(obj *Object, name string) (*Value,*Value) {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	ret := C.JSObjectGetProperty( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), &exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) ObjectGetPropertyAtIndex(obj *Object, index uint16) (*Value,*Value) {
	var exception C.JSValueRef

	ret := C.JSObjectGetPropertyAtIndex( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)), C.unsigned(index), &exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) ObjectSetProperty(obj *Object, name string, rhs *Value, attributes uint8) *Value {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	C.JSObjectSetProperty( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), C.JSValueRef(unsafe.Pointer(rhs)), 
		(C.JSPropertyAttributes)(attributes), &exception )
	if exception != nil {
		return (*Value)(unsafe.Pointer(exception))
	}

	return nil
}

func (ctx *Context) ObjectSetPropertyAtIndex(obj *Object, index uint16, rhs *Value) *Value {
	var exception C.JSValueRef

	C.JSObjectSetPropertyAtIndex( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)), C.unsigned(index), 
		C.JSValueRef(unsafe.Pointer(rhs)), &exception )
	if exception != nil {
		return (*Value)(unsafe.Pointer(exception))
	}

	return nil
}

func (ctx *Context) ObjectDeleteProperty(obj *Object, name string ) (bool,*Value) {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	ret := C.JSObjectDeleteProperty( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), &exception )
	if exception != nil {
		return false, (*Value)(unsafe.Pointer(exception))
	}

	return bool(ret), nil
}

func (obj *Object) GetPrivate() unsafe.Pointer {
	ret := C.JSObjectGetPrivate( C.JSObjectRef(unsafe.Pointer(obj)) )
	return ret
}

func (obj *Object) SetPrivate(data unsafe.Pointer) bool {
	ret := C.JSObjectSetPrivate( C.JSObjectRef(unsafe.Pointer(obj)), data )
	return bool( ret )
}

func (obj *Object) GetValue() *Value {
	return (*Value)(unsafe.Pointer(obj))
}

func (ctx *Context) IsFunction(obj *Object) bool {
	ret := C.JSObjectIsFunction( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)) )
	return bool(ret)
}

func (ctx *Context) CallAsFunction( obj *Object, thisObject *Object, parameters []*Value ) (*Value,*Value) {
	var exception C.JSValueRef

	ret := C.JSObjectCallAsFunction( ctx.value, 
		C.JSObjectRef( unsafe.Pointer(obj) ),	
		C.JSObjectRef( unsafe.Pointer(thisObject) ),
		C.size_t( len(parameters) ),
		(*C.JSValueRef)( unsafe.Pointer( &parameters[0]) ),
		&exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) IsConstructor(obj *Object) bool {
	ret := C.JSObjectIsConstructor( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)) )
	return bool(ret)
}

