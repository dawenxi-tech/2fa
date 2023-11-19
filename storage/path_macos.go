//go:build darwin && !ios && !nometal && RELEASE
// +build darwin,!ios,!nometal,RELEASE

package storage

import (
	"log/slog"
	"os"
	"path/filepath"
)

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	codePath = filepath.Join(configDir, bundleId, "codes.json")
	err = os.Mkdir(filepath.Dir(codePath), 0755)
	if err != nil {
		slog.Error("error to make config dir", slog.Any("err", err))
	}

	configurePath = filepath.Join(configDir, bundleId, "configure.json")
}
