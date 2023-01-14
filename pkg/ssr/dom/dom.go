package dom

import (
	"fmt"
	"os"

	"github.com/JanGordon/cilia-framework/pkg/component"
	"github.com/JanGordon/cilia-framework/pkg/page"
	"golang.org/x/net/html"
)

func AssembleDom(document *page.Page) *page.Page {
	for _, i := range page.GetAllDescendants(document.Dom.Node) {
		if i.Type == html.ElementNode {
			for _, c := range component.Components {
				if c.Label == i.Data {
					c.AllUsers = append(c.AllUsers, document.Path)
					c.AllUsers = append(c.AllUsers, document.AllUsers...)
					if stringInSlice(c.Path, c.AllUsers) {
						panic(fmt.Sprintf("import loop detected: %v imports self", c.Path))
					}
					// the same scope as the page it is used on must be used so it can acces variables passed to it
					file, err := os.Open(c.Path)
					if err != nil {
						panic(err)
					}
					fileText, err := os.ReadFile(c.Path)
					if err != nil {
						panic(err)
					}
					nodes, err := html.Parse(file)
					if err != nil {
						panic(err)
					}
					body := nodes.LastChild.LastChild
					// body.Parent.RemoveChild(body)

					// this js returns the html for the component
					c.Generator = fmt.Sprintf(`function (){%v}`, file)
					newDocument := AssembleDom(&page.Page{Js: document.Js, Dom: page.DomContext{Node: body}, TextContents: string(fileText), Path: c.Path, AllUsers: c.AllUsers})
					for _, v := range page.GetChildren(newDocument.Dom.Node) {
						v.Parent.RemoveChild(v)
						i.AppendChild(v)
					}
				}
			}
		}
	}
	return document
}

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if s == v {
			return true
		}
	}
	return false
}
