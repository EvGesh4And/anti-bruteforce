// Package logger предоставляет средства для логирования,
// включая уровни логов, форматирование и вывод в разные источники.
package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

var levelMap = map[string]slog.Level{
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
}

// New creates a logger with the provided level and writer.
// New creates a logger with the provided level, writer and format (JSON or human-readable).
func New(level string, out io.Writer, asJSON bool) *slog.Logger {
	var levLog slog.Level
	if lvl, ok := levelMap[level]; ok {
		levLog = lvl
	} else {
		levLog = slog.LevelDebug
	}

	var handler slog.Handler

	if asJSON {
		// JSON-логирование для парсинга в Elastic, Loki и т.п.
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: levLog,
		})
	} else {
		// Человекочитаемый вывод с цветами (если терминал)
		useColor := isStdout(out)

		handler = tint.NewHandler(out, &tint.Options{
			Level:      levLog,
			TimeFormat: time.Kitchen,
			NoColor:    !useColor,
		})
	}

	// middleware для добавления полей или обработки
	handler = NewHandlerMiddleware(handler)

	return slog.New(handler)
}

func isStdout(out io.Writer) bool {
	f, ok := out.(*os.File)
	if !ok {
		return false
	}
	return f.Fd() == os.Stdout.Fd()
}
