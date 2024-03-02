package powershell

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/google/uuid"
)

//go:embed scripts
var scripts embed.FS

type Setup func(workingDirectory string) (string, []string, error)

func PwshCwd(appId string) string {

	dir := ".koksmat/powershell"
	os.MkdirAll(dir, os.ModePerm)
	dir = path.Join(dir, fmt.Sprintf("%s-%s", appId, uuid.New()))
	os.MkdirAll(dir, os.ModePerm)

	return dir
}

type Callback func(workingDirectory string)

func CallbackMockup(workingDirectory string) {}
func Execute(appId string, fileName, args string, setEnvironment Setup, src string, callback Callback) (output []byte, err error, console string,
) {
	cmd := exec.Command("pwsh", "-nologo", "-noprofile")
	workingDirectory := PwshCwd(appId)
	initScript, environment, err := setEnvironment(workingDirectory)
	if err != nil {
		return nil, err, ""
	}

	cmd.Env = environment

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Dir = workingDirectory

	os.Remove(path.Join(cmd.Dir, "output.json"))
	ps1Code, err := scripts.ReadFile(fmt.Sprintf("scripts/connectors/%s.ps1", initScript))
	if err != nil {

		return nil, err, ""
	}

	ps2Code := []byte(src)
	if src == "" {
		ps2Code, err = scripts.ReadFile(fileName)
		if err != nil {
			return nil, err, ""
		}
	}

	err = os.WriteFile(path.Join(cmd.Dir, "run.ps1"), ps2Code, 0644)
	if err != nil {
		return nil, err, ""
	}
	err = os.WriteFile(path.Join(cmd.Dir, "init.ps1"), ps1Code, 0644)
	if err != nil {
		return nil, err, ""
	}

	pipe, _ := cmd.StdoutPipe()
	combinedOutput := []byte{}

	script := fmt.Sprintf(`
	$ErrorActionPreference = "Stop"
	. ./run.ps1  %s`, args)
	go func() {
		defer stdin.Close()
		fmt.Fprintln(stdin, ". ./init.ps1")
		fmt.Fprintln(stdin, script)

	}()

	err = cmd.Start()
	go func(p io.ReadCloser) {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')
		for err == nil {
			log.Print(line)
			combinedOutput = append(combinedOutput, []byte(line)...)
			line, err = reader.ReadString('\n')
		}
	}(pipe)
	err = cmd.Wait()

	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + string(combinedOutput))
		return nil, errors.New("Could not run PowerShell script"), string(combinedOutput)
	}

	outputJson, err := os.ReadFile(path.Join(cmd.Dir, "output.json"))

	os.WriteFile(path.Join(cmd.Dir, "output.txt"), []byte(string(combinedOutput)), 0644)
	if callback != nil {
		callback(cmd.Dir)
	}

	return outputJson, nil, string(combinedOutput)
}

func Run[R any](appId string, fileName string, args string, setup Setup, src string, callback Callback) (result *R, err error) {

	output, err, _ := Execute(appId, fileName, args, setup, src, callback)
	dataOut := new(R)
	textOutput := fmt.Sprintf("%s", output)
	if (output != nil) && (textOutput != "") {

		jsonErr := json.Unmarshal(output, &dataOut)
		if jsonErr != nil {
			s := fmt.Sprintf("[%s]", output)
			outArray := []byte(s)
			jsonErr := json.Unmarshal(outArray, &dataOut)
			if jsonErr != nil {
				log.Println("Error parsing output: ", jsonErr)
			}
		}
	}
	result = *&dataOut // fmt.Sprintf("%s", outputJson)
	return result, err
}
func RunRaw(appId string, fileName string, args string, setup Setup, src string, callback Callback) (result string, err error) {

	output, err, _ := Execute(appId, fileName, args, setup, src, callback)
	result = fmt.Sprintf("%s", output)

	return result, err
}
