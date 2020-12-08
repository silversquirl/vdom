package htmldom

import (
	"github.com/vktec/vdom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func New(node *html.Node) vdom.DOMNode {
	return htmlDOM{node}
}

type htmlDOM struct{ node *html.Node }

// Implement DOMNode
func (dom htmlDOM) Replace(newNode vdom.DOMNode) {
	node := newNode.(htmlDOM).node
	dom.node.FirstChild = node.FirstChild
	dom.node.LastChild = node.LastChild
	dom.node.Type = node.Type
	dom.node.DataAtom = node.DataAtom
	dom.node.Data = node.Data
	dom.node.Namespace = node.Namespace
	dom.node.Attr = node.Attr
}

// Implement DOMNode
func (dom htmlDOM) CreateElement(name string) vdom.DOMNode {
	return htmlDOM{&html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Lookup([]byte(name)),
		Data:     name,
	}}
}

// Implement DOMNode
func (dom htmlDOM) CreateText(text string) vdom.DOMNode {
	return htmlDOM{&html.Node{
		Type: html.TextNode,
		Data: text,
	}}
}

// Implement DOMNode
func (dom htmlDOM) FirstChild() vdom.DOMNode {
	return htmlDOM{dom.node.FirstChild}
}

// Implement DOMNode
func (dom htmlDOM) NextSibling() vdom.DOMNode {
	return htmlDOM{dom.node.NextSibling}
}

// Implement DOMNode
func (dom htmlDOM) AppendChild(child vdom.DOMNode) {
	node := child.(htmlDOM).node
	dom.node.AppendChild(node)
}

// Implement DOMNode
func (dom htmlDOM) InsertBefore(newChild, oldChild vdom.DOMNode) {
	newNode := newChild.(htmlDOM).node
	oldNode := oldChild.(htmlDOM).node
	dom.node.InsertBefore(newNode, oldNode)
}

// Implement DOMNode
func (dom htmlDOM) RemoveChild(child vdom.DOMNode) {
	node := child.(htmlDOM).node
	dom.node.RemoveChild(node)
}

// Implement DOMNode
func (dom htmlDOM) SetAttr(key, value string) {
	for i, attr := range dom.node.Attr {
		if attr.Key == key {
			dom.node.Attr[i].Val = value
			return
		}
	}
	dom.node.Attr = append(dom.node.Attr, html.Attribute{Key: key, Val: value})
}

// Implement DOMNode
func (dom htmlDOM) DelAttr(key string) {
	for i, attr := range dom.node.Attr {
		if attr.Key == key {
			dom.node.Attr = append(dom.node.Attr[:i], dom.node.Attr[i+1:]...)
			return
		}
	}
}

// Implement DOMNode
func (dom htmlDOM) SetText(text string) {
	dom.node.Data = text
}
