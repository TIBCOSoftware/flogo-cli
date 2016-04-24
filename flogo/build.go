package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo/cli"
	"github.com/TIBCOSoftware/flogo/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build [-validate]",
	Short:     "build the flogo application",
	Long: `Build the flogo application.

Options:
    -validate   validate that the project is buildable
`,
}

const fileFlowsGo string = "flows.go"

func init() {
	commandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option   *cli.OptionInfo
	validate bool
}

func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.validate), "validate", false, "only validate if buildable")
}

func (c *cmdBuild) Exec(args []string) error {

	projectDescriptor := loadProjectDescriptor()

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	gb := fgutil.NewGb(projectDescriptor.Name)

	flows := ImportFlows(projectDescriptor, dirFlows)
	createFlowsGoFile(gb.CodeSourcePath, flows)

	if len(projectDescriptor.Models) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one model.\n\n")
		os.Exit(2)
	}

	if len(projectDescriptor.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one trigger.\n\n")
		os.Exit(2)
	}

	if c.validate {
		return nil
	}

	err := gb.Build()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

