//go:build !js

package persist

import (
	"os"
	"path/filepath"
)

const appName = "blocks"

// getAppDataPath returns the path to the app's data directory.
func getAppDataPath() (string, error) {
	configDir, err := os.UserConfigDir() // Cross-platform config directory
	if err != nil {
		return "", err
	}

	appDataPath := filepath.Join(configDir, appName)

	// Create the directory if it doesn't exist
	err = os.MkdirAll(appDataPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return appDataPath, nil
}

func Store(key string, value string) error {
	path, err := getAppDataPath()
	if err != nil {
		return err
	}
	filename := filepath.Join(path, key)
	return os.WriteFile(filename, []byte(value), os.ModePerm)
}

func Load(key string) (string, error) {
	path, err := getAppDataPath()
	if err != nil {
		return "", err
	}
	filename := filepath.Join(path, key)
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
