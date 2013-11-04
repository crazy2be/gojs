package main

import (
	"github.com/crazy2be/gojs"
	"fmt"
)

const source_url = "./test.js"

func print_properties(ctx *gojs.Context, tab_count int, value *gojs.Object) {
	names := value.CopyPropertyNames()
	for lp := uint16(0); lp < names.Count(); lp++ {
		name := names.NameAtIndex(lp)
		value, _ := value.GetProperty(name)
		fmt.Printf("%s = ", name)
		print_value_ref(ctx, value)
	}
}

func print_value_ref(ctx *gojs.Context, value *gojs.Value) {
	switch t := value.Type(); true {
	case t == gojs.TypeUndefined:
		fmt.Printf("Undefined\n")
	case t == gojs.TypeNull:
		fmt.Printf("Null\n")
	case t == gojs.TypeBoolean:
		fmt.Printf("%v\n", value.ToBoolean())
	case t == gojs.TypeNumber:
		v, _ := value.ToNumber()
		fmt.Printf("%v\n", v)
	case t == gojs.TypeString:
		v, _ := value.ToString()
		fmt.Printf("%v\n", v)
	case t == gojs.TypeObject:
		fmt.Printf("{\n")
		print_properties(ctx, 1, value.ToObjectOrDie())
		fmt.Printf("}\n")
	default:
		panic(fmt.Sprintf("Unknown type for value %v", value))
	}
}

func print_result(ctx *gojs.Context, script string) {
	err := ctx.CheckScriptSyntax(script, source_url, 1)
	if err != nil {
		fmt.Printf("Syntax Error: %s\n", err)
	} else {
		result, err := ctx.EvaluateScript(script, nil, source_url, 1)
		if err != nil {
			fmt.Printf("Runtime Error: %s\n", err)
		} else {
			print_value_ref(ctx, result)
		}
	}
}

func callback(ctx *gojs.Context, obj *gojs.Object, thisObject *gojs.Object, arguments []*gojs.Value) *gojs.Value {
	fmt.Printf("In callback!\n")
	return nil
}

func main() {
	ctx := gojs.NewContext()
	defer ctx.Release()

	s := gojs.NewString("Hello from go!")
	defer s.Release()
	fmt.Printf("%v %v\n", s.Length(), s.String())
	fmt.Printf("%v %v\n", s.EqualToString("Hello"), s.EqualToString("Hello from go!"))

	obj := ctx.NewFunctionWithCallback(gojs.GoFunctionCallback(callback))
	ctx.GlobalObject().SetProperty("f", obj.ToValue(), gojs.PropertyAttributeReadOnly)
	_, err := ctx.EvaluateScript("f()", nil, "", 1)
	if err!=nil {
		panic(err)
	}

	ctx.EvaluateScript("var a = \"Go!\"", nil, "", 1)
	a, err := ctx.GlobalObject().GetProperty("a")
	fmt.Printf("%v %s %v\n", a, a.ToStringOrDie(), err)

	fmt.Printf("\nScripts...\n")
	print_result(ctx, "null")
	print_result(ctx, "false")
	print_result(ctx, "1234.123")
	print_result(ctx, "new Array(1, 2, 3)")
	print_result(ctx, "12+34\nreturn new 234 Array")
	print_result(ctx, "1/0")

	fmt.Printf("Done!\n")
}
