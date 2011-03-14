package javascriptcore_test

import(
	"testing"
	js "javascriptcore"
	"os"
)

type reflect_object struct {
	I	int
	F	float64
	S	string
}

func TestMakeFunctionWithCallback(t *testing.T) {
	var flag bool
	callback := func (ctx *js.Context, obj *js.Object, thisObject *js.Object, _ []*js.Value ) (*js.Value){
		flag = true
		return nil
	}

	ctx := js.NewContext()
	defer ctx.Release()

	fn := ctx.MakeFunctionWithCallback( callback )
	if fn == nil {
		t.Errorf( "ctx.MakeFunctionWithCallback failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.MakeFunctionWithCallback returned value that is not a function" )
	}
	if ctx.ToStringOrDie( fn.ToValue() ) != "nativecallback" {
		t.Errorf( "ctx.MakeFunctionWithCallback returned value that does not convert to property string" )
	}
	ctx.CallAsFunction( fn, nil, []*js.Value{} )
	if !flag {
		t.Errorf( "Native function did not execute" )
	}
}

func TestMakeFunctionWithCallback2(t *testing.T) {
	callback := func (ctx *js.Context, obj *js.Object, thisObject *js.Object, args []*js.Value ) (*js.Value){
		if len(args)!=2 {
			return nil
		}

		a := ctx.ToNumberOrDie( args[0] )
		b := ctx.ToNumberOrDie( args[1] )
		return ctx.NewNumberValue( a + b )
	}

	ctx := js.NewContext()
	defer ctx.Release()

	fn := ctx.MakeFunctionWithCallback( callback )
	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )
	val, err := ctx.CallAsFunction( fn, nil, []*js.Value{ a, b } )
	if err != nil || val == nil {
		t.Errorf( "Error executing native callback" )
	}
	if ctx.ToNumberOrDie(val)!=4.5 {
		t.Errorf( "Native callback did not return the correct value" )
	}
}

func TestMakeFunctionWithCallbackPanic(t *testing.T) {
	var callbacks = []js.GoFunctionCallback{}
	var error_msgs = []string{ "error from go!", os.ENOMEM.String() }

	callbacks = append( callbacks,
		func (ctx *js.Context, obj *js.Object, thisObject *js.Object, _ []*js.Value ) (*js.Value,) {
			panic( "error from go!" )
			return nil } )
	callbacks = append( callbacks,
		func (ctx *js.Context, obj *js.Object, thisObject *js.Object, _ []*js.Value ) (*js.Value,) {
			panic( os.ENOMEM )
			return nil } )

	ctx := js.NewContext()
	defer ctx.Release()

	for index, callback := range callbacks {

	fn := ctx.MakeFunctionWithCallback( callback )
	if fn == nil {
		t.Errorf( "ctx.MakeFunctionWithCallback failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.MakeFunctionWithCallback returned value that is not a function" )
	}
	if ctx.ToStringOrDie( fn.ToValue() ) != "nativecallback" {
		t.Errorf( "ctx.MakeFunctionWithCallback returned value that does not convert to property string" )
	}
	val, err := ctx.CallAsFunction( fn, nil, []*js.Value{} )
	if val != nil {
		t.Errorf( "ctx.MakeFunctionWithCallback that panicked returned a value" )
	}
	if err == nil || !ctx.IsObject( err ) {
		t.Errorf( "ctx.MakeFunctionWithCallback that panicked did not set exception" )
	}
	if ctx.ToStringOrDie(err) != "Error: " + error_msgs[index] {
		t.Errorf( "ctx.MakeFunctionWithCallback that panicked did not set exception message (%v,%v)", 
			ctx.ToStringOrDie(err), error_msgs[index] )
	}

	} // for
}

func TestNativeFunction(t *testing.T) {
	var flag bool
	callback := func () {
		flag = true
	}

	ctx := js.NewContext()
	defer ctx.Release()

	fn := ctx.MakeFunctionWithNative( callback )
	if fn == nil {
		t.Errorf( "ctx.MakeFunctionWithNative failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.MakeFunctionWithNative returned value that is not a function" )
	}
	if ctx.ToStringOrDie( fn.ToValue() ) != "nativefunction" {
		t.Errorf( "ctx.nativefunction returned value that does not convert to property string" )
	}
	ctx.CallAsFunction( fn, nil, []*js.Value{} )
	if !flag {
		t.Errorf( "Native function did not execute" )
	}
}	

func TestNativeFunction2(t *testing.T) {
	callback := func ( a float64, b float64 ) float64 {
		return a + float64(b)
	}

	ctx := js.NewContext()
	defer ctx.Release()

	fn := ctx.MakeFunctionWithNative( callback )
	if fn == nil {
		t.Errorf( "ctx.MakeFunctionWithNative failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.MakeFunctionWithNative returned value that is not a function" )
	}
	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )
	val, err := ctx.CallAsFunction( fn, nil, []*js.Value{ a, b } )
	if err != nil || val == nil {
		t.Errorf( "Error executing native function (%v)", ctx.ToStringOrDie(err) )
	}
	if ctx.ToNumberOrDie(val)!=4.5 {
		t.Errorf( "Native function did not return the correct value" )
	}
}	

func TestMakeNativeObject(t *testing.T) {
	obj := &reflect_object{ 2, 3.0, "four" }

	ctx := js.NewContext()
	defer ctx.Release()

	v := ctx.MakeNativeObject( obj )
	ctx.SetProperty( ctx.GlobalObject(), "n", v.ToValue(), 0 )

	// Following script access should be successful
	ret, err := ctx.EvaluateScript( "n.F", nil, "./testing.go", 1 )
	if err != nil || ret == nil {
		t.Errorf( "ctx.EvaluateScript returned an error (or did not return a result)" )
		return
	}
	if !ctx.IsNumber( ret ) {
		t.Errorf( "ctx.EvaluateScript did not return 'number' result when accessing native object's non-existent field." )
	}
	num := ctx.ToNumberOrDie( ret )
	if num != 3.0 {
		t.Errorf( "ctx.EvaluateScript incorrect value when accessing native object's field." )
	}

	// following script access should fail
	ret, err = ctx.EvaluateScript( "n.noexist", nil, "./testing.go", 1 )
	if err != nil || ret == nil {
		t.Errorf( "ctx.EvaluateScript returned an error (or did not return a result)" )
	}
	if !ctx.IsUndefined( ret ) {
		t.Errorf( "ctx.EvaluateScript did not return 'undefined' result when accessing native object's non-existent field." )
	}

	// following script access should succeed
	ret, err = ctx.EvaluateScript( "n.S", nil, "./testing.go", 1 )
	if err != nil || ret == nil {
		t.Errorf( "ctx.EvaluateScript returned an error (or did not return a result)" )
	}
	if !ctx.IsString( ret ) {
		t.Errorf( "ctx.EvaluateScript did not return 'string' result when accessing native object's non-existent field." )
	}
	str := ctx.ToStringOrDie( ret )
	if str != "four" {
		t.Errorf( "ctx.EvaluateScript incorrect value when accessing native object's field." )
	}
}

