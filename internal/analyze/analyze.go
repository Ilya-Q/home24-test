package analyze

import (
	"slices"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type HTMLVisitor interface {
	Visit(node *html.Node) bool
}

func Walk(root *html.Node, visitors []HTMLVisitor) {
	cur := root
	remaining := slices.Clone(visitors)
	for len(remaining) > 0 && cur != nil {
		remaining = slices.DeleteFunc(remaining, func(v HTMLVisitor) bool {
			return v.Visit(cur)
		})
		if cur.FirstChild != nil {
			cur = cur.FirstChild
		} else if cur.NextSibling != nil {
			cur = cur.NextSibling
		} else {
			for cur != nil && cur.NextSibling == nil {
				cur = cur.Parent
			}
			if cur != nil {
				cur = cur.NextSibling
			}
		}
	}
}

type TitleGetter struct {
	Title string
}

func (v *TitleGetter) Visit(node *html.Node) bool {
	if node.DataAtom == atom.Title {
		v.Title = node.FirstChild.Data
		return true
	}
	return false
}

type DoctypeGetter struct {
	Doctype string
}

func (v *DoctypeGetter) Visit(node *html.Node) bool {
	if node.Type == html.DoctypeNode {
		if len(node.Attr) == 0 {
			v.Doctype = node.Data
		} else {
			triple := node.Attr[0]
			v.Doctype = triple.Val
		}
		return true
	}
	return false
}

type HeadingCounter struct {
	H1, H2, H3, H4, H5, H6 uint64
}

func (v *HeadingCounter) Visit(node *html.Node) bool {
	switch node.DataAtom {
	case atom.H1:
		v.H1++
	case atom.H2:
		v.H2++
	case atom.H3:
		v.H3++
	case atom.H4:
		v.H4++
	case atom.H5:
		v.H5++
	case atom.H6:
		v.H6++
	}
	return false
}
