package engine

import (
	"fg"
	"flag"
	"fmt"
	"os"
)

var optList = &flogo.OptionInfo{
	Name:      "list",
	UsageLine: "list [object type]",
	Short:     "list objects configured on the engine project",
	Long: `List the object the engine project has been configured with.
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdList{option: optList})
}

type cmdList struct {
	option *flogo.OptionInfo
}

func (c *cmdList) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdList) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdList) Exec(ctx *flogo.Context, args []string) error {
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
