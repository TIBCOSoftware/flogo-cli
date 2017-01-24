package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"strings"
)

var optDel = &cli.OptionInfo{
	Name:      "del",
	UsageLine: "del <activity|model|trigger|flow> <name>",
	Short:     "remove an activity, model, trigger or flow from a flogo project",
	Long: `Remove an activity, model, trigger or flow from a flogo project
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

	projectDescriptor := loadProjectDescriptor()

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "Error: item type not specified\n\n")
		cmdUsage(c)
	}

	itemType := strings.ToLower(args[0])

	if !fgutil.IsStringInList(itemType, validItemTypes) {
		fmt.Fprintf(os.Stderr, "Error: invalid item type '%s'\n\n", itemType)
		cmdUsage(c)
	}

	if len(args) == 1 {
		if itemType == itFlow {
			fmt.Fprintf(os.Stderr, "Error: Flow name or file not specified\n\n")
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s name or path not specified\n\n", fgutil.Capitalize(itemType))
		}
		cmdUsage(c)
	}

	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	gb := fgutil.NewGb(projectDescriptor.Name)

	itemNameOrPath := args[1]

	switch itemType {
	case itActivity:
		projectDescriptor.Activities = DelFlogoItem(gb, itActivity, itemNameOrPath, projectDescriptor.Activities, c.useSrc)

	case itModel:
		projectDescriptor.Models = DelFlogoItem(gb, itModel, itemNameOrPath, projectDescriptor.Models, c.useSrc)

	case itTrigger:
		projectDescriptor.Triggers = DelFlogoItem(gb, itTrigger, itemNameOrPath, projectDescriptor.Triggers, c.useSrc)

	case itFlow:

		if (strings.HasPrefix(itemNameOrPath, "embedded://")) {

			deleted := fgutil.DeleteFilesWithPrefix(dirFlows, itemNameOrPath[11:])

			if deleted == 0 {
				fmt.Fprintf(os.Stderr, "Error: Flow '%s' not found\n", itemNameOrPath)
				os.Exit(2)
			}

		} else {

			filePath := path( dirFlows, itemNameOrPath)

			fileInfo, err := os.Stat(filePath)

			if err !=  nil {
				fmt.Fprintf(os.Stderr, "Error: Flow '%s' not found\n", filePath)
				os.Exit(2)
			}

			if fileInfo.IsDir() {
				fmt.Fprintf(os.Stderr, "Error: Flow '%s' not found\n", filePath)
				os.Exit(2)
			}

			os.Remove(filePath)
		}

		flows := ImportFlows(projectDescriptor, dirFlows)
		createFlowsGoFile(gb.CodeSourcePath, flows)

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
		os.Exit(2)
	}

	updateProjectFiles(gb, projectDescriptor)

	return nil
}
