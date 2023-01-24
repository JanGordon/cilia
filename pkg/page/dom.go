package page

import "golang.org/x/net/html"

type DomContext struct {
	Node *html.Node
}

// type Node struct {
// 	*html.Node
// }

// func (n *Node) Children() []*Node {
// 	var children []*Node
// 	for i := n.FirstChild; i != nil; i = i.NextSibling {
// 		children = append(children, &Node{i})
// 	}
// 	return children
// }

// func ParseHTML() {

// }

func GetChildren(n *html.Node) []*html.Node {
	var children []*html.Node
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		children = append(children, i)
	}
	return children
}

func GetAllDescendants(n *html.Node) []*html.Node {
	descendants := GetChildren(n)
	for i := n.FirstChild; i != nil; i = i.NextSibling {
		descendants = append(descendants, GetAllDescendants(i)...)
	}
	return descendants
}
