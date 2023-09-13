// Package: pkg is a package that contains the business logic for stree
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

	// Add a colored node to the tree based on whether the current key represents a file or a directory
	coloredKey := getColorizedKey(keys[0], len(keys) == 1)
	node := parent.Add(coloredKey)

	return addNodeWithColor(node, keys[1:], depth+1)
}

// getColorizedKey returns the colored representation of a key, with different colors for files and directories
func getColorizedKey(key string, isFile bool) string {
	if isFile {
		return color.New(color.FgGreen).SprintFunc()(key)
	}
	return color.New(color.FgBlue).SprintFunc()(key)
}
