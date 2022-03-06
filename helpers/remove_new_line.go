package helpers

import (
	"strings"
)

func RemoveNewLine(read_line string) string {
	return strings.TrimSuffix(read_line, "\n")
}
