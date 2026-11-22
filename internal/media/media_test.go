package media

import "testing"

func TestFindPicturesAndHTML(t *testing.T) {
	for md, want := range map[string]Found{
		"![The dish at Ystmark](dish.png)":                          {Images: 1},
		"Before ![a](a.png) and ![b][ref] after.\n\n[ref]: b.png\n": {Images: 2},
		"![The launch](https://example.com/launch.mp4)":             {Images: 1},
		"[![badge](badge.svg)](https://example.com)":                {Images: 1},
		"<img src=\"dish.png\" alt=\"The dish\">":                   {HTML: 1},
		"Text with <video src=\"pass.mp4\"></video> inside.":        {HTML: 2},
		"<iframe src=\"https://example.com/embed\"></iframe>\n":     {HTML: 1},
		"<div>\nan HTML block\n</div>\n":                            {HTML: 1},
		"<!-- a note to the writer -->\n":                           {HTML: 1},
		"Run go doc <package> to read it.":                          {HTML: 1},
		"A picture ![x](x.png) and a <br> line break.":              {Images: 1, HTML: 1},
	} {
		if got := Find(md); got != want {
			t.Errorf("Find(%q) = %+v, want %+v", md, got, want)
		}
	}
}

// Code is text: a lesson may teach HTML or Markdown, and its examples stay.
func TestFindLeavesCode(t *testing.T) {
	for _, md := range []string{
		"```html\n<html>\n  <body>\n    <img src=\"/rocket.png\" alt=\"Apogee rocket\" />\n  </body>\n</html>\n```\n",
		"```md\n![a rocket](ROCKET_URL)\n```\n",
		"> ```html\n> <iframe src=\"/telemetry-panel.html\"></iframe>\n> ```\n",
		"Use the `<video>` element, and `![x](y.png)`, and ``<img src=\"a`b.png\">``.\n",
		"Indented code:\n\n    page := `<html><body><img src=\"/pad.png\"></body></html>`\n",
		"```go\nfmt.Println(\"<img src=x>\")\n```\n",
	} {
		if got := Find(md); got.Total() != 0 {
			t.Errorf("Find(%q) = %+v in code, want nothing", md, got)
		}
	}
}

// Text that only looks like HTML or a picture, and links, which the app shows as text or
// follows itself, are no media.
func TestFindLeavesText(t *testing.T) {
	for _, md := range []string{
		"",
		"Plain text with no media.",
		"If a < b and b > c, the pass ends.",
		"A channel send is ch <- v, a receive <-ch.",
		"See <https://pkg.go.dev/fmt>, an autolink.",
		"The replay is [on the site](https://example.com/replay.mp4).",
		"- [ ] task\n- [x] done\n\n> [!warning]\n> Careful.\n",
		"| a | b |\n|---|---|\n| 1 | 2 |\n",
		"Tests expect `launch-replay.mp4` in the listing.",
		"An escaped \\<img> is text, as is \\![a](b.png).",
	} {
		if got := Find(md); got.Total() != 0 {
			t.Errorf("Find(%q) = %+v, want nothing", md, got)
		}
	}
}

func TestOneParagraph(t *testing.T) {
	for md, want := range map[string]bool{
		"Look at what `clock.go` defines.":                        true,
		"A hint that wraps\nonto a second line, with *emphasis*.": true,
		"See [the fmt docs](https://pkg.go.dev/fmt) for `%d`.":    true,
		"":                            false,
		"Two paragraphs.\n\nNot one.": false,
		"# A heading":                 false,
		"- a list item":               false,
		"1. a numbered item":          false,
		"> a quote":                   false,
		"```go\nfmt.Println()\n```":   false,
		"    indented code":           false,
		"<div>an HTML block</div>":    false,
		"Text\n---":                   false, // a heading, by its underline
		"---":                         false,
	} {
		if got := OneParagraph(md); got != want {
			t.Errorf("OneParagraph(%q) = %v, want %v", md, got, want)
		}
	}
}
