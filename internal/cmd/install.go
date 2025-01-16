package cmd

import (
	"qs-tools/internal/cmd/install"
)

func init() {
	RootCmd.AddCommand(install.Command())
}
