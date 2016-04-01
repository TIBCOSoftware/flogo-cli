package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optDelTrigger = &flogo.OptionInfo{
	Name:      "del-trigger",
	UsageLine: "del-trigger <trigger name>",
	Short:     "deletes a trigger from an engine project",
	Long: `Deletes a trigger from an engine project.
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdDelTrigger{option: optDelTrigger})
}

type cmdDelTrigger struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdDelTrigger) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDelTrigger) DelFlags(fs *flag.FlagSet) {
}

func (c *cmdDelTrigger) Exec(ctx *flogo.Context, args []string) error {
	//gi := func(cfg *EngineConfig) []*ItemConfig {
	//	return cfg.Triggers
	//}

	//itemConfig, engineConfig := DelEngineItem(c, "trigger", args, gi, c.useSrc)
	//engineConfig.Triggers = append(engineConfig.Triggers, itemConfig)
	//
	//updateConfigFiles(engineConfig)

	return nil
}
