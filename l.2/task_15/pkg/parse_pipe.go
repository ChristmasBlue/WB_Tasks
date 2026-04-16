package pkg

import "strings"

func ParsePipes(pipeOperators string) []string {
	result := make([]string, 0)
	if pipeOperators == "" {
		return result
	}

	var resultBuilder strings.Builder
	inQuotes := false
	escaped := false

	for i := 0; i < len(pipeOperators); i++ {

		if pipeOperators[i] == '\\' {
			escaped = true
			resultBuilder.WriteByte(pipeOperators[i])
			continue
		}

		if escaped {
			escaped = false
			resultBuilder.WriteByte(pipeOperators[i])
			continue
		}

		if pipeOperators[i] == '"' {
			inQuotes = !inQuotes

			resultBuilder.WriteByte(pipeOperators[i])
			continue
		}

		if pipeOperators[i] == '|' && i+1 < len(pipeOperators) {
			if !inQuotes {
				if pipeOperators[i+1] == '|' {
					resultBuilder.WriteByte('|')
					resultBuilder.WriteByte('|')
					i++
					continue
				} else {
					result = append(result, strings.TrimSpace(resultBuilder.String()))
					resultBuilder.Reset()
					result = append(result, "|")
					continue
				}
			}
		}

		resultBuilder.WriteByte(pipeOperators[i])
	}

	if resultBuilder.Len() > 0 {
		result = append(result, strings.TrimSpace(resultBuilder.String()))
	}

	return result
}
