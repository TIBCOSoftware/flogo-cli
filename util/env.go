package fgutil

import "strings"

func ReplaceEnvValue(env []string, envKey string, newValue string) []string {
	for key, entry := range env {

		if strings.HasPrefix(entry, envKey + "=") {
			env[key] = envKey + "=" + newValue
			break
		}
	}

	return env
}
