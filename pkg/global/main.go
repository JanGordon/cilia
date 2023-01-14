package global

import (
	"os"

	"github.com/dlclark/regexp2"
)

var ReservedNames = []string{
	"out.html",
}

var PageMatcher = regexp2.MustCompile("[^/.out/].html$", 0)
var ComponentMatcher = regexp2.MustCompile(".cell$", 0)
var ComponentNameSolver = regexp2.MustCompile(`^.*(?=(\.cell))`, 0)

var ProjectRoot, _ = os.Getwd()

// func ErrorCheck() {
// 	if err != nil {
// 		panic(err)
// 	}
// }
