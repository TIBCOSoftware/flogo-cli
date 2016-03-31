package model

import (
	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optModel = &flogo.OptionInfo{
	IsTool:    true,
	Name:      "model",
	UsageLine: "model <command>",
	Short:     "tool to manage a model project",
	Long:      "Tool for managing a model project.",
}

var toolModel *flogo.Tool

// Tool gets or creates the model tool
func Tool() *flogo.Tool {
	if toolModel == nil {
		toolModel = flogo.NewTool(optModel)
		flogo.RegisterTool(toolModel)
	}

	return toolModel
}

func init() {
	Tool()
}
