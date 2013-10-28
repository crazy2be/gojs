package gojs

import (
	"reflect"
	"testing"
)

func TestValue_GoValue(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	tests := []struct {
		jsValue     *Value
		wantGoValue interface{}
	}{
		{ctx.NewUndefinedValue(), nil},
		{ctx.NewNullValue(), nil},
		{ctx.NewBooleanValue(true), true},
		{ctx.NewBooleanValue(false), false},
		{ctx.NewNumberValue(1.5), 1.5},
		{ctx.NewNumberValue(-3.0), -3.0},
		{ctx.NewStringValue(""), ""},
		{ctx.NewStringValue("foo"), "foo"},
	}

	for _, test := range tests {
		goValue, err := test.jsValue.GoValue()
		if err != nil {
			t.Errorf("JS value %q: GoValue error: %s", test.jsValue, err)
			continue
		}
		if !reflect.DeepEqual(test.wantGoValue, goValue) {
			t.Errorf("JS value %q: want GoValue %+v, got %+v", test.jsValue, test.wantGoValue, goValue)
		}
	}
}
