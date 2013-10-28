package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "unsafe"
import "log"

type Object struct {
	ref C.JSObjectRef
	ctx *Context
}

func release_jsstringref_array(refs []C.JSStringRef) {
	for i := 0; i < len(refs); i++ {
		if refs[i] != nil {
			C.JSStringRelease(refs[i])
		}
	}
}

// Creates a new *Object given a C pointer to an JSObjectRef.
func (ctx *Context) newObject(ref C.JSObjectRef) *Object {
	obj := new(Object)
	obj.ref = ref
	obj.ctx = ctx
	return obj
}

func (ctx *Context) NewEmptyObject() *Object {
	obj := C.JSObjectMake(ctx.ref, nil, nil)
	return ctx.newObject(obj)
}

func (ctx *Context) NewObjectWithProperties(properties map[string]*Value) (*Object, error) {
	obj := ctx.NewEmptyObject()
	for name, val := range properties {
		err := ctx.SetProperty(obj, name, val, 0)
		if err != nil {
			return nil, err
		}
	}
	return obj, nil
}

func (ctx *Context) NewArray(items []*Value) (*Object, error) {
	errVal := ctx.newErrorValue()

	ret := ctx.NewEmptyObject()
	if items != nil {
		carr, carrlen := ctx.newCValueArray(items)
		ret.ref = C.JSObjectMakeArray(ctx.ref, carrlen, carr, &errVal.ref)
	} else {
		ret.ref = C.JSObjectMakeArray(ctx.ref, 0, nil, &errVal.ref)
	}
	if errVal.ref != nil {
		return nil, errVal
	}
	return ret, nil
}

func (ctx *Context) NewDate() (*Object, error) {
	errVal := ctx.newErrorValue()

	ret := C.JSObjectMakeDate(ctx.ref,
		0, nil,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewDateWithMilliseconds(milliseconds float64) (*Object, error) {
	errVal := ctx.newErrorValue()

	param := ctx.NewNumberValue(milliseconds)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), &param.ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewDateWithString(date string) (*Object, error) {
	errVal := ctx.newErrorValue()

	param := ctx.NewStringValue(date)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), &param.ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewRegExp(regex string) (*Object, error) {
	errVal := ctx.newErrorValue()

	param := ctx.NewStringValue(regex)

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(1), &param.ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewRegExpFromValues(parameters []*Value) (*Object, error) {
	errVal := ctx.newErrorValue()

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(len(parameters)), &parameters[0].ref,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewFunction(name string, parameters []string, body string, source_url string, starting_line_number int) (*Object, error) {
	Cname := NewString(name)
	defer Cname.Release()

	Cparameters := make([]C.JSStringRef, len(parameters))
	defer release_jsstringref_array(Cparameters)
	for i := 0; i < len(parameters); i++ {
		Cparameters[i] = (C.JSStringRef)(unsafe.Pointer(NewString(parameters[i])))
	}

	Cbody := NewString(body)
	defer Cbody.Release()

	var sourceRef *String
	if source_url != "" {
		sourceRef = NewString(source_url)
		defer sourceRef.Release()
	}

	errVal := ctx.newErrorValue()
	ret := C.JSObjectMakeFunction(ctx.ref,
		(C.JSStringRef)(unsafe.Pointer(Cname)),
		C.unsigned(len(Cparameters)), &Cparameters[0],
		(C.JSStringRef)(unsafe.Pointer(Cbody)),
		(C.JSStringRef)(unsafe.Pointer(sourceRef)),
		C.int(starting_line_number), &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) GetPrototype(obj *Object) *Value {
	ret := C.JSObjectGetPrototype(ctx.ref, obj.ref)
	return ctx.newValue(ret)
}

func (ctx *Context) SetPrototype(obj *Object, rhs *Value) {
	C.JSObjectSetPrototype(ctx.ref, obj.ref, rhs.ref)
}

func (ctx *Context) HasProperty(obj *Object, name string) bool {
	jsstr := NewString(name)
	defer jsstr.Release()

	ret := C.JSObjectHasProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)))
	return bool(ret)
}

func (ctx *Context) GetProperty(obj *Object, name string) (*Value, error) {
	jsstr := NewString(name)
	defer jsstr.Release()

	errVal := ctx.newErrorValue()

	ret := C.JSObjectGetProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	return ctx.newValue(ret), nil
}

func (ctx *Context) GetPropertyAtIndex(obj *Object, index uint16) (*Value, error) {
	errVal := ctx.newErrorValue()

	ret := C.JSObjectGetPropertyAtIndex(ctx.ref, obj.ref, C.unsigned(index), &errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	return ctx.newValue(ret), nil
}

func (ctx *Context) SetProperty(obj *Object, name string, rhs *Value, attributes uint8) error {
	jsstr := NewString(name)
	defer jsstr.Release()

	errVal := ctx.newErrorValue()

	C.JSObjectSetProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), rhs.ref,
		(C.JSPropertyAttributes)(attributes), &errVal.ref)
	if errVal.ref != nil {
		return errVal
	}

	return nil
}

func (ctx *Context) SetPropertyAtIndex(obj *Object, index uint16, rhs *Value) error {
	errVal := ctx.newErrorValue()

	C.JSObjectSetPropertyAtIndex(ctx.ref, obj.ref, C.unsigned(index), rhs.ref, &errVal.ref)
	if errVal.ref != nil {
		return errVal
	}

	return nil
}

func (ctx *Context) DeleteProperty(obj *Object, name string) (bool, error) {
	jsstr := NewString(name)
	defer jsstr.Release()

	errVal := ctx.newErrorValue()

	ret := C.JSObjectDeleteProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), &errVal.ref)
	if errVal.ref != nil {
		return false, errVal
	}

	return bool(ret), nil
}

// Should NOT be public. WTF.
func (obj *Object) GetPrivate() unsafe.Pointer {
	ret := C.JSObjectGetPrivate(obj.ref)
	return ret
}

// Should NOT be public. WTF.
func (obj *Object) SetPrivate(data unsafe.Pointer) bool {
	ret := C.JSObjectSetPrivate(obj.ref, data)
	return bool(ret)
}

// ToValue returns the JSValueRef wrapper for the object.
//
// Any JSObjectRef can be safely cast to a JSValueRef.
// https://lists.webkit.org/pipermail/webkit-dev/2009-May/007530.html
func (obj *Object) ToValue() *Value {
	if obj == nil {
		panic("ToValue() called on nil *Object!")
	}
	return obj.ctx.newValue(C.JSValueRef(obj.ref))
}

func (ctx *Context) IsFunction(obj *Object) bool {
	ret := C.JSObjectIsFunction(ctx.ref, obj.ref)
	return bool(ret)
}

func (ctx *Context) CallAsFunction(obj *Object, thisObject *Object, parameters []*Value) (*Value, error) {
	errVal := ctx.newErrorValue()
	cParameters, n := ctx.newCValueArray(parameters)
	if thisObject == nil {
		thisObject = ctx.newObject(nil)
		log.Println(thisObject.ref)
	}

	ret := C.JSObjectCallAsFunction(ctx.ref, obj.ref, thisObject.ref, n, cParameters, &errVal.ref)

	if errVal.ref != nil {
		return nil, errVal
	}

	return ctx.newValue(ret), nil
}

func (ctx *Context) IsConstructor(obj *Object) bool {
	ret := C.JSObjectIsConstructor(ctx.ref, obj.ref)
	return bool(ret)
}

func (ctx *Context) CallAsConstructor(obj *Object, parameters []*Value) (*Value, error) {
	errVal := ctx.newErrorValue()

	var Cparameters *C.JSValueRef
	if len(parameters) > 0 {
		Cparameters = (*C.JSValueRef)(unsafe.Pointer(&parameters[0]))
	}

	ret := C.JSObjectCallAsConstructor(ctx.ref, obj.ref,
		C.size_t(len(parameters)),
		Cparameters,
		&errVal.ref)
	if errVal.ref != nil {
		return nil, errVal
	}

	return ctx.newObject(ret).ToValue(), nil
}

//=========================================================
// PropertyNameArray
//

const (
	PropertyAttributeNone       = 0
	PropertyAttributeReadOnly   = 1 << 1
	PropertyAttributeDontEnum   = 1 << 2
	PropertyAttributeDontDelete = 1 << 3
)

const (
	ClassAttributeNone                 = 0
	ClassAttributeNoAutomaticPrototype = 1 << 1
)

type PropertyNameArray struct {
}

func (ctx *Context) CopyPropertyNames(obj *Object) *PropertyNameArray {
	ret := C.JSObjectCopyPropertyNames(ctx.ref, obj.ref)
	return (*PropertyNameArray)(unsafe.Pointer(ret))
}

func (ref *PropertyNameArray) Retain() {
	C.JSPropertyNameArrayRetain(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)))
}

func (ref *PropertyNameArray) Release() {
	C.JSPropertyNameArrayRelease(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)))
}

func (ref *PropertyNameArray) Count() uint16 {
	ret := C.JSPropertyNameArrayGetCount(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)))
	return uint16(ret)
}

func (ref *PropertyNameArray) NameAtIndex(index uint16) string {
	jsstr := C.JSPropertyNameArrayGetNameAtIndex(C.JSPropertyNameArrayRef(unsafe.Pointer(ref)), C.size_t(index))
	defer C.JSStringRelease(jsstr)
	return (*String)(unsafe.Pointer(jsstr)).String()
}
