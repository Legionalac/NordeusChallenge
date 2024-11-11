package libs

import (
	"os"

	"github.com/rs/zerolog"
)
var customLogger zerolog.Logger

func init(){
    output := zerolog.ConsoleWriter{Out: os.Stdout}
    customLogger = zerolog.New(output).With().Timestamp().Logger()
}

func GetLogger() zerolog.Logger{
	return customLogger
}
