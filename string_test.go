package javascriptcore_test

import(
	"testing"
	js "javascriptcore"
)

var(
	strtests = []string{ "a string", "unicode \u65e5\u672c\u8a9e" }
)

func TestString(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	str := ctx.NewString( "a string" )
	defer str.Release()
}

func TestString2(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	str := ctx.NewString( "a string" )
	defer str.Release()

	str.Retain()
	str.Release()
}

func TestStringString(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	for _, item := range strtests {
		str := ctx.NewString( item )
		defer str.Release()

		if str.String() != item {
			t.Errorf( "str.String() returned \"%v\", expected \"%v\"", str.String(), item )
		}
		if str.Length() != uint32(len( []int(item))) {
			t.Errorf( "str.Length() returned \"%v\", expected \"%v\"", str.Length(), len( []int(item)) )
		}
	}
}

func TestStringEqual(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	lhs := ctx.NewString( "dummy" )
	defer lhs.Release()

	for _, item := range strtests {
		str := ctx.NewString( item )
		defer str.Release()

		if lhs.Equal( str ) {
			t.Errorf( "Strings compared as equal \"%v\", and \"%v\"", lhs, str )
		}
		if str.Equal( lhs ) {
			t.Errorf( "Strings compared as equal \"%v\", and \"%v\"", str, lhs )
		}
		if !str.Equal( str ) {
			t.Errorf( "String did not compared as equal to itself \"%v\", and \"%v\"", str )
		}

		str2 := ctx.NewString( item )
		defer str2.Release()
		if !str.Equal( str2 ) {
			t.Errorf( "String did not compared as equal to itself \"%v\", and \"%v\"", str2 )
		}
	}
}

func TestStringEqualToString(t *testing.T) {
	ctx := js.NewContext()
	defer ctx.Release()

	lhs := ctx.NewString( "dummy" )
	defer lhs.Release()

	for _, item := range strtests {
		str := ctx.NewString( item )
		defer str.Release()

		if lhs.EqualToString( item ) {
			t.Errorf( "Strings compared as equal \"%v\", and \"%v\"", lhs, item )
		}
		if !str.EqualToString( item ) {
			t.Errorf( "String did not compare as equal to itself \"%v\"", item )
		}
	}
}

