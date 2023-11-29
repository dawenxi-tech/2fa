package storage

import (
	"encoding/json"
	"log/slog"
	"os"
)

type Configure struct {
	ExitWhenWindowClose bool `json:"exitWhenWindowClose"`
	ShowTray            bool `json:"showTray"`
	WindowMode          bool `json:"windowMode"`
}

func LoadConfigure() Configure {
	var conf Configure
	data, err := os.ReadFile(configurePath)
	if err == os.ErrNotExist {
		return conf
	}
	if err != nil {
		slog.Error("error to read configure file", slog.Any("err", err))
		return conf
	}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to unmarshal configure")
		return conf
	}
	return conf
}

func SaveConfigure(conf Configure) {
	data, err := json.Marshal(conf)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to marshal configure")
		return
	}
	fp, err := os.Create(configurePath)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to create configure path")
		return
	}
	_, err = fp.Write(data)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error to write configure file")
		return
	}
}
