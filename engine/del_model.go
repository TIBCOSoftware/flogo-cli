package engine

import (
	"flag"

	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optDelModel = &flogo.OptionInfo{
	Name:      "del-model",
	UsageLine: "del-model <model name>",
	Short:     "deletes a model from an engine project",
	Long: `Deletes a model from an engine project
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdDelModel{option: optDelModel})
}

type cmdDelModel struct {
	option *flogo.OptionInfo
	useSrc bool
}

func (c *cmdDelModel) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDelModel) DelFlags(fs *flag.FlagSet) {
}

func (c *cmdDelModel) Exec(ctx *flogo.Context, args []string) error {

	//gi := func(cfg *EngineConfig) []*ItemConfig {
	//	return cfg.Models
	//}

	//itemConfig, engineConfig := DelEngineItem(c, "model", args, gi, c.useSrc)
	//engineConfig.Models = append(engineConfig.Models, itemConfig)
	//
	//updateConfigFiles(engineConfig)

	return nil
}
