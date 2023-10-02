package telesearch

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	os "os"
	fmt "fmt"
)

func StartBot() {
	LoadEnv()

	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(telegramBotToken)

	if err != nil {
		panic(err)
    }

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		channel := make(chan SearchResult)

		if update.Message == nil {
            continue
        }

		go SearchGoogle(update.Message.Text, 0, channel)
		go func(update tgbotapi.Update, ch chan SearchResult) {
			chatId := update.Message.Chat.ID
			messageId := update.Message.MessageID

			text := ""
			results := (<- ch).Results

			for i := 0; i < len(results); i++ {
				text += fmt.Sprintf("%d.\t%s\n%s\n\n", i, results[i].Title, results[i].Url)
			}

			msg := tgbotapi.NewMessage(chatId, text)
			msg.BaseChat.ReplyToMessageID = messageId

			bot.Send(msg)
		}(update, channel)
	}
}
