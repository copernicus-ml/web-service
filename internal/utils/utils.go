package utils

import "log/slog"

func LoggerGoogleFormate(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.LevelKey:
		a.Key = "severity"
		level, ok := a.Value.Any().(slog.Level)
		if !ok {
			return a
		}
		a.Value = getGCPLevelValue(level)
	case slog.MessageKey:
		a.Key = "message"
	}
	return a
}

func getGCPLevelValue(level slog.Level) slog.Value {
	var value slog.Value
	switch {
	case level >= slog.LevelError:
		value = slog.StringValue("ERROR")
	case level >= slog.LevelWarn:
		value = slog.StringValue("WARNING")
	case level >= slog.LevelInfo:
		value = slog.StringValue("INFO")
	default:
		value = slog.StringValue("DEBUG")
	}
	return value
}
