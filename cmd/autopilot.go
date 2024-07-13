/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/koksmat-com/koksmat/auth"
	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var autopilotCmd = &cobra.Command{
	Use:   "auto ",
	Short: "Auto pilot mode",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,

	// Run: func(cmd *cobra.Command, args []string) {
	// 	subscriptions, err := auth.GetSubscriptions()

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	auth.InteractivlySelectSubscription(subscriptions)

	// },
}

func init() {

	rootCmd.AddCommand(autopilotCmd)
	autopilotCmd.AddCommand(&cobra.Command{
		Use: "run ",

		Args: cobra.MinimumNArgs(0),
		Long: ``,

		Run: func(cmd *cobra.Command, args []string) {

			connectionId, err := kitchen.GetAutopilotDefaultConnection()
			if err != nil {
				log.Fatal(err)
			}
			kitchen.AutoPilotRun(connectionId)
		},
	})

	autopilotCmd.AddCommand(&cobra.Command{
		Use:  "list ",
		Args: cobra.MinimumNArgs(0),
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			connections, err := kitchen.ListAutopilotConnections()

			if err != nil {
				log.Fatal(err)
			}

			for _, connection := range connections {
				fmt.Println(connection)
			}
		},
	})

	autopilotCmd.AddCommand(&cobra.Command{
		Use:  "add ",
		Args: cobra.MinimumNArgs(0),
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			token := auth.GetToken()
			if token == nil {
				log.Fatal("Error getting token")
			}

			_, err := kitchen.AddAutopilotConnection(token.AccessToken, "https://koksmat.com")
			if err != nil {
				log.Fatal(err)
			}
		},
	})

	autopilotCmd.AddCommand(&cobra.Command{
		Use:  "remove ",
		Args: cobra.MinimumNArgs(1),
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := kitchen.RemoveAutopilotConnection(args[0])
			if err != nil {
				log.Fatal(err)

			}

		},
	})

	autopilotCmd.AddCommand(&cobra.Command{
		Use:  "set ",
		Args: cobra.MinimumNArgs(1),
		Long: ``,
		Run: func(cmd *cobra.Command, args []string) {
			err := kitchen.SetAutopilotDefaultConnection(args[0])
			if err != nil {
				log.Fatal(err)
			}
		},
	})

}
