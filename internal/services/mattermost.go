package services

import (
	mmSDK "github.com/mattermost/mattermost-server/v6/model"
)

type MattermostClient struct {
	bot       *mmSDK.Client4
	channelID string
}

func NewMattermostClient(url, botToken, channelID string) (*MattermostClient, error) {
	mmBot := mmSDK.NewAPIv4Client(url)
	mmBot.SetOAuthToken(botToken)

	return &MattermostClient{bot: mmBot, channelID: channelID}, nil
}

func (mm *MattermostClient) Notify(message string) error {
	post := &mmSDK.Post{
		ChannelId: mm.channelID,
		Message:   message,
	}
	_, _, err := mm.bot.CreatePost(post)
	if err != nil {
		return err
	}

	return nil
}
