package addons

import "github.com/gomarkdown/markdown"

func markdownModifier(c string) string {
	md := []byte("## markdown document")
	output := markdown.ToHTML(md, nil, nil)
	return string(output)
}

func init() {
	a := Addon{"{>", "<}", "md", markdownModifier}
	initNewAddon(&a)
}
