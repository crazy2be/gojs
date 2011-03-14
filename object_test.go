package javascriptcore_test

import(
	"testing"
	js "javascriptcore"
)

func TestMakeObject(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	val := ctx.MakeObject()
	if val == nil {
		t.Errorf( "ctx.MakeObject returned a nil poitner" )
	}
	if !ctx.IsObject( val.ToValue() ) {
		t.Errorf( "ctx.MakeObject failed to return an object (%v)", ctx.ValueType( val.ToValue() ) )
	}
}

func TestMakeArray(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	val, err := ctx.MakeArray(nil)
	if err != nil {
		t.Errorf( "ctx.MakeArray returned an exception (%v)", ctx.ToStringOrDie(err) )
	}
	if val == nil {
		t.Errorf( "ctx.MakeArray returned a nil poitner" )
	}
	if !ctx.IsObject( val.ToValue() ) {
		t.Errorf( "ctx.MakeArray failed to return an object (%v)", ctx.ValueType( val.ToValue() ) )
	}
}	

func TestMakeArray2(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )

	val, err := ctx.MakeArray( []*js.Value{ a, b } )
	if err != nil {
		t.Errorf( "ctx.MakeArray returned an exception (%v)", ctx.ToStringOrDie(err) )
	}
	if val == nil {
		t.Errorf( "ctx.MakeArray returned a nil poitner" )
	}
	if !ctx.IsObject( val.ToValue() ) {
		t.Errorf( "ctx.MakeArray failed to return an object (%v)", ctx.ValueType( val.ToValue() ) )
	}
	prop, err := ctx.GetProperty( val, "length" )
	if err != nil || prop == nil {
		t.Errorf( "ctx.MakeArray returned object without 'length' property" )
	} else {
		if !ctx.IsNumber( prop ) {
			t.Errorf( "ctx.MakeArray return object with 'length' property not a number" )
		}
		if ctx.ToNumberOrDie( prop ) != 2 {
			t.Errorf( "ctx.MakeArray return object with 'length' not equal to 2", ctx.ToNumberOrDie( prop ) )
		}
	}
}	

func TestMakeDate(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	val, err := ctx.MakeDate()
	if err != nil {
		t.Errorf( "ctx.MakeDate returned an exception (%v)", ctx.ToStringOrDie(err) )
	}
	if val == nil {
		t.Errorf( "ctx.MakeDate returned a nil poitner" )
	}
	if !ctx.IsObject( val.ToValue() ) {
		t.Errorf( "ctx.MakeDate failed to return an object (%v)", ctx.ValueType( val.ToValue() ) )
	}
}	

func TestMakeDateWithMilliseconds(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	val, err := ctx.MakeDateWithMilliseconds( 3600000 )
	if err != nil {
		t.Errorf( "ctx.MakeDateWithMilliseconds returned an exception (%v)", ctx.ToStringOrDie(err) )
	}
	if val == nil {
		t.Errorf( "ctx.MakeDateWithMilliseconds returned a nil poitner" )
	}
	if !ctx.IsObject( val.ToValue() ) {
		t.Errorf( "ctx.MakeDateWithMilliseconds failed to return an object (%v)", ctx.ValueType( val.ToValue() ) )
	}
}	

func TestMakeDateWithString(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	val, err := ctx.MakeDateWithString( "01-Oct-2010" )
	if err != nil {
		t.Errorf( "ctx.MakeDateWithString returned an exception (%v)", ctx.ToStringOrDie(err) )
	}
	if val == nil {
		t.Errorf( "ctx.MakeDateWithString returned a nil poitner" )
	}
	if !ctx.IsObject( val.ToValue() ) {
		t.Errorf( "ctx.MakeDateWithString failed to return an object (%v)", ctx.ValueType( val.ToValue() ) )
	}
}	

func TestMakeError(t *testing.T) {
	tests := []string{ "test error 1", "test error 2" }

	ctx := js.NewContext()
	defer ctx.Release()

	for _, item := range tests {
		r, err := ctx.MakeError( item )
		if err != nil {
			t.Errorf( "ctx.MakeError failed on string %v with error %v", item, err )
		}
		v, exc := ctx.GetProperty( r, "name" )
		if exc != nil || v == nil {
			t.Errorf( "ctx.MakeError returned object without 'message' property" )
		} else {
			if !ctx.IsString( v ) {
				t.Errorf( "ctx.MakeError return object with 'message' property not a string" )
			}
			if ctx.ToStringOrDie( v ) != "Error" {
				t.Errorf( "JavaScript error object and input string don't match (%v, %v)", item, ctx.ToStringOrDie( v ) )
			}
		}
		v, exc = ctx.GetProperty( r, "message" )
		if exc != nil || v == nil {
			t.Errorf( "ctx.MakeError returned object without 'message' property" )
		} else {
			if !ctx.IsString( v ) {
				t.Errorf( "ctx.MakeError return object with 'message' property not a string" )
			}
			if ctx.ToStringOrDie( v ) != item {
				t.Errorf( "JavaScript error object and input string don't match (%v, %v)", item, ctx.ToStringOrDie( v ) )
			}
		}
	}
}

func TestMakeRegExp(t *testing.T) {
	tests := []string{ "\\bt[a-z]+\\b", "[0-9]+(\\.[0-9]*)?" }

	ctx := js.NewContext()
	defer ctx.Release()

	for _, item := range tests {
		r, err := ctx.MakeRegExp( item )
		if err != nil {
			t.Errorf( "ctx.MakeRegExp failed on string %v with error %v", item, err )
		}
		if ctx.ToStringOrDie( r.ToValue() ) != "/" + item + "/" {
			t.Errorf( "Error compling regexp %s", item )
		}
	}
}

func TestMakeRegExpFromValues(t *testing.T) {
	tests := []string{ "\\bt[a-z]+\\b", "[0-9]+(\\.[0-9]*)?" }

	ctx := js.NewContext()
	defer ctx.Release()

	for _, item := range tests {
		params := []*js.Value{ ctx.NewStringValue( item ) }
		r, err := ctx.MakeRegExpFromValues( params )
		if err != nil {
			t.Errorf( "ctx.MakeRegExp failed on string %v with error %v", item, err )
		}
		if ctx.ToStringOrDie( r.ToValue() ) != "/" + item + "/" {
			t.Errorf( "Error compling regexp %s", item )
		}
	}
}

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
	if !ctx.IsNumber( val ) {
		t.Errorf( "ctx.CallAsFunction did not compute the right value" )
	}

	num := ctx.ToNumberOrDie( val )
	if num != 4.5 {
		t.Errorf( "ctx.CallAsFunction did not compute the right value" )
	}
}

