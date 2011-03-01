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
}

func TestEvaluateScript(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	ret, err := ctx.EvaluateScript( "1.5", nil, "./testing.go", 1 )
	if err != nil {
		t.Errorf( "ctx.EvaluateScript raised an error." )
		return
	}
	if ret == nil {
		t.Errorf( "ctx.EvaluateScript failed to return a result." )
	}
	t.Logf( "Type of value is %v\n", ctx.ValueType(ret) )
	num, err := ctx.ToNumber( ret )
	if err != nil {
		t.Errorf( "ctx.EvaluateScript did not return a number as expected." )
		return
	}
	if num != 1.5 {
		t.Errorf( "ctx.EvaluateScript returned an incorrect number." )
	}
}

