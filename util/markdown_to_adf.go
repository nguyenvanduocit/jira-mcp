package util

import (
	"bytes"

	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// MarkdownToADF converts a markdown string to an Atlassian Document Format (ADF) CommentNodeScheme.
// Plain text without markdown syntax is wrapped in a paragraph node.
func MarkdownToADF(input string) *models.CommentNodeScheme {
	doc := &models.CommentNodeScheme{
		Version: 1,
		Type:    "doc",
	}

	if input == "" {
		return doc
	}

	source := []byte(input)

	md := goldmark.New(
		goldmark.WithExtensions(extension.Strikethrough),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	)

	reader := text.NewReader(source)
	tree := md.Parser().Parse(reader)

	walkChildren(tree, doc, source)

	return doc
}

func walkChildren(parent ast.Node, adfParent *models.CommentNodeScheme, source []byte) {
	for child := parent.FirstChild(); child != nil; child = child.NextSibling() {
		nodes := convertNode(child, source)
		for _, node := range nodes {
			adfParent.AppendNode(node)
		}
	}
}

func convertNode(n ast.Node, source []byte) []*models.CommentNodeScheme {
	switch node := n.(type) {
	case *ast.Paragraph:
		p := &models.CommentNodeScheme{Type: "paragraph"}
		walkInline(node, p, source, nil)
		return []*models.CommentNodeScheme{p}

	case *ast.Heading:
		h := &models.CommentNodeScheme{
			Type:  "heading",
			Attrs: map[string]interface{}{"level": float64(node.Level)},
		}
		walkInline(node, h, source, nil)
		return []*models.CommentNodeScheme{h}

	case *ast.FencedCodeBlock:
		lang := ""
		if node.Language(source) != nil {
			lang = string(node.Language(source))
		}
		cb := &models.CommentNodeScheme{
			Type: "codeBlock",
		}
		if lang != "" {
			cb.Attrs = map[string]interface{}{"language": lang}
		}
		// Extract code text from lines
		var buf bytes.Buffer
		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			buf.Write(line.Value(source))
		}
		text := buf.String()
		if text != "" {
			cb.AppendNode(&models.CommentNodeScheme{
				Type: "text",
				Text: text,
			})
		}
		return []*models.CommentNodeScheme{cb}

	case *ast.CodeBlock:
		cb := &models.CommentNodeScheme{Type: "codeBlock"}
		var buf bytes.Buffer
		for i := 0; i < node.Lines().Len(); i++ {
			line := node.Lines().At(i)
			buf.Write(line.Value(source))
		}
		text := buf.String()
		if text != "" {
			cb.AppendNode(&models.CommentNodeScheme{
				Type: "text",
				Text: text,
			})
		}
		return []*models.CommentNodeScheme{cb}

	case *ast.Blockquote:
		bq := &models.CommentNodeScheme{Type: "blockquote"}
		walkChildren(node, bq, source)
		return []*models.CommentNodeScheme{bq}

	case *ast.List:
		listType := "bulletList"
		if node.IsOrdered() {
			listType = "orderedList"
		}
		list := &models.CommentNodeScheme{Type: listType}
		if node.IsOrdered() && node.Start != 1 {
			list.Attrs = map[string]interface{}{"order": float64(node.Start)}
		}
		walkChildren(node, list, source)
		return []*models.CommentNodeScheme{list}

	case *ast.ListItem:
		li := &models.CommentNodeScheme{Type: "listItem"}
		walkChildren(node, li, source)
		return []*models.CommentNodeScheme{li}

	case *ast.ThematicBreak:
		return []*models.CommentNodeScheme{{Type: "rule"}}

	case *ast.TextBlock:
		p := &models.CommentNodeScheme{Type: "paragraph"}
		walkInline(node, p, source, nil)
		return []*models.CommentNodeScheme{p}

	default:
		// For unknown block nodes, try to process children
		var result []*models.CommentNodeScheme
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			result = append(result, convertNode(child, source)...)
		}
		return result
	}
}

func walkInline(parent ast.Node, adfParent *models.CommentNodeScheme, source []byte, marks []*models.MarkScheme) {
	for child := parent.FirstChild(); child != nil; child = child.NextSibling() {
		convertInline(child, adfParent, source, marks)
	}
}

func convertInline(n ast.Node, adfParent *models.CommentNodeScheme, source []byte, marks []*models.MarkScheme) {
	switch node := n.(type) {
	case *ast.Text:
		t := string(node.Segment.Value(source))
		if t != "" {
			textNode := &models.CommentNodeScheme{
				Type: "text",
				Text: t,
			}
			if len(marks) > 0 {
				textNode.Marks = copyMarks(marks)
			}
			adfParent.AppendNode(textNode)
		}
		if node.HardLineBreak() {
			adfParent.AppendNode(&models.CommentNodeScheme{Type: "hardBreak"})
		} else if node.SoftLineBreak() {
			// Treat soft line break as a space or hardBreak depending on context
			// In ADF, soft breaks within a paragraph are typically just spaces
		}

	case *ast.String:
		t := string(node.Value)
		if t != "" {
			textNode := &models.CommentNodeScheme{
				Type: "text",
				Text: t,
			}
			if len(marks) > 0 {
				textNode.Marks = copyMarks(marks)
			}
			adfParent.AppendNode(textNode)
		}

	case *ast.Emphasis:
		markType := "em"
		if node.Level == 2 {
			markType = "strong"
		}
		newMarks := append(copyMarks(marks), &models.MarkScheme{Type: markType})
		walkInline(node, adfParent, source, newMarks)

	case *ast.CodeSpan:
		// Collect all text content from the code span
		var buf bytes.Buffer
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			if textNode, ok := child.(*ast.Text); ok {
				buf.Write(textNode.Segment.Value(source))
			}
		}
		t := buf.String()
		if t != "" {
			newMarks := append(copyMarks(marks), &models.MarkScheme{Type: "code"})
			adfParent.AppendNode(&models.CommentNodeScheme{
				Type:  "text",
				Text:  t,
				Marks: newMarks,
			})
		}

	case *ast.Link:
		linkMark := &models.MarkScheme{
			Type:  "link",
			Attrs: map[string]interface{}{"href": string(node.Destination)},
		}
		newMarks := append(copyMarks(marks), linkMark)
		walkInline(node, adfParent, source, newMarks)

	case *ast.AutoLink:
		url := string(node.URL(source))
		linkMark := &models.MarkScheme{
			Type:  "link",
			Attrs: map[string]interface{}{"href": url},
		}
		newMarks := append(copyMarks(marks), linkMark)
		adfParent.AppendNode(&models.CommentNodeScheme{
			Type:  "text",
			Text:  url,
			Marks: newMarks,
		})

	case *east.Strikethrough:
		newMarks := append(copyMarks(marks), &models.MarkScheme{Type: "strike"})
		walkInline(node, adfParent, source, newMarks)

	default:
		// For unknown inline nodes, try to walk children
		walkInline(n, adfParent, source, marks)
	}
}

func copyMarks(marks []*models.MarkScheme) []*models.MarkScheme {
	if marks == nil {
		return nil
	}
	result := make([]*models.MarkScheme, len(marks))
	copy(result, marks)
	return result
}
