package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSValueRef.h>
import "C"
import "unsafe"

type Value struct {

}

const (
	TypeUndefined = 0
	TypeNull      = iota
	TypeBoolean   = iota
	TypeNumber    = iota
	TypeString    = iota
	TypeObject    = iota
)

func (ctx *Context) ValueType(v *Value) uint8 {
	ret := C.JSValueGetType(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return uint8(ret)
}

func (ctx *Context) IsUndefined(v *Value) bool {
	ret := C.JSValueIsUndefined(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return bool(ret)
}

func (ctx *Context) IsNull(v *Value) bool {
	ret := C.JSValueIsNull(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return bool(ret)
}

func (ctx *Context) IsBoolean(v *Value) bool {
	ret := C.JSValueIsBoolean(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return bool(ret)
}

func (ctx *Context) IsNumber(v *Value) bool {
	ret := C.JSValueIsNumber(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return bool(ret)
}

func (ctx *Context) IsString(v *Value) bool {
	ret := C.JSValueIsString(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return bool(ret)
}

func (ctx *Context) IsObject(v *Value) bool {
	ret := C.JSValueIsObject(ctx.ref, C.JSValueRef(unsafe.Pointer(v)))
	return bool(ret)
}

func (ctx *Context) IsEqual(a *Value, b *Value) (bool, *Value) {
	exception := C.JSValueRef(nil)

	ret := C.JSValueIsEqual(ctx.ref, C.JSValueRef(unsafe.Pointer(a)), C.JSValueRef(unsafe.Pointer(b)), &exception)
	if exception != nil {
		return false, (*Value)(unsafe.Pointer(exception))
	}

	return bool(ret), nil
}

func (ctx *Context) IsStrictEqual(a *Value, b *Value) bool {
	ret := C.JSValueIsStrictEqual(ctx.ref, C.JSValueRef(unsafe.Pointer(a)), C.JSValueRef(unsafe.Pointer(b)))
	return bool(ret)
}

func (ctx *Context) NewUndefinedValue() *Value {
	ref := C.JSValueMakeUndefined(ctx.ref)
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewNullValue() *Value {
	ref := C.JSValueMakeNull(ctx.ref)
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewBooleanValue(value bool) *Value {
	ref := C.JSValueMakeBoolean(ctx.ref, C.bool(value))
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewNumberValue(value float64) *Value {
	ref := C.JSValueMakeNumber(ctx.ref, C.double(value))
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) NewStringValue(value string) *Value {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	jsstr := C.JSStringCreateWithUTF8CString(cvalue)
	defer C.JSStringRelease(jsstr)
	ref := C.JSValueMakeString(ctx.ref, jsstr)
	return (*Value)(unsafe.Pointer(ref))
}

func (ctx *Context) ToBoolean(ref *Value) bool {
	ret := C.JSValueToBoolean(ctx.ref, C.JSValueRef(unsafe.Pointer(ref)))
	return bool(ret)
}

func (ctx *Context) ToNumber(ref *Value) (num float64, err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToNumber(ctx.ref, C.JSValueRef(unsafe.Pointer(ref)), &exception)
	if exception != nil {
		return float64(ret), (*Value)(unsafe.Pointer(exception))
	}

	// Successful conversion
	return float64(ret), nil
}

func (ctx *Context) ToNumberOrDie(ref *Value) float64 {
	ret, err := ctx.ToNumber(ref)
	if err != nil {
		panic(newPanicError(ctx, err))
	}
	return ret
}

func (ctx *Context) ToString(ref *Value) (str string, err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToStringCopy(ctx.ref, C.JSValueRef(unsafe.Pointer(ref)), &exception)
	if exception != nil {
		// An error occurred
		// ret should be null
		return "", (*Value)(unsafe.Pointer(exception))
	}
	defer C.JSStringRelease(ret)

	// Successful conversion
	tmp := (*String)(unsafe.Pointer(ret))
	return tmp.String(), nil
}

func (ctx *Context) ToStringOrDie(ref *Value) string {
	str, err := ctx.ToString(ref)
	if err != nil {
		panic(newPanicError(ctx, err))
	}
	return str
}

func (ctx *Context) ToObject(ref *Value) (obj *Object, err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToObject(ctx.ref, C.JSValueRef(unsafe.Pointer(ref)), &exception)
	if exception != nil {
		// An error occurred
		// ret should be null
		return nil, (*Value)(unsafe.Pointer(exception))
	}

	// Successful conversion
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) ToObjectOrDie(ref *Value) *Object {
	ret, err := ctx.ToObject(ref)
	if err != nil {
		panic(newPanicError(ctx, err))
	}
	return ret
}

func (ctx *Context) ProtectValue(ref *Value) {
	C.JSValueProtect(ctx.ref, C.JSValueRef(unsafe.Pointer(ref)))
}

func (ctx *Context) UnProtectValue(ref *Value) {
	C.JSValueProtect(ctx.ref, C.JSValueRef(unsafe.Pointer(ref)))
}
