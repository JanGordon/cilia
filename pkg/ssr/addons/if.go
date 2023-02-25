package addons

import (
	"fmt"
	"strings"

	"rogchap.com/v8go"
)

func ifModifier(c string, ctx v8go.Context, id string) (string, string) {
	script := fmt.Sprintf("var ifresult = '';%v{ifresult=`%v`};", strings.Split(c, "\n")[0], strings.Replace(c, strings.Split(c, "\n")[0], "", 1))
	fmt.Println(script)
	t, err := ctx.RunScript(script, "inlineforloop.js")
	if err != nil {
		panic(err)
	}
	fmt.Println("If result: ", t.String())
	forresult, err := ctx.RunScript("ifresult", "ifresult.js")
	if err != nil {
		panic(err)
	}

	return forresult.String(), fmt.Sprintf(`{document.querySelector('[cilia-id="%v"]').addEventListener('click', function (e) {%vdocument.querySelector('[cilia-id="%v"]').innerHTML=ifresult; console.log(ifresult)})}`, id, script, id)
}

func init() {
	a := Addon{
		"{?", "?}", // opening and closing token
		"if",       // name needs to be unique and not match any normal html
		ifModifier, // golang modifier or can just run js version
		true,       // purely ssr               // if former false, js to run on update
		true,       // should alywas be true
		0,          //should always be 0
		"",         //should always be ""
	}
	initNewAddon(&a)
}
