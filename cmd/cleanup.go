/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/koksmat-com/koksmat/restapi"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup [[service]]",
	Short: "Cleanup service ",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			restapi.Run()
			return
		}
		service := args[0]
		switch service {
		case "booking":
			// restapi.All()

		default:

			log.Fatalln("Unknown service", service)
			return
		}
		//webserver.Run()
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
