package diff

import (
	"fmt"
	"strings"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

const (
	ansiReset = "\u001b[0m"
	ansiRed   = "\u001b[31m"
	ansiGreen = "\u001b[32m"
)

func GenerateDiff(oldText, newText string, withColor bool) string {

	oldLines := strings.Split(oldText, "\n")
	newLines := strings.Split(newText, "\n")

	differ := dmp.New()

	diffs := differ.DiffMain(strings.Join(oldLines, "\n"), strings.Join(newLines, "\n"), false)

	unified := differ.PatchMake(strings.Join(oldLines, "\n"), strings.Join(newLines, "\n"), diffs)
	unifiedDiff := differ.PatchToText(unified)

	if !withColor {
		return unifiedDiff
	}

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
