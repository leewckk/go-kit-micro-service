package consul

import "github.com/sirupsen/logrus"

type ConsulLogger struct {
}

func NewConsulLogger() *ConsulLogger {
	return &ConsulLogger{}
}

func (this *ConsulLogger) Log(keyvals ...interface{}) error {
	logrus.Info(keyvals...)
	return nil
}
