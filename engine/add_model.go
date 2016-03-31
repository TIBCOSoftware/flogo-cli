package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"github.com/TIBCOSoftware/flogo-tools/fgutil"
)

var optAddModel = &flogo.OptionInfo{
	Name:      "add-model",
	UsageLine: "add-model <path>",
	Short:     "adds a model to an engine project",
	Long: `Adds a model to an engine project
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdAddModel{option: optAddModel})
}

type cmdAddModel struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdAddModel) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdAddModel) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.useSrc), "src", false, "copy contents to source (only when using local/file)")
}

func (c *cmdAddModel) Exec(ctx *flogo.Context, args []string) error {

	gi := func(cfg *EngineConfig) []*ItemConfig {
		return cfg.Models
	}

	itemConfig, engineConfig := AddEngineItem(c, "model", args, gi, c.useSrc)

	engineConfig.Models = append(engineConfig.Models, itemConfig)
	fgutil.WriteJSONtoFile(fileDescriptor, engineConfig)

	return nil
}
