package app

import (
	"flag"

	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build [-o][-e][-sp][-shim][-docker]",
	Short:     "build the flogo application",
	Long: `Build the flogo application.

Options:
    -o      optimize for directly referenced contributions
    -e      embed application configuration into executable
    -sp     skip prepare
    -shim   trigger shim
    -docker create a docker image based on Alpine Linux
`,
}

func init() {
	CommandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option      *cli.OptionInfo
	optimize    bool
	skipPrepare bool
	embedConfig bool
	shim        string
	docker      string
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

// AddFlags implementation of cli.Command.AddFlags
func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.optimize), "o", false, "optimize build")
	fs.BoolVar(&(c.embedConfig), "e", false, "embed config")
	fs.BoolVar(&(c.skipPrepare), "sp", false, "skip prepare")
	fs.StringVar(&(c.shim), "shim", "", "trigger shim")
	fs.StringVar(&(c.docker), "docker", "", "build docker")
}

// Exec implementation of cli.Command.Exec
func (c *cmdBuild) Exec(args []string) error {

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	options := &BuildOptions{SkipPrepare: c.skipPrepare, BuildDocker: c.docker, PrepareOptions: &PrepareOptions{OptimizeImports: c.optimize, EmbedConfig: c.embedConfig, Shim: c.shim}}
	return BuildApp(SetupExistingProjectEnv(appDir), options)
}
