package javascriptcore

import(
	"testing"
	"os"
)

type reflect_object struct {
	I	int
	U	uint
	F	float64
	S	string
}

func (o *reflect_object) String() string {
	return o.S
}

func (o *reflect_object) Add() float64 {
	return float64(o.I) + o.F
}

func (o *reflect_object) AddWith( op float64 ) float64 {
	return float64(o.I) + o.F + op
}

func TestNewFunctionWithCallback(t *testing.T) {
	var flag bool
	callback := func (ctx *Context, obj *Object, thisObject *Object, _ []*Value ) (*Value){
		flag = true
		return nil
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithCallback( callback )
	if fn == nil {
		t.Errorf( "ctx.NewFunctionWithCallback failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.NewFunctionWithCallback returned value that is not a function" )
	}
	if ctx.ToStringOrDie( fn.ToValue() ) != "nativecallback" {
		t.Errorf( "ctx.NewFunctionWithCallback returned value that does not convert to property string" )
	}
	ctx.CallAsFunction( fn, nil, []*Value{} )
	if !flag {
		t.Errorf( "Native function did not execute" )
	}
}

func TestNewFunctionWithCallback2(t *testing.T) {
	callback := func (ctx *Context, obj *Object, thisObject *Object, args []*Value ) (*Value){
		if len(args)!=2 {
			return nil
		}

		a := ctx.ToNumberOrDie( args[0] )
		b := ctx.ToNumberOrDie( args[1] )
		return ctx.NewNumberValue( a + b )
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithCallback( callback )
	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )
	val, err := ctx.CallAsFunction( fn, nil, []*Value{ a, b } )
	if err != nil || val == nil {
		t.Errorf( "Error executing native callback" )
	}
	if ctx.ToNumberOrDie(val)!=4.5 {
		t.Errorf( "Native callback did not return the correct value" )
	}
}

func TestNewFunctionWithCallbackPanic(t *testing.T) {
	var callbacks = []GoFunctionCallback{}
	var error_msgs = []string{ "error from go!", os.ENOMEM.String() }

	callbacks = append( callbacks,
		func (ctx *Context, obj *Object, thisObject *Object, _ []*Value ) (*Value,) {
			panic( "error from go!" )
			return nil } )
	callbacks = append( callbacks,
		func (ctx *Context, obj *Object, thisObject *Object, _ []*Value ) (*Value,) {
			panic( os.ENOMEM )
			return nil } )

	ctx := NewContext()
	defer ctx.Release()

	for index, callback := range callbacks {

	fn := ctx.NewFunctionWithCallback( callback )
	if fn == nil {
		t.Errorf( "ctx.NewFunctionWithCallback failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.NewFunctionWithCallback returned value that is not a function" )
	}
	if ctx.ToStringOrDie( fn.ToValue() ) != "nativecallback" {
		t.Errorf( "ctx.NewFunctionWithCallback returned value that does not convert to property string" )
	}
	val, err := ctx.CallAsFunction( fn, nil, []*Value{} )
	if val != nil {
		t.Errorf( "ctx.NewFunctionWithCallback that panicked returned a value" )
	}
	if err == nil || !ctx.IsObject( err ) {
		t.Errorf( "ctx.NewFunctionWithCallback that panicked did not set exception" )
	}
	if ctx.ToStringOrDie(err) != "Error: " + error_msgs[index] {
		t.Errorf( "ctx.NewFunctionWithCallback that panicked did not set exception message (%v,%v)", 
			ctx.ToStringOrDie(err), error_msgs[index] )
	}

	} // for
}

func TestNativeFunction(t *testing.T) {
	var flag bool
	callback := func () {
		flag = true
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithNative( callback )
	if fn == nil {
		t.Errorf( "ctx.NewFunctionWithNative failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.NewFunctionWithNative returned value that is not a function" )
	}
	if ctx.ToStringOrDie( fn.ToValue() ) != "nativefunction" {
		t.Errorf( "ctx.nativefunction returned value that does not convert to property string" )
	}
	ctx.CallAsFunction( fn, nil, []*Value{} )
	if !flag {
		t.Errorf( "Native function did not execute" )
	}
}	

func TestNativeFunction2(t *testing.T) {
	callback := func ( a float64, b float64 ) float64 {
		return a + float64(b)
	}

	ctx := NewContext()
	defer ctx.Release()

	fn := ctx.NewFunctionWithNative( callback )
	if fn == nil {
		t.Errorf( "ctx.NewFunctionWithNative failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.NewFunctionWithNative returned value that is not a function" )
	}
	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )
	val, err := ctx.CallAsFunction( fn, nil, []*Value{ a, b } )
	if err != nil || val == nil {
		t.Errorf( "Error executing native function (%v)", ctx.ToStringOrDie(err) )
	}
	if ctx.ToNumberOrDie(val)!=4.5 {
		t.Errorf( "Native function did not return the correct value" )
	}
}	

func TestNativeFunction3(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	callback := func ( a float64, b float64 ) *Value {
		ret := a + float64(b)
		return ctx.NewNumberValue( ret )
	}

	fn := ctx.NewFunctionWithNative( callback )
	if fn == nil {
		t.Errorf( "ctx.NewFunctionWithNative failed" )
		return
	}
	if !ctx.IsFunction( fn ) {
		t.Errorf( "ctx.NewFunctionWithNative returned value that is not a function" )
	}
	a := ctx.NewNumberValue( 1.5 )
	b := ctx.NewNumberValue( 3.0 )
	val, err := ctx.CallAsFunction( fn, nil, []*Value{ a, b } )
	if err != nil || val == nil {
		t.Errorf( "Error executing native function (%v)", ctx.ToStringOrDie(err) )
	}
	if ctx.ToNumberOrDie(val)!=4.5 {
		t.Errorf( "Native function did not return the correct value" )
	}
}	

func TestNewNativeObject(t *testing.T) {
	obj := &reflect_object{ -1, 2, 3.0, "four" }

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject( obj )
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

func TestNewNativeObjectSet(t *testing.T) {
	obj := &reflect_object{ -1, 2, 3.0, "four" }

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject( obj )
	ctx.SetProperty( ctx.GlobalObject(), "n", v.ToValue(), 0 )

	// Set the integer property
	i := ctx.NewNumberValue( -2 )
	ctx.SetProperty( v, "I", i, 0 )
	if obj.I != -2 {
		t.Errorf( "ctx.SetProperty did not set integer field correctly" )
	}

	// Set the unsigned integer property
	u := ctx.NewNumberValue( 3 )
	ctx.SetProperty( v, "U", u, 0 )
	if obj.U != 3 {
		t.Errorf( "ctx.SetProperty did not set unsigned integer field correctly" )
	}

	// Set the unsigned integer property
	u = ctx.NewNumberValue( -3 )
	err := ctx.SetProperty( v, "U", u, 0 )
	if err == nil {
		t.Errorf( "ctx.SetProperty did not set unsigned integer field correctly" )
	} else {
		t.Logf( "%v", ctx.ToStringOrDie( err ) )
	}
	if obj.U != 3 {
		t.Errorf( "ctx.SetProperty did not set unsigned integer field correctly" )
	}

	// Set the float property
	n := ctx.NewNumberValue( 4.0 )
	ctx.SetProperty( v, "F", n, 0 )
	if obj.F != 4.0 {
		t.Errorf( "ctx.SetProperty did not set float field correctly" )
	}

	s := ctx.NewStringValue( "five" )
	ctx.SetProperty( v, "S", s, 0 )
	if obj.S != "five" {
		t.Errorf( "ctx.SetProperty did not set string field correctly" )
	}
}

func TestNewNativeObjectConvert(t *testing.T) {
	obj := &reflect_object{ -1, 2, 3.0, "four" }

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject( obj )

	if ctx.ToStringOrDie( v.ToValue() ) != "four" {
		t.Errorf( "ctx.ToStringOrDie for native object did not return correct value." )
	}
}

func TestNewNativeObjectMethod(t *testing.T) {
	obj := &reflect_object{ -1, 2, 3.0, "four" }

	ctx := NewContext()
	defer ctx.Release()

	v := ctx.NewNativeObject( obj )
	ctx.SetProperty( ctx.GlobalObject(), "n", v.ToValue(), 0 )

	// Following script access should be successful
	ret, err := ctx.EvaluateScript( "n.Add()", nil, "./testing.go", 1 )
	if err != nil || ret == nil {
		t.Errorf( "ctx.EvaluateScript returned an error (or did not return a result)" )
		return
	}
	if !ctx.IsNumber( ret ) {
		t.Errorf( "ctx.EvaluateScript did not return 'number' result when calling method 'Add'." )
	}
	num := ctx.ToNumberOrDie( ret )
	if num != 2.0 {
		t.Errorf( "ctx.EvaluateScript incorrect value when accessing native object's field." )
	}

	// Following script access should be successful
	ret, err = ctx.EvaluateScript( "n.AddWith(0.5)", nil, "./testing.go", 1 )
	if err != nil || ret == nil {
		t.Errorf( "ctx.EvaluateScript returned an error (or did not return a result)" )
		return
	}
	if !ctx.IsNumber( ret ) {
		t.Errorf( "ctx.EvaluateScript did not return 'number' result when calling method 'AddWith'." )
	}
	num = ctx.ToNumberOrDie( ret )
	if num != 2.5 {
		t.Errorf( "ctx.EvaluateScript incorrect value when accessing native object's field." )
	}
}

