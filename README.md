# meridian
A simple CLI tool that gives you information about your physical location. It
can display things like your country, city, ISP, longitude/latitude, IP address
and more. For a full list of fields see the _Available Data_ section below or 
run `meridian info`.

## Installation
### Homebrew
The recommended installation method for meridian is using 
[Homebrew](https://brew.sh). Run the following commands to install.

```bash
$ brew tap mrflynn/cider
$ brew install mrflynn/cider/meridian
```

### Manually
Download the proper build for your OS + architecture combination from the 
releases page and run the following commands.

```bash
$ cd $HOME/Downloads # or where ever your downloads folder is.
$ tar -xzf meridian_<version>_<os>_<arch>.tar.gz
$ mv meridian /usr/bin
```

## Usage
To see the help menu just run `meridian --help`. You can also run
`meridian <section> --help` to get help for individual subcommands.

This tool can display the location data in three different modes which are
explained below.

### Default
The default mode is to diplay location information in a human readable format.
For example, running `meridian` will display your country, region (state,
territory, etc.), city, longitude, latitude, and IP address:

```bash
$ meridian
Country: United States
RegionName: California
City: San Francisco
Latitude: 37.7749
Longitude: -122.4194
IP: 0.0.0.0
```

You can also customize the output using the `--fields` or `--ip` flags.
* `--fields`: Configures which data fields are displayed. See the _Available
Data_ section for a list of valid fields.
* `--ip`: Allows the command to mimick as if it were in a different location.
By passing a IP address or domain name, the command will return location data 
for that IP or domain.

### JSON
This subcommand returns the location data as JSON. This is useful if you wish
to parse or augment the data using tools like [jq](https://stedolan.github.io/jq/).

Running this command with no flags returns your country, region (state,
territory, etc.), city, longitude, latitude, UTC offset (in seconds), whether 
or not your using a proxy and on a mobile connection, and IP address.

```bash
$ meridian json | jq
{
  "country": "United States",
  "regionName": "California",
  "city": "San Francisco",
  "lat": 37.7749,
  "lon": -122.4194,
  "offset": 0,
  "mobile": false,
  "proxy": false,
  "query": "0.0.0.0"
}
```

Both `--fields` and `--ip` work with this subcommand and behave the same
way.

### Exec
This subcommand allows you to execute external commands using location data.
This works by using [Go's template engine](https://golang.org/pkg/text/template/) 
to parse and execute commands. This is a very powerful subcommand that allows
you to use Go's very powerful template library to transform data and use it 
with external tools.

To use this subcommand, just surround the target command in quotes and use any
valid field or template directive you would like anywhere in the command. For 
example, you could use it to form a request to external API.

```bash
$ meridian exec "curl -sSL https://example.com/api/{{ .Latitude }}/{{ .Longitude }}"
```

Only the `--ip` flag works in this mode and it behaves the same as the default
command.

Pipes are not supported in this mode. However, you can use pipes after calling
this command as stdout and stderr returned.

## Available Data
The following location data fields are available through this tool:
* **Continent**: Full name of continent, _ex. North America_.
* **ContinentCode**: Shorthand name of continent, _ex. NA_.
* **Country**: Full name of country, _ex. United States_.
* **CountryCode**: Shorthand name of country, _ex. US_.
* **Region**: Shorthand name of region, state, etc., _ex. CA_.
* **RegionName**: Full name of region, state, etc., _ex. California_.
* **City**: Full name of city, _ex. San Francisco_.
* **District**: Full name of city district, _ex. South of Market_.
* **ZIP**: Postal code, _ex. 94103_.
* **Latitude**.
* **Longitude**.
* **Timezone**: tzdata name of timezone, _ex. America/Los\_Angeles_.
* **TimezoneOffset**: Offset in seconds from UTC, _ex. -28800 for 
America/Los\_Angeles_.
* **ISP**: Name of ISP, _ex. Comcast Cable Communications, LLC_.
* **ORG**: Organizational owner of IP, usually ISP _ex. Comcast Cable 
Communications, Inc_.
* **ASN**: Number and name of AS for current IP, _ex. AS7922 Comcast Cable 
Communications, LLC_.
* **Mobile**: Whether or not you are on a mobile network.
* **Proxy**: Whether or not you are using a proxy.
* **IP**: Current IP address, _ex. 0.0.0.0_.

If you wish to displays all fields in either the default or json modes, just
pass `--fields=All`.

## Contributing
If you have a new feature you would like to add or have a bugfix, please open a
pull request and I'll review it. Make sure you test everything you add as I need
to ensure that it works properly and how it works.

## Maintainers
* [Nick Pleatsikas](https://github.com/mrflynn)

## License
[MIT](LICENSE)