// Package pkg provides the core functionality of the program.
package pkg

import (
	"fmt"
	"strings"

	"github.com/ddddddO/gtree"
	"github.com/fatih/color"
)

// createFullPath creates a full path for each keys
func createFullPath(bucket string, keys [][]string) [][]string {
	for i := range keys {
		keys[i][len(keys[i])-1] = fmt.Sprintf("%s/%s", bucket, strings.Join(keys[i], "/"))
	}
	return keys
}

// BuildTreeWithColor builds a tree with colored nodes
func BuildTreeWithColor(root *gtree.Node, bucket string, keys [][]string, f bool) *gtree.Node {
	if f {
		keys = createFullPath(bucket, keys)
	}

	for _, key := range keys {
		addNodeWithColor(root, key, 0)
	}
	return root
}

// BuildTreeWithoutColor builds a tree without colored nodes
func BuildTreeWithoutColor(root *gtree.Node, bucket string, keys [][]string, f bool) *gtree.Node {
	if f {
		keys = createFullPath(bucket, keys)
	}
	for _, key := range keys {
		addNodeWithoutColor(root, key, 0)
	}
	return root
}

func addNodeWithoutColor(parent *gtree.Node, keys []string, depth int) *gtree.Node {
	if len(keys) == 0 {
		return nil
	}

	// Skip adding empty strings as nodes
	if len(keys) == 1 && keys[0] == "" {
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

	// Skip adding empty strings as nodes
	if len(keys) == 1 && keys[0] == "" {
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
