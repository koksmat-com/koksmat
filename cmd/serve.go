/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/koksmat-com/koksmat/service"
	"github.com/spf13/cobra"
)

func serve() {
	service.Serve()
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve ",
	Short: "Auto pilot mode",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Service started on port 8080")
		serve()
		// restapi.Sail()
		//webserver.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().BoolP("port", "p", false, "Port to listen on")
}
