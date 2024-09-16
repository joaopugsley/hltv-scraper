package utils

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	c "github.com/chromedp/chromedp"
)

func hasClass(node *cdp.Node, class string) bool {
	parts := strings.Split(node.AttributeValue("class"), " ")
	for _, part := range parts {
		if part == class {
			return true
		}
	}
	return false
}

func GetFirstChildByClass(node *cdp.Node, class string) *cdp.Node {
	for _, child := range node.Children {
		if hasClass(child, class) {
			return child
		}
	}
	return nil
}

func GetChildrenByClass(node *cdp.Node, class string) []*cdp.Node {
	children := []*cdp.Node{}
	for _, child := range node.Children {
		if hasClass(child, class) {
			children = append(children, child)
		}
	}
	return children
}

// recursively search for the first child node with the given class ><
func DeepGetFirstChildByClass(node *cdp.Node, class string) *cdp.Node {
	if node == nil {
		return nil
	}

	if hasClass(node, class) {
		return node
	}

	for _, child := range node.Children {
		if result := DeepGetFirstChildByClass(child, class); result != nil {
			return result
		}
	}

	return nil
}

func GetElementText(ctx context.Context, node *cdp.Node) string {
	var text string
	err := c.Run(ctx,
		c.Text(node.FullXPath(), &text, c.NodeVisible),
	)

	if err != nil {
		log.Fatal(err)
	}

	return text
}
