package global

import (
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

var ComponentContext = v8go.NewContext()

// func ErrorCheck() {
// 	if err != nil {
// 		panic(err)
// 	}
// }
