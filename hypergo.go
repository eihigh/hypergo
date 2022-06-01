// hypergo builds HTML in pure Go.
package hypergo

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

var (
	sp      = []byte(" ")
	lt      = []byte("<")
	gt      = []byte(">")
	eqQuot  = []byte("=\"")
	quot    = []byte("\"")
	ltSlash = []byte("</")
	gtLf    = []byte(">\n")
	lf      = []byte("\n")
)

type Node struct {
	root     *Node
	t        string // Tag or Text
	attrs    []attr
	children []*Node
}

type attr struct {
	key, value string
}

func (a *attr) isBoolean() bool { return a.key == "" }

func IsEmptyTag(tag string) bool {
	switch tag {
	case "!DOCTYPE", "area", "base", "basefont", "br", "col", "frame",
		"hr", "img", "input", "isindex", "link", "meta", "param":
		return true
	}
	return false
}

func HTML5(html *Node) *Node {
	root := Element("")
	root.Append(
		Element("!DOCTYPE", "html"),
		html,
	)
	return root
}

func Element(tag string, attrs ...string) *Node {
	n := &Node{
		t:        tag,
		children: make([]*Node, 0, 1),
		attrs:    make([]attr, 0, len(attrs)),
	}
	n.root = n

	for _, a := range attrs {
		before, after, found := strings.Cut(a, "=")
		if found {
			n.attrs = append(n.attrs, attr{before, after})
		} else {
			n.attrs = append(n.attrs, attr{"", before})
		}
	}
	return n
}

func Text(text string) *Node {
	n := &Node{
		t:        text,
		children: nil,
		attrs:    nil,
	}
	n.root = n
	return n
}

func (n *Node) IsText() bool {
	return n.children == nil
}

func (n *Node) Append(childNodes ...*Node) *Node {
	for _, c := range childNodes {
		if c.root == n.root {
			panic("cannot share root nodes")
		}
		n.children = append(n.children, c.root)
	}
	return n
}

func (n *Node) Element(tag string, attrs ...string) *Node {
	c := Element(tag, attrs...)
	c.root = n.root
	n.children = append(n.children, c)
	return c
}

func (n *Node) Text(text string) *Node {
	c := Text(text)
	c.root = n.root
	n.children = append(n.children, c)
	return c
}

func (n *Node) Textv(args ...any) *Node {
	return n.Text(fmt.Sprint(args...))
}

func (n *Node) Textf(format string, args ...any) *Node {
	return n.Text(fmt.Sprintf(format, args...))
}

// ============================================================
// Rendering
// ============================================================

func (n *Node) Render() string {
	b := strings.Builder{}
	n.root.render(&b)
	return b.String()
}

func (n *Node) FRender(w io.Writer) {
	n.root.render(w)
}

func (n *Node) RenderIndent(indentStr string) string {
	b := strings.Builder{}
	n.root.renderIndent(&b, 0, []byte(indentStr))
	return b.String()
}

func (n *Node) FRenderIndent(w io.Writer, indentStr string) {
	n.root.renderIndent(w, 0, []byte(indentStr))
}

func (n *Node) render(w io.Writer) {
	if n == nil { // Zero value of *Node
		return
	}

	if n.t == "" { // Fragment node
		for _, c := range n.children {
			c.render(w)
		}
		return
	}

	if n.IsText() { // Text node
		template.HTMLEscape(w, []byte(n.t))
		return
	}

	w.Write(lt)
	template.HTMLEscape(w, []byte(n.t))

	for _, a := range n.attrs {
		w.Write(sp)
		if a.isBoolean() {
			template.HTMLEscape(w, []byte(a.value)) // <tag value
		} else {
			template.HTMLEscape(w, []byte(a.key))
			w.Write(eqQuot)
			template.HTMLEscape(w, []byte(a.value))
			w.Write(quot) // <tag value key="value"
		}
	}
	w.Write(gt) // <tag value key="value">

	for _, c := range n.children {
		c.render(w)
	} // <tag value key="value"> children

	if len(n.children) == 0 && IsEmptyTag(n.t) {
		return
	}

	w.Write(ltSlash)
	template.HTMLEscape(w, []byte(n.t))
	w.Write(gt) // <tag value key="value"> children </tag>
}

func (n *Node) renderIndent(w io.Writer, indent int, indentStr []byte) {
	if n == nil { // Zero value of *Node
		return
	}

	if n.t == "" { // Fragment node
		for _, c := range n.children {
			c.renderIndent(w, indent, indentStr)
		}
		return
	}

	for i := 0; i < indent; i++ {
		w.Write(indentStr)
	}

	if n.IsText() { // Text node
		template.HTMLEscape(w, []byte(n.t))
		w.Write(lf)
		return
	}

	w.Write(lt)
	template.HTMLEscape(w, []byte(n.t)) // \t\t<tag

	for _, a := range n.attrs {
		w.Write(sp)
		if a.isBoolean() {
			template.HTMLEscape(w, []byte(a.value)) // \t\t<tag value
		} else {
			template.HTMLEscape(w, []byte(a.key))
			w.Write(eqQuot)
			template.HTMLEscape(w, []byte(a.value))
			w.Write(quot) // \t\t<tag value key="value"
		}
	}
	w.Write(gt) // \t\t<tag value key="value">

	switch {
	case len(n.children) == 0 && IsEmptyTag(n.t):
		w.Write(lf) // \t\t<tag value key="value">\n
		return

	case len(n.children) == 0 || (len(n.children) == 1 && n.children[0].IsText()):
		for _, c := range n.children {
			c.render(w)
		} // \t\t<tag value key="value"> children

	default:
		w.Write(lf) // \t\t<tag value key="value">\n

		for _, c := range n.children {
			c.renderIndent(w, indent+1, indentStr)
		} // \t\t<tag value key="value">\n\t\t\t children \n

		for i := 0; i < indent; i++ {
			w.Write(indentStr)
		} // \t\t<tag value key="value">\n\t\t\t children \n\t\t
	}

	w.Write(ltSlash)
	template.HTMLEscape(w, []byte(n.t))
	w.Write(gtLf)
	// \t\t<tag value key="value"> children </tag>\n
	// OR \t\t<tag value key="value">\n\t\t\t children \n\t\t</tag>\n
}
