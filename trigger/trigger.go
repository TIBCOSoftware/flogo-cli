package trigger

import (
	"github.com/TIBCOSoftware/flogo/fg"
)

var optTrigger = &flogo.OptionInfo{
	IsTool:    true,
	Name:      "trigger",
	UsageLine: "trigger <command>",
	Short:     "tool to manage a trigger project",
	Long:      "Tool for managing a trigger project.",
}

var toolTrigger *flogo.Tool

// Tool gets or creates the tir
func Tool() *flogo.Tool {
	if toolTrigger == nil {
		toolTrigger = flogo.NewTool(optTrigger)
		flogo.RegisterTool(toolTrigger)
	}

	return toolTrigger
}

func init() {
	Tool()
}
