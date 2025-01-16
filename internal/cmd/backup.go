package cmd

import (
	"qs-tools/internal/cmd/backup"
)

func init() {
	RootCmd.AddCommand(backup.Command())
}
