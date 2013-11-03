package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSValueRef.h>
import "C"
import (
	"encoding/json"
	"fmt"
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

func (ctx *Context) newValue(ref C.JSValueRef) *Value {
	if ref == nil {
		return nil
	}
	val := new(Value)
	val.ctx = ctx
	val.ref = ref
	return val
}

type RawValue C.JSValueRef

func (ctx *Context) NewValueFrom(raw RawValue) *Value {
	return ctx.newValue(C.JSValueRef(raw))
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

func (val *Value) String() string {
	str, err := val.ToString()
	if err != nil {
		return "Error:" + err.Error()
	}
	return str
}

// GoVal converts a JavaScript value to a Go value. TODO(sqs): might it be
// easier to just have JavaScriptCore serialize this to JSON and then
// deserialize it in Go?
func (v *Value) GoValue() (goval interface{}, err error) {
	switch v.Type() {
	case TypeUndefined, TypeNull:
		return nil, nil
	case TypeBoolean:
		return v.ToBoolean(), nil
	case TypeNumber:
		return v.ToNumber()
	case TypeString:
		return v.ToString()
	case TypeObject:
		jsonData, err := v.ToJSON()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonData, &goval)
		return goval, err
	}
	return nil, fmt.Errorf("JS value type %d is not convertible to a Go value", v.Type())
}

func (v *Value) Type() uint8 {
	return uint8(C.JSValueGetType(v.ctx.ref, v.ref))
}

func (v *Value) IsUndefined() bool {
	return bool(C.JSValueIsUndefined(v.ctx.ref, v.ref))
}

func (v *Value) IsNull() bool {
	return bool(C.JSValueIsNull(v.ctx.ref, v.ref))
}

func (v *Value) IsBoolean() bool {
	return bool(C.JSValueIsBoolean(v.ctx.ref, v.ref))
}

func (v *Value) IsNumber() bool {
	return bool(C.JSValueIsNumber(v.ctx.ref, v.ref))
}

func (v *Value) IsString() bool {
	return bool(C.JSValueIsString(v.ctx.ref, v.ref))
}

func (v *Value) IsObject() bool {
	return bool(C.JSValueIsObject(v.ctx.ref, v.ref))
}

func (v *Value) Equals(b *Value) bool {
	return bool(C.JSValueIsStrictEqual(v.ctx.ref, v.ref, b.ref))
}

// JavaScript ==
func (v *Value) LooseEquals(b *Value) (bool, error) {
	errVal := v.ctx.newErrorValue()
	ret := C.JSValueIsEqual(v.ctx.ref, v.ref, b.ref, &errVal.ref)
	if errVal.ref != nil {
		return false, errVal
	}

	return bool(ret), nil
}

func (v *Value) ToBoolean() bool {
	return bool(C.JSValueToBoolean(v.ctx.ref, v.ref))
}

func (v *Value) ToNumber() (num float64, err error) {
	errVal := v.ctx.newErrorValue()
	ret := C.JSValueToNumber(v.ctx.ref, v.ref, &errVal.ref)
	if errVal.ref != nil {
		return float64(ret), errVal
	}

	// Successful conversion
	return float64(ret), nil
}

// TODO(crazy2be): Should this return NaN instead of panicing?
func (v *Value) ToNumberOrDie() float64 {
	ret, err := v.ToNumber()
	if err != nil {
		panic(err)
	}
	return ret
}

func (v *Value) ToString() (str string, err error) {
	errVal := v.ctx.newErrorValue()
	ret := C.JSValueToStringCopy(v.ctx.ref, v.ref, &errVal.ref)
	if errVal.ref != nil {
		return "", errVal
	}
	defer C.JSStringRelease(ret)
	return newStringFromRef(ret).String(), nil
}

func (v *Value) ToStringOrDie() string {
	str, err := v.ToString()
	if err != nil {
		panic(err)
	}
	return str
}

func (v *Value) ToObject() (*Object, error) {
	errVal := v.ctx.newErrorValue()
	ret := C.JSValueToObject(v.ctx.ref, v.ref, &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return v.ctx.newObject(ret), nil
}

func (v *Value) ToObjectOrDie() *Object {
	ret, err := v.ToObject()
	if err != nil {
		panic(err)
	}
	return ret
}

// JSON returns the JSON representation of the JavaScript value.
func (v *Value) ToJSON() ([]byte, error) {
	errVal := v.ctx.newErrorValue()
	jsstr := C.JSValueCreateJSONString(v.ctx.ref, v.ref, 0, &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	defer C.JSStringRelease(jsstr)
	return (*String)(unsafe.Pointer(jsstr)).Bytes(), nil
}

func (v *Value) Protect() {
	C.JSValueProtect(v.ctx.ref, v.ref)
}

func (v *Value) UnProtect() {
	C.JSValueProtect(v.ctx.ref, v.ref)
}
