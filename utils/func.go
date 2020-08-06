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
		logrus.Error(ErrLogCode{LogType: logType, Message: message, Code: code, Err: err})
	}
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
