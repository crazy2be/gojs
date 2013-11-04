package gojs

import (
	"testing"
)

func TestNewValueWithNil(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(nil)
	if val.Type() != TypeNull {
		t.Errorf("ctx.ValueType did not return TypeNull")
	}
	if !val.IsNull() {
		t.Errorf("ctx.IsNull did not return true")
	}
}

func TestNewValueWithInt(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(int(4))
	if val.Type() != TypeNumber {
		t.Errorf("ctx.ValueType did not return TypeNumber")
	}
	if !val.IsNumber() {
		t.Errorf("ctx.IsNumber did not return true")
	}
	if val.ToNumberOrDie() != 4 {
		t.Errorf("ctx.ToNumberOrDie did not return correct value")
	}
}

func TestNewValueWithUint(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(uint(4))
	if val.Type() != TypeNumber {
		t.Errorf("ctx.ValueType did not return TypeNumber")
	}
	if !val.IsNumber() {
		t.Errorf("ctx.IsNumber did not return true")
	}
	if val.ToNumberOrDie() != 4 {
		t.Errorf("ctx.ToNumberOrDie did not return correct value")
	}
}

func TestNewValueWithFloat(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(4.5)
	if val.Type() != TypeNumber {
		t.Errorf("ctx.ValueType did not return TypeNumber")
	}
	if !val.IsNumber() {
		t.Errorf("ctx.IsNumber did not return true")
	}
	if val.ToNumberOrDie() != 4.5 {
		t.Errorf("ctx.ToNumberOrDie did not return correct value")
	}
}

func TestNewValueWithString(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue("Some text.")
	if val.Type() != TypeString {
		t.Errorf("ctx.ValueType did not return TypeString")
	}
	if !val.IsString() {
		t.Errorf("ctx.IsString did not return true")
	}
	if val.ToStringOrDie() != "Some text." {
		t.Errorf("ctx.ToStringOrDie did not return correct value")
	}
}

func TestNewValueWithFunc(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	val := ctx.NewValue(func() int { return 1 })
	if val.Type() != TypeObject {
		t.Errorf("ctx.ValueType did not return TypeObject")
	}
	if !val.IsObject() {
		t.Errorf("ctx.IsObject did not return true")
	}
	if !val.ToObjectOrDie().IsFunction() {
		t.Errorf("ctx.IsFunction did not return true")
	}

	val2, err := val.ToObjectOrDie().CallAsFunction(nil, nil)
	if err != nil || val2 == nil {
		t.Errorf("Error executing native function (%v)", err)
	}
	if val2.ToNumberOrDie() != 1 {
		t.Errorf("Native function did not return the correct value")
	}
}

func TestNewValueWithObject(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	obj := &reflect_object{-1, 2, 3.5, "four"}

	val := ctx.NewValue(obj)
	if val.Type() != TypeObject {
		t.Errorf("ctx.ValueType did not return TypeObject")
	}
	if !val.IsObject() {
		t.Errorf("ctx.IsObject did not return true")
	}
}
