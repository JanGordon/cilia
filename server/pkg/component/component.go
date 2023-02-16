package component

import (
	"io/fs"
	"path/filepath"

	"github.com/JanGordon/cilia-framework/server/pkg/global"
)

var Components []Component

type Component struct {
	Label     string
	Path      string
	Used      bool
	Generator string
	AllUsers  []string
}

func (c *Component) JS() {

}

func SyncComponents() {
	Components = nil
	filepath.WalkDir(filepath.Join(global.ProjectRoot, "components"), addComponent)
}

func addComponent(path string, d fs.DirEntry, err error) error {
	isC, err := global.ComponentMatcher.MatchString(d.Name())
	if err != nil {
		panic(err)
	}
	if !d.IsDir() && isC {
		label, err := global.ComponentNameSolver.FindStringMatch(d.Name())
		if err != nil {
			panic(err)
		}
		Components = append(Components, Component{label.String(), path, false, "undefined", nil})
	}

	return nil
}

// func StringToNodes(n io.Reader) html.Node {
// 	nodes, err := html.Parse(n)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return
// }
