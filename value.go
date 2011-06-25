package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSValueRef.h>
import "C"
import "unsafe"
import "fmt"

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

func (ctx *Context) ValueType(v *Value) uint8 {
	fmt.Println(ctx, v)
	fmt.Printf("%#v %#v\n", ctx, v.ctx)
	fmt.Printf("%#v %#v\n", ctx.ref, v.ctx.ref)
	fmt.Printf("%#v\n", v.ref)
	fmt.Printf("%d\n", (*v.ref))
	ret := C.JSValueGetType(v.ctx.ref, v.ref)
	return uint8(ret)
}

func (ctx *Context) IsUndefined(v *Value) bool {
	ret := C.JSValueIsUndefined(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsNull(v *Value) bool {
	ret := C.JSValueIsNull(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsBoolean(v *Value) bool {
	ret := C.JSValueIsBoolean(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsNumber(v *Value) bool {
	ret := C.JSValueIsNumber(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsString(v *Value) bool {
	ret := C.JSValueIsString(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsObject(v *Value) bool {
	ret := C.JSValueIsObject(ctx.ref, v.ref)
	return bool(ret)
}

func (ctx *Context) IsEqual(a *Value, b *Value) (bool, *Value) {
	exception := C.JSValueRef(nil)

	ret := C.JSValueIsEqual(ctx.ref, a.ref, b.ref, &exception)
	if exception != nil {
		return false, ctx.newValue(exception)
	}

	return bool(ret), nil
}

func (ctx *Context) IsStrictEqual(a *Value, b *Value) bool {
	ret := C.JSValueIsStrictEqual(ctx.ref, a.ref, b.ref)
	return bool(ret)
}

func (ctx *Context) newValue(ref C.JSValueRef) *Value {
	val := new(Value)
	val.ctx = ctx
	val.ref = ref
	return val
}

func (ctx *Context) NewUndefinedValue() *Value {
	ref := C.JSValueMakeUndefined(ctx.ref)
	return ctx.newValue(ref)
}

func (ctx *Context) NewNullValue() *Value {
	ref := C.JSValueMakeNull(ctx.ref)
	return ctx.newValue(ref)
}

func (ctx *Context) NewBooleanValue(value bool) *Value {
	ref := C.JSValueMakeBoolean(ctx.ref, C.bool(value))
	return ctx.newValue(ref)
}

func (ctx *Context) NewNumberValue(value float64) *Value {
	ref := C.JSValueMakeNumber(ctx.ref, C.double(value))
	return ctx.newValue(ref)
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
	ret := C.JSValueToBoolean(ctx.ref, ref.ref)
	return bool(ret)
}

func (ctx *Context) ToNumber(ref *Value) (num float64, err *Value) {
	var exception C.JSValueRef
	ret := C.JSValueToNumber(ctx.ref, ref.ref, &exception)
	if exception != nil {
		return float64(ret), ctx.newValue(exception)
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
	ret := C.JSValueToStringCopy(ctx.ref, ref.ref, &exception)
	if exception != nil {
		// An error occurred
		// ret should be null
		return "", ctx.newValue(exception)
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
	ret := C.JSValueToObject(ctx.ref, ref.ref, &exception)
	if exception != nil {
		// An error occurred
		// ret should be null
		return nil, ctx.newValue(exception)
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
	C.JSValueProtect(ctx.ref, ref.ref)
}

func (ctx *Context) UnProtectValue(ref *Value) {
	C.JSValueProtect(ctx.ref, ref.ref)
}
