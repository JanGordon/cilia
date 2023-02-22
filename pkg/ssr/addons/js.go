package addons

import (
	"fmt"

	"rogchap.com/v8go"
)

func jsModifier(c string, ctx v8go.Context) string {
	output, err := ctx.RunScript(c, "resultinline.js")
	if err != nil {
		panic(err)
	}
	fmt.Println(output.String(), c)
	return output.String()
}

func init() {
	a := Addon{
		"{!", "!}", // opening and closing token
		"js",       // name needs to be unique and not match any normal html
		jsModifier, // golang modifier or can just run js version
		true,       // purely ssr               // if former false, js to run on update
		true,       // should alywas be true
		0,          //should always be 0
		"",         //should always be ""
	}
	initNewAddon(&a)
}
