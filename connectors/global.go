package connectors

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"
)

type Connector struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
	Url         string `json:"url"`
	JSON        any    `json:"json"`
	IsCurrent   bool   `json:"isCurrent"`

	// Name of the connector
}

type Options struct {
	Channel string
	Dir     string
	Env     []string
}

type MessageData struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	IsError   bool      `json:"isError"`
}
type Message struct {
	Channel string      `json:"channel"`
	Data    MessageData `json:"data"`
}

func Execute(program string, options Options, args ...string) (output []byte, err error,
) {

	cmd := exec.Command(program, args...)
	if options.Dir != "" {
		cmd.Dir = options.Dir
	}
	if options.Env != nil {
		cmd.Env = options.Env
	}
	pipe, _ := cmd.StdoutPipe()
	combinedOutput := []byte{}

	err = cmd.Start()
	if err != nil {
		return nil, errors.New(fmt.Sprint(err))
	}
	go func(p io.ReadCloser) {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')

		for err == nil {
			//log.Print(line)
			combinedOutput = append(combinedOutput, []byte(line)...)
			line, err = reader.ReadString('\n')
		}
		//
	}(pipe)
	err = cmd.Wait()

	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + string(combinedOutput))
		return nil, errors.New(fmt.Sprint(err))
	}

	return combinedOutput, nil
}
