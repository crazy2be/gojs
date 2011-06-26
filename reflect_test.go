package gojs

import (
	"testing"
)

func TestNewValueWithNil(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(nil)
	if ctx.ValueType(val) != TypeNull {
		t.Errorf("ctx.ValueType did not return TypeNull")
	}
	if !ctx.IsNull(val) {
		t.Errorf("ctx.IsNull did not return true")
	}
}

func TestNewValueWithInt(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(int(4))
	if ctx.ValueType(val) != TypeNumber {
		t.Errorf("ctx.ValueType did not return TypeNumber")
	}
	if !ctx.IsNumber(val) {
		t.Errorf("ctx.IsNumber did not return true")
	}
	if ctx.ToNumberOrDie(val) != 4 {
		t.Errorf("ctx.ToNumberOrDie did not return correct value")
	}
}

func TestNewValueWithUint(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(uint(4))
	if ctx.ValueType(val) != TypeNumber {
		t.Errorf("ctx.ValueType did not return TypeNumber")
	}
	if !ctx.IsNumber(val) {
		t.Errorf("ctx.IsNumber did not return true")
	}
	if ctx.ToNumberOrDie(val) != 4 {
		t.Errorf("ctx.ToNumberOrDie did not return correct value")
	}
}

func TestNewValueWithFloat(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(4.5)
	if ctx.ValueType(val) != TypeNumber {
		t.Errorf("ctx.ValueType did not return TypeNumber")
	}
	if !ctx.IsNumber(val) {
		t.Errorf("ctx.IsNumber did not return true")
	}
	if ctx.ToNumberOrDie(val) != 4.5 {
		t.Errorf("ctx.ToNumberOrDie did not return correct value")
	}
}

func TestNewValueWithString(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue("Some text.")
	if ctx.ValueType(val) != TypeString {
		t.Errorf("ctx.ValueType did not return TypeString")
	}
	if !ctx.IsString(val) {
		t.Errorf("ctx.IsString did not return true")
	}
	if ctx.ToStringOrDie(val) != "Some text." {
		t.Errorf("ctx.ToStringOrDie did not return correct value")
	}
}

func TestNewValueWithFunc(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(func() int { return 1 })
	if ctx.ValueType(val) != TypeObject {
		t.Errorf("ctx.ValueType did not return TypeObject")
	}
	if !ctx.IsObject(val) {
		t.Errorf("ctx.IsObject did not return true")
	}
	if !ctx.IsFunction(ctx.ToObjectOrDie(val)) {
		t.Errorf("ctx.IsFunction did not return true")
	}

	val2, err := ctx.CallAsFunction(ctx.ToObjectOrDie(val), nil, nil)
	if err != nil || val2 == nil {
		t.Errorf("Error executing native function (%v)", err)
	}
	if ctx.ToNumberOrDie(val2) != 1 {
		t.Errorf("Native function did not return the correct value")
	}
}

func TestNewValueWithObject(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	obj := &reflect_object{-1, 2, 3.5, "four"}

	val := ctx.NewValue(obj)
	if ctx.ValueType(val) != TypeObject {
		t.Errorf("ctx.ValueType did not return TypeObject")
	}
	if !ctx.IsObject(val) {
		t.Errorf("ctx.IsObject did not return true")
	}
}
