package hypergo

import (
	"testing"
)

const (
	noesc = `"Fran & Freddie's Diner" <tasty@example.com>`
	esc   = `&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt;`
)

func TestRender(t *testing.T) {
	tests := []struct {
		name       string
		node       func() *Node
		want       string
		wantIndent string
	}{
		{
			name: "helloWorld",
			node: func() *Node {
				return P().Text("Hello World!")
			},
			want:       "<p>Hello World!</p>",
			wantIndent: "<p>Hello World!</p>\n",
		},
		{
			name: "fullHTML",
			node: func() *Node {
				head := Head().Title().Text("Title")
				body := Body().H1().Text("Heading")
				return HTML5(Html().Append(head, body))
			},
			want: "<!DOCTYPE html><html><head><title>Title</title></head><body><h1>Heading</h1></body></html>",
			wantIndent: `<!DOCTYPE html>
<html>
	<head>
		<title>Title</title>
	</head>
	<body>
		<h1>Heading</h1>
	</body>
</html>
`,
		},
		{
			name: "for",
			node: func() *Node {
				ul := Ul()
				for i := 0; i < 3; i++ {
					ul.Li().Textv("Count: ", i)
				}
				return ul
			},
			want: "<ul><li>Count: 0</li><li>Count: 1</li><li>Count: 2</li></ul>",
			wantIndent: `<ul>
	<li>Count: 0</li>
	<li>Count: 1</li>
	<li>Count: 2</li>
</ul>
`,
		},
		{
			name: "append",
			node: func() *Node {
				ul := Ul().Append(
					Li().Text("0"),
					Li().Text("1"),
				)
				ul.Append(
					Li().Text("2"),
					Li().Text("3"),
				)
				ul.Li().Text("4")
				return ul
			},
			want: `<ul><li>0</li><li>1</li><li>2</li><li>3</li><li>4</li></ul>`,
			wantIndent: `<ul>
	<li>0</li>
	<li>1</li>
	<li>2</li>
	<li>3</li>
	<li>4</li>
</ul>
`,
		},
		{
			name: "appendChild",
			node: func() *Node {
				child := Li()
				ul := Ul().Append(
					child,
					Li().Text("1"),
				)
				child.Append(
					Span().Text("foo"),
				)
				return ul
			},
			want: `<ul><li><span>foo</span></li><li>1</li></ul>`,
			wantIndent: `<ul>
	<li>
		<span>foo</span>
	</li>
	<li>1</li>
</ul>
`,
		},
		{
			name: "escape",
			node: func() *Node {
				return P("attr="+noesc, noesc, "id=foo").P().Text(noesc)
			},
			want: `<p attr="` + esc + `" ` + esc + ` id="foo"><p>` + esc + `</p></p>`,
			wantIndent: `<p attr="` + esc + `" ` + esc + ` id="foo">
	<p>` + esc + `</p>
</p>
`,
		},
		{
			name: "attrs",
			node: func() *Node {
				return Input("id=foo bar", " class = foo    bar ", "disabled ")
			},
			want: `<input id="foo bar"  class =" foo    bar " disabled >`,
			wantIndent: `<input id="foo bar"  class =" foo    bar " disabled >
`,
		},
		{
			name: "emptyAttr",
			node: func() *Node {
				return Input("id=", "=foo", "disabled")
			},
			want: `<input id="" foo disabled>`,
			wantIndent: `<input id="" foo disabled>
`,
		},
	}

	for _, tt := range tests {
		got := tt.node().Render()
		if got != tt.want {
			t.Errorf("%s: got %q; want %q", tt.name, got, tt.want)
		}
		gotIndent := tt.node().RenderIndent("\t")
		if gotIndent != tt.wantIndent {
			t.Errorf("%s: got %q; want %q", tt.name, gotIndent, tt.wantIndent)
		}
	}
}
