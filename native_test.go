package gojs

import (
	"log"
	"syscall"
	"testing"
	"unsafe"
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

func PrettyPrintValArr(t *testing.T, values []*Value) {
	for _, val := range values {
		tlog(t, val)
	}
}

func ValArrToStrings(t *testing.T, values []*Value) {
	for _, val := range values {
		tlog(t, val.ctx.ToStringOrDie(val))
	}
}

func checkExpected(t *testing.T, val *Value, str string) {
	str2 := val.ctx.ToStringOrDie(val)
	if str2 != str {
		t.Fatalf("Got %s, expected %s", str2, str)
	}
}

func checkArrsEqual(t *testing.T, vals []*Value, expectedVals []*Value) {
	for i, _ := range vals {
		valstr := vals[i].ctx.ToStringOrDie(vals[i])
		expvalstr := expectedVals[i].ctx.ToStringOrDie(expectedVals[i])
		if valstr != expvalstr {
			t.Fatalf("Got %s, expected %s", valstr, expvalstr)
		}
	}
}

// t.Log doesn't print things immediately, this does if TESTING_DEBUG_LOG is set to true. Useful when you have pointer crashes and faults such as are common with cgo code.
const TESTING_DEBUG_LOG = false

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
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func TestNewCValueArray(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	valarr := make([]*Value, 5)
	valstrs := make([]string, 5)
	valarr[0] = ctx.NewValue(0)
	valstrs[0] = "0"
	valarr[1] = ctx.NewNumberValue(1.3)
	valstrs[1] = "1.3"
	valarr[2] = ctx.NewValue(nil)
	valstrs[2] = "null"
	valarr[3] = ctx.NewValue(2309240)
	valstrs[3] = "2309240"
	valarr[4] = ctx.NewValue(0x934)
	valstrs[4] = "2356"

	cptr, size := ctx.newCValueArray(valarr)

	// Make sure the C array is correct
	ptrs := unsafe.Sizeof(uintptr(0))
	uptr := uintptr(unsafe.Pointer(cptr))
	ptr0 := unsafe.Pointer(uptr + 0)
	val0 := ctx.ptrValue(ptr0)
	ptr1 := unsafe.Pointer(uptr + (ptrs * 1))
	val1 := ctx.ptrValue(ptr1)
	ptr2 := unsafe.Pointer(uptr + (ptrs * 2))
	val2 := ctx.ptrValue(ptr2)
	ptr3 := unsafe.Pointer(uptr + (ptrs * 3))
	val3 := ctx.ptrValue(ptr3)
	ptr4 := unsafe.Pointer(uptr + (ptrs * 4))
	val4 := ctx.ptrValue(ptr4)

	checkExpected(t, val0, valstrs[0])
	checkExpected(t, val1, valstrs[1])
	checkExpected(t, val2, valstrs[2])
	checkExpected(t, val3, valstrs[3])
	checkExpected(t, val4, valstrs[4])

	origarray := ctx.newGoValueArray(unsafe.Pointer(cptr), uint(size))

	checkArrsEqual(t, valarr, origarray)
}

func TestNewFunctionWithCallback(t *testing.T) {
	var flag bool
	callback := func(ctx *Context, obj *Object, thisObject *Object, _ []*Value) *Value {
		tlog(t, "In callback function!")
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
	tlog(t, "Made new function with callback")
	if !ctx.IsFunction(fn) {
		t.Errorf("ctx.NewFunctionWithCallback returned value that is not a function")
	}
	tlog(t, "Function is a callback function")
	if ctx.ToStringOrDie(fn.ToValue()) != "nativecallback" {
		t.Errorf("ctx.NewFunctionWithCallback returned value that does not convert to property string")
	}
	tlog(t, "Successfully converted to property string")
	tlog(t, "Calling as function...")
	ctx.CallAsFunction(fn, nil, []*Value{})
	if !flag {
		t.Errorf("Native function did not execute")
	}
	tlog(t, "Called as function.")
}

func TestNewFunctionWithCallback2(t *testing.T) {
	callback := func(ctx *Context, obj *Object, thisObject *Object, args []*Value) *Value {
		tlog(t, "In callback function!")
		if len(args) != 2 {
			return nil
		}
		tlog(t, "Attempting to convert args to numbers...", args)
		tlog(t, ctx)
		tlog(t, args[0], args[1])
		tlog(t, args[0].ctx, args[1].ctx)
		tlog(t, args[0].ref, args[1].ref)
		//a, err := args[0].ctx.ToNumber(args[0])
		//tlog(t, a, err)
		//return ctx.NewNumberValue(2)
		a, err := ctx.ToNumber(args[0])
		tlog(t, err)
		tlog(t, "Converted first arg...")
		b := ctx.ToNumberOrDie(args[1])
		return ctx.NewNumberValue(a + b)
		return ctx.NewNumberValue(1)
		return ctx.newValue(nil)
	}

	tlog(t, "Acquiring context!")
	ctx := NewContext()
	defer ctx.Release()

	tlog(t, "Creating a new function with callback")
	fn := ctx.NewFunctionWithCallback(callback)
	tlog(t, "Ceating new number values")
	numarr := make([]*Value, 2)
	numarr[0] = ctx.NewNumberValue(1.5)
	numarr[1] = ctx.NewNumberValue(3.0)
	tlog(t, "Calling callback as function")
	val, err := ctx.CallAsFunction(fn, nil, numarr)
	tlog(t, "Called callback as function!")
	tlog(t, err, val)
	if err != nil || val == nil {
		t.Fatalf("Error executing native callback")
	}
	if ctx.ToNumberOrDie(val) != 4.5 {
		t.Fatalf("Native callback did not return the correct value")
	}
}

func TestNewFunctionWithCallbackPanic(t *testing.T) {
	var callbacks = []GoFunctionCallback{}
	var error_msgs = []string{"error from go!", syscall.ENOMEM.Error()}

	callbacks = append(callbacks,
		func(ctx *Context, obj *Object, thisObject *Object, _ []*Value) *Value {
			panic("error from go!")
			return nil
		},
	)
	callbacks = append(callbacks,
		func(ctx *Context, obj *Object, thisObject *Object, _ []*Value) *Value {
			panic(syscall.ENOMEM)
			return nil
		},
	)

	ctx := NewContext()
	defer ctx.Release()

	for index, callback := range callbacks {
		fn := ctx.NewFunctionWithCallback(callback)
		if fn == nil {
			t.Fatalf("ctx.NewFunctionWithCallback failed")
			return
		}
		if !ctx.IsFunction(fn) {
			t.Fatalf("ctx.NewFunctionWithCallback returned value that is not a function")
		}
		if ctx.ToStringOrDie(fn.ToValue()) != "nativecallback" {
			t.Fatalf("ctx.NewFunctionWithCallback returned value that does not convert to property string")
		}
		val, err := ctx.CallAsFunction(fn, nil, nil)
		if val != nil {
			t.Fatalf("ctx.NewFunctionWithCallback that panicked returned a value")
		}
		if err == nil {
			t.Fatalf("want err != nil for function that panicked")
		}
		if err.Error() != error_msgs[index] {
			t.Fatalf("want exception message %q, got %q", err, error_msgs[index])
		}
	}
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
	ctx.CallAsFunction(fn, nil, nil)
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
	args := make([]*Value, 2)
	args[0] = ctx.NewNumberValue(1.5)
	args[1] = ctx.NewNumberValue(3.0)
	val, err := ctx.CallAsFunction(fn, nil, args)
	tlog(t, "Called as function")
	if err != nil || val == nil {
		t.Errorf("Error executing native function (%v)", err)
	}

	wantNum := 4.5
	gotNum, err := ctx.ToNumber(val)
	if err != nil {
		t.Error("ToNumber", err)
	} else if wantNum != gotNum {
		t.Errorf("want return value %v, got %v", wantNum, gotNum)
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
	if err != nil || val == nil {
		t.Errorf("Error executing native function (%v)", err)
	}
	if ctx.ToNumberOrDie(val) != 4.5 {
		t.Errorf("Native function did not return the correct value")
	}
}

func TestNativeFunctionPanic(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	callbacks := []func(){
		func() { panic("Panic!") }, func() { panic(syscall.ENOMEM) }}

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
		tlog(t, "Called as function")
		tlog(t, val, err)
		if err == nil || val != nil {
			t.Fatalf("ctx.CallAsFunction did not panic as expected")
		}
		tlog(t, "About to convert to string!")
		if err.Error() == "" {
			t.Errorf("want non-empty error message, got %q", err)
		} else {
			t.Logf("ctx.CallAsFunction paniced as expected (%v)", err)
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

	tlog(t, "Creating new native object from", obj)

	v := ctx.NewNativeObject(obj)

	tlog(t, "Setting property n")

	ctx.SetProperty(ctx.GlobalObject(), "n", v.ToValue(), 0)

	tlog(t, "Creating value for property I")

	// Set the integer property
	i := ctx.NewNumberValue(-2)
	tlog(t, "Setting property I", i)
	tlog(t, ctx.ToStringOrDie(i))
	ctx.SetProperty(v, "I", i, 0)
	tlog(t, "Set property I, checking for errors.")
	if obj.I != -2 {
		t.Errorf("ctx.SetProperty did not set integer field correctly")
	}

	tlog(t, "Setting property U")

	// Set the unsigned integer property
	u := ctx.NewNumberValue(3)
	ctx.SetProperty(v, "U", u, 0)
	if obj.U != 3 {
		t.Fatalf("ctx.SetProperty did not set unsigned integer field correctly")
	}

	// 	t.Error("Skipping setting property U to invalid value, it currently causes a fault.")

	// Set the unsigned integer property
	u = ctx.NewNumberValue(-3)
	tlog(t, "Setting property U to invalid value")
	err := ctx.SetProperty(v, "U", u, 0)
	tlog(t, "Set property, checking for correctness...")
	if err == nil {
		t.Errorf("ctx.SetProperty did not set unsigned integer field correctly: No error was returned")
	} else {
		t.Logf("%v", err)
	}
	if obj.U != 3 {
		t.Errorf("ctx.SetProperty did not set unsigned integer field correctly: Incorrect value")
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

	tlog(t, "Testing n.Add()")

	// Following script access should be successful
	ret, err := ctx.EvaluateScript("n.Add()", nil, "./testing.go", 1)

	tlog(t, "Evaluated Script")
	if err != nil {
		t.Errorf("ctx.EvaluateScript returned an error: %#v", err)
		return
	}
	tlog(t, "No error")
	if ret == nil {
		t.Errorf("ctx.EvaluateScript did not return a result! (no error)")
		return
	}
	tlog(t, "Result was returned")
	if !ctx.IsNumber(ret) {
		t.Errorf("ctx.EvaluateScript did not return 'number' result when calling method 'Add'.")
	}
	tlog(t, "sucessfully checked that value was a number")
	num := ctx.ToNumberOrDie(ret)
	if num != 2.0 {
		t.Errorf("ctx.EvaluateScript incorrect value when accessing native object's field.")
	}

	tlog(t, "Testing n.AddWith()")

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
