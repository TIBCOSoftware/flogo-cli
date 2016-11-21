package device

import (
	"text/template"
	"strings"
	"strconv"
)


type SettingsConfig struct {
	Settings         map[string]string
	EndpointSettings []map[string]string
}

var DeviceFuncMap = template.FuncMap {

	"isDigital": func(pinName string) string {
		upPinName := strings.ToUpper(pinName)
		isDigital := strings.HasPrefix(upPinName, "D")
		return strconv.FormatBool(isDigital)
	},

	"getPinNum": func(pinName string) string {
		return pinName[2:];
	},

	"setting": func(ctx interface{}, key string) string {
		config := ctx.(*SettingsConfig)
		return config.Settings[key]
	},

	"epsetting": func(ctx interface{}, id int, key string) string {
		config := ctx.(*SettingsConfig)
		return config.EndpointSettings[id][key]
	},

	"debug": func(ctx interface{}) bool {
		config := ctx.(*SettingsConfig)
		debug := config.Settings["device:debug"]
		return debug == "true"
	},
}
