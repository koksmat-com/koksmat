/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runCommand(command string, args ...string) {
	execCmd := exec.Command(command, args...)

	output, err := execCmd.CombinedOutput()

	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(output))
}

func sail() {
	kitchenRoot := viper.GetString("KITCHENROOT")
	packagePath := path.Join(kitchenRoot, ".koksmat", "packages")
	runCommand("open", "http://localhost:3000")
	execCmd := exec.Command("pnpm", "start")
	execCmd.Dir = path.Join(packagePath, "koksmat-mate", ".koksmat", "web")
	pipe, _ := execCmd.StdoutPipe()
	err := execCmd.Start()
	go func(p io.ReadCloser) {
		reader := bufio.NewReader(pipe)
		line, err := reader.ReadString('\n')
		for err == nil {
			log.Print(line)

			line, err = reader.ReadString('\n')
		}
	}(pipe)
	err = execCmd.Wait()

	if err != nil {
		log.Fatal(err)
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
