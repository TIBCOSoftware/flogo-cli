package activity

import (
	"github.com/TIBCOSoftware/flogo-cli/cli"
)

var optActivity = &cli.OptionInfo{
	IsTool:    true,
	Name:      "activity",
	UsageLine: "activity <command>",
	Short:     "tool to manage an activity project",
	Long:      "Tool for managing an activity project.",
}

var activityTool *cli.Tool

// Tool gets or create the activity tool
func Tool() *cli.Tool {
	if activityTool == nil {
		activityTool = cli.NewTool(optActivity)
		cli.RegisterTool(activityTool)
	}

	return activityTool
}

func init() {
	Tool()
}
