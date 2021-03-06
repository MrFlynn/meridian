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
		Args:  cobra.NoArgs,
		RunE:  printJSON,
	}
)

func init() {
	jsonCmd.Flags().StringSliceVarP(
		&fields,
		"fields",
		"f",
		[]string{"Latitude", "Longitude", "City", "RegionName", "Country", "IP"},
		"Location data to display. See info subcommand for options (required)\n",
	)

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
