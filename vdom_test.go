package vdom_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/vktec/vdom"
	"github.com/vktec/vdom/testdom"
)

func testTree() vdom.Node {
	return &vdom.Element{
		Name:  "body",
		Attrs: map[string]string{"charset": "utf-8"},
		Children: []vdom.Node{
			&vdom.Element{Name: "h1", Children: []vdom.Node{&vdom.Text{"Hello, world!"}}},
			&vdom.Element{
				Name: "p",
				Children: []vdom.Node{
					&vdom.Text{"foo"},
					&vdom.Element{Name: "br"},
					&vdom.Text{"bar"},
					&vdom.Element{Name: "br"},
					&vdom.Text{"baz"},
					&vdom.Element{Name: "br"},
					&vdom.Text{"quux"},
				},
			},
			&vdom.Element{
				Name:     "p",
				Attrs:    map[string]string{"data-foo": `"c<d'&;`},
				Children: []vdom.Node{&vdom.Text{"frob"}},
			},
		},
	}
}

func TestClone(t *testing.T) {
	tree := testTree()
	clone := tree.Clone()
	// Change some text
	tree.(*vdom.Element). // body
				Children[0].(*vdom.Element). // h1
				Children[0].(*vdom.Text).    // text
				Text = "Hi everyone!"
	if reflect.DeepEqual(tree, clone) {
		t.Error("tree and clone are equal")
	}
}

func TestHTML(t *testing.T) {
	r := strings.NewReplacer("\t", "", "\n", "")
	expect := r.Replace(`
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

	html := testTree().HTML()
	if html != expect {
		t.Errorf("Generated HTML does not match:\nexpected: %s\ngot:      %s", expect, html)
	}
}

func TestConstruct(t *testing.T) {
	dom := testdom.NewTestDOM()
	expect := testTree()
	node := *expect.Construct(dom).(*testdom.TestDOM).Node
	if !reflect.DeepEqual(node, expect) {
		t.Errorf("Generated node does not match:\nexpected: %s\ngot:      %s", expect.HTML(), node.HTML())
	}
}

func TestPatch(t *testing.T) {
	dom := testdom.NewTestDOM()
	tree := testTree()
	var prev vdom.Node

	body := tree.(*vdom.Element)
	h1 := body.Children[0].(*vdom.Element)
	p_0 := body.Children[1].(*vdom.Element)
	p_1 := body.Children[2].(*vdom.Element)

	tree.Patch(dom, prev)
	prev = tree.Clone()
	if !reflect.DeepEqual(*dom.Node, tree) {
		t.Errorf("Patched node does not match:\nexpected: %s\ngot:      %s", tree.HTML(), (*dom.Node).HTML())
	}

	// Change some text
	h1.Children[0].(*vdom.Text).Text = "Hi everyone!"

	tree.Patch(dom, prev)
	prev = tree.Clone()
	if !reflect.DeepEqual(*dom.Node, tree) {
		t.Errorf("Patched node does not match:\nexpected: %s\ngot:      %s", tree.HTML(), (*dom.Node).HTML())
	}

	// Change some attributes
	h1.Attrs = map[string]string{"class": "title"}
	body.Attrs["charset"] = "ascii"
	delete(p_1.Attrs, "data-foo")

	tree.Patch(dom, prev)
	prev = tree.Clone()
	if !reflect.DeepEqual(*dom.Node, tree) {
		t.Errorf("Patched node does not match:\nexpected: %s\ngot:      %s", tree.HTML(), (*dom.Node).HTML())
	}

	// Move some children around
	text := p_0.Children[2]
	p_0.Children = append(p_0.Children[:2], p_0.Children[3:]...)
	p_1.Children = append(p_1.Children, text)

	tree.Patch(dom, prev)
	prev = tree.Clone()
	if !reflect.DeepEqual(*dom.Node, tree) {
		t.Errorf("Patched node does not match:\nexpected: %s\ngot:      %s", tree.HTML(), (*dom.Node).HTML())
	}

	// Change an element's name
	p_0.Name = "div"

	tree.Patch(dom, prev)
	prev = tree.Clone()
	if !reflect.DeepEqual(*dom.Node, tree) {
		t.Errorf("Patched node does not match:\nexpected: %s\ngot:      %s", tree.HTML(), (*dom.Node).HTML())
	}
}
