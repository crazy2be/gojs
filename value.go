package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSValueRef.h>
import "C"
import (
	"unsafe"
)

type Value struct {
	ref C.JSValueRef
	ctx *Context
}

const (
	TypeUndefined = 0
	TypeNull      = iota
	TypeBoolean   = iota
	TypeNumber    = iota
	TypeString    = iota
	TypeObject    = iota
)

func (val *Value) String() string {
	str, err := val.ctx.ToString(val)
	if err != nil {
		return "Error:" + err.Error()
	}
	return str
}

func (ctx *Context) ValueType(v *Value) uint8 {
	return uint8(C.JSValueGetType(ctx.ref, v.ref))
}

func (ctx *Context) IsUndefined(v *Value) bool {
	return bool(C.JSValueIsUndefined(ctx.ref, v.ref))
}

func (ctx *Context) IsNull(v *Value) bool {
	return bool(C.JSValueIsNull(ctx.ref, v.ref))
}

func (ctx *Context) IsBoolean(v *Value) bool {
	return bool(C.JSValueIsBoolean(ctx.ref, v.ref))
}

func (ctx *Context) IsNumber(v *Value) bool {
	return bool(C.JSValueIsNumber(ctx.ref, v.ref))
}

func (ctx *Context) IsString(v *Value) bool {
	return bool(C.JSValueIsString(ctx.ref, v.ref))
}

func (ctx *Context) IsObject(v *Value) bool {
	ret := C.JSValueIsObject(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsEqual(a *Value, b *Value) (bool, error) {
	errVal := ctx.newErrorValue()
	ret := C.JSValueIsEqual(ctx.ref, a.ref, b.ref, &errVal.ref)
	if errVal.ref != nil {
		return false, errVal
	}

	return bool(ret), nil
}

func (ctx *Context) IsStrictEqual(a *Value, b *Value) bool {
	return bool(C.JSValueIsStrictEqual(ctx.ref, a.ref, b.ref))
}

func (ctx *Context) newValue(ref C.JSValueRef) *Value {
	if ref == nil {
		return nil
	}
	val := new(Value)
	val.ctx = ctx
	val.ref = ref
	return val
}

func (ctx *Context) NewUndefinedValue() *Value {
	return ctx.newValue(C.JSValueMakeUndefined(ctx.ref))
}

func (ctx *Context) NewNullValue() *Value {
	return ctx.newValue(C.JSValueMakeNull(ctx.ref))
}

func (ctx *Context) NewBooleanValue(value bool) *Value {
	return ctx.newValue(C.JSValueMakeBoolean(ctx.ref, C.bool(value)))
}

func (ctx *Context) NewNumberValue(value float64) *Value {
	return ctx.newValue(C.JSValueMakeNumber(ctx.ref, C.double(value)))
}

func (ctx *Context) NewStringValue(value string) *Value {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	jsstr := C.JSStringCreateWithUTF8CString(cvalue)
	defer C.JSStringRelease(jsstr)
	ref := C.JSValueMakeString(ctx.ref, jsstr)
	return ctx.newValue(ref)
}

// TODO: Move to Value struct
func (ctx *Context) ToBoolean(ref *Value) bool {
	return bool(C.JSValueToBoolean(ctx.ref, ref.ref))
}

func (ctx *Context) ToNumber(ref *Value) (num float64, err error) {
	errVal := ctx.newErrorValue()
	ret := C.JSValueToNumber(ctx.ref, ref.ref, &errVal.ref)
	if errVal.ref != nil {
		return float64(ret), errVal
	}

	// Successful conversion
	return float64(ret), nil
}

func (ctx *Context) ToNumberOrDie(ref *Value) float64 {
	ret, err := ctx.ToNumber(ref)
	if err != nil {
		panic(err)
	}
	return ret
}

func (ctx *Context) ToString(ref *Value) (str string, err error) {
	errVal := ctx.newErrorValue()
	ret := C.JSValueToStringCopy(ctx.ref, ref.ref, &errVal.ref)
	if errVal.ref != nil {
		return "", errVal
	}
	defer C.JSStringRelease(ret)
	return newStringFromRef(ret).String(), nil
}

func (ctx *Context) ToStringOrDie(ref *Value) string {
	str, err := ctx.ToString(ref)
	if err != nil {
		panic(err)
	}
	return str
}

func (ctx *Context) ToObject(ref *Value) (*Object, error) {
	errVal := ctx.newErrorValue()
	ret := C.JSValueToObject(ctx.ref, ref.ref, &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	// Successful conversion
	return ctx.newObject(ret), nil
}

func (ctx *Context) ToObjectOrDie(ref *Value) *Object {
	ret, err := ctx.ToObject(ref)
	if err != nil {
		panic(err)
	}
	return ret
}

func (ctx *Context) ProtectValue(ref *Value) {
	C.JSValueProtect(ctx.ref, ref.ref)
}

func (ctx *Context) UnProtectValue(ref *Value) {
	C.JSValueProtect(ctx.ref, ref.ref)
}
