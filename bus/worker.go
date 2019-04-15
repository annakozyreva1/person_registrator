package bus

import (
	"context"
	"time"
)

const (
	MaxPublishTries       = 3
	TryTimeout            = time.Second * 10
	IdleConnectionTimeout = time.Minute
)

func worker(ctx context.Context, url string, tasks chan task, limit chan struct{}) {
	logger.Debug("started bus worker")
	pub := newPublisher(url)
	defer func() {
		pub.Close()
		<-limit
		logger.Debug("closed bus worker")
	}()
	var err error
	for {
		select {
		case task := <-tasks:
			{
				for try := 1; try <= MaxPublishTries; try++ {
					err = pub.Publish(task.Queue, task.ContentType, task.Body)
					if err == nil {
						break
					}
					time.Sleep(TryTimeout)
				}
				if err != nil {
					logger.Errorf("failed to publish in %s: %s", task.Queue, err.Error())
					task.Failure()
				} else {
					task.Success()
				}
			}
		case <-time.After(IdleConnectionTimeout):
			{
				return
			}
		case <-ctx.Done():
			{
				return
			}
		}
	}
}
