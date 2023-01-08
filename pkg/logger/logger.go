package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	logEnvDebug    = false
	logger         = log.New(os.Stdout, "", log.LstdFlags)
	extraSkipCount = 0
)

func writeLog(level string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	_, fname, line, _ := runtime.Caller(2 + extraSkipCount)

	logger.Printf("%s:%d [%s] %s\n", fname, line, level, msg)
}

func writeCustomLog(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	logger.Printf("%s\n", msg)
}

// InitLogger initialize variables for logger
func InitLogger(debugMode bool, fileName string) error {
	logEnvDebug = debugMode
	if fileName != "" {
		file, err := os.Create(fileName)
		if err != nil {
			return err
		}

		//do not call file.Close() because logger write log through file.Writer
		logger.SetOutput(file)
	}
	return nil
}

// Debug method outputs log as DEBUG Level
func Debug(format string, a ...interface{}) {
	if logEnvDebug {
		writeLog("DEBUG", format, a...)
	}
}

// Info method outputs log as INFO Level
func Info(format string, a ...interface{}) {
	writeLog("INFO", format, a...)
}

// Error method outputs log as ERROR Level
func Error(format string, a ...interface{}) {
	writeLog("ERROR", format, a...)
}

// ErrorCustom outputs custom format log as ERROR Level
func ErrorCustom(format string, a ...interface{}) {
	writeCustomLog(format, a...)
}

func SetExtraSkipCount(n int) {
	extraSkipCount = n
}
