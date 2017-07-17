package device

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
)

var optInstall = &cli.OptionInfo{
	Name:      "install",
	UsageLine: "install [-v version] contribution",
	Short:     "install a flogo device contribution",
	Long: `Installs a flogo device contribution.

Options:
    -v specify the version
`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdInstall{option: optInstall})
}

type cmdInstall struct {
	option   *cli.OptionInfo
	version string
	palette bool
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdInstall) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdInstall) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&(c.version), "v", "", "version")

}

// Exec implementation of cli.Command.Exec
func (c *cmdInstall) Exec(args []string) error {

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "Error: contribution not specified\n\n")
		cmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprint(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	contribPath, version := splitVersion(args[0])

	if c.version != "" {
		version = c.version
	}

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	return InstallContribution(SetupExistingProjectEnv(appDir), contribPath, version)
}

