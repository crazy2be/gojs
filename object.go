package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "unsafe"

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

func (ctx *Context) NewObject(ref C.JSObjectRef) *Object {
	obj := new(Object)
	obj.ref = ref
	obj.ctx = ctx
	return obj
	//ret := 
	//return (*Object)(unsafe.Pointer(ret))
}

func (ctx *Context) NewEmptyObject() *Object {
	obj := C.JSObjectMake(ctx.ref, nil, nil)
	return ctx.NewObject(obj)
}

func (ctx *Context) NewArray(items []*Value) (*Object, *Exception) {
	var exception = ctx.NewException()

	ret := ctx.NewEmptyObject()
	if items != nil {
		ret.ref = C.JSObjectMakeArray(ctx.ref,
			C.size_t(len(items)), (*C.JSValueRef)(unsafe.Pointer(&items[0])),
			&exception.val.ref)
	} else {
		ret.ref = C.JSObjectMakeArray(ctx.ref,
			0, nil,
			&exception.val.ref)
	}
	if exception != nil {
		return nil, exception
	}
	return ret, nil
}

func (ctx *Context) NewDate() (*Object, *Exception) {
	var exception = ctx.NewException()

	ret := C.JSObjectMakeDate(ctx.ref,
		0, nil,
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return ctx.NewObject(ret), nil
}

func (ctx *Context) NewDateWithMilliseconds(milliseconds float64) (*Object, *Exception) {
	var exception = ctx.NewException()

	param := ctx.NewNumberValue(milliseconds)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), &param.ref,
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return ctx.NewObject(ret), nil
}

func (ctx *Context) NewDateWithString(date string) (*Object, *Exception) {
	var exception = ctx.NewException()

	param := ctx.NewStringValue(date)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), (*C.JSValueRef)(unsafe.Pointer(&param)),
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return ctx.NewObject(ret), nil
}

// Used for reporting errors in go code to javascript.
func (ctx *Context) NewError(message string) (*Object, *Exception) {
	exception := ctx.NewException()

	param := ctx.NewStringValue(message)

	ret := C.JSObjectMakeError(ctx.ref,
		C.size_t(1), (*C.JSValueRef)(unsafe.Pointer(&param)),
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) NewRegExp(regex string) (*Object, *Exception) {
	exception := ctx.NewException()

	param := ctx.NewStringValue(regex)

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(1), (*C.JSValueRef)(unsafe.Pointer(&param)),
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) NewRegExpFromValues(parameters []*Value) (*Object, *Exception) {
	exception := ctx.NewException()

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(len(parameters)), (*C.JSValueRef)(unsafe.Pointer(&parameters[0])),
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) NewFunction(name string, parameters []string, body string, source_url string, starting_line_number int) (*Object, *Exception) {
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

	exception := ctx.NewException()

	ret := C.JSObjectMakeFunction(ctx.ref,
		(C.JSStringRef)(unsafe.Pointer(Cname)),
		C.unsigned(len(Cparameters)), &Cparameters[0],
		(C.JSStringRef)(unsafe.Pointer(Cbody)),
		(C.JSStringRef)(unsafe.Pointer(sourceRef)),
		C.int(starting_line_number), &exception.val.ref)
	if exception != nil {
		return nil, exception
	}
	return (*Object)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) GetPrototype(obj *Object) *Value {
	ret := C.JSObjectGetPrototype(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)))
	return (*Value)(unsafe.Pointer(ret))
}

func (ctx *Context) SetPrototype(obj *Object, rhs *Value) {
	C.JSObjectSetPrototype(ctx.ref,
		C.JSObjectRef(unsafe.Pointer(obj)), C.JSValueRef(unsafe.Pointer(rhs)))
}

func (ctx *Context) HasProperty(obj *Object, name string) bool {
	jsstr := NewString(name)
	defer jsstr.Release()

	ret := C.JSObjectHasProperty(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)))
	return bool(ret)
}

func (ctx *Context) GetProperty(obj *Object, name string) (*Value, *Exception) {
	jsstr := NewString(name)
	defer jsstr.Release()

	exception := ctx.NewException()

	ret := C.JSObjectGetProperty(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), &exception.val.ref)
	if exception != nil {
		return nil, exception
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) GetPropertyAtIndex(obj *Object, index uint16) (*Value, *Exception) {
	exception := ctx.NewException()

	ret := C.JSObjectGetPropertyAtIndex(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)), C.unsigned(index), &exception.val.ref)
	if exception != nil {
		return nil, exception
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) SetProperty(obj *Object, name string, rhs *Value, attributes uint8) *Exception {
	jsstr := NewString(name)
	defer jsstr.Release()

	exception := ctx.NewException()

	C.JSObjectSetProperty(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), C.JSValueRef(unsafe.Pointer(rhs)),
		(C.JSPropertyAttributes)(attributes), &exception.val.ref)
	if exception != nil {
		return exception
	}

	return nil
}

func (ctx *Context) SetPropertyAtIndex(obj *Object, index uint16, rhs *Value) *Exception {
	exception := ctx.NewException()

	C.JSObjectSetPropertyAtIndex(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)), C.unsigned(index),
		C.JSValueRef(unsafe.Pointer(rhs)), &exception.val.ref)
	if exception != nil {
		return exception
	}

	return nil
}

func (ctx *Context) DeleteProperty(obj *Object, name string) (bool, *Exception) {
	jsstr := NewString(name)
	defer jsstr.Release()

	exception := ctx.NewException()

	ret := C.JSObjectDeleteProperty(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)), C.JSStringRef(unsafe.Pointer(jsstr)), &exception.val.ref)
	if exception != nil {
		return false, exception
	}

	return bool(ret), nil
}

func (obj *Object) GetPrivate() unsafe.Pointer {
	ret := C.JSObjectGetPrivate(C.JSObjectRef(unsafe.Pointer(obj)))
	return ret
}

func (obj *Object) SetPrivate(data unsafe.Pointer) bool {
	ret := C.JSObjectSetPrivate(C.JSObjectRef(unsafe.Pointer(obj)), data)
	return bool(ret)
}

func (obj *Object) ToValue() *Value {
	return (*Value)(unsafe.Pointer(obj))
}

func (ctx *Context) IsFunction(obj *Object) bool {
	ret := C.JSObjectIsFunction(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)))
	return bool(ret)
}

func (ctx *Context) CallAsFunction(obj *Object, thisObject *Object, parameters []*Value) (*Value, *Exception) {
	exception := ctx.NewException()

	var Cparameters *C.JSValueRef
	if len(parameters) > 0 {
		Cparameters = (*C.JSValueRef)(unsafe.Pointer(&parameters[0]))
	}

	ret := C.JSObjectCallAsFunction(ctx.ref,
		C.JSObjectRef(unsafe.Pointer(obj)),
		C.JSObjectRef(unsafe.Pointer(thisObject)),
		C.size_t(len(parameters)),
		Cparameters,
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}

	return (*Value)(unsafe.Pointer(ret)), nil
}

func (ctx *Context) IsConstructor(obj *Object) bool {
	ret := C.JSObjectIsConstructor(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)))
	return bool(ret)
}

func (ctx *Context) CallAsConstructor(obj *Object, parameters []*Value) (*Value, *Exception) {
	exception := ctx.NewException()

	var Cparameters *C.JSValueRef
	if len(parameters) > 0 {
		Cparameters = (*C.JSValueRef)(unsafe.Pointer(&parameters[0]))
	}

	ret := C.JSObjectCallAsConstructor(ctx.ref,
		C.JSObjectRef(unsafe.Pointer(obj)),
		C.size_t(len(parameters)),
		Cparameters,
		&exception.val.ref)
	if exception != nil {
		return nil, exception
	}

	return (*Value)(unsafe.Pointer(ret)), nil
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
	ret := C.JSObjectCopyPropertyNames(ctx.ref, C.JSObjectRef(unsafe.Pointer(obj)))
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
