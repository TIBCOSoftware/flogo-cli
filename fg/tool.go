package flogo

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/TIBCOSoftware/flogo/fgutil"
)

// Tool is a
type Tool struct {
	commandsMu  sync.Mutex
	optionInfo  *OptionInfo
	registry    *CommandRegistry
	TplUsage    string
	TplCmdUsage string
	TplCmdHelp  string
}

// NewTool creates a new tool
func NewTool(optionInfo *OptionInfo) *Tool {
	return &Tool{
		optionInfo:  optionInfo,
		registry:    NewCommandRegistry(),
		TplUsage:    tplToolUsage,
		TplCmdUsage: tplCmdUsage,
		TplCmdHelp:  tplCmdHelp,
	}
}

// OptionInfo implements HasOptionInfo
func (t *Tool) OptionInfo() *OptionInfo {
	return t.optionInfo
}

// CommandRegistry gets the command registry for the tool
func (t *Tool) CommandRegistry() *CommandRegistry {
	return t.registry
}

// Usage prints the usage details of the tool and exits with error
func (t *Tool) Usage() {
	t.PrintUsage(os.Stderr)
	os.Exit(2)
}

// PrintUsage prints the usage details of the tool
func (t *Tool) PrintUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	data := struct {
		Name        string
		OptionInfos []*OptionInfo
	}{
		t.optionInfo.Name,
		t.registry.CommandOptionInfos(),
	}

	fgutil.RenderTemplate(bw, t.TplUsage, data)
	bw.Flush()
}

// CmdUsage prints the usage details of the specified Command and
// exits with error
func (t *Tool) CmdUsage(command Command) {
	t.PrintCmdUsage(os.Stderr, command)
	os.Exit(2)
}

// PrintCmdUsage prints the usage details of the specified Command
func (t *Tool) PrintCmdUsage(w io.Writer, command Command) {
	bw := bufio.NewWriter(w)

	data := struct {
		ToolName     string
		CmdUsageLine string
	}{
		t.optionInfo.Name,
		command.OptionInfo().UsageLine,
	}

	fgutil.RenderTemplate(bw, t.TplCmdUsage, data)
	bw.Flush()
}

// PrintCmdHelp prints the help details of the specified Command
func (t *Tool) PrintCmdHelp(command Command) {
	bw := bufio.NewWriter(os.Stdout)

	data := struct {
		ToolName     string
		CmdUsageLine string
		CmdLong      string
	}{
		t.optionInfo.Name,
		command.OptionInfo().UsageLine,
		command.OptionInfo().Long,
	}

	fgutil.RenderTemplate(bw, t.TplCmdHelp, data)
	bw.Flush()
}

var tplToolUsage = `Usage:

    flogo {{.Name}} <command> [arguments]

Commands:
{{range .OptionInfos}}
    {{.Name | printf "%-20s"}} {{.Short}}{{end}}

`
var tplCmdUsage = `Usage:

    flogo {{.ToolName}} {{.CmdUsageLine}}

`
var tplCmdHelp = `Usage:

    flogo {{.ToolName}} {{.CmdUsageLine}}

{{.CmdLong | trim}}

`
