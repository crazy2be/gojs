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
		{ctx.NewEmptyObject().ToValue(), map[string]interface{}{}},
		// TODO(sqs): convert JavaScript Date objects to Go time.Time (currently
		// they are converted to a Go string)
		// {jsObjectToJSValue(ctx.NewDateWithMilliseconds(123)), time.Unix(0, 123000000).In(time.UTC)},
		{
			jsObjectToJSValue(ctx.NewArray([]*Value{ctx.NewStringValue("foo")})),
			[]interface{}{"foo"},
		},
		{
			jsObjectToJSValue(ctx.NewObjectWithProperties(map[string]*Value{"foo": ctx.NewStringValue("bar")})),
			map[string]interface{}{"foo": "bar"},
		},
	}

	for _, test := range tests {
		goValue, err := test.jsValue.GoValue()
		if err != nil {
			t.Errorf("JS value %q: GoValue error: %s", test.jsValue, err)
			continue
		}
		if !reflect.DeepEqual(test.wantGoValue, goValue) {
			t.Errorf("JS value %q: want GoValue %+v (type %T), got %+v (type %T)", test.jsValue, test.wantGoValue, test.wantGoValue, goValue, goValue)
		}
	}
}

func jsObjectToJSValue(obj *Object, err error) *Value {
	if err != nil {
		panic("object creation failed: " + err.Error())
	}
	return obj.ToValue()
}

func TestNewValueFrom(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	wantString := "foo"
	jsval := ctx.NewStringValue(wantString)
	rawval := RawValue(jsval.ref)
	jsval2 := ctx.NewValueFrom(rawval)

	if gotString := ctx.ToStringOrDie(jsval2); wantString != gotString {
		t.Errorf("want string %q, got %q", wantString, gotString)
	}
}
