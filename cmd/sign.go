package cmd

import (
	"github.com/spf13/cobra"
)

var signRootCmd = &cobra.Command{
	Use:   "sign ",
	Short: "Digital signing",
	Args:  cobra.MinimumNArgs(0),
	Long:  ``,
}

func signSubCmd(use string, short string, long string, minargs int, run func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(minargs),
		Run:   run,
	}

	return cmd

}

/*

 Sign and validate files
 https://betterprogramming.pub/exploring-cryptography-in-go-signing-vs-encryption-f19534334ad

*/

func sign(cmd *cobra.Command, args []string) {

}

func validate(cmd *cobra.Command, args []string) {

}

func init() {

	rootCmd.AddCommand(signRootCmd)
	signRootCmd.AddCommand(signSubCmd("sign [file]", "Sign file", "", 1, sign))
	signRootCmd.AddCommand(signSubCmd("validate [file]", "Validate signature file", "", 1, validate))

}
