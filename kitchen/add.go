package kitchen

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/koksmat-com/koksmat/input"
)

func checkCurrentDirectory() error {
	// Check if package.json file exists
	if _, err := os.Stat("package.json"); os.IsNotExist(err) {
		return errors.New("not a koksmat project: package.json file is missing")
	} else if err != nil {
		return err
	}

	// Check if components directory exists
	if stat, err := os.Stat("components"); os.IsNotExist(err) {
		return errors.New("not a koksmat project: components directory is missing")
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return errors.New("not a koksmat project: components is not a directory")
	}

	return nil
}
func AddIngredients(key string) error {
	// Add ingredients to your project

	// Check if the current directory is a Koksmat project
	err := checkCurrentDirectory()
	if err != nil {
		return err
	}

	// Make the HTTP GET request
	url := "http://localhost:4339/api/package/" + key
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	name := input.GetString("Component", "Enter the component name")

	log.Println(name, string(body))

	// Return the raw JSON as a byte slice
	return nil
}
