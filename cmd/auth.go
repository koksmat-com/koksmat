/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/koksmat-com/koksmat/auth"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var authCmd = &cobra.Command{
	Use:   "auth ",
	Short: "Auto pilot mode",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		subscriptions, err := auth.GetSubscriptions()

		if err != nil {
			log.Fatal(err)
		}

		auth.InteractivlySelectSubscription(subscriptions)

	},
}

func init() {
	rootCmd.AddCommand(authCmd)

}
