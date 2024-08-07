package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"strings"
)

type AccessToken struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

func GetToken() *AccessToken {
	var token AccessToken
	retry := true

	for retry {
		retry = false
		cmd := exec.Command("az", "account", "get-access-token")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		if err != nil {
			//fmt.Println("Error executing command:", err)
			errString := stderr.String()
			if strings.Contains(errString, "AADSTS50078") {
				fmt.Println("Error: You need to reauthenticate")
				login := exec.Command("az", "login", "--scope", "https://management.core.windows.net//.default")
				loginOutput, err := login.CombinedOutput()
				if err != nil {
					fmt.Println("Could not reauthenticate:", err)
					fmt.Println("Output:", string(loginOutput))
					return nil
				}
				retry = true

			} else {
				fmt.Println("Error executing command:", errString)
			}

			return nil
		}

		if err := json.Unmarshal(output, &token); err != nil {
			fmt.Println("Error parsing JSON:", err)

			return nil
		}
	}
	// fmt.Println("Access Token:", token.AccessToken)
	// fmt.Println("Token Type:", token.TokenType)
	// fmt.Println("Expires In:", token.ExpiresIn)
	return &token
}
