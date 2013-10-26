package gojs

import (
	"testing"
)

type BaseTests struct {
	script    string
	valuetype uint8
	result    string
}

var (
	basetests = []BaseTests{
		{"return 2341234 \"asdf\"", TypeUndefined, ""}, // syntax error
		{"1.5", TypeNumber, "1.5"},
		{"1.5 + 3.0", TypeNumber, "4.5"},
		{"'a' + 'b'", TypeString, "ab"},
		{"new Object()", TypeObject, "[object Object]"},
		{"var obj = {}; obj", TypeObject, "[object Object]"},
		{"var obj = function () { return 1;}; obj", TypeObject, "function () { return 1;}"},
		{"function test() { return 1;}; test", TypeObject, "function test() { return 1;}"},
	}
)

func TestBase(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()
}

func TestEvaluateScript(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	for index, item := range basetests {
		tlog(t, "On item", item, "index", index)
		ret, err := ctx.EvaluateScript(item.script, nil, "./testing.go", 1)
		tlog(t, "Evaluated Script")
		if item.result != "" {
			if err != nil {
				t.Errorf("ctx.EvaluateScript raised an error (script %v)", index)
			} else if ret == nil {
				t.Errorf("ctx.EvaluateScript failed to return a result (script %v)", index)
			} else {
				tlog(t, "No error, and there was a return result.")
				tlog(t, "Type of value is", ctx.ValueType(ret))
				valuetype := ctx.ValueType(ret)
				if valuetype != item.valuetype {
					terrf(t, "ctx.EvaluateScript did not return the expected type (%v instead of %v).", valuetype, item.valuetype)
				}
				stringval := ctx.ToStringOrDie(ret)
				if stringval != item.result {
					terrf(t, "ctx.EvaluateScript returned an incorrect value (%v instead of %v).", stringval, item.result)
				}
			}
		} else {
			// Script has a syntax error
			if err == nil {
				t.Errorf("ctx.EvaluateScript did not raise an error on an invalid script")
			}
			if ret != nil {
				t.Errorf("ctx.EvaluateScript returned a result on an invalid script")
			}
		}
	}
}

func TestCheckScript(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	for _, item := range basetests {
		err := ctx.CheckScriptSyntax(item.script, "./testing.go", 1)
		if err != nil && item.result != "" {
			t.Errorf("ctx.CheckScriptSyntax raised an error but script is good")
		}
		if err == nil && item.result == "" {
			t.Errorf("ctx.CheckScriptSyntax failed to raise an error but script is bad")
		}
	}
}

func TestGarbageCollect(t *testing.T) {
	ctx := NewContext()
	defer ctx.Release()

	ctx.GarbageCollect()
}
