package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonCmd = &cobra.Command{
		Use:   "json",
		Short: "Outputs location information as JSON",
		Long:  "Outputs location query results as JSON for use with tools like jq",
		RunE:  printJSON,
	}
)

func init() {
	rootCmd.AddCommand(jsonCmd)
}

func printJSON(cmd *cobra.Command, args []string) error {
	jsonOutput, err := recievedInfo.ToJSON(fields...)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonOutput))
	return nil
}
