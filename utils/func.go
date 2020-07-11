package utils

import (
	"github.com/sirupsen/logrus"
)

type ErrLogCode struct {
	LogType string
	Message string
	Code    int
	Err     error
}

func DeferRecoverLog(logType string, message string, code int, err error) {
	if panicErr := recover(); panicErr != nil {
		logrus.Warn(ErrLogCode{LogType: logType, Message: message, Code: code, Err: err})
	}
}
