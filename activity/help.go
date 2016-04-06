package activity

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo/fg"
)

var optHelp = &flogo.OptionInfo{
	Name:      "help",
	UsageLine: "help [command]",
	Short:     "get help for an activity tool command",
	Long:      "Gets help for activity tool commands.",
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdHelp{option: optHelp})
}

type cmdHelp struct {
	option *flogo.OptionInfo
}

func (c *cmdHelp) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdHelp) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdHelp) Exec(ctx *flogo.Context, args []string) error {

	if len(args) == 0 {
		Tool().PrintUsage(os.Stdout)
		return nil
	}

	if len(args) != 1 {
		Tool().CmdUsage(c)
	}

	arg := args[0]

	cmd, exists := activityTool.CommandRegistry().Command(arg)

	if exists {
		Tool().PrintCmdHelp(cmd)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Unknown flogo activity command %#q. Run 'flogo activity help'.\n", arg)
	os.Exit(2)

	return nil
}
