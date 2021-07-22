package uuidvalidator

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var uuidRegexp = regexp.MustCompile(`^[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}$`)

// ValidUUID checks that the String held in the attribute is a valid UUID
func ValidUUID() validator.String {
	return stringvalidator.RegexMatches(uuidRegexp, "value must be a valid UUID")
}
