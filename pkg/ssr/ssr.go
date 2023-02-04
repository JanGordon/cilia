package ssr

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/JanGordon/cilia-framework/pkg/component"
	"github.com/JanGordon/cilia-framework/pkg/global"
	"github.com/JanGordon/cilia-framework/pkg/page"
	"github.com/JanGordon/cilia-framework/pkg/ssr/dom"
	"golang.org/x/net/html"
	"rogchap.com/v8go"
)

var ssr bool
var req string
var prebuiltPages []*page.Page
var pageMap map[string]*page.Page

func Compile(path string, isSSR bool, request string) map[string]*page.Page {
	ssr = isSSR
	req = request
	pageMap = make(map[string]*page.Page)
	component.SyncComponents()
	filepath.WalkDir(path, processPage)
	return pageMap
}

func processPage(path string, d fs.DirEntry, err error) error {
	isPage, err := global.PageMatcher.MatchString(d.Name())
	if err != nil {
		panic(err)
	}
	if !d.IsDir() && isPage {
		var newPage page.Page
		if ssr {
			newPage = *getPageFromPath(path, prebuiltPages)
			newPage.Js.Ctx.RunScript(fmt.Sprintf("var REQ = '%v'", req), "request.js")

		} else {
			b, err := os.ReadFile(path)
			if err != nil {
				panic(err)
			}
			jsCtx := v8go.NewContext()
			newPage = page.Page{Js: page.JsContext{Path: path, Ctx: jsCtx}, Dom: page.DomContext{Node: nil}, TextContents: string(b), Path: path}
			newPage.Js.Ctx.RunScript(fmt.Sprintf("var REQ = 'this is not ssr'", req), "requestdef.js")

		}
		// js.RunJS(&newPage)
		newPage = *dom.AssembleDom(&newPage, true, ssr)
		// writeFile, err := os.OpenFile(path+".out", os.O_WRONLY|os.O_CREATE, 0600)
		// no longer writing to file because of ssr
		pageBuf := bytes.NewBuffer([]byte{})
		child := newPage.Dom.Node.FirstChild
		lastChild := newPage.Dom.Node.LastChild
		for {
			if child.Type != html.CommentNode {
				if err = html.Render(pageBuf, child); err != nil {
					panic(err)
				}
			}
			if child != lastChild {
				child = child.NextSibling
			} else {
				break
			}
		}

		newPage.TextContents = pageBuf.String()
		if ssr {
			pageMap[path+".out"] = &newPage
		} else {
			prebuiltPages = append(prebuiltPages, &newPage)
		}
	}

	return nil
}

func getPageFromPath(path string, pages []*page.Page) *page.Page {
	for _, page := range pages {
		if page.Path == path {
			return page
		}
	}
	return nil
}
