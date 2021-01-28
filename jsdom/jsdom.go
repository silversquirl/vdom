// +build js

package jsdom

import (
	"syscall/js"

	"github.com/vktec/vdom"
)

// Create a new DOMNode from a JavaScript Node object
func NewDOM(val js.Value) vdom.DOMNode {
	return jsDOM{val}
}

type jsDOM struct{ js.Value }

func (node jsDOM) Replace(newNode vdom.DOMNode) {
	newJS := newNode.(jsDOM)
	node.Get("parentNode").Call("replaceChild", newJS.Value, node.Value)
}

var document = js.Global().Get("document")

func (node jsDOM) CreateElement(name string) vdom.DOMNode {
	return NewDOM(document.Call("createElement", name))
}
func (node jsDOM) CreateText(text string) vdom.DOMNode {
	return NewDOM(document.Call("createTextNode", text))
}

func (node jsDOM) FirstChild() vdom.DOMNode {
	return NewDOM(node.Get("firstChild"))
}
func (node jsDOM) NextSibling() vdom.DOMNode {
	return NewDOM(node.Get("nextSibling"))
}
func (node jsDOM) AppendChild(child vdom.DOMNode) {
	node.Call("appendChild", child.(jsDOM).Value)
}
func (node jsDOM) InsertBefore(newChild, oldChild vdom.DOMNode) {
	node.Call("insertChild", newChild.(jsDOM).Value, oldChild.(jsDOM).Value)
}
func (node jsDOM) RemoveChild(child vdom.DOMNode) {
	node.Call("removeChild", child.(jsDOM).Value)
}

func (node jsDOM) SetAttr(attr, value string) {
	node.Call("setAttribute", attr, value)
}
func (node jsDOM) DelAttr(attr string) {
	node.Call("removeAttribute", attr)
}
func (node jsDOM) SetText(text string) {
	node.Set("nodeValue", text)
}
