package engine

import (
	"flag"
	"fmt"
	"os"

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
}

func (c *cmdDelActivity) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdDelActivity) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdDelActivity) Exec(ctx *flogo.Context, args []string) error {
	//if len(args) == 0 {
	//	Tool().PrintUsage(os.Stdout)
	//	return nil
	//}
	//if len(args) != 1 {
	//	fmt.Fprintf(os.Stderr, "usage: flogo activity add-activity command\n\nToo many arguments given.\n")
	//	os.Exit(2)
	//}
	//

	arg := args[0]

	//
	//cmd, exists := activityTool.CommandRegistry().Command(arg)
	//
	//if exists {
	//	fgutil.RenderTemplate(os.Stdout, add-activityTpl, cmd.OptionInfo())
	//	return nil
	//}

	fmt.Fprintf(os.Stderr, "Unknown flogo engine add-activity command %#q. Run 'flogo engine help add-activity'.\n", arg)
	os.Exit(2)

	return nil
}

//Error: Current working directory is not a flogo-based project.
