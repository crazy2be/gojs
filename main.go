package main

import js "javascriptcore"
import "fmt"
import "os"

func print_value_ref( ctx *js.Context, value *js.ValueRef ) {
	switch t := ctx.ValueType( value ); true {
		case t == js.TypeUndefined:
			fmt.Printf( "Undefined.\n" )
		case t == js.TypeNull:
			fmt.Printf( "Null.\n" )
		case t == js.TypeBoolean:
			fmt.Printf( "%v\n", ctx.ToBoolean( value ) )
		case t == js.TypeNumber:
			v, _ := ctx.ToNumber( value )
			fmt.Printf( "%v\n", v )
		case t == js.TypeString:
			v, _ := ctx.ToString( value )
			fmt.Printf( "%v\n", v )
		case t == js.TypeObject:
			fmt.Printf( "{}\n" )
		default:
			panic( os.EEXIST )
	}
}

func print_result( ctx *js.Context, script string ) {
	err := ctx.CheckScriptSyntax( script, "", 1 )
	if err!=nil {
		fmt.Printf( "Syntax Error:\n" )
		print_value_ref( ctx, err )
	} else {
		result, err := ctx.EvaluateScript( script, nil, "", 1 )
		if err!=nil {
			fmt.Printf( "Runtime Error:\n" )
			print_value_ref( ctx, err )
		} else {
			print_value_ref( ctx, result )
		}
	}
}

func main() {
	ctx := js.NewContext()
	defer ctx.Release()

	s := ctx.NewString( "Hello from go!" )
	defer s.Release()
	fmt.Printf( "%v %v\n", s.Length(), s.String() )
	fmt.Printf( "%v %v\n", s.EqualToString("Hello"), s.EqualToString("Hello from go!") )

	fmt.Printf("\nScripts...\n" )
	print_result( ctx, "null" )
	print_result( ctx, "false" )
	print_result( ctx, "1234.123" )
	print_result( ctx, "new Array" )	
	print_result( ctx, "return new 234 Array" )

	fmt.Printf( "Done!\n" )
}

