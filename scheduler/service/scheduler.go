package service

import (
	"github.com/go-co-op/gocron"
	"github.com/r0bertson/msg-scheduler/common/db"
	"time"
)

type Service struct {
	DB        *db.Client
	scheduler *gocron.Scheduler
}

func sendPendingMessages(db *db.Client) {

}

func Init(db *db.Client) *Service {
	s := gocron.NewScheduler(time.UTC)
	return &Service{
		DB:        db,
		scheduler: s,
	}

}

func (s *Service) RunCronJobs() {
	s.scheduler.Every(1).Minutes().Do(func() {
		sendPendingMessages(s.DB)
	})
	s.scheduler.StartBlocking()
}
