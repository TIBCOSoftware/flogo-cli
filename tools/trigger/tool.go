package trigger

import (
	"github.com/TIBCOSoftware/flogo/cli"
)

var optTrigger = &cli.OptionInfo{
	IsTool:    true,
	Name:      "trigger",
	UsageLine: "trigger <command>",
	Short:     "tool to manage a trigger project",
	Long:      "Tool for managing a trigger project.",
}

var toolTrigger *cli.Tool

// Tool gets or creates the tir
func Tool() *cli.Tool {
	if toolTrigger == nil {
		toolTrigger = cli.NewTool(optTrigger)
		cli.RegisterTool(toolTrigger)
	}

	return toolTrigger
}

func init() {
	Tool()
}
