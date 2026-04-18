package util

import (
	"encoding/json"
	"testing"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
)

func TestMarkdownToADF_EmptyString(t *testing.T) {
	result := MarkdownToADF("")
	if result.Type != "doc" || result.Version != 1 {
		t.Errorf("expected doc node, got %+v", result)
	}
	if len(result.Content) != 0 {
		t.Errorf("expected no content for empty string, got %d nodes", len(result.Content))
	}
}

func TestMarkdownToADF_PlainText(t *testing.T) {
	result := MarkdownToADF("Hello world")
	assertDocWithParagraph(t, result)
	p := result.Content[0]
	if len(p.Content) != 1 || p.Content[0].Text != "Hello world" {
		t.Errorf("expected text 'Hello world', got %+v", p.Content)
	}
}

func TestMarkdownToADF_Headings(t *testing.T) {
	tests := []struct {
		input string
		level float64
	}{
		{"# H1", 1},
		{"## H2", 2},
		{"### H3", 3},
		{"#### H4", 4},
		{"##### H5", 5},
		{"###### H6", 6},
	}
	for _, tt := range tests {
		result := MarkdownToADF(tt.input)
		if len(result.Content) == 0 {
			t.Fatalf("no content for %q", tt.input)
		}
		h := result.Content[0]
		if h.Type != "heading" {
			t.Errorf("expected heading, got %s for %q", h.Type, tt.input)
		}
		if h.Attrs["level"] != tt.level {
			t.Errorf("expected level %v, got %v for %q", tt.level, h.Attrs["level"], tt.input)
		}
	}
}

func TestMarkdownToADF_Bold(t *testing.T) {
	result := MarkdownToADF("**bold**")
	p := result.Content[0]
	assertMark(t, p.Content[0], "bold", "strong")
}

func TestMarkdownToADF_Italic(t *testing.T) {
	result := MarkdownToADF("*italic*")
	p := result.Content[0]
	assertMark(t, p.Content[0], "italic", "em")
}

func TestMarkdownToADF_InlineCode(t *testing.T) {
	result := MarkdownToADF("`code`")
	p := result.Content[0]
	assertMark(t, p.Content[0], "code", "code")
}

func TestMarkdownToADF_Strikethrough(t *testing.T) {
	result := MarkdownToADF("~~strike~~")
	p := result.Content[0]
	assertMark(t, p.Content[0], "strike", "strike")
}

func TestMarkdownToADF_Link(t *testing.T) {
	result := MarkdownToADF("[click](https://example.com)")
	p := result.Content[0]
	if len(p.Content) == 0 {
		t.Fatal("no content in paragraph")
	}
	textNode := p.Content[0]
	if textNode.Text != "click" {
		t.Errorf("expected text 'click', got %q", textNode.Text)
	}
	found := false
	for _, m := range textNode.Marks {
		if m.Type == "link" {
			found = true
			if m.Attrs["href"] != "https://example.com" {
				t.Errorf("expected href https://example.com, got %v", m.Attrs["href"])
			}
		}
	}
	if !found {
		t.Error("expected link mark not found")
	}
}

func TestMarkdownToADF_BulletList(t *testing.T) {
	result := MarkdownToADF("- item 1\n- item 2\n- item 3")
	if len(result.Content) == 0 {
		t.Fatal("no content")
	}
	list := result.Content[0]
	if list.Type != "bulletList" {
		t.Errorf("expected bulletList, got %s", list.Type)
	}
	if len(list.Content) != 3 {
		t.Errorf("expected 3 list items, got %d", len(list.Content))
	}
	for _, li := range list.Content {
		if li.Type != "listItem" {
			t.Errorf("expected listItem, got %s", li.Type)
		}
	}
}

func TestMarkdownToADF_OrderedList(t *testing.T) {
	result := MarkdownToADF("1. first\n2. second")
	if len(result.Content) == 0 {
		t.Fatal("no content")
	}
	list := result.Content[0]
	if list.Type != "orderedList" {
		t.Errorf("expected orderedList, got %s", list.Type)
	}
	if len(list.Content) != 2 {
		t.Errorf("expected 2 list items, got %d", len(list.Content))
	}
}

func TestMarkdownToADF_FencedCodeBlock(t *testing.T) {
	input := "```go\nfunc main() {}\n```"
	result := MarkdownToADF(input)
	if len(result.Content) == 0 {
		t.Fatal("no content")
	}
	cb := result.Content[0]
	if cb.Type != "codeBlock" {
		t.Errorf("expected codeBlock, got %s", cb.Type)
	}
	if cb.Attrs["language"] != "go" {
		t.Errorf("expected language 'go', got %v", cb.Attrs["language"])
	}
	if len(cb.Content) == 0 || cb.Content[0].Text != "func main() {}\n" {
		t.Errorf("unexpected code content: %+v", cb.Content)
	}
}

func TestMarkdownToADF_CodeBlockNoLang(t *testing.T) {
	input := "```\nhello\n```"
	result := MarkdownToADF(input)
	cb := result.Content[0]
	if cb.Type != "codeBlock" {
		t.Errorf("expected codeBlock, got %s", cb.Type)
	}
	if cb.Attrs != nil {
		t.Errorf("expected no attrs, got %v", cb.Attrs)
	}
}

func TestMarkdownToADF_Blockquote(t *testing.T) {
	result := MarkdownToADF("> quoted text")
	if len(result.Content) == 0 {
		t.Fatal("no content")
	}
	bq := result.Content[0]
	if bq.Type != "blockquote" {
		t.Errorf("expected blockquote, got %s", bq.Type)
	}
	// Blockquote should contain a paragraph
	if len(bq.Content) == 0 || bq.Content[0].Type != "paragraph" {
		t.Errorf("expected paragraph inside blockquote, got %+v", bq.Content)
	}
}

func TestMarkdownToADF_HorizontalRule(t *testing.T) {
	result := MarkdownToADF("above\n\n---\n\nbelow")
	found := false
	for _, node := range result.Content {
		if node.Type == "rule" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected rule node not found")
	}
}

func TestMarkdownToADF_MixedContent(t *testing.T) {
	input := "# Title\n\nSome **bold** text.\n\n- item 1\n- item 2\n\n```js\nconsole.log('hi')\n```"
	result := MarkdownToADF(input)
	if len(result.Content) < 4 {
		t.Fatalf("expected at least 4 nodes, got %d", len(result.Content))
	}
	if result.Content[0].Type != "heading" {
		t.Errorf("expected heading first, got %s", result.Content[0].Type)
	}
	if result.Content[1].Type != "paragraph" {
		t.Errorf("expected paragraph second, got %s", result.Content[1].Type)
	}
	if result.Content[2].Type != "bulletList" {
		t.Errorf("expected bulletList third, got %s", result.Content[2].Type)
	}
	if result.Content[3].Type != "codeBlock" {
		t.Errorf("expected codeBlock fourth, got %s", result.Content[3].Type)
	}
}

func TestMarkdownToADF_BoldItalic(t *testing.T) {
	result := MarkdownToADF("***bold italic***")
	p := result.Content[0]
	if len(p.Content) == 0 {
		t.Fatal("no content in paragraph")
	}
	textNode := p.Content[0]
	hasStrong := false
	hasEm := false
	for _, m := range textNode.Marks {
		if m.Type == "strong" {
			hasStrong = true
		}
		if m.Type == "em" {
			hasEm = true
		}
	}
	if !hasStrong || !hasEm {
		t.Errorf("expected both strong and em marks, got %+v", textNode.Marks)
	}
}

// Regression test for the Jira ADF constraint: the "code" mark is exclusive —
// combining it with strong/em/strike produces invalid ADF that Jira rejects.
// See util/markdown_to_adf.go:convertInline for the CodeSpan branch.
func TestMarkdownToADF_CodeMarkIsExclusiveInsideBold(t *testing.T) {
	// "**bold `code` bold**" — the code span sits inside strong emphasis.
	// Without the fix, the inner text would carry [strong, code].
	result := MarkdownToADF("**bold `code` bold**")
	if len(result.Content) == 0 {
		t.Fatal("no content")
	}
	p := result.Content[0]

	// Find the text node whose Text is "code".
	var codeNode *models.CommentNodeScheme
	for _, n := range p.Content {
		if n.Text == "code" {
			codeNode = n
			break
		}
	}
	if codeNode == nil {
		t.Fatalf("expected a text node with Text=%q, got content=%+v", "code", p.Content)
	}

	if len(codeNode.Marks) != 1 {
		t.Fatalf("code span must carry exactly one mark, got %d: %+v", len(codeNode.Marks), codeNode.Marks)
	}
	if codeNode.Marks[0].Type != "code" {
		t.Errorf("code span mark type = %q, want %q", codeNode.Marks[0].Type, "code")
	}
}

func TestMarkdownToADF_CodeMarkIsExclusiveInsideLink(t *testing.T) {
	// "[text with `code` inside](https://example.com)" — code span inside link.
	// Code mark must stand alone; surrounding text keeps the link mark.
	result := MarkdownToADF("[text with `code` inside](https://example.com)")
	if len(result.Content) == 0 {
		t.Fatal("no content")
	}
	p := result.Content[0]

	var codeNode *models.CommentNodeScheme
	for _, n := range p.Content {
		if n.Text == "code" {
			codeNode = n
			break
		}
	}
	if codeNode == nil {
		t.Fatalf("expected a text node with Text=%q, got content=%+v", "code", p.Content)
	}

	if len(codeNode.Marks) != 1 || codeNode.Marks[0].Type != "code" {
		t.Errorf("code span must carry only [code] mark, got %+v", codeNode.Marks)
	}
}

func TestMarkdownToADF_ValidJSON(t *testing.T) {
	input := "# Hello\n\nWorld **bold** and *italic*\n\n- list\n\n```go\ncode\n```"
	result := MarkdownToADF(input)
	_, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}
}

// Helpers

func assertDocWithParagraph(t *testing.T, doc *models.CommentNodeScheme) {
	t.Helper()
	if doc.Type != "doc" {
		t.Fatalf("expected doc, got %s", doc.Type)
	}
	if len(doc.Content) == 0 {
		t.Fatal("expected at least one child node")
	}
	if doc.Content[0].Type != "paragraph" {
		t.Fatalf("expected paragraph, got %s", doc.Content[0].Type)
	}
}

func assertMark(t *testing.T, node *models.CommentNodeScheme, expectedText string, markType string) {
	t.Helper()
	if node == nil {
		t.Fatal("node is nil")
	}
	if node.Text != expectedText {
		t.Errorf("expected text %q, got %q", expectedText, node.Text)
	}
	found := false
	for _, m := range node.Marks {
		if m.Type == markType {
			found = true
		}
	}
	if !found {
		t.Errorf("expected mark %q not found in %+v", markType, node.Marks)
	}
}
