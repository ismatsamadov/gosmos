package main

// The articles come as HTML: go.dev's documents are written in it by hand, or in Markdown that
// goldmark turns into it. This file parses that HTML into a tree and writes the tree as the
// Markdown the app renders (app/Sources/GosmosKit/Readme/Readme.swift: CommonMark with GitHub's
// tables, callouts like > [!NOTE], and ^superscript^).

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

// A node is an element of a parsed document, or a run of its text.
type node struct {
	tag      string            // the element's name in lower case; "" for text
	attr     map[string]string // values with their entities decoded
	text     string            // a text node's text, entities decoded
	children []*node
}

func tagSet(names string) map[string]bool {
	m := map[string]bool{}
	for _, n := range strings.Fields(names) {
		m[n] = true
	}
	return m
}

var (
	voidTags  = tagSet("area base br col embed hr img input link meta source track wbr")
	blockTags = tagSet("address article aside blockquote center dd details div dl dt figcaption figure footer form " +
		"h1 h2 h3 h4 h5 h6 header hr li main nav ol p pre section table tbody td tfoot th thead tr ul")
	// Elements whose content never shows: pictures and other media, scripts, forms.
	droppedTags = tagSet("audio button canvas embed iframe img input noscript object picture script select style svg template textarea video")
	// Elements that only hold other blocks, whose blocks are read as the document's own.
	containerTags = tagSet("#document article body center div footer header html main nav section")
	headingLevel  = map[string]int{"h1": 1, "h2": 2, "h3": 3, "h4": 4, "h5": 5, "h6": 6}
	// HTML's button scope: an open p element above any of these stays open.
	pScope = tagSet("#document applet button caption html marquee object table td template th")
)

// parseHTML parses a document the way a browser would for the HTML the articles use: elements
// HTML lets a document leave open (p, li, dt, dd, tr, td, th) are closed where the next one
// starts, an end tag closes whatever is open inside it, and comments, scripts and styles are
// dropped.
func parseHTML(src string) *node {
	root := &node{tag: "#document"}
	stack := []*node{root}
	top := func() *node { return stack[len(stack)-1] }
	// closeOpen closes the innermost open element named in names, and all inside it, unless
	// an element named in scope comes first.
	closeOpen := func(names, scope map[string]bool) {
		for i := len(stack) - 1; i > 0; i-- {
			if names[stack[i].tag] {
				stack = stack[:i]
				return
			}
			if scope[stack[i].tag] {
				return
			}
		}
	}
	text := func(s string) {
		if s == "" {
			return
		}
		p := top()
		if k := len(p.children); k > 0 && p.children[k-1].tag == "" {
			p.children[k-1].text += s
			return
		}
		p.children = append(p.children, &node{text: s})
	}
	for i := 0; i < len(src); {
		lt := strings.IndexByte(src[i:], '<')
		if lt != 0 {
			if lt < 0 {
				lt = len(src) - i
			}
			text(html.UnescapeString(src[i : i+lt]))
			i += lt
			continue
		}
		rest := src[i:]
		switch {
		case strings.HasPrefix(rest, "<!--"):
			end := strings.Index(rest[4:], "-->")
			if end < 0 {
				return root
			}
			i += 4 + end + 3
		case strings.HasPrefix(rest, "</") && len(rest) > 2 && isLetter(rest[2]):
			end := strings.IndexByte(rest, '>')
			if end < 0 {
				return root
			}
			closeOpen(map[string]bool{strings.ToLower(strings.TrimSpace(rest[2:end])): true}, nil)
			i += end + 1
		case len(rest) > 1 && isLetter(rest[1]):
			name, attr, selfClosing, n := parseTag(rest)
			if n == 0 {
				text("<")
				i++
				break
			}
			i += n
			if name == "script" || name == "style" || name == "textarea" {
				end := strings.Index(strings.ToLower(src[i:]), "</"+name)
				if end < 0 {
					return root
				}
				i += end
				if gt := strings.IndexByte(src[i:], '>'); gt >= 0 {
					i += gt + 1
				}
				break
			}
			switch name {
			case "li":
				closeOpen(tagSet("li"), tagSet("ul ol"))
			case "dt", "dd":
				closeOpen(tagSet("dt dd"), tagSet("dl"))
			case "tr":
				closeOpen(tagSet("tr"), tagSet("table thead tbody tfoot"))
			case "td", "th":
				closeOpen(tagSet("td th"), tagSet("tr table"))
			}
			if blockTags[name] {
				closeOpen(tagSet("p"), pScope)
			}
			el := &node{tag: name, attr: attr}
			top().children = append(top().children, el)
			if !voidTags[name] && !selfClosing {
				stack = append(stack, el)
			}
		case strings.HasPrefix(rest, "<!") || strings.HasPrefix(rest, "<?"):
			end := strings.IndexByte(rest, '>')
			if end < 0 {
				return root
			}
			i += end + 1
		default:
			text("<")
			i++
		}
	}
	return root
}

func isLetter(c byte) bool { return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' }

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' }

// parseTag reads the start tag at the beginning of s and returns its name, attributes, whether
// it ends with />, and its length; a length of 0 if s holds no whole tag.
func parseTag(s string) (name string, attr map[string]string, selfClosing bool, n int) {
	i := 1
	for i < len(s) && (isLetter(s[i]) || '0' <= s[i] && s[i] <= '9' || s[i] == '-' || s[i] == ':') {
		i++
	}
	name = strings.ToLower(s[1:i])
	attr = map[string]string{}
	for i < len(s) {
		switch {
		case isSpace(s[i]):
			i++
		case s[i] == '>':
			return name, attr, selfClosing, i + 1
		case s[i] == '/':
			selfClosing = true
			i++
		default:
			selfClosing = false
			j := i
			for j < len(s) && !isSpace(s[j]) && s[j] != '=' && s[j] != '>' && !(s[j] == '/' && j+1 < len(s) && s[j+1] == '>') {
				j++
			}
			key := strings.ToLower(s[i:j])
			i = j
			for i < len(s) && isSpace(s[i]) {
				i++
			}
			value := ""
			if i < len(s) && s[i] == '=' {
				i++
				for i < len(s) && isSpace(s[i]) {
					i++
				}
				if i < len(s) && (s[i] == '"' || s[i] == '\'') {
					end := strings.IndexByte(s[i+1:], s[i])
					if end < 0 {
						return "", nil, false, 0
					}
					value = s[i+1 : i+1+end]
					i += end + 2
				} else {
					k := i
					for k < len(s) && !isSpace(s[k]) && s[k] != '>' {
						k++
					}
					value = s[i:k]
					i = k
				}
			}
			if key != "" {
				attr[key] = html.UnescapeString(value)
			}
			if i == j && key == "" { // a character no attribute can start with
				i++
			}
		}
	}
	return "", nil, false, 0
}

// hasClass reports whether the element n has class c.
func (n *node) hasClass(c string) bool {
	for _, x := range strings.Fields(n.attr["class"]) {
		if strings.EqualFold(x, c) {
			return true
		}
	}
	return false
}

func (n *node) isBlock() bool { return blockTags[n.tag] || containerTags[n.tag] }

// A converter writes parsed HTML as Markdown. link returns the godoc: or godoc-article: URL an
// href leads to, or "" when it leads out of the reference: such a link becomes its text.
type converter struct {
	link func(href string) string
}

// markdown writes nodes as Markdown blocks separated by blank lines.
func (c *converter) markdown(nodes []*node) string {
	return strings.Join(c.blocks(nodes), "\n\n")
}

// blocks writes nodes as Markdown blocks: each block element as one or more blocks, and each
// run of text and inline elements between them as a paragraph.
func (c *converter) blocks(nodes []*node) []string {
	var out []string
	var inline []*node
	flush := func() {
		if p := c.paragraph(inline); p != "" {
			out = append(out, p)
		}
		inline = nil
	}
	for _, n := range nodes {
		if n.tag == "" || !n.isBlock() {
			inline = append(inline, n)
			continue
		}
		flush()
		out = append(out, c.block(n)...)
	}
	flush()
	return out
}

// block writes one block element.
func (c *converter) block(n *node) []string {
	switch n.tag {
	case "p":
		p := c.paragraph(n.children)
		switch {
		case p == "":
			return nil
		case n.hasClass("rule"): // the memory model's rules, set apart on go.dev
			return []string{quote(p)}
		}
		return []string{p}
	case "pre":
		return []string{c.pre(n)}
	case "ul", "ol":
		return c.list(n)
	case "dl":
		return c.definitions(n)
	case "blockquote":
		if inner := c.markdown(n.children); inner != "" {
			return []string{quote(inner)}
		}
		return nil
	case "h1", "h2", "h3", "h4", "h5", "h6": // a heading inside another block
		if p := c.paragraph(n.children); p != "" {
			return []string{"**" + p + "**"}
		}
		return nil
	case "hr":
		return []string{"---"}
	case "table", "thead", "tbody", "tfoot":
		return c.table(n)
	}
	if n.hasClass("note") { // Effective Go's note, drawn as a callout
		children := n.children
		for len(children) > 0 && children[0].tag == "" && strings.TrimSpace(children[0].text) == "" {
			children = children[1:]
		}
		if len(children) > 0 && (children[0].tag == "strong" || children[0].tag == "b") &&
			strings.TrimSuffix(strings.TrimSpace(plain(children[0])), ":") == "Note" {
			children = children[1:]
		}
		if inner := c.markdown(children); inner != "" {
			return []string{"> [!NOTE]\n" + quote(inner)}
		}
		return nil
	}
	return c.blocks(n.children) // a container, or li, dt, dd, tr, td or th out of place
}

// paragraph writes inline content as one paragraph: HTML's runs of white space become one
// space, and a line that would start a Markdown block is escaped.
func (c *converter) paragraph(nodes []*node) string {
	s := reSpaces.ReplaceAllString(c.inline(nodes), " ")
	var lines []string
	for line := range strings.SplitSeq(strings.ReplaceAll(s, "\x00", "\n"), "\n") {
		if line = strings.TrimSpace(line); line != "" && line != `\` {
			lines = append(lines, escapeLineStart(line))
		}
	}
	return strings.TrimSuffix(strings.Join(lines, "\n"), `\`) // a line break at the very end
}

var reSpaces = regexp.MustCompile(`[ \t\r\n\f]+`)

// lineBreak is a Markdown line break for <br>; its NUL stands for the newline until paragraph
// has collapsed the white space around it.
const lineBreak = "\\\x00"

// inline writes text and inline elements as inline Markdown. Block elements inside them are
// read as inline content.
func (c *converter) inline(nodes []*node) string {
	var b strings.Builder
	for _, n := range nodes {
		switch n.tag {
		case "":
			b.WriteString(escapeText(n.text))
		case "br":
			b.WriteString(lineBreak)
		case "code", "tt", "kbd", "samp":
			code := codeSpan(strings.TrimSpace(reSpaces.ReplaceAllString(plain(n), " ")))
			if target := c.innerLink(n); target != "" && code != "" {
				code = "[" + code + "](" + target + ")"
			}
			b.WriteString(code)
		case "i", "em", "dfn", "var", "cite":
			b.WriteString(emphasize("*", c.inline(n.children)))
		case "b", "strong":
			b.WriteString(emphasize("**", c.inline(n.children)))
		case "a":
			inner := c.inline(n.children)
			if target := c.target(n); target != "" && strings.TrimSpace(inner) != "" {
				lead, text, trail := splitSpace(inner)
				b.WriteString(lead + "[" + text + "](" + target + ")" + trail)
			} else {
				b.WriteString(inner)
			}
		case "sub":
			b.WriteString(escapeText(subscript(plain(n))))
		case "sup":
			s := strings.TrimSpace(plain(n))
			if s != "" && !strings.ContainsAny(s, " \t\n^") {
				s = "^" + s + "^"
			}
			b.WriteString(escapeText(s))
		default:
			if !droppedTags[n.tag] {
				b.WriteString(c.inline(n.children))
			}
		}
	}
	return b.String()
}

// target is where the link a leads, or "".
func (c *converter) target(a *node) string {
	href, ok := a.attr["href"]
	if !ok || c.link == nil {
		return ""
	}
	return c.link(href)
}

// innerLink is where a link inside the code element n leads, for <code><a href=…>x</a></code>.
func (c *converter) innerLink(n *node) string {
	for _, ch := range n.children {
		if ch.tag == "a" {
			return c.target(ch)
		}
	}
	return ""
}

// emphasize wraps s in mark, outside any white space at its ends, which would keep the mark
// from working.
func emphasize(mark, s string) string {
	lead, text, trail := splitSpace(s)
	if text == "" {
		return s
	}
	return lead + mark + text + mark + trail
}

func splitSpace(s string) (lead, text, trail string) {
	text = strings.TrimLeft(s, " \t\r\n\f")
	lead = s[:len(s)-len(text)]
	trimmed := strings.TrimRight(text, " \t\r\n\f")
	return lead, trimmed, text[len(trimmed):]
}

// codeSpan writes s as inline code, fenced with more backticks than it holds in a row.
func codeSpan(s string) string {
	if s == "" {
		return ""
	}
	fence := "`"
	for strings.Contains(s, fence) {
		fence += "`"
	}
	if s[0] == '`' || s[len(s)-1] == '`' {
		s = " " + s + " "
	}
	return fence + s + fence
}

// escapeText escapes the characters that would read as Markdown markup in a line of text.
var escapeText = strings.NewReplacer(`\`, `\\`, "`", "\\`", "*", `\*`, "_", `\_`, "[", `\[`, "]", `\]`,
	"<", `\<`, ">", `\>`, "~", `\~`, "&", `\&`, "|", `\|`).Replace

var reOrdered = regexp.MustCompile(`^(\d+)([.)])`)

// escapeLineStart escapes what would start a heading, list, block quote or setext underline
// at the start of a line.
func escapeLineStart(line string) string {
	if line == "" {
		return line
	}
	switch line[0] {
	case '#', '+', '-', '=':
		return `\` + line
	}
	return reOrdered.ReplaceAllString(line, `$1\$2`)
}

// subscripts are the Unicode subscripts of the characters that have one.
var subscripts = map[rune]rune{
	'0': '₀', '1': '₁', '2': '₂', '3': '₃', '4': '₄', '5': '₅', '6': '₆', '7': '₇', '8': '₈', '9': '₉',
	'a': 'ₐ', 'e': 'ₑ', 'h': 'ₕ', 'i': 'ᵢ', 'j': 'ⱼ', 'k': 'ₖ', 'l': 'ₗ', 'm': 'ₘ', 'n': 'ₙ', 'o': 'ₒ',
	'p': 'ₚ', 'r': 'ᵣ', 's': 'ₛ', 't': 'ₜ', 'u': 'ᵤ', 'v': 'ᵥ', 'x': 'ₓ', '+': '₊', '-': '₋', '=': '₌',
	'(': '₍', ')': '₎',
}

// subscript writes s in Unicode subscripts, as a₁ for a<sub>1</sub>, when each of its
// characters has one, and as it is otherwise.
func subscript(s string) string {
	s = strings.TrimSpace(s)
	var b strings.Builder
	for _, r := range s {
		sub, ok := subscripts[r]
		if !ok {
			return s
		}
		b.WriteRune(sub)
	}
	return b.String()
}

// plain is the text of n, without markup; subscripts as subscript does them.
func plain(n *node) string {
	switch {
	case n.tag == "":
		return n.text
	case droppedTags[n.tag]:
		return ""
	case n.tag == "br":
		return "\n"
	}
	var b strings.Builder
	for _, ch := range n.children {
		if ch.tag == "sub" {
			b.WriteString(subscript(plain(ch)))
		} else {
			b.WriteString(plain(ch))
		}
	}
	return b.String()
}

// pre writes a <pre> element as a fenced code block: ```go for Go, and a plain fence for the
// spec's grammar (class ebnf or grammar), shell sessions and the other languages goldmark
// names.
func (c *converter) pre(n *node) string {
	text := plain(n)
	text = strings.TrimPrefix(strings.TrimPrefix(text, "\r"), "\n") // HTML drops the newline after <pre>
	text = strings.TrimRight(text, " \t\r\n")
	for strings.HasPrefix(text, "\n") {
		text = text[1:]
	}
	lang := "go"
	class := n.attr["class"]
	for _, ch := range n.children {
		if ch.tag == "code" {
			class += " " + ch.attr["class"]
		}
	}
	for _, x := range strings.Fields(class) {
		switch {
		case x == "ebnf" || x == "grammar":
			lang = ""
		case strings.HasPrefix(x, "language-") && x != "language-go":
			lang = ""
		}
	}
	if strings.HasPrefix(text, "$ ") || strings.Contains(text, "\n$ ") { // a shell session
		lang = ""
	}
	return fence(lang, text)
}

// fence writes text as a fenced code block, fenced with more backticks than it holds in a row.
func fence(lang, text string) string {
	f := "```"
	for strings.Contains(text, f) {
		f += "`"
	}
	return f + lang + "\n" + text + "\n" + f
}

// list writes a <ul> or <ol> as a Markdown list. Its items are loose, with blank lines between
// them, when one of them holds more than one block.
func (c *converter) list(n *node) []string {
	start := 1
	if s, err := strconv.Atoi(n.attr["start"]); err == nil {
		start = s
	}
	var items [][]string
	for _, ch := range n.children {
		switch {
		case ch.tag == "li":
			items = append(items, c.blocks(ch.children))
		case ch.tag == "" && strings.TrimSpace(ch.text) == "":
		default: // content outside any item joins the one before it
			if len(items) == 0 {
				items = append(items, nil)
			}
			items[len(items)-1] = append(items[len(items)-1], c.blocks([]*node{ch})...)
		}
	}
	if len(items) == 0 {
		return nil
	}
	markers := make([]string, len(items))
	for i := range items {
		markers[i] = "- "
		if n.tag == "ol" {
			markers[i] = strconv.Itoa(start+i) + ". "
		}
	}
	return []string{listItems(markers, items)}
}

// definitions writes a <dl> as a list whose items start with their term in bold.
func (c *converter) definitions(n *node) []string {
	var items [][]string
	for _, ch := range n.children {
		switch ch.tag {
		case "dt":
			items = append(items, []string{emphasize("**", c.paragraph(ch.children))})
		case "dd":
			if len(items) == 0 {
				items = append(items, nil)
			}
			items[len(items)-1] = append(items[len(items)-1], c.blocks(ch.children)...)
		}
	}
	if len(items) == 0 {
		return nil
	}
	markers := make([]string, len(items))
	for i := range markers {
		markers[i] = "- "
	}
	return []string{listItems(markers, items)}
}

// listItems writes list items, each a marker and its blocks, the lines after the first indented
// to the marker's width.
func listItems(markers []string, items [][]string) string {
	loose := false
	for _, it := range items {
		loose = loose || len(it) > 1
	}
	out := make([]string, len(items))
	for i, it := range items {
		indent := strings.Repeat(" ", len(markers[i]))
		lines := strings.Split(strings.Join(it, "\n\n"), "\n")
		for j, line := range lines {
			switch {
			case j == 0:
				lines[j] = strings.TrimRight(markers[i]+line, " ")
			case line != "":
				lines[j] = indent + line
			}
		}
		out[i] = strings.Join(lines, "\n")
	}
	if loose {
		return strings.Join(out, "\n\n")
	}
	return strings.Join(out, "\n")
}

// quote writes Markdown as a block quote.
func quote(md string) string {
	lines := strings.Split(md, "\n")
	for i, l := range lines {
		if l == "" {
			lines[i] = ">"
		} else {
			lines[i] = "> " + l
		}
	}
	return strings.Join(lines, "\n")
}

// table writes a <table> as a GitHub table, which the app draws; its first row is the header.
// A cell's blocks are read as inline content.
func (c *converter) table(n *node) []string {
	var rows [][]string
	var walk func(*node)
	walk = func(x *node) {
		for _, ch := range x.children {
			switch ch.tag {
			case "tr":
				var row []string
				for _, cell := range ch.children {
					if cell.tag == "td" || cell.tag == "th" {
						s := strings.TrimSpace(reSpaces.ReplaceAllString(c.inline(cell.children), " "))
						row = append(row, strings.ReplaceAll(s, lineBreak, "<br>"))
					}
				}
				rows = append(rows, row)
			case "thead", "tbody", "tfoot":
				walk(ch)
			}
		}
	}
	walk(n)
	width := 0
	for _, r := range rows {
		width = max(width, len(r))
	}
	if width == 0 {
		return nil
	}
	var b strings.Builder
	for i, r := range rows {
		for len(r) < width {
			r = append(r, "")
		}
		b.WriteString("| " + strings.Join(r, " | ") + " |\n")
		if i == 0 {
			b.WriteString(strings.Repeat("| --- ", width) + "|\n")
		}
	}
	return []string{strings.TrimSuffix(b.String(), "\n")}
}

// plainText is what search reads of nodes: their text without markup, a line per paragraph
// and per line of code.
func plainText(nodes []*node) string {
	var b strings.Builder
	var walk func(n *node, pre bool)
	walk = func(n *node, pre bool) {
		switch {
		case n.tag == "" && pre:
			b.WriteString(n.text)
		case n.tag == "":
			b.WriteString(reSpaces.ReplaceAllString(n.text, " "))
		case droppedTags[n.tag]:
		case n.tag == "br":
			b.WriteString("\n")
		case n.tag == "sub":
			b.WriteString(subscript(plain(n)))
		case n.tag == "sup":
			b.WriteString("^" + strings.TrimSpace(plain(n)))
		default:
			if n.isBlock() {
				b.WriteString("\n")
			}
			for _, ch := range n.children {
				walk(ch, pre || n.tag == "pre")
			}
			if n.isBlock() {
				b.WriteString("\n")
			}
		}
	}
	for _, n := range nodes {
		walk(n, false)
	}
	var lines []string
	for line := range strings.SplitSeq(b.String(), "\n") {
		if line = strings.TrimSpace(reSpaces.ReplaceAllString(line, " ")); line != "" {
			lines = append(lines, line)
		}
	}
	return strings.Join(lines, "\n")
}
