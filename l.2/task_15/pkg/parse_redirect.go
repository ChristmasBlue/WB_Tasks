package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRedirects парсит редиректы, если перед редиректом стоит цифра(без пробельных символов), то учитывает её
func ParseRedirects(redirectOperators string) ([]string, error) {
	result := make([]string, 0)
	if redirectOperators == "" {
		return result, fmt.Errorf("line is empty")
	}

	lastIndex := 0
	file := ""

	inQuotes := false
	escaped := false

	for i := 0; i < len(redirectOperators); i++ {

		if redirectOperators[i] == '\\' {
			escaped = true
			continue
		}

		if escaped {
			escaped = false
			continue
		}

		if redirectOperators[i] == '"' {
			inQuotes = !inQuotes
			continue
		}

		if isRedirect(string(redirectOperators[i])) && i+1 < len(redirectOperators) {
			if !inQuotes {
				if redirectOperators[i+1] == '>' || redirectOperators[i+1] == '<' {
					if isRedirect(redirectOperators[i : i+2]) {
						index := lastSymbolIndex(redirectOperators, i)
						if index == -1 {
							addRes := strings.TrimSpace(redirectOperators[lastIndex:i])
							if addRes == "" {
								result = append(result, redirectOperators[i:i+2])
								lastIndex = i + 2
								i++
								file, lastIndex = searchFile(redirectOperators, lastIndex)
								if file != "" {
									result = append(result, file)
									file = ""
								}
								continue
							}
							result = append(result, strings.TrimSpace(redirectOperators[lastIndex:i]))
							result = append(result, redirectOperators[i:i+2])
							lastIndex = i + 2
							i++
							file, lastIndex = searchFile(redirectOperators, lastIndex)
							if file != "" {
								result = append(result, file)
								file = ""
							}
							continue
						}
						if redirectOperators[index:i] != "" {
							//если не получится конвертировать в число, возвращаем ошибку
							num, err := strconv.Atoi(redirectOperators[index:i])
							if err != nil {
								return nil, fmt.Errorf("parse redirect operators error:%v", err)
							}
							//если число меньше 0, тоже верну ошибку, редирект не работает с отрицательными значениями
							if num < 0 {
								return nil, fmt.Errorf("parse redirect operators less than 0")
							}
						}

						addRes := strings.TrimSpace(redirectOperators[lastIndex:index])
						if addRes == "" {
							result = append(result, strings.TrimSpace(redirectOperators[index:i+2]))
							lastIndex = i + 2
							i++
							file, lastIndex = searchFile(redirectOperators, lastIndex)
							if file != "" {
								result = append(result, file)
								file = ""
							}
							continue
						}

						result = append(result, strings.TrimSpace(redirectOperators[lastIndex:index]))
						result = append(result, strings.TrimSpace(redirectOperators[index:i+2]))
						lastIndex = i + 2
						i++
						file, lastIndex = searchFile(redirectOperators, lastIndex)
						if file != "" {
							result = append(result, file)
							file = ""
						}
						continue
					}
					return nil, fmt.Errorf("invalid redirect operator")
				}

				index := lastSymbolIndex(redirectOperators, i)
				if index == -1 {
					addRes := strings.TrimSpace(redirectOperators[lastIndex:i])
					if addRes == "" {
						result = append(result, string(redirectOperators[i]))
						lastIndex = i + 1
						file, lastIndex = searchFile(redirectOperators, lastIndex)
						if file != "" {
							result = append(result, file)
							file = ""
						}
						continue
					}
					result = append(result, strings.TrimSpace(redirectOperators[lastIndex:i]))
					result = append(result, string(redirectOperators[i]))
					lastIndex = i + 1
					file, lastIndex = searchFile(redirectOperators, lastIndex)
					if file != "" {
						result = append(result, file)
						file = ""
					}
					continue
				}

				if redirectOperators[index:i] != "" {
					//если не получится конвертировать в число, возвращаем ошибку
					num, err := strconv.Atoi(string(redirectOperators[index:i]))
					if err != nil {
						return nil, fmt.Errorf("parse redirect operators error:%v", err)
					}
					//если число меньше 0, тоже верну ошибку, редирект не работает с отрицательными значениями
					if num < 0 {
						return nil, fmt.Errorf("parse redirect operators less than 0")
					}
				}

				addRes := strings.TrimSpace(redirectOperators[lastIndex:index])
				if addRes == "" {
					result = append(result, strings.TrimSpace(redirectOperators[index:i+1]))
					lastIndex = i + 1
					file, lastIndex = searchFile(redirectOperators, lastIndex)
					if file != "" {
						result = append(result, file)
						file = ""
					}
					continue
				}

				result = append(result, strings.TrimSpace(redirectOperators[lastIndex:index]))
				result = append(result, strings.TrimSpace(redirectOperators[index:i+1]))
				lastIndex = i + 1
				file, lastIndex = searchFile(redirectOperators, lastIndex)
				if file != "" {
					result = append(result, file)
					file = ""
				}
				continue
			}
		}
	}

	if lastIndex < len(redirectOperators) {
		result = append(result, strings.TrimSpace(redirectOperators[lastIndex:]))
	}

	return result, nil
}

func isRedirect(s string) bool {
	switch s {
	case ">":
		return true
	case "<":
		return true

	case ">>":
		return true

	default:
		return false
	}
}

func lastSymbolIndex(s string, indexRedirect int) int {
	isNumber := false
	for i := indexRedirect - 1; i >= 0; i-- {
		if s[i] == ' ' || s[i] == '\t' {
			return i + 1
		}
		if s[i] >= '0' && s[i] <= '9' {
			isNumber = true
		}
	}
	if isNumber {
		return 0
	}
	return -1
}

func searchFile(s string, index int) (string, int) {
	inQuotes := false
	result := ""

	for index < len(s) && s[index] == ' ' {
		index++
	}
	end := index

	for i := index; i < len(s); i++ {
		if s[i] == '\\' {
			i++
			end = i
			continue
		}

		if s[i] == '"' {
			inQuotes = !inQuotes
			end = i
			continue
		}

		if !inQuotes && s[i] == ' ' {
			result = s[index:i]
			index = i
			break
		}
		end = i
	}

	if end > index {
		result = s[index : end+1]
		return result, end + 1
	}

	return result, index
}
