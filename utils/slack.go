package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slack-go/slack"
)

const (
	InitialMessage     string = "Sucessfully configured and saved iptables rules."
	AlarmMessage       string = "Alarm, configuration changed outside of app."
	DisabledMessage    string = "Unix-defender is disabled"
	StartMessage       string = "Unix-defender is started"
	ReconfigureMessage string = "IpTables rules are reconfigured."
	GreenColor         string = "#36a64f"
	RedColor           string = "#FF0000"
)

func SendMessageToSlack(text, color string) {
	config, err := LoadConfigEnv(EnvFile)
	if err != nil {
		log.Fatal("Cannot load environment config:", err)
	}
	if config.SlackEnable {
		hostname, err := os.Hostname()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		token := config.SlackAuthToken
		channelID := config.SlackChannelId
		client := slack.New(token, slack.OptionDebug(false))
		attachment := slack.Attachment{
			Title: "Unix-defender Notification",
			Text:  text,
			Color: color,
			Fields: []slack.AttachmentField{
				{
					Title: "Host",
					Value: hostname,
				},
				{
					Title: "Event time",
					Value: time.Now().Format("2006.01.02 15:04:05"),
				},
			},
		}
		_, timestamp, err := client.PostMessage(
			channelID,
			slack.MsgOptionAttachments(attachment),
		)
		_ = timestamp
		if err != nil {
			panic(err)
		}
		if config.SlackSendFiles && color == GreenColor {
			channelArr := []string{config.SlackChannelId}
			fileArr := []string{config.RulesBackupV4, config.RulesBackupV6}
			for i := 0; i < len(fileArr); i++ {
				params := slack.FileUploadParameters{
					Channels: channelArr,
					File:     fileArr[i],
				}
				file, err := client.UploadFile(params)
				if err != nil {
					panic(err)
				}
				fmt.Println("Files is uploaded:", file.Name)
			}
		}
	}

}
