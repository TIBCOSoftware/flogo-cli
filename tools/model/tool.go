package model

import (
	"github.com/TIBCOSoftware/flogo/cli"
)

var optModel = &cli.OptionInfo{
	IsTool:    true,
	Name:      "model",
	UsageLine: "model <command>",
	Short:     "tool to manage a model project",
	Long:      "Tool for managing a model project.",
}

var toolModel *cli.Tool

// Tool gets or creates the model tool
func Tool() *cli.Tool {
	if toolModel == nil {
		toolModel = cli.NewTool(optModel)
		cli.RegisterTool(toolModel)
	}

	return toolModel
}

func init() {
	Tool()
}
