package provider

import (
	"fmt"
	"strings"
)

func getFormattedValueAsTableRow(val string, description string) string {
	return fmt.Sprint("\n  | ", strings.ReplaceAll(val, "|", "\\|"), " | ", description, " | ")
}
