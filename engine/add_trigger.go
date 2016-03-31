package engine

import (
	"fg"
	"flag"
	"fgutil"
)

var optAddTrigger = &flogo.OptionInfo{
	Name:      "add-trigger",
	UsageLine: "add-trigger <trigger name>",
	Short:     "adds a trigger to an engine project",
	Long: `Adds a trigger to an engine project.
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdAddTrigger{option: optAddTrigger})
}

type cmdAddTrigger struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdAddTrigger) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdAddTrigger) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.useSrc), "src", false, "copy contents to source (only when using local/file)")
}

func (c *cmdAddTrigger) Exec(ctx *flogo.Context, args []string) error {
	gi := func(cfg *EngineConfig) []*ItemConfig {
		return cfg.Models
	}

	itemConfig, engineConfig := AddEngineItem(c, "trigger", args, gi, c.useSrc)

	engineConfig.Models = append(engineConfig.Models, itemConfig)
	fgutil.WriteJsonToFile(fileDescriptor, engineConfig)

	return nil
}
