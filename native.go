package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"syscall"
	"unsafe"
)

type object_data struct {
	typ    reflect.Type
	val    reflect.Value
	method int
}

var (
	nativecallback C.JSClassRef
	nativefunction C.JSClassRef
	nativeobject   C.JSClassRef
	nativemethod   C.JSClassRef
	objects        map[uintptr]*object_data
)

type Stringer interface {
	String() string
}

func init() {
	// Create the class definition for JavaScriptCore
	nativecallback = C.JSClassDefinition_NativeCallback()
	if nativecallback == nil {
		panic(syscall.ENOMEM)
	}

	// Create the class definition for JavaScriptCore
	nativeobject = C.JSClassDefinition_NativeObject()
	if nativeobject == nil {
		panic(syscall.ENOMEM)
	}

	// Create the class definition for JavaScriptCore
	nativefunction = C.JSClassDefinition_NativeFunction()
	if nativefunction == nil {
		panic(syscall.ENOMEM)
	}

	// Create the class definition for JavaScriptCore
	nativemethod = C.JSClassDefinition_NativeMethod()
	if nativemethod == nil {
		panic(syscall.ENOMEM)
	}

	// Create map for native objects
	objects = make(map[uintptr]*object_data)
}

// Given a slice of go-style Values, this function allocates a new array of c-style values and returns a pointer to the first element in the array, along with the length of the array.
func (ctx *Context) newCValueArray(val []*Value) (*C.JSValueRef, C.size_t) {
	if len(val) == 0 {
		return nil, 0
	}
	arr := make([]C.JSValueRef, len(val))
	for i := 0; i < len(val); i++ {
		arr[i] = val[i].ref
	}
	return &arr[0], C.size_t(len(arr))
}

func (ctx *Context) ptrValue(ptr unsafe.Pointer) *Value {
	return ctx.newValue(*(*C.JSValueRef)(ptr))
}

func (ctx *Context) newGoValueArray(ptr unsafe.Pointer, size uint) []*Value {
	// TODO(sqs): use technique from https://code.google.com/p/go-wiki/wiki/cgo
	if uintptr(ptr) == 0x00000000 {
		return nil
	}
	ptrs := unsafe.Sizeof(uintptr(0))
	goarr := make([]*Value, size)
	for i := uint(0); i < size; i++ {
		goarr[i] = ctx.ptrValue(ptr)
		// Increment the pointer by one space
		ptr = unsafe.Pointer(uintptr(ptr) + ptrs)
	}
	return goarr
}

// Given a reflect.Value, this function examines the type and returns a javascript value that best represents the given value. If no acceptable conversion can be found, it panics.
func (ctx *Context) reflectToJSValue(value reflect.Value) *Value {
	// Allows functions to return JavaScriptCore values and objects
	// directly.  These we can return without conversion.
	if value.Type() == reflect.TypeOf((*Value)(nil)) {
		// Type is already a JavaScriptCore value
		return value.Interface().(*Value)
	}
	if value.Type() == reflect.TypeOf((*Object)(nil)) {
		// Type is already a JavaScriptCore object
		// nearly there
		return value.Interface().(*Object).ToValue()
	}

	// Handle simple types directly.  These can be identified by their
	// types in the package 'reflect'.
	switch value.Kind() {
	case (reflect.Int):
		r := value.Int()
		return ctx.NewNumberValue(float64(r))
	case (reflect.Uint):
		r := value.Uint()
		return ctx.NewNumberValue(float64(r))
	case (reflect.Float64), (reflect.Float32):
		r := value.Float()
		return ctx.NewNumberValue(r)
	case (reflect.String):
		r := value.String()
		return ctx.NewStringValue(r)
	case (reflect.Func):
		r := value.Interface()
		return ctx.NewFunctionWithNative(r).ToValue()
	//case (reflect.Struct):
	//	r := value.Interface()
	//	return ctx.NewNativeObject(r).ToValue()
	case (reflect.Ptr):
		if value.IsNil() {
			return ctx.NewNullValue()
		}
		r := value.Elem()
		if r.Kind() == reflect.Struct {
			ret := ctx.NewNativeObject(value.Interface())
			return ret.ToValue()
		}
		if r.Kind() == reflect.Array {
			panic("Called reflectToJSValue() with a pointer to an array or slice. Most likely, this is a native javascript object you are passing by accident, and meant to pass to newValue() rather than NewValue(). If you really are trying to convert a pointer to an array into a native javascript object, dereference it first (slices are pointers internally anyway, so no significant loss in efficiency).")
			//log.Println("About to make new native object from *[0]uint8")
			//ret := ctx.NewNativeObject(value.Interface())
			//log.Println("Made new native object from *[0]uint8")
			//return ret.ToValue()
		}
	}
	// No acceptable conversion found.
	panic("Parameter can not be converted from Go native type. Type is " + value.Kind().String() + ", value is " + value.String())
}

func panicArgToJSString(ctx *Context, r interface{}) *Value {
	var msg string
	switch r := r.(type) {
	case error:
		msg = r.Error()
	case string:
		msg = r
	default:
		msg = fmt.Sprintf("unhandled Go panic: %v", r)
	}
	return ctx.NewStringValue(msg)
}

func (ctx *Context) jsValuesToReflect(param []*Value) []reflect.Value {
	ret := make([]reflect.Value, len(param))

	for index, item := range param {
		var goval interface{}
		log.Println(index, item)

		switch ctx.ValueType(item) {
		case TypeBoolean:
			goval = ctx.ToBoolean(item)
		case TypeNumber:
			goval = ctx.ToNumberOrDie(item)
		case TypeString:
			goval = ctx.ToStringOrDie(item)
		default:
			panic("Parameter can not be converted to Go native type.")
		}

		ret[index] = reflect.ValueOf(goval)
	}

	return ret
}

func setNativeFieldFromJSValue(field reflect.Value, ctx *Context, value *Value) (err error) {
	switch field.Kind() {
	case reflect.String:
		var str string
		str, err = ctx.ToString(value)
		if err == nil {
			field.SetString(str)
		} else {
			return
		}

	case reflect.Float32, reflect.Float64:
		var flt float64
		flt, err = ctx.ToNumber(value)
		if err == nil {
			field.SetFloat(flt)
		} else {
			return
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var flt float64
		flt, err = ctx.ToNumber(value)
		if err == nil {
			field.SetInt(int64(flt))
		} else {
			return
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		log.Println("Dealing with uint type of some sort...")
		var flt float64
		flt, err = ctx.ToNumber(value)
		log.Println("Got value!")
		if err != nil {
			return
		}
		if flt >= 0 {
			field.SetUint(uint64(flt))
		} else {
			err = errors.New("number must be greater than or equal to zero")
			return
		}

	default:
		panic("Parameter can not be converted to Go native type.")
	}
	return
}

//=========================================================
// Finalizer from JavaScriptCore for all native objects
//---------------------------------------------------------

func register(data *object_data) {
	id := uintptr(unsafe.Pointer(data))
	objects[id] = data
}

//export finalize_go
func finalize_go(data unsafe.Pointer) {
	// Called from JavaScriptCore finalizer methods
	id := uintptr(data)
	delete(objects, id)
}

//=========================================================
// Native Callback
//---------------------------------------------------------

type GoFunctionCallback func(ctx *Context, obj *Object, thisObject *Object, arguments []*Value) (ret *Value)

func (ctx *Context) NewFunctionWithCallback(callback GoFunctionCallback) *Object {
	// Register the native Go object
	data := &object_data{
		reflect.TypeOf(callback),
		reflect.ValueOf(callback),
		0}
	register(data)

	ret := C.JSObjectMake(ctx.ref, nativecallback, unsafe.Pointer(data))
	return ctx.newObject(ret)
}

//export nativecallback_CallAsFunction_go
func nativecallback_CallAsFunction_go(data_ptr unsafe.Pointer, rawCtx C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount uint, arguments unsafe.Pointer, exception *C.JSValueRef) unsafe.Pointer {
	ctx := NewContextFrom(RawContext(rawCtx))
	defer func() {
		if r := recover(); r != nil {
			*exception = panicArgToJSString(ctx, r).ref
		}
	}()

	data := (*object_data)(data_ptr)
	ret := data.val.Interface().(GoFunctionCallback)(
		ctx, ctx.newObject(function), ctx.newObject(thisObject), ctx.newGoValueArray(arguments, argumentCount) /*(*[1 << 14]*Value)(arguments)[0:argumentCount]*/)
	if ret == nil {
		return unsafe.Pointer(nil)
	}
	return unsafe.Pointer(ret.ref)
}

//=========================================================
// Native Function
//---------------------------------------------------------

func (ctx *Context) NewFunctionWithNative(fn interface{}) *Object {
	// Sanity checks on the function
	if typ := reflect.TypeOf(fn); typ.NumOut() > 1 {
		panic("Bad native function:  too many output parameters")
	}

	// Create Go-side registration
	data := &object_data{
		reflect.TypeOf(fn),
		reflect.ValueOf(fn),
		0}
	register(data)

	ret := C.JSObjectMake(ctx.ref, nativefunction, unsafe.Pointer(data))
	return ctx.newObject(ret)
}

func docall(ctx *Context, val reflect.Value, argumentCount uint, arguments unsafe.Pointer) *Value {
	// Step one, convert the JavaScriptCore array of arguments to
	// an array of reflect.Values.
	var in []reflect.Value
	if argumentCount != 0 {
		valarr := ctx.newGoValueArray(arguments, argumentCount)
		log.Println("Converted pointer to go-style array", valarr)
		log.Println("Converting to relfect.Value s")
		in = ctx.jsValuesToReflect(valarr)
	}

	log.Println("Converted arguments to native go reflect types. About to actually call the callback function...")

	log.Println(in)

	// Step two, perform the call
	out := val.Call(in)

	// Step three, convert the function return value back to JavaScriptCore
	if len(out) == 0 {
		return nil
	}
	// len(out) should be equal to 1
	return ctx.reflectToJSValue(out[0])
}

//export nativefunction_CallAsFunction_go
func nativefunction_CallAsFunction_go(data_ptr unsafe.Pointer, rawCtx C.JSContextRef, _ unsafe.Pointer, _ unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *C.JSValueRef) unsafe.Pointer {
	ctx := NewContextFrom(RawContext(rawCtx))
	defer func() {
		if r := recover(); r != nil {
			*exception = panicArgToJSString(ctx, r).ref
		}
	}()

	// recover the object
	data := (*object_data)(data_ptr)
	typ := data.typ
	val := data.val

	// Do the number of input parameters match?
	if typ.NumIn() != int(argumentCount) {
		panic("Incorrect number of function arguments")
	}

	log.Println("About to docall()!")

	ret := docall(ctx, val, argumentCount, arguments)
	if ret == nil {
		return nil
	}
	return unsafe.Pointer(ret.ref)
}

//=========================================================
// Native Object
//---------------------------------------------------------

func (ctx *Context) NewNativeObject(obj interface{}) *Object {
	// The obj must be a pointer to a struct
	// TODO:  add error checking code

	data := &object_data{
		reflect.TypeOf(obj),
		reflect.ValueOf(obj),
		0}
	register(data)

	ret := C.JSObjectMake(ctx.ref, nativeobject, unsafe.Pointer(data))
	return ctx.newObject(ret)
}

//export nativeobject_GetProperty_go
func nativeobject_GetProperty_go(data_ptr, uctx, _, propertyName unsafe.Pointer, exception *unsafe.Pointer) unsafe.Pointer {
	ctx := NewContextFrom(RawContext(uctx))
	// Get name of property as a go string
	name := (*String)(propertyName).String()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Drill down through reflect to find the property
	val := data.val
	if ptrvalue := val; ptrvalue.Kind() == reflect.Ptr {
		val = ptrvalue.Elem()
	}
	struct_val := val
	if struct_val.Kind() != reflect.Struct {
		return nil
	}

	// Can we locate a field with the proper name?
	field := struct_val.FieldByName(name)
	if field.IsValid() {
		return unsafe.Pointer(ctx.reflectToJSValue(field).ref)
	}

	// Can we locate a method with the proper name?
	typ := data.typ
	for lp := 0; lp < typ.NumMethod(); lp++ {
		if typ.Method(lp).Name == name {
			ret := newNativeMethod(ctx, data, lp)
			return unsafe.Pointer(ret.ref)
		}
	}

	// No matches found
	return nil
}

//export nativeobject_SetProperty_go
func nativeobject_SetProperty_go(data_ptr unsafe.Pointer, rawCtx C.JSContextRef, _, propertyName C.JSStringRef, value C.JSValueRef, exception *C.JSValueRef) C.char {
	ctx := NewContextFrom(RawContext(rawCtx))
	// Get name of property as a go string
	name := newStringFromRef(propertyName).String()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Drill down through reflect to find the property
	val := data.val
	if ptrvalue := val; ptrvalue.Kind() == reflect.Ptr {
		val = ptrvalue.Elem()
	}
	struct_val := val
	if struct_val.Kind() != reflect.Struct {
		*exception = ctx.newErrorOrPanic("object is not a Go struct")
		return 0
	}

	field := struct_val.FieldByName(name)
	if !field.IsValid() {
		return 0
	}

	err := setNativeFieldFromJSValue(field, ctx, ctx.newValue(C.JSValueRef(value)))
	if err != nil {
		*exception = ctx.newErrorOrPanic(err.Error())
		return 0
	}
	return 1
}

//export nativeobject_ConvertToString_go
func nativeobject_ConvertToString_go(data_ptr, ctx, obj unsafe.Pointer) unsafe.Pointer {
	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Can we get a string?
	if stringer, ok := data.val.Interface().(Stringer); ok {
		str := stringer.String()
		ret := NewString(str)
		return unsafe.Pointer(ret)
	}

	return nil
}

//=========================================================
// Native Method
//---------------------------------------------------------

func newNativeMethod(ctx *Context, obj *object_data, method int) *Object {
	data := &object_data{
		obj.typ,
		obj.val,
		method}
	register(data)

	ret := C.JSObjectMake(ctx.ref, nativemethod, unsafe.Pointer(data))
	return ctx.newObject(ret)
}

//export nativemethod_CallAsFunction_go
func nativemethod_CallAsFunction_go(data_ptr unsafe.Pointer, rawCtx C.JSContextRef, function C.JSObjectRef, thisObject C.JSObjectRef, argumentCount uint, arguments unsafe.Pointer, exception *C.JSValueRef) unsafe.Pointer {
	ctx := NewContextFrom(RawContext(rawCtx))
	defer func() {
		if r := recover(); r != nil {
			*exception = panicArgToJSString(ctx, r).ref
		}
	}()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Get the method
	method := data.val.Method(data.method)

	// Do the number of input parameters match?
	if method.Type().NumIn() != int(argumentCount) {
		panic(fmt.Sprintf("Incorrect number of function arguments! Got %d, expected %d!", method.Type().NumIn(), int(argumentCount)))
	}

	// Perform the call
	ret := docall(ctx, method, argumentCount, arguments)
	return unsafe.Pointer(ret.ref)
}
