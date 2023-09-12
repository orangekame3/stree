package pkg

import (
	"github.com/ddddddO/gtree"
	"github.com/fatih/color"
)

func addNodeWithColor(parent *gtree.Node, keys []string, depth int) *gtree.Node {
	if len(keys) == 0 {
		return nil
	}
	var coloredKey string
	if len(keys) == 1 {
		// This is a file
		coloredKey = colorFile(keys[0])
	} else {
		// This is a directory
		coloredKey = colorDirectory(keys[0])
	}

	node := parent.Add(coloredKey)
	return addNodeWithColor(node, keys[1:], depth+1)
}

func colorFile(key string) string {
	return color.New(color.FgGreen).SprintFunc()(key)
}

func colorDirectory(key string) string {
	return color.New(color.FgBlue).SprintFunc()(key)
}

func addNodeWithoutColor(parent *gtree.Node, keys []string, depth int) *gtree.Node {
	if len(keys) == 0 {
		return nil
	}

	coloredKey := keys[0]
	node := parent.Add(coloredKey)
	return addNodeWithoutColor(node, keys[1:], depth+1)
}
