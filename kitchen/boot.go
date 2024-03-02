package kitchen

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/koksmat-com/koksmat/global"
	"github.com/spf13/viper"
)

var verbose bool = false

var yellow = color.New(color.FgYellow).PrintlnFunc()

var green = color.New(color.FgGreen).PrintlnFunc()

func verboseLog(s ...interface{}) {
	if verbose {
		yellow(s)
	}
}

func walkPath(startingPath string, fileToLookFor string) (string, error) {

	verboseLog("Looking for file: ", fileToLookFor, " in path: ", startingPath)

	filePath := path.Join(startingPath, fileToLookFor)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if startingPath == "/" {
			return "file does not exist", err
		}

		return walkPath(path.Dir(startingPath), fileToLookFor)
	}

	verboseLog("Found at ", startingPath)

	return startingPath, nil
}

/**
 * Check if a .env file exists in the current directory or any of its parent directories
 */
func envFileCheck() (string, error) {
	filename := ".env"
	startingPath, err := os.Getwd()
	if err != nil {
		return "error", err
	}
	foundOn, err := walkPath(startingPath, filename)
	if err != nil {
		yellow("No default .env file")
		return "", nil

	}
	return path.Join(foundOn, filename), nil

}

func createDefaultEnvFile(envFileName string) (string, error) {
	verboseLog("Creating a default .env file")

	defaultEnvironment := "# This is a default .env file\n"

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	kitchenRoot := path.Dir(wd)
	verboseLog("Kitchen root will be", kitchenRoot)
	defaultEnvironment += "KITCHENROOT=" + kitchenRoot + "\n"

	defaultEnvFilepath := path.Join(kitchenRoot, envFileName)

	verboseLog("Default path will be", defaultEnvFilepath)

	if _, err := os.Stat(defaultEnvFilepath); os.IsNotExist(err) {
		if err != nil {
			verboseLog("Creating file")
			yellow("Creating default .env file at", defaultEnvFilepath)
			os.WriteFile(path.Join(kitchenRoot, envFileName), []byte(defaultEnvironment), 0644)

		}

	} else {
		verboseLog("File already exists")
	}

	return defaultEnvFilepath, nil
}

type BootOptions struct {
	Verbose bool
}

func logo() {
	fmt.Println(global.Logo)
}

func checkEnv() {
	envFilePath, err := envFileCheck()
	if err != nil {
		log.Fatal(err)
	}
	if envFilePath == "" {
		configPath, err := createDefaultEnvFile(".env")
		if err != nil {
			log.Fatal(err)
		}
		viper.SetConfigFile(configPath)
	}
}

func setup(options BootOptions) {
	verbose = options.Verbose

}
func Boot(options BootOptions) error {
	setup(options)
	checkEnv()
	kitchenRoot := viper.GetString("KITCHENROOT")
	verboseLog("Kitchen root is", kitchenRoot)

	return nil
}
