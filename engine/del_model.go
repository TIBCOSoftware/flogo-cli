package engine

import (
	"flag"
	"fmt"
	"os"

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
}

func (c *cmdDelModel) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDelModel) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdDelModel) Exec(ctx *flogo.Context, args []string) error {
	//if len(args) == 0 {
	//	Tool().PrintUsage(os.Stdout)
	//	return nil
	//}
	//if len(args) != 1 {
	//	fmt.Fprintf(os.Stderr, "usage: flogo activity add-model command\n\nToo many arguments given.\n")
	//	os.Exit(2)
	//}
	//

	arg := args[0]

	//
	//cmd, exists := activityTool.CommandRegistry().Command(arg)
	//
	//if exists {
	//	fgutil.RenderTemplate(os.Stdout, add-modelTpl, cmd.OptionInfo())
	//	return nil
	//}

	fmt.Fprintf(os.Stderr, "Unknown flogo engine add-model option %#q. Run 'flogo engine help add-model'.\n", arg)
	os.Exit(2)

	return nil
}
