package domain

import (
	"fmt"
	"strings"
)

func (c *Conditionals) ParseCommands(commands []string) {
	if len(commands) == 0 {
		return
	}

	for i := 0; i < len(commands); i++ {
		switch {
		case commands[i] == "&&":
			if i == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: && can`t stay first"))
			}

			if len(c.Subsequence) == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: subsequence is empty"))
			}

			c.Subsequence[len(c.Subsequence)-1].And = true
			continue

		case commands[i] == "||":
			if i == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: || can`t stay first"))
			}

			if len(c.Subsequence) == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: subsequence is empty"))
			}

			c.Subsequence[len(c.Subsequence)-1].Or = true
			continue

		case strings.HasSuffix(commands[i], "<") || strings.HasSuffix(commands[i], ">"):
			if i == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: redirect can`t stay first"))
			}

			if len(c.Subsequence) == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: subsequence is empty"))
			}

			if i >= len(commands)-1 || commands[i+1] == "&&" || commands[i+1] == "||" || commands[i+1] == "|" {
				c.Subsequence[len(c.Subsequence)-1].Redirect(commands[i], "")
			} else {
				c.Subsequence[len(c.Subsequence)-1].Redirect(commands[i], commands[i+1])
				i++
			}
			continue

		case commands[i] == "|":
			if i == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: | can`t stay first"))
			}

			if len(c.Subsequence) == 0 {
				c.Errs = append(c.Errs, fmt.Errorf("invalid command: subsequence is empty"))
			}

			c.Subsequence[len(c.Subsequence)-1].Pipe = true

		default:
			c.Subsequence = append(c.Subsequence, NewCommand(commands[i]))
		}
	}
}
