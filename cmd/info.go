package cmd

import (
	"fmt"

	"github.com/mrflynn/meridian/geolocation"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display program information",
	Long:  `Display program information as well as a list of valid field values.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("meridian v%s\n", viper.GetString("version"))
		fmt.Printf("built: %s, commit: %s\n", viper.GetTime("date"), viper.GetString("commit"))

		info := &geolocation.Info{}
		fmt.Printf("\n%s", info.ToDescription())
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
