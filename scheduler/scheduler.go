package scheduler

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/zelenin/go-tdlib/client"
)

func scheduleMessage(
	c *client.Client,
	chatId int32,
	text string,
	sendDate time.Time,
) (*client.Message, error) {
	msg, err := c.SendMessage(
		&client.SendMessageRequest{
			ChatId: int64(chatId),
			Options: &client.MessageSendOptions{
				SchedulingState: &client.MessageSchedulingStateSendAtDate{
					SendDate: int32(sendDate.Unix()),
				},
			},
			InputMessageContent: &client.InputMessageText{
				Text: &client.FormattedText{Text: text},
			},
		},
	)
	if err != nil {
		return msg, err
	}

	return msg, nil
}

func scheduleMessagesWithCooldown(
	c *client.Client,
	chatId int64,
	text string,
	nTimes uint,
	startTime time.Time,
	cooldown time.Duration,
) ([]*client.Message, error) {
	var messages []*client.Message
	cooldown += time.Minute
	scheduleTime := startTime

	for i := 0; i < int(nTimes); i++ {
		msg, err := c.SendMessage(
			&client.SendMessageRequest{
				ChatId: chatId,
				Options: &client.MessageSendOptions{
					SchedulingState: &client.MessageSchedulingStateSendAtDate{
						SendDate: int32(scheduleTime.Unix()),
					},
				},
				InputMessageContent: &client.InputMessageText{
					Text: &client.FormattedText{Text: text},
				},
			},
		)
		if err != nil {
			return messages, err
		}

		messages = append(messages, msg)
		log.Println(msg)

		scheduleTime = scheduleTime.Add(cooldown)
		time.Sleep(time.Second)
	}

	return messages, nil
}

func ScheduleFeedings(
	c *client.Client,
	nTimes uint,
	cooldownHours uint,
	startTime time.Time,
) ([]*client.Message, error) {
	var messages []*client.Message

	feedMessageText := "покормить жабу"
	feedCooldown := time.Duration(cooldownHours) * time.Hour

	chatId, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	if err != nil {
		return messages, err
	}

	messages, err = scheduleMessagesWithCooldown(
		c,
		int64(chatId),
		feedMessageText,
		nTimes,
		startTime,
		feedCooldown,
	)

	return messages, nil
}

func ScheduleWork(
	c *client.Client,
	nTimes uint,
	startTime time.Time,
) ([]*client.Message, error) {
	var messages []*client.Message

	startWorkingMessageText := "отправить жабу на работу"
	finishWorkingMessageText := "завершить работу"
	workDuration := 2 * time.Hour
	workCooldown := 6 * time.Hour

	chatId, err := strconv.Atoi(os.Getenv("CHAT_ID"))
	if err != nil {
		return messages, err
	}

	startWorkingMessages, err := scheduleMessagesWithCooldown(
		c,
		int64(chatId),
		startWorkingMessageText,
		nTimes,
		startTime,
		workDuration + workCooldown,
	)
	if err != nil {
		return messages, err
	}

	finishWorkingMessages, err := scheduleMessagesWithCooldown(
		c,
		int64(chatId),
		finishWorkingMessageText,
		nTimes,
		startTime.Add(workDuration),
		workDuration + workCooldown,
	)
	if err != nil {
		return messages, err
	}

	messages = append(startWorkingMessages, finishWorkingMessages[:]...)

	return messages, nil
}
