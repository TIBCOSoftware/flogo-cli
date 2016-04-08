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

//////////////////////////////////////////////////////////////
// ProjectConfig

// EngineProjectConfig is engine project configuration object
type EngineProjectConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Models     []*ItemConfig `json:"models"`
	Activities []*ItemConfig `json:"activities"`
	Triggers   []*ItemConfig `json:"triggers"`
}

// ItemConfig is configuration for a model, activity or trigger
type ItemConfig struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Local   bool   `json:"local"`
}

type TriggerProjectConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Config     []*ConfigValue `json:"config"`
}

type ConfigValue struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

///////////////////////////////////////////////////////////////
// Engine Config

type EngineConfig struct {
	LogLevel        string           `json:"loglevel"`
	StateServiceURI string           `json:"state_service"`
	WorkerConfig    *WorkerConfig    `json:"engine"`
	Triggers        []*TriggerConfig `json:"triggers"`
}

type TriggerConfig struct {
	Name   string `json:"name"`
	Config map[string]string  `json:"config"`
}

type WorkerConfig struct {
	NumWorkers    int `json:"workers_count"`
	WorkQueueSize int `json:"workqueue_size"`
	MaxStepCount  int `json:"stepcount_max"`
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


