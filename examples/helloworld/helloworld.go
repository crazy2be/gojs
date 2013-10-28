package main

import (
	"fmt"
	"github.com/crazy2be/gojs"
)

func main() {
	ctx := gojs.NewContext()
	defer ctx.Release()

	ret, err := ctx.EvaluateScript("['hello', 'world'].join(' ')", nil, ".", 0)

	if err != nil {
		fmt.Println("Script had an error :(", err)
		return
	}

	if ret == nil {
		fmt.Println("Nothing returned...")
		return
	}

	retstr := ctx.ToStringOrDie(ret)

	fmt.Println(retstr)
}
