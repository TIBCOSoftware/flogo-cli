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
	UsageLine: "build [-o]",
	Short:     "build the flogo application",
	Long: `Build the flogo application.

Options:
    -o   optimize for embedded flows
`,
}

const fileFlowsGo string = "flows.go"

func init() {
	commandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option   *cli.OptionInfo
	optimize bool
}

func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.optimize), "o", false, "optimize build")
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

	if c.optimize {

		//todo optimize triggers

		activityTypes := getAllActivityTypes(dirFlows)

		var activities []*ItemDescriptor

		for  _, activity := range projectDescriptor.Activities {

			if _, ok := activityTypes[activity.Name]; ok {
				activities = append(activities, activity)
			}
		}

		projectDescriptor.Activities = activities

		createImportsGoFile(gb.CodeSourcePath, projectDescriptor)
	}

	err := gb.Build()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

