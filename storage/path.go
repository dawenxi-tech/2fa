package storage

import (
	"log/slog"
	"os"
)

func init() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func codePath() string {
	return "codes.json"
}

func configurePath() string {
	return "configure.json"
}
