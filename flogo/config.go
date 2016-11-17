package main

//////////////////////////////////////////////////////////////
// ProjectDescriptor

// FlogoProjectDescriptor is the flogo project descriptor object
type FlogoProjectDescriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Models     []*ItemDescriptor `json:"models"`
	Activities []*ItemDescriptor `json:"activities"`
	Triggers   []*ItemDescriptor `json:"triggers"`
}

// FlogoPaletteDescriptor is the flogo palette descriptor object

type FlogoExtensions struct {
	Models     []*ItemDescriptor `json:"models"`
	Activities []*ItemDescriptor `json:"activities"`
	Triggers   []*ItemDescriptor `json:"triggers"`
}

type FlogoPaletteDescriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	FlogoExtensions *FlogoExtensions `json:"extensions"`
}

// ItemDescriptor is configuration for a model, activity or trigger
type ItemDescriptor struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Version   string `json:"version"`
	LocalPath string `json:"localpath,omitempty"`
}

func (d *ItemDescriptor) Local() bool {
	return len(d.LocalPath) > 0
}

// TriggerProjectDescriptor is the trigger project descriptor object
type TriggerProjectDescriptor struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Settings []*ConfigValue `json:"settings"`
	Outputs  []*ConfigValue `json:"outputs"`

	Endpoint *EndpointDescriptor `json:"endpoint"`
}

// EndpointDescriptor is the trigger endpoint descriptor object
type EndpointDescriptor struct {
	Settings []*ConfigValue `json:"settings"`
}

// ConfigValue struct describes a configuration value
type ConfigValue struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

///////////////////////////////////////////////////////////////
// Engine Config

// todo: consolidate with config from flogo-lib

// EngineConfig is the engine configuration object
type EngineConfig struct {
	LogLevel     string           `json:"loglevel"`
	RunnerConfig *RunnerConfig    `json:"actionRunner"`
	Triggers     []*TriggerConfig `json:"triggers,omitempty"`
	Services     []*ServiceConfig `json:"services"`
}

// TriggersConfig is the triggers configuration object
type TriggersConfig struct {
	Triggers []*TriggerConfig `json:"triggers"`
}

// RunnerConfig is the runner configuration object
type RunnerConfig struct {
	Type   string        `json:"type"`
	Pooled *PooledConfig `json:"pooled,omitempty"`
	Direct *DirectConfig `json:"direct,omitempty"`
}

// DirectConfig  is the configuration object for a Direct Runner
type DirectConfig struct {
}

// PooledConfig  is the configuration object for a Pooled Runner
type PooledConfig struct {
	NumWorkers    int `json:"numWorkers"`
	WorkQueueSize int `json:"workQueueSize"`
}

// TriggerConfig is the trigger configuration object
type TriggerConfig struct {
	Name      string            `json:"name"`
	Settings  map[string]string `json:"settings"`
	Endpoints []*EndpointConfig `json:"endpoints"`
}

// EndpointConfig is the endpoint configuration object
type EndpointConfig struct {
	ActionType string            `json:"actionType"`
	ActionURI  string            `json:"actionURI"`
	Settings   map[string]string `json:"settings"`
}

// ServiceConfig is the service configuration object
type ServiceConfig struct {
	Name     string            `json:"name"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings,omitempty"`
}

// DefaultEngineConfig returns the default engine configuration
func DefaultEngineConfig() *EngineConfig {

	var ec EngineConfig

	ec.LogLevel = "INFO"
	ec.RunnerConfig = &RunnerConfig{Type: "pooled", Pooled: &PooledConfig{NumWorkers: 5, WorkQueueSize: 50}}
	ec.Services = make([]*ServiceConfig, 0)

	ec.Services = append(ec.Services, &ServiceConfig{Name: "stateRecorder", Enabled: false, Settings: map[string]string{"host": "", "port": ""}})
	ec.Services = append(ec.Services, &ServiceConfig{Name: "flowProvider", Enabled: true})
	ec.Services = append(ec.Services, &ServiceConfig{Name: "engineTester", Enabled: true, Settings: map[string]string{"port": "8080"}})

	return &ec
}

// DefaultTriggersConfig returns the default triggers configuration
func DefaultTriggersConfig() *TriggersConfig {

	var tc TriggersConfig
	tc.Triggers = make([]*TriggerConfig, 0)

	return &tc
}
