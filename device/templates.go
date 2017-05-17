package device

import (
	"text/template"
)

type SettingsConfig struct {
	DeviceName       string
	Settings         map[string]string
}

var DeviceFuncMap = template.FuncMap{

	"setting": func(ctx interface{}, key string) string {
		config := ctx.(*SettingsConfig)
		return config.Settings[key]
	},

	"debug": func(ctx interface{}) bool {
		config := ctx.(*SettingsConfig)
		debug := config.Settings["device:debug"]
		return debug == "true"
	},
}
