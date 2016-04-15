package model

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo/cli"
)

var optHelp = &cli.OptionInfo{
	Name:      "help",
	UsageLine: "help [command]",
	Short:     "get help for a model tool command",
	Long:      "Gets help for model tool commands.",
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdHelp{option: optHelp})
}

type cmdHelp struct {
	option *cli.OptionInfo
}

func (c *cmdHelp) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdHelp) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdHelp) Exec(args []string) error {
	if len(args) == 0 {
		Tool().PrintUsage(os.Stdout)
		return nil
	}
	if len(args) != 1 {
		Tool().CmdUsage(c)
	}

	arg := args[0]

	cmd, exists := toolModel.CommandRegistry().Command(arg)

	if exists {
		Tool().PrintCmdHelp(cmd)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Unknown flogo model command %#q. Run 'flogo model help'.\n", arg)
	os.Exit(2)

	return nil
}
