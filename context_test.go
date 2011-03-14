package javascriptcore_test

import(
	"testing"
	js "javascriptcore"
)

func TestContext(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()
}

func TestContext2(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	ctx.Retain()
	defer ctx.Release()
}

func TestContextGlobalObject(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	obj := ctx.GlobalObject()
	if obj == nil {
		t.Errorf( "ctx.GlobalObject() returned nil" )
	}
	if ctx.ValueType(obj.ToValue()) != js.TypeObject {
		t.Errorf( "ctx.GlobalObject() did not return a javascript object" )
	}
}

