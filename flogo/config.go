package main

//////////////////////////////////////////////////////////////
// ProjectConfig

// FlogoProjectConfig is the flogo project configuration object
type FlogoProjectConfig struct {
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

// TriggerProjectConfig is the trigger project configuration object
type TriggerProjectConfig struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`

	Config []*ConfigValue `json:"config"`
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
	RunnerConfig *RunnerConfig    `json:"flowRunner"`
	Triggers     []*TriggerConfig `json:"triggers"`
	Services     []*ServiceConfig `json:"services"`
}

// RunnerConfig is the runner configuration object
type RunnerConfig struct {
	Type   string        `json:"type"`
	Pooled *PooledConfig `json:"pooled,omitempty"`
	Direct *DirectConfig `json:"direct,omitempty"`
}

// DirectConfig  is the configuration object for a Direct Runner
type DirectConfig struct {
	MaxStepCount int `json:"maxStepCount"`
}

// PooledConfig  is the configuration object for a Pooled Runner
type PooledConfig struct {
	NumWorkers    int `json:"numWorkers"`
	WorkQueueSize int `json:"workQueueSize"`
	MaxStepCount  int `json:"maxStepCount"`
}

// TriggerConfig is the trigger configuration object
type TriggerConfig struct {
	Name      string            `json:"name"`
	Settings  map[string]string `json:"settings"`
	Endpoints []*EndpointConfig `json:"endpoints"`
}

// EndpointConfig is the endpoint configuration object
type EndpointConfig struct {
	FlowURI  string            `json:"flowURI"`
	Settings map[string]string `json:"settings"`
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
	ec.RunnerConfig = &RunnerConfig{Type: "pooled", Pooled: &PooledConfig{NumWorkers: 5, WorkQueueSize: 50, MaxStepCount: 32000}}
	ec.Triggers = make([]*TriggerConfig, 0)
	ec.Services = make([]*ServiceConfig, 0)

	ec.Services = append(ec.Services, &ServiceConfig{Name: "stateRecorder", Enabled: false, Settings: map[string]string{"host": "", "port": ""}})
	ec.Services = append(ec.Services, &ServiceConfig{Name: "flowProvider", Enabled: true})
	ec.Services = append(ec.Services, &ServiceConfig{Name: "engineTester", Enabled: true, Settings: map[string]string{"port": "8080"}})

	return &ec
}
