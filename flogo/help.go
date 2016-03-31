// Command, OptionInfo and command execution pattern derived from
// github.com/constabulary/gb, released under MIT license
// https://github.com/constabulary/gb/blob/master/LICENSE
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"github.com/TIBCOSoftware/flogo-tools/fgutil"
)

var optHelp = &flogo.OptionInfo{
	Name:      "help",
	UsageLine: "help [command]",
	Short:     "Get help for a command or tool",
	Long: `Get help for a flogo command or tool.

`,
}

func init() {
	commandRegistry.RegisterCommand(&cmdHelp{option: optHelp})
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
		printUsage(os.Stdout)
		return nil
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: flogo help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	arg := args[0]

	cmd, exists := commandRegistry.Command(arg)

	if exists {
		fgutil.RenderTemplate(os.Stdout, helpTpl, cmd.OptionInfo())
		return nil
	}

	tool, exists := flogo.GetTool(arg)

	if exists {
		fgutil.RenderTemplate(os.Stdout, "{{.Long}}\n\n", tool.OptionInfo())
		tool.PrintUsage(os.Stdout)
		return nil
	}

	fmt.Fprintf(os.Stderr, "Unknown help command %#q. Run 'flogo help'.\n", arg)
	os.Exit(2)

	return nil
}

var helpTpl = `usage: flogo {{.UsageLine}}

{{.Long | trim}}
`
