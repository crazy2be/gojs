package gojs

import (
	"testing"
)

func catch(fn func()) (ret interface{}) {
	defer func() {
		if r := recover(); r != nil {
			ret = r
		}
	}()
	fn()
	return nil
}

func TestPanicNumber(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	a := ctx.NewNumberValue(1.5)

	r := catch(func() { panic(newPanicError(ctx, a)) })
	if r == nil {
		t.Errorf("Callback did not panic")
		return
	}
	err, ok := r.(*Error)
	if !ok {
		t.Errorf("Type conversion to *Error failed")
	} else {
		if err.Context != ctx {
			t.Errorf("err.Context not set correctly")
		}
		if err.Value != a {
			t.Errorf("err.Value not set correctly")
		}
		if err.Name != "Error" {
			t.Errorf("err.Name not set correctly")
		}
		if err.Message != "1.5" {
			t.Errorf("err.Message not set correctly")
		}
	}
}

func TestPanicString(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	a := ctx.NewStringValue("my custom string")

	r := catch(func() { panic(newPanicError(ctx, a)) })
	if r == nil {
		t.Errorf("Callback did not panic")
		return
	}
	err, ok := r.(*Error)
	if !ok {
		t.Errorf("Type conversion to *Error failed")
	} else {
		if err.Context != ctx {
			t.Errorf("err.Context not set correctly")
		}
		if err.Value != a {
			t.Errorf("err.Value not set correctly")
		}
		if err.Name != "Error" {
			t.Errorf("err.Name not set correctly")
		}
		if err.Message != "my custom string" {
			t.Errorf("err.Message not set correctly")
		}
	}
}

func TestPanicError(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	a, _ := ctx.NewError("my custom string")

	r := catch(func() { panic(newPanicError(ctx, a.ToValue())) })
	if r == nil {
		t.Errorf("Callback did not panic")
		return
	}
	err, ok := r.(*Error)
	if !ok {
		t.Errorf("Type conversion to *Error failed")
	} else {
		if err.Context != ctx {
			t.Errorf("err.Context not set correctly")
		}
		if err.Value != a.ToValue() {
			t.Errorf("err.Value not set correctly")
		}
		if err.Name != "Error" {
			t.Errorf("err.Name not set correctly")
		}
		if err.Message != "my custom string" {
			t.Errorf("err.Message not set correctly")
		}
	}
}
