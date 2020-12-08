package testdom

import (
	"github.com/vktec/vdom"
)

func NewTestDOM() *TestDOM {
	return makeTestDOM(&vdom.Element{})
}
func makeTestDOM(node vdom.Node) *TestDOM {
	return &TestDOM{&node}
}

type TestDOM struct{ Node *vdom.Node }

func (dom *TestDOM) Element() *vdom.Element {
	return (*dom.Node).(*vdom.Element)
}
func (dom *TestDOM) Text() *vdom.Text {
	return (*dom.Node).(*vdom.Text)
}

// Implement DOMNode
func (dom *TestDOM) Replace(newNode vdom.DOMNode) {
	*dom.Node = *newNode.(*TestDOM).Node
}

// Implement DOMNode
func (dom *TestDOM) CreateElement(name string) vdom.DOMNode {
	return makeTestDOM(&vdom.Element{Name: name})
}

// Implement DOMNode
func (dom *TestDOM) CreateText(text string) vdom.DOMNode {
	return makeTestDOM(&vdom.Text{text})
}

// Implement DOMNode
func (dom *TestDOM) NChild() int {
	return len(dom.Element().Children)
}

// Implement DOMNode
func (dom *TestDOM) Child(i int) vdom.DOMNode {
	return &TestDOM{&dom.Element().Children[i]}
}

// Implement DOMNode
func (dom *TestDOM) AppendChild(child vdom.DOMNode) {
	c := &dom.Element().Children
	*c = append(*c, *child.(*TestDOM).Node)
}

// Implement DOMNode
func (dom *TestDOM) InsertChild(child vdom.DOMNode, i int) {
	c := &dom.Element().Children
	*c = append(append((*c)[:i], *child.(*TestDOM).Node), (*c)[i:]...)
}

// Implement DOMNode
func (dom *TestDOM) RemoveChild(i int) {
	c := &dom.Element().Children
	*c = append((*c)[:i], (*c)[i+1:]...)
}

// Implement DOMNode
func (dom *TestDOM) SetAttr(attr, value string) {
	elem := dom.Element()
	if elem.Attrs == nil {
		elem.Attrs = make(map[string]string)
	}
	elem.Attrs[attr] = value
}

// Implement DOMNode
func (dom *TestDOM) DelAttr(attr string) {
	delete(dom.Element().Attrs, attr)
}

// Implement DOMNode
func (dom *TestDOM) SetText(text string) {
	dom.Text().Text = text
}
