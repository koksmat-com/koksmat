package cmd

import (
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync ",
	Short: "Syncronization handling",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,
}

func init() {

	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(KitchenCompareCmd())
	syncCmd.AddCommand(KitchenUpdateCmd())

}
