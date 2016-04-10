package engine

import (
	"github.com/TIBCOSoftware/flogo/fg"
	"strings"
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



// ContainsItemPath determines if the path exists in  list of ItemConfigs
func ContainsItemPath(list []*ItemConfig, path string) bool {
	for _, v := range list {
		if v.Path == path {
			return true
		}
	}
	return false
}

// ContainsItemPath determines if the path exists in  list of ItemConfigs
func ContainsItemName(list []*ItemConfig, name string) bool {
	for _, v := range list {
		if v.Name == name {
			return true
		}
	}
	return false
}

// GetItemConfig gets the item config for the specified path or name
func GetItemConfig(list []*ItemConfig, itemNameOrPath string) (int, *ItemConfig) {

	isPath := strings.Contains(itemNameOrPath, "/")

	for i, v := range list {
		if (isPath && v.Path == itemNameOrPath) ||  (!isPath && v.Name == itemNameOrPath){
			return i, v
		}
	}
	return -1, nil
}


