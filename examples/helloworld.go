package main

import (
	"gojs"
	"fmt"
)

func main() {
	ctx := gojs.NewContext()
	defer ctx.Release()

	ret, err := ctx.EvaluateScript("['hello', 'world'].join(' ')", nil, ".", 0)

	if err != nil {
		fmt.Println("Script had an error :(", ctx.ToStringOrDie(err))
		return
	}

	if ret == nil {
		fmt.Println("Nothing returned...")
		return
	}

	retstr := ctx.ToStringOrDie(ret)

	fmt.Println(retstr)
}