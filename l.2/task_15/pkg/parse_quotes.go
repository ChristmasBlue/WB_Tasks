package pkg

import "strings"

func ParseQuotes(s string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	escaped := false

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if escaped {
			if ch == '"' || ch == '\'' || ch == '\\' {
				current.WriteByte(ch)
				escaped = false
				continue
			} else {
				current.WriteByte('\\')
				current.WriteByte(ch)
				escaped = false
				continue

			}
		}

		switch ch {
		case '"':
			if inQuotes {
				result = append(result, current.String())
				current.Reset()
				inQuotes = false
			} else {
				inQuotes = true
			}

		case ' ':

			if !inQuotes {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			} else {
				current.WriteByte(ch)
			}

		case '\\':
			escaped = true
		default:
			current.WriteByte(ch)
		}

	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
