/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func sail() {
	kitchenRoot := viper.GetString("KITCHENROOT")
	packagePath := path.Join(kitchenRoot, ".koksmat", "packages")

	execCmd := exec.Command("pnpm", "start")
	execCmd.Dir = path.Join(packagePath, "koksmat-mate", ".koksmat", "web")
	execResult := execCmd.Run()
	if execResult.Error() != "" {
		log.Fatal(execResult.Error())
		return
	}

	//filename := UnEscape(args[0])
}

// serveCmd represents the serve command
var sailCmd = &cobra.Command{
	Use:   "sail ",
	Short: "Auto pilot mode",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Sailing")
		sail()
		// restapi.Sail()
		//webserver.Run()
	},
}

func init() {
	rootCmd.AddCommand(sailCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
