package domain

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (cmd *Command) Redirect(redirect string, file string) error {

	if len(redirect) == 1 {
		switch redirect {
		case ">":
			cmd.RedirectOut(file)

		case "<":
			cmd.RedirectIn(file)

		default:
			return fmt.Errorf("invalid redirect: %s", redirect)
		}

		return nil
	}

	if len(redirect) == 2 && redirect[0] != '>' {

		if redirect[0] == '<' {
			return fmt.Errorf("invalid redirect: %s", redirect)
		}

		val, err := strconv.Atoi(redirect[:1])
		if err != nil {
			return fmt.Errorf("invalid redirect: %s", redirect)
		}

		switch val {
		case 1:
			cmd.RedirectOut(file)

		case 2:
			cmd.RedirectErr(file)

		default:
			return fmt.Errorf("invalid redirect: %s", redirect)
		}

		return nil
	}

	if strings.HasSuffix(redirect, ">>") {
		val := -1

		if len(redirect) > 2 {
			v, err := strconv.Atoi(redirect[:len(redirect)-2])
			if err != nil {
				return err
			}
			if 0 <= v && v <= 2 {
				val = v
			} else {
				return fmt.Errorf("invalid redirect: %s", redirect)
			}

		}

		switch val {
		case 1:
			cmd.RedirectAddToEnd(file)

		case 2:
			cmd.RedirectAddToEndErr(file)

		case -1:
			cmd.RedirectAddToEnd(file)

		default:
			return fmt.Errorf("invalid redirect: %s", redirect)
		}

		return nil
	}

	return fmt.Errorf("invalid redirect: %s", redirect)
}

func (cmd *Command) RedirectIn(file string) {
	cmd.StdIn = file
	cmd.OpenFlagIn = os.O_RDONLY
}

func (cmd *Command) RedirectOut(file string) {
	cmd.StdOut = file
	cmd.OpenFlagOut = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
}

func (cmd *Command) RedirectErr(file string) {
	cmd.StdErr = file
	cmd.OpenFlagErr = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
}

func (cmd *Command) RedirectAddToEndErr(file string) {
	cmd.StdErr = file
	cmd.OpenFlagErr = os.O_WRONLY | os.O_CREATE | os.O_APPEND
}

func (cmd *Command) RedirectAddToEnd(file string) {
	cmd.StdOut = file
	cmd.OpenFlagOut = os.O_WRONLY | os.O_CREATE | os.O_APPEND
}
