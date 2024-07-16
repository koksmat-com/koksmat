package autopilot

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

type Options struct {
	Timeout int
	Cwd     string
}

func ExecuteCommand(cmd string, args []string, options Options, callback func(isStdOut bool, output string)) (*string, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(options.Timeout)*time.Second)
	defer cancel()

	// Create the command with context
	command := exec.CommandContext(ctx, cmd, args...)
	if options.Cwd != "" {
		command.Dir = options.Cwd
	}

	stdoutPipe, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPipe, err := command.StderrPipe()
	if err != nil {
		return nil, err
	}

	var combinedOutput bytes.Buffer

	// Start the command
	err = command.Start()
	if err != nil {
		return nil, err
	}

	// Function to read from the pipe and send to the callback
	readPipe := func(pipe *bufio.Scanner, isStdOut bool) {
		for pipe.Scan() {
			line := pipe.Text()
			if callback != nil {
				callback(isStdOut, line)
			}
			combinedOutput.WriteString(line + "\n")
		}
	}

	// Create scanners for stdout and stderr
	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stderrScanner := bufio.NewScanner(stderrPipe)

	// Start reading stdout and stderr concurrently
	go readPipe(stdoutScanner, true)
	go readPipe(stderrScanner, false)

	// Wait for the command to finish
	err = command.Wait()
	if err != nil {
		// Handle the timeout error
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.New("command timed out")
		}
	}

	// Combine the output
	combinedOutputStr := combinedOutput.String()

	// Return the combined output and error
	return &combinedOutputStr, err
}

/*
*
  - Execute a command with arguments and options

sample usage:

	 cmd := "ls"
		args := []string{"-la"}
		options := Options{
			Timeout: 5,
			Cwd:     "",
		}
*/
func Execute(cmd string, args []string, options Options, callback func(isStdOut bool, output string)) (*string, error) {

	return ExecuteCommand(cmd, args, options, callback)

}
