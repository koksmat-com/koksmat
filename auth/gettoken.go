package auth

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type AccessToken struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

func GetToken() *AccessToken {
	cmd := exec.Command("az", "account", "get-access-token")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return nil
	}

	var token AccessToken
	if err := json.Unmarshal(output, &token); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// fmt.Println("Access Token:", token.AccessToken)
	// fmt.Println("Token Type:", token.TokenType)
	// fmt.Println("Expires In:", token.ExpiresIn)
	return &token
}
