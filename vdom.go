package vdom

import (
	"html"
	"strings"
)

var voidElems = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

func IsVoid(elemName string) bool {
	return voidElems[elemName]
}

type DOMNode interface {
	// Replace the node with a new node.
	Replace(newNode DOMNode)
	// Create a new Element node.
	CreateElement(name string) DOMNode
	// Create a new Text node.
	CreateText(text string) DOMNode

	// Return the number of children of an Element node.
	// May panic if the node is not an Element node.
	NChild() int
	// Return the ith child of an Element node.
	// May panic if the node is not an Element node.
	Child(i int) DOMNode
	// Add a child to the end of an Element node.
	// May panic if the node is not an Element node.
	AppendChild(child DOMNode)
	// Insert a child into an Element node at the given position.
	// May panic if the node is not an Element node.
	InsertChild(child DOMNode, i int)
	// Remove the given child from an Element node.
	// May panic if the node is not an Element node.
	RemoveChild(i int)

	// Set an attribute on an Element node.
	// May panic if the node is not an Element node.
	SetAttr(attr, value string)
	// Remove an attribute from an Element node.
	// May panic if the node is not an Element node.
	DelAttr(attr string)
	// Set the text of a Text node
	// May panic if the node is not a Text node.
	SetText(text string)
}

type Node interface {
	Clone() Node
	HTML() string
	Construct(dom DOMNode) DOMNode
	Patch(dom DOMNode, old Node)
}

type Text struct{ Text string }

func (t Text) Clone() Node {
	return &t
}

func (t *Text) HTML() string {
	return t.Text
}

func (t *Text) Construct(dom DOMNode) DOMNode {
	return dom.CreateText(t.Text)
}
func (t *Text) Patch(dom DOMNode, old Node) {
	if t0, ok := old.(*Text); ok {
		if t != t0 {
			dom.SetText(t.Text)
		}
	} else {
		dom.Replace(t.Construct(dom))
	}
}

type Element struct {
	Name     string
	Attrs    map[string]string
	Children []Node
}

func (e *Element) Clone() Node {
	e2 := &Element{
		e.Name, make(map[string]string),
		make([]Node, len(e.Children)),
	}
	for k, v := range e.Attrs {
		e2.Attrs[k] = v
	}
	for i, child := range e.Children {
		e2.Children[i] = child.Clone()
	}
	return e2
}

func (e *Element) HTML() string {
	b := strings.Builder{}

	// Write opening tag
	b.WriteByte('<')
	// TODO: validate tag name
	b.WriteString(e.Name)
	var attrs bool
	for k, v := range e.Attrs {
		attrs = true
		b.WriteByte(' ')
		// TODO: validate attribute name
		b.WriteString(html.EscapeString(k))
		b.WriteString(`="`)
		b.WriteString(html.EscapeString(v))
		b.WriteByte('"')
	}
	if IsVoid(e.Name) {
		if attrs {
			b.WriteByte(' ')
		}
		b.WriteByte('/')
	}
	b.WriteByte('>')

	if !IsVoid(e.Name) {
		// Write contents
		for _, child := range e.Children {
			b.WriteString(child.HTML())
		}

		// Write closing tag
		b.WriteString("</")
		b.WriteString(e.Name)
		b.WriteByte('>')
	}

	return b.String()
}

func (e *Element) Construct(dom DOMNode) DOMNode {
	node := dom.CreateElement(e.Name)
	for k, v := range e.Attrs {
		node.SetAttr(k, v)
	}
	for _, child := range e.Children {
		node.AppendChild(child.Construct(dom))
	}
	return node
}
func (e *Element) Patch(dom DOMNode, old Node) {
	if e0, ok := old.(*Element); ok && e.Name == e0.Name {
		// Update/add attributes
		for k, v := range e.Attrs {
			if v0, ok := e0.Attrs[k]; !ok || v != v0 {
				dom.SetAttr(k, v)
			}
		}

		// Remove attributes
		for k := range e0.Attrs {
			if _, ok := e.Attrs[k]; !ok {
				dom.DelAttr(k)
			}
		}

		// Update children
		// TODO: improve this algorithm to make use of InsertChild
		n := dom.NChild()
		commonLen := n
		if len(e.Children) < commonLen {
			commonLen = len(e.Children)
		}
		if len(e0.Children) < commonLen {
			commonLen = len(e0.Children)
		}
		for i := 0; i < commonLen; i++ {
			e.Children[i].Patch(dom.Child(i), e0.Children[i])
		}
		// Add new children
		for i := n; i < len(e.Children); i++ {
			dom.AppendChild(e.Children[i].Construct(dom))
		}
		// Remove old children
		for i := len(e.Children); i < n; i++ {
			dom.RemoveChild(i)
		}
	} else {
		dom.Replace(e.Construct(dom))
	}
}
