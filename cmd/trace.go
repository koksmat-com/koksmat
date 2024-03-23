package cmd

import (
	"github.com/koksmat-com/koksmat/tracing"
	"github.com/spf13/cobra"
)

func init() {

	var tracecmd = &cobra.Command{
		Use:              "trace ",
		Short:            "Log handling",
		TraverseChildren: true,
		Long:             ``,
	}
	tracecmd.Example = `  koksmat trace log 'hello world'`

	tracecmd.AddCommand(

		&cobra.Command{
			Use:   "log [message]",
			Short: "Trace log ",

			Args: cobra.MinimumNArgs(1),
			Long: `
The trace log will ship data according to environment variables.

Currently the following environment variables are supported

NATS_SUBJECT - The subject to publish to send log message to NATS
			
			`,

			Run: func(cmd *cobra.Command, args []string) {
				tracing.Log(args)

			},
		})
	rootCmd.AddCommand(tracecmd)

}
