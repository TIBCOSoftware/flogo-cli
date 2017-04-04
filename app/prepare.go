package app

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"os"
	"fmt"
)

var optPrepare = &cli.OptionInfo{
	Name:      "prepare",
	UsageLine: "prepare [-o][-i]",
	Short:     "prepare the flogo application",
	Long: `Prepare the flogo application.

Options:
    -o   optimize for embedded flows
    -i   incorporate config into application
`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdPrepare{option: optPrepare})
}

type cmdPrepare struct {
	option     *cli.OptionInfo
	optimize   bool
	includeCfg bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdPrepare) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdPrepare) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.optimize), "o", false, "optimize prepare")
	fs.BoolVar(&(c.includeCfg), "i", false, "include config")
}

// Exec implementation of cli.Command.Exec
func (c *cmdPrepare) Exec(args []string) error {

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	options := &PrepareOptions{OptimizeImports:c.optimize}
	return PrepareApp(SetupExistingProjectEnv(appDir), options)
}
