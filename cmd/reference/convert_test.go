package main

import (
	"net/url"
	"strings"
	"testing"
)

// testLinker knows a little of the standard library and the spec's Types section, in which
// also stands an element with the id map_example.
func testLinker() *linker {
	return &linker{
		symbols: map[string]map[string]string{
			"fmt":     {"Println": "func"},
			"strings": {"Builder": "type", "Builder.Grow": "method"},
		},
		articles: map[string]*article{
			"spec": {slug: "spec", ids: map[string]string{"Types": "Types", "Struct_types": "Struct_types", "map_example": "Types"}},
			"mem":  {slug: "mem", ids: map[string]string{}},
		},
	}
}

func convertHTML(t *testing.T, src string) string {
	t.Helper()
	l := testLinker()
	base := &url.URL{Scheme: "https", Host: "go.dev", Path: "/ref/spec"}
	c := &converter{link: func(href string) string { return l.target(base, href) }}
	return c.markdown(parseHTML(src).children)
}

func TestConvert(t *testing.T) {
	for _, tc := range []struct{ name, html, want string }{
		{"paragraphs left open", "<p>one\n<p>two", "one\n\ntwo"},
		{"inline code, emphasis, strong", "<p>a <code>x  y</code> <i>it</i> <em> em </em><b>b</b> <strong>s</strong> <dfn>d</dfn></p>",
			"a `x y` *it* *em* **b** **s** *d*"},
		{"code holding backticks", "<p><code>``</code> and <code>`a</code></p>", "``` `` ``` and `` `a ``"},
		{"entities and Markdown's characters", "<p>a*b_c [d] &lt;e&gt; ~f &amp; g|h \\ `i`</p>",
			`a\*b\_c \[d\] \<e\> \~f \& g\|h \\ ` + "\\`i\\`"},
		{"lines that would start blocks", "<p>1. not a list</p><p># not a heading</p><p>- nor this<br>+ nor that</p>",
			"1\\. not a list\n\n\\# not a heading\n\n\\- nor this\\\n\\+ nor that"},
		{"a line break", "<p>one<br>\n  two<br></p>", "one\\\ntwo"},
		{"subscripts and superscripts", "<p>a<sub>i</sub>, x<sub>10</sub>, ≡<sub>A</sub>, 2<sup>32</sup>, 2<sup><i>n</i></sup></p>", "aᵢ, x₁₀, ≡A, 2^32^, 2^n^"},
		{"pictures", `<p>a <img src="gopher.png" alt="a gopher"> b</p><picture><img src="x.png"></picture><p><video src="v.mp4"></video></p>`, "a b"},
		{"Go", "<pre>\nfunc main() {\n\tx := 1 &lt;&lt; 2 // <b>shifted</b>\n}\n</pre>", "```go\nfunc main() {\n\tx := 1 << 2 // shifted\n}\n```"},
		{"the spec's grammar", `<pre class="ebnf">
Block = "{" StatementList "}" .
</pre><pre class="grammar">
|   alternation
</pre>`, "```\nBlock = \"{\" StatementList \"}\" .\n```\n\n```\n|   alternation\n```"},
		{"a shell session", "<pre>\n# set up\n$ go mod init example/hello\n</pre>", "```\n# set up\n$ go mod init example/hello\n```"},
		{"another language goldmark names", `<pre><code class="language-sh">ls -l
</code></pre>`, "```\nls -l\n```"},
		{"code holding a fence", "<pre>x := \"```\"</pre>", "````go\nx := \"```\"\n````"},
		{"a list, items left open", "<ul><li>one<li>two <code>2</code></ul>", "- one\n- two `2`"},
		{"a numbered list from 3", `<ol start="3"><li><p>three</p></li><li><p>four</p><pre>x := 4</pre></li></ol>`,
			"3. three\n\n4. four\n\n   ```go\n   x := 4\n   ```"},
		{"a list in a list", "<ul><li>outer<ul><li>inner</li></ul></li><li>next</li></ul>", "- outer\n\n  - inner\n\n- next"},
		{"a definition list", "<dl><dt>Indentation</dt><dd>Use tabs.</dd><dt>Width</dt><dd><p>None.</p></dd></dl>",
			"- **Indentation**\n\n  Use tabs.\n\n- **Width**\n\n  None."},
		{"a block quote", "<blockquote><p>one</p><p>two</p></blockquote>", "> one\n>\n> two"},
		{"a note", `<div class="Note"><strong>Note:</strong> This is old.</div>`, "> [!NOTE]\n> This is old."},
		{"a rule", `<p class="rule">A send happens before the receive.</p>`, "> A send happens before the receive."},
		{"a table", "<table><tr><th>Op</th><th>Means</th></tr><tr><td><code>|</code></td><td>or<br>else</td></tr></table>",
			"| Op | Means |\n| --- | --- |\n| `|` | or<br>else |"},
		{"a heading inside a block", "<blockquote><h3>Aside</h3><p>text</p></blockquote>", "> **Aside**\n>\n> text"},
		{"links within the article", `<p><a href="#Types">types</a>, <a href="#map_example">maps</a>, <a href="#Nowhere">nowhere</a></p>`,
			"[types](godoc-article:spec#Types), [maps](godoc-article:spec#Types), [nowhere](godoc-article:spec)"},
		{"links to other articles", `<p><a href="/ref/mem">mem</a>, <a href="https://golang.org/ref/spec#Struct_types">structs</a>, <a href="/doc/faq#why">faq</a></p>`,
			"[mem](godoc-article:mem), [structs](godoc-article:spec#Struct_types), faq"},
		{"links to packages", `<p><a href="https://pkg.go.dev/strings#Builder.Grow">Grow</a>, <a href="/pkg/fmt/#Println">Println</a>, ` +
			`<a href="http://golang.org/pkg/fmt/#hdr-Printing">printing</a>, <a href="/pkg/strings/#Builder.Len">Len</a>, ` +
			`<a href="https://pkg.go.dev/strings@go1.22">strings</a>, <a href="/pkg/nowhere/">nowhere</a>, <a href="/pkg/">all</a></p>`,
			"[Grow](godoc:strings.Builder.Grow), [Println](godoc:fmt.Println), [printing](godoc:fmt), [Len](godoc:strings.Builder.Len), [strings](godoc:strings), nowhere, all"},
		{"links out of the app", `<p><a href="https://example.com/x">an essay</a>, <a href="/blog/go1.22">the blog</a>, <a href="mailto:a@b.c">mail</a>, <a name="here">here</a></p>`,
			"an essay, the blog, mail, here"},
		{"code that links", `<p><a href="#Types"><code>T</code></a> and <code><a href="/pkg/fmt/#Println">Println</a></code></p>`,
			"[`T`](godoc-article:spec#Types) and [`Println`](godoc:fmt.Println)"},
		{"comments, styles and scripts", "<!-- a comment --><style>p { color: red }</style><script>alert(1)</script><p>text</p>", "text"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := convertHTML(t, tc.html); got != tc.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

// An article splits at each heading, of whatever level and however deep in the page, and each
// element id belongs to the section it is in.
func TestSplit(t *testing.T) {
	a := &article{}
	a.split(parseHTML(`<div class="Note">Old.</div>
<p id="top">Lead.</p>
<div><h2 id="Types">Types</h2><p>Types are <a id="kinds">kinds</a>.</p><pre class="ebnf">
Type     = TypeName .
TypeName = identifier .
</pre>
<h3 id="Struct_types">Struct <code>types</code></h3></div>
<h3>No id here</h3><p>Body.</p><h4 id="Types">Twice</h4>`))
	var got []string
	for _, s := range a.sections {
		got = append(got, s.anchor+" "+string(rune('0'+s.level))+" "+s.title)
	}
	want := []string{"Types 2 Types", "Struct_types 3 Struct types", "No_id_here 3 No id here", "Twice 4 Twice"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("sections = %q, want %q", got, want)
	}
	c := &converter{}
	got = []string{c.markdown(a.lead)}
	for _, s := range a.sections {
		got = append(got, c.markdown(s.nodes))
	}
	want = []string{"> [!NOTE]\n> Old.\n\nLead.", "Types are kinds.\n\n```\nType     = TypeName .\nTypeName = identifier .\n```", "", "Body.", ""}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("the lead and the sections' bodies = %q, want %q", got, want)
	}
	for id, want := range map[string]string{"top": "", "Types": "Types", "kinds": "Types", "TypeName": "Types", "Struct_types": "Struct_types", "No_id_here": "No_id_here"} {
		if got, ok := a.ids[id]; !ok || got != want {
			t.Errorf("ids[%s] = %q, %v; want %q", id, got, ok, want)
		}
	}
}

func TestCodeLines(t *testing.T) {
	src := "package main\n\nfunc main() {\n\tprintln(1) // OMIT\n\tprintln(2)\n}\n"
	for _, tc := range []struct {
		args []any
		want string
	}{
		{nil, "package main\n\nfunc main() {\n    println(2)\n}\n"},
		{[]any{"/^func/", "/^}/"}, "func main() {\n    println(2)\n}\n"},
		{[]any{"/println.2/"}, "    println(2)\n"},
		{[]any{"/^func/", "$", "HLx"}, "func main() {\n    println(2)\n}\n"},
		{[]any{3, 6}, "func main() {\n    println(2)\n}\n"},
	} {
		got, err := codeLines(src, tc.args)
		if err != nil || got != tc.want {
			t.Errorf("codeLines(%v) = %q, %v; want %q", tc.args, got, err, tc.want)
		}
	}
	if _, err := codeLines(src, []any{"/nothing/"}); err == nil {
		t.Error("a pattern that matches no line was accepted")
	}
}
