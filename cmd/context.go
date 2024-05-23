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
	Use:   "init [service]",
	Short: "init ",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No service specified")
			return
		}
		err := testContext()
		if err != nil {
			log.Fatalln("Cannot setup context file", err)
		}
		connectionType := args[0]
		switch connectionType {
		case "exchange":
			log.Println("Exchange")
			scriptPath, err := MakeConnectionScript(connectionType, `

$EXCHAPPID = $env:EXCHAPPID
$EXCHORGANIZATION = $env:EXCHORGANIZATION
$EXCHCERTIFICATEPASSWORD = $env:EXCHCERTIFICATEPASSWORD
$EXCHCERTIFICATEPATH = "$PSScriptRoot/certificate.pfx"
$bytes = [Convert]::FromBase64String($ENV:EXCHCERTIFICATE)
[IO.File]::WriteAllBytes($EXCHCERTIFICATEPATH, $bytes)

Write-Output "Connecting to Exchange for $EXCHORGANIZATION"

if (($EXCHCERTIFICATEPASSWORD -ne $null) -and ($EXCHCERTIFICATEPASSWORD -ne "") ){
	Connect-ExchangeOnline -CertificateFilePath $EXCHCERTIFICATEPATH  -AppID $EXCHAPPID -Organization $EXCHORGANIZATION -ShowBanner:$false -CertificatePassword (ConvertTo-SecureString -String $EXCHCERTIFICATEPASSWORD -AsPlainText -Force)
}else{
	Connect-ExchangeOnline -CertificateFilePath $EXCHCERTIFICATEPATH  -AppID $EXCHAPPID -Organization $EXCHORGANIZATION -ShowBanner:$false #   -BypassMailboxAnchoring:$true

}							
				`)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Connection script created at ", scriptPath)
		case "sharepoint":
			log.Println("Sharepoint")
			scriptPath, err := MakeConnectionScript(connectionType, `

$PNPAPPID=$env:PNPAPPID
$PNPTENANTID=$env:PNPTENANTID
$PNPCERTIFICATEPATH = "$($PSScriptRoot)/pnp.pfx"
$PNPSITE=$env:PNPSITE
$bytes = [Convert]::FromBase64String($ENV:PNPCERTIFICATE)
[IO.File]::WriteAllBytes($PNPCERTIFICATEPATH, $bytes)

write-output "Connecting to $PNPSITE"
Connect-PnPOnline -Url $PNPSITE  -ClientId $PNPAPPID -Tenant $PNPTENANTID -CertificatePath "$PNPCERTIFICATEPATH"
			
							
				`)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Connection script created at ", scriptPath)

		default:

			log.Fatalln("Unknown ", connectionType)
			return
		}
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
