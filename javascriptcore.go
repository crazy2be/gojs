package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSContextRef.h>
// #include <JavaScriptCore/JSStringRef.h>
import "C"
import "os"
import "unsafe"

//=========================================================
// StringRef
//

type StringRef struct {
	value C.JSStringRef
}

func (ref *StringRef) Retain() {
	C.JSStringRetain( ref.value )
}

func (ref *StringRef) Release() {
	C.JSStringRelease( ref.value )
}

func string_js_2_go( ref C.JSStringRef ) string {
	// Conversion 1, null-terminate UTF-8 string
	len := C.JSStringGetMaximumUTF8CStringSize( ref )
	buffer := C.malloc( len )
	if buffer==nil {
		panic( os.ENOMEM )
	}
	defer C.free( buffer )
	C.JSStringGetUTF8CString( ref, (*C.char)(buffer), len )

	// Conversion 2, Go string
	ret := C.GoString( (*C.char)(buffer) )
	return ret
}

func (ref *StringRef) String() string {
	return string_js_2_go( ref.value )
}

func (ref *StringRef) Length() uint32 {
	ret := C.JSStringGetLength( ref.value )
	return uint32( ret )
}

func (ref *StringRef) Equal( rhs *StringRef ) bool {
	ret := C.JSStringIsEqual( ref.value, rhs.value )
	return bool( ret )
}

func (ref *StringRef) EqualToString( rhs string ) bool {
	crhs := C.CString( rhs )
	defer C.free( unsafe.Pointer(crhs) )
	ret := C.JSStringIsEqualToUTF8CString( ref.value, crhs )
	return bool( ret )
}

//=========================================================
// ContextRef
//

type Context struct {
	value C.JSGlobalContextRef
}

type ObjectRef struct {
	value C.JSObjectRef
}

type ValueRef struct {
	value C.JSValueRef
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
	return &Context{ ctx }
}

func (ctx *Context) Retain() {
	C.JSGlobalContextRetain( ctx.value )
}

func (ctx *Context) Release() {
	C.JSGlobalContextRelease( ctx.value )
}

func (ctx *Context) GlobalObject() *ObjectRef {
	ret := C.JSContextGetGlobalObject( ctx.value )
	return &ObjectRef{ ret }
}

func (ctx *Context) EvaluateScript( script string, object_ref *ObjectRef, source_url string, startingLineNumber int ) (*ValueRef, *ValueRef) {
	scriptRef := ctx.NewString( script )
	defer scriptRef.Release()

	var obj C.JSObjectRef
	if object_ref != nil {
		obj = object_ref.value
	}

	var sourceRef *StringRef
	if source_url != "" {
		sourceRef = ctx.NewString( source_url )
		defer sourceRef.Release()
	} else {
		sourceRef = &StringRef{ nil }
	}

	var exception C.JSValueRef

	ret := C.JSEvaluateScript( ctx.value, scriptRef.value, obj, sourceRef.value,
		C.int(startingLineNumber), &exception )
	if ret == nil {
		// An error occurred
		// Error information should be stored in exception
		return nil, &ValueRef{ exception }
	}

	// Successful evaluation
	return &ValueRef{ ret }, nil
}

func (ctx *Context) CheckScriptSyntax( script string, source_url string, startingLineNumber int ) *ValueRef {
	scriptRef := ctx.NewString( script )
	defer scriptRef.Release()

	var sourceRef *StringRef
	if source_url != "" {
		sourceRef = ctx.NewString( source_url )
		defer sourceRef.Release()
	} else {
		sourceRef = &StringRef{ nil }
	}

	var exception C.JSValueRef

	ret := C.JSCheckScriptSyntax( ctx.value, scriptRef.value, sourceRef.value, 
		C.int(startingLineNumber), &exception )
	if !ret {
		// A syntax error was found
		// exception should be non-nil
		return &ValueRef{ exception }
	}

	// exception should be nil
	return nil
}

func (ctx *Context) GarbageCollect() {
	C.JSGarbageCollect( ctx.value )
}

//=========================================================
// ValueRef
//

func (ctx *Context) ValueIsUndefined( v *ValueRef ) bool {
	ret := C.JSValueIsUndefined( ctx.value, v.value )
	return bool( ret )
}

func (ctx *Context) ValueIsNull( v *ValueRef ) bool {
	ret := C.JSValueIsNull( ctx.value, v.value )
	return bool( ret )
}

func (ctx *Context) ValueIsBoolean( v *ValueRef ) bool {
	ret := C.JSValueIsBoolean( ctx.value, v.value )
	return bool( ret )
}

func (ctx *Context) ValueIsString( v *ValueRef ) bool {
	ret := C.JSValueIsString( ctx.value, v.value )
	return bool( ret )
}

func (ctx *Context) ValueIsObject( v *ValueRef ) bool {
	ret := C.JSValueIsObject( ctx.value, v.value )
	return bool( ret )
}

func (ctx *Context) ValueType( v *ValueRef ) uint8 {
	ret := C.JSValueGetType( ctx.value, v.value )
	return uint8( ret )
}

func (ctx *Context) IsEqual( a *ValueRef, b *ValueRef ) (bool, *ValueRef) {
	var exception C.JSValueRef

	ret := C.JSValueIsEqual( ctx.value, a.value, b.value, &exception )
	if exception != nil {
		return false, &ValueRef{ exception }
	}

	return bool(ret), nil
}

func (ctx *Context) IsStrictEqual( a *ValueRef, b *ValueRef ) bool {
	ret := C.JSValueIsStrictEqual( ctx.value, a.value, b.value )
	return bool(ret)
}

func (ctx *Context) NewUndefinedValue() *ValueRef {
	ref := C.JSValueMakeUndefined( ctx.value )
	return &ValueRef{ ref }
}

func (ctx *Context) NewNullValue() *ValueRef {
	ref := C.JSValueMakeNull( ctx.value )
	return &ValueRef{ ref }
}

func (ctx *Context) NewBooleanValue( value bool ) *ValueRef {
	ref := C.JSValueMakeBoolean( ctx.value, C.bool(value) )
	return &ValueRef{ ref }
}

func (ctx *Context) NewNumberValue( value float64 ) *ValueRef {
	ref := C.JSValueMakeNumber( ctx.value, C.double(value) )
	return &ValueRef{ ref }
}

func (ctx *Context) NewStringValue( value string ) *ValueRef {
	cvalue := C.CString( value )
	defer C.free( unsafe.Pointer(cvalue) )
	jsstr := C.JSStringCreateWithUTF8CString( cvalue )
	defer C.JSStringRelease( jsstr )
	ref := C.JSValueMakeString( ctx.value, jsstr )
	return &ValueRef{ ref }
}

func (ctx *Context) NewString( value string ) *StringRef {
	cvalue := C.CString( value )
	defer C.free( unsafe.Pointer(cvalue) )
	ref := C.JSStringCreateWithUTF8CString( cvalue )
	return &StringRef{ ref }
}

func (ctx *Context) ProtectValue( ref *ValueRef ) {
	C.JSValueProtect( ctx.value, ref.value )
}

func (ctx *Context) UnProtectValue( ref *ValueRef ) {
	C.JSValueProtect( ctx.value, ref.value )
}

func (ctx *Context) ToBoolean( ref *ValueRef ) bool {
	ret := C.JSValueToBoolean( ctx.value, ref.value )
	return bool( ret )
}

func (ctx *Context) ToNumber( ref *ValueRef ) (num float64,err *ValueRef) {
	var exception C.JSValueRef
	ret := C.JSValueToNumber( ctx.value, ref.value, &exception )
	if exception != nil {
		return float64(ret), &ValueRef{ exception }
	}

	// Successful conversion
	return float64(ret), nil
}

func (ctx *Context) ToString( ref *ValueRef ) (str string, err *ValueRef) {
	var exception C.JSValueRef
	ret := C.JSValueToStringCopy( ctx.value, ref.value, &exception )
	if exception != nil {
		// An error occurred
		// ret should be null
		return "", &ValueRef{exception}
	}

	// Successful conversion
	tmp := &StringRef{ ret }
	return tmp.String(), nil
}

func (ctx *Context) ToStringOrDie( ref *ValueRef ) string {
	str, err := ctx.ToString( ref )
	if err!=nil {
		panic( err )
	}
	return str
}

//=========================================================
// ObjectRef
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

func (ctx *Context) ObjectHasProperty(obj *ObjectRef, name string) bool {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	ret := C.JSObjectHasProperty( ctx.value, obj.value, jsstr.value )
	return bool(ret)
}

func (ctx *Context) ObjectGetProperty(obj *ObjectRef, name string) (*ValueRef,*ValueRef) {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	ret := C.JSObjectGetProperty( ctx.value, obj.value, jsstr.value, &exception )
	if exception != nil {
		return nil, &ValueRef{ exception }
	}

	return &ValueRef{ ret }, nil
}

func (ctx *Context) ObjectSetProperty(obj *ObjectRef, name string, rhs *ValueRef, attributes uint8) *ValueRef {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	C.JSObjectSetProperty( ctx.value, obj.value, jsstr.value, rhs.value, 
		(C.JSPropertyAttributes)(attributes), &exception )
	if exception != nil {
		return &ValueRef{ exception }
	}

	return nil
}

func (ctx *Context) ObjectDeleteProperty(obj *ObjectRef, name string ) (bool,*ValueRef) {
	jsstr := ctx.NewString( name )
	defer jsstr.Release()

	var exception C.JSValueRef

	ret := C.JSObjectDeleteProperty( ctx.value, obj.value, jsstr.value, &exception )
	if exception != nil {
		return false, &ValueRef{ exception }
	}

	return bool(ret), nil
}

func (ctx *Context) IsFunction(obj *ObjectRef) bool {
	ret := C.JSObjectIsFunction( ctx.value, obj.value )
	return bool(ret)
}

func (ctx *Context) IsConstructor(obj *ObjectRef) bool {
	ret := C.JSObjectIsConstructor( ctx.value, obj.value )
	return bool(ret)
}

