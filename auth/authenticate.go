// auth/auth.go
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/cqroot/prompt"
)

type Subscription struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CheckErr(err error) {
	if err != nil {
		if errors.Is(err, prompt.ErrUserQuit) {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		} else {
			panic(err)
		}
	}
}
func GetSubscriptions() ([]Subscription, error) {
	cmd := exec.Command("az", "account", "list", "--output", "json", "--only-show-errors")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var subscriptions []Subscription
	if err := json.Unmarshal(output, &subscriptions); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func ListSubscriptions() {
	subscriptions, err := GetSubscriptions()
	if err != nil {
		fmt.Println("Error fetching subscriptions:", err)
		return
	}

	fmt.Println("Subscriptions:")
	for _, sub := range subscriptions {
		fmt.Printf("ID: %s, Name: %s\n", sub.ID, sub.Name)
	}

	// Other CLI functionalities can be implemented here
}

func InteractivlySelectSubscription(subs []Subscription) {
	choices := []string{}
	for _, s := range subs {
		choices = append(choices, fmt.Sprintf("%s | %s", s.Name, s.ID))
	}

	val1, err := prompt.New().Ask("Choose:").
		Choose(choices)
	CheckErr(err)

	fmt.Printf("{ %s }\n", val1)

}
