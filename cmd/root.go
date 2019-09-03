package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print more verbose information")

	convertCmd.PersistentFlags().String("base-currency", "EUR", "Currency the account is based on")

	RootCmd.AddCommand(convertCmd)
}

var RootCmd = &cobra.Command{
	Use: "",

	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
		debug, err := ccmd.Flags().GetBool("verbose")

		if err != nil {
			log.WithError(err).Error("failed to parse verbose flag")
			os.Exit(1)
			return
		}

		if debug {
			log.SetLevel(log.DebugLevel)
		}
	},

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
