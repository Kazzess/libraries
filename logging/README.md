# logging

## how to use

```golang

package main

import (
	"git.adapticode.com/libraries/golang/logging"
)

func main() {
	// logger from context
	logger := logging.L(ctx)
	// logger without access to context (level debug)
	logger = logging.Logger()

	// create logger with level
	logger = logging.LoggerWLevel("info")

	// add field to logger and put to context
	logger = logger.With(slog.String("key", "value"))
	ctx = logging.ContextWithLogger(ctx, l)
	
	// create logger with fields instantly
	logger = logging.WithAttrs(ctx, slog.String("key", "value"))
}

```