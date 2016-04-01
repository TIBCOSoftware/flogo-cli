package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optAddActivity = &flogo.OptionInfo{
	Name:      "add-activity",
	UsageLine: "add-activity <activity name>",
	Short:     "adds an activity to an engine project",
	Long: `Adds an activity to an engine project
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdAddActivity{option: optAddActivity})
}

type cmdAddActivity struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdAddActivity) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdAddActivity) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.useSrc), "src", false, "copy contents to source (only when using local/file)")
}

func (c *cmdAddActivity) Exec(ctx *flogo.Context, args []string) error {

	gi := func(cfg *EngineConfig) []*ItemConfig {
		return cfg.Activities
	}

	itemConfig, engineConfig := AddEngineItem(c, "activity", args, gi, c.useSrc)
	engineConfig.Activities = append(engineConfig.Activities, itemConfig)

	updateConfigFiles(engineConfig)

	return nil
}
