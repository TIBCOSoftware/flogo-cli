package engine

import (
	"github.com/TIBCOSoftware/flogo-tools/fg"
)

var optEngine = &flogo.OptionInfo{
	IsTool:    true,
	Name:      "engine",
	UsageLine: "engine [command]",
	Short:     "tool to manage an engine project",
	Long:      "Tool for managing an engine project.",
}

var toolEngine *flogo.Tool

// Tool gets or creates the engine tool
func Tool() *flogo.Tool {
	if toolEngine == nil {
		toolEngine = flogo.NewTool(optEngine)
		flogo.RegisterTool(toolEngine)
	}

	return toolEngine
}

func init() {
	Tool()
}

// EngineConfig is engine project configuration object
type EngineConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Models     []*ItemConfig `json:"models"`
	Activities []*ItemConfig `json:"activities"`
	Triggers   []*ItemConfig `json:"triggers"`
}

// ItemConfig is configuration for a model, activity or trigger
type ItemConfig struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

// ContainsItem determines if the path exists in  list of ItemConfigs
func ContainsItem(path string, list []*ItemConfig) bool {
	for _, v := range list {
		if v.Path == path {
			return true
		}
	}
	return false
}
