package exec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/mrflynn/meridian/geolocation"
)

// ParseCommandString takes a command string with Golang template directives and
// returns a space separated slice command with all fields filled out.
func ParseCommandString(cmdString string, info *geolocation.Info) ([]string, error) {
	tmpl, err := template.New("exec").Parse(cmdString)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, info)
	if err != nil {
		return nil, err
	}

	return strings.Split(buf.String(), " "), nil
}

// Run executes the given command and outputs the result to stdout and stderr, respectively.
func Run(args []string) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
