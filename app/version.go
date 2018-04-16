package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
)

var optVersion = &cli.OptionInfo{
	Name:      "version",
	UsageLine: "version",
	Short:     "Displays the version of Flogo CLI",
	Long: `Get the current version number of the CLI.

`,
}

var tag = ""
var hash = ""

func init() {
	CommandRegistry.RegisterCommand(&cmdVersion{option: optVersion})
}

type cmdVersion struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdVersion) OptionInfo() *cli.OptionInfo {
	return c.option
}

// Exec implementation of cli.Command.Exec
func (c *cmdVersion) AddFlags(fs *flag.FlagSet) {
	//op op
}

// Exec implementation of cli.Command.Exec
func (c *cmdVersion) Exec(args []string) error {

	line := fmt.Sprintf("Flogo CLI version [%s] and commithash [%s]\n\n", tag, hash)
	fmt.Fprint(os.Stdout, line)

	return nil
}
