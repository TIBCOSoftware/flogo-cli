package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"fmt"
	"os"
)

var optDel = &flogo.OptionInfo{
	Name:      "del",
	UsageLine: "del <activity|model|trigger> <name>",
	Short:     "deletes an activity, model, or trigger from an engine project",
	Long: `Deletes an activity, model or trigger from an engine project
`,
}


func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdDel{option: optDel})
}

type cmdDel struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdDel) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDel) AddFlags(fs *flag.FlagSet) {
}

func (c *cmdDel) Exec(ctx *flogo.Context, args []string) error {

	itemType := args[0]

	var engineConfig *EngineProjectConfig
	var toRemove int

	switch itemType {
	case "activity":
		gi := func(cfg *EngineProjectConfig) []*ItemConfig {
			return cfg.Activities
		}
		toRemove, engineConfig = DelEngineItem(c, "activity", args[1:], gi, c.useSrc)
		if toRemove > -1 {
			engineConfig.Activities = append(engineConfig.Activities[:toRemove], engineConfig.Activities[toRemove +1:]...)
		}
	case "model":
		gi := func(cfg *EngineProjectConfig) []*ItemConfig {
			return cfg.Models
		}
		toRemove, engineConfig = DelEngineItem(c, "model", args[1:], gi, c.useSrc)
		if toRemove > -1 {
			engineConfig.Models = append(engineConfig.Models[:toRemove], engineConfig.Models[toRemove + 1:]...)
		}
	case "trigger":
		gi := func(cfg *EngineProjectConfig) []*ItemConfig {
			return cfg.Triggers
		}
		toRemove, engineConfig = DelEngineItem(c, "trigger", args[1:], gi, c.useSrc)
		if toRemove > -1 {
			engineConfig.Triggers = append(engineConfig.Triggers[:toRemove], engineConfig.Triggers[toRemove + 1:]...)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
		os.Exit(2)
	}

	updateProjectConfigFiles(engineConfig)

	return nil
}
