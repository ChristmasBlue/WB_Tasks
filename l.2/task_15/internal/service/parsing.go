package service

import "task_15/pkg"

func ParseCommand(commandLine string) ([]string, error) {
	expanded := pkg.ParseEnvVar(commandLine) //заменили переменные окружения
	pipes := pkg.ParsePipes(expanded)        //разбили на пайпы
	result := make([]string, 0)
	for _, pipe := range pipes {
		if pipe == "|" {
			result = append(result, pipe)
			continue
		}
		cond := pkg.ParseConditionalOperators(pipe) //разбиваем на условные операторы
		for _, c := range cond {
			if c == "&&" || c == "||" {
				result = append(result, c)
				continue
			}
			redirects, err := pkg.ParseRedirects(c) //разбиваем на редиректы
			if err != nil {
				return nil, err
			}
			result = append(result, redirects...) //сразу сохраняем всё в результат
		}
	}

	return result, nil
}
