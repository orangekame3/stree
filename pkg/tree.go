package pkg

import (
	"github.com/ddddddO/gtree"
	"github.com/fatih/color"
)

func BuildTreeWithColor(root *gtree.Node, keys [][]string) *gtree.Node {
	for _, key := range keys {
		addNodeWithColor(root, key, 0)
	}
	return root
}

func BuildTreeWithoutColor(root *gtree.Node, keys [][]string) *gtree.Node {
	for _, key := range keys {
		addNodeWithoutColor(root, key, 0)
	}
	return root
}

func addNodeWithoutColor(parent *gtree.Node, keys []string, depth int) *gtree.Node {
	if len(keys) == 0 {
		return nil
	}
	// Add the current key (without color) as a node to the parent
	node := parent.Add(keys[0])
	// Recursively add the remaining keys
	return addNodeWithoutColor(node, keys[1:], depth+1)
}

func addNodeWithColor(parent *gtree.Node, keys []string, depth int) *gtree.Node {
	if len(keys) == 0 {
		return nil
	}
	// Determine the correct color for the node based on its type (file or directory)
	coloredKey := colorizeKey(keys[0], len(keys) == 1)
	// Add node to the tree
	node := parent.Add(coloredKey)
	return addNodeWithColor(node, keys[1:], depth+1)
}

func colorizeKey(key string, isFile bool) string {
	if isFile {
		return colorFile(key)
	}
	return colorDirectory(key)
}

func colorFile(key string) string {
	return color.New(color.FgGreen).SprintFunc()(key)
}

func colorDirectory(key string) string {
	return color.New(color.FgBlue).SprintFunc()(key)
}
