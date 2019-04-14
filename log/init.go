package log

import (
	"sync"
	log "github.com/sirupsen/logrus"
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
		level, err := log.ParseLevel("debug")
		if err != nil {
			level = log.InfoLevel
		}
		Logger.SetLevel(level)
	})
}
