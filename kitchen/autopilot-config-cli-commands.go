package kitchen

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*
CliCommands

This file contains the Cobra CLI command implementations for managing
KoksmatAutoPilotHosts configurations and testing connections. It provides
commands for listing, adding, updating, and deleting host configurations,
as well as for testing connections using the AutopilotConnectionTester.
*/

var (
	configFile string
	jwtToken   string
)

// NewRootCommand creates the root command for the CLI
func NewAutopilotConfigRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "config",
		Short: "Koksmat CLI for managing AutoPilot hosts",
		Long:  `A CLI tool for managing Koksmat AutoPilot host configurations and testing connections.`,
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "koksmat_config.yaml", "config file (default is koksmat_config.yaml)")
	rootCmd.PersistentFlags().StringVar(&jwtToken, "token", "", "JWT token for authentication")

	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newAddCommand())
	rootCmd.AddCommand(newUpdateCommand())
	rootCmd.AddCommand(newDeleteCommand())
	rootCmd.AddCommand(newTestCommand())

	return rootCmd
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured hosts",
		Run: func(cmd *cobra.Command, args []string) {
			hosts := NewKoksmatAutoPilotHosts()
			err := hosts.LoadFromFile(configFile)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			fmt.Println("Configured hosts:")
			for tag, host := range hosts.config.Hosts {
				fmt.Printf("- %s: %s\n", tag, host.Href)
			}
			fmt.Printf("Default host: %s\n", hosts.config.DefaultHost)
		},
	}
}

func newAddCommand() *cobra.Command {
	var tag, href, key string
	var setDefault bool

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new host configuration",
		Run: func(cmd *cobra.Command, args []string) {
			hosts := NewKoksmatAutoPilotHosts()
			err := hosts.LoadFromFile(configFile)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			hosts.AddHost(tag, href, key)
			if setDefault {
				hosts.SetDefaultHost(tag)
			}

			err = hosts.SaveToFile(configFile)
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}

			fmt.Printf("Host '%s' added successfully\n", tag)
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "", "Tag for the new host")
	cmd.Flags().StringVar(&href, "href", "", "HREF for the new host")
	cmd.Flags().StringVar(&key, "key", "", "Key for the new host")
	cmd.Flags().BoolVar(&setDefault, "default", false, "Set as default host")
	cmd.MarkFlagRequired("tag")
	cmd.MarkFlagRequired("href")
	cmd.MarkFlagRequired("key")

	return cmd
}

func newUpdateCommand() *cobra.Command {
	var tag, href, key string
	var setDefault bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an existing host configuration",
		Run: func(cmd *cobra.Command, args []string) {
			hosts := NewKoksmatAutoPilotHosts()
			err := hosts.LoadFromFile(configFile)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			if _, ok := hosts.config.Hosts[tag]; !ok {
				fmt.Printf("Host '%s' not found\n", tag)
				return
			}

			hosts.AddHost(tag, href, key)
			if setDefault {
				hosts.SetDefaultHost(tag)
			}

			err = hosts.SaveToFile(configFile)
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}

			fmt.Printf("Host '%s' updated successfully\n", tag)
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "", "Tag of the host to update")
	cmd.Flags().StringVar(&href, "href", "", "New HREF for the host")
	cmd.Flags().StringVar(&key, "key", "", "New key for the host")
	cmd.Flags().BoolVar(&setDefault, "default", false, "Set as default host")
	cmd.MarkFlagRequired("tag")

	return cmd
}

func newDeleteCommand() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a host configuration",
		Run: func(cmd *cobra.Command, args []string) {
			hosts := NewKoksmatAutoPilotHosts()
			err := hosts.LoadFromFile(configFile)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			if _, ok := hosts.config.Hosts[tag]; !ok {
				fmt.Printf("Host '%s' not found\n", tag)
				return
			}

			delete(hosts.config.Hosts, tag)
			if hosts.config.DefaultHost == tag {
				hosts.config.DefaultHost = ""
			}

			err = hosts.SaveToFile(configFile)
			if err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				return
			}

			fmt.Printf("Host '%s' deleted successfully\n", tag)
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "", "Tag of the host to delete")
	cmd.MarkFlagRequired("tag")

	return cmd
}

func newTestCommand() *cobra.Command {
	var tag string

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Test connection to a host",
		Run: func(cmd *cobra.Command, args []string) {
			hosts := NewKoksmatAutoPilotHosts()
			err := hosts.LoadFromFile(configFile)
			if err != nil {
				fmt.Printf("Error loading config: %v\n", err)
				return
			}

			hostConfig, err := hosts.GetHost(tag)
			if err != nil {
				fmt.Printf("Error getting host config: %v\n", err)
				return
			}

			tester := NewAutopilotConnectionTester(hostConfig, jwtToken)

			fmt.Println("Testing connection...")
			err = tester.Ping()
			if err != nil {
				fmt.Printf("Ping failed: %v\n", err)
				return
			}
			fmt.Println("Ping successful")

			fmt.Println("Registering connection...")
			err = tester.RegisterConnection()
			if err != nil {
				fmt.Printf("Registration failed: %v\n", err)
				return
			}
			fmt.Println("Registration successful")

			fmt.Println("Getting status...")
			status, err := tester.GetStatus()
			if err != nil {
				fmt.Printf("Failed to get status: %v\n", err)
				return
			}
			fmt.Printf("Status: %s\n", status)
		},
	}

	cmd.Flags().StringVar(&tag, "tag", "", "Tag of the host to test")
	cmd.MarkFlagRequired("tag")

	return cmd
}
