package log

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	Logger *log.Logger // Глобальный объект логгера
	once   sync.Once
)

func init() {
	once.Do(func() {
		Logger = log.StandardLogger()
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: false,
			TimestampFormat:  "2006-01-02 12:00:00",
		})
		Logger.SetLevel(log.TraceLevel)
	})
}
