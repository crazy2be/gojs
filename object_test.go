package javascriptcore_test

import(
	"testing"
	js "javascriptcore"
)

func TestMakeFunction(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	fn, err := ctx.MakeFunction( "myfun", []string{ "a", "b" }, "return a+b;", "./testing.go", 1 )
	if err != nil {
		t.Errorf( "ctx.MakeFunction failed with %v", err )
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.MakeFunction did not return a function object" )
	}
}

func TestMakeCallAsFunction(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	fn, err := ctx.MakeFunction( "myfun", []string{ "a", "b" }, "return a+b;", "./testing.go", 1 )
	if err != nil {
		t.Errorf( "ctx.MakeFunction failed with %v", err )
	}
	
	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )
	val, err := ctx.CallAsFunction( fn, nil, []*js.Value{ a, b } )
	if err != nil {
		t.Errorf( "ctx.CallAsFunction failed with %v", err )
	}
	if !ctx.ValueIsNumber( val ) {
		t.Errorf( "ctx.CallAsFunction did not compute the right value" )
	}

	num := ctx.ToNumberOrDie( val )
	if num != 4.5 {
		t.Errorf( "ctx.CallAsFunction did not compute the right value" )
	}
}

