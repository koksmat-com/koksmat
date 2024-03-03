package kitchen

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/koksmat-com/koksmat/connectors"
	"github.com/spf13/viper"
)

type Folder struct {
	Path string `json:"path"`
}
type VSCodeWorkspace struct {
	Folders []Folder `json:"folders"`
}

func GetPath(kitchenName string) (string, error) {
	root := viper.GetString("KITCHENROOT")
	pathName := path.Join(root, kitchenName)
	// if !fileExists(pathName) {
	// 	return "", fmt.Errorf("kitchen %s does not exist", kitchenName)
	// }
	return pathName, nil
}
func CreateKitchen(name string) error {

	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, name)

	err := os.MkdirAll(kitchenPath, 0755)
	if err != nil {
		return err
	}

	readmeFilePath := path.Join(kitchenPath, "readme.md")

	if FileExists(readmeFilePath) {
		return fmt.Errorf("kitchen %s already exists", name)
	}
	data := fmt.Sprintf(`---
title: %s
description: Describe the main purpose of this kitchen
---

# %s
`, name, name)

	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.WriteFile(readmeFilePath, []byte(data), 0644)
	if err != nil {
		return err
	}
	os.Chdir(kitchenPath)
	return CreateWorkspaceFile()
}

func BuildKitchen(name string) error {

	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, name)

	goFilePath := path.Join(kitchenPath, "main.go")

	if FileExists(goFilePath) {
		_, err := connectors.Execute("go", *&connectors.Options{Dir: kitchenPath}, "install")
		if err != nil {
			return err
		}
	}
	/*output, err := connectors.Execute(name, *&connectors.Options{}, "")
	if err != nil {
		return err
	}
	fmt.Println(string(output))*/

	return nil
}

func AddTemplatesToKitchen(name string) error {

	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, name)

	_, err := connectors.Execute("go", *&connectors.Options{Dir: kitchenPath}, "mod", "init", fmt.Sprintf("github.com/365admin/%s", name))
	if err != nil {
		return err
	}
	return nil
}

type KitchenStatus struct {
	Name     string
	Branch   string
	Revision string
	Status   string
}

func GetKitchenStatus(name string) (KitchenStatus, error) {

	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, name)

	branch, err := connectors.Execute("git", *&connectors.Options{Dir: kitchenPath}, "branch")
	if err != nil {
		return KitchenStatus{}, err
	}
	revision, err := connectors.Execute("git", *&connectors.Options{Dir: kitchenPath}, "rev-parse", "HEAD")
	if err != nil {
		return KitchenStatus{}, err
	}
	status, err := connectors.Execute("git", *&connectors.Options{Dir: kitchenPath}, "status", "-s")
	if err != nil {
		return KitchenStatus{}, err
	}
	return KitchenStatus{Name: name, Branch: string(branch), Revision: string(revision), Status: string(status)}, nil
}

func MakeRelease(name string, tagname string) error {

	root := viper.GetString("KITCHENROOT")
	kitchenPath := path.Join(root, name)

	_, err := connectors.Execute("git", *&connectors.Options{Dir: kitchenPath}, "tag", tagname)
	if err != nil {
		return err
	}
	_, err = connectors.Execute("git", *&connectors.Options{Dir: kitchenPath}, "push", "origin", tagname)
	if err != nil {
		return err
	}
	_, err = connectors.Execute("gh", *&connectors.Options{Dir: kitchenPath}, "release", "create", tagname, "-t", tagname)
	if err != nil {
		return err
	}

	return nil
}

func CreateWorkspaceFile() error {
	root := viper.GetString("KITCHENROOT")
	wsPath := path.Join(root, "kitchens.code-workspace")
	workspaces := VSCodeWorkspace{}
	dirs, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if dir.IsDir() && !strings.HasPrefix(dir.Name(), ".") {
			n := dir.Name()
			workspaces.Folders = append(workspaces.Folders, Folder{Path: n})
		}
	}

	data, err := json.MarshalIndent(workspaces, "", " ")
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = os.WriteFile(wsPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
