package main

import "strings"

func Pad(s string, indent int) string {
	lines := strings.Split(s, "\n")
	longestLine := -1
	for _, l := range lines {
		if len(l) > longestLine {
			longestLine = len(l)
		}
	}
	for i, l := range lines {
		lines[i] = strings.Repeat(" ", indent) + l + strings.Repeat(" ", longestLine-len(l))
	}
	return strings.Join(lines, "\n")
}
