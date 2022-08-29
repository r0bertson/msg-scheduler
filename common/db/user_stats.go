package db

import (
	"github.com/r0bertson/msg-scheduler/common/models"
)

func (c *Client) AllPendingUsersStats() (*[]models.UserStats, error) {
	//fetch pending users
	users, err := c.PendingUsers()
	if err != nil {
		return nil, err
	}
	if len(*users) == 0 {
		return &[]models.UserStats{}, nil
	}
	var userIds []uint
	for _, user := range *users {
		userIds = append(userIds, user.ID)
	}
	//fetch messages sent of each user
	msgsSent, err := c.MessagesSentFor(userIds)
	if err != nil {
		return nil, err
	}
	msgsByUserId := make(map[uint][]uint)
	for _, msgSent := range *msgsSent {
		msgsByUserId[msgSent.UserID] = append(msgsByUserId[msgSent.UserID], []uint{msgSent.MessageID}...)
	}
	var stats []models.UserStats
	//join users and statistics together
	for _, user := range *users {
		stats = append(stats, models.UserStats{
			User:  user,
			Stats: models.Stats{SentMessages: msgsByUserId[user.ID]},
		})
	}
	return &stats, nil
}
