package logger

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/EvGesh4And/anti-bruteforce/config"
)

// NewSlogLogger creates a new slog.Logger instance with the given configuration.
func NewSlogLogger(cfg config.LoggerConfig) (*slog.Logger, io.Closer, error) {
	var out io.WriteCloser
	switch cfg.Mod {
	case "console", "":
		out = os.Stdout
	case "file":
		filePath := cfg.Path
		if filePath == "" {
			filePath = "calendar.log"
		}
		// #nosec G304
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			log.Printf("error opening log file %s: %s", filePath, err)
			return nil, nil, err
		}
		out = f
	default:
		log.Printf("unknown logger mode: %s, using console", cfg.Mod)
		out = os.Stdout
	}

	l := New(cfg.Level, out, cfg.JSON)
	var closer io.Closer
	if c, ok := out.(io.Closer); ok && c != os.Stdout {
		closer = c
	}
	return l, closer, nil
}
