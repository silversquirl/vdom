package vdom

import (
	"fmt"

	"golang.org/x/net/html"
)

type DOMNode interface {
	// Replace the node with a new node, returning the new node.
	Replace(newNode DOMNode) DOMNode
	// Create a new Element node.
	CreateElement(name string) DOMNode
	// Create a new Text node.
	CreateText(text string) DOMNode

	// Return the first child of an Element node.
	// May panic if the node is not an Element node.
	FirstChild() DOMNode
	// Return the next sibling of a node.
	NextSibling() DOMNode
	// Add a child to the end of an Element node.
	// May panic if the node is not an Element node.
	AppendChild(child DOMNode)
	// Insert a child into an Element node at the given position.
	// May panic if the node is not an Element node.
	InsertBefore(newChild, oldChild DOMNode)
	// Remove the given child from an Element node.
	// May panic if the node is not an Element node.
	RemoveChild(child DOMNode)

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

func Clone(node *html.Node) *html.Node {
	var attr []html.Attribute
	if node.Attr != nil {
		attr = make([]html.Attribute, len(node.Attr))
		copy(attr, node.Attr)
	}

	firstChild, lastChild := clone(node.FirstChild)
	return &html.Node{
		Parent:      node.Parent,
		FirstChild:  firstChild,
		LastChild:   lastChild,
		PrevSibling: node.PrevSibling,
		NextSibling: node.NextSibling,

		Type:      node.Type,
		DataAtom:  node.DataAtom,
		Data:      node.Data,
		Namespace: node.Namespace,
		Attr:      attr,
	}
}
func clone(node *html.Node) (me, lastSibling *html.Node) {
	if node == nil {
		return nil, nil
	}

	var attr []html.Attribute
	if node.Attr != nil {
		attr = make([]html.Attribute, len(node.Attr))
		copy(attr, node.Attr)
	}

	nextSibling, lastSibling := clone(node.NextSibling)
	firstChild, lastChild := clone(node.FirstChild)
	me = &html.Node{
		Parent:      node.Parent,
		FirstChild:  firstChild,
		LastChild:   lastChild,
		PrevSibling: node.PrevSibling,
		NextSibling: nextSibling,

		Type:      node.Type,
		DataAtom:  node.DataAtom,
		Data:      node.Data,
		Namespace: node.Namespace,
		Attr:      attr,
	}
	return
}

func Construct(node *html.Node, dom DOMNode) DOMNode {
	switch node.Type {
	case html.TextNode:
		return dom.CreateText(node.Data)

	case html.ElementNode:
		domNode := dom.CreateElement(node.Data)
		for _, attr := range node.Attr {
			domNode.SetAttr(attr.Key, attr.Val)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			domNode.AppendChild(Construct(child, dom))
		}
		return domNode

	default:
		panic(fmt.Sprint("Cannot construct node of type ", node.Type))
	}
}

func Patch(dom DOMNode, node, oldNode *html.Node) DOMNode {
	if oldNode == nil || node.Type != oldNode.Type {
		return dom.Replace(Construct(node, dom))
	}

	switch node.Type {
	case html.TextNode:
		if node.Data != oldNode.Data {
			dom.SetText(node.Data)
		}

	case html.ElementNode:
		if node.Data != oldNode.Data {
			return dom.Replace(Construct(node, dom))
		}

		// Update/add attributes
		oldAttr := make(map[string]string)
		for _, attr := range oldNode.Attr {
			oldAttr[attr.Key] = attr.Val
		}
		for _, attr := range node.Attr {
			if oldVal, ok := oldAttr[attr.Key]; !ok || attr.Val != oldVal {
				dom.SetAttr(attr.Key, attr.Val)
			}
			delete(oldAttr, attr.Key)
		}
		// Remove old attributes
		for k := range oldAttr {
			dom.DelAttr(k)
		}

		// Update children
		// TODO: improve this algorithm to make use of InsertChild
		child := node.FirstChild
		oldChild := oldNode.FirstChild
		domChild := dom.FirstChild()
		for child != nil && oldChild != nil && domChild != nil {
			domChild = Patch(domChild, child, oldChild)
			child = child.NextSibling
			oldChild = oldChild.NextSibling
			domChild = domChild.NextSibling()
		}
		// Remove old children
		for oldChild != nil && domChild != nil {
			dom.RemoveChild(domChild)
			oldChild = oldChild.NextSibling
			domChild = domChild.NextSibling()
		}
		// Add new children
		for child != nil {
			dom.AppendChild(Construct(child, dom))
			child = child.NextSibling
		}
	}

	return dom
}
