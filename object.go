package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "os"
import "runtime"
import "unsafe"

func release_jsstringref_array( refs []C.JSStringRef ) {
	for i:=0; i<len(refs); i++ {
		if refs[i] != nil {
			C.JSStringRelease( refs[i] )
		}
	}
}

var nativecallback C.JSClassRef
var nativecallback_typ runtime.Type

func init() {
	// Create the class definition for JavaScriptCore
	nativecallback = C.JSClassDefinition_NativeCallback()
	if nativecallback == nil {
		panic( os.ENOMEM )
	}

	// Get the Go type information to recreate the callback
	var dummy GoFunctionCallback
	nativecallback_typ, _ = unsafe.Reflect( dummy )
}

type GoFunctionCallback func(ctx *Context, obj *Object, thisObject *Object, arguments []*Value) (ret *Value, err *Value)

func (ctx *Context) MakeFunctionWithCallback( callback GoFunctionCallback ) *Object {
	_, addr := unsafe.Reflect( callback )

	ret := C.JSObjectMake( C.JSContextRef(unsafe.Pointer(ctx)), nativecallback, addr )
	return (*Object)(unsafe.Pointer(ret))
}

//export nativecallback_CallAsFunction_go
func nativecallback_CallAsFunction_go( data unsafe.Pointer, ctx unsafe.Pointer, obj unsafe.Pointer, thisObject unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer ) unsafe.Pointer {
	ret, err := unsafe.Unreflect( nativecallback_typ, data ).(GoFunctionCallback)(
		(*Context)(ctx), (*Object)(obj), (*Object)(thisObject), (*[1<<14]*Value)(arguments)[0:argumentCount] )
	if err != nil {
		*exception = unsafe.Pointer(err)
		return nil
	}

	return unsafe.Pointer(ret)
}

func (ctx *Context) MakeRegExp( regex string ) (*Object,*Value) {
	var exception C.JSValueRef

	param := ctx.NewStringValue( regex )

	ret := C.JSObjectMakeRegExp( C.JSContextRef(unsafe.Pointer(ctx)), 
		C.size_t(1), (*C.JSValueRef)( unsafe.Pointer( &param ) ),
		&exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) MakeRegExpFromValues( parameters []*Value ) (*Object,*Value) {
	var exception C.JSValueRef

	ret := C.JSObjectMakeRegExp( C.JSContextRef(unsafe.Pointer(ctx)), 
		C.size_t(len(parameters)), (*C.JSValueRef)( unsafe.Pointer( &parameters[0] ) ),
		&exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) MakeFunction(name string, parameters []string, body string, source_url string, starting_line_number int ) (*Object,*Value) {
	Cname := NewString( name )
	defer Cname.Release()

	Cparameters := make( []C.JSStringRef, len(parameters) )
	defer release_jsstringref_array( Cparameters )
	for i:=0; i<len(parameters); i++ {
		Cparameters[i] = (C.JSStringRef)(unsafe.Pointer(NewString( parameters[i] ) ) )
	}

	Cbody := NewString( body )
	defer Cbody.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = NewString( source_url )
		defer sourceRef.Release()
	}

	var exception C.JSValueRef

	ret := C.JSObjectMakeFunction( C.JSContextRef(unsafe.Pointer(ctx)), 
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
	ret := C.JSObjectGetPrototype( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)) )
	return (*Value)(unsafe.Pointer( ret ))
}

func (ctx *Context) SetPrototype(obj *Object, rhs *Value) {
	C.JSObjectSetPrototype( C.JSContextRef(unsafe.Pointer(ctx)), 
		C.JSObjectRef(unsafe.Pointer(obj)), C.JSValueRef(unsafe.Pointer(rhs)) )
}

type PropertyNameArray struct {
}

func (ctx *Context) ObjectHasProperty(obj *Object, name string) bool {
	jsstr := NewString( name )
	defer jsstr.Release()

	ret := C.JSObjectHasProperty( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)) )
	return bool(ret)
}

func (ctx *Context) ObjectGetProperty(obj *Object, name string) (*Value,*Value) {
	jsstr := NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	ret := C.JSObjectGetProperty( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), &exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) ObjectGetPropertyAtIndex(obj *Object, index uint16) (*Value,*Value) {
	var exception C.JSValueRef

	ret := C.JSObjectGetPropertyAtIndex( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)), C.unsigned(index), &exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) ObjectSetProperty(obj *Object, name string, rhs *Value, attributes uint8) *Value {
	jsstr := NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	C.JSObjectSetProperty( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), C.JSValueRef(unsafe.Pointer(rhs)), 
		(C.JSPropertyAttributes)(attributes), &exception )
	if exception != nil {
		return (*Value)(unsafe.Pointer(exception))
	}

	return nil
}

func (ctx *Context) ObjectSetPropertyAtIndex(obj *Object, index uint16, rhs *Value) *Value {
	var exception C.JSValueRef

	C.JSObjectSetPropertyAtIndex( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)), C.unsigned(index), 
		C.JSValueRef(unsafe.Pointer(rhs)), &exception )
	if exception != nil {
		return (*Value)(unsafe.Pointer(exception))
	}

	return nil
}

func (ctx *Context) ObjectDeleteProperty(obj *Object, name string ) (bool,*Value) {
	jsstr := NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	ret := C.JSObjectDeleteProperty( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), &exception )
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
	ret := C.JSObjectIsFunction( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)) )
	return bool(ret)
}

func (ctx *Context) CallAsFunction( obj *Object, thisObject *Object, parameters []*Value ) (*Value,*Value) {
	var exception C.JSValueRef

	var Cparameters *C.JSValueRef
	if len(parameters)>0 {
		Cparameters = (*C.JSValueRef)( unsafe.Pointer( &parameters[0]) )
	}

	ret := C.JSObjectCallAsFunction( C.JSContextRef(unsafe.Pointer(ctx)), 
		C.JSObjectRef( unsafe.Pointer(obj) ),	
		C.JSObjectRef( unsafe.Pointer(thisObject) ),
		C.size_t( len(parameters) ),
		Cparameters,
		&exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) IsConstructor(obj *Object) bool {
	ret := C.JSObjectIsConstructor( C.JSContextRef(unsafe.Pointer(ctx)), C.JSObjectRef(unsafe.Pointer(obj)) )
	return bool(ret)
}

func (ctx *Context) CallAsConstructor( obj *Object, parameters []*Value ) (*Value,*Value) {
	var exception C.JSValueRef

	var Cparameters *C.JSValueRef
	if len(parameters)>0 {
		Cparameters = (*C.JSValueRef)( unsafe.Pointer( &parameters[0]) )
	}

	ret := C.JSObjectCallAsConstructor( C.JSContextRef(unsafe.Pointer(ctx)), 
		C.JSObjectRef( unsafe.Pointer(obj) ),	
		C.size_t( len(parameters) ),
		Cparameters,
		&exception )
	if exception != nil {
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

