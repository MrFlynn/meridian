package cmd

import (
	"fmt"
	"os"

	"github.com/mrflynn/meridian/geolocation"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "meridian",
		Short: "An application for displaying information about yourlocation",
		Long: `Meridian is a CLI application for displaying information about your location like
latitude, longitude, timezone, country, etc.`,
		PersistentPreRunE: setup,
		Args:              cobra.NoArgs,
		Run:               printDefault,
	}

	fields       []string
	location     string
	recievedInfo *geolocation.Info
)

func init() {
	rootCmd.Flags().StringSliceVarP(
		&fields,
		"fields",
		"f",
		[]string{"Latitude", "Longitude", "City", "RegionName", "Country", "IP"},
		"Location data to display. See info subcommand for options (required)\n",
	)
	rootCmd.PersistentFlags().StringVarP(
		&location, "ip", "p", "", "IP address to use in query. Defaults to current location",
	)
}

func setup(cmd *cobra.Command, args []string) error {
	recievedInfo = &geolocation.Info{}

	err := recievedInfo.ValidateFields(fields...)
	if err != nil {
		return err
	}

	return recievedInfo.New(location)
}

func printDefault(cmd *cobra.Command, args []string) {
	fmt.Println(recievedInfo.ToString(fields...))
}

// Execute runs the main root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
