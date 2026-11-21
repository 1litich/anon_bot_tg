package main

import (
	"context"
	"fmt"

	"anon_bottg/database"

	"github.com/mymmrac/telego"
)

func sticker(update telego.Update, bot *telego.Bot) {
	if update.Message.Sticker != nil {
		userID := update.Message.From.ID
		stickerID := update.Message.Sticker.FileID

		if update.Message.ReplyToMessage != nil {
			messageID := update.Message.ReplyToMessage.MessageID

			senderID, err := database.GetSender(messageID)
			if err != nil {
				fmt.Println("Ошибка! ID получателя не найден!", err)
				return
			}

			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: senderID},
				Text:      "💬 <b>Пришёл ответ стикером:</b>",
				ParseMode: telego.ModeHTML,
			})

			_, _ = bot.SendSticker(context.Background(), &telego.SendStickerParams{
				ChatID:  telego.ChatID{ID: senderID},
				Sticker: telego.InputFile{FileID: stickerID},
			})

			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: userID},
				Text:      "✅ <b>Ответ отправлен!</b>",
				ParseMode: telego.ModeHTML,
			})
		} else {
			targetID, err := database.GetTarget(userID)
			if err != nil {
				_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
					ChatID:    telego.ChatID{ID: userID},
					Text:      "⚠️ <b>Кому отправить?</b>\nСначала перейди по чьей-либо персональной ссылке или свайпни входящее сообщение для ответа.",
					ParseMode: telego.ModeHTML,
				})
				return
			}

			sentMsg, err := bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: targetID},
				Text:      "🎭 <b>Новый анонимный стикер!</b>\n<i>↩️ Чтобы ответить, свайпни это сообщение или сам стикер</i>",
				ParseMode: telego.ModeHTML,
			})
			if err != nil {
				fmt.Println("Ошибка отправки сообщения(2)", err)
				return
			}

			err = database.SaveSenderID(sentMsg.MessageID, userID)
			if err != nil {
				fmt.Println("Ошибка записи в бд senderid!", err)
				return
			}

			sentSticker, err := bot.SendSticker(context.Background(), &telego.SendStickerParams{
				ChatID:  telego.ChatID{ID: targetID},
				Sticker: telego.InputFile{FileID: stickerID},
			})
			if err != nil {
				fmt.Println("Ошибка отправки сообщений(3)!", err)
				return
			}

			err = database.SaveSenderID(sentSticker.MessageID, userID)
			if err != nil {
				fmt.Println("Ошибка записи в бд senderid!", err)
				return
			}

			_, _ = bot.SendMessage(context.Background(), &telego.SendMessageParams{
				ChatID:    telego.ChatID{ID: userID},
				Text:      "📨 <b>Доставлено!</b> Ожидай ответа в этом чате.",
				ParseMode: telego.ModeHTML,
			})
		}
	}
}
