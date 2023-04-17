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

func Compile(path string, isSSR bool, request string, json string) map[string]*page.Page {
	fmt.Println("Started compile", isSSR)
	ssr = isSSR
	req = request
	pageMap = make(map[string]*page.Page)
	if !ssr {
		prebuiltPages = nil
		pageMap = nil
	}
	component.SyncComponents()
	filepath.WalkDir(path, processPage)
	fmt.Println("ended compile succesfully", isSSR)
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
			jsCtx := v8go.NewContext(global.GlobalIsolate)
			newPage = page.Page{Js: page.JsContext{Path: path, Ctx: jsCtx}, Dom: page.DomContext{Node: nil}, TextContents: string(b), Path: path}
			newPage.Js.Ctx.RunScript(fmt.Sprintf("var REQ = 'this is not ssr'", req), "requestdef.js")

		}
		// js.RunJS(&newPage)
		jsfile := page.JsFile{Contents: "", PathOfRoute: path}
		newPage = *dom.AssembleDom(&newPage, true, ssr, &jsfile)
		// writeFile, err := os.OpenFile(path+".out", os.O_WRONLY|os.O_CREATE, 0600)
		// no longer writing to file because of ssr
		if ssr {
			scriptNode := &html.Node{
				Data: "script",
				Type: html.ElementNode,
			}
			scriptNode.AppendChild(&html.Node{
				Data: jsfile.Contents,
				Type: html.TextNode,
			})
			newPage.Dom.Node.LastChild.AppendChild(scriptNode)

		}
		pageBuf := bytes.NewBuffer([]byte{})
		child := newPage.Dom.Node
		lastChild := newPage.Dom.Node
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
		// now the page has run its server side rendering
		// we run preact ssr
		// preactRender(pageBuf)
		newPage.TextContents = pageBuf.String()
		if ssr {
			pageMap[path+".out"] = &newPage
		} else {
			fmt.Println("adding prebuilt page........")
			prebuiltPages = append(prebuiltPages, &newPage)
		}
	}

	return nil
}

func preactRender(page *bytes.Buffer) {
	preactSSRContext := v8go.NewContext(global.GlobalIsolate)
	window := preactSSRContext.Global()
	window.Set("pageString", page.String())
	if global.SSRScript != "" {
		fmt.Println("Found embeded script!!")
		preactSSRContext.RunScript(global.SSRScript, "ssrscript.js")
	} else {
		// probably running with go run main.go
		file, err := os.ReadFile(global.SSRScriptPath)
		if err != nil {
			panic(fmt.Errorf("%v, either incorrectly built (use the script on github) or the ssr script has moved", err))
		}
		preactSSRContext.RunScript(string(file), "ssrscript.js")
		fmt.Println("didnt find ebebed script")

	}
}

func FlushCache() {
	prebuiltPages = nil
	pageMap = nil
}

func getPageFromPath(path string, pages []*page.Page) *page.Page {
	for _, page := range pages {
		if page.Path == path {
			return page
		}
	}
	return nil
}
