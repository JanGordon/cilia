package addons

import (
	"fmt"
	"strings"

	"rogchap.com/v8go"
)

func forModifier(c string, ctx v8go.Context) string {
	ctx.RunScript(fmt.Sprintf("var forresult = '';%v{forresult+='%v'}", strings.Split(c, "\n")[0], strings.Replace(c, strings.Split(c, "\n")[0], "", 1)), "inlineforloop.js")
	forresult, err := ctx.RunScript("forresult", "forresult.js")
	if err != nil {
		panic(err)
	}

	return forresult.String()
}

func init() {
	a := Addon{
		"{.", ".}", // opening and closing token
		"for",       // name needs to be unique and not match any normal html
		forModifier, // golang modifier or can just run js version
		true,        // purely ssr               // if former false, js to run on update
		true,        // should alywas be true
		0,           //should always be 0
		"",          //should always be ""
	}
	initNewAddon(&a)
}
