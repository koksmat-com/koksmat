package kitchen

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/koksmat-com/koksmat/connectors"
	"github.com/spf13/viper"
)

func Ship(kitchenName string) (string, error) {
	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, kitchenName)
	goModPath := path.Join(kitchenPath, "go.mod")

	if !FileExists(goModPath) {
		return "", fmt.Errorf("go.mod not found in %s", kitchenPath)
	}
	fileContent, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(fileContent), "\n")
	name := strings.Split(lines[0], "module ")

	execname := filepath.Base(name[1])

	_, err = connectors.Execute("go", *&connectors.Options{Dir: kitchenPath}, "install")
	if err != nil {
		return "", err
	}

	bytes, err := connectors.Execute(execname, *&connectors.Options{}, "")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`Build %s - %s`, execname, string(bytes)), nil
}
