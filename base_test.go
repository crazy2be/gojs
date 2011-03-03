package javascriptcore_test

import(
	"testing"
	js "javascriptcore"
)

type BaseTests struct {
	script string
	valuetype uint8
	result string
}

var(
	basetests = []BaseTests{
		{ "return 2341234 \"asdf\"", js.TypeUndefined, "" },	// syntax error
		{ "1.5", js.TypeNumber, "1.5" },
		{ "1.5 + 3.0", js.TypeNumber, "4.5" },
		{ "'a' + 'b'", js.TypeString, "ab" } }
)

func TestBase(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()
}

func TestEvaluateScript(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	for index, item := range basetests {
		ret, err := ctx.EvaluateScript( item.script, nil, "./testing.go", 1 )
		if item.result != "" {
			if err != nil {
				t.Errorf( "ctx.EvaluateScript raised an error (script %v)", index )
			} else if ret == nil {
				t.Errorf( "ctx.EvaluateScript failed to return a result (script %v)", index )
			}
			t.Logf( "Type of value is %v\n", ctx.ValueType(ret) )
			valuetype := ctx.ValueType(ret)
			if valuetype != item.valuetype {
				t.Errorf( "ctx.EvaluateScript did not return the expected type (%v instead of %v).", valuetype, item.valuetype )
			}
			if ctx.ToStringOrDie(ret) != item.result {
				t.Errorf( "ctx.EvaluateScript returned an incorrect value." )
			}
		} else {
			// Script has a syntax error
			if err == nil {
				t.Errorf( "ctx.EvaluateScript did not raise an error on an invalid script" )
			}
			if ret != nil {
				t.Errorf( "ctx.EvaluateScript returned a result on an invalid script" )
			}
		}
	}
}

func TestCheckScript(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	for _, item := range basetests {
		err := ctx.CheckScriptSyntax( item.script, "./testing.go", 1 )
		if err != nil && item.result!="" {
			t.Errorf( "ctx.CheckScriptSyntax raised an error but script is good" )
		} 
		if err == nil && item.result=="" {
			t.Errorf( "ctx.CheckScriptSyntax failed to raise an error but script is bad" )
		}
	}
}

func TestGarbageCollect(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	ctx.GarbageCollect()
}

