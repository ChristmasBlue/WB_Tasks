package pkg

import "strings"

func ParseConditionalOperators(conditionalOperators string) []string {
	result := make([]string, 0)
	if conditionalOperators == "" {
		return result
	}

	var resultBuilder strings.Builder
	inQuotes := false
	escaped := false

	for i := 0; i < len(conditionalOperators); i++ {

		if conditionalOperators[i] == '\\' {
			escaped = true
			resultBuilder.WriteByte(conditionalOperators[i])
			continue
		}

		if escaped {
			escaped = false
			resultBuilder.WriteByte(conditionalOperators[i])
			continue
		}

		if conditionalOperators[i] == '"' {
			inQuotes = !inQuotes

			resultBuilder.WriteByte(conditionalOperators[i])
			continue
		}

		if conditionalOperators[i] == '&' && i+1 < len(conditionalOperators) {
			if conditionalOperators[i+1] == '&' && !inQuotes {
				result = append(result, strings.TrimSpace(resultBuilder.String()))
				resultBuilder.Reset()
				result = append(result, "&&")
				i++
				continue
			}
		}

		if conditionalOperators[i] == '|' && i+1 < len(conditionalOperators) {
			if conditionalOperators[i+1] == '|' && !inQuotes {
				result = append(result, strings.TrimSpace(resultBuilder.String()))
				resultBuilder.Reset()
				result = append(result, "||")
				i++
				continue
			}
		}

		resultBuilder.WriteByte(conditionalOperators[i])
	}

	if resultBuilder.Len() > 0 {
		result = append(result, strings.TrimSpace(resultBuilder.String()))
	}

	return result
}
