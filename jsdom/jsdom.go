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

var document = js.Global().Get("document")

func CreateElement(name string) vdom.DOMNode {
	return NewDOM(document.Call("createElement", name))
}
func CreateText(text string) vdom.DOMNode {
	return NewDOM(document.Call("createTextNode", text))
}

type jsDOM struct{ js.Value }

func (node jsDOM) Replace(newNode vdom.DOMNode) {
	newJS := newNode.(jsDOM)
	node.Get("parentNode").Call("replaceChild", newJS.Value, node.Value)
}

func (node jsDOM) CreateElement(name string) vdom.DOMNode {
	return CreateElement(name)
}
func (node jsDOM) CreateText(text string) vdom.DOMNode {
	return CreateText(text)
}

func (node jsDOM) NChild() int {
	return node.Get("children").Get("length").Int()
}
func (node jsDOM) child(i int) js.Value {
	return node.Get("children").Call("item", i)
}
func (node jsDOM) Child(i int) vdom.DOMNode {
	return NewDOM(node.child(i))
}
func (node jsDOM) AppendChild(child vdom.DOMNode) {
	node.Call("appendChild", child.(jsDOM).Value)
}
func (node jsDOM) InsertChild(child vdom.DOMNode, i int) {
	node.Call("insertChild", child.(jsDOM).Value, node.child(i))
}
func (node jsDOM) RemoveChild(i int) {
	node.Call("removeChild", node.child(i))
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
