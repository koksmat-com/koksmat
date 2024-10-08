package autopilot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/koksmat-com/koksmat/auth"
	"github.com/spf13/viper"
)

type Request struct {
	Action       string   `json:"action"`
	Command      string   `json:"command"`
	ReplyTo      string   `json:"reply_to"`
	Args         []string `json:"args"`
	Cwd          string   `json:"cwd"`
	Errormessage string   `json:"errormessage"`
}

type PartialResponse struct {
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
	ReplyTo   string `json:"reply_to"`
	Body      string `json:"body"`
}

type ErrorResponse struct {
	SessionID    string `json:"session_id"`
	Type         string `json:"type"`
	ReplyTo      string `json:"reply_to"`
	ErrorMessage string `json:"error_message"`
	Body         string `json:"body"`
}

type CompleteResponse struct {
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
	ReplyTo   string `json:"reply_to"`
	Body      string `json:"body"`
}

func List() ([]string, error) {
	token := auth.GetToken()
	if token == nil {
		log.Fatal("Error getting token")
	}
	bearerToken := token.AccessToken
	rooturl := viper.GetString("STUDIO_URL")
	url := fmt.Sprintf("%s/api/workspace", rooturl)
	//log.Println("Studio URL:", url)
	requestResponse, err := makeRequest(url, bearerToken)
	if err != nil {
		//log.Println("Error getting workspaces", err)
		return nil, errors.New("Error getting workspaces")
	}
	response := []string{}
	err = json.Unmarshal(requestResponse, &response)
	return response, err
}
func Run(sessionId string, rooturl string) {
	log.Println("Running auto pilot mode with id:", sessionId)
	if rooturl == "" {
		rooturl = viper.GetString("STUDIO_URL")
	}
	url := fmt.Sprintf("%s/api/autopilot/session/%s", rooturl, sessionId)
	log.Println("Studio URL:", url)
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

	for {
		requestResponse, err := makeRequest(url, bearerToken)
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
		case "ping":
			go handlePing(sessionId, request, rooturl, bearerToken)
		case "execute":
			go handleExecute(sessionId, request, rooturl, bearerToken)
		case "execute-nostream":
			go handleExecuteNoStream(sessionId, request, rooturl, bearerToken)
		// case "powershellsession":
		// 	go handlePowerShellHost(sessionId, request, rooturl, bearerToken)
		case "write":
			go handleWrite(sessionId, request, rooturl, bearerToken)
		default:

			log.Println("Unknown command:", request.Command)
		}

		time.Sleep(1 * time.Second) // Loop every second
	}
}

func makeRequest(url, bearerToken string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	return body, nil
}

func handleExecute(sessionId string, request Request, rooturl string, bearerToken string) {
	// Implement your logic for COMMAND_A here
	log.Println("Handling request for session:", sessionId)

	callback := func(isStdOut bool, output string) {
		if isStdOut {
			//log.Println(output)
			postPartialResponse(rooturl, bearerToken, PartialResponse{
				Type:      "append",
				SessionID: sessionId,
				Body:      output,
			},
			)
		} else {
			log.Println("Error:", output)
			postPartialResponse(rooturl, bearerToken, PartialResponse{
				Type:      "append",
				SessionID: sessionId,
				Body:      output,
			},
			)

		}
	}

	result, err := Execute(request.Command, request.Args, Options{Timeout: 30, Cwd: ""}, callback)

	if err != nil {
		log.Println("Error executing command:", err)
		postErrorResponse(rooturl, bearerToken, ErrorResponse{
			Type:         "error",
			SessionID:    sessionId,
			ErrorMessage: fmt.Sprintf("Error executing command: %s", err),
		})
		return
	}
	postResponse(rooturl, bearerToken, CompleteResponse{
		Type:      "done",
		SessionID: sessionId,
		Body:      *result,
	})

}
func handleExecuteNoStream(sessionId string, request Request, rooturl string, bearerToken string) {
	// Implement your logic for COMMAND_A here
	log.Println("Handling request for session:", sessionId)

	// callback := func(isStdOut bool, output string) {
	// 	//return
	// 	url := fmt.Sprintf("%s/api/autopilot", rooturl)
	// 	if isStdOut {
	// 		log.Println(output)
	// 		postPartialResponse(url, bearerToken, PartialResponse{
	// 			Type:      "append",
	// 			SessionID: sessionId,
	// 			ReplyTo:   request.ReplyTo + ".echo",
	// 			Body:      output,
	// 		},
	// 		)
	// 	} else {
	// 		log.Println("Error:", output)
	// 		postPartialResponse(url, bearerToken, PartialResponse{
	// 			Type:      "append",
	// 			SessionID: sessionId,
	// 			ReplyTo:   request.ReplyTo + ".echo",
	// 			Body:      output,
	// 		},
	// 		)

	// 	}
	// }
	result, err := Execute(request.Command, request.Args, Options{Timeout: 30, Cwd: request.Cwd}, nil) //callback)
	url := fmt.Sprintf("%s/api/autopilot", rooturl)

	if err != nil {
		log.Println("Error executing command:", err)
		postErrorResponse(rooturl, bearerToken, ErrorResponse{
			Type:         "error",
			SessionID:    sessionId,
			ReplyTo:      request.ReplyTo,
			Body:         "",
			ErrorMessage: fmt.Sprintf("Error executing command: %s", err),
		})
		return
	}
	postResponse(url, bearerToken, CompleteResponse{
		Type:      "done",
		ReplyTo:   request.ReplyTo,
		SessionID: sessionId,
		Body:      *result,
	})

}

func handlePing(sessionId string, request Request, rooturl string, bearerToken string) {

	log.Println("Handling ping request for session:", sessionId)

	postResponse(rooturl, bearerToken, CompleteResponse{
		Type:      "done",
		SessionID: sessionId,
		ReplyTo:   request.ReplyTo,
		Body:      "pong",
	})

}

func handleWrite(sessionId string, request Request, rooturl string, bearerToken string) {

	log.Println("Handling write request for session:", sessionId)
	req := WriteRequest{}
	json.Unmarshal([]byte(request.Args[0]), &req)

	statuses, _ := handleFileWriteRequest(req)

	body, _ := json.Marshal(statuses)
	postResponse(rooturl, bearerToken, CompleteResponse{
		Type:      "done",
		SessionID: sessionId,
		ReplyTo:   request.ReplyTo,
		Body:      string(body),
	})

}
func postPartialResponse(url, bearerToken string, response PartialResponse) error {

	return postResponseHelper(url, bearerToken, response)
}

func postErrorResponse(url, bearerToken string, response ErrorResponse) error {
	return postResponseHelper(url, bearerToken, response)
}

func postResponse(url, bearerToken string, response CompleteResponse) error {
	return postResponseHelper(url, bearerToken, response)
}

func postResponseHelper(url, bearerToken string, response interface{}) error {
	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status code: %d, response: %s", resp.StatusCode, body)
	}

	return nil
}
