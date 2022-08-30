package service

import (
	"github.com/go-co-op/gocron"
	"github.com/r0bertson/msg-scheduler/common/db"
	"github.com/r0bertson/msg-scheduler/common/messaging"
	"github.com/r0bertson/msg-scheduler/common/models"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

type Service struct {
	DB         *db.Client
	scheduler  *gocron.Scheduler
	msgService messaging.MsgService
}

// this is very naive and not optimal at all
func notIn(list []uint, in []uint) []uint {
	listDict := map[uint]bool{}

	for _, e := range list {
		listDict[e] = true
	}

	for _, i := range in {
		listDict[i] = false
	}
	notIn := []uint{}
	for key, value := range listDict {
		if value {
			notIn = append(notIn, key)
		}
	}
	return notIn
}

func (s *Service) sendMessages() {
	log.Info().Msg("Started sending messages.")
	messages, err := s.DB.Messages()
	if err != nil {
		return //skipping error because it'll be retried later
	}
	messageIds := []uint{}
	messageDict := make(map[uint]models.Message)

	for _, msg := range *messages {
		messageIds = append(messageIds, msg.ID)
		messageDict[msg.ID] = msg
	}

	pendingUserStats, err := s.DB.AllPendingUsersStats()
	if err != nil {
		return
	}

	for _, stats := range *pendingUserStats {
		pendingMessages := notIn(messageIds, stats.Stats.SentMessages)
		if len(pendingMessages) == 0 {
			//in case the flag is still true
			stats.User.ShouldSendMessages = false
			s.DB.UpdateUserStatus(&stats.User)
			continue
		}
		//this is not truly random, but since the data is different every time, it is not a problem
		nextMsgId := pendingMessages[rand.Intn(len(pendingMessages))]
		nextMsg := messageDict[nextMsgId]

		log.Info().Msgf("Sending message #%d to user %d", nextMsgId, stats.User.ID)
		if err := s.msgService.SendMessage(messaging.EmailPayload{
			To:      stats.Email,
			Subject: nextMsg.Subject,
			Content: nextMsg.Content,
		}); err != nil {
			//ignoring error because it'll be retried one minute later
			continue
		}

		sentMessage := &models.SentMessage{
			UserID:    stats.User.ID,
			MessageID: nextMsgId,
		}

		if sentMessage, err = s.DB.CreateSentMessage(sentMessage); err != nil {
			//todo: retry mechanism would be good here
			log.Info().Msg("Couldn't insert sent message into DB.")
			continue
		}

		if len(pendingMessages) == 1 {
			stats.User.ShouldSendMessages = false
			s.DB.UpdateUserStatus(&stats.User)
			log.Info().Msgf("All messages were sent to user %d", stats.User.ID)
			//no need to handle errors here, because we check for their side effect in the beginning of this loop
		}
	}
	log.Info().Msg("Finished sending messages.")

}

func Init(db *db.Client, msgService messaging.MsgService) *Service {
	s := gocron.NewScheduler(time.UTC)
	return &Service{
		DB:         db,
		scheduler:  s,
		msgService: msgService,
	}

}

func (s *Service) RunCronJobs() {
	s.scheduler.Every(1).Minutes().Do(func() {
		s.sendMessages()
	})
	s.scheduler.StartBlocking()
}
