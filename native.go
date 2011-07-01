package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "os"
import "reflect"
import "unsafe"

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
		panic(os.ENOMEM)
	}

	// Create the class definition for JavaScriptCore
	nativeobject = C.JSClassDefinition_NativeObject()
	if nativeobject == nil {
		panic(os.ENOMEM)
	}

	// Create the class definition for JavaScriptCore
	nativefunction = C.JSClassDefinition_NativeFunction()
	if nativefunction == nil {
		panic(os.ENOMEM)
	}

	// Create the class definition for JavaScriptCore
	nativemethod = C.JSClassDefinition_NativeMethod()
	if nativemethod == nil {
		panic(os.ENOMEM)
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

// Given a reflect.Value, this function examines the type and returns a javascript value that best represents the given value. If no acceptable conversion can be found, it panics.
func (ctx *Context) reflectToJSValue(value reflect.Value) *Value {
	panic("")
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
			panic("Called reflectToJSValue() with a pointer to an array or slice. Most likely, this is a native javascript object you are passing by accident, and meant to pass to newValue() rather than NewValue(). If you really are trying to convert a pointer to an array into a native javascript object, dereference it first (slices are pointers internally anyway, so no loss in efficiency).")
			//log.Println("About to make new native object from *[0]uint8")
			//ret := ctx.NewNativeObject(value.Interface())
			//log.Println("Made new native object from *[0]uint8")
			//return ret.ToValue()
		}
	}
	// No acceptable conversion found.
	panic("Parameter can not be converted from Go native type. Type is "+value.Kind().String()+", value is "+value.String())
}

func recover_to_javascript(ctx *Context, r interface{}) *Value {
	if re, ok := r.(os.Error); ok {
		// TODO:  Check for error return from NewError
		ret, _ := ctx.NewError(re.String())
		return ret.ToValue()
	}
	if str, ok := r.(Stringer); ok {
		ret, _ := ctx.NewError(str.String())
		return ret.ToValue()
	}
	if str := reflect.ValueOf(r); str.Kind() == reflect.String {
		ret, _ := ctx.NewError(str.String())
		return ret.ToValue()
	}

	// Don't know how to convert this panic into a JavaScript error.
	// TODO:  Check for error return from NewError
	ret, _ := ctx.NewError("Unknown panic from within Go.")
	return ret.ToValue()
}

func (ctx *Context) jsValuesToReflect(param []*Value) []reflect.Value {
	ret := make([]reflect.Value, len(param))

	for index, item := range param {
		var goval interface{}

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

func javascript_to_value(field reflect.Value, ctx *Context, value *Value, exception *unsafe.Pointer) {
	switch field.Kind() {
	case reflect.String:
		str, err := ctx.ToString(value)
		if err == nil {
			field.SetString(str)
		} else {
			*exception = unsafe.Pointer(err)
		}

	case reflect.Float32, reflect.Float64:
		flt, err := ctx.ToNumber(value)
		if err == nil {
			field.SetFloat(flt)
		} else {
			*exception = unsafe.Pointer(err)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		flt, err := ctx.ToNumber(value)
		if err == nil {
			field.SetInt(int64(flt))
		} else {
			*exception = unsafe.Pointer(err)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		flt, err := ctx.ToNumber(value)
		if err == nil {
			if flt >= 0 {
				field.SetUint(uint64(flt))
			} else {
				err1, err := ctx.NewError("Number must be greater than or equal to zero.")
				if err == nil {
					*exception = unsafe.Pointer(err1)
				} else {
					*exception = unsafe.Pointer(err)
				}
			}
		} else {
			*exception = unsafe.Pointer(err)
		}

	default:
		panic("Parameter can not be converted to Go native type.")
	}
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
	objects[id] = nil, false
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
func nativecallback_CallAsFunction_go(data_ptr unsafe.Pointer, ctx unsafe.Pointer, obj unsafe.Pointer, thisObject unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer) unsafe.Pointer {
	defer func() {
		if r := recover(); r != nil {
			*exception = unsafe.Pointer(recover_to_javascript((*Context)(ctx), r))
		}
	}()

	data := (*object_data)(data_ptr)
	ret := data.val.Interface().(GoFunctionCallback)(
		(*Context)(ctx), (*Object)(obj), (*Object)(thisObject), (*[1 << 14]*Value)(arguments)[0:argumentCount])
	return unsafe.Pointer(ret)
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
		in = ctx.jsValuesToReflect((*[1 << 14]*Value)(arguments)[0:argumentCount])
	}

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
func nativefunction_CallAsFunction_go(data_ptr unsafe.Pointer, ctx unsafe.Pointer, _ unsafe.Pointer, _ unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer) unsafe.Pointer {
	defer func() {
		if r := recover(); r != nil {
			*exception = unsafe.Pointer(recover_to_javascript((*Context)(ctx), r))
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

	ret := docall((*Context)(ctx), val, argumentCount, arguments)
	return unsafe.Pointer(ret)
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
func nativeobject_GetProperty_go(data_ptr, ctx, _, propertyName unsafe.Pointer, exception *unsafe.Pointer) unsafe.Pointer {
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
		return unsafe.Pointer((*Context)(ctx).reflectToJSValue(field).ref)
	}

	// Can we locate a method with the proper name?
	typ := data.typ
	for lp := 0; lp < typ.NumMethod(); lp++ {
		if typ.Method(lp).Name == name {
			ret := newNativeMethod((*Context)(ctx), data, lp)
			return unsafe.Pointer(ret)
		}
	}

	// No matches found
	return nil
}

func internal_go_error(ctx *Context) *Value {
	param := ctx.NewStringValue("Internal Go error.")

	exception := C.JSValueRef(unsafe.Pointer(nil))
	ret := C.JSObjectMakeError(ctx.ref,
		C.size_t(1), &param.ref,
		&exception)
	if ret != nil {
		return ctx.NewObject(ret).ToValue()
	}
	if exception != nil {
		return ctx.newValue(exception)
	}
	panic("Internal Go error.")
}

//export nativeobject_SetProperty_go
func nativeobject_SetProperty_go(data_ptr, ctx, _, propertyName, value unsafe.Pointer, exception *unsafe.Pointer) C.char {
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
		*exception = unsafe.Pointer(internal_go_error((*Context)(ctx)))
		return 1
	}

	field := struct_val.FieldByName(name)
	if !field.IsValid() {
		return 0
	}

	javascript_to_value(field, (*Context)(ctx), (*Value)(value), exception)
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
	return (*Object)(unsafe.Pointer(ret))
}

//export nativemethod_CallAsFunction_go
func nativemethod_CallAsFunction_go(data_ptr, ctx, obj, thisObject unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer) unsafe.Pointer {
	defer func() {
		if r := recover(); r != nil {
			*exception = unsafe.Pointer(recover_to_javascript((*Context)(ctx), r))
		}
	}()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Get the method
	method := data.val.Method(data.method)

	// Do the number of input parameters match?
	if method.Type().NumIn() != int(argumentCount)+1 {
		panic("Incorrect number of function arguments")
	}

	// Perform the call
	ret := docall((*Context)(ctx), method, argumentCount, arguments)
	return unsafe.Pointer(ret)
}
