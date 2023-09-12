// color.go
package pkg

import (
	"github.com/ddddddO/gtree"
	"github.com/fatih/color"
)

func AddNodeWithColor(parent *gtree.Node, keys []string, depth int, noColor bool) {
	if len(keys) == 0 {
		return
	}

	var coloredKey string
	if noColor {
		coloredKey = keys[0]
	} else if len(keys) == 1 {
		// This is a file
		coloredKey = colorFile(keys[0])
	} else {
		// This is a directory
		coloredKey = colorDirectory(keys[0])
	}

	node := parent.Add(coloredKey)
	AddNodeWithColor(node, keys[1:], depth+1, noColor)
}

func colorFile(key string) string {
	return color.New(color.FgGreen).SprintFunc()(key)
}

func colorDirectory(key string) string {
	return color.New(color.FgBlue).SprintFunc()(key)
}
