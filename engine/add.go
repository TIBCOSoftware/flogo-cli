package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo/fg"
	"os"
	"fmt"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo/fgutil"
)

var optAdd = &flogo.OptionInfo{
	Name:      "add",
	UsageLine: "add <activity|model|trigger> <path>",
	Short:     "adds an activity, model or trigger to an engine project",
	Long: `Adds an activity, model or trigger to an engine project

Options:
    -src   copy contents to source (only when using local/file)
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdAdd{option: optAdd})
}

type cmdAdd struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdAdd) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdAdd) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.useSrc), "src", false, "copy contents to source (only when using local/file)")
}

func (c *cmdAdd) Exec(ctx *flogo.Context, args []string) error {


	projectConfigFile, err := os.Open(fileProjectConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, "Error: %s item type not specified\n\n")
		Tool().CmdUsage(c)
	}

	itemType := args[0]

	if len(args) == 1 {
		fmt.Fprintf(os.Stderr, "Error: %s path not specified\n\n", fgutil.Capitalize(itemType))
		Tool().CmdUsage(c)
	}

	if len(args) > 2 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	projectConfig := &EngineProjectConfig{}
	jsonParser := json.NewDecoder(projectConfigFile)

	if err = jsonParser.Decode(projectConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse engine.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	projectConfigFile.Close()

	var itemConfigPath string
	var itemConfig *ItemConfig

	switch itemType {
	case "activity":
		gi := func(cfg *EngineProjectConfig) []*ItemConfig {
			return cfg.Activities
		}
		itemConfig, itemConfigPath = AddEngineItem(c, projectConfig, "activity", args[1:], gi, c.useSrc)
		projectConfig.Activities = append(projectConfig.Activities, itemConfig)
	case "model":
		gi := func(cfg *EngineProjectConfig) []*ItemConfig {
			return cfg.Models
		}
		itemConfig, itemConfigPath = AddEngineItem(c, projectConfig, "model", args[1:], gi, c.useSrc)
		projectConfig.Models = append(projectConfig.Models, itemConfig)
	case "trigger":
		gi := func(cfg *EngineProjectConfig) []*ItemConfig {
			return cfg.Triggers
		}
		itemConfig, itemConfigPath = AddEngineItem(c, projectConfig, "trigger", args[1:], gi, c.useSrc)
		projectConfig.Triggers = append(projectConfig.Triggers, itemConfig)

		//read trigger.json
		triggerConfigFile, err := os.Open(itemConfigPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to find '%s'\n\n", itemConfigPath)
			os.Exit(2)
		}

		triggerProjectConfig := &TriggerProjectConfig{}
		jsonParser = json.NewDecoder(triggerConfigFile)

		if err = jsonParser.Decode(triggerProjectConfig); err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to parse trigger.json, file may be corrupted.\n\n")
			os.Exit(2)
		}

		triggerConfigFile.Close()

		//read engine config.json
		engineConfigPath := path("bin",fileEngineConfig)
		engineConfigFile, err := os.Open(engineConfigPath)

		engineConfig := &EngineConfig{}
		jsonParser = json.NewDecoder(engineConfigFile)

		if err = jsonParser.Decode(engineConfig); err != nil {
			fmt.Fprint(os.Stderr, "Error: Unable to parse engine config.json, file may be corrupted.\n\n")
			os.Exit(2)
		}

		engineConfigFile.Close()

		if engineConfig.Triggers == nil {
			engineConfig.Triggers = make([]*TriggerConfig,0)
		}

		if !ContainsTriggerConfig(engineConfig.Triggers, itemConfig.Name) {

			triggerConfig := &TriggerConfig{Name:itemConfig.Name, Settings:make(map[string]string)}

			for _, v := range triggerProjectConfig.Config {

				triggerConfig.Settings[v.Name] = v.Value
			}

			engineConfig.Triggers = append(engineConfig.Triggers, triggerConfig)

			fgutil.WriteJSONtoFile(engineConfigPath, engineConfig)
		}

	// add config
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
		os.Exit(2)
	}

	updateProjectConfigFiles(projectConfig)

	return nil
}

func ContainsTriggerConfig(list []*TriggerConfig, triggerName string) bool {
	for _, v := range list {
		if v.Name == triggerName {
			return true
		}
	}
	return false
}