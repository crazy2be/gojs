package gojs

import (
	"bytes"
	"testing"
)

var (
	strtests = []string{"a string", "unicode \u65e5\u672c\u8a9e"}
)

func TestString(t *testing.T) {
	str := NewString("a string")
	defer str.Release()
}

func TestString2(t *testing.T) {
	str := NewString("a string")
	defer str.Release()

	str.Retain()
	str.Release()
}

func TestStringString(t *testing.T) {
	for _, item := range strtests {
		str := NewString(item)
		defer str.Release()

		if str.String() != item {
			t.Errorf("str.String() returned \"%v\", expected \"%v\"", str.String(), item)
		}
		if str.Length() != uint32(len([]rune(item))) {
			t.Errorf("str.Length() returned \"%v\", expected \"%v\"", str.Length(), len([]rune(item)))
		}
	}
}

func TestStringBytes(t *testing.T) {
	for _, item := range strtests {
		str := NewString(item)
		defer str.Release()

		wantBytes := []byte(item)
		gotBytes := str.Bytes()
		if !bytes.Equal(wantBytes, gotBytes) {
			t.Errorf("%q: want Bytes %q, got %q", item, wantBytes, gotBytes)
		}
	}
}

func TestStringEqual(t *testing.T) {
	lhs := NewString("dummy")
	defer lhs.Release()

	for _, item := range strtests {
		str := NewString(item)
		defer str.Release()

		if lhs.Equal(str) {
			t.Errorf("Strings compared as equal \"%v\", and \"%v\"", lhs, str)
		}
		if str.Equal(lhs) {
			t.Errorf("Strings compared as equal \"%v\", and \"%v\"", str, lhs)
		}
		if !str.Equal(str) {
			t.Errorf("String did not compared as equal to itself \"%v\", and \"%v\"", str)
		}

		str2 := NewString(item)
		defer str2.Release()
		if !str.Equal(str2) {
			t.Errorf("String did not compared as equal to itself \"%v\", and \"%v\"", str2)
		}
	}
}

func TestStringEqualToString(t *testing.T) {
	lhs := NewString("dummy")
	defer lhs.Release()

	for _, item := range strtests {
		str := NewString(item)
		defer str.Release()

		if lhs.EqualToString(item) {
			t.Errorf("Strings compared as equal \"%v\", and \"%v\"", lhs, item)
		}
		if !str.EqualToString(item) {
			t.Errorf("String did not compare as equal to itself \"%v\"", item)
		}
	}
}
