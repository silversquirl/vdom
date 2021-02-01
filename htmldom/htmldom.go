package htmldom

import (
	"github.com/vktec/vdom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func New(node *html.Node) DOM {
	return DOM{node}
}

type DOM struct{ *html.Node }

// Implement DOMNode
func (dom DOM) Replace(newDOM vdom.DOMNode) vdom.DOMNode {
	if dom.Node != nil && dom.Node.Parent != nil {
		node := newDOM.(DOM).Node
		dom.Node.Parent.InsertBefore(node, dom.Node)
		dom.Node.Parent.RemoveChild(dom.Node)
	}
	return newDOM
}

// Implement DOMNode
func (dom DOM) CreateElement(name string) vdom.DOMNode {
	return DOM{&html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Lookup([]byte(name)),
		Data:     name,
	}}
}

// Implement DOMNode
func (dom DOM) CreateText(text string) vdom.DOMNode {
	return DOM{&html.Node{
		Type: html.TextNode,
		Data: text,
	}}
}

// Implement DOMNode
func (dom DOM) FirstChild() vdom.DOMNode {
	return DOM{dom.Node.FirstChild}
}

// Implement DOMNode
func (dom DOM) NextSibling() vdom.DOMNode {
	return DOM{dom.Node.NextSibling}
}

// Implement DOMNode
func (dom DOM) AppendChild(child vdom.DOMNode) {
	node := child.(DOM).Node
	dom.Node.AppendChild(node)
}

// Implement DOMNode
func (dom DOM) InsertBefore(newChild, oldChild vdom.DOMNode) {
	newNode := newChild.(DOM).Node
	oldNode := oldChild.(DOM).Node
	dom.Node.InsertBefore(newNode, oldNode)
}

// Implement DOMNode
func (dom DOM) RemoveChild(child vdom.DOMNode) {
	node := child.(DOM).Node
	dom.Node.RemoveChild(node)
}

// Implement DOMNode
func (dom DOM) SetAttr(key, value string) {
	for i, attr := range dom.Node.Attr {
		if attr.Key == key {
			dom.Node.Attr[i].Val = value
			return
		}
	}
	dom.Node.Attr = append(dom.Node.Attr, html.Attribute{Key: key, Val: value})
}

// Implement DOMNode
func (dom DOM) DelAttr(key string) {
	for i, attr := range dom.Node.Attr {
		if attr.Key == key {
			dom.Node.Attr = append(dom.Node.Attr[:i], dom.Node.Attr[i+1:]...)
			return
		}
	}
}

// Implement DOMNode
func (dom DOM) SetText(text string) {
	dom.Node.Data = text
}
