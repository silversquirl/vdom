package vdom_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/vktec/vdom"
	"github.com/vktec/vdom/htmldom"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var testFragment = strings.NewReplacer("\t", "", "\n", "").Replace(`
	<body charset="utf-8">
		<h1>Hello, world!</h1>
		<p>
			foo<br/>
			bar<br/>
			baz<br/>
			quux
		</p>
		<p data-foo="&#34;c&lt;d&#39;&amp;;">
			frob
		</p>
	</body>
`)

func testTree() *html.Node {
	nodes, err := html.ParseFragment(strings.NewReader(testFragment), nil)
	if err != nil {
		panic(err)
	}
	// html -> body
	return nodes[0].LastChild
}

func checkNodes(t *testing.T, thing string, expected, got *html.Node) {
	if nodesEqual(expected, got) {
		return
	}

	b := strings.Builder{}
	html.Render(&b, expected)
	estr := b.String()
	b.Reset()

	html.Render(&b, got)
	gstr := b.String()

	t.Log(got, expected)
	t.Errorf("%s does not match:\nexpected: %s\ngot:      %s", thing, estr, gstr)
}
func nodesEqual(a, b *html.Node) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.Type != b.Type || a.DataAtom != b.DataAtom || a.Data != b.Data || a.Namespace != b.Namespace {
		return false
	}
	if len(a.Attr) != len(b.Attr) {
		return false
	}
	for i := range a.Attr {
		if a.Attr[i] != b.Attr[i] {
			return false
		}
	}
	return nodesEqual(a.NextSibling, b.NextSibling) && nodesEqual(a.FirstChild, b.FirstChild)
}

func TestClone(t *testing.T) {
	tree := testTree()
	clone := vdom.Clone(tree)
	// Change some text
	tree.FirstChild.FirstChild.Data = "Hi everyone!"
	if reflect.DeepEqual(tree, clone) {
		t.Error("tree and clone are equal")
	}
}

func TestConstruct(t *testing.T) {
	node := &html.Node{}
	dom := htmldom.New(node)
	expect := testTree()
	dom.AppendChild(vdom.Construct(expect, dom))
	checkNodes(t, "Generated node", expect, node.FirstChild)
}

func TestPatch(t *testing.T) {
	node := &html.Node{}
	dom := htmldom.New(node)
	tree := testTree()
	var prev *html.Node

	body := tree
	h1 := body.FirstChild
	p_0 := h1.NextSibling
	p_1 := p_0.NextSibling

	vdom.Patch(tree, prev, dom)
	prev = vdom.Clone(tree)
	checkNodes(t, "[init] Patched node", tree, node)

	// Change some text
	h1.FirstChild.Data = "Hi everyone!"

	vdom.Patch(tree, prev, dom)
	prev = vdom.Clone(tree)
	checkNodes(t, "[text] Patched node", tree, node)

	// Change some attributes
	h1.Attr = append(h1.Attr, html.Attribute{Key: "class", Val: "title"})
	body.Attr[0].Val = "ascii"
	p_1.Attr = p_1.Attr[:0]

	vdom.Patch(tree, prev, dom)
	prev = vdom.Clone(tree)
	checkNodes(t, "[attr] Patched node", tree, node)

	// Move some children around
	text := p_0.FirstChild.NextSibling.NextSibling
	p_0.RemoveChild(text)
	p_1.AppendChild(text)

	vdom.Patch(tree, prev, dom)
	prev = vdom.Clone(tree)
	checkNodes(t, "[child] Patched node", tree, node)

	// Change an element's name
	p_0.DataAtom = atom.Div
	p_0.Data = "div"

	vdom.Patch(tree, prev, dom)
	prev = vdom.Clone(tree)
	checkNodes(t, "[name] Patched node", tree, node)
}
