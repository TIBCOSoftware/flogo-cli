package main

import (
	"flag"

	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo/cli"
	"github.com/TIBCOSoftware/flogo/util"
)

var optDel = &cli.OptionInfo{
	Name:      "del",
	UsageLine: "del <activity|model|trigger> <name>",
	Short:     "remove an activity, model, or trigger from a flogo project",
	Long: `Remove an activity, model or trigger from a flogo project
`,
}

func init() {
	commandRegistry.RegisterCommand(&cmdDel{option: optDel})
}

type cmdDel struct {
	option *cli.OptionInfo
	useSrc bool
}

func (c *cmdDel) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdDel) AddFlags(fs *flag.FlagSet) {
}

func (c *cmdDel) Exec(args []string) error {

	projectConfig := loadProjectConfig()

	itemType := args[0]

	if len(args) == 1 {
		fmt.Fprintf(os.Stderr, "Error: %s name or path not specified\n\n", fgutil.Capitalize(itemType))
		cmdUsage(c)
	}

	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	gb := fgutil.NewGb(projectConfig.Name)

	itemNameOrPath := args[1]

	switch itemType {
	case itActivity:
		projectConfig.Activities = DelFlogoItem(gb, itActivity, itemNameOrPath, projectConfig.Activities, c.useSrc)

	case itModel:
		projectConfig.Models = DelFlogoItem(gb, itModel, itemNameOrPath, projectConfig.Models, c.useSrc)

	case itTrigger:
		projectConfig.Triggers = DelFlogoItem(gb, itTrigger, itemNameOrPath, projectConfig.Triggers, c.useSrc)

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
		os.Exit(2)
	}

	updateProjectConfigFiles(gb, projectConfig)

	return nil
}
