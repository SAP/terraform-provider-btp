package provider

import (
	"fmt"
	"strings"
)

func getFormattedValueAsTableRow(val string, description string) string {
	return fmt.Sprintf("\n  | %s | %s | ", strings.ReplaceAll(val, "|", "\\|"), strings.ReplaceAll(description, "|", "\\|"))
}
