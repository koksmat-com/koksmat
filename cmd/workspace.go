/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var workspaceCmd = &cobra.Command{
	Use:   "workspace [workspace [path]",
	Short: "workspace",
	Long:  `workspace`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No service specified")
			return
		}
		workspace := args[0]
		path := args[1]
		switch path {
		case "status":
			printData(kitchen.GetStatus(workspace, true))
		// case "github/organisations":
		// 	printData(connectors.GithubOrgs())
		// case "github/repositories":
		// 	printData(connectors.GithubRepos())
		// case "github/codespaces":
		// 	printData(connectors.GithubCodespaces())
		// case "azure/subscriptions":
		// 	printData(connectors.AzureSubscriptions())
		// case "azure/storageaccounts":
		// 	printData(connectors.AzureStorageAccounts())
		// case "kubernetes/clusters":
		// 	printData(connectors.KubernetesClusters())
		// case "kubernetes/namespaces":
		// 	printData(connectors.KubernetesNamespaces())
		// case "kubernetes/pods":
		// 	printData(connectors.KubernetesPods())
		// case "sharepoint/tenants":
		// 	printData(connectors.SharePointTenants())
		// case "microsoft365/tenants":
		// 	printData(connectors.M365Context())
		// case "microsoft365/sites":
		// 	printData(connectors.M365Sites())
		default:

			log.Fatalln("Unknown ", path)
			return
		}
		//webserver.Run()
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
