package gojs

import (
	"testing"
)

func TestNewObject(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewEmptyObject()
	if val == nil {
		t.Errorf("ctx.NewObject returned a nil poitner")
	}
	if !ctx.IsObject(val.ToValue()) {
		t.Errorf("ctx.NewObject failed to return an object (%v)", ctx.ValueType(val.ToValue()))
	}
}

func TestNewArray(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val, err := ctx.NewArray(nil)
	tlog(t, val)
	if err != nil {
		t.Fatalf("ctx.NewArray returned an exception (%v)", err)
	}
	if val == nil {
		t.Fatalf("ctx.NewArray returned a nil poitner")
	}
	if !ctx.IsObject(val.ToValue()) {
		t.Fatalf("ctx.NewArray failed to return an object (%v)", ctx.ValueType(val.ToValue()))
	}
}

func TestNewArray2(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	a := ctx.NewNumberValue(1.5)
	b := ctx.NewNumberValue(3.0)

	val, err := ctx.NewArray([]*Value{a, b})
	if err != nil {
		t.Fatalf("ctx.NewArray returned an exception (%v)", err)
	}
	if val == nil {
		t.Fatalf("ctx.NewArray returned a nil poitner")
	}
	if !ctx.IsObject(val.ToValue()) {
		t.Fatalf("ctx.NewArray failed to return an object (%v)", ctx.ValueType(val.ToValue()))
	}
	prop, err := ctx.GetProperty(val, "length")
	if err != nil || prop == nil {
		t.Fatalf("ctx.NewArray returned object without 'length' property")
	} else {
		if !ctx.IsNumber(prop) {
			t.Errorf("ctx.NewArray return object with 'length' property not a number")
		}
		if ctx.ToNumberOrDie(prop) != 2 {
			t.Errorf("ctx.NewArray return object with 'length' not equal to 2", ctx.ToNumberOrDie(prop))
		}
	}
}

func TestNewDate(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val, err := ctx.NewDate()
	if err != nil {
		t.Fatalf("ctx.NewDate returned an exception (%v)", err)
	}
	if val == nil {
		t.Fatalf("ctx.NewDate returned a nil poitner")
	}
	if !ctx.IsObject(val.ToValue()) {
		t.Fatalf("ctx.NewDate failed to return an object (%v)", ctx.ValueType(val.ToValue()))
	}
}

func TestNewDateWithMilliseconds(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val, err := ctx.NewDateWithMilliseconds(3600000)
	if err != nil {
		t.Errorf("ctx.NewDateWithMilliseconds returned an exception (%v)", err)
	}
	if val == nil {
		t.Errorf("ctx.NewDateWithMilliseconds returned a nil poitner")
	}
	if !ctx.IsObject(val.ToValue()) {
		t.Errorf("ctx.NewDateWithMilliseconds failed to return an object (%v)", ctx.ValueType(val.ToValue()))
	}
}

func TestNewDateWithString(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val, err := ctx.NewDateWithString("01-Oct-2010")
	if err != nil {
		t.Errorf("ctx.NewDateWithString returned an exception (%v)", err)
	}
	if val == nil {
		t.Errorf("ctx.NewDateWithString returned a nil poitner")
	}
	if !ctx.IsObject(val.ToValue()) {
		t.Errorf("ctx.NewDateWithString failed to return an object (%v)", ctx.ValueType(val.ToValue()))
	}
}

func TestNewError(t *testing.T) {
	tests := []string{"test error 1", "test error 2"}

	ctx := NewContext()
	defer ctx.Release()

	for _, item := range tests {
		r, err := ctx.NewError(item)
		if err != nil {
			t.Errorf("ctx.NewError failed on string %v with error %v", item, err)
		}
		v, exc := ctx.GetProperty(r, "name")
		if exc != nil || v == nil {
			t.Errorf("ctx.NewError returned object without 'message' property")
		} else {
			if !ctx.IsString(v) {
				t.Errorf("ctx.NewError return object with 'message' property not a string")
			}
			if ctx.ToStringOrDie(v) != "Error" {
				t.Errorf("JavaScript error object and input string don't match (%v, %v)", item, ctx.ToStringOrDie(v))
			}
		}
		v, exc = ctx.GetProperty(r, "message")
		if exc != nil || v == nil {
			t.Errorf("ctx.NewError returned object without 'message' property")
		} else {
			if !ctx.IsString(v) {
				t.Errorf("ctx.NewError return object with 'message' property not a string")
			}
			if ctx.ToStringOrDie(v) != item {
				t.Errorf("JavaScript error object and input string don't match (%v, %v)", item, ctx.ToStringOrDie(v))
			}
		}
	}
}

func TestNewRegExp(t *testing.T) {
	tests := []string{"\\bt[a-z]+\\b", "[0-9]+(\\.[0-9]*)?"}

	ctx := NewContext()
	defer ctx.Release()

	for _, item := range tests {
		r, err := ctx.NewRegExp(item)
		if err != nil {
			t.Errorf("ctx.NewRegExp failed on string %v with error %v", item, err)
		}
		if ctx.ToStringOrDie(r.ToValue()) != "/"+item+"/" {
			t.Errorf("Error compling regexp %s", item)
		}
	}
}

func TestNewRegExpFromValues(t *testing.T) {
	tests := []string{"\\bt[a-z]+\\b", "[0-9]+(\\.[0-9]*)?"}

	ctx := NewContext()
	defer ctx.Release()

	for _, item := range tests {
		params := []*Value{ctx.NewStringValue(item)}
		r, err := ctx.NewRegExpFromValues(params)
		if err != nil {
			t.Errorf("ctx.NewRegExp failed on string %v with error %v", item, err)
		}
		if ctx.ToStringOrDie(r.ToValue()) != "/"+item+"/" {
			t.Errorf("Error compling regexp %s", item)
		}
	}
}

func TestNewFunction(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	fn, err := ctx.NewFunction("myfun", []string{"a", "b"}, "return a+b;", "./testing.go", 1)
	if err != nil {
		t.Errorf("ctx.NewFunction failed with %v", err)
	}
	if !ctx.IsFunction(fn) {
		t.Errorf("ctx.NewFunction did not return a function object")
	}
}

func TestNewCallAsFunction(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	fn, err := ctx.NewFunction("myfun", []string{"a", "b"}, "return a+b;", "./testing.go", 1)
	if err != nil {
		t.Errorf("ctx.NewFunction failed with %v", err)
	}

	a := ctx.NewNumberValue(1.5)
	b := ctx.NewNumberValue(3.0)
	val, err := ctx.CallAsFunction(fn, nil, []*Value{a, b})
	if err != nil {
		t.Errorf("ctx.CallAsFunction failed with %v", err)
	}
	if !ctx.IsNumber(val) {
		t.Errorf("ctx.CallAsFunction did not compute the right value")
	}

	num := ctx.ToNumberOrDie(val)
	if num != 4.5 {
		t.Errorf("ctx.CallAsFunction did not compute the right value")
	}
}
