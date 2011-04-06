package javascriptcore

// #include <stdlib.h>
import "C"
import "reflect"

func (ctx *Context) NewValue( value interface{} ) *Value {
	// Handle simple case right off
	if value==nil {
		return ctx.NewNullValue()
	}

	ret := value_to_javascript( ctx, reflect.NewValue( value ) )
	return ret
}

