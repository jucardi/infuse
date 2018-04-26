package cli

import (
	"github.com/jucardi/go-infuse/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// FromCommand sets values to the global configuration by obtaining values from the command flags.
func FromCommand(cmd *cobra.Command) {
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug level enabled")
		config.Get().Verbose = true
	}
}
