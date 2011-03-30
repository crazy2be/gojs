package javascriptcore

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
// #include <JavaScriptCore/JSObjectRef.h>
// #include "callback.h"
import "C"
import "os"
import "reflect"
import "unsafe"

type object_data struct {
	typ reflect.Type
	val reflect.Value
	method int
}

var (
	nativecallback C.JSClassRef
	nativefunction C.JSClassRef
	nativeobject C.JSClassRef
	nativemethod C.JSClassRef
	objects map[uintptr] *object_data
)

type Stringer interface {
	String() string
}

func init() {
	// Create the class definition for JavaScriptCore
	nativecallback = C.JSClassDefinition_NativeCallback()
	if nativecallback == nil {
		panic( os.ENOMEM )
	}

	// Create the class definition for JavaScriptCore
	nativeobject = C.JSClassDefinition_NativeObject()
	if nativeobject == nil {
		panic( os.ENOMEM )
	}

	// Create the class definition for JavaScriptCore
	nativefunction = C.JSClassDefinition_NativeFunction()
	if nativefunction == nil {
		panic( os.ENOMEM )
	}

	// Create the class definition for JavaScriptCore
	nativemethod = C.JSClassDefinition_NativeMethod()
	if nativemethod == nil {
		panic( os.ENOMEM )
	}

	// Create map for native objects
	objects = make( map[uintptr] *object_data )
}

func value_to_javascript( ctx *Context, value reflect.Value ) *Value {

	// Allows functions to return JavaScriptCore values and objects
	// directly.  These we can return without conversion.
	if value.Type() == reflect.Typeof( (*Value)(nil) ) {
		// Type is already a JavaScriptCore value
		return value.Interface().(*Value)
	}
	if value.Type() == reflect.Typeof( (*Object)(nil) ) {
		// Type is already a JavaScriptCore object
		// nearly there
		return value.Interface().(*Object).ToValue()
	}

	// Handle simple types directly.  These can be identified by their
	// types in the package 'reflect'.
	switch value.(type) {
		case (*reflect.IntValue):
			r := value.(*reflect.IntValue).Get()
			return ctx.NewNumberValue( float64(r) )
		case (*reflect.UintValue):
			r := value.(*reflect.UintValue).Get()
			return ctx.NewNumberValue( float64(r) )
		case (*reflect.FloatValue):
			r := value.(*reflect.FloatValue).Get()
			return ctx.NewNumberValue( r )
		case (*reflect.StringValue):
			r := value.(*reflect.StringValue).Get()
			return ctx.NewStringValue( r )
		case (*reflect.PtrValue):
			r := value.(*reflect.PtrValue)
			_, ok := r.Elem().(*reflect.StructValue)
			if ok {
				ret := ctx.NewNativeObject( value.Interface() )
				return ret.ToValue()
			}
	}

	// No acceptable conversion found.
	panic( "Parameter can not be converted from Go native type." )
}

func recover_to_javascript( ctx *Context, r interface{} ) *Value {
	if re, ok := r.(os.Error); ok {
		// TODO:  Check for error return from NewError
		ret, _ := ctx.NewError( re.String() )		
		return (*Value)(unsafe.Pointer(ret))
	}
	if str, ok := r.(Stringer); ok {
		ret, _ := ctx.NewError( str.String() )		
		return (*Value)(unsafe.Pointer(ret))
	}
	if str, ok := reflect.NewValue(r).(*reflect.StringValue); ok {
		ret, _ := ctx.NewError( str.Get() )		
		return (*Value)(unsafe.Pointer(ret))
	}

	// Don't know how to convert this panic into a JavaScript error.
	// TODO:  Check for error return from NewError
	ret, _ := ctx.NewError( "Unknown panic from within Go." )		
	return (*Value)(unsafe.Pointer(ret))
}

func javascript_to_reflect( ctx *Context, param []*Value ) []reflect.Value {
	ret := make( []reflect.Value, len(param) )

	for index, item := range param {
		var goval interface{}

		switch ctx.ValueType( item ) {
		case TypeBoolean:
			goval = ctx.ToBoolean( item )
		case TypeNumber:
			goval = ctx.ToNumberOrDie( item )
		case TypeString:
			goval = ctx.ToStringOrDie( item )
		default:
			panic( "Parameter can not be converted to Go native type." )
		}

		ret[index] = reflect.NewValue( goval )
	}

	return ret
}

func javascript_to_value( field reflect.Value, ctx *Context, value *Value, exception *unsafe.Pointer )  {
	switch field.(type) {
		case *reflect.StringValue:
			str, err := ctx.ToString( value )
			if err == nil {
				field.(*reflect.StringValue).Set( str )
			} else {
				*exception = unsafe.Pointer( err )
			}

		case *reflect.FloatValue:
			flt, err := ctx.ToNumber( value )
			if err == nil {
				field.(*reflect.FloatValue).Set( flt )
			} else {
				*exception = unsafe.Pointer( err )
			}

		case *reflect.IntValue:
			flt, err := ctx.ToNumber( value )
			if err == nil {
				field.(*reflect.IntValue).Set( int64(flt) )
			} else {
				*exception = unsafe.Pointer( err )
			}

		case *reflect.UintValue:
			flt, err := ctx.ToNumber( value )
			if err == nil {
				if flt >= 0 {
					field.(*reflect.UintValue).Set( uint64(flt) )
				} else {
					err1, err := ctx.NewError( "Number must be greater than or equal to zero." )
					if err == nil {
						*exception = unsafe.Pointer( err1 )
					} else {
						*exception = unsafe.Pointer( err )
					}
				}
			} else {
				*exception = unsafe.Pointer( err )
			}

		default:
			panic( "Parameter can not be converted to Go native type." )
	}
}

//=========================================================
// Finalizer from JavaScriptCore for all native objects
//---------------------------------------------------------

func register( data *object_data ) {
	id := uintptr(unsafe.Pointer(data))
	objects[ id ] = data
}

//export finalize_go
func finalize_go( data unsafe.Pointer ) {
	// Called from JavaScriptCore finalizer methods
	id := uintptr( data )
	objects[id] = nil, false
}

//=========================================================
// Native Callback
//---------------------------------------------------------

type GoFunctionCallback func(ctx *Context, obj *Object, thisObject *Object, arguments []*Value) (ret *Value)

func (ctx *Context) NewFunctionWithCallback( callback GoFunctionCallback ) *Object {
	// Register the native Go object
	data := &object_data{
		reflect.Typeof( callback ),
		reflect.NewValue( callback ),
		0 }
	register( data )

	ret := C.JSObjectMake( C.JSContextRef(unsafe.Pointer(ctx)), nativecallback, unsafe.Pointer(data) )
	return (*Object)(unsafe.Pointer(ret))
}

//export nativecallback_CallAsFunction_go
func nativecallback_CallAsFunction_go( data_ptr unsafe.Pointer, ctx unsafe.Pointer, obj unsafe.Pointer, thisObject unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer ) unsafe.Pointer {
	defer func() {
		if r := recover(); r != nil {
			*exception = unsafe.Pointer( recover_to_javascript( (*Context)(ctx), r ) )
		}
	}()

	data := (*object_data)(data_ptr)
	ret := data.val.Interface().(GoFunctionCallback)(
		(*Context)(ctx), (*Object)(obj), (*Object)(thisObject), (*[1<<14]*Value)(arguments)[0:argumentCount] )
	return unsafe.Pointer(ret)
}

//=========================================================
// Native Function
//---------------------------------------------------------

func (ctx *Context) NewFunctionWithNative( fn interface{} ) *Object {
	// Sanity checks on the function
	if typ := reflect.Typeof( fn ).(*reflect.FuncType); typ.NumOut() > 1 {
		panic( "Bad native function:  too many output parameters" )
	}

	// Create Go-side registration
	data := &object_data{
		reflect.Typeof(fn),
		reflect.NewValue(fn),
		0 }
	register( data )

	ret := C.JSObjectMake( C.JSContextRef(unsafe.Pointer(ctx)), nativefunction, unsafe.Pointer(data) )
	return (*Object)(unsafe.Pointer(ret))
}

func docall( ctx *Context, val *reflect.FuncValue, argumentCount uint, arguments unsafe.Pointer ) (*Value) {
	// Step one, convert the JavaScriptCore array of arguments to 
	// an array of reflect.Values.  
	var in []reflect.Value
	if argumentCount!=0 {
		in = javascript_to_reflect( (*Context)(ctx), (*[1<<14]*Value)(arguments)[0:argumentCount] )
	}

	// Step two, perform the call
	out := val.Call( in )

	// Step three, convert the function return value back to JavaScriptCore
	if len(out) == 0 {
		return nil
	}
	// len(out) should be equal to 1
	return value_to_javascript( ctx, out[0] )
}

//export nativefunction_CallAsFunction_go
func nativefunction_CallAsFunction_go( data_ptr unsafe.Pointer, ctx unsafe.Pointer, _ unsafe.Pointer, _ unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer ) unsafe.Pointer {
	defer func() {
		if r := recover(); r != nil {
			*exception = unsafe.Pointer( recover_to_javascript( (*Context)(ctx), r ) )
		}
	}()

	// recover the object
	data := (*object_data)(data_ptr)
	typ := data.typ.(*reflect.FuncType)
	val := data.val.(*reflect.FuncValue)

	// Do the number of input parameters match?
	if typ.NumIn() != int(argumentCount) {
		panic( "Incorrect number of function arguments" )
	}

	ret := docall( (*Context)(ctx), val, argumentCount, arguments )
	return unsafe.Pointer( ret )
}

//=========================================================
// Native Object
//---------------------------------------------------------

func (ctx *Context) NewNativeObject( obj interface{} ) *Object {
	// The obj must be a pointer to a struct
	// TODO:  add error checking code

	data := &object_data{
		reflect.Typeof( obj ),
		reflect.NewValue( obj ),
		0 }
	register( data )

	ret := C.JSObjectMake( C.JSContextRef(unsafe.Pointer(ctx)), nativeobject, unsafe.Pointer(data) )
	return (*Object)(unsafe.Pointer(ret))
}

//export nativeobject_GetProperty_go
func nativeobject_GetProperty_go( data_ptr, ctx, _, propertyName unsafe.Pointer, exception *unsafe.Pointer ) unsafe.Pointer {
	// Get name of property as a go string
	name := (*String)(propertyName).String()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Drill down through reflect to find the property
	val := data.val
	if ptrvalue, ok := val.(*reflect.PtrValue); ok {
		val = ptrvalue.Elem()
	}
	struct_val, ok := val.(*reflect.StructValue)
	if !ok {
		return nil
	}

	// Can we locate a field with the proper name?
	field := struct_val.FieldByName( name )
	if field != nil {
		return unsafe.Pointer( value_to_javascript( (*Context)(ctx), field ) )
	}

	// Can we locate a method with the proper name?
	typ := data.typ
	for lp:=0; lp<typ.NumMethod(); lp++ {
		if typ.Method(lp).Name == name {
			ret := newNativeMethod( (*Context)(ctx), data, lp )
			return unsafe.Pointer(ret)
		}
	}

	// No matches found
	return nil
}

func internal_go_error( ctx *Context ) *Value {
	param := ctx.NewStringValue( "Internal Go error." )
	
	exception := (*Value)(nil)
	ret := C.JSObjectMakeError( C.JSContextRef(unsafe.Pointer(ctx)), 
		C.size_t(1), (*C.JSValueRef)( unsafe.Pointer( &param ) ),
		(*C.JSValueRef)(unsafe.Pointer(&exception)) )
	if ret != nil {
		return (*Value)(unsafe.Pointer(ret))
	}
	if exception != nil{
		return exception
	}
	panic( "Internal Go error." )
}

//export nativeobject_SetProperty_go
func nativeobject_SetProperty_go( data_ptr, ctx, _, propertyName, value unsafe.Pointer, exception *unsafe.Pointer ) C.char {
	// Get name of property as a go string
	name := (*String)(propertyName).String()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Drill down through reflect to find the property
	val := data.val
	if ptrvalue, ok := val.(*reflect.PtrValue); ok {
		val = ptrvalue.Elem()
	}
	struct_val, ok := val.(*reflect.StructValue)
	if !ok {
		*exception = unsafe.Pointer( internal_go_error( (*Context)(ctx) ) )
		return 1
	}

	field := struct_val.FieldByName( name )
	if field == nil {
		return 0
	}

	javascript_to_value( field, (*Context)(ctx), (*Value)(value), exception )
	return 1
}

//export nativeobject_ConvertToString_go
func nativeobject_ConvertToString_go( data_ptr, ctx, obj unsafe.Pointer ) unsafe.Pointer {
	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Can we get a string?
	if stringer, ok := data.val.Interface().(Stringer); ok {
		str := stringer.String()
		ret := NewString( str )
		return unsafe.Pointer( ret )
	}

	return nil
}

//=========================================================
// Native Method
//---------------------------------------------------------

func newNativeMethod( ctx *Context, obj *object_data, method int ) *Object {
	data := &object_data{
		obj.typ,
		obj.val,
		method }
	register( data )

	ret := C.JSObjectMake( C.JSContextRef(unsafe.Pointer(ctx)), nativemethod, unsafe.Pointer(data) )
	return (*Object)(unsafe.Pointer(ret))
}

//export nativemethod_CallAsFunction_go
func nativemethod_CallAsFunction_go( data_ptr, ctx, obj, thisObject unsafe.Pointer, argumentCount uint, arguments unsafe.Pointer, exception *unsafe.Pointer ) unsafe.Pointer {
	defer func() {
		if r := recover(); r != nil {
			*exception = unsafe.Pointer( recover_to_javascript( (*Context)(ctx), r ) )
		}
	}()

	// Reconstruct the object interface
	data := (*object_data)(data_ptr)

	// Get the method
	method := data.val.Method( data.method )

	// Do the number of input parameters match?
	if method.Type().(*reflect.FuncType).NumIn() != int(argumentCount)+1 {
		panic( "Incorrect number of function arguments" )
	}

	// Perform the call
	ret := docall( (*Context)(ctx), method, argumentCount, arguments )
	return unsafe.Pointer( ret )
}

