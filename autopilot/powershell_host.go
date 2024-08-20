package autopilot

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"time"

	"github.com/koksmat-com/koksmat/auth"
	"github.com/spf13/viper"
)

// FetchCommand fetches the next command from the external service
func FetchCommand(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching command: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func stripANSI(s string) string {
	// Define the ANSI escape sequence pattern
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	re := regexp.MustCompile(ansi)
	// Replace all matches with an empty string
	return re.ReplaceAllString(s, "")
}

func PowerShellHost(bootCommands []string, sessionId string) error {
	rooturl := viper.GetString("STUDIO_URL")
	commandsUrl := fmt.Sprintf("%s/api/autopilot/session/%s", rooturl, sessionId)
	log.Println("Studio URL:", commandsUrl)
	var bearerToken string
	ticker := time.NewTicker(10 * time.Minute)
	quit := make(chan struct{})

	// Initial token retrieval
	token := auth.GetToken()
	if token == nil {
		log.Fatal("Error getting token")
	}
	bearerToken = token.AccessToken

	go func() {
		for {
			select {
			case <-ticker.C:
				token := auth.GetToken()
				if token == nil {
					log.Fatal("Error getting token")
				}
				bearerToken = token.AccessToken
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	callback := func(isStdOut bool, output string) {

		log.Println(isStdOut, output)
		// url := fmt.Sprintf("%s/api/autopilot", rooturl)

		// if !isStdOut {
		// 	postErrorResponse(rooturl, bearerToken, ErrorResponse{
		// 		Type:      "error",
		// 		SessionID: sessionId,
		// 		ReplyTo:   request.ReplyTo,
		// 		Body:      output,
		// 	})
		// 	return
		// }
		// postResponse(url, bearerToken, CompleteResponse{
		// 	Type:      "done",
		// 	ReplyTo:   request.ReplyTo,
		// 	SessionID: sessionId,
		// 	Body:      output,
		// })
	}
	//PowerShellHost(request.Command, request.Args, Options{Timeout: 30, Cwd: request.Cwd}, callback)

	ctx := context.Background()

	// Start a new PowerShell process
	cmd := exec.CommandContext(ctx, "pwsh", "-NoExit", "-Command", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {

		return err

	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {

		return err
	}

	// Start the PowerShell process
	if err := cmd.Start(); err != nil {

		return err
	}
	log.Println("pid", cmd.Process.Pid)
	// Create a scanner to read the output
	scanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	_, err = fmt.Fprintln(stdin, `
$PSStyle.OutputRendering = [System.Management.Automation.OutputRendering]::PlainText;

 	$ErrorActionPreference = "Continue";
	# $ProgressPreference = "SilentlyContinue"
	# $VerbosePreference = "SilentlyContinue"
	# $DebugPreference = "SilentlyContinue"

	`)

	if err != nil {
		return err
	}
	for _, command := range bootCommands {
		_, err := fmt.Fprintln(stdin, command)
		if err != nil {
			return err
		}
		//time.Sleep(2 * time.Second) // Adjust the delay as necessary
	}

	// Read the output
	go func() {
		for scanner.Scan() {
			t := scanner.Text()
			var justText = stripANSI(t)

			callback(false, justText)
			//fmt.Println(scanner.Text())
		}
	}()

	go func() {
		for errScanner.Scan() {
			t := scanner.Text()
			var justText = stripANSI(t)
			callback(true, justText)
			//fmt.Fprintln(os.Stderr, errScanner.Text())
		}
	}()
	// Fetch and execute commands
	for {
		requestResponse, err := makeRequest(commandsUrl, bearerToken)
		if err != nil {
			log.Println("Error making request:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var request Request

		err = json.Unmarshal(requestResponse, &request)
		if err != nil {
			log.Println("Error unmarshalling response:", err)
			log.Println("Response:", string(requestResponse))
			time.Sleep(1 * time.Second)
			continue
		}
		if request.Errormessage == "timeout" {
			log.Println("No request found, waiting for new request")
			continue
		}
		log.Println("Request:", string(requestResponse))
		switch request.Action {
		case "execute-nostream":

			// iterate over requests.Args and send them to stdin

			for _, arg := range request.Args {
				//_, err =
				log.Println("Sending command:", arg)
				fmt.Fprintln(stdin, arg)
				// if err != nil {
				// 	return err
				// }
			}

			// cmd := strings.Join(request.Args, "\n")
			// _, err = fmt.Fprintln(stdin, cmd)
		// case "ping":
		// 	go handlePing(sessionId, request, rooturl, bearerToken)
		// case "execute":
		// 	go handleExecute(sessionId, request, rooturl, bearerToken)
		// case "execute-nostream":
		// 	go handleExecuteNoStream(sessionId, request, rooturl, bearerToken)
		// // case "powershellsession":
		// // 	go handlePowerShellHost(sessionId, request, rooturl, bearerToken)
		// case "write":
		// 	go handleWrite(sessionId, request, rooturl, bearerToken)
		default:

			log.Println("Unknown command:", request.Command)
		}

		time.Sleep(1 * time.Second) // Loop every second
	}

	// Close stdin to signal that no more input will be sent
	if err := stdin.Close(); err != nil {

		return err
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return err

	}
	return nil
}
