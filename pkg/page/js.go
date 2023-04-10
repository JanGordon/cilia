package page

import "rogchap.com/v8go"

type JsContext struct {
	Path string
	Ctx  *v8go.Context
}

type JsFile struct {
	Contents    string
	PathOfRoute string
}

var Ctxs []JsContext
