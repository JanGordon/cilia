package page

// a type defintion for a page

type Page struct {
	Js           JsContext
	Dom          DomContext
	TextContents string
	Path         string
	AllUsers     []string
}
