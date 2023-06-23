package provider

import (
	"fmt"
	"strings"
)

func getFormattedValue(val string, description string) string {
	return fmt.Sprint("\n  | ", strings.ReplaceAll(val, "|", "\\|"), " | ", description, " | ")
}
