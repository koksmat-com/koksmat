/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func testContext() error {
	kitchenRoot := viper.GetString("KITCHENROOT")
	contextFile := path.Join(kitchenRoot, "mate.json")

	if kitchen.FileExists(contextFile) {
		return nil
	}
	context := `
{
	"current": {
		"tenant": "default"
	},
	
	"sharepoint": [
	],
	"mongo": [
	]
	}
		  
	`
	err := os.WriteFile(contextFile, []byte(context), 0755)

	return err

}

func getConnectionDir(connectionType string) string {
	kitchenRoot := viper.GetString("KITCHENROOT")
	packagePath := path.Join(kitchenRoot, ".koksmat", "tenants", "default", connectionType)
	kitchen.CreateIfNotExists(packagePath, 0755)
	return packagePath
}
func MakeConnectionScript(connectionType string, script string) (string, error) {
	packagePath := getConnectionDir(connectionType)
	psFilePath := path.Join(packagePath, "connect.ps1")
	err := os.WriteFile(psFilePath, []byte(script), 0755)

	return psFilePath, err
}

// serveCmd represents the serve command
var initCmd = &cobra.Command{
	Use:   "init ",
	Short: "init ",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

		err := testContext()
		if err != nil {
			log.Fatalln("Cannot setup context file", err)
		}
		root := viper.GetString("KITCHENROOT")

		SyncConnectorsWithMaster(root)
		//webserver.Run()
	},
}

var kitchenRootCmd = &cobra.Command{
	Use:   "kitchenRoot",
	Short: "kitchenRoot ",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		kitchenRoot := viper.GetString("KITCHENROOT")
		fmt.Println(kitchenRoot)
	},
}

func init() {
	var contextCmd = &cobra.Command{
		Use:   "context",
		Short: "context ",
		Long:  ``}
	rootCmd.AddCommand(contextCmd)
	contextCmd.AddCommand(initCmd)
	contextCmd.AddCommand(kitchenRootCmd)
}
