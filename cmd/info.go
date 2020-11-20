package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Display program information",
		Long:  `Display program information as well as a list of valid field values.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("meridian v%s\n", viper.GetString("version"))
			fmt.Printf("built: %s, commit: %s\n", viper.GetTime("date"), viper.GetString("commit"))
			fmt.Printf("\n%s\n", availableOptions)
		},
	}

	availableOptions = `Available data fields include:
- Continent: Full name of continent, ex. North America.
- ContinentCode: Shorthand name of continent, ex. NA.
- Country: Full name of country, ex. United States.
- CountryCode: Shorthand name of country, ex. US.
- Region: Shorthand name of region, state, etc., ex. CA.
- RegionName: Full name of region, state, etc., ex. California.
- City: Full name of city, ex. San Francisco.
- District: Full name of city district, ex. South of Market.
- ZIP: Postal code, ex. 94103.
- Latitude.
- Longitude.
- Timezone: tzdata name of timezone, ex. America/Los_Angeles.
- TimezoneOffset: Offset in seconds from UTC, ex. -28800 for America/Los_Angeles.
- ISP: Name of ISP, ex. Comcast Cable Communications, LLC.
- ORG: Organizational owner of IP, usually ISP ex. Comcast Cable Communications, Inc.
- ASN: Number and name of AS for current IP, ex. AS7922 Comcast Cable Communications, LLC.
- Mobile: Whether or not you are on a mobile network.
- Proxy: Whether or not you are using a proxy.
- IP: Current IP address, ex. 0.0.0.0.`
)

func init() {
	rootCmd.AddCommand(infoCmd)
}
