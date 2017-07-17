package device

import (
	"text/template"
	"strconv"
)

type SettingsConfig struct {
	DeviceName       string
	Id               string
	Settings         map[string]string
}

func (s *SettingsConfig) GetSetting(key string) string {
	return s.Settings[key]
}

type WithSettings interface {
	GetSetting(key string) string
}

var DeviceFuncMap = template.FuncMap{

	"val": func(name string, value interface{}) map[string]interface{} {

		val := make(map[string]interface{}, 1)
		val[name] = value

		return val
	},

	"setting": func(ctx interface{}, key string) string {
		config := ctx.(WithSettings)
		return config.GetSetting(key)
	},

	"settingb": func(ctx interface{}, key string) bool {
		config := ctx.(WithSettings)
		val,_ := strconv.ParseBool(config.GetSetting(key))
		return val
	},

	"debug": func(ctx interface{}) bool {
		config := ctx.(WithSettings)
		debug := config.GetSetting("device:debug")
		return debug == "true"
	},
}
