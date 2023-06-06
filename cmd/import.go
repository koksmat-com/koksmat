/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/koksmat-com/koksmat/db"
	"github.com/koksmat-com/koksmat/io"
	model "github.com/koksmat-com/koksmat/model/exchange"
	"github.com/spf13/cobra"
)

var inputFile string
var domain string
var subject string

func readAndSave[K any]() {
	data := io.Readfile[K](inputFile)
	db.Save[K](domain, subject, data)
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Add a JSON file to the import queue",
	Long:  `Add a JSON file to the import queue for further processing `,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Importing")

		switch subject {
		case "recipients":
			readAndSave[model.RecipientType]()
		case "rooms":
			readAndSave[model.RoomType]()
		default:

			log.Fatalln("Unknown subject", subject)
			return
		}

		log.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&inputFile, "inputFile", "i", "", "Input file (required)")
	importCmd.MarkFlagRequired("inputFile")
	importCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain (required)")
	importCmd.MarkFlagRequired("domain")
	importCmd.Flags().StringVarP(&subject, "subject", "s", "", "Subject (required)")
	importCmd.MarkFlagRequired("subject")
}