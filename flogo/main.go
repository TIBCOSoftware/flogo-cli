// Command, OptionInfo, help and command execution pattern derived from
// github.com/constabulary/gb, released under MIT license
// https://github.com/constabulary/gb/blob/master/LICENSE

package main

import (
	"bufio"
	"fg"
	"fgutil"
	"flag"
	"fmt"
	"io"
	"os"

	_ "activity"
	_ "engine"
	_ "model"
	_ "trigger"
)

var (
	commandRegistry = flogo.NewCommandRegistry()
	fs              = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
)

func init() {
	fs.Usage = usage
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) < 2 || args[1] == "-h" {
		usage()
	}
	name := args[1]

	var remainingArgs []string

	cmd, exists := commandRegistry.Command(name)

	if !exists {

		tool, toolExists := flogo.GetTool(name)

		if !toolExists {
			fmt.Fprintf(os.Stderr, "FATAL: unknown command or tool %q\n\n", name)
			usage()
		}

		if len(args) < 3 {
			tool.Usage()
		}

		cmd, exists = tool.CommandRegistry().Command(args[2])

		if !exists {
			fmt.Fprintf(os.Stderr, "FATAL: unknown command %q\n\n", args[2])
			tool.Usage()
		}

		name = name + ":" + tool.OptionInfo().Name
		remainingArgs = args[3:]
	} else {
		remainingArgs = args[2:]
	}

	if err := flogo.ExecCommand(nil, fs, cmd, remainingArgs); err != nil {
		fatalf("command %q failed: %v", name, err)
	}

	os.Exit(0)
}

var usageTpl = `Usage:

    flogo command/tool [arguments]

Commands:
{{range .}}{{if not .IsTool}}
    {{.Name | printf "%-12s"}} {{.Short}}{{end}}{{end}}

Tools:
{{range .}}{{if .IsTool}}
    {{.Name | printf "%-12s"}} {{.Short}}{{end}}{{end}}

`

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	options := commandRegistry.CommandOptionInfos()
	options = append(options, flogo.GetToolOptionInfos()...)

	fgutil.RenderTemplate(bw, usageTpl, options)
	bw.Flush()
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}
