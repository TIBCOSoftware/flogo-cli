package engine

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

// todo: consolidate with config from flogo-lib

type EngineConfig struct {
	LogLevel     string           `json:"loglevel"`
	RunnerConfig *RunnerConfig    `json:"processRunner"`
	Triggers     []*TriggerConfig `json:"triggers"`
	Services     []*ServiceConfig `json:"services"`
}

type RunnerConfig struct {
	Type   string        `json:"type"`
	Pooled *PooledConfig `json:"pooled,omitempty"`
	Direct *DirectConfig `json:"direct,omitempty"`
}

type DirectConfig struct {
	MaxStepCount int `json:"maxStepCount"`
}

type PooledConfig struct {
	NumWorkers    int `json:"numWorkers"`
	WorkQueueSize int `json:"workQueueSize"`
	MaxStepCount  int `json:"maxStepCount"`
}

type TriggerConfig struct {
	Name      string            `json:"name"`
	Settings  map[string]string `json:"settings"`
	Endpoints []*EndpointConfig `json:"endpoints"`
}

type EndpointConfig struct {
	ProcessURI string `json:"processURI"`
	ConfigData string `json:"configData"` // if string, the trigger can unmarshall its own config
}

type ServiceConfig struct {
	Name     string            `json:"name"`
	Enabled  bool              `json:"enabled"`
	Settings map[string]string `json:"settings,omitempty"`
}

// DefaultConfig returns the default engine configuration
func DefaultEngineConfig() *EngineConfig {

	var ec EngineConfig

	ec.LogLevel = "INFO"
	ec.RunnerConfig =  &RunnerConfig{Type:"pooled", Pooled:&PooledConfig{NumWorkers:5, WorkQueueSize:50, MaxStepCount:32000}}
	ec.Triggers = make([]*TriggerConfig,0)
	ec.Services = make([]*ServiceConfig,0)

	ec.Services = append(ec.Services, &ServiceConfig{Name:"stateRecorder", Enabled: true, Settings: map[string]string{"host": ""}})
	ec.Services = append(ec.Services, &ServiceConfig{Name:"processProvider", Enabled: true})
	ec.Services = append(ec.Services, &ServiceConfig{Name: "engineTester", Enabled: true, Settings: map[string]string{"port": "8080"}})

	return &ec
}
