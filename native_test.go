package gojs

import (
	"testing"
	"log"
	"os"
)

type reflect_object struct {
	I int
	U uint
	F float64
	S string
}

func (o *reflect_object) String() string {
	return o.S
}

func (o *reflect_object) Add() float64 {
	return float64(o.I) + o.F
}

func (o *reflect_object) AddWith(op float64) float64 {
	return float64(o.I) + o.F + op
}

func (o *reflect_object) Self() *reflect_object {
	return o
}

func (o *reflect_object) Null() *reflect_object {
	return nil
}

func TestNewFunctionWithCallback(t *testing.T) {
	var flag bool
	callback := func(ctx *Context, obj *Object, thisObject *Object, _ []*Value) *Value {
		flag = true
		return nil
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithCallback(callback)
	if fn == nil {
		t.Errorf("ctx.NewFunctionWithCallback failed")
		return
	}
	if !ctx.IsFunction(fn) {
		t.Errorf("ctx.NewFunctionWithCallback returned value that is not a function")
	}
	if ctx.ToStringOrDie(fn.ToValue()) != "nativecallback" {
		t.Errorf("ctx.NewFunctionWithCallback returned value that does not convert to property string")
	}
	ctx.CallAsFunction(fn, nil, []*Value{})
	if !flag {
		t.Errorf("Native function did not execute")
	}
}

// t.Log doesn't print things immediately, this does if TESTING_DEBUG_LOG is set to true. Useful when you have pointer crashes and faults such as are common with cgo code.
const TESTING_DEBUG_LOG = true

func tlog(t *testing.T, v ...interface{}) {
	if TESTING_DEBUG_LOG {
		log.Println(v...)
	} else {
		t.Log(v...)
	}
	return
}

func terrf(t *testing.T, format string, v ...interface{}) {
	if TESTING_DEBUG_LOG {
		log.Printf(format, v...)
		t.Fail()
	} else {
		t.Errorf(format, v...)
	}
}

func init() {
	log.SetFlags(log.Ltime|log.Lshortfile)
}

func TestNewFunctionWithCallback2(t *testing.T) {
	callback := func(ctx *Context, obj *Object, thisObject *Object, args []*Value) *Value {
		tlog(t, "In callback function!")
		if len(args) != 2 {
			return nil
		}
		tlog(t, "Attempting to convert args to numbers...", args)
		a, err := ctx.ToNumber(args[0])
		tlog(t, err)
		tlog(t, "Converted first arg...")
		b := ctx.ToNumberOrDie(args[1])
		return ctx.NewNumberValue(a + b)
	}

	tlog(t, "Acquiring context!")
	ctx := NewContext()
	defer ctx.Release()

	tlog(t, "Creating a new function with callback")
	fn := ctx.NewFunctionWithCallback(callback)
	tlog(t, "Ceating new number values")
	a := ctx.NewNumberValue(1.5)
	b := ctx.NewNumberValue(3.0)
	tlog(t, "Calling callback as function")
	val, err := ctx.CallAsFunction(fn, nil, []*Value{a, b})
	if err != nil || val == nil {
		t.Errorf("Error executing native callback")
	}
	if ctx.ToNumberOrDie(val) != 4.5 {
		t.Errorf("Native callback did not return the correct value")
	}
}

func TestNewFunctionWithCallbackPanic(t *testing.T) {
	var callbacks = []GoFunctionCallback{}
	var error_msgs = []string{"error from go!", os.ENOMEM.String()}

	callbacks = append(callbacks,
		func(ctx *Context, obj *Object, thisObject *Object, _ []*Value) *Value {
			panic("error from go!")
			return nil
		})
	callbacks = append(callbacks,
		func(ctx *Context, obj *Object, thisObject *Object, _ []*Value) *Value {
			panic(os.ENOMEM)
			return nil
		})

	ctx := NewContext()
	defer ctx.Release()

	for index, callback := range callbacks {

		fn := ctx.NewFunctionWithCallback(callback)
		if fn == nil {
			t.Errorf("ctx.NewFunctionWithCallback failed")
			return
		}
		if !ctx.IsFunction(fn) {
			t.Errorf("ctx.NewFunctionWithCallback returned value that is not a function")
		}
		if ctx.ToStringOrDie(fn.ToValue()) != "nativecallback" {
			t.Errorf("ctx.NewFunctionWithCallback returned value that does not convert to property string")
		}
		val, err := ctx.CallAsFunction(fn, nil, []*Value{})
		if val != nil {
			t.Errorf("ctx.NewFunctionWithCallback that panicked returned a value")
		}
		if err == nil || !ctx.IsObject(err.val) {
			t.Errorf("ctx.NewFunctionWithCallback that panicked did not set exception")
		}
		if ctx.ToStringOrDie(err.val) != "Error: "+error_msgs[index] {
			t.Errorf("ctx.NewFunctionWithCallback that panicked did not set exception message (%v,%v)",
				ctx.ToStringOrDie(err.val), error_msgs[index])
		}

	} // for
}

func TestNativeFunction(t *testing.T) {
	var flag bool
	callback := func() {
		flag = true
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithNative(callback)
	if fn == nil {
		t.Errorf("ctx.NewFunctionWithNative failed")
		return
	}
	if !ctx.IsFunction(fn) {
		t.Errorf("ctx.NewFunctionWithNative returned value that is not a function")
	}
	if ctx.ToStringOrDie(fn.ToValue()) != "nativefunction" {
		t.Errorf("ctx.nativefunction returned value that does not convert to property string")
	}
	ctx.CallAsFunction(fn, nil, []*Value{})
	if !flag {
		t.Errorf("Native function did not execute")
	}
}

func TestNativeFunction2(t *testing.T) {
	callback := func(a float64, b float64) float64 {
		return a + float64(b)
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithNative(callback)
	if fn == nil {
		t.Errorf("ctx.NewFunctionWithNative failed")
		return
	}
	if !ctx.IsFunction(fn) {
		t.Errorf("ctx.NewFunctionWithNative returned value that is not a function")
	}
	a := ctx.NewNumberValue(1.5)
	b := ctx.NewNumberValue(3.0)
	val, err := ctx.CallAsFunction(fn, nil, []*Value{a, b})
	if err.val != nil || val == nil {
		t.Errorf("Error executing native function (%v)", ctx.ToStringOrDie(err.val))
	}
	if ctx.ToNumberOrDie(val) != 4.5 {
		t.Errorf("Native function did not return the correct value")
	}
}

func TestNativeFunction3(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	callback := func(a float64, b float64) *Value {
		ret := a + float64(b)
		return ctx.NewNumberValue(ret)
	}

	fn := ctx.NewFunctionWithNative(callback)
	if fn == nil {
		t.Errorf("ctx.NewFunctionWithNative failed")
		return
	}
	if !ctx.IsFunction(fn) {
		t.Errorf("ctx.NewFunctionWithNative returned value that is not a function")
	}
	a := ctx.NewNumberValue(1.5)
	b := ctx.NewNumberValue(3.0)
	val, err := ctx.CallAsFunction(fn, nil, []*Value{a, b})
	if err.val != nil || val == nil {
		t.Errorf("Error executing native function (%v)", ctx.ToStringOrDie(err.val))
	}
	if ctx.ToNumberOrDie(val) != 4.5 {
		t.Errorf("Native function did not return the correct value")
	}
}

func TestNativeFunctionPanic(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	callbacks := []func(){
		func() { panic("Panic!") }, func() { panic(os.ENOMEM) }}

	for _, callback := range callbacks {

		fn := ctx.NewFunctionWithNative(callback)
		if fn == nil {
			t.Errorf("ctx.NewFunctionWithNative failed")
			return
		}
		if !ctx.IsFunction(fn) {
			t.Errorf("ctx.NewFunctionWithNative returned value that is not a function")
		}
		val, err := ctx.CallAsFunction(fn, nil, nil)
		if err.val == nil || val != nil {
			t.Errorf("ctx.CallAsFunction did not panic as expected")
		}
		msg := ctx.ToStringOrDie(err.val)
		if msg[0:7] != "Error: " {
			t.Errorf("ctx.CallAsFunction did return expected error object (%v)", msg)
		} else {
			t.Logf("ctx.CallAsFunction paniced as expected (%v)", msg)
		}

	}
}

func TestNewNativeObject(t *testing.T) {
	obj := &reflect_object{-1, 2, 3.0, "four"}

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject(obj)
	ctx.SetProperty(ctx.GlobalObject(), "n", v.ToValue(), 0)

	// Following script access should be successful
	ret, err := ctx.EvaluateScript("n.F", nil, "./testing.go", 1)
	if err != nil {
		t.Errorf("ctx.EvaluateScript returned an error: %#v", err)
		return
	}
	if ret == nil {
		t.Errorf("ctx.EvaluateScript did not return a result (no error specified)!")
		return
	}
	if !ctx.IsNumber(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'number' result when accessing native object's non-existent field.")
	}
	num := ctx.ToNumberOrDie(ret)
	if num != 3.0 {
		t.Errorf("ctx.EvaluateScript incorrect value when accessing native object's field.")
	}

	// following script access should fail
	ret, err = ctx.EvaluateScript("n.noexist", nil, "./testing.go", 1)
	if err != nil || ret == nil {
		t.Errorf("ctx.EvaluateScript returned an error (or did not return a result)")
	}
	if !ctx.IsUndefined(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'undefined' result when accessing native object's non-existent field.")
	}

	// following script access should succeed
	ret, err = ctx.EvaluateScript("n.S", nil, "./testing.go", 1)
	if err != nil || ret == nil {
		t.Errorf("ctx.EvaluateScript returned an error (or did not return a result)")
	}
	if !ctx.IsString(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'string' result when accessing native object's non-existent field.")
	}
	str := ctx.ToStringOrDie(ret)
	if str != "four" {
		t.Errorf("ctx.EvaluateScript incorrect value when accessing native object's field.")
	}
}

func TestNewNativeObjectSet(t *testing.T) {
	obj := &reflect_object{-1, 2, 3.0, "four"}

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject(obj)
	ctx.SetProperty(ctx.GlobalObject(), "n", v.ToValue(), 0)

	// Set the integer property
	i := ctx.NewNumberValue(-2)
	ctx.SetProperty(v, "I", i, 0)
	if obj.I != -2 {
		t.Errorf("ctx.SetProperty did not set integer field correctly")
	}

	// Set the unsigned integer property
	u := ctx.NewNumberValue(3)
	ctx.SetProperty(v, "U", u, 0)
	if obj.U != 3 {
		t.Errorf("ctx.SetProperty did not set unsigned integer field correctly")
	}

	// Set the unsigned integer property
	u = ctx.NewNumberValue(-3)
	err := ctx.SetProperty(v, "U", u, 0)
	if err == nil {
		t.Errorf("ctx.SetProperty did not set unsigned integer field correctly")
	} else {
		t.Logf("%v", err)
	}
	if obj.U != 3 {
		t.Errorf("ctx.SetProperty did not set unsigned integer field correctly")
	}

	// Set the float property
	n := ctx.NewNumberValue(4.0)
	ctx.SetProperty(v, "F", n, 0)
	if obj.F != 4.0 {
		t.Errorf("ctx.SetProperty did not set float field correctly")
	}

	s := ctx.NewStringValue("five")
	ctx.SetProperty(v, "S", s, 0)
	if obj.S != "five" {
		t.Errorf("ctx.SetProperty did not set string field correctly")
	}
}

func TestNewNativeObjectConvert(t *testing.T) {
	obj := &reflect_object{-1, 2, 3.0, "four"}

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject(obj)

	if ctx.ToStringOrDie(v.ToValue()) != "four" {
		t.Errorf("ctx.ToStringOrDie for native object did not return correct value.")
	}
}

func TestNewNativeObjectMethod(t *testing.T) {
	obj := &reflect_object{-1, 2, 3.0, "four"}

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject(obj)
	ctx.SetProperty(ctx.GlobalObject(), "n", v.ToValue(), 0)

	// Following script access should be successful
	ret, err := ctx.EvaluateScript("n.Add()", nil, "./testing.go", 1)
	if err != nil {
		t.Errorf("ctx.EvaluateScript returned an error: %#v", *err)
		return
	}
	if ret == nil {
		t.Errorf("ctx.EvaluateScript did not return a result! (no error)")
		return
	}
	if !ctx.IsNumber(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'number' result when calling method 'Add'.")
	}
	num := ctx.ToNumberOrDie(ret)
	if num != 2.0 {
		t.Errorf("ctx.EvaluateScript incorrect value when accessing native object's field.")
	}

	// Following script access should be successful
	ret, err = ctx.EvaluateScript("n.AddWith(0.5)", nil, "./testing.go", 1)
	if err != nil || ret == nil {
		t.Errorf("ctx.EvaluateScript returned an error (or did not return a result)")
		return
	}
	if !ctx.IsNumber(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'number' result when calling method 'AddWith'.")
	}
	num = ctx.ToNumberOrDie(ret)
	if num != 2.5 {
		t.Errorf("ctx.EvaluateScript incorrect value when accessing native object's field.")
	}

	// Following script access should be successful
	ret, err = ctx.EvaluateScript("n.Self()", nil, "./testing.go", 1)
	if err != nil || ret == nil {
		t.Errorf("ctx.EvaluateScript returned an error (or did not return a result)")
		return
	}
	if !ctx.IsObject(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'object' result when calling method 'Self'.")
	}

	// Following script access should be successful
	ret, err = ctx.EvaluateScript("n.Null()", nil, "./testing.go", 1)
	if err != nil || ret == nil {
		t.Errorf("ctx.EvaluateScript 'n.Null()' returned an error (or did not return a result)")
		t.Logf("Error:  %s", err)
		return
	}
	if !ctx.IsNull(ret) {
		t.Errorf("ctx.EvaluateScript 'n.Null()'did not return a javascript null value.")
	}
}
