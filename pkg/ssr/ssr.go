package ssr

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/JanGordon/cilia-framework/pkg/component"
	"github.com/JanGordon/cilia-framework/pkg/global"
	"github.com/JanGordon/cilia-framework/pkg/page"
	"github.com/JanGordon/cilia-framework/pkg/ssr/dom"
	"github.com/JanGordon/cilia-framework/pkg/ssr/js"
	"golang.org/x/net/html"
	"rogchap.com/v8go"
)

func Compile() {
	component.SyncComponents()
	filepath.WalkDir(global.ProjectRoot, processPage)
}

func processPage(path string, d fs.DirEntry, err error) error {
	isPage, err := global.PageMatcher.MatchString(d.Name())
	if err != nil {
		panic(err)
	}
	if !d.IsDir() && isPage {
		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		f, _ := os.Open(path)
		jsCtx := v8go.NewContext()
		domCtx, err := html.Parse(f)
		newPage := page.Page{Js: page.JsContext{Path: path, Ctx: jsCtx}, Dom: page.DomContext{Node: domCtx}, TextContents: string(b), Path: path}
		js.RunJS(&newPage)
		newPage = *dom.AssembleDom(&newPage, true)
		writeFile, err := os.OpenFile(path+".out", os.O_WRONLY|os.O_CREATE, 0600)

		child := newPage.Dom.Node.FirstChild
		lastChild := newPage.Dom.Node.LastChild
		for {
			if child.Type != html.CommentNode {
				if err = html.Render(writeFile, child); err != nil {
					panic(err)
				}
			}
			if child != lastChild {
				child = child.NextSibling
			} else {
				break
			}
		}
	}

	return nil
}
