package logger

import "github.com/rs/zerolog"


// SetLogLevel parses and sets the global logging level.
func SetLogLevel(levelStr string) {
	level := zerolog.NoLevel

	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.NoLevel
	}

	zerolog.SetGlobalLevel(level)

}