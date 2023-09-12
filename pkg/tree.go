package pkg

import (
	"github.com/ddddddO/gtree"
)

func BuildTree(root *gtree.Node, keys [][]string, noColor bool) *gtree.Node {
	if noColor {
		for _, key := range keys {
			addNodeWithoutColor(root, key, 0)
		}
		return root
	}
	for _, key := range keys {
		addNodeWithColor(root, key, 0)
	}
	return root
}
