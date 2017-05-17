package device

import (
	"github.com/TIBCOSoftware/flogo-cli/util"
	"os"
	"fmt"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"bufio"
	"io"
)

var (
	CommandRegistry = cli.NewCommandRegistry()
)

func SetupNewProjectEnv() Project {
	return NewPlatformIoProject()
}

func SetupExistingProjectEnv(appDir string) Project {

	project := NewPlatformIoProject()

	if err := project.Init(appDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing flogo device project: %s\n\n", err.Error())
		os.Exit(2)
	}

	if err := project.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "Error opening flogo device project: %s\n\n", err.Error())
		os.Exit(2)
	}

	return project
}

func Usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func cmdUsage(command cli.Command) {
	cli.CmdUsage("", command)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	options := CommandRegistry.CommandOptionInfos()
	options = append(options, cli.GetToolOptionInfos()...)

	fgutil.RenderTemplate(bw, usageTpl, options)
	bw.Flush()
}

var usageTpl = `Usage:

    flogodevice <command> [arguments]

Commands:
{{range .}}
    {{.Name | printf "%-12s"}} {{.Short}}{{end}}
`
