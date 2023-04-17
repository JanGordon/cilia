package addons

import (
	"fmt"
	"strings"

	"rogchap.com/v8go"
)

func forModifier(c string, ctx v8go.Context, id string) (string, string) {
	script := fmt.Sprintf("var forresult = '';%v{forresult+='%v'}", strings.Split(c, "\n")[0], strings.Replace(c, strings.Split(c, "\n")[0], "", 1))
	ctx.RunScript(script, "inlineforloop.js")
	forresult, err := ctx.RunScript("forresult", "forresult.js")
	if err != nil {
		panic(err)
	}
	fmt.Println("There si  afro loop")
	return forresult.String(), fmt.Sprintf("{document.getElementById('%v').addEventListener('click', function (e) {%vdocument.getElementById('%v').innerText=forresult})}", id, script, id)
}

func init() {
	a := Addon{
		"{.", ".}", // opening and closing token
		"for",       // name needs to be unique and not match any normal html
		forModifier, // golang modifier or can just run js version
		false,       // purely ssr               // if former false, js to run on update
		true,        // should alywas be true
		0,           //should always be 0
		"",          //should always be ""
	}
	initNewAddon(&a)
}
