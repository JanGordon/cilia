package global

import (
	"net/http"
	"os"

	"github.com/dlclark/regexp2"
	"rogchap.com/v8go"
)

var ReservedNames = []string{
	"out.html",
}

var PageMatcher = regexp2.MustCompile("[^/.out/].html$", 0)
var BuiltPageMatcher = regexp2.MustCompile(".html.out$", 0)
var ComponentMatcher = regexp2.MustCompile(".cell$", 0)
var ComponentNameSolver = regexp2.MustCompile(`^.*(?=(\.cell))`, 0)

var ProjectRoot, _ = os.Getwd()

var Server *http.Server

var GlobalIsolate = v8go.NewIsolate()

var ComponentContext = v8go.NewContext(GlobalIsolate)

var SSRScriptPath = "./pkg/ssr/preactssr.js"

var SSRScript = ""

// func ErrorCheck() {
// 	if err != nil {
// 		panic(err)
// 	}
// }
