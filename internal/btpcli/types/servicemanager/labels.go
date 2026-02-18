package servicemanager

import (
	"regexp"
	"strings"
)

type ServiceManagerLabels map[string][]string

func (s *ServiceManagerLabels) UnmarshalJSON(data []byte) error {
	*s = make(ServiceManagerLabels)

	r := regexp.MustCompile(`([a-zA-Z0-9_\-]*)\s*=\s*((?:[^,;]+)(?:\s*,\s*(?:[^,;]+))*)`)

	matches := r.FindAllStringSubmatch(string(data), -1)

	for _, match := range matches {
		valuesMap := make(map[string]struct{})

		for val := range strings.SplitSeq(match[2], ",") {
			valuesMap[strings.TrimSpace(val)] = struct{}{}
		}

		values := []string{}
		for val := range valuesMap {
			values = append(values, val)
		}

		(*s)[match[1]] = values
	}

	return nil
}
