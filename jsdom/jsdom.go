// +build js

package jsdom

import (
	"syscall/js"

	"github.com/vktec/vdom"
)

// Create a new DOMNode from a JavaScript Node object
func NewDOM(val js.Value) vdom.DOMNode {
	return DOM{val}
}

type DOM struct{ js.Value }

func (node DOM) Replace(newNode vdom.DOMNode) vdom.DOMNode {
	newJS := newNode.(DOM)
	node.Get("parentNode").Call("replaceChild", newJS.Value, node.Value)
	return newNode
}

var document = js.Global().Get("document")

func (node DOM) CreateElement(name string) vdom.DOMNode {
	return NewDOM(document.Call("createElement", name))
}
func (node DOM) CreateText(text string) vdom.DOMNode {
	return NewDOM(document.Call("createTextNode", text))
}

func (node DOM) FirstChild() vdom.DOMNode {
	return NewDOM(node.Get("firstChild"))
}
func (node DOM) NextSibling() vdom.DOMNode {
	return NewDOM(node.Get("nextSibling"))
}
func (node DOM) AppendChild(child vdom.DOMNode) {
	node.Call("appendChild", child.(DOM).Value)
}
func (node DOM) InsertBefore(newChild, oldChild vdom.DOMNode) {
	node.Call("insertChild", newChild.(DOM).Value, oldChild.(DOM).Value)
}
func (node DOM) RemoveChild(child vdom.DOMNode) {
	node.Call("removeChild", child.(DOM).Value)
}

func (node DOM) SetAttr(attr, value string) {
	node.Call("setAttribute", attr, value)
}
func (node DOM) DelAttr(attr string) {
	node.Call("removeAttribute", attr)
}
func (node DOM) SetText(text string) {
	node.Set("nodeValue", text)
}
