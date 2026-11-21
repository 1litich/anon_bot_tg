package main

import (
	"context"
	"fmt"
	"html"

	"anon_bottg/database"

	"github.com/mymmrac/telego"
)

func video(bot *telego.Bot, update telego.Update) {
	if update.Message.Video != nil {
		userID := update.Message.From.ID
		videoID := update.Message.Video.FileID
		textVideo := update.Message.Caption

		safeCaption := "<i>Без подписи</i>"
		if textVideo != "" {
			safeCaption = html.EscapeString(textVideo)
		}

		if update.Message.ReplyToMessage != nil {
			messageID := update.Message.ReplyToMessage.MessageID

			senderID, err := database.GetSender(messageID)
			if err != nil {
				fmt.Println("Ошибка! ID пользователя не найден!", err)
				return
			}

			textFromAdmin := fmt.Sprintf(
				"💬 <b>Пришёл ответ:</b>\n\n"+
					"<blockquote>%s</blockquote>",
				safeCaption,
			)

			_, _ = bot.SendVideo(context.Background(), &telego.SendVideoParams{
				ChatID:    telego.ChatID{ID: senderID},
				Video:     telego.InputFile{FileID: videoID},
				Caption:   textFromAdmin,
				ParseMode: telego.ModeHTML,
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

			videoForAdmin := fmt.Sprintf(
				"🎬 <b>Новое анонимное видео!</b>\n\n"+
					"<blockquote>%s</blockquote>\n\n"+
					"<i>↩️ Чтобы ответить, свайпни это сообщение</i>",
				safeCaption,
			)

			sentMsg, err := bot.SendVideo(context.Background(), &telego.SendVideoParams{
				ChatID:    telego.ChatID{ID: targetID},
				Video:     telego.InputFile{FileID: videoID},
				Caption:   videoForAdmin,
				ParseMode: telego.ModeHTML,
			})
			if err != nil {
				fmt.Println("Ошибка отправки сообщений!", err)
				return
			}

			err = database.SaveSenderID(sentMsg.MessageID, userID)
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
