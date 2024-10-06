/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/koksmat-com/koksmat/auth"
	"github.com/koksmat-com/koksmat/input"
	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var autopilotCmd = &cobra.Command{
	Use:   "auto ",
	Short: "Auto pilot mode",
	Args:  cobra.MinimumNArgs(0),
	Long: `
	
	
	
	`,

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
	autopilotCmd.AddCommand(kitchen.NewAutopilotConfigRootCommand())
	autopilotCmd.AddCommand(&cobra.Command{
		Use: "run [connectionId] [[studiourl]]",

		Args:    cobra.MinimumNArgs(0),
		Long:    ``,
		Example: `koksmat auto run`,

		Run: func(cmd *cobra.Command, args []string) {
			studioUrl := ""
			connectionId := ""
			if len(args) == 0 {

				defaultConnectionId, err := kitchen.GetAutopilotDefaultConnection()
				connectionId = defaultConnectionId
				if err != nil {
					log.Fatal(err)
				}
			} else {
				connectionId = args[0]
			}
			if len(args) > 1 {
				studioUrl = args[1]
			}
			kitchen.AutoPilotRun(connectionId, studioUrl)
		},
	})

	pwshCmd := &cobra.Command{
		Use:              "pwsh ",
		TraverseChildren: true,
	}

	var connectionId string
	pwshCmd.PersistentFlags().StringVarP(&connectionId, "connectionId", "c", "", "Connection id")
	autopilotCmd.AddCommand(pwshCmd)

	pwshCmd.AddCommand(&cobra.Command{
		Use: "host",

		Args:    cobra.MinimumNArgs(0),
		Long:    ``,
		Example: `koksmat auto pwsh host connect-exchangeonline -connectionId 124`,

		Run: func(cmd *cobra.Command, args []string) {

			kitchen.PowerShellHost(connectionId, args)
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
			input.Confirm()
			// _, err := kitchen.AddAutopilotConnection(token.AccessToken, "https://koksmat.com")
			// if err != nil {
			// 	log.Fatal(err)
			// }
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
