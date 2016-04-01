package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optDelActivity = &flogo.OptionInfo{
	Name:      "del-activity",
	UsageLine: "del-activity <activity name>",
	Short:     "deletes an activity from an engine project",
	Long: `Deletes an activity from an engine project
`,
}


func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdDelActivity{option: optDelActivity})
}

type cmdDelActivity struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdDelActivity) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDelActivity) DelFlags(fs *flag.FlagSet) {
}

func (c *cmdDelActivity) Exec(ctx *flogo.Context, args []string) error {

	//gi := func(cfg *EngineConfig) []*ItemConfig {
	//	return cfg.Activities
	//}

	//itemConfig, engineConfig := DelEngineItem(c, "activity", args, gi, c.useSrc)
	//engineConfig.Activities = append(engineConfig.Activities, itemConfig)

	//updateConfigFiles(engineConfig)

	return nil
}
