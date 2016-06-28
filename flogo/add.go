package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optAdd = &cli.OptionInfo{
	Name:      "add",
	UsageLine: "add <activity|model|trigger|flow> <path>",
	Short:     "add an activity, model or trigger to a flogo project",
	Long: `Add an activity, model, trigger or flow to a flogo project

Options:
    -src   copy contents to source (only when using file url)
`,
}

var validItemTypes = []string{itActivity, itTrigger, itModel, itFlow}

func init() {
	commandRegistry.RegisterCommand(&cmdAdd{option: optAdd})
}

type cmdAdd struct {
	option   *cli.OptionInfo
	addToSrc bool
}

func (c *cmdAdd) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdAdd) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.addToSrc), "src", false, "copy contents to source (only when using local/file)")
}

func (c *cmdAdd) Exec(args []string) error {

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
		fmt.Fprintf(os.Stderr, "Error: %s path not specified\n\n", fgutil.Capitalize(itemType))
		cmdUsage(c)
	}

	itemPath := args[1]

	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	installItem(projectDescriptor, itemType, itemPath, c.addToSrc)

	return nil
}

func installItem(projectDescriptor *FlogoProjectDescriptor, itemType string, itemPath string, addToSrc bool) {

	gb := fgutil.NewGb(projectDescriptor.Name)

	var itemConfigPath string
	var itemConfig *ItemDescriptor
	updateFiles := true

	switch itemType {
	case itActivity:
		itemConfig, itemConfigPath = AddFlogoItem(gb, itActivity, itemPath, projectDescriptor.Activities, addToSrc)
		projectDescriptor.Activities = append(projectDescriptor.Activities, itemConfig)
	case itModel:
		itemConfig, itemConfigPath = AddFlogoItem(gb, itModel, itemPath, projectDescriptor.Models, addToSrc)
		projectDescriptor.Models = append(projectDescriptor.Models, itemConfig)
	case itTrigger:
		itemConfig, itemConfigPath = AddFlogoItem(gb, itTrigger, itemPath, projectDescriptor.Triggers, addToSrc)
		projectDescriptor.Triggers = append(projectDescriptor.Triggers, itemConfig)

		//read trigger.json
		triggerConfigFile, err := os.Open(itemConfigPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to find '%s'\n\n", itemConfigPath)
			os.Exit(2)
		}

		triggerProjectDescriptor := &TriggerProjectDescriptor{}
		jsonParser := json.NewDecoder(triggerConfigFile)

		if err = jsonParser.Decode(triggerProjectDescriptor); err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to parse trigger.json, file may be corrupted.\n\n")
			os.Exit(2)
		}

		triggerConfigFile.Close()

		//read triggers.json
		triggersConfigPath := gb.NewBinFilePath(fileTriggersConfig)
		triggersConfigFile, err := os.Open(triggersConfigPath)

		triggersConfig := &TriggersConfig{}
		jsonParser = json.NewDecoder(triggersConfigFile)

		if err = jsonParser.Decode(triggersConfig); err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to parse application triggers.json, file may be corrupted.\n\n")
			os.Exit(2)
		}

		triggersConfigFile.Close()

		if triggersConfig.Triggers == nil {
			triggersConfig.Triggers = make([]*TriggerConfig, 0)
		}

		if !ContainsTriggerConfig(triggersConfig.Triggers, itemConfig.Name) {

			triggerConfig := &TriggerConfig{Name: itemConfig.Name, Settings: make(map[string]string)}

			for _, v := range triggerProjectDescriptor.Settings {

				triggerConfig.Settings[v.Name] = v.Value
			}

			triggersConfig.Triggers = append(triggersConfig.Triggers, triggerConfig)

			fgutil.WriteJSONtoFile(triggersConfigPath, triggersConfig)
		}

	case itFlow:
		updateFiles = false

		pathInfo, err := fgutil.GetPathInfo(itemPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid path '%s'\n", itemPath)
			os.Exit(2)
		}

		if pathInfo.IsLocal {

			if !pathInfo.IsFile {
				fmt.Fprintf(os.Stderr, "Error: Path '%s' is not a file\n", itemPath)
				os.Exit(2)
			}
		}

		ValidateFlow(projectDescriptor, itemPath, pathInfo.IsURL)

		if (pathInfo.IsLocal) {
			fgutil.CopyFile(pathInfo.FilePath, path("flows", pathInfo.FileName))
		} else if (pathInfo.IsURL) {
			fgutil.CopyRemoteFile(pathInfo.FileURL.String(), path("flows", pathInfo.FileName))
		}

		flows := ImportFlows(projectDescriptor, dirFlows)
		createFlowsGoFile(gb.CodeSourcePath, flows)

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
		os.Exit(2)
	}

	if updateFiles {
		updateProjectFiles(gb, projectDescriptor)
	}
}

// ContainsTriggerConfig determines if the list of TriggerConfigs contains the specified one
func ContainsTriggerConfig(list []*TriggerConfig, triggerName string) bool {
	for _, v := range list {
		if v.Name == triggerName {
			return true
		}
	}
	return false
}
