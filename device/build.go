package device

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"os"
	"fmt"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build [-o][-i][-bp]",
	Short:     "build the flogo device application",
	Long: `Build the flogo device application.

Options:
    -sp  skip prepare
`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option      *cli.OptionInfo
	skipPrepare bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.skipPrepare), "sp", false, "skip prepare")
}

// Exec implementation of cli.Command.Exec
func (c *cmdBuild) Exec(args []string) error {

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	options := &BuildOptions{SkipPrepare:c.skipPrepare, PrepareOptions:&PrepareOptions{}}
	return BuildDevice(SetupExistingProjectEnv(appDir), options)
}
