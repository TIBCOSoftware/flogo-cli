package activity

import (
	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optActivity = &flogo.OptionInfo{
	IsTool:    true,
	Name:      "activity",
	UsageLine: "activity <command>",
	Short:     "tool to manage an activity project",
	Long:      "Tool for managing an activity project.",
}

var activityTool *flogo.Tool

// Tool gets or create the activity tool
func Tool() *flogo.Tool {
	if activityTool == nil {
		activityTool = flogo.NewTool(optActivity)
		flogo.RegisterTool(activityTool)
	}

	return activityTool
}

func init() {
	Tool()
}
