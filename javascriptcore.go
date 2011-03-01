package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include "callback.h"
import "C"
import "unsafe"

//=========================================================
// ContextRef
//

type Context struct {
	value C.JSGlobalContextRef
	callbacks []Function
}

type Object struct {
}

type Value struct {
}

const (
	TypeUndefined = 0
	TypeNull = iota
	TypeBoolean = iota
	TypeNumber = iota
	TypeString = iota
	TypeObject = iota
)

func NewContext() *Context {
	const c_nil = unsafe.Pointer( uintptr(0) )
	ctx := C.JSGlobalContextCreate( (*[0]uint8)(c_nil) );
	return &Context{ ctx, []Function{} }
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

//=========================================================
// *Value
//

func (ctx *Context) ValueIsUndefined( v *Value ) bool {
	ret := C.JSValueIsUndefined( ctx.value, C.JSValueRef(unsafe.Pointer(v)) )
	return bool( ret )
}

func (ctx *Context) ValueIsNull( v *Value ) bool {
	ret := C.JSValueIsNull( ctx.value, C.JSValueRef(unsafe.Pointer(v)) )
	return bool( ret )
}

func (ctx *Context) ValueIsBoolean( v *Value ) bool {
	ret := C.JSValueIsBoolean( ctx.value, C.JSValueRef(unsafe.Pointer(v)) )
	return bool( ret )
}

func (ctx *Context) ValueIsString( v *Value ) bool {
	ret := C.JSValueIsString( ctx.value, C.JSValueRef(unsafe.Pointer(v)) )
	return bool( ret )
}

func (ctx *Context) ValueIsObject( v *Value ) bool {
	ret := C.JSValueIsObject( ctx.value, C.JSValueRef(unsafe.Pointer(v)) )
	return bool( ret )
}

func (ctx *Context) ValueType( v *Value ) uint8 {
	ret := C.JSValueGetType( ctx.value, C.JSValueRef(unsafe.Pointer(v)) )
	return uint8( ret )
}

func (ctx *Context) IsEqual( a *Value, b *Value ) (bool, *Value) {
	exception := C.JSValueRef(nil)

	ret := C.JSValueIsEqual( ctx.value, C.JSValueRef(unsafe.Pointer(a)), C.JSValueRef(unsafe.Pointer(b)), &exception )
	if exception != nil {
		return false, (*Value)(unsafe.Pointer(exception))
	}

	return bool(ret), nil
}

func (ctx *Context) IsStrictEqual( a *Value, b *Value ) bool {
	ret := C.JSValueIsStrictEqual( ctx.value, C.JSValueRef(unsafe.Pointer(a)), C.JSValueRef(unsafe.Pointer(b)) )
	return bool(ret)
}

func (ctx *Context) NewUndefinedValue() *Value {
	ref := C.JSValueMakeUndefined( ctx.value )
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewNullValue() *Value {
	ref := C.JSValueMakeNull( ctx.value )
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewBooleanValue( value bool ) *Value {
	ref := C.JSValueMakeBoolean( ctx.value, C.bool(value) )
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewNumberValue( value float64 ) *Value {
	ref := C.JSValueMakeNumber( ctx.value, C.double(value) )
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewStringValue( value string ) *Value {
	cvalue := C.CString( value )
	defer C.free( unsafe.Pointer(cvalue) )
	jsstr := C.JSStringCreateWithUTF8CString( cvalue )
	defer C.JSStringRelease( jsstr )
	ref := C.JSValueMakeString( ctx.value, jsstr )
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewString( value string ) *String {
	cvalue := C.CString( value )
	defer C.free( unsafe.Pointer(cvalue) )
	ref := C.JSStringCreateWithUTF8CString( cvalue )
	return (*String)( unsafe.Pointer( ref ) )
}

func (ctx *Context) ProtectValue( ref *Value ) {
	C.JSValueProtect( ctx.value, C.JSValueRef(unsafe.Pointer(ref)) )
}

func (ctx *Context) UnProtectValue( ref **Value ) {
	C.JSValueProtect( ctx.value, C.JSValueRef(unsafe.Pointer(ref)) )
}

func (ctx *Context) ToBoolean( ref *Value ) bool {
	ret := C.JSValueToBoolean( ctx.value, C.JSValueRef(unsafe.Pointer(ref)) )
	return bool( ret )
}

func (ctx *Context) ToNumber( ref *Value ) (num float64,err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToNumber( ctx.value, C.JSValueRef(unsafe.Pointer(ref)), &exception )
	if exception != nil {
		return float64(ret), (*Value)(unsafe.Pointer(exception))
	}

	// Successful conversion
	return float64(ret), nil
}

func (ctx *Context) ToString( ref *Value ) (str string, err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToStringCopy( ctx.value, C.JSValueRef(unsafe.Pointer(ref)), &exception )
	if exception != nil {
		// An error occurred
		// ret should be null
		return "", (*Value)(unsafe.Pointer(exception))
	}
	defer C.JSStringRelease( ret )

	// Successful conversion
	tmp := (*String)( unsafe.Pointer( ret ) )
	return tmp.String(), nil
}

func (ctx *Context) ToStringOrDie( ref *Value ) string {
	str, err := ctx.ToString( ref )
	if err!=nil {
		panic( err )
	}
	return str
}

func (ctx *Context) ToObject( ref *Value ) (obj *Object, err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToObject( ctx.value, C.JSValueRef(unsafe.Pointer(ref)), &exception )
	if exception != nil {
		// An error occurred
		// ret should be null
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	// Successful conversion
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) ToObjectOrDie( ref *Value ) *Object {
	var exception C.JSValueRef
	ret := C.JSValueToObject( ctx.value, C.JSValueRef(unsafe.Pointer(ref)), &exception )
	if exception != nil {
		panic( (*Value)(unsafe.Pointer(exception)) )
	}

	// Successful conversion
	return (*Object)(unsafe.Pointer(ret))
}

//=========================================================
// *Object
//

const (
	PropertyAttributeNone = 0
	PropertyAttributeReadOnly     = 1 << 1
	PropertyAttributeDontEnum     = 1 << 2
	PropertyAttributeDontDelete   = 1 << 3
)

const (
	ClassAttributeNone = 0
	ClassAttributeNoAutomaticPrototype = 1 << 1
)

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

func (ctx *Context) IsConstructor(obj *Object) bool {
	ret := C.JSObjectIsConstructor( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)) )
	return bool(ret)
}

func (ctx *Context) CopyPropertyNames(obj *Object) *PropertyNameArray {
	ret := C.JSObjectCopyPropertyNames( ctx.value, C.JSObjectRef(unsafe.Pointer(obj)) )
	return (*PropertyNameArray)(unsafe.Pointer( ret ))
}

func (ref *PropertyNameArray) Retain() {
	C.JSPropertyNameArrayRetain( C.JSPropertyNameArrayRef(unsafe.Pointer(ref)) )
}

func (ref *PropertyNameArray) Release() {
	C.JSPropertyNameArrayRelease( C.JSPropertyNameArrayRef(unsafe.Pointer(ref)) )
}

func (ref *PropertyNameArray) Count() uint16 {
	ret := C.JSPropertyNameArrayGetCount( C.JSPropertyNameArrayRef(unsafe.Pointer(ref)) )
	return uint16( ret )
}

func (ref *PropertyNameArray) NameAtIndex( index uint16 ) string {
	jsstr := C.JSPropertyNameArrayGetNameAtIndex( C.JSPropertyNameArrayRef(unsafe.Pointer(ref)), C.size_t(index) )
	defer C.JSStringRelease( jsstr )
	return string_js_2_go( jsstr )
}

type Function interface {
	Callback( ctx *Context, obj *Object, thisObject *Object, exception **Value ) *Value
}

func (ctx *Context) MakeFunction( name string, f Function ) *Object {
	stringRef := (*String)(nil)
	if name != "" {
		stringRef = ctx.NewString( name )
		defer stringRef.Release()
	}

	ctx.callbacks = append( ctx.callbacks, f )

	tmp := C.JSObjectMakeFunctionWithCallback_wka( ctx.value, C.JSStringRef(unsafe.Pointer(stringRef)), unsafe.Pointer( &f ) )
	return (*Object)(unsafe.Pointer(tmp))
}

//export JSObjectCallAsFunctionCallback_go
func JSObjectCallAsFunctionCallback_go( callback unsafe.Pointer, ret *unsafe.Pointer ) {
//	ctx C.JSContextRef, function C.JS*Object, thisObject C.JS*Object, argumentCount uint16, JS*Value* exception ) C.JSValueRef {

	f := *(*Function)( callback )
	var exception *Value
	value := f.Callback( nil, nil, nil, &exception )
	*ret = unsafe.Pointer( value )
}

