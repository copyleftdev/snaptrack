package diff

import (
	"fmt"
	"strings"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

// ANSI escape sequences for basic colors
const (
	ansiReset = "\u001b[0m"
	ansiRed   = "\u001b[31m"
	ansiGreen = "\u001b[32m"
)

// GenerateDiff returns a unified, line-by-line diff between oldText and newText.
// If withColor is true, it includes ANSI color codes for +/– lines.
func GenerateDiff(oldText, newText string, withColor bool) string {
	// Split the input into lines
	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	// Create a new diff-match-patch instance
	differ := dmp.New()

	// Instead of the non-existent DiffMainLines, use DiffMain or a lines-to-chars approach:
	diffs := differ.DiffMain(strings.Join(oldLines, "\n"), strings.Join(newLines, "\n"), false)
	// diffs is a slice of Diff objects representing additions, deletions, or equals.

	// Convert to a unified diff
	unified := differ.PatchMake(strings.Join(oldLines, "\n"), strings.Join(newLines, "\n"), diffs)
	unifiedDiff := differ.PatchToText(unified)

	// If no color is requested, return the raw unified diff
	if !withColor {
		return unifiedDiff
	}

	// Otherwise, colorize lines starting with + or –
	var sb strings.Builder
	lines := strings.Split(unifiedDiff, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			sb.WriteString(fmt.Sprintf("%s%s%s\n", ansiGreen, line, ansiReset))
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			sb.WriteString(fmt.Sprintf("%s%s%s\n", ansiRed, line, ansiReset))
		} else {
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}
