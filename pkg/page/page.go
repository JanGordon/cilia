package page

// a type defintion for a page

type Page struct {
	Js           JsContext
	Script       []string
	Dom          DomContext
	TextContents string
	Path         string
	AllUsers     []string
}
