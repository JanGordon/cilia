package dom

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JanGordon/cilia-framework/pkg/component"
	"github.com/JanGordon/cilia-framework/pkg/global"
	"github.com/JanGordon/cilia-framework/pkg/page"
	"github.com/JanGordon/cilia-framework/pkg/ssr/addons"
	"github.com/JanGordon/cilia-framework/pkg/url"
	"golang.org/x/net/html"
	"rogchap.com/v8go"
)

func AssembleDom(document *page.Page, root bool, ssr bool, jsFile *page.JsFile) *page.Page {
	newDocument := addons.ReplaceAddons(document, ssr, jsFile)
	document = &newDocument
	newNode, err := html.Parse(strings.NewReader(document.TextContents))
	if err != nil {
		panic(err)
	}
	document.Dom.Node = newNode.LastChild
	for _, i := range page.GetAllDescendants(document.Dom.Node) {
		if i.Type == html.ElementNode {
			for _, c := range component.Components {
				if hasAttribute("ssr", "", i) {
					if !ssr {
						continue
					}
					// this means it shouldnt be built as it is build
				} else {
					if ssr {
						continue
					}
				}
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
					// body.Parent.RemoveChild(body)

					// this js returns the html for the component
					// var args []string
					// //get expected args
					// for _, v := range page.GetChildren(nodes.LastChild.FirstChild) {
					// 	if v.Data == "script" {
					// 		args = getComponentArgs(v.FirstChild.Data)
					// 	}
					// }
					// //getting attributes passed
					// attributes := make(map[string]string)

					// for _, attribute := range i.Attr {
					// 	renderedAttribute, err := document.Js.Ctx.RunScript(fmt.Sprintf("%v", attribute.Val), "attrscript")
					// 	if err != nil {
					// 		attributes[attribute.Key] = ""

					// 	} else {
					// 		attributes[attribute.Key] = renderedAttribute.String()

					// 	}
					// }

					argText := ""
					argDefinitions := "let _"
					for _, arg := range i.Attr {
						argText += fmt.Sprintf("%v:'%v',", arg.Key, arg.Val)
						argDefinitions += fmt.Sprintf(",%v = args.%v", arg.Key, arg.Key)
					}
					funcName := strings.TrimSuffix(filepath.Base(file.Name()), filepath.Ext(file.Name()))
					// before we run the generator we need to make sure there are no other components to have a genertor made for them
					jsCtx := v8go.NewContext()
					newDocument := AssembleDom(&page.Page{Js: page.JsContext{Path: c.Path, Ctx: jsCtx}, Dom: page.DomContext{Node: nodes.LastChild.LastChild}, TextContents: string(fileText), Path: c.Path, AllUsers: c.AllUsers, ExternalScripts: document.ExternalScripts}, false, ssr, jsFile)

					// we need to render the output as a string to pass to the js:
					var bytesOfComponent = bytes.NewBuffer([]byte{})
					for _, node := range page.GetChildren(newDocument.Dom.Node) {
						html.Render(bytesOfComponent, node)
					}
					c.Generator = fmt.Sprintf("function %v (args){%v\n return `%v`}", funcName, argDefinitions, string(bytesOfComponent.Bytes()))

					// we need to make sure all the other js on the page has been run
					global.ComponentContext.RunScript(string(c.Generator), file.Name())
					// fmt.Println(string(c.Generator), funcName+"()")
					componentHTML, err := global.ComponentContext.RunScript(fmt.Sprintf("%v({%v})", funcName, argText), file.Name()+"ssr")
					if err != nil {
						panic(err)
					}
					// fmt.Println(componentHTML)
					// need to rerun assemble dom to mkae sure all returned components are resolved

					newDocument = AssembleDom(&page.Page{Js: page.JsContext{Path: c.Path, Ctx: jsCtx}, Dom: page.DomContext{Node: nil}, TextContents: componentHTML.String(), Path: c.Path, AllUsers: c.AllUsers, ExternalScripts: document.ExternalScripts}, false, ssr, jsFile)
					// for _, v := range page.GetAllDescendants(newDocument.Dom.Node) {
					// 	fmt.Println("nodes in compoentn: ", v.Data)

					// }
					for _, v := range page.GetChildren(newDocument.Dom.Node) {
						v.Parent.RemoveChild(v)
						i.AppendChild(v)

					}
				} else if i.Data == "script" && hasAttribute("src", "*", i) {
					scriptPath, err := url.ResolvePath(getAttribute("src", i), document.Path)
					if err != nil {
						panic(fmt.Sprintf("could not resolve path: %v", getAttribute("src", i)))
					}
					fmt.Println("Docuemtn paht", document.ExternalScripts, scriptPath)
					if document.ExternalScripts == nil {
						externScript := []string{}
						document.ExternalScripts = &externScript
					}
					f := append(*document.ExternalScripts, scriptPath)
					document.ExternalScripts = &f
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

func hasAttribute(attr string, val string, node *html.Node) bool {
	for _, v := range node.Attr {
		if attr == v.Key {
			if val == v.Val || val == "*" {
				return true
			}
		}
	}
	return false
}

func getAttribute(attr string, node *html.Node) string {
	for _, v := range node.Attr {
		if attr == v.Key {
			return v.Val
		}
	}
	return ""
}

func getComponentArgs(s string) []string {
	componentArgs := []string{}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		// fmt.Println(line)
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//use") {
			componentArgs = append(componentArgs, strings.TrimPrefix(line, "//use "))
		}
	}

	return componentArgs
}
