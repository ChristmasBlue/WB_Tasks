package domain

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"task_15/pkg"
)

func (c *Conditionals) Execute() {
	isAnd := false
	isPrevOk := true
	for i := 0; i < len(c.Subsequence); i++ {
		if len(c.Errs) != 0 {
			return
		}

		switch {
		case c.Subsequence[i].And:
			if !isPrevOk {
				continue
			}
			isAnd = true

		case c.Subsequence[i].Or:
			isAnd = false
			isPrevOk = true

		case c.Subsequence[i].Pipe:
			if !isPrevOk {
				continue
			}
			index, errs := c.ExecutePipe(i)
			if index >= len(c.Subsequence)-1 {
				return
			}
			i = index
			isAnd = c.Subsequence[i].And
			if errs != nil {
				c.Errs = append(c.Errs, errs...)
				isPrevOk = false
				continue
			}

			continue
		}

		err := c.Subsequence[i].Execute()
		if err != nil {
			c.Errs = append(c.Errs, err)
			isPrevOk = false
			continue
		}

		if !isAnd && isPrevOk {
			return
		}
	}
}

func (c *Command) Execute() error {
	args := pkg.ParseQuotes(c.Text)
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "exit":
		os.Exit(0) // или возвращать специальную ошибку
		return nil

	case "echo":
		output := strings.Join(args[1:], " ")

		// Поддержка редиректа stdout
		if c.PipeWriter != nil {
			c.PipeWriter.Write([]byte(output + "\n"))
			c.PipeWriter.Close()
		} else if c.StdOut != "" {
			file, err := os.OpenFile(c.StdOut, c.OpenFlagOut, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			file.WriteString(output + "\n")
		} else {
			fmt.Println(output)
		}

		return nil

	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if c.StdOut != "" {
			file, err := os.OpenFile(c.StdOut, c.OpenFlagOut, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			file.WriteString(dir + "\n")
		} else {
			fmt.Println(dir)
		}
		return nil

	case "cd":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "cd: missing argument")
			return nil
		}
		if err := os.Chdir(args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "cd: %v\n", err)
			return nil // не возвращаем ошибку, просто выводим
		}
		return nil

	case "ps":
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("tasklist")
		} else {
			cmd = exec.Command("ps", "aux")
		}

		if c.PipeReader != nil {
			cmd.Stdin = c.PipeReader
		} else if c.StdIn != "" {
			file, err := os.OpenFile(c.StdIn, c.OpenFlagIn, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stdin = file
		} else {
			cmd.Stdin = os.Stdin
		}

		if c.PipeWriter != nil {
			cmd.Stdout = c.PipeWriter
		} else if c.StdOut != "" {
			file, err := os.OpenFile(c.StdOut, c.OpenFlagOut, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stdout = file
		} else {
			cmd.Stdout = os.Stdout
		}

		if c.StdErr != "" {
			file, err := os.OpenFile(c.StdErr, c.OpenFlagErr, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stderr = file
		} else {
			cmd.Stderr = os.Stderr
		}

		cmd.Stderr = os.Stderr
		return cmd.Run()

	case "kill":

		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "kill: missing PID")
			return nil
		}
		pid, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "kill: invalid PID: %v\n", err)
			return nil
		}
		if runtime.GOOS == "windows" {
			cmd := exec.Command("taskkill", "/PID", args[1])
			if c.PipeReader != nil {
				cmd.Stdin = c.PipeReader
			} else if c.StdIn != "" {
				file, err := os.OpenFile(c.StdIn, c.OpenFlagIn, 0644)
				if err != nil {
					return err
				}
				defer file.Close()
				cmd.Stdin = file
			} else {
				cmd.Stdin = os.Stdin
			}

			if c.PipeWriter != nil {
				cmd.Stdout = c.PipeWriter
			} else if c.StdOut != "" {
				file, err := os.OpenFile(c.StdOut, c.OpenFlagOut, 0644)
				if err != nil {
					return err
				}
				defer file.Close()
				cmd.Stdout = file
			} else {
				cmd.Stdout = os.Stdout
			}

			if c.StdErr != "" {
				file, err := os.OpenFile(c.StdErr, c.OpenFlagErr, 0644)
				if err != nil {
					return err
				}
				defer file.Close()
				cmd.Stderr = file
			} else {
				cmd.Stderr = os.Stderr
			}
			return cmd.Run()
		} else {
			process, err := os.FindProcess(pid)
			if err != nil {
				fmt.Fprintf(os.Stderr, "kill: %v\n", err)
				return nil
			}
			if err := process.Signal(os.Interrupt); err != nil {
				fmt.Fprintf(os.Stderr, "kill: %v\n", err)
			}
			return nil
		}

	default:
		cmd := exec.Command(args[0], args[1:]...)
		if c.PipeReader != nil {
			cmd.Stdin = c.PipeReader
		} else if c.StdIn != "" {
			file, err := os.OpenFile(c.StdIn, c.OpenFlagIn, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stdin = file
		} else {
			cmd.Stdin = os.Stdin
		}

		if c.PipeWriter != nil {
			cmd.Stdout = c.PipeWriter
		} else if c.StdOut != "" {
			file, err := os.OpenFile(c.StdOut, c.OpenFlagOut, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stdout = file
		} else {
			cmd.Stdout = os.Stdout
		}

		if c.StdErr != "" {
			file, err := os.OpenFile(c.StdErr, c.OpenFlagErr, 0644)
			if err != nil {
				return err
			}
			defer file.Close()
			cmd.Stderr = file
		} else {
			cmd.Stderr = os.Stderr
		}

		var result error

		if c.PipeChanKill != nil {
			err := cmd.Start()
			if err != nil {
				return err
			}
			ch := make(chan error)

			go func(ch chan error) {
				ch <- cmd.Wait()

			}(ch)

			select {
			case <-c.PipeChanKill:

				if cmd.Process != nil {
					cmd.Process.Kill()
				}

			case result = <-ch:

			}
		} else {
			return cmd.Run()
		}

		return result
	}
}

type executeErrs struct {
	errs []error
	mu   sync.Mutex
}

func (c *Conditionals) ExecutePipe(indexStart int) (int, []error) {
	pipes := make([]*Command, 0)
	errs := &executeErrs{
		errs: make([]error, 0),
	}
	indexEnd := indexStart
	for i := indexStart; i < len(c.Subsequence); i++ {
		if c.Subsequence[i].Pipe {
			indexEnd++
			reader, writer := io.Pipe()
			if i+1 >= len(c.Subsequence) {
				errs.mu.Lock()
				errs.errs = append(errs.errs, fmt.Errorf("out of range"))
				errs.mu.Unlock()
				return indexEnd, errs.errs
			}
			c.Subsequence[i].PipeWriter = writer
			c.Subsequence[i+1].PipeReader = reader
			pipes = append(pipes, c.Subsequence[i])
		} else {
			break
		}
	}
	if indexEnd >= len(c.Subsequence) {
		errs.mu.Lock()
		errs.errs = append(errs.errs, fmt.Errorf("out of range"))
		errs.mu.Unlock()
		return indexEnd, errs.errs
	}

	pipes = append(pipes, c.Subsequence[indexEnd])

	chKill := make(chan struct{})

	for i, cmd := range pipes {
		cmd.PipeChanKill = chKill

		if i < len(pipes)-1 {
			go func(cmd *Command) {
				err := cmd.Execute()
				if cmd.PipeReader != nil {
					cmd.PipeReader.Close()
				}
				if cmd.PipeWriter != nil {
					cmd.PipeWriter.Close()
				}

				if err != nil {
					errs.mu.Lock()
					errs.errs = append(errs.errs, err)
					errs.mu.Unlock()
				}
			}(cmd)
		}
	}

	err := c.Subsequence[indexEnd].Execute()
	c.Subsequence[indexEnd].PipeReader.Close()

	if err != nil {
		errs.mu.Lock()
		errs.errs = append(errs.errs, err)
		errs.mu.Unlock()
	}

	close(chKill)

	return indexEnd, errs.errs
}
