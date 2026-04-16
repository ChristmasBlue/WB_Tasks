package pkg

import (
	"os"
	"strings"
)

func ParseEnvVar(str string) string {
	var result strings.Builder

	for i := 0; i < len(str); i++ {
		if str[i] == '$' && i+1 < len(str) {
			start := i + 1

			if '0' <= str[start] && str[start] <= '9' {
				result.WriteByte(str[i])
				continue
			}

			if str[start] == '{' {
				end := start + 1
				for ; end < len(str); end++ {
					if str[end] == '}' {
						break
					}
				}

				envVar := str[start+1 : end]

				if end < len(str) && str[end] == '}' {
					result.WriteString(os.Getenv(envVar))
				} else {
					result.WriteByte('$')
					result.WriteByte('{')
					result.WriteString(envVar)
				}

				i = end
				continue
			}

			if isAlphaNum(str[start]) || str[start] == '_' {
				end := start + 1
				for ; end < len(str); end++ {
					if !(isAlphaNum(str[end]) || str[end] == '_') {
						break
					}
				}

				envVar := str[start:end]
				result.WriteString(os.Getenv(envVar))
				i = end - 1
				continue
			}
		}
		result.WriteByte(str[i])
	}

	return result.String()
}

func isAlphaNum(c byte) bool {
	return ('a' <= c && 'z' >= c) ||
		('A' <= c && 'Z' >= c) ||
		('0' <= c && '9' >= c)
}
