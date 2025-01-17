package cmd

import (
	"qs-tools/internal/cmd/apply"
)

func init() {
	RootCmd.AddCommand(apply.Command())
}
