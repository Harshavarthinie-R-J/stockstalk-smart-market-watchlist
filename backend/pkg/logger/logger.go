package logger

import (
	"log"
	"os"
)

// App logger.
func New() *log.Logger {
	return log.New(os.Stdout, "[STOCKSTALK] ", log.LstdFlags|log.LUTC)
}
