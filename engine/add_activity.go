package engine

import (
	"fg"
	"flag"
	"fgutil"
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
		return cfg.Models
	}

	itemConfig, engineConfig := AddEngineItem(c, "activity", args, gi, c.useSrc)

	engineConfig.Models = append(engineConfig.Models, itemConfig)
	fgutil.WriteJsonToFile(fileDescriptor, engineConfig)

	return nil
}

