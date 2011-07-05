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

// DEPRECATED! Use ctx.newObject() instead, this should be private!
func (ctx *Context) NewObject(ref C.JSObjectRef) *Object {
	log.Println("Warning: Use of depricated method NewObject!")
	return ctx.newObject(ref)
	//ret := 
	//return (*Object)(unsafe.Pointer(ret))
}

func (ctx *Context) NewEmptyObject() *Object {
	obj := C.JSObjectMake(ctx.ref, nil, nil)
	return ctx.newObject(obj)
}

func (ctx *Context) NewArray(items []*Value) (*Object, *Exception) {
	var exception = ctx.NewException()

	ret := ctx.NewEmptyObject()
	log.Println(exception)
	if items != nil {
		carr, carrlen := ctx.newCValueArray(items)
		ret.ref = C.JSObjectMakeArray(ctx.ref, carrlen, carr, &exception.val.ref)
	} else {
		ret.ref = C.JSObjectMakeArray(ctx.ref, 0, nil,  &exception.val.ref)
	}
	if exception.val != nil {
		return nil, exception
	}
	return ret, nil
}

func (ctx *Context) NewDate() (*Object, *Exception) {
	var exception = ctx.NewException()

	ret := C.JSObjectMakeDate(ctx.ref,
		0, nil,
		&exception.val.ref)
	if exception.val != nil {
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
	if exception.val != nil {
		return nil, exception
	}
	return ctx.NewObject(ret), nil
}

func (ctx *Context) NewDateWithString(date string) (*Object, *Exception) {
	var exception = ctx.NewException()

	param := ctx.NewStringValue(date)

	ret := C.JSObjectMakeDate(ctx.ref,
		C.size_t(1), &param.ref,
		&exception.val.ref)
	if exception.val != nil {
		return nil, exception
	}
	return ctx.NewObject(ret), nil
}

func (ctx *Context) NewRegExp(regex string) (*Object, *Exception) {
	exception := ctx.NewException()

	param := ctx.NewStringValue(regex)

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(1), &param.ref,
		&exception.val.ref)
	if exception.val != nil {
		return nil, exception
	}
	return ctx.newObject(ret), nil
}

func (ctx *Context) NewRegExpFromValues(parameters []*Value) (*Object, *Exception) {
	exception := ctx.NewException()

	ret := C.JSObjectMakeRegExp(ctx.ref,
		C.size_t(len(parameters)), &parameters[0].ref,
		&exception.val.ref)
	if exception.val != nil {
		return nil, exception
	}
	return ctx.newObject(ret), nil
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
	if exception.val != nil {
		return nil, exception
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

func (ctx *Context) GetProperty(obj *Object, name string) (*Value, *Exception) {
	jsstr := NewString(name)
	defer jsstr.Release()

	exception := ctx.NewException()

	ret := C.JSObjectGetProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), &exception.val.ref)
	if exception.val != nil {
		return nil, exception
	}

	return ctx.newValue(ret), nil
}

func (ctx *Context) GetPropertyAtIndex(obj *Object, index uint16) (*Value, *Exception) {
	exception := ctx.NewException()

	ret := C.JSObjectGetPropertyAtIndex(ctx.ref, obj.ref, C.unsigned(index), &exception.val.ref)
	if exception != nil {
		return nil, exception
	}

	return ctx.newValue(ret), nil
}

func (ctx *Context) SetProperty(obj *Object, name string, rhs *Value, attributes uint8) *Exception {
	jsstr := NewString(name)
	defer jsstr.Release()

	exception := ctx.NewException()

	C.JSObjectSetProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), rhs.ref,
		(C.JSPropertyAttributes)(attributes), &exception.val.ref)
	if exception.val != nil {
		return exception
	}

	return nil
}

func (ctx *Context) SetPropertyAtIndex(obj *Object, index uint16, rhs *Value) *Exception {
	exception := ctx.NewException()

	C.JSObjectSetPropertyAtIndex(ctx.ref, obj.ref, C.unsigned(index), rhs.ref, &exception.val.ref)
	if exception != nil {
		return exception
	}

	return nil
}

func (ctx *Context) DeleteProperty(obj *Object, name string) (bool, *Exception) {
	jsstr := NewString(name)
	defer jsstr.Release()

	exception := ctx.NewException()

	ret := C.JSObjectDeleteProperty(ctx.ref, obj.ref, C.JSStringRef(unsafe.Pointer(jsstr)), &exception.val.ref)
	if exception != nil {
		return false, exception
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

// Does this work properly?
func (obj *Object) ToValue() *Value {
	log.Println(obj)
	log.Println("In ToValue() function...")
	if obj == nil {
		panic("ToValue() called on nil *Object!")
	}
	val := obj.ctx.newValue(C.JSValueRef(obj.ref))
	log.Println("Converted to value!", val)
	return val
}

func (ctx *Context) IsFunction(obj *Object) bool {
	ret := C.JSObjectIsFunction(ctx.ref, obj.ref)
	return bool(ret)
}

func (ctx *Context) CallAsFunction(obj *Object, thisObject *Object, parameters []*Value) (*Value, *Exception) {
	exception := ctx.NewException()

	Cparameters, n := ctx.newCValueArray(parameters)
	
	if thisObject == nil {
		thisObject = ctx.newObject(nil)
		log.Println(thisObject.ref)
	}
	if len(parameters) > 0 {
		str1, err := ctx.ToString(parameters[0])
		log.Println(str1, err)
		str2, err := ctx.ToString(parameters[1])
		log.Println(str2, err)
		str, err :=  ctx.ToString(ctx.newValue(*Cparameters))
		log.Println(str, err)
// 		val2 := ctx.newValue(C.JSValueRef(int(uintptr(unsafe.Pointer(Cparameters)) + 0x1)))
// 		str2, err := ctx.ToString(val2)
// 		log.Println(str2, err)
	}
	log.Println("In CallAsFunction, about to enter C mode...")
	log.Println(obj, thisObject, parameters, Cparameters, n, exception)
	
	ret := C.JSObjectCallAsFunction(ctx.ref, obj.ref, thisObject.ref, n, Cparameters, &exception.val.ref)
	
	log.Println("Successfully exited C mode...")
	log.Println(ret)
	
	if exception.val != nil {
		return nil, exception
	}

	return ctx.newValue(ret), nil
}

func (ctx *Context) IsConstructor(obj *Object) bool {
	ret := C.JSObjectIsConstructor(ctx.ref, obj.ref)
	return bool(ret)
}

func (ctx *Context) CallAsConstructor(obj *Object, parameters []*Value) (*Value, *Exception) {
	exception := ctx.NewException()

	var Cparameters *C.JSValueRef
	if len(parameters) > 0 {
		Cparameters = (*C.JSValueRef)(unsafe.Pointer(&parameters[0]))
	}

	ret := C.JSObjectCallAsConstructor(ctx.ref, obj.ref,
		C.size_t(len(parameters)),
		Cparameters,
		&exception.val.ref)
	if exception != nil {
		return nil, exception
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
