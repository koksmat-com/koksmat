package cmd

import (
	"log"

	"github.com/koksmat-com/koksmat/kitchen"
	"github.com/spf13/cobra"
)

var installRootCmd = &cobra.Command{
	Use:   "add ",
	Short: "Add ingredients to your project",
	Args:  cobra.MinimumNArgs(1),
	Long:  `This command will add ingredients to your project based on the id provided.`,
	Example: `
	koksmat add V1StGXR8_Z5jdHi6B-myT
	
	Adds a snippet which you have found using Koksmat Studio
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := kitchen.AddIngredients(args[0])
		if err != nil {
			log.Fatalln(err)

		}
		log.Println("Ingredients added successfully")

	},
}

func init() {

	rootCmd.AddCommand(installRootCmd)
	//installRootCmd.AddCommand(installSubCmd("snippet", "Install snippet", "i", 1, install))

}
