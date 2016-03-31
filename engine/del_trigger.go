package engine

import (
	"flag"
	"fmt"
	"os"

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
}

func (c *cmdDelTrigger) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDelTrigger) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdDelTrigger) Exec(ctx *flogo.Context, args []string) error {
	//if len(args) == 0 {
	//	Tool().PrintUsage(os.Stdout)
	//	return nil
	//}
	//if len(args) != 1 {
	//	fmt.Fprintf(os.Stderr, "usage: flogo activity add-trigger command\n\nToo many arguments given.\n")
	//	os.Exit(2)
	//}
	//

	arg := args[0]

	//
	//cmd, exists := activityTool.CommandRegistry().Command(arg)
	//
	//if exists {
	//	fgutil.RenderTemplate(os.Stdout, add-triggerTpl, cmd.OptionInfo())
	//	return nil
	//}

	fmt.Fprintf(os.Stderr, "Unknown flogo engine add-trigger option %#q. Run 'flogo engine help add-trigger'.\n", arg)
	os.Exit(2)

	return nil
}
