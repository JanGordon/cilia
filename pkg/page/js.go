package page

import "rogchap.com/v8go"

type JsContext struct {
	Path string
	Ctx  *v8go.Context
}

var Ctxs []JsContext
