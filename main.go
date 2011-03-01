package main

import js "javascriptcore"
import "fmt"
import "os"

const source_url = "./test.js"

func print_properties( ctx *js.Context, tab_count int, value *js.Object ) {
	names := ctx.CopyPropertyNames( value )
	for lp := uint16(0); lp < names.Count(); lp++ {
		name := names.NameAtIndex(lp)
		value, _ := ctx.ObjectGetProperty( value, name )
		fmt.Printf( "%s = ", name )
		print_value_ref( ctx, value )
	}
}

func print_value_ref( ctx *js.Context, value *js.Value ) {
	switch t := ctx.ValueType( value ); true {
		case t == js.TypeUndefined:
			fmt.Printf( "Undefined\n" )
		case t == js.TypeNull:
			fmt.Printf( "Null\n" )
		case t == js.TypeBoolean:
			fmt.Printf( "%v\n", ctx.ToBoolean( value ) )
		case t == js.TypeNumber:
			v, _ := ctx.ToNumber( value )
			fmt.Printf( "%v\n", v )
		case t == js.TypeString:
			v, _ := ctx.ToString( value )
			fmt.Printf( "%v\n", v )
		case t == js.TypeObject:
			fmt.Printf( "{\n" )
			print_properties( ctx, 1, ctx.ToObjectOrDie(value) )
			fmt.Printf( "}\n" )
		default:
			panic( os.EEXIST )
	}
}

func print_result( ctx *js.Context, script string ) {
	err := ctx.CheckScriptSyntax( script, source_url, 1 )
	if err!=nil {
		fmt.Printf( "Syntax Error:\n" )
		print_value_ref( ctx, err )
	} else {
		result, err := ctx.EvaluateScript( script, nil, source_url, 1 )
		if err!=nil {
			fmt.Printf( "Runtime Error:\n" )
			print_value_ref( ctx, err )
		} else {
			print_value_ref( ctx, result )
		}
	}
}

type dummy struct {
}

func (d dummy) Callback( ctx *js.Context, obj *js.Object, thisObject *js.Object, exception **js.Value ) *js.Value {
	fmt.Printf( "In callback!\n" )
	*exception = nil
	return nil
}

func main() {
	ctx := js.NewContext()
	defer ctx.Release()

	s := ctx.NewString( "Hello from go!" )
	defer s.Release()
	fmt.Printf( "%v %v\n", s.Length(), s.String() )
	fmt.Printf( "%v %v\n", s.EqualToString("Hello"), s.EqualToString("Hello from go!") )

	//obj := ctx.MakeFunction( "f", dummy{} )
	//ctx.ObjectSetProperty( ctx.GlobalObject(), "f", obj.GetValue(), js.PropertyAttributeReadOnly )
	//_, err := ctx.EvaluateScript( "f()", nil, "", 1 )
	//if err!=nil {
	//	panic(err)
	//}

	ctx.EvaluateScript( "var a = \"Go!\"", nil, "", 1 )
	a, err := ctx.ObjectGetProperty( ctx.GlobalObject(), "a" )
	fmt.Printf( "%v %s %v\n", a, ctx.ToStringOrDie(a), err )

	fmt.Printf("\nScripts...\n" )
	print_result( ctx, "null" )
	print_result( ctx, "false" )
	print_result( ctx, "1234.123" )
	print_result( ctx, "new Array" )	
	print_result( ctx, "12+34\nreturn new 234 Array" )
	print_result( ctx, "1/0" )

	fmt.Printf( "Done!\n" )
}

