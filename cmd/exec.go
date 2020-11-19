package cmd

import (
	"github.com/mrflynn/meridian/internal/exec"
	"github.com/spf13/cobra"
)

var (
	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "Execute another command that uses location information",
		Long: `Execute external commands using location information through predefined template
variables.`,
		Args: cobra.ExactArgs(1),
		RunE: executeCommand,
	}
)

func init() {
	rootCmd.AddCommand(execCmd)
}

func executeCommand(cmd *cobra.Command, args []string) error {
	cmdString, err := exec.ParseCommandString(args[0], recievedInfo)
	if err != nil {
		return err
	}

	return exec.Run(cmdString)
}
